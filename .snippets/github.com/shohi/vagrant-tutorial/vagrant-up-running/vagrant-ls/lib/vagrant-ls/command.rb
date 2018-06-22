module Vagrant
    module Ls
        class Command < Vagrant.plugin('2', :command)
            def execute
                exec('VBoxManage list vms')
                0
            end
        end
    end
end