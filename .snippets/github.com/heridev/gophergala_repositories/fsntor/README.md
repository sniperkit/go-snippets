fsntor
=======

![first logo to be able to submit the project to gopher-gala.challengepost.com](https://raw.githubusercontent.com/gophergala/fsntor/master/assets/logo-readme.png)

## Purpose

The concept is building a binary that allows to pass a configuration file with different file/directory patterns which will be watched through [fsnotify](https://github.com/go-fsnotify/fsnotify) and executing a command which can be any available in PATH.

There are a bunch of tools that do this implemented in different languages but I would to have a binary without any environment as a dependency (e.g. Ruby, node, etc.). I will refer this binary as `fsntor` henceforth.

There are also implementations in GO, however my main goal is to use [fsnotify](https://github.com/go-fsnotify/fsnotify) to be more compatible through different Operative Systems.

As a kickstart of the project I came up with this brief features to implement insofar as possible, however after thinking for a while and provided the first stuff I got two groups, one for the first alpha version and the second one to implement after it.

1. To implement for the first working alpha version:
  * Define several file/directory patterns to "watch" each one with the command to execute.
  * Execute a command which doesn't end with a result and keeps running until `fsntor` process ends or files match by the pattern and action requires to restart the process. I will refer them as long-processes.
  * Ability to only execute the command under some events (e.g. File change, file removed, etc.).
  * ~~Ability to define what signal to send when long-process must restart and define timeout to wait for it finishes otherwise kill it.~~ Discarded due Go only have support for [two signals](http://golang.org/pkg/os/#Signal), then no choice.

2. To implement after first alpha version:
  * Ability to specify if the command must be executed when `fsntor` starts.
  * Other features that I know that are supported for this Go packages:
    * [romanoff/gow](https://github.com/romanoff/gow)
    * [parkghost/watchf](https://github.com/parkghost/watchf)

Not really sure if all of these things are  already supported for the different GO package (listed above) and other similar tools in other languages which I've seen, however the project has started mainly for two reasons:

* Learning purpose
* A [Gopher Gala](http://gophergala.com/) Project (January 2015), however with aim to carry on after it.

## License

Just MIT, read the [LICENSE](LICENSE) file for more information.
