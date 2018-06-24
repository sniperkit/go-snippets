package main

import (
	"encoding/json"
	"fmt"
	"github.com/datatogether/core"
)

// ClientReqActions is a list of all actions a client may request
var ClientReqActions = []ClientAction{
	CreateUserAct{},
	SaveUserAct{},
	SessionLoginAct{},
	SessionLogoutAct{},
	SessionKeysAct{},
	MsgReqAct{},
	SearchReqAct{},
	FetchUrlAct{},
	FetchCollectionsAction{},
	FetchInboundLinksAct{},
	FetchOutboundLinksAct{},
	FetchContentUrlsAction{},
	FetchMetadataAction{},
	SaveMetadataAction{},
	FetchPrimersAction{},
	FetchPrimerAction{},
	FetchSourcesAction{},
	FetchSourceAction{},
	FetchSourceUrlsAction{},
	FetchSourceAttributedUrlsAction{},
	FetchConsensusAction{},
	FetchCollectionAction{},
	UserCollectionsAction{},
	FetchCollectionsAction{},
	SaveCollectionAction{},
	DeleteCollectionAction{},
	MetadataByKeyRequest{},
	FetchRecentContentUrlsAction{},
	TasksRequestAct{},
	TaskEnqueueAct{},
	CollectionItemsAction{},
	SaveCollectionItemsAction{},
	DeleteCollectionItemsAction{},
}

// Action is a collection of typed events for exchange between client & server
type Action interface {
	Type() string
}

type ClientAction interface {
	Action
	Parse(string, json.RawMessage) ClientRequestAction
}

type ClientRequestAction interface {
	Action
	SuccessType() string
	FailureType() string
	Exec() *ClientResponse
}

// ServerRequestAction is an action from the server to send to the client
type ServerRequestAction interface {
	Action
	Send()
}

