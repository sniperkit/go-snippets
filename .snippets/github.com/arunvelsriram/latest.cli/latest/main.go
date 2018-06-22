package main

import (
  "fmt"
  "os"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "time"
  "sort"
  "strings"

  "github.com/urfave/cli"
)

// Package : struct representing a package
type Package struct {
  Name string
  Version string
}

// Command : command interface for all latest* command implementations
type Command interface {
  Execute() error
}

// LatestGem : LatestGem command
type LatestGem struct {
  repo string
  name string
}

// Execute : gets latest version of gem
func (latestGem LatestGem) Execute() (error) {
  fmt.Println()
  fmt.Println("Ruby Gem:")

  pkg, err := getPackage( fmt.Sprintf("%s/%s.json", latestGem.repo, latestGem.name))
  if err != nil {
    return err
  }

  fmt.Println(pkg.Name, pkg.Version)

  return nil
}

// LatestNodeModule : LatestNodeModule command
type LatestNodeModule struct {
  repo string
  name string
}

// Execute : gets latest version of node module
func (latestNodeModule LatestNodeModule) Execute() (error) {
  fmt.Println()
  fmt.Println("Node Module:")

  pkg, err := getPackage(fmt.Sprintf("%s/%s/latest", latestNodeModule.repo, latestNodeModule.name))
  if err != nil {
    return err
  }

  fmt.Println(pkg.Name, pkg.Version)

  return nil
}

// NotFoundError : an error implementation used when a package is not availabe in the source repository
type NotFoundError struct {
  Message string
}

func (error *NotFoundError) Error() (string) {
  return error.Message
}


func getPackage(url string) (Package, error) {
  var pkg Package
  httpClient := &http.Client{
    Timeout: time.Second * 10,
  }

  resp, err := httpClient.Get(url)
  if err != nil {
    return pkg, fmt.Errorf("Failed to get response from %s", url)
  }

  if(resp.StatusCode == http.StatusNotFound) {
    return pkg, &NotFoundError{"Not found!"}
  }

  jsonBlob, err := ioutil.ReadAll(resp.Body)
  resp.Body.Close()
  if err != nil {
    return pkg, fmt.Errorf("Failed to read response body")
  }

  err = json.Unmarshal(jsonBlob, &pkg)
  if err != nil {
    return pkg, fmt.Errorf("Failed to parse JSON")
  }

  return pkg, nil
}

func latestGemCommand(gemName string) (LatestGem) {
  return LatestGem{"https://rubygems.org/api/v1/gems", gemName}
}

func latestNodeModuleCommand(nodeModuleName string) (LatestNodeModule) {
  return LatestNodeModule{"https://registry.npmjs.org", nodeModuleName}
}

func isEmpty(content string) bool {
  if strings.TrimSpace(content) == "" {
    return true
  }

  return false
}

func exitStatus(err error) (int) {
  if _, ok := err.(*NotFoundError); ok {
    return 0
  }

  return 1
}

func main() {
  app := cli.NewApp()
  app.Name = "latest"
  app.Version = "v1.0.2"
  app.Usage = "A CLI to find the latest version of a Ruby Gem, Node module, Java JAR etc."
  app.ArgsUsage = "<name>"

  before :=  func(cliContext *cli.Context) (error) {
    name := cliContext.Args().Get(0)
    if ok := isEmpty(name); ok {
      err := fmt.Errorf("name not given")
      return cli.NewExitError(err, exitStatus(err))
    }

    return nil
  }

  app.Commands = []cli.Command {
    {
      Name:  "gem",
      Aliases: []string{"g"},
      Usage: "get latest version of ruby gem <name>",
      ArgsUsage: "<name>",
      Before: before,
      Action: func(cliContext *cli.Context) error {
       name := cliContext.Args().Get(0)
       command := latestGemCommand(name)
       if err := command.Execute(); err != nil {
         return cli.NewExitError(err, exitStatus(err))
       }

        return nil
      },
    },
    {
      Name: "node-module",
      Aliases: []string{"n"},
      Usage: "get latest version of node module <name>",
      ArgsUsage: "<name>",
      Before: before,
      Action: func(cliContext *cli.Context) error {
        name := cliContext.Args().Get(0)
        command := latestNodeModuleCommand(name)
        if err := command.Execute(); err != nil {
          return cli.NewExitError(err, exitStatus(err))
        }

        return nil
      },
    },
  }

  app.Before = before
  app.Action = func (cliContext *cli.Context) error {
    var errs []error
    computedExitStatus := 0
    name := cliContext.Args().Get(0)
    commands := []Command{latestGemCommand(name), latestNodeModuleCommand(name)}

    for _, command := range commands {
      if err := command.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        errs = append(errs, err)
      }
    }

    for _, err := range errs {
      computedExitStatus += exitStatus(err)
    }

    os.Exit(computedExitStatus) 
    
    return nil
  }

  sort.Sort(cli.CommandsByName(app.Commands))

  app.Run(os.Args)
}
