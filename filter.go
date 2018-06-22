package insight

import (
	"fmt"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/homebot/core/utils"
	insightV1 "github.com/homebot/protobuf/pkg/api/insight/v1"
)

// Builder builds message filters and serves a filter chain
type Builder struct {
	chain []MessageFilter
}

// First returns the first MessageFilter in the chain
func (b *Builder) First() MessageFilter {
	if len(b.chain) == 0 {
		return nil
	}

	return b.chain[0]
}

// Match checks if the given message matches all filters in the chain
func (b *Builder) Match(msg *insightV1.LogMessage) bool {
	for _, m := range b.chain {
		if !m.Match(msg) {
			return false
		}
	}

	return true
}

// Len returns the number of entries in the filter chain
func (b *Builder) Len() int {
	return len(b.chain)
}

func (b *Builder) String() string {
	var res = ""

	for _, c := range b.chain {
		if res != "" {
			res += " and "
		}

		res += c.String()
	}

	return res
}

// WithService addes a filter based on the service name
func (b *Builder) WithService(s string) *Builder {
	b.chain = append(b.chain, ServiceFilter{
		Service: s,
	})
	return b
}

// WithServiceType adds a filter based on the service type
func (b *Builder) WithServiceType(t string) *Builder {
	b.chain = append(b.chain, ServiceTypeFilter{t})
	return b
}

// WithIdentity adds a filter based on the identity
func (b *Builder) WithIdentity(i string) *Builder {
	b.chain = append(b.chain, IdentityFilter{i})
	return b
}

// WithResource adds a filter based on the resource name
func (b *Builder) WithResource(r string) *Builder {
	b.chain = append(b.chain, ResourceFilter{r})
	return b
}

// WithSeverity adds a filter based on the message severity
func (b *Builder) WithSeverity(l insightV1.SeverityLevel) *Builder {
	b.chain = append(b.chain, SeverityFilter{l})
	return b
}

// WithNotBefore adds a filter that matches messages that happend NOT BEFORE t
func (b *Builder) WithNotBefore(t time.Time) *Builder {
	b.chain = append(b.chain, TimeFilter{
		NotBefore:  true,
		BeforeTime: t,
	})
	return b
}

// WithBefore adds a filter that matches messages that happend BEFORE t
func (b *Builder) WithBefore(t time.Time) *Builder {
	b.chain = append(b.chain, TimeFilter{
		BeforeTime: t,
	})
	return b
}

// WithNotAfter adds a filter that matches messages that happend NOT AFTER t
func (b *Builder) WithNotAfter(t time.Time) *Builder {
	b.chain = append(b.chain, TimeFilter{
		NotAfter:  true,
		AfterTime: t,
	})
	return b
}

// WithAfter adds a filter that matches messages that happend AFTER t
func (b *Builder) WithAfter(t time.Time) *Builder {
	b.chain = append(b.chain, TimeFilter{
		AfterTime: t,
	})
	return b
}

// WithDetails adds a filter that only matches messages that have equal
// key value pairs
func (b *Builder) WithDetails(values utils.ValueMap) *Builder {
	b.chain = append(b.chain, DetailsFilter{values})
	return b
}

// NewFilter creates a new filter builder
func NewFilter() *Builder {
	return &Builder{}
}

// ServiceFilter filters messages based on the service name
type ServiceFilter struct {
	Service string
}

// Match implements MessageFilter
func (sf ServiceFilter) Match(msg *insightV1.LogMessage) bool {
	return msg.GetService() == sf.Service
}

func (sf ServiceFilter) String() string {
	return fmt.Sprintf("Service == %q", sf.Service)
}

// ServiceTypeFilter filters messages based on the service type
type ServiceTypeFilter struct {
	ServiceType string
}

// Match implements MessageFilter
func (st ServiceTypeFilter) Match(msg *insightV1.LogMessage) bool {
	return msg.GetServiceType() == st.ServiceType
}

func (st ServiceTypeFilter) String() string {
	return fmt.Sprintf("ServiceType == %q", st.ServiceType)
}

// IdentityFilter filters messages based on the identity
type IdentityFilter struct {
	Identity string
}

// Match implements MessageFilter
func (i IdentityFilter) Match(msg *insightV1.LogMessage) bool {
	return msg.GetIdentity() == i.Identity
}

func (i IdentityFilter) String() string {
	return fmt.Sprintf("Identity == %q", i.Identity)
}

