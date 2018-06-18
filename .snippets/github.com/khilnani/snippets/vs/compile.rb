#!/usr/bin/ruby

def compile (expr="*")
  puts "Preparing ..." 

  system "rm -rf classes"
  system "mkdir classes"

  cmd = String.new("")

  if expr == "*"
    cmd = "scalac -d classes src/*.scala"
  elsif expr.end_with? ".scala"
    cmd = "scalac -d classes #{expr}"
  else
    cmd = "scalac -d classes #{expr}.scala" 
  end

  cmd += " -feature"

  puts "Compiling: #{cmd}"
  system cmd

  puts "Done."
end

if ARGV.length == 0
  puts "No File specified. Using *"
  compile
else
  compile ARGV[0]
end
