package lib

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	npmUserSessionTimeout = 3600 //seconds
)

//ProxyTarget ...
type ProxyTarget string

//Scheduler ...
type Scheduler struct {
	pool       *RuntimePool
	imageStore *ImageStore
	executor   *Executor
	packer     *Packer
	ctx        context.Context
	drivers    map[string]ScheduleDriver
	exitChan   chan struct{}
	doneChan   chan struct{}
}

//NewScheduler ...
func NewScheduler(ctx context.Context) *Scheduler {
	return &Scheduler{
		pool:       NewRuntimePool(),
		imageStore: NewImageStore(),
		executor:   NewExecutor(Config.Dockerd.Host, Config.Dockerd.Port, Config.Harbor.Host),
		packer:     NewPacker(Config.Dockerd.Host, Config.Dockerd.Port, Config.Harbor.Host),
		ctx:        ctx,
		exitChan:   make(chan struct{}, 1),
		doneChan:   make(chan struct{}, 1),
	}
}

//Start ...
func (s *Scheduler) Start() {
	go s.sweepRuntimes()
	go s.sweepImages()

	registryAPI := fmt.Sprintf("%s://%s/api", Config.Harbor.Protocol, Config.Harbor.Host)
	s.drivers = make(map[string]ScheduleDriver)
	s.drivers[registryTypeNpm] = NewNpmScheduleDriver(registryAPI, Config.NpmRegistry.Namespace)
	s.drivers[registryTypePip] = NewPipScheduleDriver(registryAPI, Config.PipRegistry.Namespace)

	log.Println("Scheduler is started")
}

func (s *Scheduler) sweepRuntimes() {
	defer func() {
		log.Println("Runtime sweeper exit")
		s.doneChan <- struct{}{}
	}()

	tk := time.NewTicker(30 * time.Second)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			//Garbage collection
			garbages := s.pool.Garbages()
			if len(garbages) > 0 {
				for _, v := range garbages {
					//Clear
					if err := s.executor.Destroy(v.ID); err != nil {
						log.Fatalf("garbage collection %s error: %s\n", v.ID, err)
					} else {
						log.Printf("Destroy container instance: %s\n", v.ID)
					}
				}
			}
		case <-s.ctx.Done():
			return
		case <-s.exitChan:
			return
		}
	}
}

func (s *Scheduler) sweepImages() {
	defer func() {
		log.Println("Image sweeper exit")
		s.doneChan <- struct{}{}
	}()

	tk := time.NewTicker(30 * time.Second)
	defer tk.Stop()

	for {
		select {
		case <-tk.C:
			images := s.imageStore.Garbage()
			if len(images) > 0 {
				for _, image := range images {
					theImage := fmt.Sprintf("%s:%s", image.Name, image.Tag)
					if err := s.packer.RMImage(theImage); err != nil {
						log.Printf("Failed to sweep outdated image '%s' with error:%s\n", theImage, err)
					}
				}
			}
		case <-s.ctx.Done():
			return
		case <-s.exitChan:
			return
		}
	}
}

//Stop scheduler
func (s *Scheduler) Stop() {
	defer log.Println("Scheduler is stopped")
	//Stop 1st loop
	s.exitChan <- struct{}{}
	<-s.doneChan
	//Stop 2nd loop
	s.exitChan <- struct{}{}
	<-s.doneChan
}

