package cove

// Copyright 2015 MediaMath <http://www.mediamath.com>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestNoMissing(t *testing.T) {
	miss, err := parseString(noMissing)
	if err != nil {
		t.Errorf("%v", err)
	}

	if len(miss) > 0 {
		t.Errorf("%v", miss)
	}
}

func TestHasMissing(t *testing.T) {

	miss, err := parseString(hasMissing)
	if err != nil {
		t.Errorf("%v", err)
	}

	if err := has(miss,
		"github.com/BurntSushi/toml",
		"github.com/mattn/go-sqlite3",
		"github.com/MediaMath/fpg",
		"github.com/MediaMath/fpg/boot",
		"github.com/MediaMath/fpg/config",
		"github.com/MediaMath/fpg/entity",
		"github.com/MediaMath/fpg/fee/generate",
		"github.com/MediaMath/fpg/fee/privacy_compliance_cost",
		"github.com/MediaMath/fpg/util",
		"github.com/MediaMath/fpg/util/logger",
		"github.com/op/go-logging"); err != nil {
		t.Errorf("%v", err)
	}
}

func has(packs []Package, dep ...string) error {
	miss := PackagesAsStrings(packs)

	if len(miss) < len(dep) {
		return fmt.Errorf("Expected is missing items:%v %v", miss, dep)
	}

	if len(dep) < len(miss) {
		return fmt.Errorf("Expected has too many items:%v %v", miss, dep)
	}

	sort.Strings(miss)
	sort.Strings(dep)
	for i := range dep {
		if miss[i] != dep[i] {
			return fmt.Errorf("Position %i mismatch: %v %v", miss, dep)
		}
	}

	return nil
}

func parseString(toParse string) ([]Package, error) {
	var parsed missing
	if err := json.NewDecoder(strings.NewReader(toParse)).Decode(&parsed); err != nil {
		return []Package{}, err
	}

	return missingFromParsed(&parsed)
}

