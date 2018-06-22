                                                                          


#Sends a request - this need to be executed from the Raspberry Pi or any remote devices 

import httplib, urllib
import time
while True:
     conn = httplib.HTTPConnection("192.168.1.9:8079")
     conn.request("GET", "/devices/001-iamhash1234")
     response = conn.getresponse()
     print response.status, response.reason
     time.sleep(60)
     conn.close()









