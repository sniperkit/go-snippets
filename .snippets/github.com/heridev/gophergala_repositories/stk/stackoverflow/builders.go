package stackoverflow

import (
	"github.com/lann/builder"
)

// Search Request Builder
type searchReqBuilder builder.Builder

func (b searchReqBuilder) Query(v string) searchReqBuilder {
	return builder.Set(b, "Query", v).(searchReqBuilder)
}

func (b searchReqBuilder) AddTag(v string) searchReqBuilder {
	return builder.Append(b, "Tags", v).(searchReqBuilder)
}

func (b searchReqBuilder) SiteID(v string) searchReqBuilder {
	return builder.Set(b, "SiteID", v).(searchReqBuilder)
}

func (b searchReqBuilder) Sort(v string) searchReqBuilder {
	return builder.Set(b, "Sort", v).(searchReqBuilder)
}

func (b searchReqBuilder) Accepted(flag bool) searchReqBuilder {
	return builder.Set(b, "Accepted", flag).(searchReqBuilder)
}

func (b searchReqBuilder) Tags(tags []string) searchReqBuilder {
	return builder.Set(b, "Tags", tags).(searchReqBuilder)
}

func (b searchReqBuilder) Build() SearchRequest {
	return builder.GetStruct(b).(SearchRequest)
}

// Answer Request Builder
type answersReqBuilder builder.Builder

func (b answersReqBuilder) AddAnswerID(v int) answersReqBuilder {
	return builder.Append(b, "AnswerIDS", v).(answersReqBuilder)
}

func (b answersReqBuilder) SiteID(v string) answersReqBuilder {
	return builder.Set(b, "SiteID", v).(answersReqBuilder)
}

func (b answersReqBuilder) Build() AnswerRequest {
	return builder.GetStruct(b).(AnswerRequest)
}

// Register Builders
var SearchRequestBuilder = builder.Register(searchReqBuilder{}, SearchRequest{}).(searchReqBuilder)
var AnswerRequestBuilder = builder.Register(answersReqBuilder{}, AnswerRequest{}).(answersReqBuilder)
