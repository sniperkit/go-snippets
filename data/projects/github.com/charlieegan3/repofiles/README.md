#RepoFiles

Simple library to recursively list files in a GitHub repo and collect the file contents for each. Makes use of the [`trees`](https://developer.github.com/v3/git/trees/) endpoint.

## Usage

```go
package main

import (
  "fmt"

  "github.com/charlieegan3/repofiles"
)

func main() {
	credentials = Credentials{User: os.Getenv("GITHUB_USER"), Token: os.Getenv("GITHUB_TOKEN")}
    repo := repofiles.NewRepo("username", "repoName", "master or SHA")
    repo.List(credentials) //collect the list files
    files := repo.Files(credentials) //request the contents of each file

    fmt.Println(files)
}

```

Output:

```
[{file1.ext file1content} {file2.ext file1content} ...]
```
