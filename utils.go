package gogurt

import (
	// "log"
	"os/exec"
)

// Just call sed.
func ReplaceInFile(fileName, pattern, repl string) error {
	cmd := exec.Command(
		"sed",
		"--in-place",
		"--regexp-extended",
		"s/" + pattern + "/" + repl + "/g",
		fileName,
	)
	// log.Println("Sed command: ", cmd)

	return cmd.Run()

	// patt, err := regexp.Compile("\n" + pattern + "\n")
	// if err != nil {
	// 	return err
	// }

	// stat, err := os.Stat(fileName)
	// if err != nil {
	// 	return err
	// }

	// fileContents, err := ioutil.ReadFile(fileName)
	// if err != nil {
	// 	return err
	// }

	// replacedContents := patt.ReplaceAll(fileContents, []byte(repl))
	// ioutil.WriteFile(fileName, replacedContents, stat.Mode())

	// return nil
}

// CopyFile recursively copies the src to the dest, overwriting it if it already exists.
// TODO: Make into a purely go function, not shelled out to 'cp -rf'.
func CopyFile(src, dest string) error {
	cmd := exec.Command("cp", "--recursive", "--force" , src, dest)
	return cmd.Run()
}
