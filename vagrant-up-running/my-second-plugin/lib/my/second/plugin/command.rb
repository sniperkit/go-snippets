require "my/second/plugin/version"

module My
  module Second
    module Plugin
      class Command < Vagrant.plugin("2", "command")
        def execute 
          puts "Hello Again! Version: " + My::Second::Plugin::VERSION
          return 0
        end 
      end
    end
  end
end