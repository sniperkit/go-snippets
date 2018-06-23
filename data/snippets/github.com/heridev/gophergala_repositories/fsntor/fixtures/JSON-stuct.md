fsntor configuration file
============================================

## Current JSON structure

The current JSON structure configuration is clunky, but I have some reasons over other which may look simpler.

First of all it is an array to allow to define different actions for a set of patterns, I will refer them from now on "tasks", however it doesn't justify the awkward structure of a task.

A task could be pairs of one pattern and one command, however that would require to identify the patterns which have several tasks (i.e. they have an action per event) to avoid to register (`Watch`) the files and directories more than one to avoid performance issues, moreover it would have to group the commands under the same event for the same task, then it gets the implementation more complicated.

Discarding the one to one association mentioned in the previous paragraph we could define a pattern for a set of commands, i.e. actions as it does the current structure, but then the configuration could have same actions for different patterns, because they have to execute exactly the same, then I discarded it.

After all of this we can focus in the current first level structure ("patterns" and "actions" properties) and to argue the "action" structure.

Although an "action" could be a simple executor, however I think that is worthwhile to have an array of executors to be able to execute several actions, nonetheless it could be an simple array of commands (strings), but then we would loose the ability to tune the executor to have some management over the process created for the command execution, then an object is more flexible, moreover that I've already defined two properties to control the process time and termination.

Actions array define each action by the name(s) event(s), when an action is for several events, each event's name must be comma separated with not blank spaces; actions mustn't be duplicated.


## Execution flow matched with current configuration file

In the time being, the execution of the commands has been though to execute in series, it means, when finishes one, then execute the following one, then the two specified properties should have the following role:

* `timeout`: Set a maximum time to terminate the process if it doesn't do before this. It is a string to allow to support different time scale milliseconds (ms), seconds (s), etc. It will be parsed [time.ParseDuration](http://golang.org/pkg/time/#ParseDuration). By default, process run infinitely and it has to terminate itself, this applies if the value define a 0 or negative value.
* `waitFor`: Set the time to wait after Interrupt signal is sent, if the process doesn't terminate during this, then `fsntor` will kill the process. By default, the process can take an infinitely amount of time to terminate, the same applies if the value is negative and if it is 0 then kill the process straightaway without sending the Interrupt signal.
