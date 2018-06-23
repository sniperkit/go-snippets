#!/usr/bin/env ruby
require "thor"

class Build < Thor
  desc "binary", "Builda go binary za linux amd64"
  def binary(name)
    puts "#{name}: building binary za linux amd64"
    path_cmd("../03-consul/#{name}", "env GOOS=linux GOARCH=amd64 go build")
    path_cmd(".", "mv ../03-consul/#{name}/#{name} ./images/#{name}/.")
  end

  desc "binary_all", "Builda sve go binaries"
  def binary_all()
    binary('sensor')
    binary('worker')
    binary('app')
  end

  desc "image", "Builda docker image"
  def image(name)
    puts "#{name}: building docker image"
    path_cmd("./images/#{name}", "docker build -t #{name} .")
  end

  desc "image_all", "Builda sve docker images"
  def image_all()
    image('sensor')
    image('worker')
    image('app')
  end

  private
  def path_cmd(path, cmd)
    Dir.chdir(path) {
      puts "#{path} :: #{cmd}"
      system cmd
    }
  end

end

Build.start(ARGV)
