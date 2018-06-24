package lib

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	pb "github.com/keelerh/omniscience/protos"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
)

const (
	GoogleDriveWebViewLink = "https://docs.google.com/document/d/"
	service                = "gdrive"
)

type GoogleDriveService struct {
	cfg *jwt.Config
}

func NewGoogleDrive(srvAccountCfg *jwt.Config) *GoogleDriveService {
	return &GoogleDriveService{
		cfg: srvAccountCfg,
	}
}

func (g *GoogleDriveService) Fetch(modifiedSince time.Time) ([]*pb.Document, error) {
	var allDocuments []*pb.Document

	modifiedSinceProto, err := ptypes.TimestampProto(modifiedSince)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse modified since timestamp as proto")
	}

	ctx := context.Background()
	svc, err := drive.New(g.cfg.Client(ctx))
	if err != nil {
		return nil, errors.Wrap(err, "failed to instantiate new Google Drive service from config")
	}

	pageToken := ""
	for {
		q := svc.Files.List()
		// If we have a pageToken set, apply it to the query.
		if pageToken != "" {
			q = q.PageToken(pageToken)
		}
		r, err := q.Do()
		if err != nil {
			return nil, err
		}
		for _, f := range r.Files {
			// TODO: Figure out how to download all formats of text files.
			// Only attempt to download text files and gdocs.
			isGoogleDoc := isGoogleDoc(f.MimeType)
			if !(isTextMime(f.MimeType) || isGoogleDoc) {
				continue
			}
			content, err := downloadFile(svc, f.Id, isGoogleDoc)
			if err != nil {
				log.Warningf("unable to download file: FileId(%v) %v", f.Id, err)
				continue
			}
			// TODO: Only retrieve files modified after the last modified time specified in the request.
			allDocuments = append(allDocuments, &pb.Document{
				Id:          &pb.DocumentId{Id: f.Id},
				Title:       f.Name,
				Description: f.Description,
				Service:     service,
				Content:     content,
				Url:         GoogleDriveWebViewLink + f.Id,
				// TODO: This should be using the ModifiedTime of the file returned by Google.
				// but that field currently appears to be null.
				LastModified: modifiedSinceProto,
			})
		}
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return allDocuments, nil
}

func downloadFile(svc *drive.Service, fileId string, isGoogleDoc bool) (string, error) {
	var resp *http.Response
	var err error
	if isGoogleDoc {
		resp, err = svc.Files.Export(fileId, "text/plain").Download()
	} else {
		resp, err = svc.Files.Get(fileId).Download()
	}
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	content := standardizeSpaces(string(body))

	return content, nil
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func isTextMime(mimeType string) bool {
	return strings.HasPrefix("text", mimeType)
}

func isGoogleDoc(mimeType string) bool {
	return mimeType == "application/vnd.google-apps.document"
}
