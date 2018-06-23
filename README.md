<!--![MEDKIT](medkit.png "MEDKIT")-->
# MEDKIT
> MEDKIT (Multi-Environment Dotfiles Kit) is the dotfile management solution for the obsessive compulsive.

[![GitHub issues][github-issues-image]][github-issues-url]
[![Go Report Card](https://goreportcard.com/badge/github.com/LeadPipeSoftware/medkit)](https://goreportcard.com/report/github.com/LeadPipeSoftware/medkit)
[![CircleCI](https://circleci.com/gh/LeadPipeSoftware/medkit.svg?style=shield)](https://circleci.com/gh/LeadPipeSoftware/medkit)

### Complete Command List

[commands](docs/commands.md)

### Overview

MEDKIT is a tool to help you take control of your local environment configuration including:

* Dotfile Management
* Run Command (rc) Script Management
* PATH Variable Management
* Software Installation

MEDKIT also has a native concept of [environments][environments] which allow you to categorize your settings to match
the systems you're using. This lets you customize each environment with ease.

### Manage Your Dotfiles
The trouble with dotfiles is that, unless you have just one login on one computer, you'll have lots of them spread all
over the place. If you're like most people, you want your preferences and configuration to be consistent no matter where
you log in.

MEDKIT follows a "Write Once, Use Anywhere" philosophy. With a little help from your favorite version control tool, you
can keep a "master" copy of your dotfiles and apply them anywhere you go.

- Install the dotfiles in your repository
- Scan for new dotfiles
- Add dotfiles to your repository

### Organize Your Startup
Startup files often start simple, but turn into a nightmare to manage. It gets worse when tools and installers try to be
helpful and make changes to them without asking. You can take advantage of MEDKIT's version control integration and
templates to get your startup scripts under control once and for all.

- Start with the MEDKIT template or roll your own
- Run `*.source` scripts at startup

### Organize Your PATH
Much like startup scripts, keeping your environment's `PATH` in order can be a real challenge. MEDKIT's convention-based
system can help!

- Run `*.path` scripts
- Determine the order of each PATH element

### Install Your Software
Getting all your tools installed is often one of the biggest challenges when you start using a new computer for the
first time. With MEDKIT, you can script your installation scripts and run them as-needed.

- Run `*.installer` scripts
- Remember what has already been run
- Reset if needed

### Environments <a id="environments"></a>
MEDKIT allows you to define one or more environments, but right now it's a big mystery.

<!--
## Philosophy
Bringing forth the notion of topics is one of the simple yet genius things Zach brought to his dotfiles. This makes it
super easy to organize things logically and keep things tidy. For example, if you start using the greatest editor of all
time, vim, then you create a `vim` folder, drop your `.vimrc` file (named as `.vimrc.symlink`) in there, and run `redot`
to handle your vim settings. Everything vim related will live in that folder. If you lose your freaking mind and decide
to use a different editor then cleaning up is as simple as removing the `vim` folder and re-running `redot`.

## Conventions Rule!
I didn't want to have to edit the "base" files any time I made a change or added something new. To accomplish that, I
adopted and tweaked Zach's convention-based setup:

### Global Conventions
- `bin/`: This directory is added to $PATH and is where to put useful scripts
- `homebrew/Brewfile`: This is a [Homebrew bundle file][homebrew-bundle] that gets executed if you elect to do so when you run `redot`

### Topic Conventions
- `<topic>/path.sh`: Any file with this name will be sourced during `$PATH` setup
- `<topic>/install.sh`: Any file with this name will be executed if you elect to do so when you run `redot`
- `<topic>/*.symlink`: Any file with the `.symlink` extension will be symlinked into your `$HOME` directory
- `<topic>/*.source`: Any file with the `.source` extension will be sourced when you run `redot`
-->

## Installation
macOS:
```
TBD
```
Linux:
```
TBD
```
Windows WSL:
```
TBD
```

## Usage

### Dotfiles structure
MEDKIT operates on a convention-based directory structure consisting of a root, and N bundle directories:

```
dotfiles/
├ bundles
│   ├ go
│   │   ├ Brewfile
│   │   └ path.sh
│   └ macos
│       └ Brewfile
├ homebrew
│   └ Brewfile
├ vim
│   └ .vimrc.symlink
└ zsh
    ├ install.sh
        └ .zshrc.symlink
```
Files you intend to share across all systems should be organized at the root level of your dotfile directory.  

Files that you only want to use on some systems can be organized under the bundles directory. Each directory under bundles/ will act as a bundle, and can be optionally installed by specifying it at the command line, or in the .medkit config file.

In the example directory structure above, you can see 3 instances of the Brewfile.  homebrew/Brewfile exists at the root level, and will install software every time medkit is run.  Under the bundles directory, there is a Brewfile for the go bundle, and another for the macos bundle.  The latter two brewfiles will only be run if specifically requested.

### Initialization
First, we need to set up a ~/.medkit config file specifying the location of our dotfiles (defaults to ~/dotfiles), and any bundles we would like to enable on this machine.
```yaml
dotfiles-directory: /home/marvin/dotfiles
active-bundles: [go,macos]
```

If you already have a dotfiles repo, clone it into dotfiles directory specified in your config file, and ensure that it follows the conventions described above.

If you do not yet have your own dotfiles, let medkit help you create one.  Create and init your dotfiles directory:
```sh
mkdir /home/marvin/dotfiles
medkit init
```

### TODO: Adding Dotfiles
Now, let's add an existing dotfile to MEDKIT. Of course, you'll replace the example path shown below with something real
on your computer.
```sh
medkit add dotfile -f /home/marvin/.vimrc -d vim
```

This command will:

- Create a new folder in your dotfiles repo (if it doesn't already exist) named `vim`
- Move the specified dotfile into the folder
- Symlink the dotfile back to the original location

You can view all of the dotfiles MEDKIT is managing like this.
```sh
medkit get dotfiles
```
At this point you're ready to make your repo available in your other environments. MEDKIT is agnostic when it comes to
how you do this, but using something like GitHub is highly recommended.

TODO: Provide basic GitHub instructions.

### Install Dotfiles
With your MEDKIT repo now under version control, you can use your repo anywhere. Let's say you have a new computer. Just
clone your GitHub repository to the new computer and install your dotfiles.
```sh
medkit install dotfiles
```
That's it!

### Updating Dotfiles
When you make a change to a dotfile, you'll probably want to make that change available on all your computers. The
specifics of how to synchronize your MEDKIT repo will depend on the tools you choose. If you're using GitHub, for
example, you'll need to make your changes, push them to GitHub, then pull them on all your other computers. No matter
which process you prefer, your dotfiles will be automatically updated.

## How to Contribute
Contributions are welcome! Check out [this link][contributing] on how you can help!

## Release History
MEDKIT's full release history can be found [here][changelog].

## Credits
* [Zach Holman][zach-holman-github-url] and his excellent ["dotfiles are meant to be forked" blog post][zach-holman-blog-url]

[zach-holman-github-url]: zach@zachholman.com
[zach-holman-blog-url]: https://zachholman.com/2010/08/dotfiles-are-meant-to-be-forked/

## License

Distributed under the MIT license. See the [LICENSE][license] file for more information.

<!-- Markdown link & img definitions -->
[environments]: #environments
[homebrew-bundle]: https://coderwall.com/p/afmnbq/homebrew-s-new-feature-brewfiles
[changelog]: https://github.com/LeadPipeSoftware/medkit/blob/master/CHANGELOG.md
[authors]: https://github.com/LeadPipeSoftware/medkit/blob/master/AUTHORS.md
[contributing]: https://github.com/LeadPipeSoftware/medkit/blob/master/CONTRIBUTING.md
[security]: https://github.com/LeadPipeSoftware/medkit/blob/master/SECURITY.md
[license]: https://github.com/LeadPipeSoftware/medkit/blob/master/LICENSE
[github-issues-image]: https://img.shields.io/github/issues/badges/shields.svg
[github-issues-url]: https://github.com/LeadPipeSoftware/medkit/issues
[wiki]: https://github.com/LeadPipeSoftware/medkit/wiki
