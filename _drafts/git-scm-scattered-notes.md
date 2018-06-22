---
layout: page
title: "Git SCM : Scattered notes..."
description: ""
---

_Git, shake dreams from your hair my pretty tool..._

List files from a given branch/commit :

    git ls-files --with-tree=pg

    git ls-files --with-tree=pg~10

    git ls-files --with-tree=origin/develop

List tags sorted by version :

    git tag | sort -V

Git merge branches that have unrelated histories :

    git pull --allow-unrelated-histories <remote> <branch>

Generate a __patch__ for the currently unstaged content :

    git diff > unstaged-stuff.patch

    git apply unstaged-stuff.patch

Generate a __patch__ for a whole `some_topic_branch` branch from the point in
history where it “first meets” the `some_trunk_branch` (note that we ommitted
the `HEAD` argument to `diff ...` here, assuming we're already on the topic
branch) :

    git diff $(git merge-base <some_trunk_branch> <some_topic_branch>) > <output_patch_file>



## Git split tree

    $ git subtree split -b pg --prefix=playground/

## ...

* [Github desktop app.](https://desktop.github.com/)

## Initial setup

    $ git config --global user.name "`id -un`"
    $ git config --global user.email "`id -un`@`hostname -f`"

    $ git config --global color.ui auto

    $ git config --global merge.tool vimdiff

    $ git config --global -l

## Basic day-to-day usage

    git init .

    git clone [--recursive]

    git branch [-a]

    git show-branch master origin/master

    git log --graph ...

    gl HEAD..FETCH_HEAD

    git submodule update --init --recursive

    git remote -v

    git remote set-head origin -a

    git reset

    git reset --hard

    git checkout <file>

    git push

    git push origin master

    git push -u origin my_branch
    git push -f origin my_branch

    git push origin :my_branch

## Misc.

    git config ...share group ???

## Documentation pointers

git-scm.org + the book.

