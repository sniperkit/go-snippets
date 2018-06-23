from __future__ import print_function
import httplib2
import os
import sys
import email
import base64
import csv

from apiclient import discovery
from oauth2client import client
from oauth2client import tools
from oauth2client.file import Storage

#try:
#    import argparse
#    flags = argparse.ArgumentParser(parents=[tools.argparser]).parse_args()
#except ImportError:
#    flags = None

# If modifying these scopes, delete your previously saved credentials
# at ~/.credentials/gmail-python-quickstart.json
SCOPES = 'https://www.googleapis.com/auth/gmail.readonly'
CLIENT_SECRET_FILE = 'client_secret.json'
APPLICATION_NAME = 'Gmail API Python Quickstart'
user_id = 'me'


def get_credentials():
    """Gets valid user credentials from storage.

    If nothing has been stored, or if the stored credentials are invalid,
    the OAuth2 flow is completed to obtain the new credentials.

    Returns:
        Credentials, the obtained credential.
    """
    home_dir = os.path.expanduser('~')
    credential_dir = os.path.join(home_dir, '.credentials')
    if not os.path.exists(credential_dir):
        os.makedirs(credential_dir)
    credential_path = os.path.join(credential_dir,
                                   'gmail-python-quickstart.json')

    store = Storage(credential_path)
    credentials = store.get()
    if not credentials or credentials.invalid:
        flow = client.flow_from_clientsecrets(CLIENT_SECRET_FILE, SCOPES)
        flow.user_agent = APPLICATION_NAME
        if flags:
            credentials = tools.run_flow(flow, store, flags)
        else: # Needed only for compatibility with Python 2.6
            credentials = tools.run(flow, store)
        print('Storing credentials to ' + credential_path)
    return credentials

def main():
    """Shows basic usage of the Gmail API.

    Creates a Gmail API service object and outputs a list of label names
    of the user's Gmail account.
    """
    credentials = get_credentials()
    http = credentials.authorize(httplib2.Http())
    service = discovery.build('gmail', 'v1', http=http)

    #response = service.users().messages().list(userId=user_id, labelIds='Label_5').execute()
    response = service.users().messages().list(userId=user_id, q='label:inbox AND NOT label:'+sys.argv[2]).execute()
    messages = []

    if 'messages' in response:
        messages.extend(response['messages'])

    while 'nextPageToken' in response:
        print(messages)
        page_token = response['nextPageToken']
        #response = service.users().messages().list(userId=user_id, labelIds='Label_5', pageToken=page_token).execute()
        response = service.users().messages().list(userId=user_id, q='label:inbox AND NOT label:'+sys.argv[2], pageToken=page_token).execute()
        messages.extend(response['messages'])
        if len(messages) >= 1000:
            break


    f = open(sys.argv[1], 'wt')
    try:
        writer = csv.writer(f)
        #writer.writerow( ('Domain Name', 'Subject') )

        for message in messages:
            m = service.users().messages().get(userId=user_id, id=message['id'],
                                                 format='raw').execute()
            msg_str = base64.urlsafe_b64decode(m['raw'].encode('ASCII'))
            msg = email.message_from_string(msg_str)
            msg_from = msg.get('from')
            writer.writerow( (msg_from.split('@')[1].split('.')[0], msg.get('subject'), "not_hr") )
            print( (msg_from.split('@')[1].split('.')[0], msg.get('subject')) )
#            for part in msg.walk():
#                if part.get_content_type() == 'text/plain':
#                    print(part.get_payload())
    finally:
        f.close()
        


if __name__ == '__main__':
    main()
