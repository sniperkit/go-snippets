package sls

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strconv"

	lz4 "github.com/cloudflare/golz4"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
)

// LogStore defines LogStore struct
type LogStore struct {
	Name       string `json:"logstoreName"`
	TTL        int
	ShardCount int

	CreateTime     uint32
	LastModifyTime uint32

	project *LogProject
}

// Shard defines shard struct
type Shard struct {
	ShardID int `json:"shardID"`
}

// ListShards returns shard id list of this logstore.
func (s *LogStore) ListShards() (shardIDs []int, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
	}
	uri := fmt.Sprintf("/logstores/%v/shards", s.Name)
	r, err := request(s.project, "GET", uri, h, nil)
	if err != nil {
		return nil, NewClientError(err.Error())
	}
	buf, _ := ioutil.ReadAll(r.Body)
	if r.StatusCode != http.StatusOK {
		err := &Error{}
		json.Unmarshal(buf, err)
		return nil, err
	}

	var shards []*Shard
	json.Unmarshal(buf, &shards)
	for _, v := range shards {
		shardIDs = append(shardIDs, v.ShardID)
	}
	return shardIDs, nil
}

// PutLogs put logs into logstore.
// The callers should transform user logs into LogGroup.
func (s *LogStore) PutLogs(lg *LogGroup) (err error) {
	if len(lg.Logs) == 0 {
		// empty log group
		return nil
	}

	body, err := proto.Marshal(lg)
	if err != nil {
		return NewClientError(err.Error())
	}

	// Compresse body with lz4
	// out := make([]byte, lz4.CompressBound(body))
	// n, err := lz4.Compress(body, out)
	// if err != nil {
	// 	return NewClientError(err.Error())
	// }

	h := map[string]string{
		// "x-log-compresstype": "lz4",
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/x-protobuf",
	}

	uri := fmt.Sprintf("/logstores/%v", s.Name)
	r, err := request(s.project, "POST", uri, h, body)
	if err != nil {
		return NewClientError(err.Error())
	}

	body, _ = ioutil.ReadAll(r.Body)
	if r.StatusCode != http.StatusOK {
		err := new(Error)
		json.Unmarshal(body, err)
		return err
	}
	return nil
}

// GetCursor gets log cursor of one shard specified by shardId.
// The from can be in three form: a) unix timestamp in seccond, b) "begin", c) "end".
// For more detail please read: http://gitlab.alibaba-inc.com/sls/doc/blob/master/api/shard.md#logstore
func (s *LogStore) GetCursor(shardID int, from string) (cursor string, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
	}
	uri := fmt.Sprintf("/logstores/%v/shards/%v?type=cursor&from=%v",
		s.Name, shardID, from)
	r, err := request(s.project, "GET", uri, h, nil)
	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	if r.StatusCode != http.StatusOK {
		errMsg := &Error{}
		err = json.Unmarshal(buf, errMsg)
		if err != nil {
			err = fmt.Errorf("failed to get cursor")
			dump, _ := httputil.DumpResponse(r, true)
			if glog.V(1) {
				glog.Error(string(dump))
			}
			return
		}
		err = fmt.Errorf("%v:%v", errMsg.Code, errMsg.Message)
		return
	}

	type Body struct {
		Cursor string
	}
	body := &Body{}

	err = json.Unmarshal(buf, body)
	if err != nil {
		return
	}
	cursor = body.Cursor
	return
}

// GetLogsBytes gets logs binary data from shard specified by shardId according cursor and endCursor.
// The logGroupMaxCount is the max number of logGroup could be returned.
// The nextCursor is the next curosr can be used to read logs at next time.
func (s *LogStore) GetLogsBytes(shardID int, cursor, endCursor string,
	logGroupMaxCount int) (out []byte, nextCursor string, err error) {
	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Accept":            "application/x-protobuf",
		"Accept-Encoding":   "lz4",
	}

	uri := ""
	if endCursor == "" {
		uri = fmt.Sprintf("/logstores/%v/shards/%v?type=logs&cursor=%v&count=%v",
			s.Name, shardID, cursor, logGroupMaxCount)
	} else {
		uri = fmt.Sprintf("/logstores/%v/shards/%v?type=logs&cursor=%v&end_cursor=%v&count=%v",
			s.Name, shardID, cursor, endCursor, logGroupMaxCount)
	}

	r, err := request(s.project, "GET", uri, h, nil)
	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	if r.StatusCode != http.StatusOK {
		errMsg := &Error{}
		err = json.Unmarshal(buf, errMsg)
		if err != nil {
			err = fmt.Errorf("failed to get cursor")
			dump, _ := httputil.DumpResponse(r, true)
			if glog.V(1) {
				glog.Error(string(dump))
			}
			return
		}
		err = fmt.Errorf("%v:%v", errMsg.Code, errMsg.Message)
		return
	}
	v, ok := r.Header["X-Log-Compresstype"]
	if !ok || len(v) == 0 {
		err = fmt.Errorf("can't find 'x-log-compresstype' header")
		return
	}
	if v[0] != "lz4" {
		err = fmt.Errorf("unexpected compress type:%v", v[0])
		return
	}

	v, ok = r.Header["X-Log-Cursor"]
	if !ok || len(v) == 0 {
		err = fmt.Errorf("can't find 'x-log-cursor' header")
		return
	}
	nextCursor = v[0]

	v, ok = r.Header["X-Log-Bodyrawsize"]
	if !ok || len(v) == 0 {
		err = fmt.Errorf("can't find 'x-log-bodyrawsize' header")
		return
	}
	bodyRawSize, err := strconv.Atoi(v[0])
	if err != nil {
		return
	}

	out = make([]byte, bodyRawSize)
	err = lz4.Uncompress(buf, out)
	if err != nil {
		return
	}

	return
}

