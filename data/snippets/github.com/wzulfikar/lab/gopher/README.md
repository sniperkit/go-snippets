This `gopher` directory is symlink-ed to `/Users/strawhat/go/src/local-dev` so you can develop go packages locally and import using `local-dev/{package-name}`.

Once pushed to public repo (github, gopkg, etc), refactor the `local-dev` to the respective repo. ie. from `local-dev/go-nem-clerk` to `github.com/wzulfiar/go-nem-clerk`.