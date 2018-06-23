#!/usr/bin/env bash
set -euo pipefail
set -x
# generates test fixtures

# toplevel
#	hash1
#	hash2
#	hash3
#	hash4
#	hashB
#	hashC
#	generate.sh
#	- repo_c
#		- <file>
#	- repo_b
#		- <file>
#		- <submodule repo_c>
#	- repo_a
#		- <file>
#		- <submodule repo_b>

export GIT_CONFIG_NOSYSTEM=true
export HOME='/dev/null'
export GIT_ASKPASS='/bin/true'

export GIT_AUTHOR_NAME='rio'
export GIT_AUTHOR_EMAIL='rio@polydawn.net'
export GIT_COMMITTER_NAME='rio'
export GIT_COMMITTER_EMAIL='rio@polydawn.net'

CWD=${PWD}

TOPLEVEL=$(mktemp -d)
cd $TOPLEVEL
REPO_C='repo-c'
REPO_B='repo-b'
REPO_A='repo-a'

echo $PWD > ./pwd

git init -- "$REPO_C"
cd "$REPO_C"
echo 'some text here' > 'a file name'
git add -- '.'
git commit -m 'sub sub repo commit'
git rev-parse HEAD > ../hashC

cd "$TOPLEVEL"
git init -- "$REPO_B"
cd "$REPO_B"
echo 'yet another test file' > 'this_is_a_file'
git add -- '.'
git submodule add -- "${TOPLEVEL}/${REPO_C}"
git commit -m 'sub repo commit'
git rev-parse HEAD > ../hashB

cd "$TOPLEVEL"
git init -- "$REPO_A"
cd "$REPO_A"

git commit --allow-empty -m 'initial commit'
git rev-parse HEAD > ../hash1

echo 'abcd' > 'file1'
git add -- '.'
git commit -m 'commit file1'
git rev-parse HEAD > ../hash2

echo 'efghi' > 'file2'
git add -- '.'
git commit -m 'commit file2'
git rev-parse HEAD > ../hash3

git submodule add -- "${TOPLEVEL}/${REPO_B}"
git commit -m 'add submodule'
git rev-parse HEAD > ../hash4

cd $TOPLEVEL
git clone --bare $REPO_A bare

tar -czf 'git_test_fixture.tar.gz' -T <(ls)
mv './git_test_fixture.tar.gz' "${CWD}/git_test_fixture.tar.gz"
