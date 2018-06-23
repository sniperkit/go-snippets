package main

import (
  "github.com/koding/kite"
  "fmt"
)

func main() {
  // Create a kite
  k := kite.New("math", "1.0.0")

  // Add pre handler method.
  k.PreHandleFunc(func(r *kite.Request) (interface{}, error) {
    fmt.Println("\nThis pre handler is executed before the method is executed")
    resp := "hello from pre handler!"

    // let us return an hello to base square method!
    r.Context.Set("response", resp)
    return resp, nil
  })

  // Add post handler method.
  k.PostHandleFunc(func(r *kite.Request) (interface{}, error) {
    fmt.Println("This post handler is executed after the method is executed")

    // Pass the response from the previous square method back to the
    // client, this is imporant if you use post handler.
    return r.Context.Get("response")
  })

  // Add our handler method with the name "square"
  k.HandleFunc("square", func(r *kite.Request) (interface{}, error) {
      a := r.Args.One().MustFloat64()
      result := a * a    // calculate the square
      return result, nil // send back the result
  }).DisableAuthentication()

  // Attach to a server with port 8080 and run it
  k.Config.Port = 8080
  k.Run()
}