// LogsBytesDecode decodes logs binary data returned by GetLogsBytes API
func LogsBytesDecode(data []byte) (gl *LogGroupList, err error) {

	gl = &LogGroupList{}
	err = proto.Unmarshal(data, gl)
	if err != nil {
		return nil, err
	}

	return gl, nil
}

// PullLogs gets logs from shard specified by shardId according cursor and endCursor.
// The logGroupMaxCount is the max number of logGroup could be returned.
// The nextCursor is the next cursor can be used to read logs at next time.
func (s *LogStore) PullLogs(shardID int, cursor, endCursor string,
	logGroupMaxCount int) (gl *LogGroupList, nextCursor string, err error) {

	out, nextCursor, err := s.GetLogsBytes(shardID, cursor, endCursor, logGroupMaxCount)
	if err != nil {
		return nil, "", err
	}

	gl, err = LogsBytesDecode(out)
	if err != nil {
		return nil, "", err
	}

	return gl, nextCursor, nil
}

// GetHistograms query logs with [from, to) time range
func (s *LogStore) GetHistograms(topic string, from int64, to int64, queryExp string) (*GetHistogramsResponse, error) {

	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Accept":            "application/json",
	}

	uri := fmt.Sprintf("/logstores/%v?type=histogram&topic=%v&from=%v&to=%v&query=%v", s.Name, topic, from, to, queryExp)

	r, err := request(s.project, "GET", uri, h, nil)
	if err != nil {
		return nil, NewClientError(err.Error())
	}

	body, _ := ioutil.ReadAll(r.Body)
	if r.StatusCode != http.StatusOK {
		err := new(Error)
		json.Unmarshal(body, err)
		return nil, err
	}

	histograms := []SingleHistogram{}
	err = json.Unmarshal(body, &histograms)
	if err != nil {
		return nil, err
	}

	count, err := strconv.ParseInt(r.Header[GetLogsCountHeader][0], 10, 32)
	if err != nil {
		return nil, err
	}
	getHistogramsResponse := GetHistogramsResponse{
		Progress:   r.Header[ProgressHeader][0],
		Count:      count,
		Histograms: histograms,
	}

	return &getHistogramsResponse, nil
}

// GetLogs query logs with [from, to) time range
func (s *LogStore) GetLogs(topic string, from int64, to int64, queryExp string,
	maxLineNum int64, offset int64, reverse bool) (*GetLogsResponse, error) {

	h := map[string]string{
		"x-log-bodyrawsize": "0",
		"Accept":            "application/json",
	}

	uri := fmt.Sprintf("/logstores/%v?type=log&topic=%v&from=%v&to=%v&query=%v&line=%v&offset=%v&reverse=%v", s.Name, topic, from, to, queryExp, maxLineNum, offset, reverse)

	r, err := request(s.project, "GET", uri, h, nil)
	if err != nil {
		return nil, NewClientError(err.Error())
	}

	body, _ := ioutil.ReadAll(r.Body)
	if r.StatusCode != http.StatusOK {
		err := new(Error)
		json.Unmarshal(body, err)
		return nil, err
	}

	logs := []map[string]string{}
	err = json.Unmarshal(body, &logs)
	if err != nil {
		return nil, err
	}

	count, err := strconv.ParseInt(r.Header[GetLogsCountHeader][0], 10, 32)
	if err != nil {
		return nil, err
	}

	getLogsResponse := GetLogsResponse{
		Progress: r.Header[ProgressHeader][0],
		Count:    count,
		Logs:     logs,
	}

	return &getLogsResponse, nil
}

// CreateIndex ...
func (s *LogStore) CreateIndex(index Index) error {
	body, err := json.Marshal(index)
	if err != nil {
		return err
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
		"Accept-Encoding":   "deflate", // TODO: support lz4
	}

	uri := fmt.Sprintf("/logstores/%s/index", s.Name)
	_, err = request(s.project, "POST", uri, h, body)
	return err
}

// UpdateIndex ...
func (s *LogStore) UpdateIndex(index Index) error {
	body, err := json.Marshal(index)
	if err != nil {
		return err
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
		"Accept-Encoding":   "deflate", // TODO: support lz4
	}

	uri := fmt.Sprintf("/logstores/%s/index", s.Name)
	_, err = request(s.project, "PUT", uri, h, body)
	return err
}

// DeleteIndex ...
func (s *LogStore) DeleteIndex() error {
	type Body struct {
		project string `json:"projectName"`
		store   string `json:"logstoreName"`
	}

	body, err := json.Marshal(Body{
		project: s.project.Name,
		store:   s.Name,
	})
	if err != nil {
		return err
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
		"Accept-Encoding":   "deflate", // TODO: support lz4
	}

	uri := fmt.Sprintf("/logstores/%s/index", s.Name)
	_, err = request(s.project, "DELETE", uri, h, body)
	return err
}

func (s *LogStore) GetIndex() (*Index, error) {
	type Body struct {
		project string `json:"projectName"`
		store   string `json:"logstoreName"`
	}

	body, err := json.Marshal(Body{
		project: s.project.Name,
		store:   s.Name,
	})
	if err != nil {
		return nil, err
	}

	h := map[string]string{
		"x-log-bodyrawsize": fmt.Sprintf("%v", len(body)),
		"Content-Type":      "application/json",
		"Accept-Encoding":   "deflate", // TODO: support lz4
	}

	uri := fmt.Sprintf("/logstores/%s/index", s.Name)
	resp, err := request(s.project, "GET", uri, h, body)
	if err != nil {
	}

	index := &Index{}
	data, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(data, index)
	if err != nil {
		return nil, err
	}

	return index, err
}
