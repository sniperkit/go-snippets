package main

import (
  "fmt"
  "github.com/codegangsta/cli"
  "github.com/gophergala/sbuca/server"
  "github.com/gophergala/sbuca/pkix"
  "os"
  "net/http"
  "net/url"
  "encoding/json"
)

func main() {

  app := cli.NewApp()
  app.Name = "sbuca"
  app.Usage = "Simple But Useful CA"
  app.Commands = []cli.Command{

    {
      Name: "server",
      Usage: "Run a CA server",
      Action: func(c *cli.Context) {
        host := os.Getenv("HOST")
        port := "8600"
        server.Run(host + ":" + port)
      },
    },

    {
      Name: "genkey",
      Usage: "Generate a RSA Private Key to STDOUT",
      Action: func (c *cli.Context) {
        key, err := pkix.NewKey()
        if err != nil {
          panic(err)
        }

        pem, err := key.ToPEM()
        if err != nil {
          panic(err)
        }

        fmt.Print(string(pem))
      },
    },

    {
      Name: "gencsr",
      Usage: "Generate a Certificate Request to STDOUT",
      Flags: []cli.Flag {
        cli.StringFlag {
          Name: "key",
          Usage: "RSA Private Key",
        },
      },
      Action: func (c *cli.Context) {
        keyName := c.String("key")
        if keyName == "" {
          fmt.Fprintln(os.Stderr, "[ERROR] Requere private key as parameter")
          return
        }

        key, err := pkix.NewKeyFromPrivateKeyPEMFile(keyName) 
        if err != nil {
          fmt.Fprintln(os.Stderr, "[ERROR] Failed to generate CSR: " + err.Error())
          return
        }

        csr, err := pkix.NewCertificateRequest(key)
        if err != nil {
          fmt.Fprintln(os.Stderr, "[ERROR] Failed to generate CSR: " + err.Error())
          return
        }

        pem, err := csr.ToPEM()
        if err != nil {
          fmt.Fprintln(os.Stderr, "[ERROR] Failed to generate CSR: " + err.Error())
          return
        }

        fmt.Print(string(pem))
      },
    },

    {
      Name: "submitcsr",
      Usage: "Submit a Certificate Request to CA",
      Flags: []cli.Flag {
        cli.StringFlag {
          Name: "host",
          Usage: "Host ip & port",
        },
        cli.StringFlag {
          Name: "format",
          Value: "cert",
          Usage: "output cert or id",
        },
      },
      Action: func (c *cli.Context){
        host := c.String("host")
        if host == "" {
          fmt.Fprintln(os.Stderr, "[ERROR] Requere host as parameter")
          return
        }

        format := c.String("format")
        if format != "cert" && format != "id" {
          fmt.Fprintln(os.Stderr, "[ERROR] format should be 'cert' or 'id'")
          return
        }

        args := c.Args()
        if len(args) == 0 {
          fmt.Fprintln(os.Stderr, "[ERROR] Should provide csr")
          return
        }
        csrName := c.Args().First()

        csr, err := pkix.NewCertificateRequestFromPEMFile(csrName)
        if err != nil {
          fmt.Fprintln(os.Stderr, "[ERROR] Failed to parse CSR: " + err.Error())
          return
        }

        //resp, err := http.Post("http://example.com/upload", "application/json", &buf)
        //var data interface{}
        pem, err := csr.ToPEM()
        if err != nil {
          fmt.Fprintln(os.Stderr, "[ERROR] Failed to parse CSR: " + err.Error())
          return
        }

        data := make(url.Values)
        data.Add("csr", string(pem))

        //resp, err := http.Post("http://" + host + "/certificates", "text/plain", values.Encode())
        resp, err := http.PostForm("http://" + host + "/certificates", data)
        if err != nil {
          fmt.Fprintln(os.Stderr, "[ERROR] Failed to request: " + err.Error())
          return
        }
        decoder := json.NewDecoder(resp.Body)
        respData := make(map[string]map[string]interface{})
        if err := decoder.Decode(&respData); err != nil {
          panic(err)
        }

        if format == "cert" {
          fmt.Print(respData["certificate"]["crt"])
        }
        if format == "id" {
          fmt.Println(respData["certificate"]["id"])
        }
      },
    },

    {
      Name: "getcrt",
      Flags: []cli.Flag {
        cli.StringFlag {
          Name: "host",
          Usage: "Host ip & port",
        },
      },
      Usage: "Get a Certificate from CA and output to STDOUT",
      Action: func (c *cli.Context){

        host := c.String("host")
        if host == "" {
          fmt.Fprintln(os.Stderr, "[ERROR] Requere host as parameter")
          return
        }

        args := c.Args()
        if len(args) == 0 {
          fmt.Fprintln(os.Stderr, "[ERROR] Should provide id (same as serial number)")
          return
        }
        id := c.Args().First()

        resp, err := http.Get("http://" + host + "/certificates/" + id)
        if err != nil {
          fmt.Fprintln(os.Stderr, "[ERROR] Failed to request: " + err.Error())
          return
        }

        decoder := json.NewDecoder(resp.Body)
        respData := make(map[string]map[string]interface{})
        if err := decoder.Decode(&respData); err != nil {
          panic(err)
        }

        fmt.Print(respData["certificate"]["crt"])

      },
    },

    {
      Name: "getcacrt",
      Flags: []cli.Flag {
        cli.StringFlag {
          Name: "host",
          Usage: "Host ip & port",
        },
      },
      Usage: "Get CA's Certificate and output to STDOUT",
      Action: func (c *cli.Context){

        host := c.String("host")
        if host == "" {
          fmt.Fprintln(os.Stderr, "[ERROR] Requere host as parameter")
          return
        }

        resp, err := http.Get("http://" + host + "/ca/certificate")
        if err != nil {
          fmt.Fprintln(os.Stderr, "[ERROR] Failed to request: " + err.Error())
          return
        }

        decoder := json.NewDecoder(resp.Body)
        respData := make(map[string]map[string]interface{})
        if err := decoder.Decode(&respData); err != nil {
          panic(err)
        }

        fmt.Print(respData["ca"]["crt"])

      },
    },

  }

  app.Run(os.Args)

}
