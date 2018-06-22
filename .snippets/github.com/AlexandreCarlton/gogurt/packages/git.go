package packages

import (
	"fmt"
	"github.com/alexandrecarlton/gogurt"
)

// Note: this requires the Perl module ExtUtils::MakeMaker.
// Generated variables go into config.mak.autogen

type Git struct{}

func (git Git) Name() string {
	return "git"
}

func (git Git) URL(version string) string {
	return fmt.Sprintf("https://www.kernel.org/pub/software/scm/git/git-%s.tar.gz", version)
}

func (git Git) Build(config gogurt.Config) error {

	makeConfigure := gogurt.MakeCmd{Args: []string{"configure"}}.Cmd()
	if err := makeConfigure.Run(); err != nil {
		return err
	}
	gogurt.ReplaceInFile("Makefile", "-lssl", "-lssl -ldl -lz")
	gogurt.ReplaceInFile("Makefile", "-lcrypto", "-lcrypto -ldl -lz")

	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(git),
		Args: []string{
			"--with-curl=" + config.InstallDir(Curl{}),
			"--with-libpcre=" + config.InstallDir(Pcre{}),
			"--with-openssl=" + config.InstallDir(OpenSSL{}),
			"--with-zlib=" + config.InstallDir(Zlib{}),
		},
		// We need to add these so that curl's dependencies (ssl, libz) are picked up,
		// otherwise we don't get cloning via http[s].
		CFlags: []string{
			"-I" + config.IncludeDir(Curl{}),
			"-I" + config.IncludeDir(OpenSSL{}),
			"-I" + config.IncludeDir(Zlib{}),
		},
		LdFlags: []string{
			"-static",
			"-L" + config.LibDir(Curl{}),
			"-L" + config.LibDir(OpenSSL{}),
			"-L" + config.LibDir(Zlib{}),
		},
		Libs: []string{
			"-lcurl",
			"-lssl",
			"-lcrypto",
			"-ldl",
			"-lz",
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
		Args: []string{
			"NEEDS_SSL_WITH_CURL=YesPlease", // so we get -lssl with curl
			"NEEDS_CRYPTO_WITH_SSL=YesPlease",
			"NEEDS_SSL_WITH_CRYPTO=YesPlease",
		},
	}.Cmd()
	return make.Run()
}

func (git Git) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (git Git) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Curl{},
		OpenSSL{},
		Pcre{},
		Zlib{},
	}
}
