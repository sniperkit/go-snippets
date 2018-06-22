module My
  module Second
    module Plugin
      class Plugin < Vagrant.plugin('2')
        name "my-second-plugin"
        
        command 'hello-again' do
          require_relative "plugin/command"
          Command
        end
      end
    end
  end
end
