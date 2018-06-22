package transcoder_test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	. "github.com/gophergala/abbita/transcoder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var dbName string
	var session *DatabaseSession
	var server Server
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	BeforeEach(func() {
		dbName = "transcoder_test"
		session = NewSession(dbName)
		server = NewServer(session)

		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		session.DB(dbName).DropDatabase()
	})

	Describe("GET /", func() {
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/", nil)
		})
		It("returns a status code of 200", func() {
			server.ServeHTTP(recorder, request)
			Expect(recorder.Code).To(Equal(200))
		})
	})

	Describe("POST /sessions/new", func() {
		BeforeEach(func() {
			os.Mkdir("media", 0777)
		})
		AfterEach(func() {
			os.RemoveAll("media")
		})
		Context("with file and name", func() {
			BeforeEach(func() {
				extraParams := map[string]string{
					"name": "mpthreetest.mp3",
				}
				multipartPayload, mpWriter := makeMultiPartPayload("../test_files/mpthreetest.mp3", "fileUpload", extraParams)
				var err error
				request, err = http.NewRequest("POST", "/api/sessions/new", multipartPayload)
				if err != nil {
					panic(err)
				}

				request.Header.Add("Content-Type", mpWriter.FormDataContentType())

				err = mpWriter.Close()
				if err != nil {
					panic(err)
				}
			})
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns the session as json", func() {
				server.ServeHTTP(recorder, request)
				sessionJson := mapFromJSON(recorder.Body.Bytes())
				Expect(sessionJson["name"]).To(Equal("mpthreetest.mp3"))
			})
			It("saves the session to the db", func() {
				db := session.DB(dbName)
				connection := db.C("sessions")
				initialCount, _ := connection.Find(nil).Count()
				server.ServeHTTP(recorder, request)
				finalCount, _ := connection.Find(nil).Count()
				Expect(finalCount).To(Equal(initialCount + 1))
			})

		})

	})

})

func makeMultiPartPayload(path string, paramName string, params map[string]string) (*bytes.Buffer, *multipart.Writer) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(part, file)
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		panic(err)
	}
	return body, writer
}

func mapFromJSON(data []byte) map[string]interface{} {
	var result interface{}
	json.Unmarshal(data, &result)
	return result.(map[string]interface{})
}
