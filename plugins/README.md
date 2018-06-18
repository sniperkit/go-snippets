In the following application we have:

1. a main process that gets events and instructs the Boss
2. a Boss that managers the workers until:
    1. Work is done
    2. Process gives Boss new work
3. Workers/Plugins that work until:
    1. They are Finished
    2. Boss tells them quit
