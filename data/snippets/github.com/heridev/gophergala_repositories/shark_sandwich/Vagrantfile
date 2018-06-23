# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  # All Vagrant configuration is done here. The most common configuration
  # options are documented and commented below. For a complete reference,
  # please see the online documentation at vagrantup.com.

  config.vm.box = "hashicorp/precise64"

  config.vm.synced_folder ".", "/home/vagrant/gocode/src/shark-sandwich"

  config.vm.provision "shell", inline: "sudo apt-get -qq update"
  config.vm.provision "shell", inline: "sudo apt-get -y -qq install git mercurial make cmake pkg-config"
  config.vm.provision "shell", privileged: false, inline: "hg clone -u go1.4.1 https://code.google.com/p/go"
  config.vm.provision "shell", privileged: false, inline: "cd go/src && ./make.bash"
  config.vm.provision "shell", privileged: false, inline: "echo export GOPATH=$HOME/gocode >> $HOME/.profile"
  config.vm.provision "shell", privileged: false, inline: "echo export PATH=$HOME/gocode/bin:$HOME/go/bin:$PATH >> $HOME/.profile"
  config.vm.provision "shell", inline: "chown -R vagrant:vagrant /home/vagrant/gocode"
  config.vm.provision "shell", privileged: false, inline: "go get github.com/libgit2/git2go"
  config.vm.provision "shell", privileged: false, inline: "cd $HOME/gocode/src/github.com/libgit2/git2go && git submodule update --init && make install"

end
