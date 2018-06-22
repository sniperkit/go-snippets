# The bash-driver two main dependencies are:
#
# 1. IntelliJ Idea SDK
# 2. BashSupport plugin
#
# None of them have Maven repositories.  We can handle 2 easily by
# using jitpack.io with the Github url of BashSupport.  But for 1 things
# are a little bit more complicated.
#
# This makefile downloads, extract and install the required IntelliJ
# Idea jars into the vendor directory, so we can declare them as
# local dependencies in the build.gradle file.

# Configure these to point to the desired IntelliJ Idea version,
# tar.gz file and the URL to download it:
IDEA_VERSION := idea-IC-171.4694.23
IDEA_TAR_FILE := ideaIC-2017.1.4-no-jdk.tar.gz
IDEA_URL=https://download.jetbrains.com/idea/${IDEA_TAR_FILE}

# Configure IDEA_JARS to include the jars you need from the IntelliJ Idea
# distribution:
IDEA_JARS := util.jar
IDEA_JARS += idea.jar
IDEA_JARS += openapi.jar
IDEA_JARS += extensions.jar
IDEA_JARS += picocontainer.jar
IDEA_JARS += trove4j.jar
IDEA_JARS += guava-19.0.jar

# Don't modify any of these unless you know what you are doing.
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.SUFFIXES:
vendor_dir := vendor
idea_dir := ${vendor_dir}/${IDEA_VERSION}
idea_lib_dir := ${idea_dir}/lib
packaged_jars := $(addprefix ${IDEA_VERSION}/lib/, ${IDEA_JARS})
installed_jars := $(addprefix ${idea_lib_dir}/, ${IDEA_JARS})

.PHONY: all
all: ${installed_jars}

.PHONY: clean
clean:
	- rm -rf ${idea_dir}

%.jar: | ${vendor_dir}
	$(eval TMPDIR := $(shell mktemp -d))
	$(eval TMPFILE := ${TMPDIR}/${IDEA_TAR_FILE})
	curl --location --output ${TMPFILE} ${IDEA_URL}
	tar -xzf ${TMPFILE} -C ${vendor_dir} ${packaged_jars}
	rm -rf ${TMPDIR}

${vendor_dir}:
	mkdir $@
