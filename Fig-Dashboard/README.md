# Fig-Dashboard

This is a Fig Dashboard web app written in go for the GopherGala 2015 hackathon

NAME:
   figdash - fig dashboard

USAGE:
   figdash [global options] command [command options] [arguments...]

VERSION:
   0.0.1

COMMANDS:
   kill		Force stop service containers.
   ps		List containers
   rm		Remove stopped service containers.
   start	Start existing containers for a service.
   stop		Stop existing containers without removing them.
   web		Enable web monitoring on http://localhost/1984/PROJECTNAME.
   help, h	Shows a list of commands or help for one command
   
GLOBAL OPTIONS:
   --verbose			Show more output
   --file, -f "fig.yml"		Specify an alternate fig file [$FIG_FILE]
   --projectname, -p "notset"	Specify an alternate project name [$FIG_PROJECT_NAME]
   --help, -h			show help
   --version, -v		print the version
