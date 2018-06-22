# lib/vagrant-ls.rb
require 'bundler'

begin
  require 'vagrant'
rescue LoadError
  Bundler.require(:default, :development)
end

require 'vagrant-ls/plugin'
require 'vagrant-ls/command'