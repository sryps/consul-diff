package gitutil

import (
	"log"
	"os"
	"time"

	"consuldiff/state"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	// Add the b64 file
	_, err = w.Add(s.GitConfig.Filename + ".b64")
	if err != nil {
		return err
	}

	// Check if there are changes to commit
	status, err := w.Status()
	if err != nil {
		return err
	}
	if status.IsClean() {
		log.Printf("No changes to commit in %s", s.GitConfig.RepoPath)
		return nil
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
		log.Printf("Error pushing changes to %s: %v", s.GitConfig.RemoteURL, err)
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
			_, err := git.PlainClone(s.GitConfig.RepoPath, false, &git.CloneOptions{
				URL: s.GitConfig.RemoteURL,
				Auth: &http.BasicAuth{
					Username: "git",             // GitHub accepts anything
					Password: s.GitConfig.Token, // Personal Access Token
				},
				ReferenceName: plumbing.ReferenceName("refs/heads/" + s.GitConfig.Branch),
				SingleBranch:  true,
				Depth:         1,
			})
			if err != nil {
				log.Fatalf("Error cloning repo: %v", err)
				return nil, err
			}
		}
	}
	log.Printf("Opening existing repo at %s", s.GitConfig.RepoPath)
	return git.PlainOpen(s.GitConfig.RepoPath)
}
