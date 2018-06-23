package scheme_test

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/encoding/base62"
	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/scheme"
	"github.com/distributed-vision/go-resources/ids/schemeformat"
	"github.com/distributed-vision/go-resources/ids/schemevisibility"
	"github.com/distributed-vision/go-resources/init/schemeinit"
)

func TestMain(m *testing.M) {

	GetUninitialisedResolver()
	schemeinit.Init()

	os.Exit(m.Run())
}

func ignore(value interface{}) {
}

func GetUninitialisedResolver() {

	schemeId, _ := base62.Decode("0")
	cres, cerr := scheme.Resolve(context.Background(), scheme.Selector{Id: schemeId})

	if cres == nil || cerr == nil {
		panic("GetUninitialisedResolver Failed: schemeResolver.Resolve channels are undefined")
	}

	select {
	case err := <-cerr:
		ignore(err)
	case scheme := <-cres:
		if scheme != nil {
			panic("GetUninitialisedResolver Failed: resolution should fail")
		}
	}
}

func TestGetUntypedScheme(t *testing.T) {
	schemeId, _ := base62.Decode("0")
	schemeNameSelector := scheme.Selector{Name: "untyped",
		Opts: scheme.SelectorOpts{
			IgnoreCase:       true,
			IgnoreWhitespace: true}}

	idScheme, err := scheme.Get(context.Background(), scheme.Selector{Id: schemeId})

	if err != nil {
		t.Errorf("ResolveUntypedScheme Failed With Error: %s", err)
	} else {
		if idScheme == nil {
			t.Errorf("ResolveUntypedScheme Failed: scheme is undefined")
		} else {
			if !bytes.Equal(idScheme.Id(), schemeId) {
				t.Errorf("ResolveUntypedScheme Failed: Id: expected: '%v' got '%v'", schemeId, idScheme.Id())
			}

			if idScheme.Name() != "Untyped" {
				t.Errorf("ResolveUntypedScheme Failed: Name: expected: '%s' got '%s'", "Untyped", idScheme.Name())
			}

			if idScheme.Visibility() != schemevisibility.PUBLIC {
				t.Errorf("ResolveUntypedScheme Failed: Visibility: expected: '%v' got '%v'", schemevisibility.PUBLIC, idScheme.Visibility())

			}

			if idScheme.Format() != schemeformat.FIXED {
				t.Errorf("ResolveUntypedScheme Failed: Format: expected: '%v' got '%v'", schemeformat.FIXED, idScheme.Format())
			}
		}
	}

	idScheme, err = scheme.Get(context.Background(), schemeNameSelector)

	if err != nil {
		t.Errorf("ResolveUntypedScheme Failed With Error: %s", err)
	} else {
		if idScheme == nil {
			t.Errorf("ResolveUntypedScheme Failed: scheme is undefined")
		} else {
			if !bytes.Equal(idScheme.Id(), schemeId) {
				t.Errorf("ResolveUntypedScheme Failed: Id: expected: '%v' got '%v'", schemeId, idScheme.Id())
			}

			if idScheme.Name() != "Untyped" {
				t.Errorf("ResolveUntypedScheme Failed: Name: expected: '%s' got '%s'", "Untyped", idScheme.Name())
			}

			if idScheme.Visibility() != schemevisibility.PUBLIC {
				t.Errorf("ResolveUntypedScheme Failed: Visibility: expected: '%v got '%v'", schemevisibility.PUBLIC, idScheme.Visibility())

			}

			if idScheme.Format() != schemeformat.FIXED {
				t.Errorf("ResolveUntypedScheme Failed: Format: expected: '%v' got '%v'", schemeformat.FIXED, idScheme.Format())
			}
		}
	}
}

func TestSchemeIdFormatting(t *testing.T) {
	schemeId, err := scheme.DecodeId(encodertype.BASE62, "0", "123")

	if err != nil {
		t.Errorf("TestSchemeIdFormatting Failed With Error: %s", err)
	}

	if base62.Encode(schemeId) != "1AfsQd" {
		t.Errorf("TestSchemeIdFormatting unexpected schemeId: expected %s, got: %s", "1AfsQd", base62.Encode(schemeId))
	}

	if schemeId[0] != (1 << 6) {
		t.Errorf("TestSchemeIdFormatting base encoding failed: expected %v, got: %v", (1 << 6), schemeId[0])
	}

	if schemeId[1] != 2 {
		t.Errorf("TestSchemeIdFormatting encoding length incorrect: expected %v, got: %v", 2, schemeId[1])
	}

	if domain.RawSchemeLength(schemeId) != 4 {
		t.Errorf("TestSchemeIdFormatting scheme length incorrect: expected %v, got: %v", 4, domain.SchemeLength(schemeId))
	}

	if !bytes.Equal(schemeId[2:2+schemeId[1]], base62.MustDecode("123")) {
		t.Errorf("TestSchemeIdFormatting extension encoding failed: expected %v, got: %v", base62.MustDecode("123"), schemeId[2:2+schemeId[1]])
	}

}
