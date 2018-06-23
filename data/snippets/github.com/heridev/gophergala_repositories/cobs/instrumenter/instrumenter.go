package instrumenter

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/gophergala/cobs/types"
	"github.com/iron-io/iron_go/mq"
)

var rc redis.Conn

func Run(imageId string) {
	var err error
	log.Println("Connecting to Redis")
	rc, err = redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatalf("Issues with redis: %s", err)
	}
	defer rc.Close()

	df, _ := redis.String(rc.Do("GET", "dockerfile-"+imageId))
	info, _ := redis.Bytes(rc.Do("GET", "info-"+imageId))
	var image types.ImageInfo
	json.Unmarshal(info, &image)

	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	var files = []struct {
		Name, Body string
	}{
		{"Dockerfile.txt", df},
	}

	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Size: int64(len(file.Body)),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			log.Fatalln(err)
		}
		if _, err := tw.Write([]byte(file.Body)); err != nil {
			log.Fatalln(err)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatalln(err)
	}

	rc.Do("SET", "tarball-"+imageId, buf.Bytes())
	queue := mq.New("builder-" + image.Architecture)
	id, err := queue.PushString(imageId)
	if err != nil {
		log.Fatalf("error pushing to queue: %s", err)
	}
	log.Printf("build queued: %s\n", id)
}
