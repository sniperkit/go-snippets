#!/usr/bin/env bash

# add entry
# deb https://mirrors.tuna.tsinghua.edu.cn/CRAN/bin/linux/ubuntu trusty/
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys E084DAB9

# Make sure the system is up-to-date
sudo apt-get update

# upgrade takes too much time
# apt-get upgrade -y 

# install java
# Add Java PPA
sudo apt-get install python-software-properties
sudo add-apt-repository ppa:webupd8team/java
sudo apt-get install oracle-java8-installer　oracle-java8-set-default

# Install development dependencies, dependencies needed for --check-as-cran,
# or libraries that R packages use (libcurl, gdal)
sudo apt-get build-dep -y r-base-core
sudo apt-get install -y libcurl4-openssl-dev libxml2-dev qpdf texlive-full build-essential gdal-bin libgdal-dev libproj-dev libxml2-dev git

# Install basic R
sudo apt-get install -y r-base r-base-core r-base-dev

# Permissions for R's package library
sudo adduser vagrant staff

# Install R packages
Rscript /vagrant/install-r-packages.R

