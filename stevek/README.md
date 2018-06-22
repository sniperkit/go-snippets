# Aims

I wanted to build a representation of a computer network and, given configuration files for each network device, establish where messages could reach. For example, what ports can the User access on the Linux host?
```
       User
         +
         |
         v
      Internet
         +
         |
         |
   +-----+-----+         +------------+
   |           |         |            |
   |  Firewall +---------+ Linux host |
   |           |         |            |
   +-----------+         +------------+
```

The intention was that by using channels and goroutines this could be achieved though simulation and then extended to other uses.

My hope was that I could do something functional with a nice web frontend.

# Where did I get to?

The usual - not as far as would've been nice...

What was done:
 + Defining pipelines of channels directly in code
 + Filtering messages passed through the pipeline
  + With very limited and messy parsing of example configuration lines from real-world devices 
 + Putting objects onto a channel to describe result of filtering
 + Using HTML5 Server-Sent events to send the results to a web browser

What was not done:
 + Making this anything more than a one-shot demo
 + Idiomatic go
  + There is plenty of tidying to do and more interfaces would really help
 + Nice in-browser visualisation
 + Use of a graph representation of a network
 + Complete (or semi-complete) parsing of config files

# References/resources

These are things that I think may be either relevant background or relevant in terms of implementation design.

Background:
 + I remember that *ipchains* could check things for you, something I think was left out of *iptables*
 > **Using ipchains**
 > ...
 > The final (and perhaps the most useful) function allows you to check what would happen to a given packet if it were to traverse a given chain.

 (from http://www.tldp.org/HOWTO/IPCHAINS-HOWTO-4.html)
 + Cisco ASAs have a *packet-tracer* function that is fantastic for showing the route that a packet would take through the device, including firewalling decisions, VPN establishment, etc.
  + https://supportforums.cisco.com/document/29601/troubleshooting-access-problems-using-packet-tracer
 + The *--dry-run* option to the *patch* tool
 > **--dry-run**
 > Print the results of applying the patches without actually changing any files.

 (from man pages e.g. http://unixhelp.ed.ac.uk/CGI/man-cgi?patch+1)

Design/implementation:
 + When I first saw the [Cayley](https://github.com/google/cayley) graph database [visualisation](http://cayley-graph.appspot.com/ui/visualize) it looked like the kind of thing that could really help with this project
 + When I saw [Andrew Gerrand](https://twitter.com/enneff)'s [Sigourney](https://github.com/nf/sigourney) modular synthesizer I could see many parallels with what I was wanting to do.
 + Blog on piplines in go: http://blog.golang.org/pipelines
 + Functional programming: http://golang.org/doc/codewalk/functions/

# Future

A nice way of interacting through a browser would be great.

Also, I think a graph database would be a very useful underlying structure for representing the network and channel structure.

Other things:
 + Use *cgo* to actually use *iptables*/*netfilter* libraries to parse *iptables* rules
 + The app should suggest an *nmap* or *telnet* command that the user could run to verify
 + The app should suggest a firewall rule to increase depth of security (belt-and-braces)  
