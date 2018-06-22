package main

import "fmt"

func main() {
	handler := &MemHandler{
		Processes: make([]Process, 256),
	}
	handler.getProcesses()


	/*
	 * Dumping processes and their modules for the fun of it.
	 */
	for _, v := range handler.Processes{
		if v.getPid() != 0{
			fmt.Println(v.getPid(), ":", v.getName())
		}

	 	v.getModules()
		for _, module := range v.ModuleList{
			if module.ProcessID != 0 {
				fmt.Println("\t\t", module.getName())
			}
	 	}



	}

}
