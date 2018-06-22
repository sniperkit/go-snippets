package backend

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/garyburd/redigo/redis"

	"code.google.com/p/go-uuid/uuid"
	"github.com/codegangsta/negroni"
	"github.com/gophergala/cobs/hunter"
	"github.com/gophergala/cobs/types"
	"github.com/gorilla/mux"
)

var rc redis.Conn

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Nothing here yet. Try: \n\n  curl -X POST -F repository=redis -F arch=armv7l http://cobs.aas.io/search"))
}

func BuildTarballHandler(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//	case "POST":
	//		//var buf []byte
	//		file, _, _ := r.FormFile("file")
	//		defer file.Close()
	//		data, _ := ioutil.ReadAll(file)
	//		//file.Read(buf)
	//		rc.Do("SET", mux.Vars(r)["imageid"], data)
	default:
		data, _ := redis.Bytes(rc.Do("GET", mux.Vars(r)["imageid"]))
		rw.Write(data)
	}
}

func BuildDockerfileHandler(rw http.ResponseWriter, r *http.Request) {
	imageId := mux.Vars(r)["imageid"]
	data, _ := redis.Bytes(rc.Do("GET", "dockerfile-"+imageId))
	rw.Write(data)
}

func ImageInfoHandler(rw http.ResponseWriter, r *http.Request) {
	imageId := mux.Vars(r)["imageid"]
	data, _ := redis.Bytes(rc.Do("GET", "info-"+imageId))
	rw.Write(data)
}

func BuildStatusHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte(mux.Vars(r)["imageid"]))
}

func MakeNewImageName(name, arch string) string {
	n := strings.Replace(name, "/", "-", -1)
	return "cobs" + arch + "/" + n
}

func RequestImageBuild(repository, arch, tag string) string {
	imageId := uuid.New()
	rc.Do("SET", "url-"+imageId, repository)
	rc.Do("SET", "name-"+repository+"-arch-"+arch+"-tag-"+tag, imageId)
	newName := MakeNewImageName(repository, arch)
	data, _ := json.Marshal(types.ImageInfo{repository, newName, arch, tag})
	rc.Do("SET", "info-"+imageId, data)
	go hunter.GoHunting(imageId)

	return imageId
}

//func BuildHandler(rw http.ResponseWriter, r *http.Request) {
//	//imageid := mux.Vars(r)["imageid"]
//	switch r.Method {
//	case "POST":
//		repository := r.FormValue("repository")
//		tag := r.FormValue("tag")
//		arch := r.FormValue("arch")
//		imageId := RequestImageBuild(repository, arch, tag)
//		rw.Write([]byte(imageId))
//	default:
//		data, _ := redis.Bytes(rc.Do("GET", "tarball"))
//		rw.Write(data)
//	}
//}

func RepoSearch(repo string) string {
	res := hunter.SearchDockerRegistry(repo)
	return res[0].Name

}

func SearchHandler(rw http.ResponseWriter, r *http.Request) {
	//imageid := mux.Vars(r)["imageid"]
	switch r.Method {
	case "POST":
		repository := r.FormValue("repository")
		tag := r.FormValue("tag")
		arch := r.FormValue("arch")

		if tag == "" {
			tag = "latest"
		}
		if arch == "" {
			arch = "x8664"
		}

		imageId, _ := redis.String(rc.Do("GET", "name-"+repository+"-arch-"+arch+"-tag-"+tag))
		if len(imageId) == 0 {
			name := RepoSearch(repository)
			imageId, _ = redis.String(rc.Do("GET", "name-"+name+"-arch-"+arch+"-tag-"+tag))
			if len(imageId) == 0 {
				imageId = RequestImageBuild(name, arch, tag)
			}
		}
		rw.Write([]byte("http://cobs.aas.io/api/v1/info/" + imageId))
	default:
		data, _ := redis.Bytes(rc.Do("GET", "tarball"))
		rw.Write(data)
	}
}

func FakeHandler(rw http.ResponseWriter, r *http.Request) {
	repository := r.FormValue("repository")
	newName := r.FormValue("new")
	tag := r.FormValue("tag")
	arch := r.FormValue("arch")

	imageId := uuid.New()
	rc.Do("SET", "url-"+imageId, repository)
	rc.Do("SET", "name-"+repository+"-arch-"+arch+"-tag-"+tag, imageId)
	data, _ := json.Marshal(types.ImageInfo{repository, newName, arch, tag})
	rc.Do("SET", "info-"+imageId, data)

	log.Println(data)
	rw.Write(data)

}

func Run() {
	var err error
	log.Println("Connecting to Redis")
	rc, err = redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatalf("Issues with redis: %s", err)
	}
	defer rc.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)

	r.HandleFunc("/api/v1/build/{imageid}/tarball", BuildTarballHandler)
	r.HandleFunc("/api/v1/build/{imageid}/dockerfile", BuildDockerfileHandler)
	r.HandleFunc("/api/v1/build/{imageid}", BuildStatusHandler)
	r.HandleFunc("/api/v1/info/{imageid}", ImageInfoHandler)
	r.HandleFunc("/api/v1/build/fake/", FakeHandler)

	//r.HandleFunc("/api/v1/build/", BuildHandler)
	r.HandleFunc("/search", SearchHandler)

	n := negroni.Classic()
	n.UseHandler(r)

	// shouldn't this be automatic?
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	port = ":" + port
	n.Run(port)
}
