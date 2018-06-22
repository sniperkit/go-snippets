module MyVagrantPlugin
    class Plugin < Vagrant.plugin("2")
        name "my vagrant plugin"
        command "free-memory" do
            require_relative "my-vagrant-plugin/command"
            Command
        end
    end
end