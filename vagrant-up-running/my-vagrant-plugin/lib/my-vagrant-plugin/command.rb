module MyVagrantPlugin
    class Command < Vagrant.plugin("2", "command")
        def execute
            
            machine = @env.machine(:default, :virtualbox)
            if machine.state.id != :running
                @env.ui.error("Machine must be running.")
                return 1 
            end

            command = "grep MemFree /proc/meminfo | awk '{print $2}'" 
            machine.communicate.execute(command) do |type, data|
                @env.ui.info(data)
            end
            
            return 0
        end 
    end
end