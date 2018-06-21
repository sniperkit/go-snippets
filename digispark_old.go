/*
<!--
Copyright (c) 2017 Christoph Berger. Some rights reserved.

Use of the text in this file is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

Use of the code in this file is governed by a BSD 3-clause license that can be found
in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2017-05-21"
draft = "false"
+++


# NOTE - this is the previous version of the installation part of https://appliedgo.net/digispark.

**This is outdated information that only serves documentation purposes.**

<!--more-->


## Controlling a Digispark board from Go code

With a few more (software) ingredients, the I/O ports can be controlled from Go code running on a PC/Mac/LinuxBox/etc. These are:

* LittleWire
* The `micronucleus` cli app
* Gobot

[LittleWire](http://littlewire.github.io) is basically a script that adds USB communication capabilities to the Digispark. With an accompanying library at the USB host's end, the board's I/O ports can be remote-controlled from an app running on the host.

[`micronucleus`](https://github.com/micronucleus/micronucleus) is a bootloader that comes preinstalled with Digisparks. (It needs a little space in the microcontroller's flash memory, leaving about 6KB memory available for scripts.) It has a companion CLI app that is needed for installing LittleWire on the device.

[Gobot](https://gobot.io/) is a Go robotics framework that connects Go apps to a large array of electronic devices, from the little Digispark to Arduino, Raspberry Pi, and even Quadrocopters.

And here is how I used all this to make the onboard status LED blink, and a servo motor move.

## The steps


### Step 1: Install the micronucleus CLI app

I compiled the `micronucleus` CLI tool from its source in two steps:


#### Clone the source

A simple `git clone` fetches the source code from GitHub:

    git clone https://github.com/micronucleus/micronucleus


#### Compile the command

After this, I cd'ed into `micronucleus/commandline` and ran `make`:

```
commandline $ make
Building library: micronucleus_lib...
gcc -I/usr/local/Cellar/libusb-compat/0.1.5_1/include -L/usr/local/Cellar/libusb-compat/0.1.5_1/lib -lusb -Ilibrary -O -g -D MAC_OS -c library/micronucleus_lib.c
```

On my Mac, I got a couple of warnings but as the `micronucleus` command turned out to work fine, I consider these warnings as benign. On my Lubuntu netbook, `make` finished without any warnings.


### Step 2: Install LittleWire

Again, only two steps are required.

#### Download littlewire_13.hex

The "firmware" script LittleWire is available at [the LittleWire download page](http://littlewire.github.io/downloads.html)


#### Install LittleWire on the device

Then I ran `micronucleus`:

    commandline $  ./micronucleus --run ../../LittleWire/littlewire_v13.hex

(I had LittleWire in a directory next to the Git workspace of micronucleus; if you want to repeat these steps, adapt the path as needed.)

`micronucleus` then asked to plug in the device into an USB port. (The Digispark, when powered on, listens on the USB port for some seconds, and if the host does not send anything, the device disables the USB port, because two of the I/O pins have a double use as USB channels, and disabling the USB port makes these two ports available for hardware wired to the board. This is why the board must be plugged in only when the host app is ready.)

I plugged in the board, and `micronucleus` uploaded the LittleWire firmware in a few seconds.


### Step 3: Install Gobot

For this, I simply followed the instructions from the [Gobot documentation](https://gobot.io/documentation/platforms/digispark/).

For Mac and Linux, some USB libraries need to be installed. On the Mac, this is a one-liner if you have [Homebrew](https://brew.sh) installed.

    brew install libusb && brew install libusb-compat

On Ubuntu (and surely also on other Debian-based distributions), a simple `apt-get` does the trick.

    sudo apt-get install libusb-dev

Now I was ready to fetch Gobot and install the Digispark platform package. (If you are new to Go: the ellipsis at the end of the `go get` line advise the `go get` tool to also download all subprojects, even if the top project does not import them. The `-d` is also important as this prevents the `go get` tool from installing anything at this point. The only thing we want to install is the Digispark package.)

      go get -d -u gobot.io/x/gobot/...

The next step installs the `digispark` library. I first ran into an error ("undefined: lw") until I realized that I had `cgo` disabled. On a Bash shell, the prefix `CGO_ENABLED=1` ensures that go install uses `cgo`.

      CGO_ENABLED=1 go install gobot.io/x/gobot/platforms/digispark

While this worked on the Mac, on Linux I got a "package libusb was not found" error message. Apparently the error was raised by a tool named `pkg-config`. I wasn't able to figure out the reason for this, so my Linux-based test ended here.

*/