# omniscience

## Github

### Index Github

Index Github by running

```
$ go run main.go --github_api_token_file_path=<github-api-token> --github_organization=<my-organization>
```

from within the `cmd/omniscience_client/github` directory.

### Generate Github credentials

For testing purposes, you can generate a Github Personal access token.

1. Go to Settings, then Developer settings
2. Select Personal access token from the menu on the left-hand side
2. Generate a token with the permission `read:org`

