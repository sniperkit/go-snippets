package main

import (
	"regexp"
	"strings"
	"testing"
)

func TestFilename(t *testing.T) {
	for k, v := range map[string]string{
		"*7* johab.py gen_probe.ko palmtx.h macpath.py tzp dm-dirty-log.h bh1770.h pktloc faillog.8.gz zconf.gperf":                                                                                             "bh1770.h",
		"*[0123456789]*[auoei]* IBM1008_420.so zfgrep limits.conf.5.gz tc.h ilogb.3.gz limits.conf CyrAsia-TerminusBold28x14.psf.gz nf_conntrack_sip.ko DistUpgradeViewNonInteractive.pyc NFKDQC":               "IBM1008_420.so",
		"*.??? max_user_watches arptables.h net_namespace Kannada.pl menu_no_no.utf-8.vim shtags.1 unistd_32_ia32.h gettext-tools.mo ntpdate.md5sums linkat.2.gz":                                               "menu_no_no.utf-8.vim",
		"*.pdf OldItali.pl term.log plymouth-upstart-bridge rand.so libipw.ko jisfreq.pyc impedance-analyzer xmon.h 1.5.0.3.txt bank":                                                                           "-",
		"g*.* 56b8a0b6.0 sl.vim digctl.h groff-base.conffiles python-software-properties.md5sums CountryInformation.py use_zero_page session-noninteractive d2i_RSAPublicKey.3ssl.gz container-detect.log.4.gz": "groff-base.conffiles"} {
		if r := filename(k); r != v {
			t.Errorf("failed: filename %s is %s, got %s",
				k, v, r)
		}
	}
}

func filename(q string) string {
	var r []string
	f := strings.Fields(q)
	p := regexp.QuoteMeta(f[0])
	p = strings.Replace(p, `\?`, ".", -1)
	p = strings.Replace(p, `\*`, ".*", -1)
	p = strings.Replace(p, `\[`, "[", -1)
	p = strings.Replace(p, `\]`, "]", -1)
	p = "^" + p + "$"
	pattern := regexp.MustCompile(p)
	for i := 1; i < len(f); i++ {
		if pattern.MatchString(f[i]) {
			r = append(r, f[i])
		}
	}
	if len(r) > 0 {
		return strings.Join(r, " ")
	} else {
		return "-"
	}
}
