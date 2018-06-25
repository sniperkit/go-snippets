#!/usr/bin/env ruby
# filename: ~/config.rb

require 'inifile'

ini = IniFile::load('/etc/NetworkManager/NetworkManager.conf')
ini['main']['dns'] = 'dnsmasq'