//Schedule ...
func (s *Scheduler) Schedule(meta RequestMeta) (ServeEnvironment, error) {
	driver, ok := s.drivers[meta.RegistryType]
	if !ok {
		return ServeEnvironment{}, fmt.Errorf("registry type %s not support", meta.RegistryType)
	}

	policy := driver.Schedule(meta)
	if len(policy.ReuseIdentity) > 0 {
		key := fmt.Sprintf("%s:%s", meta.RegistryType, policy.ReuseIdentity)
		_, yes := s.pool.Index(key)
		if yes {
			r, err := s.pool.Use(key)
			if err != nil {
				return ServeEnvironment{}, err
			}
			log.Printf("Reuse %s: %s\n", r.ID, r.Target)
			if policy.Rebuild != nil {
				policy.Rebuild.BaseContainer = r.ID
			}
			return ServeEnvironment{
				Target:      r.Target,
				Rebuild:     policy.Rebuild,
				InstanceKey: key,
			}, nil
		}
	}

	//Create
	s.executor.SetNamespace(policy.Namespace)
	imageKey := fmt.Sprintf("%s:%s", policy.Image, policy.SessionTag)
	if _, ok := s.imageStore.Get(imageKey); ok {
		policy.Tag = policy.SessionTag
	}
	env, err := s.executor.Exec(policy)
	if err != nil {
		return ServeEnvironment{}, err
	}

	log.Printf("Start new service instance: %s\n", env.RuntimeID)

	key := env.RuntimeID //Just for garbage collection
	if len(policy.ReuseIdentity) > 0 {
		key = policy.ReuseIdentity
	}
	key = fmt.Sprintf("%s:%s", meta.RegistryType, key)

	r := &Runtime{
		ID:         env.RuntimeID,
		Target:     env.Target,
		ActiveTime: time.Now().Unix(),
		Image:      fmt.Sprintf("%s:%s", policy.Image, policy.Tag),
	}
	if err := s.pool.Put(key, r); err != nil {
		//let's see if problem will appear
		log.Printf("Pool error: %s\n", err)
	}

	if policy.Rebuild != nil {
		policy.Rebuild.BaseContainer = env.RuntimeID
	}
	return ServeEnvironment{
		Target:      env.Target,
		Rebuild:     policy.Rebuild,
		InstanceKey: key,
	}, nil
}

//Rebuild ...
func (s *Scheduler) Rebuild(policy *BuildPolicy) error {
	if policy == nil {
		return errors.New("nil build policy")
	}

	if len(policy.Image) == 0 || len(policy.Tag) == 0 {
		return errors.New("target image or tag is invalid")
	}

	if len(policy.BaseContainer) == 0 {
		return errors.New("no base container for build")
	}

	if policy.NeedPush {
		s.packer.SetNamespace(policy.Namespace)
		return s.packer.Build(policy.BaseContainer, policy.Image, policy.Tag)
	}

	return s.packer.BuildLocal(policy.BaseContainer, policy.Image, policy.Tag)
}

//FreeRuntime mark the instance as idle
func (s *Scheduler) FreeRuntime(key string) error {
	go func() {
		<-time.After(2 * time.Second)
		s.pool.SetIdle(key)
	}()

	return nil
}

//GetRuntimes get all runtimes including the destroyed ones
func (s *Scheduler) GetRuntimes() []*Runtime {
	return s.pool.GetAll()
}

//StoreImage ...
func (s *Scheduler) StoreImage(image, tag string) {
	s.imageStore.Put(image, tag)
}

//ServeEnvironment ...
type ServeEnvironment struct {
	Target      ProxyTarget
	Rebuild     *BuildPolicy
	InstanceKey string
}

//SchedulePolicy ...
type SchedulePolicy struct {
	Image         string
	Tag           string
	SessionTag    string
	UseHub        bool
	ReuseIdentity string
	BoundPorts    []int
	Rebuild       *BuildPolicy
	EnvVars       map[string]string
	Namespace     string
}

//BuildPolicy ...
type BuildPolicy struct {
	BaseContainer string `json:"base_container"`
	Image         string `json:"image"`
	Tag           string `json:"tag"`
	NeedPush      bool   `json:"need_push"`
	Namespace     string `json:"namespace"`
	NeedStore     bool   `json:"need_store"`
}

//Encode ...
func (bp *BuildPolicy) Encode() (string, error) {
	bytes, err := json.Marshal(bp)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

//Decode ...
func (bp *BuildPolicy) Decode(data string) error {
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, bp)
}

//ScheduleDriver ...
type ScheduleDriver interface {
	//Schedule ...
	Schedule(meta RequestMeta) *SchedulePolicy
}

//PipScheduleDriver ...
type PipScheduleDriver struct {
	registryAPI       string
	registryNamespace string
	httpClient        *http.Client
}

