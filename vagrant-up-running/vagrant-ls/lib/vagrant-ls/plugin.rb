module Vagrant
    module Ls
        class Plugin < Vagrant.plugin('2')
            name "List"
            description <<-DESC
            This plugin 'ls' all of the vms running inside the machine
            DESC
    
            command 'ls' do
                require_relative 'command'
                Command
            end
        end
    end
end