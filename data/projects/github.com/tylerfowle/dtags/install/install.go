package install

import (
	"log"
	"os"
	"os/user"
)

const script = `#!/usr/bin/env bash
dtags=dtags
declare -a arr=("" "add" "del" "tags" "list" "ls" "completion", "install")
found=0

# check if arg is a single arg command in arr array
for i in "${arr[@]}"; do
    if [[ $i == $1 ]]; then
        found=1
        ${dtags} $@
        false 1
    fi
done

if [[ -d "$(${dtags} $1)" && ${found} -ne 1 ]]; then
    cd "$(${dtags} $1)"
elif [[ ${found} -ne 1 ]]; then
    echo "no directory found for tag [$1]"
    false 1
fi`

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func WriteFile() (string, error) {

	// get user info
	u, err := user.Current()
	check(err)

	file := u.HomeDir + "/.config/dtags/dt"

	// create file
	f, err := os.Create(file)
	check(err)
	err = os.Chmod(file, 0755)
	check(err)
	defer f.Close()

	// write script to file
	f.Write([]byte(script))

	return script, nil
}
