package df

import (
	"reflect"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		dockerfile string
		want       *Dockerfile
		wantErr    bool
	}{
		{
			dockerfile: `FROM foo`,
			want:       &Dockerfile{From: "foo"},
		},
		{
			dockerfile: `fRoM foo`,
			want:       &Dockerfile{From: "foo"},
		},
		{
			dockerfile: `FROM foo:bar`,
			want:       &Dockerfile{From: "foo:bar"},
		},
		{
			dockerfile: `x`,
			wantErr:    true,
		},
	}
	for _, test := range tests {
		got, err := Decode(strings.NewReader(test.dockerfile))
		if test.wantErr != (err != nil) {
			t.Errorf("got error %v, wantErr = %v", err, test.wantErr)
		}
		if err != nil {
			continue
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got Dockerfile %+v, want %+v", got, test.want)
		}
	}
}
