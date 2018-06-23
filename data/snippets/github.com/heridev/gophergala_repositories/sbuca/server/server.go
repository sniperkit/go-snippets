package server

import (
  "fmt"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/render"
  //"github.com/gophergala/sbuca/x509util"
  "net/http"
  "github.com/gophergala/sbuca/pkix"
  "github.com/gophergala/sbuca/ca"
  "strconv"
)


func Run(addr string) {
  fmt.Print("start...")

  m := martini.Classic()
  m.Use(render.Renderer())

  //FIXME
  ca.NewCA(".")

  m.Get("/", func() string {
    return "Hello sbuca"
  })
  m.Get("/ca/certificate", func(req *http.Request, params martini.Params, r render.Render) {

    format := req.URL.Query().Get("format")

    newCA, err := ca.NewCA(".")
    if err != nil {
      panic(err)
    }

    pem, err := newCA.Certificate.ToPEM()
    if err != nil {
      panic(err)
    }

    if format == "file" {
      r.Data(200, pem)
    } else {
      r.JSON(200, map[string]interface{}{
        "ca": map[string]interface{}{
          "crt": string(pem),
        },
      })
    }
  })
  m.Get("/certificates/:id", func(req *http.Request, params martini.Params, r render.Render) {

    format := req.URL.Query().Get("format")

    newCA, err := ca.NewCA(".")
    if err != nil {
      panic(err)
    }

    id := params["id"]
    idInt, err := strconv.Atoi(id)
    if err != nil {
      r.JSON(401, map[string]interface{}{
        "result": "wrong id",
      })
      return
    }
    cert, err := newCA.GetCertificate(int64(idInt))
    if err != nil {
      r.JSON(401, map[string]interface{}{
        "result": "cannot get cert",
      })
      return
    }

    pem, err := cert.ToPEM()
    if err != nil {
      r.JSON(401, map[string]interface{}{
        "result": "cannot get cert",
      })
      return
    }

    if format == "file" {
      r.Data(200, pem)
    } else {
      r.JSON(200, map[string]interface{}{
        "certificate": map[string]interface{}{
          "id": cert.GetSerialNumber().Int64(),
          "crt": string(pem),
          //"csr": csr,
        },
      })
    }

  })
  m.Post("/certificates", func(req *http.Request, params martini.Params, r render.Render) {

    csrString := req.PostFormValue("csr")
    format := req.URL.Query().Get("format")

    csr, err := pkix.NewCertificateRequestFromPEM([]byte(csrString))
    if err != nil {
      panic(err)
    }

    newCA, err := ca.NewCA(".")
    if err != nil {
      panic(err)
    }

    cert, err := newCA.IssueCertificate(csr)
    if err != nil {
      panic(err)
    }

    certPem, err := cert.ToPEM()
    if err != nil {
      panic(err)
    }
    if format == "file" {
      r.Data(200, certPem)
    } else {
      r.JSON(200, map[string]interface{}{
        "certificate": map[string]interface{}{
          "id": cert.GetSerialNumber().Int64(),
          "crt": string(certPem),
          //"csr": csr,
        },
      })
    }
  })

  m.RunOnAddr(addr)

}
