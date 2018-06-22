package insight

import (
	"testing"

	insightV1 "github.com/homebot/protobuf/pkg/api/insight/v1"
)

func TestFilter(t *testing.T) {
	cases := []struct {
		B *Builder
		M *insightV1.LogMessage
		E bool
	}{
		{
			B: NewFilter(),
			M: &insightV1.LogMessage{
				Service:  "foobar",
				Severity: insightV1.SeverityLevel_DEBUG,
				Message:  "case1",
			},
			E: true,
		},
		{
			B: NewFilter().
				WithService("foobar"),
			M: &insightV1.LogMessage{
				Service:  "foobar",
				Severity: insightV1.SeverityLevel_DEBUG,
				Message:  "case1",
			},
			E: true,
		},
		{
			B: NewFilter().
				WithSeverity(insightV1.SeverityLevel_ERROR),
			M: &insightV1.LogMessage{
				Service:  "foobar",
				Severity: insightV1.SeverityLevel_ERROR,
				Message:  "case1",
			},
			E: true,
		},
		{
			B: NewFilter().
				WithServiceType("foo"),
			M: &insightV1.LogMessage{
				ServiceType: "foo",
				Severity:    insightV1.SeverityLevel_ERROR,
				Message:     "case1",
			},
			E: true,
		},
		{
			B: NewFilter().
				WithIdentity("admin"),
			M: &insightV1.LogMessage{
				Identity: "admin",
				Severity: insightV1.SeverityLevel_ERROR,
				Message:  "case1",
			},
			E: true,
		},
	}

	for idx, c := range cases {
		if c.B.Match(c.M) != c.E {
			t.Fatalf("Case %d: expected %v but got %v\n\t%s", idx, c.E, !c.E, c.B.String())
		}
	}
}
