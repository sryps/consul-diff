package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"consuldiff/consulkv"
	"consuldiff/gitutil"
	"consuldiff/state"
	"consuldiff/storage"

	"github.com/hashicorp/consul/api"
)

func main() {

	// Initialize the state
	s := state.InitState()

	// Setup git repository configuration
	gitutil.SetupGitRepo(*s)

	config := api.DefaultConfig()
	if s.TLSSkipVerify {
		// Custom HTTP client with TLS skip verify
		config.HttpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}

	log.Println("Starting Consul KV Diff Tool...")
	log.Printf("Polling Consul KV every %s", s.PollInterval)

	// Init Consul client
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	// Check if the storage file exists and read it
	fileCheck, err := storage.ReadMapFromFile(s.GitConfig.Filename)
	checkBytes, err := json.Marshal(fileCheck)
	if err != nil {
		// handle error
	}
	check := string(checkBytes)
	if check == "null" {
		log.Printf("No previous KV state found, creating new file at %s", s.GitConfig.Filename)
		firstRun, err := consulkv.FetchKV(client, s.KeyPrefix)
		if err != nil {
			log.Fatalf("Error fetching initial KV state: %v", err)
		}
		err = storage.WriteMapToFile(s.GitConfig.Filename, firstRun)
	}
	if err != nil {
		log.Fatalf("Error initializing storage file: %v", err)
	}

	log.Printf("Using storage file: %s", s.GitConfig.Filename)
	log.Println("Starting KV polling for diff...")

	for {
		s.Current, err = consulkv.FetchKV(client, s.KeyPrefix)
		if err != nil {
			log.Printf("Error fetching KV: %v", err)
		}

		filepath := s.GitConfig.RepoPath + "/" + s.GitConfig.Filename
		s.Previous, err = storage.ReadMapFromFile(filepath)
		if err != nil {
			log.Printf("Error reading previous KV state: %v", err)
			s.Previous = make(map[string]string) // Initialize if file doesn't exist
		}

		if s.Previous != nil {
			consulkv.DiffKV(s.Previous, s.Current, filepath)
		}

		if s.GitConfig.Enabled {
			gitutil.GitCommitAndPush(*s)
		}

		time.Sleep(s.PollInterval)
	}
}