type ClientResponse struct {
	Type        string      `json:"type"`
	RequestId   string      `json:"requestId"`
	Error       string      `json:"error,omitempty"`
	SilentError bool        `json:"silentError,omitempty"`
	Message     string      `json:"message,omitempty"`
	Schema      string      `json:"schema,omitempty"`
	Page        int         `json:"page,omitempty"`
	PageSize    int         `json:"pageSize,omitempty"`
	Id          string      `json:"id,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}

type ReqAction struct {
	RequestId   string `json:"requestId"`
	SilentError bool   `json:"silentError"`
	err         error
}

type MsgReqAct struct {
	ReqAction
	Message string
}

func (m MsgReqAct) Type() string        { return "MESSAGE_REQUEST" }
func (m MsgReqAct) SuccessType() string { return "MESSAGE_SUCCESS" }
func (m MsgReqAct) FailureType() string { return "MESSAGE_FAILURE" }

func (MsgReqAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &MsgReqAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}
func (a *MsgReqAct) Exec() (res *ClientResponse) {
	return &ClientResponse{
		Type:    a.SuccessType(),
		Message: fmt.Sprintf("oh really? %s", a.Message),
		Data: map[string]string{
			"message": fmt.Sprintf("oh really? %s", a.Message),
		},
	}
}

type SearchReqAct struct {
	ReqAction
	Query    string
	Page     int
	PageSize int
}

func (SearchReqAct) Type() string        { return "SEARCH_REQUEST" }
func (SearchReqAct) SuccessType() string { return "SEARCH_SUCCESS" }
func (SearchReqAct) FailureType() string { return "SEARCH_FAILURE" }

func (SearchReqAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &SearchReqAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}
func (s *SearchReqAct) Exec() (res *ClientResponse) {
	if s.Page > 0 {
		s.Page = s.Page - 1
	}
	results, err := core.Search(appDB, s.Query, s.PageSize, s.Page*s.PageSize)
	if err != nil {
		return &ClientResponse{
			Type:      s.FailureType(),
			Error:     err.Error(),
			RequestId: s.RequestId,
		}
	}
	return &ClientResponse{
		Type:      s.SuccessType(),
		RequestId: s.RequestId,
		Schema:    "SEARCH_RESULT_ARRAY",
		Data:      results,
	}
}

// FetchUrlAct fetches a url from the DB
type FetchUrlAct struct {
	ReqAction
	Url string
}

func (FetchUrlAct) Type() string        { return "URL_FETCH_REQUEST" }
func (FetchUrlAct) SuccessType() string { return "URL_FETCH_SUCCESS" }
func (FetchUrlAct) FailureType() string { return "URL_FETCH_FAILURE" }

func (FetchUrlAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchUrlAct{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchUrlAct) Exec() (res *ClientResponse) {
	u := &core.Url{Url: a.Url}
	if err := u.Read(store); err != nil {
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		Schema:    "URL",
		RequestId: a.RequestId,
		Data:      u,
	}
}

// FetchInboundLinksAct fetches a url's outbound links
type FetchInboundLinksAct struct {
	ReqAction
	Url string
}

func (FetchInboundLinksAct) Type() string        { return "URL_FETCH_INBOUND_LINKS_REQUEST" }
func (FetchInboundLinksAct) SuccessType() string { return "URL_FETCH_INBOUND_LINKS_SUCCESS" }
func (FetchInboundLinksAct) FailureType() string { return "URL_FETCH_INBOUND_LINKS_FAILURE" }

func (FetchInboundLinksAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchInboundLinksAct{}
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchInboundLinksAct) Exec() (res *ClientResponse) {
	links, err := core.ReadSrcLinks(appDB, &core.Url{Url: a.Url})
	if err != nil {
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "LINK_ARRAY",
		Data:      links,
	}
}

// FetchOutboundLinksAct fetches a url's outbound links
type FetchOutboundLinksAct struct {
	ReqAction
	Url string
}

func (FetchOutboundLinksAct) Type() string        { return "URL_FETCH_OUTBOUND_LINKS_REQUEST" }
func (FetchOutboundLinksAct) SuccessType() string { return "URL_FETCH_OUTBOUND_LINKS_SUCCESS" }
func (FetchOutboundLinksAct) FailureType() string { return "URL_FETCH_OUTBOUND_LINKS_FAILURE" }

func (FetchOutboundLinksAct) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchOutboundLinksAct{}
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchOutboundLinksAct) Exec() (res *ClientResponse) {
	links, err := core.ReadDstLinks(appDB, &core.Url{Url: a.Url})
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "LINK_ARRAY",
		Data:      links,
	}
}

// FetchRecentContentUrlsAction grabs a page of recently getted (no, "getted" is not a word)
// urls that lead to content
type FetchRecentContentUrlsAction struct {
	ReqAction
	Page     int
	PageSize int
}

func (FetchRecentContentUrlsAction) Type() string        { return "CONTENT_RECENT_URLS_REQUEST" }
func (FetchRecentContentUrlsAction) SuccessType() string { return "CONTENT_RECENT_URLS_SUCCESS" }
func (FetchRecentContentUrlsAction) FailureType() string { return "CONTENT_RECENT_URLS_FAILURE" }

func (FetchRecentContentUrlsAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchRecentContentUrlsAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchRecentContentUrlsAction) Exec() (res *ClientResponse) {
	urls, err := core.ContentUrls(appDB, a.PageSize, a.PageSize*(a.Page-1))
	if err != nil {
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "URL_ARRAY",
		Page:      a.Page,
		PageSize:  a.PageSize,
		Data:      urls,
	}
}

// FetchContentUrlsAction triggers archiving a url
type FetchContentUrlsAction struct {
	ReqAction
	Hash string
}

func (FetchContentUrlsAction) Type() string        { return "CONTENT_URLS_REQUEST" }
func (FetchContentUrlsAction) SuccessType() string { return "CONTENT_URLS_SUCCESS" }
func (FetchContentUrlsAction) FailureType() string { return "CONTENT_URLS_FAILURE" }

func (FetchContentUrlsAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchContentUrlsAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchContentUrlsAction) Exec() (res *ClientResponse) {
	urls, err := core.UrlsForHash(appDB, a.Hash)
	if err != nil {
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}
	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "URL_ARRAY",
		Data:      urls,
	}
}

// FetchMetadataAction triggers archiving a url
type FetchMetadataAction struct {
	ReqAction
	KeyId   string `json:"keyId"`
	Subject string `json:"subject"`
}

func (FetchMetadataAction) Type() string        { return "METADATA_REQUEST" }
func (FetchMetadataAction) SuccessType() string { return "METADATA_SUCCESS" }
func (FetchMetadataAction) FailureType() string { return "METADATA_FAILURE" }

func (FetchMetadataAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchMetadataAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchMetadataAction) Exec() (res *ClientResponse) {
	m, err := core.LatestMetadata(appDB, a.KeyId, a.Subject)
	if err != nil {
		if err == core.ErrNotFound {
			return &ClientResponse{
				Type:      a.SuccessType(),
				RequestId: a.RequestId,
			}
		}

		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "METADATA",
		Data:      m,
	}
}

// SaveMetadataAction triggers archiving a url
type SaveMetadataAction struct {
	ReqAction
	KeyId   string                 `json:"keyId"`
	Subject string                 `json:"subject"`
	Meta    map[string]interface{} `json:"meta"`
}

func (SaveMetadataAction) Type() string        { return "METADATA_SAVE_REQUEST" }
func (SaveMetadataAction) SuccessType() string { return "METADATA_SAVE_SUCCESS" }
func (SaveMetadataAction) FailureType() string { return "METADATA_SAVE_FAILURE" }

func (SaveMetadataAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &SaveMetadataAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *SaveMetadataAction) Exec() (res *ClientResponse) {
	m, err := core.NextMetadata(appDB, a.KeyId, a.Subject)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	m.Meta = a.Meta
	if err := m.Write(store); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "METADATA",
		Data:      a,
	}
}

// FetchPrimersAction grabs a page of primers
type FetchPrimersAction struct {
	ReqAction
	BaseOnly bool
	Page     int
	PageSize int
}

func (FetchPrimersAction) Type() string        { return "PRIMERS_FETCH_REQUEST" }
func (FetchPrimersAction) SuccessType() string { return "PRIMERS_FETCH_SUCCESS" }
func (FetchPrimersAction) FailureType() string { return "PRIMERS_FETCH_FAILURE" }

func (FetchPrimersAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchPrimersAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchPrimersAction) Exec() (res *ClientResponse) {
	var (
		primers []*core.Primer
		err     error
	)
	if a.BaseOnly {
		primers, err = core.BasePrimers(appDB, 50, 0)
	} else {
		primers, err = core.ListPrimers(store, 50, 0)
	}
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}
	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "PRIMER_ARRAY",
		Data:      primers,
	}
}

// FetchPrimerAction grabs a page of primers
type FetchPrimerAction struct {
	ReqAction
	Id string
}

func (FetchPrimerAction) Type() string        { return "PRIMER_FETCH_REQUEST" }
func (FetchPrimerAction) SuccessType() string { return "PRIMER_FETCH_SUCCESS" }
func (FetchPrimerAction) FailureType() string { return "PRIMER_FETCH_FAILURE" }

func (FetchPrimerAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchPrimerAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchPrimerAction) Exec() (res *ClientResponse) {
	p := &core.Primer{Id: a.Id}
	if err := p.Read(store); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	if err := p.ReadSubPrimers(appDB); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	if err := p.ReadSources(appDB); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "PRIMER",
		Data:      p,
	}
}

// FetchSourcesAction grabs a page of primers
type FetchSourcesAction struct {
	ReqAction
	Page     int
	PageSize int
}

func (FetchSourcesAction) Type() string        { return "SOURCES_FETCH_REQUEST" }
func (FetchSourcesAction) SuccessType() string { return "SOURCES_FETCH_SUCCESS" }
func (FetchSourcesAction) FailureType() string { return "SOURCES_FETCH_FAILURE" }

func (FetchSourcesAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchSourcesAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchSourcesAction) Exec() (res *ClientResponse) {
	s, err := core.ListSources(store, a.PageSize, (a.Page-1)*a.PageSize)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "SOURCE_ARRAY",
		Data:      s,
		Page:      a.Page,
		PageSize:  a.PageSize,
	}
}

// FetchSourceAction grabs a page of subprimers for a given primer id
type FetchSourceAction struct {
	ReqAction
	Id string
}

func (FetchSourceAction) Type() string        { return "SOURCE_FETCH_REQUEST" }
func (FetchSourceAction) SuccessType() string { return "SOURCE_FETCH_SUCCESS" }
func (FetchSourceAction) FailureType() string { return "SOURCE_FETCH_FAILURE" }

func (FetchSourceAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchSourceAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchSourceAction) Exec() (res *ClientResponse) {
	s := &core.Source{Id: a.Id}
	if err := s.Read(store); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	if err := s.Primer.Read(store); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	// TODO - hook this up to a cron-based que
	go func() {
		if err := s.CalcStats(appDB); err != nil {
			log.Info(err.Error())
			// return &ClientResponse{
			// 	Type:      a.FailureType(),
			// 	RequestId: a.RequestId,
			// 	Error:     err.Error(),
			// }
		}
	}()

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "SOURCE",
		Data:      s,
	}
}

// FetchSourceAction grabs a page of primers
type FetchSourceUrlsAction struct {
	ReqAction
	Id       string
	Page     int
	PageSize int
}

func (FetchSourceUrlsAction) Type() string        { return "SOURCE_URLS_REQUEST" }
func (FetchSourceUrlsAction) SuccessType() string { return "SOURCE_URLS_SUCCESS" }
func (FetchSourceUrlsAction) FailureType() string { return "SOURCE_URLS_FAILURE" }

func (FetchSourceUrlsAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchSourceUrlsAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchSourceUrlsAction) Exec() (res *ClientResponse) {
	s := &core.Source{Id: a.Id}
	if err := s.Read(store); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	urls, err := s.UndescribedContent(appDB, 100, 0)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "URL_ARRAY",
		Id:        a.Id,
		Page:      a.Page,
		PageSize:  a.PageSize,
		Data:      urls,
	}
}

type FetchSourceAttributedUrlsAction struct {
	ReqAction
	Id       string
	Page     int
	PageSize int
}

func (FetchSourceAttributedUrlsAction) Type() string {
	return "SOURCE_ATTRIBUTED_URLS_REQUEST"
}
func (FetchSourceAttributedUrlsAction) SuccessType() string {
	return "SOURCE_ATTRIBUTED_URLS_SUCCESS"
}
func (FetchSourceAttributedUrlsAction) FailureType() string {
	return "SOURCE_ATTRIBUTED_URLS_FAILURE"
}

func (FetchSourceAttributedUrlsAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchSourceAttributedUrlsAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchSourceAttributedUrlsAction) Exec() (res *ClientResponse) {
	s := &core.Source{Id: a.Id}
	if err := s.Read(store); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	urls, err := s.DescribedContent(appDB, 100, 0)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "URL_ARRAY",
		Id:        a.Id,
		Page:      a.Page,
		PageSize:  a.PageSize,
		Data:      urls,
	}
}

// FetchConsensusAction fetches a url from the DB
type FetchConsensusAction struct {
	ReqAction
	Subject string
}

func (FetchConsensusAction) Type() string        { return "CONSENSUS_REQUEST" }
func (FetchConsensusAction) SuccessType() string { return "CONSENSUS_SUCCESS" }
func (FetchConsensusAction) FailureType() string { return "CONSENSUS_FAILURE" }

func (FetchConsensusAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchConsensusAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchConsensusAction) Exec() (res *ClientResponse) {
	blocks, err := core.MetadataBySubject(appDB, a.Subject)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	c, values, err := core.SumConsensus(a.Subject, blocks)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	md, err := c.Metadata(values)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		Schema:    "CONSENSUS",
		RequestId: a.RequestId,
		Data: map[string]interface{}{
			"subject": a.Subject,
			"data":    md,
		},
	}
}

// FetchCollectionsAction grabs a page of collections
type FetchCollectionsAction struct {
	ReqAction
	Page     int
	PageSize int
}

func (FetchCollectionsAction) Type() string        { return "COLLECTIONS_FETCH_REQUEST" }
func (FetchCollectionsAction) SuccessType() string { return "COLLECTIONS_FETCH_SUCCESS" }
func (FetchCollectionsAction) FailureType() string { return "COLLECTIONS_FETCH_FAILURE" }

func (FetchCollectionsAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchCollectionsAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchCollectionsAction) Exec() (res *ClientResponse) {
	collections, err := core.ListCollections(store, 50, 0)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "COLLECTION_ARRAY",
		Data:      collections,
	}
}

// UserCollectionsAction grabs a page of a user's collections
type UserCollectionsAction struct {
	ReqAction
	Creator  string
	Page     int
	PageSize int
}

func (UserCollectionsAction) Type() string        { return "USER_COLLECTIONS_REQUEST" }
func (UserCollectionsAction) SuccessType() string { return "USER_COLLECTIONS_SUCCESS" }
func (UserCollectionsAction) FailureType() string { return "USER_COLLECTIONS_FAILURE" }

func (UserCollectionsAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &UserCollectionsAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *UserCollectionsAction) Exec() (res *ClientResponse) {
	collections, err := core.CollectionsByCreator(store, a.Creator, "created DESC", a.PageSize, (a.Page-1)*a.PageSize)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "COLLECTION_ARRAY",
		Data:      collections,
	}
}

// FetchCollectionAction grabs a page of collections
type FetchCollectionAction struct {
	ReqAction
	Id string
}

func (FetchCollectionAction) Type() string        { return "COLLECTION_FETCH_REQUEST" }
func (FetchCollectionAction) SuccessType() string { return "COLLECTION_FETCH_SUCCESS" }
func (FetchCollectionAction) FailureType() string { return "COLLECTION_FETCH_FAILURE" }

func (FetchCollectionAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &FetchCollectionAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *FetchCollectionAction) Exec() (res *ClientResponse) {
	c := &core.Collection{Id: a.Id}
	if err := c.Read(store); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "COLLECTION",
		Data:      c,
	}
}

// SaveCollectionAction triggers archiving a url
type SaveCollectionAction struct {
	ReqAction
	Collection *core.Collection `json:"collection"`
}

func (SaveCollectionAction) Type() string        { return "COLLECTION_SAVE_REQUEST" }
func (SaveCollectionAction) SuccessType() string { return "COLLECTION_SAVE_SUCCESS" }
func (SaveCollectionAction) FailureType() string { return "COLLECTION_SAVE_FAILURE" }

func (SaveCollectionAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &SaveCollectionAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *SaveCollectionAction) Exec() (res *ClientResponse) {
	log.Info(a.Collection)
	if err := a.Collection.Save(store); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "COLLECTION",
		Data:      a.Collection,
	}
}

// DeleteCollectionAction triggers archiving a url
type DeleteCollectionAction struct {
	ReqAction
	// Collection *core.Collection `json:"collection"`
	Id string `json:"id"`
}

func (DeleteCollectionAction) Type() string        { return "COLLECTION_DELETE_REQUEST" }
func (DeleteCollectionAction) SuccessType() string { return "COLLECTION_DELETE_SUCCESS" }
func (DeleteCollectionAction) FailureType() string { return "COLLECTION_DELETE_FAILURE" }

func (DeleteCollectionAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &DeleteCollectionAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *DeleteCollectionAction) Exec() (res *ClientResponse) {
	c := &core.Collection{Id: a.Id}
	if err := c.Delete(store); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "COLLECTION",
		Data:      c,
	}
}

// CollectionItemsAction grabs a page of collection items
type CollectionItemsAction struct {
	ReqAction
	CollectionId string
	Page         int
	PageSize     int
}

func (CollectionItemsAction) Type() string        { return "COLLECTION_ITEMS_REQUEST" }
func (CollectionItemsAction) SuccessType() string { return "COLLECTION_ITEMS_SUCCESS" }
func (CollectionItemsAction) FailureType() string { return "COLLECTION_ITEMS_FAILURE" }

func (CollectionItemsAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &CollectionItemsAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *CollectionItemsAction) Exec() (res *ClientResponse) {
	c := core.Collection{Id: a.CollectionId}

	items, err := c.ReadItems(store, "created DESC", a.PageSize, (a.Page-1)*a.PageSize)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "COLLECTION_ITEM_ARRAY",
		Data:      items,
		Id:        a.CollectionId,
		Page:      a.Page,
		PageSize:  a.PageSize,
	}
}

// SaveCollectionItemsAction grabs a page of collection items
type SaveCollectionItemsAction struct {
	ReqAction
	CollectionId string
	Items        []*core.CollectionItem
}

func (SaveCollectionItemsAction) Type() string        { return "COLLECTION_SAVE_ITEMS_REQUEST" }
func (SaveCollectionItemsAction) SuccessType() string { return "COLLECTION_SAVE_ITEMS_SUCCESS" }
func (SaveCollectionItemsAction) FailureType() string { return "COLLECTION_SAVE_ITEMS_FAILURE" }

func (SaveCollectionItemsAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &SaveCollectionItemsAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *SaveCollectionItemsAction) Exec() (res *ClientResponse) {
	c := core.Collection{Id: a.CollectionId}
	if err := c.SaveItems(store, a.Items); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "COLLECTION_ITEM_ARRAY",
		Data:      a.Items,
		// need to provide collection id & pagination info to
		// pull items out of pagination client-side
		Id:       c.Id,
		Page:     1,
		PageSize: len(a.Items),
	}
}

// DeleteCollectionItemsAction grabs a page of collection items
type DeleteCollectionItemsAction struct {
	ReqAction
	CollectionId string
	Items        []*core.CollectionItem
}

func (DeleteCollectionItemsAction) Type() string        { return "COLLECTION_DELETE_ITEMS_REQUEST" }
func (DeleteCollectionItemsAction) SuccessType() string { return "COLLECTION_DELETE_ITEMS_SUCCESS" }
func (DeleteCollectionItemsAction) FailureType() string { return "COLLECTION_DELETE_ITEMS_FAILURE" }

func (DeleteCollectionItemsAction) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &DeleteCollectionItemsAction{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *DeleteCollectionItemsAction) Exec() (res *ClientResponse) {
	c := core.Collection{Id: a.CollectionId}
	if err := c.DeleteItems(store, a.Items); err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "COLLECTION_ITEM_ARRAY",
		Data:      a.Items,
		// need to provide collection id & pagination info to
		// pull items out of pagination client-side
		Id:       c.Id,
		Page:     1,
		PageSize: len(a.Items),
	}
}

// MetadataByKeyRequest triggers archiving a url
type MetadataByKeyRequest struct {
	ReqAction
	Key      string `json:"key"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

func (MetadataByKeyRequest) Type() string        { return "METADATA_BY_KEY_REQUEST" }
func (MetadataByKeyRequest) SuccessType() string { return "METADATA_BY_KEY_SUCCESS" }
func (MetadataByKeyRequest) FailureType() string { return "METADATA_BY_KEY_FAILURE" }

func (MetadataByKeyRequest) Parse(reqId string, data json.RawMessage) ClientRequestAction {
	a := &MetadataByKeyRequest{}
	a.RequestId = reqId
	a.err = json.Unmarshal(data, a)
	return a
}

func (a *MetadataByKeyRequest) Exec() (res *ClientResponse) {
	results, err := core.MetadataByKey(appDB, a.Key, a.PageSize, (a.Page-1)*a.PageSize)
	if err != nil {
		log.Info(err.Error())
		return &ClientResponse{
			Type:      a.FailureType(),
			RequestId: a.RequestId,
			Error:     err.Error(),
		}
	}

	return &ClientResponse{
		Type:      a.SuccessType(),
		RequestId: a.RequestId,
		Schema:    "METADATA_ARRAY",
		Data:      results,
	}
}