// ResourceFilter filters messages based on the resource name
type ResourceFilter struct {
	Resource string
}

// Match implements MessageFilter
func (i ResourceFilter) Match(msg *insightV1.LogMessage) bool {
	return msg.GetTargetResource() == i.Resource
}

func (i ResourceFilter) String() string {
	return fmt.Sprintf("Resource == %q", i.Resource)
}

// SeverityFilter filters messages based on the severity
type SeverityFilter struct {
	Level insightV1.SeverityLevel
}

// Match implements MessageFilter
func (i SeverityFilter) Match(msg *insightV1.LogMessage) bool {
	return msg.GetSeverity() == i.Level
}

func (i SeverityFilter) String() string {
	return fmt.Sprintf("Serverity == %q", i.Level.String())
}

// TimeFilter filters messages based on the message timestamp
type TimeFilter struct {
	NotBefore  bool
	BeforeTime time.Time

	NotAfter  bool
	AfterTime time.Time
}

// Match implements MessageFilter
func (m TimeFilter) Match(msg *insightV1.LogMessage) bool {
	mt, _ := ptypes.Timestamp(msg.GetCreatedTime())

	if !m.BeforeTime.Equal(time.Unix(0, 0)) {
		if m.NotBefore {
			if mt.Before(m.BeforeTime) {
				return false
			}
		} else {
			if mt.After(m.BeforeTime) {
				return false
			}
		}
	}

	if !m.AfterTime.Equal(time.Unix(0, 0)) {
		if m.NotAfter {
			if mt.After(m.AfterTime) {
				return false
			}
		} else {
			if mt.Before(m.AfterTime) {
				return false
			}
		}
	}

	return true
}

func (t TimeFilter) String() string {
	res := ""

	if !t.BeforeTime.Equal(time.Unix(0, 0)) {
		what := "<"
		if t.NotBefore {
			what = ">="
		}

		res = fmt.Sprintf("Time %s %q", what, t.BeforeTime.Format(time.RFC822))
	}

	if !t.AfterTime.Equal(time.Unix(0, 0)) {
		what := ">"
		if t.NotAfter {
			what = "<="
		}

		if res != "" {
			res += " and "
		}

		res += fmt.Sprintf("Time %s %q", what, t.AfterTime.Format(time.RFC822))
	}

	return res
}

// DetailsFilter filters messages based on optional message details
type DetailsFilter struct {
	Details utils.ValueMap
}

// Match implements MessageFilter
func (m DetailsFilter) Match(msg *insightV1.LogMessage) bool {
	details := utils.ValueMapFrom(msg.GetDetails())

	for key, value := range m.Details {
		v, ok := details[key]
		if !ok {
			return false
		}

		if !reflect.DeepEqual(value, v) {
			return false
		}
	}

	return true
}

func (m DetailsFilter) String() string {
	return "details"
}

// BuilderFromProto creates a Filter chain that applies search conditions on
// log messages based on a SearchMessageRequest
func BuilderFromProto(in *insightV1.SearchMessagesRequest) (*Builder, error) {
	b := NewFilter()

	if in.GetServiceType() != "" {
		b.WithServiceType(in.GetServiceType())
	}

	if in.GetService() != "" {
		b.WithService(in.GetService())
	}

	if in.GetIdentity() != "" {
		b.WithIdentity(in.GetIdentity())
	}

	// FIXME
	if in.GetSeverity() != insightV1.SeverityLevel_DEBUG {
		b.WithSeverity(in.GetSeverity())
	}

	if in.GetDetails() != nil {
		d := utils.ValueMapFrom(in.GetDetails())

		b.WithDetails(d)
	}

	if in.GetBefore() != nil {
		ts, err := ptypes.Timestamp(in.GetBefore())
		if err != nil {
			return nil, err
		}

		b.WithBefore(ts)
	}

	if in.GetNotBefore() != nil {
		ts, err := ptypes.Timestamp(in.GetNotBefore())
		if err != nil {
			return nil, err
		}

		b.WithNotBefore(ts)
	}

	if in.GetAfter() != nil {
		ts, err := ptypes.Timestamp(in.GetAfter())
		if err != nil {
			return nil, err
		}

		b.WithAfter(ts)
	}

	if in.GetNotAfter() != nil {
		ts, err := ptypes.Timestamp(in.GetNotAfter())
		if err != nil {
			return nil, err
		}

		b.WithNotAfter(ts)
	}

	return b, nil
}
