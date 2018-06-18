#!/usr/bin/env python

import sys, requests, json, getpass

#-------------------------------------------

GITHUB_API = "https://api.github.com"

#-------------------------------------------

def auth (): 
  result=''
  username = raw_input("Username:")
  password = getpass.getpass()
  try:
    request = requests.get( GITHUB_API + "/user", auth=(username, password) )
    if request.status_code != 200:
      print 'Github Authentication error: %s' % str(request.status_code)
      print 'Please check your user name and/or password.'
      sys.exit(0)
    json_dict = json.loads( request.text )
    public_gists = json_dict['public_gists']
    private_gists = json_dict['private_gists']
    print 'You have %i Private Gists and %i Public Gists' % (private_gists, public_gists)
  except Exception as e:
    print 'Oops. We had a slight problem with the GitHub SSO: ' + str( e ) 
    sys.exit(0)
    
#-------------------------------------------
  
if __name__ == "__main__":
  auth()
  
  
