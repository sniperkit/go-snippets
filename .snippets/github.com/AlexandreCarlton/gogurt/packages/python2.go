package packages

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"github.com/alexandrecarlton/gogurt"
)

type Python2 struct{}

func (python2 Python2) Name() string {
	return "python2"
}

func (python2 Python2) URL(version string) string {
	return fmt.Sprintf("https://www.python.org/ftp/python/%s/Python-%s.tar.xz", version, version)
}

func (python2 Python2) Build(config gogurt.Config) error {
	// Termcap has been succeeded by terminfo.
	gogurt.ReplaceInFile("Modules/Setup.dist", "termcap", "tinfow")

	gogurt.ReplaceInFile("Modules/Setup.dist", "^#SSL=.*", "SSL=" + strings.Replace(config.InstallDir(OpenSSL{}), "/", "\\/", -1))
	uncommentModule("_ssl")
	uncommentModule("\t-DUSE_SSL")
	uncommentModule("\t-L\\$\\(SSL\\)") // TODO: Redo this section, these aren't modules.
	gogurt.ReplaceInFile("Modules/Setup.dist", "-lcrypto", "-lcrypto -ldl")

	addModule("future_builtins", []string{"future_builtins.c"})
	addModule("_multiprocessing", []string{
		"_multiprocessing/multiprocessing.c",
		"_multiprocessing/semaphore.c",
		"_multiprocessing/socket_connection.c",
		"-IModules/_multiprocessing",
	})
	addModule("_ctypes", []string{
		"_ctypes/callbacks.c",
		"_ctypes/callproc.c",
		"_ctypes/cfield.c",
		"_ctypes/_ctypes.c",
		"_ctypes/_ctypes_test.c",
		"_ctypes/malloc_closure.c",
		"_ctypes/stgdict.c",
		"-IModules/_ctypes",
		"-lffi",
	})
	uncommentModule("_bisect")
	uncommentModule("_collections")
	uncommentModule("_csv")
	uncommentModule("_elementtree")
	uncommentModule("_functools")
	uncommentModule("_heapq")
	uncommentModule("_io")
	uncommentModule("_md5")
	uncommentModule("_random")
	uncommentModule("_sha")
	uncommentModule("_sha256")
	uncommentModule("_sha512")
	uncommentModule("_socket")
	uncommentModule("_struct")

	uncommentModule("array")
	uncommentModule("binascii")
	uncommentModule("cmath")
	uncommentModule("cPickle")
	uncommentModule("cStringIO")
	uncommentModule("datetime")
	uncommentModule("fcntl")
	uncommentModule("itertools")
	uncommentModule("math")
	uncommentModule("operator")
	uncommentModule("pyexpat")
	uncommentModule("readline")
	uncommentModule("select")
	uncommentModule("termios")
	uncommentModule("time")
	uncommentModule("unicodedata")
	uncommentModule("zlib")

	gogurt.ReplaceInFile("Modules/Setup.dist", "#\\*shared\\*", "\\*static\\*")


	configure := gogurt.ConfigureCmd{
		Prefix: config.InstallDir(python2),
		Args: []string{
			"--disable-shared",
			"--enable-unicode=ucs4",
			"--with-system-expat",
			"--with-system-libffi",
			"--with-ensurepip=install",
		},
		CFlags: []string{
			"-I" + config.IncludeDir(Expat{}),
			"-I" + config.IncludeDir(LibFFI{}),
			"-I" + config.IncludeDir(Ncurses{}),
			"-I" + config.IncludeDir(OpenSSL{}),
			"-I" + config.IncludeDir(ReadLine{}),
			"-I" + config.IncludeDir(Zlib{}),
		},
		LdFlags: []string{
			"-L" + config.LibDir(Expat{}),
			"-L" + config.LibDir(LibFFI{}),
			"-L" + config.LibDir(Ncurses{}),
			"-L" + config.LibDir(OpenSSL{}),
			"-L" + config.LibDir(ReadLine{}),
			"-L" + config.LibDir(Zlib{}),
		},
		// While not necessary to build Python, these are recorded used in building applications using Python (e.g. GDB).
		Libs: []string{
			"-lffi",
			"-lexpat",
			"-lreadline",
			"-ltinfow",
			"-lssl",
			"-lcrypto",
			"-lz",
		},
	}.Cmd();
	configure.Env = append(
		configure.Env,
		"LINKFORSHARED= ",
		"DYNLOADFILE=dynload_stub.o",
		"CFLAGSFORSHARED=-fPIC")
	if err := configure.Run(); err != nil {
		return err
	}

	make := gogurt.MakeCmd{
		Jobs: config.NumCores,
	}.Cmd()

	return make.Run();
}

// TODO: Optimise - Taking too slow
// Perhaps only open it once, and modify it then.

func uncommentModule(module string) error {
	return gogurt.ReplaceInFile(
		"Modules/Setup.dist",
		"^#(" + module + ".*)",
		"\\1",
	)
}

func addModule(module string, sources []string) error {
	setupDistContents, _ := ioutil.ReadFile("Modules/Setup.dist")
	if bytes.Contains(setupDistContents, []byte(module)) {
		log.Printf("Modules/Setup.dist already contains %s.\n", module)
	} else {
		log.Printf("Adding module %s to Modules/Setup.dist.\n", module)
		setupDist, err := os.OpenFile("Modules/Setup.dist", os.O_WRONLY | os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err)
			return err
		}
		writer := bufio.NewWriter(setupDist)
		fmt.Fprintf(writer, "%s %s\n", module, strings.Join(sources, " "))
		writer.Flush()
		setupDist.Close()
	}
	return nil
}

func (python2 Python2) Install(config gogurt.Config) error {
	makeInstall := gogurt.MakeCmd{
		Args: []string{"install"},
	}.Cmd()
	return makeInstall.Run();
}

func (python2 Python2) Dependencies() []gogurt.Package {
	return []gogurt.Package{
		Expat{},
		LibFFI{},
		Ncurses{},
		OpenSSL{},
		ReadLine{},
		Zlib{},
	}
}

