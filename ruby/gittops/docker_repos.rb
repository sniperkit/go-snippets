#!/usr/bin/env ruby

require_relative 'gittops'
require 'json'
include GitTops

class DockerRepo
  @token = ''
  @query = 'stars:>1 openstack'
end

repos = GitTopRepos.new(YOUR_GIT_TOKEN, 'stars:>1 docker', 10)
repos.get_repo.each do |repo|
  puts repo.name + ": " + repo.html_url
end
