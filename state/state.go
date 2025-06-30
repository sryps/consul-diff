package state

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type State struct {
	Previous      map[string]string
	Current       map[string]string
	GitConfig     GitConfig
	PollInterval  time.Duration
	KeyPrefix     string
	TLSSkipVerify bool
}

type GitConfig struct {
	Enabled     bool
	RepoPath    string
	RemoteURL   string
	Filename    string
	Message     string
	AuthorName  string
	AuthorEmail string
	Token       string
}

func InitState() *State {

	// Ensure the TLS_SKIP_VERIFY environment variable is set
	var TLSSkipVerify bool
	var err error
	if os.Getenv("TLS_SKIP_VERIFY") == "" {
		log.Println("Environment variable TLS_SKIP_VERIFY is not set, defaulting to false")
		TLSSkipVerify = false
	} else {
		TLSSkipVerify, err = strconv.ParseBool(os.Getenv("TLS_SKIP_VERIFY"))
		if err != nil {
			panic("Environment variable TLS_SKIP_VERIFY is set but not a valid boolean, please set it to 'true' or 'false'.")
		}
	}

	// Ensure the GIT_ENABLED environment variable is set
	gitEnabled := os.Getenv("GIT_ENABLED")
	if gitEnabled == "" {
		panic("Environment variable GIT_ENABLED is not set, please set it to 'true' or 'false'.")
	}
	gitEnabledBool, err := strconv.ParseBool(gitEnabled)
	if err != nil {
		panic("Environment variable GIT_ENABLED is set but not a valid boolean, please set it to 'true' or 'false'.")
	}

	// Ensure the GIT_REPO_PATH environment variable is set
	repoPath := os.Getenv("GIT_REPO_PATH")
	if repoPath == "" {
		panic("Environment variable GIT_REPO_PATH is not set, please set it to the path of your Git repository.")
	}
	// Ensure the GIT_REMOTE_URL environment variable is set
	if os.Getenv("GIT_REMOTE_URL") == "" {
		panic("Environment variable GIT_REMOTE_URL is not set, please set it to the remote URL of your Git repository.")
	}
	//if GIT_REPO_PATH is not a valid directory
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		panic("GIT_REPO_PATH does not point to a valid directory, please check the path.")
	}
	// if GIT_REPO_PATH ends in a slash, remove it
	repoPath = strings.TrimSuffix(repoPath, "/")

	// Ensure the GIT_AUTHOR_NAME environment variable is set
	if os.Getenv("GIT_AUTHOR_NAME") == "" {
		log.Println("Environment variable GIT_AUTHOR_NAME is not set, using default 'git'")
		os.Setenv("GIT_AUTHOR_NAME", "git")
	}
	// Ensure the GIT_AUTHOR_EMAIL environment variable is set
	if os.Getenv("GIT_AUTHOR_EMAIL") == "" {
		log.Println("Environment variable GIT_AUTHOR_EMAIL is not set, using default 'consuldiff@consuldiff.com'")
		os.Setenv("GIT_AUTHOR_EMAIL", "consuldiff@consuldiff.com")
	}
	// Ensure the GIT_TOKEN environment variable is set
	if os.Getenv("GIT_TOKEN") == "" {
		panic("Environment variable GIT_TOKEN is not set, please set it to your Git token.")
	}
	// Ensure the POLL_INTERVAL_MINUTES environment variable is set
	if os.Getenv("POLL_INTERVAL_MINUTES") == "" {
		panic("Environment variable POLL_INTERVAL_MINUTES is not set, please set it to the desired polling interval in minutes.")
	}
	// Ensure the CONSUL_HTTP_ADDR environment variable is set
	if os.Getenv("CONSUL_HTTP_ADDR") == "" {
		panic("Environment variable CONSUL_HTTP_ADDR is not set, please set it to the address of your Consul server.")
	}
	// Ensure the GIT_COMMIT_MESSAGE environment variable is set
	if os.Getenv("GIT_COMMIT_MESSAGE") == "" {
		log.Println("Environment variable GIT_COMMIT_MESSAGE is not set, using default 'Consul KV Diff Update'")
		os.Setenv("GIT_COMMIT_MESSAGE", "Consul KV Diff Update")
	}

	keyPrefix := os.Getenv("CONSUL_KV_PREFIX")
	if keyPrefix == "" {
		log.Println("Environment variable CONSUL_KV_PREFIX is not set, using default '*/'")
	}

	// get poll interval from environment variable
	pollIntervalEnv, err := strconv.Atoi(os.Getenv("POLL_INTERVAL_MINUTES"))
	if pollIntervalEnv == 0 || err != nil {
		panic("Environment variable POLL_INTERVAL_MINUTES is not set or invalid")
	}
	pollInterval := time.Duration(pollIntervalEnv) * time.Minute

	// Ensure the STORAGE_FILENAME environment variable is set
	if os.Getenv("STORAGE_FILENAME") == "" {
		log.Println("Environment variable STORAGE_FILENAME is not set, using default 'consul_kv_diff.json'")
		os.Setenv("STORAGE_FILENAME", "consul_kv_diff.json")
	}

	returnData := &State{
		Previous: make(map[string]string),
		Current:  make(map[string]string),
		GitConfig: GitConfig{
			Enabled:     gitEnabledBool,
			RepoPath:    repoPath,
			RemoteURL:   os.Getenv("GIT_REMOTE_URL"),
			Filename:    os.Getenv("STORAGE_FILENAME"),
			Message:     os.Getenv("GIT_COMMIT_MESSAGE"),
			AuthorName:  os.Getenv("GIT_AUTHOR_NAME"),
			AuthorEmail: os.Getenv("GIT_AUTHOR_EMAIL"),
			Token:       os.Getenv("GIT_TOKEN"),
		},
		PollInterval:  pollInterval,
		KeyPrefix:     keyPrefix,
		TLSSkipVerify: TLSSkipVerify,
	}

	LogRedactedState(*returnData)

	return returnData
}

func LogRedactedState(s State) {
	returnDataRedacted := State{
		Previous: make(map[string]string),
		Current:  make(map[string]string),
		GitConfig: GitConfig{
			RepoPath:    s.GitConfig.RepoPath,
			RemoteURL:   s.GitConfig.RemoteURL,
			Filename:    s.GitConfig.Filename,
			Message:     s.GitConfig.Message,
			AuthorName:  s.GitConfig.AuthorName,
			AuthorEmail: s.GitConfig.AuthorEmail,
			Token:       "[REDACTED]",
		},
		PollInterval:  s.PollInterval,
		KeyPrefix:     s.KeyPrefix,
		TLSSkipVerify: s.TLSSkipVerify,
	}
	returnDataBytes, err := json.MarshalIndent(returnDataRedacted, "", "  ")
	if err != nil {
		log.Printf("Error marshalling state: %v", err)
	}
	returnDataString := string(returnDataBytes)
	log.Printf("Initialized State: %s", returnDataString)
}
