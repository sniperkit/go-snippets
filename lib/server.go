package lib

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

//ProxyServer serves the requests
type ProxyServer struct {
	server     *http.Server
	proxy      *httputil.ReverseProxy
	context    context.Context
	reqParser  *ParserChain
	scheduler  *Scheduler
	apiHandler *APIHandler
}

//NewProxyServer create new server instance
func NewProxyServer(ctx context.Context) *ProxyServer {
	commandList := NewCommandList()
	scheduler := NewScheduler(ctx)
	apiHandler := &APIHandler{
		scheduler:   scheduler,
		commandList: commandList,
	}
	parser := &ParserChain{
		commandList: commandList,
	}

	return &ProxyServer{
		apiHandler: apiHandler,
		scheduler:  scheduler,
		context:    ctx,
		reqParser:  parser,
	}
}

//Start the proxy server
func (ps *ProxyServer) Start() error {
	if err := ps.reqParser.Init(); err != nil {
		return err
	}

	ps.scheduler.Start()

	if ps.proxy == nil {
		t := &http.Transport{
			//Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			MaxIdleConnsPerHost:   100,
			MaxIdleConns:          100,
			IdleConnTimeout:       120 * time.Second,
			TLSHandshakeTimeout:   20 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		ps.proxy = &httputil.ReverseProxy{
			Transport: t,
			Director: func(req *http.Request) {
				log.Printf("INCOMING REQ: %s %s\n", req.Method, req.URL.String())
				//Log sessions
				session := []string{}
				for _, c := range req.Cookies() {
					session = append(session, c.String())
				}
				//For npm
				npmSession := req.Header.Get("Npm-Session")
				if len(npmSession) > 0 {
					session = append(session, fmt.Sprintf("%s:%s", "Npm-Session", npmSession))
				}
				log.Printf("SESSION: %s\n", strings.Join(session, "; "))
				log.Printf("HEADER: %#-v\n", req.Header)

				//Parse request
				if ps.reqParser != nil {
					meta, err := ps.reqParser.Parse(req)
					if err != nil {
						log.Printf("[ERROR]: Parse error: %s\n", err)
						return
					}

					if meta.HasHit {
						var rawTarget string
						if meta.RegistryType == registryTypeNpm || meta.RegistryType == registryTypePip {
							env, err := ps.scheduler.Schedule(meta)
							if err != nil {
								log.Printf("[ERROR]: schedule error: %s\n", err)
								return
							}
							rawTarget = fmt.Sprintf("%s%s", "http://", env.Target)

							if env.Rebuild != nil {
								h, err := env.Rebuild.Encode()
								if err != nil {
									log.Printf("[ERROR]: set rebuild header failed: %s", err)
									return
								}
								req.Header.Set("registry-factory", h)
							}

							//Set instance key for status updating
							if len(env.InstanceKey) > 0 {
								req.Header.Set("instance-key", env.InstanceKey)
							}
						} else {
							//Treat as management/harbor
							rawTarget = fmt.Sprintf("%s://%s", Config.Harbor.Protocol, Config.Harbor.Host)
						}

						target, err := url.Parse(rawTarget)
						if err != nil {
							log.Printf("[ERROR]: Url parse error: %s\n", err)
							return
						}
						targetQuery := target.RawQuery
						req.URL.Scheme = target.Scheme
						req.URL.Host = target.Host
						req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
						if targetQuery == "" || req.URL.RawQuery == "" {
							req.URL.RawQuery = targetQuery + req.URL.RawQuery
						} else {
							req.URL.RawQuery = targetQuery + req.URL.RawQuery
						}
						if _, ok := req.Header["User-Agent"]; !ok {
							// explicitly disable User-Agent so it's not set to default value
							req.Header.Set("User-Agent", "")
						}

						log.Printf("PROXY TO: %s\n", req.URL.String())
					}
				}
				//do nothing
			},

			ModifyResponse: func(res *http.Response) error {
				log.Printf("RESPONSE: %s\n", res.Status)
				//Request served
				//Do not care the response status code
				instanceKey := res.Request.Header.Get("instance-key")
				if len(instanceKey) > 0 {
					//delay several seconds
					go func() {
						time.Sleep(5 * time.Second)
						ps.scheduler.FreeRuntime(instanceKey)
					}()
				}
				if res.StatusCode >= http.StatusOK && res.StatusCode <= http.StatusAccepted {
					rebuildPolicyHeader := res.Request.Header.Get("registry-factory")
					if len(rebuildPolicyHeader) > 0 {
						//Use async way to improve pref, this may cause inconsistent case
						//client get success code
						//but the package image may not be pushed to harbor
						go func() {
							rebuildPolicy := &BuildPolicy{}
							err := rebuildPolicy.Decode(rebuildPolicyHeader)
							if err != nil {
								log.Printf("[ERROR]: Failed to decode rebuild policy with error:%s\n", err)
							}

							log.Printf("Rebuild image (base container): %s:%s (%s)\n", rebuildPolicy.Image, rebuildPolicy.Tag, rebuildPolicy.BaseContainer)

							if err := ps.scheduler.Rebuild(rebuildPolicy); err != nil {
								log.Printf("[ERROR]: Failed to rebuild image: %s:%s\n", rebuildPolicy.Image, rebuildPolicy.Tag)
							}

							//Store image for future use
							if rebuildPolicy.NeedStore {
								ps.scheduler.StoreImage(rebuildPolicy.Image, rebuildPolicy.Tag)
								log.Printf("Store image: %s:%s\n", rebuildPolicy.Image, rebuildPolicy.Tag)
							}
						}()
					}
				}

				return nil
			},
		}
	}

	if ps.server == nil {
		ps.server = &http.Server{
			Addr: fmt.Sprintf("%s:%d", Config.Host, Config.Port),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if ps.apiHandler.IsMatchedRequests(r) {
					ps.apiHandler.ServeHTTP(w, r)
					return
				}
				ps.proxy.ServeHTTP(w, r)
			}),
		}
	}

	return ps.server.ListenAndServe()
}

//Stop the proxy server
func (ps *ProxyServer) Stop() error {
	if ps.server == nil {
		return errors.New("No server existing")
	}

	//Stop scheduler
	ps.scheduler.Stop()

	ctx, cancel := context.WithTimeout(ps.context, 30*time.Second)
	defer cancel()

	return ps.server.Shutdown(ctx)
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
