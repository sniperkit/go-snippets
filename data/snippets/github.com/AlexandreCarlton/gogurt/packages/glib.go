package packages

// This is breaking because of automake (which is also broken... :( )
// Why is it invoking automake though?
// Should probably fix automake though... god it's such a kick in the teeth.

import (
	"fmt"
	"strings"
	"github.com/alexandrecarlton/gogurt"
)

type GLib struct{}

func (glib GLib) Name() string {
	return "glib"
}

func (glib GLib) URL(version string) string {
	majorVersion := strings.Join(strings.Split(version, ".")[:2], ".")
	return fmt.Sprintf("http://ftp.acc.umu.se/pub/gnome/sources/glib/%s/glib-%s.tar.xz", majorVersion, version)
}

func (glib GLib) Build(config gogurt.Config) error {
	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(glib),
		Args: []string{
			"--disable-shared",
			"--disable-gtk-doc",
			"--disable-gtk-doc-html",
			"--disable-gtk-doc-pdf",
			"--disable-libmount", // TODO: enable this if possible, libmount found in util-linux
			"--enable-static",
			"--with-pcre=system",
		},
		PkgConfigPaths: []string{
			config.PkgConfigLibDir(LibFFI{}),
			config.PkgConfigLibDir(Pcre{}),
			config.PkgConfigLibDir(Zlib{}),
		},
		// pkgconfig points to <lib>/libffi-3.2.1/include, which we've tweaked.
		// Maybe we should undo that? Only Python and Vim need it.
		CFlags: []string{
			"-I" + config.IncludeDir(LibFFI{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(LibFFI{}),
		},
	}.Cmd()
	if err := configure.Run(); err != nil {
		return err
	}
	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()
	return make.Run()
}

func (glib GLib) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{Args: []string{"install"}}.Cmd()
	return makeInstall.Run()
}

func (glib GLib) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		LibFFI{},
		Pcre{},
		Zlib{},
	}
}