//NewPipScheduleDriver ...
func NewPipScheduleDriver(registryAPI, registryNamespace string) *PipScheduleDriver {
	return &PipScheduleDriver{
		registryAPI:       registryAPI,
		registryNamespace: registryNamespace,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

//Schedule ...
func (psd *PipScheduleDriver) Schedule(meta RequestMeta) *SchedulePolicy {
	if !meta.HasHit || meta.RegistryType != registryTypePip {
		return nil
	}
	if meta.Metadata["command"] == "install" {
		//TODO: change
		image := fmt.Sprintf("pip-project/pypi-%s", meta.Metadata["package"])
		//Default policy
		policy := &SchedulePolicy{
			Image:         image,
			Tag:           "dev",
			BoundPorts:    []int{80},
			ReuseIdentity: meta.Metadata["package"],
			EnvVars:       map[string]string{"PYPI_EXTRA": "--disable-fallback", "PYPI_ROOT": "/pypi"},
			Namespace:     psd.registryNamespace,
		}
		return policy

	}

	log.Printf("Unknown command for pip package: %s\n", meta.Metadata["command"])

	return nil
}

//NpmScheduleDriver ...
type NpmScheduleDriver struct {
	registryAPI       string
	registryNamespace string
	httpClient        *http.Client
}

//NewNpmScheduleDriver ...
func NewNpmScheduleDriver(registryAPI, registryNamespace string) *NpmScheduleDriver {
	return &NpmScheduleDriver{
		registryAPI:       registryAPI,
		registryNamespace: registryNamespace,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

//Schedule ...
func (nsd *NpmScheduleDriver) Schedule(meta RequestMeta) *SchedulePolicy {
	if !meta.HasHit || len(meta.Metadata) == 0 || meta.RegistryType != registryTypeNpm {
		return nil
	}

	//Default policy
	policy := &SchedulePolicy{
		Image:      Config.NpmRegistry.BaseImage,
		Tag:        Config.NpmRegistry.BaseImageTag,
		UseHub:     true,
		BoundPorts: []int{80},
		Rebuild: &BuildPolicy{
			Image:     Config.NpmRegistry.BaseImage,
			Tag:       Config.NpmRegistry.BaseImageTag,
			Namespace: nsd.registryNamespace,
		},
		Namespace: nsd.registryNamespace,
	}

	//If has reuseIdentity
	session := meta.Metadata["session"]
	if len(session) > 0 {
		policy.ReuseIdentity = session
	}

	//Special cases
	requestPath := meta.Metadata["path"]
	command := meta.Metadata["command"]
	if command == "view" || command == "install" {
		extraInfo := meta.Metadata["extra"]
		repo := strings.TrimPrefix(requestPath, "/")
		tag := strings.TrimSpace(strings.TrimPrefix(extraInfo, repo+"@"))
		if checkImageExisting(nsd.registryAPI, nsd.registryNamespace, repo, tag, nsd.httpClient) {
			policy.Image = repo
			policy.Tag = tag
			policy.UseHub = false
		}
		policy.Rebuild = nil
	}

	if command == "login" || command == "adduser" || command == "add-user" {
		if strings.Contains(requestPath, "org.couchdb.user:") &&
			!strings.Contains(requestPath, "/-rev/") { //for login
			policy.Rebuild.Tag = meta.Metadata["basic_auth"]
			//Store it for next actions with same auth info
			policy.Rebuild.NeedStore = true
		}
	}

	if command == "publish" {
		repo := strings.TrimPrefix(requestPath, "/")
		tag := meta.Metadata["extra"]
		log.Printf("PUBLISH: %s@%s", repo, tag)
		if checkImageExisting(nsd.registryAPI, nsd.registryNamespace, repo, tag, nsd.httpClient) {
			policy.Image = repo
			policy.Tag = tag
			policy.UseHub = false
		} else {
			sessionTag := meta.Metadata["basic_auth"]
			if len(sessionTag) > 0 {
				policy.SessionTag = sessionTag
			}
		}
		policy.Rebuild.Image = strings.TrimPrefix(requestPath, "/")
		policy.Rebuild.Tag = meta.Metadata["extra"]
		policy.Rebuild.NeedPush = true
	}

	return policy
}

func checkImageExisting(registryAPI, registryNamespace, image, tag string, client *http.Client) bool {
	url := fmt.Sprintf("%s%s%s%s%s%s%s", registryAPI, "/repositories/", registryNamespace, "/", image, "/tags/", tag)
	resp, err := client.Get(url)
	if err == nil {
		if resp.StatusCode == http.StatusOK {
			log.Printf("Image %s:%s existing\n", image, tag)
			return true
		}

		log.Printf("Image %s:%s not existing\n", image, tag)
		return false
	}

	log.Printf("Failed to Check image %s:%s existing: %s\n", image, tag, err)
	return false
}