const hasMissing = `
{
	"Dir": "/Users/kklipsch/go/src/github.com/MediaMath/fpg",
	"ImportPath": "github.com/MediaMath/fpg",
	"Name": "main",
	"Target": "/Users/kklipsch/go/bin/fpg",
	"Stale": true,
	"Root": "/Users/kklipsch/go",
	"GoFiles": [
		"main.go"
	],
	"Imports": [
		"flag",
		"fmt",
		"github.com/MediaMath/Keryx/client",
		"github.com/MediaMath/fpg/adama",
		"github.com/MediaMath/fpg/boot",
		"github.com/MediaMath/fpg/config",
		"github.com/MediaMath/fpg/entity",
		"github.com/MediaMath/fpg/fee/aggregate",
		"github.com/MediaMath/fpg/loop",
		"github.com/MediaMath/fpg/statsd",
		"github.com/MediaMath/fpg/util",
		"github.com/MediaMath/fpg/util/logger",
		"os"
	],
	"Deps": [
		"bufio",
		"bytes",
		"container/list",
		"crypto",
		"crypto/aes",
		"crypto/cipher",
		"crypto/des",
		"crypto/dsa",
		"crypto/ecdsa",
		"crypto/elliptic",
		"crypto/hmac",
		"crypto/md5",
		"crypto/rand",
		"crypto/rc4",
		"crypto/rsa",
		"crypto/sha1",
		"crypto/sha256",
		"crypto/sha512",
		"crypto/subtle",
		"crypto/tls",
		"crypto/x509",
		"crypto/x509/pkix",
		"database/sql",
		"database/sql/driver",
		"encoding",
		"encoding/asn1",
		"encoding/base64",
		"encoding/binary",
		"encoding/csv",
		"encoding/gob",
		"encoding/hex",
		"encoding/json",
		"encoding/pem",
		"errors",
		"flag",
		"fmt",
		"github.com/BurntSushi/toml",
		"github.com/MediaMath/Keryx/client",
		"github.com/MediaMath/Keryx/files",
		"github.com/MediaMath/Keryx/message",
		"github.com/MediaMath/Keryx/pg",
		"github.com/MediaMath/Keryx/pg/xlog",
		"github.com/MediaMath/Keryx/streamer",
		"github.com/MediaMath/Keryx/testingUtils",
		"github.com/MediaMath/Keryx/utils",
		"github.com/MediaMath/fpg/adama",
		"github.com/MediaMath/fpg/boot",
		"github.com/MediaMath/fpg/config",
		"github.com/MediaMath/fpg/entity",
		"github.com/MediaMath/fpg/fee",
		"github.com/MediaMath/fpg/fee/active_ids",
		"github.com/MediaMath/fpg/fee/ad_serving_cost",
		"github.com/MediaMath/fpg/fee/ad_verification_cost",
		"github.com/MediaMath/fpg/fee/agency_margin_pct",
		"github.com/MediaMath/fpg/fee/aggregate",
		"github.com/MediaMath/fpg/fee/byo_adserve_flag",
		"github.com/MediaMath/fpg/fee/byo_adver_flag",
		"github.com/MediaMath/fpg/fee/byo_dyncre_flag",
		"github.com/MediaMath/fpg/fee/byo_pmp_flag",
		"github.com/MediaMath/fpg/fee/byo_privacy_flag",
		"github.com/MediaMath/fpg/fee/byo_udi_flag",
		"github.com/MediaMath/fpg/fee/client_exchange_cost",
		"github.com/MediaMath/fpg/fee/collect_byos_flag",
		"github.com/MediaMath/fpg/fee/collect_udi_data_cost_cpm",
		"github.com/MediaMath/fpg/fee/collect_udi_data_cost_pct",
		"github.com/MediaMath/fpg/fee/contextual_cost_cpm",
		"github.com/MediaMath/fpg/fee/donotbill_adserve_flag",
		"github.com/MediaMath/fpg/fee/donotbill_byo_pmp_flag",
		"github.com/MediaMath/fpg/fee/donotbill_byos_eos_flag",
		"github.com/MediaMath/fpg/fee/donotbill_byos_rtb_flag",
		"github.com/MediaMath/fpg/fee/donotbill_dyncre_flag",
		"github.com/MediaMath/fpg/fee/donotbill_margin_flag",
		"github.com/MediaMath/fpg/fee/donotbill_privacy_flag",
		"github.com/MediaMath/fpg/fee/donotbill_udi_flag",
		"github.com/MediaMath/fpg/fee/dynamic_creative_cost",
		"github.com/MediaMath/fpg/fee/edata",
		"github.com/MediaMath/fpg/fee/esupply",
		"github.com/MediaMath/fpg/fee/generate",
		"github.com/MediaMath/fpg/fee/mm_data_cost_cpm",
		"github.com/MediaMath/fpg/fee/mm_exchange_cost",
		"github.com/MediaMath/fpg/fee/mm_managed_service_fee",
		"github.com/MediaMath/fpg/fee/mm_optimization_fee",
		"github.com/MediaMath/fpg/fee/mm_platform_fee",
		"github.com/MediaMath/fpg/fee/mm_profit_share_pct",
		"github.com/MediaMath/fpg/fee/ms_flag",
		"github.com/MediaMath/fpg/fee/oes_flag",
		"github.com/MediaMath/fpg/fee/opto_flag",
		"github.com/MediaMath/fpg/fee/pmp_noopto",
		"github.com/MediaMath/fpg/fee/pmp_opto",
		"github.com/MediaMath/fpg/fee/privacy_compliance_cost",
		"github.com/MediaMath/fpg/fee/track_adserve_flag",
		"github.com/MediaMath/fpg/fee/track_adver_flag",
		"github.com/MediaMath/fpg/fee/track_byo_oes_flag",
		"github.com/MediaMath/fpg/fee/track_byo_pmp_flag",
		"github.com/MediaMath/fpg/fee/track_byos_rtb_flag",
		"github.com/MediaMath/fpg/fee/track_dyncre_flag",
		"github.com/MediaMath/fpg/fee/track_privacy_flag",
		"github.com/MediaMath/fpg/fee/track_udi_flag",
		"github.com/MediaMath/fpg/keryx",
		"github.com/MediaMath/fpg/loop",
		"github.com/MediaMath/fpg/reports",
		"github.com/MediaMath/fpg/statsd",
		"github.com/MediaMath/fpg/timer",
		"github.com/MediaMath/fpg/util",
		"github.com/MediaMath/fpg/util/logger",
		"github.com/go-fsnotify/fsnotify",
		"github.com/jmoiron/sqlx",
		"github.com/jmoiron/sqlx/reflectx",
		"github.com/lib/pq",
		"github.com/lib/pq/oid",
		"github.com/mattn/go-sqlite3",
		"github.com/op/go-logging",
		"github.com/peterbourgon/g2s",
		"hash",
		"hash/adler32",
		"io",
		"io/ioutil",
		"log",
		"math",
		"math/big",
		"math/rand",
		"net",
		"net/url",
		"os",
		"os/exec",
		"os/signal",
		"os/user",
		"path",
		"path/filepath",
		"reflect",
		"regexp",
		"regexp/syntax",
		"runtime",
		"runtime/cgo",
		"runtime/pprof",
		"sort",
		"strconv",
		"strings",
		"sync",
		"sync/atomic",
		"syscall",
		"testing",
		"text/tabwriter",
		"time",
		"unicode",
		"unicode/utf16",
		"unicode/utf8",
		"unsafe"
	],
	"Incomplete": true,
	"DepsErrors": [
		{
			"ImportStack": [
				"github.com/MediaMath/fpg",
				"github.com/MediaMath/fpg/boot",
				"github.com/MediaMath/fpg/fee/generate",
				"github.com/MediaMath/fpg/fee/privacy_compliance_cost",
				"github.com/MediaMath/fpg/util",
				"github.com/MediaMath/fpg/util/logger",
				"github.com/MediaMath/fpg/config",
				"github.com/BurntSushi/toml"
			],
			"Pos": "../fpg/config/config.go:5:2",
			"Err": "cannot find package \"github.com/BurntSushi/toml\" in any of:\n\t/usr/local/Cellar/go/1.4.2/libexec/src/github.com/BurntSushi/toml (from $GOROOT)\n\t/Users/kklipsch/go/src/github.com/BurntSushi/toml (from $GOPATH)"
		},
		{
			"ImportStack": [
				"github.com/MediaMath/fpg",
				"github.com/MediaMath/fpg/boot",
				"github.com/MediaMath/fpg/entity",
				"github.com/mattn/go-sqlite3"
			],
			"Pos": "../fpg/entity/store.go:6:2",
			"Err": "cannot find package \"github.com/mattn/go-sqlite3\" in any of:\n\t/usr/local/Cellar/go/1.4.2/libexec/src/github.com/mattn/go-sqlite3 (from $GOROOT)\n\t/Users/kklipsch/go/src/github.com/mattn/go-sqlite3 (from $GOPATH)"
		},
		{
			"ImportStack": [
				"github.com/MediaMath/fpg",
				"github.com/MediaMath/fpg/boot",
				"github.com/MediaMath/fpg/fee/generate",
				"github.com/MediaMath/fpg/fee/privacy_compliance_cost",
				"github.com/MediaMath/fpg/util",
				"github.com/MediaMath/fpg/util/logger",
				"github.com/op/go-logging"
			],
			"Pos": "../fpg/util/logger/logger.go:5:2",
			"Err": "cannot find package \"github.com/op/go-logging\" in any of:\n\t/usr/local/Cellar/go/1.4.2/libexec/src/github.com/op/go-logging (from $GOROOT)\n\t/Users/kklipsch/go/src/github.com/op/go-logging (from $GOPATH)"
		}
	],
	"TestGoFiles": [
		"main_test.go"
	],
	"TestImports": [
		"testing"
	]
}
`

