# latest.cli

[![Build Status](https://travis-ci.org/arunvelsriram/latest.cli.svg?branch=master)](https://travis-ci.org/arunvelsriram/latest.cli)

A CLI to find the latest version of a Ruby Gem, Node module, Java JAR etc.

## Usage

### Prerequisites

On **Linux** `curl` dependencies are required.

##### Debian based Linux:

`apt-get install curl`

##### RPM based Linux:

`yum install curl`

### Install using Go

```
go get github.com/arunvelsriram/latest.cli/latest
```

### Install from releases

#### Linux

Based on your Linux variant download `deb` or `rpm` file from the [latest release](https://github.com/arunvelsriram/latest.cli/releases/latest).

##### Debian based Linux:

`dpkg -i latest.cli.deb`

##### RPM based Linux:

`rpm -i latest.cli.rpm`

#### MacOS and Linux (using binary)

1. Download tar file from the [latest release](https://github.com/arunvelsriram/latest.cli/releases/latest)
2. Extract the tar file
3. Move the binary to `/usr/local/bin`
4. Run `latest`

#### Windows

1. Download zip file from the [latest release](https://github.com/arunvelsriram/latest.cli/releases/latest)
2. Extract the zip file
3. Open the extracted directory in Command Prompt
4. Run `latest`
5. Optinally add `latest.exe` to system path

## TODO
- [ ] Make latest installable through brew
- [ ] Support Python packages
- [ ] Findest latest version of a JAR  
- [ ] Pass multiple names (eg. `latest rails puma rspec`)
