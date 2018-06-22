package transcoder

import (
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/vimeo/go-magic/magic"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Server *martini.ClassicMartini
type FileForm struct {
	Name       string                `form:"name"`
	FileUpload *multipart.FileHeader `form:"fileUpload"`
}

type TranscodingSession struct {
	Name string        `json:"name"`
	Type string        `json:"type"`
	Path string        `json:"path"`
	Id   bson.ObjectId `json:"id" bson:"_id,omitempty"`
}

func NewServer(session *DatabaseSession) Server {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(session.Database())
	setupRoutes(m)
	return m
}

func setupRoutes(m Server) {
	m.Get("/", Home)
	m.Post("/api/sessions/new", binding.MultipartForm(FileForm{}), NewTranscodingSession)
	m.Get("/api/sessions/:id", GetSession)
}

func Home() string {
	return "Hello abbita!"
}

func NewTranscodingSession(fileForm FileForm, r render.Render, db *mgo.Database) {
	file, err := fileForm.FileUpload.Open()
	HandleError(err, r)
	stringFilePrefix := generate_id()
	fileName := fmt.Sprintf("media/%s-%s", stringFilePrefix, fileForm.Name)
	out, err := os.Create(fileName)
	HandleError(err, r)
	defer out.Close()
	_, err = io.Copy(out, file)
	HandleError(err, r)
	fileType := magic.MimeFromFile(fileName)

	session := &TranscodingSession{
		Name: fileForm.Name,
		Path: fileName,
		Type: fileType,
		Id:   bson.NewObjectId(),
	}
	connection := db.C("sessions")
	err = connection.Insert(session)
	HandleError(err, r)
	r.JSON(200, session)
}

func GetSession(r render.Render, db *mgo.Database, params martini.Params) {
	collection := db.C("sessions")
	id := params["id"]
	session := &TranscodingSession{}
	collection.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(session)
	r.JSON(200, session)
}

func HandleError(err error, r render.Render) {
	if err != nil {
		r.JSON(400, map[string]string{
			"error": err.Error(),
		})
	}
}

func generate_id() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
