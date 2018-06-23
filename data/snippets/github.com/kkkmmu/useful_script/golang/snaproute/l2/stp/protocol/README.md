It should me noted that the tests must be run as root as the tests use the pcap lib that needs root access

Tests can be run via go test

example:

go test -v 


As of 3/22/16 because tests use interface l0 to in inject and listen for STP packets how long it takes the tests to run can take longer.  I have seen the tests run in as little as 6 seconds but as high as 30 seconds.   Will create a secondary loopback interface specifically used by the test so that run time is more consistant

