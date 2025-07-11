package main

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"consuldiff/consulkv"
	"consuldiff/gitutil"
	"consuldiff/kvtypes"
	"consuldiff/state"
	"consuldiff/storage"

	"github.com/hashicorp/consul/api"
)

func main() {

	// Initialize the state
	s := state.InitState()

	// Setup git repository configuration
	if s.GitConfig.Enabled {
		gitutil.SetupGitRepo(*s)
	}

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
		if err != nil {
			log.Fatalf("Error writing initial KV state to file: %v", err)
		}
	}

	log.Printf("Using storage file: %s", s.GitConfig.Filename)
	log.Println("Starting KV polling for diff...")

	for {
		log.Println("Fetching current KV state from Consul...")

		s.Current, err = consulkv.FetchKV(client, s.KeyPrefix)
		if err != nil {
			log.Fatalf("Error fetching KV: %v", err)
		}

		filepath := s.GitConfig.RepoPath + "/" + s.GitConfig.Filename
		s.Previous, err = storage.ReadMapFromFile(filepath)
		if err != nil {
			log.Fatalf("Error reading previous KV state: %v", err)
			s.Previous = []kvtypes.KVExportEntry{}

		}

		if s.Previous != nil {
			log.Println("Comparing current KV state with previous state from file: ", filepath)
			consulkv.LogKVDiff(s.Previous, s.Current, filepath)
		}

		log.Println("Writing current base64 KV state to file...")
		kvB64, err := consulkv.FetchKVBase64(client, s.KeyPrefix)
		if err != nil {
			log.Fatalf("Error fetching KV Base64: %v", err)
		}
		storage.WriteMapToFile(filepath+".b64", kvB64)

		if s.GitConfig.Enabled {
			log.Println("Git is enabled, committing changes if any...")
			err := gitutil.GitCommitAndPush(*s)
			if err != nil {
				log.Fatalf("Error committing and pushing changes to Git: %v", err)
			}
		}

		time.Sleep(s.PollInterval)
	}
}
