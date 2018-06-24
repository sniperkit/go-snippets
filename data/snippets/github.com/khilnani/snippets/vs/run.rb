#!/usr/bin/ruby

def run (expr)
  cmd = String.new("")

  if expr.end_with? ".scala"
    expr = expr.split('.')[0]
  end
  if expr.include? "/"
    expr = expr.split('/').last
  end

  cmd = "scala -classpath classes #{expr}"

  puts "Running: #{cmd}"
  system cmd
end

if ARGV.length == 0
  puts "USAGE: #{__FILE__} CLASSNAME"
else
  run ARGV[0]
end
