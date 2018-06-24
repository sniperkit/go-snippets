package ingestion_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	"github.com/keelerh/omniscience/cmd/omniscience_server/ingestion"
	"github.com/keelerh/omniscience/cmd/omniscience_server/ingestion/mocks"
	pb "github.com/keelerh/omniscience/protos"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIngester_Ingest_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"_shards" : {
				"total" : 2,
				"failed" : 0,
				"successful" : 2
			},
			"_index" : "omniscience",
			"_type" : "_doc",
			"_id" : "123",
			"_version" : 1,
			"_seq_no" : 0,
			"_primary_term" : 1,
			"result" : "created"
		}`)
	}))
	defer ts.Close()

	elasticClient, err := mockElasticClient(ts.URL)
	require.NoError(t, err)

	ingester := ingestion.NewIngester(elasticClient)

	doc := &pb.Document{
		Id:          &pb.DocumentId{Id: "123"},
		Title:       "Fake document",
		Description: "",
		Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit," +
			"sed do eiusmod tempor incididunt ut labore et dolore magna aliqua." +
			"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris" +
			"nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in" +
			"reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla" +
			"nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt" +
			"in culpa qui officia deserunt mollit anim id est laborum.",
		Url:          "https://abc.xyz",
		Service:      "fake_service",
		LastModified: ptypes.TimestampNow(),
	}
	stream := ingester_mocks.NewMockIngester_IngestServer(mockCtrl)
	stream.EXPECT().Context().Return(context.TODO()).Times(1)
	stream.EXPECT().Recv().Return(doc, nil).Times(1)
	stream.EXPECT().Recv().Return(nil, io.EOF).Times(1)

	err = ingester.Ingest(stream)
	assert.NoError(t, err)
}

func TestIngester_Ingest_StreamError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	elasticClient, err := mockElasticClient(ts.URL)
	require.NoError(t, err)

	ingester := ingestion.NewIngester(elasticClient)

	stream := ingester_mocks.NewMockIngester_IngestServer(mockCtrl)
	stream.EXPECT().Recv().Return(nil, errors.New("stream failed")).Times(1)

	err = ingester.Ingest(stream)
	assert.EqualError(t, err, "stream failed")
}

func TestIngester_ElasticSearchIndexExists_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, `{}`)
	}))
	defer ts.Close()

	elasticClient, err := mockElasticClient(ts.URL)
	require.NoError(t, err)

	ingester := ingestion.NewIngester(elasticClient)

	stream := ingester_mocks.NewMockIngester_IngestServer(mockCtrl)
	stream.EXPECT().Context().Return(context.TODO()).Times(1)
	stream.EXPECT().Recv().Return(nil, nil).Times(1)

	err = ingester.Ingest(stream)
	assert.Error(t, err)
}

func mockElasticClient(url string) (*elastic.Client, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	return client, nil
}
