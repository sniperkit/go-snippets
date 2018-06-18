from apiclient.discovery import build
from httplib2 import Http
from oauth2client import file, client, tools

SCOPES = "https://www.googleapis.com/auth/gmail.readonly"
CLIENT_SECRET = "client_secret.json"

store = file.Storage('storage.json')
creds = store.get()

if creds is None or creds.invalide:
    flow = client.flow_from_clientsecrets(CLIENT_SECRET, SCOPES)
    creds = tools.run(flow, store)
GMAIL = build('gmail', 'v1', http=creds.authorize(Http()))

threads = GMAIL.users().threads().list(userId='me', q='hao.1.wang@gmail.com').execute().get('threads', [])
for thread in threads:
    tdata = GMAIL.users().threads().get(userId='me', id=thread['id']).execute()
    print tdata['messages']