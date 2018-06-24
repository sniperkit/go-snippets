#!/usr/bin/env ruby
require "thor"

class Build < Thor
  desc "binary", "Builda go binary"
  def binary(name)
    puts "#{name}: building binary"
    path_cmd(name, "go build")
  end

  desc "binary_all", "Builda sve go binaries"
  def binary_all()
    binary('sensor')
    binary('worker')
    binary('app')
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
