#!/usr/bin/env python

import sys, urllib2, base64, json, getpass

#-------------------------------------------

GITHUB_API = "https://api.github.com"

#-------------------------------------------

def auth (): 
  result=''
  username = raw_input("Username:")
  password = getpass.getpass()
  try:
    request = urllib2.Request( GITHUB_API + "/user" )
    base64string = base64.encodestring('%s:%s' % (username, password)).replace('\n', '') 
    request.add_header("Authorization", "Basic %s" % base64string) 
    response = urllib2.urlopen(request) 
    print( response.info() )
    result = response.read()
    response.close()
  except urllib2.HTTPError as httpe:
    print 'Github Authentication error: %s' % str(httpe)
    print 'Please check your user name and/or password.'
    sys.exit(0)
  except Exception as e:
    print 'Oops. We had a slight problem with the GitHub SSO: ' + str( e ) 
    sys.exit(0)
  print( result )
  json_dict = json.loads( result )
  public_gists = json_dict['public_gists']
  private_gists = json_dict['private_gists']
  print 'You have %i Private Gists and %i Public Gists' % (private_gists, public_gists)
  
if __name__ == "__main__":
  auth()
  
  