const noMissing = `{
	"Dir": "/usr/local/Cellar/go/1.4.2/libexec/src/os/exec",
	"ImportPath": "os/exec",
	"Name": "exec",
	"Doc": "Package exec runs external commands.",
	"Target": "/usr/local/Cellar/go/1.4.2/libexec/pkg/darwin_amd64/os/exec.a",
	"Goroot": true,
	"Standard": true,
	"Root": "/usr/local/Cellar/go/1.4.2/libexec",
	"GoFiles": [
		"exec.go",
		"lp_unix.go"
	],
	"IgnoredGoFiles": [
		"lp_plan9.go",
		"lp_windows.go",
		"lp_windows_test.go"
	],
	"Imports": [
		"bytes",
		"errors",
		"io",
		"os",
		"path/filepath",
		"runtime",
		"strconv",
		"strings",
		"sync",
		"syscall"
	],
	"Deps": [
		"bytes",
		"errors",
		"io",
		"math",
		"os",
		"path/filepath",
		"runtime",
		"sort",
		"strconv",
		"strings",
		"sync",
		"sync/atomic",
		"syscall",
		"time",
		"unicode",
		"unicode/utf8",
		"unsafe"
	],
	"TestGoFiles": [
		"lp_test.go",
		"lp_unix_test.go"
	],
	"TestImports": [
		"io/ioutil",
		"os",
		"testing"
	],
	"XTestGoFiles": [
		"example_test.go",
		"exec_test.go"
	],
	"XTestImports": [
		"bufio",
		"bytes",
		"encoding/json",
		"fmt",
		"io",
		"io/ioutil",
		"log",
		"net",
		"net/http",
		"net/http/httptest",
		"os",
		"os/exec",
		"path/filepath",
		"runtime",
		"strconv",
		"strings",
		"testing",
		"time"
	]
}
`
