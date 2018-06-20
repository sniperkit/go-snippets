package lib

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const (
	registryTypePip   = "pip"
	registryTypeNpm   = "npm"
	registryTypeImage = "harbor"
)

//RequestMeta ...
type RequestMeta struct {
	RegistryType string
	HasHit       bool
	Metadata     map[string]string
}

type npmPackMeta struct {
	Tags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
}

type npmLoginMeta struct {
	Username string `json:"name"`
	Password string `json:"password"`
}

//Parser ...
type Parser func(req *http.Request) (RequestMeta, error)

//PipParser ...
func PipParser(req *http.Request) (RequestMeta, error) {
	meta := RequestMeta{}
	userAgent := req.Header.Get("User-Agent")
	if strings.Contains(userAgent, "pip") {
		if req.Method == http.MethodGet {
			path := req.URL.Path
			pkg := ""
			if strings.HasPrefix(path, "/packages/") && path != "/packages/" {
				//TODO:Very rough guess
				p := strings.TrimPrefix(path, "/packages/")
				pkg = strings.Split(p, "-")[0]
			} else {
				if strings.HasPrefix(path, "/simple") && path != "/simple/" {
					path = strings.TrimPrefix(path, "/simple")
				}
				pkg = strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
			}
			fmt.Printf("DEBUG: PIP install pkg: %s\n", pkg)
			meta.RegistryType = registryTypePip
			meta.HasHit = true
			meta.Metadata = map[string]string{
				"package":      pkg,
				"command":      "install",
				"full_command": fmt.Sprintf("%s %s", "pip install", pkg),
			}
		}

	}
	return meta, nil
}

//NpmParser ...
func NpmParser(req *http.Request) (RequestMeta, error) {
	userAgent := req.Header.Get("User-Agent")
	if strings.Contains(userAgent, "npm") {
		npmCmd := req.Header.Get("Referer")
		if len(npmCmd) > 0 {
			//Hit only when the command existing
			meta := RequestMeta{
				RegistryType: registryTypeNpm,
				HasHit:       true,
				Metadata:     make(map[string]string),
			}
			commands := strings.Split(npmCmd, " ")
			command := strings.TrimSpace(commands[0])
			meta.Metadata["command"] = command
			meta.Metadata["path"] = req.URL.String()
			meta.Metadata["extra"] = strings.TrimSpace(strings.TrimPrefix(npmCmd, command))
			meta.Metadata["session"] = req.Header.Get("Npm-Session")
			meta.Metadata["basic_auth"] = hex.EncodeToString([]byte(strings.TrimPrefix(req.Header.Get("Authorization"), "Basic ")))
			meta.Metadata["full_command"] = fmt.Sprintf("%s %s", "npm", npmCmd)

			//Read more info
			if command == "publish" || command == "adduser" {
				if req.Body != nil && req.ContentLength > 0 {
					buf, err := ioutil.ReadAll(req.Body)
					if err != nil {
						return RequestMeta{}, err
					}

					if command == "publish" {
						npmMetaJSON := &npmPackMeta{}
						if err := json.Unmarshal(buf, npmMetaJSON); err != nil {
							return RequestMeta{}, err
						}

						meta.Metadata["extra"] = npmMetaJSON.Tags.Latest
					} else if command == "adduser" {
						npmLoginJSON := &npmLoginMeta{}
						if err := json.Unmarshal(buf, npmLoginJSON); err != nil {
							return RequestMeta{}, err
						}
						meta.Metadata["basic_auth"] = hex.EncodeToString([]byte(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", npmLoginJSON.Username, npmLoginJSON.Password)))))
					}

					body := ioutil.NopCloser(bytes.NewBuffer(buf))
					req.Body = body
					req.ContentLength = int64(len(buf))
					req.Header.Set("Content-Length", strconv.Itoa(len(buf)))
				}
			}

			return meta, nil
		}
	}

	return RequestMeta{}, nil
}

//HarborParser ...
//Treat as deafault now
func HarborParser(req *http.Request) (RequestMeta, error) {
	return RequestMeta{
		RegistryType: registryTypeImage,
		HasHit:       true, //default handler
	}, nil
}

//ParserChain ...
type ParserChain struct {
	head        *parserWrapper
	tail        *parserWrapper
	commandList *CommandList
}

//ParserWrapper ...
type parserWrapper struct {
	parser Parser
	next   *parserWrapper
}

//Parse ...
func (pc *ParserChain) Parse(req *http.Request) (RequestMeta, error) {
	if pc.head == nil {
		return RequestMeta{}, errors.New("no parsers")
	}

	var errs []string
	p := pc.head
	for p != nil && p.parser != nil {
		if meta, err := p.parser(req); err != nil {
			errs = append(errs, err.Error())
		} else {
			if meta.HasHit {
				if len(meta.Metadata["full_command"]) > 0 {
					pc.commandList.Log(meta.Metadata["full_command"])
				}
				return meta, nil
			}
		}

		//next
		p = p.next
	}

	//No hit
	return RequestMeta{}, fmt.Errorf("%s:%s", "no hit", strings.Join(errs, ";"))
}

//Init ...
func (pc *ParserChain) Init() error {
	pc.head = nil
	pc.tail = nil

	if err := pc.Register(NpmParser); err != nil {
		return err
	}
	if err := pc.Register(PipParser); err != nil {
		return err
	}

	return pc.Register(HarborParser)
}

//Register ...
func (pc *ParserChain) Register(parser Parser) error {
	if parser == nil {
		return errors.New("nil parser")
	}

	if pc.head == nil {
		pc.head = &parserWrapper{parser, nil}
		pc.tail = pc.head

		return nil
	}

	pc.tail.next = &parserWrapper{parser, nil}
	pc.tail = pc.tail.next

	return nil
}
