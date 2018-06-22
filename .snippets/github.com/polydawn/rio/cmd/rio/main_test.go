package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"go.polydawn.net/go-timeless-api/rio"
	"go.polydawn.net/rio/fs"
	"go.polydawn.net/rio/fs/osfs"
	"go.polydawn.net/rio/fsOp"
	"go.polydawn.net/rio/testutil"
)

func stdBuffers() (stdin, stdout, stderr *bytes.Buffer) {
	return &bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{}
}

type unpackTest struct {
	Name           string
	Args           []string
	ExpectedExit   int
	ExpectedStdout string
	ExpectedStderr string
}

func TestWithoutArgs(t *testing.T) {
	Convey("rio: usage printed to stderr", t, func() {
		args := []string{"rio"}
		stdin, stdout, stderr := stdBuffers()

		ctx := context.Background()
		exitCode := Main(ctx, args, stdin, stdout, stderr)
		So(string(stdout.Bytes()), ShouldBeBlank)
		So(string(stderr.Bytes()), ShouldNotBeBlank)
		firstLine, err := stderr.ReadString('\n')
		So(err, ShouldBeNil)
		So(string(firstLine), ShouldContainSubstring, "usage: rio [<flags>] <command> [<args> ...]")
		So(string(stderr.Bytes()), ShouldNotContainSubstring, "usage: rio [<flags>] <command> [<args> ...]")
		So(exitCode, ShouldEqual, rio.ExitCodeForCategory(rio.ErrUsage))
	})
}

func TestUnpackBogusFlag(t *testing.T) {
	Convey("rio: usage printed to stderr", t, func() {
		args := []string{"rio", "unpack", "--bogus"}
		stdin, stdout, stderr := stdBuffers()

		ctx := context.Background()
		exitCode := Main(ctx, args, stdin, stdout, stderr)
		So(string(stdout.Bytes()), ShouldBeBlank)
		So(string(stderr.Bytes()), ShouldResemble, "error parsing args: unknown long flag '--bogus'\n")
		So(exitCode, ShouldEqual, rio.ExitCodeForCategory(rio.ErrUsage))
	})
}

/*
	Tests against pre-generated, known fixtures of tar binary blobs.

	These tests allow us to cover compat with other tar impls, compression, etc.
*/
func TestTarFixtureUnpack(t *testing.T) {
	Convey("rio: unpacking of tar fixtures", t,
		testutil.Requires(testutil.RequiresCanManageOwnership, func() {
			testutil.WithTmpdir(func(tmpDir fs.AbsolutePath) {
				Convey("Unpack a fixture from gnu tar which includes a base dir", func() {
					ctx := context.Background()
					wareID := "tar:5y6NvK6GBPQ6CcuNyJyWtSrMAJQ4LVrAcZSoCRAzMSk5o53pkTYiieWyRivfvhZwhZ"
					source := "file://../../transmat/tar/fixtures/tar_withBase.tgz"
					for _, fixture := range []unpackTest{
						{"UnpackBasic",
							[]string{
								"rio",
								"unpack",
								wareID,
								tmpDir.String(),
								"--uid=keep", "--gid=keep",
								fmt.Sprintf("--placer=%s", rio.Placement_Direct),
								fmt.Sprintf("--source=%s", source),
							},
							0,
							wareID,
							fmt.Sprintf(`log: lvl=info msg=read for ware "%s" opened from warehouse "%s"`, wareID, source),
						},
						{"UnpackJsonFormat",
							[]string{
								"rio",
								"unpack",
								wareID,
								tmpDir.String(),
								"--uid=keep", "--gid=keep",
								fmt.Sprintf("--placer=%s", rio.Placement_Direct),
								fmt.Sprintf("--source=%s", source),
								fmt.Sprintf("--format=%s", format_Json),
							},
							0,
							fmt.Sprintf(`{"log":null,"prog":null,"result":{"wareID":"%s","error":null}}`, wareID),
							"",
						},
					} {
						Convey(fmt.Sprintf("- test %q", fixture.Name), func() {
							stdin, stdout, stderr := stdBuffers()
							exitCode := Main(ctx, fixture.Args, stdin, stdout, stderr)
							So(lastLine(string(stdout.Bytes())), ShouldEqual, fixture.ExpectedStdout)
							So(lastLine(string(stderr.Bytes())), ShouldEqual, fixture.ExpectedStderr)
							So(exitCode, ShouldEqual, fixture.ExpectedExit)
							if fixture.ExpectedExit != 0 {
								Convey("The filesystem should not have things", func() {
									afs := osfs.New(tmpDir)
									_, err := afs.LStat(fs.MustRelPath("."))
									So(err, ShouldNotBeNil)
								})
							} else {
								Convey("The filesystem contains the correct unpacked fixture", func() {
									var err error
									fmeta, reader, err := fsOp.ScanFile(osfs.New(tmpDir), fs.MustRelPath("ab"))
									So(err, ShouldBeNil)
									So(fmeta.Name, ShouldResemble, fs.MustRelPath("ab"))
									So(fmeta.Type, ShouldResemble, fs.Type_File)
									So(fmeta.Uid, ShouldEqual, 7000)
									So(fmeta.Gid, ShouldEqual, 7000)
									So(fmeta.Mtime.UTC(), ShouldResemble, time.Date(2015, 05, 30, 19, 53, 35, 0, time.UTC))
									body, err := ioutil.ReadAll(reader)
									So(string(body), ShouldResemble, "")

									fmeta, reader, err = fsOp.ScanFile(osfs.New(tmpDir), fs.MustRelPath("bc"))
									So(err, ShouldBeNil)
									So(fmeta.Name, ShouldResemble, fs.MustRelPath("bc"))
									So(fmeta.Type, ShouldResemble, fs.Type_Dir)
									So(fmeta.Mtime.UTC(), ShouldResemble, time.Date(2015, 05, 30, 19, 53, 35, 0, time.UTC))
									So(reader, ShouldBeNil)

									fmeta, reader, err = fsOp.ScanFile(osfs.New(tmpDir), fs.MustRelPath("."))
									So(err, ShouldBeNil)
									So(fmeta.Name, ShouldResemble, fs.MustRelPath("."))
									So(fmeta.Type, ShouldResemble, fs.Type_Dir)
									So(fmeta.Mtime.UTC(), ShouldResemble, time.Date(2015, 05, 30, 19, 53, 35, 0, time.UTC))
									So(reader, ShouldBeNil)
								})
							}
						})
					}

				})
			})
		}),
	)
}

func lastLine(str string) string {
	str = strings.TrimRight(str, "\n")
	ss := strings.Split(str, "\n")
	fmt.Printf("::goddamnit \n\t%#v\n", ss)
	return ss[len(ss)-1]
}
