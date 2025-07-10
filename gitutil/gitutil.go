package gitutil

import (
	"log"
	"os"
	"time"

	"consuldiff/state"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GitCommitAndPush stages, commits, and pushes a file to a Git repo.
func GitCommitAndPush(s state.State) error {
	// Open the repo
	repo, err := git.PlainOpen(s.GitConfig.RepoPath)
	if err != nil {
		return err
	}

	// Get the working tree
	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	// Add the file
	_, err = w.Add(s.GitConfig.Filename)
	if err != nil {
		return err
	}

	// Commit
	_, err = w.Commit(s.GitConfig.Message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  s.GitConfig.AuthorName,
			Email: s.GitConfig.AuthorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	// Push with token (GitHub, GitLab, etc.)
	err = repo.Push(&git.PushOptions{
		Auth: &http.BasicAuth{
			Username: s.GitConfig.AuthorName,
			Password: s.GitConfig.Token, // Use a PAT (Personal Access Token)
		},
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	log.Printf("Git changes pushed to %s", s.GitConfig.RemoteURL)
	return nil
}

func SetupGitRepo(s state.State) (*git.Repository, error) {
	log.Printf("Setting up Git repository at %s for %s", s.GitConfig.RepoPath, s.GitConfig.RemoteURL)

	if _, err := os.Stat(s.GitConfig.RepoPath + "/.git"); os.IsNotExist(err) {
		if s.GitConfig.RemoteURL != "" {
			log.Printf("Cloning repo from %s into %s", s.GitConfig.RemoteURL, s.GitConfig.RepoPath)
			return git.PlainClone(s.GitConfig.RepoPath, false, &git.CloneOptions{
				URL: s.GitConfig.RemoteURL,
				Auth: &http.BasicAuth{
					Username: "git",             // GitHub accepts anything
					Password: s.GitConfig.Token, // Personal Access Token
				},
			})
		} else {
			log.Printf("Initializing new git repo at %s", s.GitConfig.RepoPath)
			repo, err := git.PlainInit(s.GitConfig.RepoPath, false)
			if err != nil {
				return nil, err
			}
			_, err = repo.CreateRemote(&config.RemoteConfig{
				Name: "origin",
				URLs: []string{s.GitConfig.RemoteURL},
			})
			if err != nil {
				return nil, err
			}
			return repo, nil
		}
	}
	log.Printf("Opening existing repo at %s", s.GitConfig.RepoPath)
	return git.PlainOpen(s.GitConfig.RepoPath)
}
