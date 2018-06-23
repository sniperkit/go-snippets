#!/usr/bin/env bash

# Ask for the administrator password upfront
sudo -v

# Keep-alive: update existing `sudo` time stamp until has finished
while true; do sudo -n true; sleep 60; kill -0 "$$" || exit; done 2>/dev/null &

###
# Install stuff
###
sudo dnf install cowsay fortune-mod ruby rubygems ruby-devel cmake ninja-build vim python3-devel python-devel llvm llvm-libs clang boost-devel boost-python3-devel vim htop zsh libreoffice-langpack-de thermald rigprep exa -y
sudo dnf groupinstall "Development Tools" "C Development Tools and Libraries" -y
sudo dnf install automake gcc gcc-c++ glib-devel dbus-glib-devel libxml2-devel -y

gem install lolcat

# TRIM
sudo systemctl enable --now fstrim.timer

###
# RPMfusion
###

# Add repos
sudo dnf install "https://download1.rpmfusion.org/free/fedora/rpmfusion-free-release-$(rpm -E %fedora).noarch.rpm" -y
sudo dnf install "https://download1.rpmfusion.org/nonfree/fedora/rpmfusion-nonfree-release-$(rpm -E %fedora).noarch.rpm" -y

# Install stuff
sudo dnf install -y ffmpeg

###
# Google Chrome
###
sudo rpm --import https://dl-ssl.google.com/linux/linux_signing_key.pub
sudo dnf -y install "https://dl.google.com/linux/direct/google-chrome-stable_current_$(uname -i).rpm"

###
# TeX stuff
###
sudo dnf install texlive-scheme-medium texlive-minted texstudio -y

###
# Node.js
###
sudo dnf install nodejs nodejs-devel npm -y

sudo rpm --import https://dl.yarnpkg.com/rpm/pubkey.gpg
sudo wget https://dl.yarnpkg.com/rpm/yarn.repo -O /etc/yum.repos.d/yarn.repo
sudo dnf install yarn

# Get zsh plugins
sudo npm i -g zsh-goggles
zgoggles

###
# Syncthing
###

# Install Syncthing
sudo dnf install -y syncthing syncthing-gtk syncthing-inotify nautilus-python

# Enable and start Syncthing
systemctl enable --user --now syncthing.service

###
# PlayOnLinux
###
sudo dnf install liberation-fonts wine wine-mono wine-fonts mingw32-wine-gecko mingw64-wine-gecko -y
sudo dnf install playonlinux -y
sudo sed -i 's/windows_add_size = 25/windows_add_size = 250/' /usr/share/playonlinux/python/lib/Variables.py

