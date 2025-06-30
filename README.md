# Consul KV Diff Tool

A simple tool that logs the diff between changes in Consul KV store at a given interval.
It logs the changes and also commits them to a Git repository.

## Features

- Monitors Consul KV store for changes
- Logs changes to a file
- Commits changes to a Git repository
- Configurable interval for checking changes

## Required Envirnment Variables

- `CONSUL_HTTP_ADDR`: The address of the Consul HTTP API (e.g., `http://localhost:8500`).
- `GIT_ENABLED`: Set to `true` to enable Git operations. If set to `false`, the tool will only log changes without committing them.
- `GIT_REPO_PATH`: The path to the Git repository where changes will be committed.
- `GIT_COMMIT_MESSAGE`: The commit message to use when committing changes to the Git repository.
- `GIT_REMOTE_URL`: The remote URL of the Git repository.
- `GIT_AUTHOR_NAME`: The name of the author for the Git commits.
- `GIT_AUTHOR_EMAIL`: The email of the author for the Git commits.
- `GIT_TOKEN`: The Git PAT token for authentication with the remote repository. Requires CONTENT permissions for the repository.
- `POLL_INTERVAL_MINUTES`: The interval in minutes at which to check for changes in the Consul KV store.
- `STORAGE_FILENAME`: The file where the last known state of the KV store will be saved.
- `TLS_SKIP_VERIFY`: Set to `true` to skip TLS verification (not recommended for production).
- `KEY_PREFIX`: The prefix for the keys in the Consul KV store to monitor (optional). Defaults to all paths "\*/".

## Installation

1. Clone the repository:
   ```bash
   git clone ${GIT_REPO_URL}
   cd consul-diff
   go build
   ```
2. Set the required environment variables in your shell or `.env` file.
3. Run the tool:
   ```bash
   source .env
   ./consuldiff
   ```

## Important Notes

- Ensure that the Consul server is running and accessible at the specified `CONSUL_HTTP_ADDR`.
- Make sure the POLL_INTERVAL_MINUTES is set to a reasonable value to avoid excessive load on the Consul server. Especially if the Consul KV store is large.
- The tool uses the `GIT_ENABLED` environment variable to determine whether to commit changes to the Git repository. If set to `false`, it will only log changes without committing.
- The tool will create a file specified by `GIT_REPO_PATH`/`STORAGE_FILENAME` to store the last known state of the KV store. Ensure that the path is writable.
- The tool uses the `KEY_PREFIX` environment variable to filter keys in the Consul KV store. If not set, it will monitor all keys.
- The tool uses the `TLS_SKIP_VERIFY` environment variable to skip TLS verification. This is not recommended for production environments as it can expose you to security risks.
- The tool requires a Git repository with the necessary permissions to commit changes. Ensure that the `GIT_TOKEN` has the required permissions for the repository.
- The tool uses the `GIT_REMOTE_URL` to push changes to the remote repository. Ensure that the URL is correct and accessible.

> ⚠️ This project is in early testing. Feedback is welcome, but it is not yet suitable for production environments.
