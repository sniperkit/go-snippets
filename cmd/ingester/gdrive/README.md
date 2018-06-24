# omniscience

## Google Drive

### Index Google Drive

Index Google Drive by running

```
$ go run main.go --google_service_account_file_path=<google_service_account.json>
```

from within the `cmd/omniscience_client/gdrive` directory.

### Generate a Google service account

1. Go to [cloud.google.com/console](http://cloud.google.com/console)
2. Click Create Project
3. Enter a project name and click Create
4. Once you have created your project click on the options menu icon in the top left corner of the page, then 'APIs & Services', then 'Credentials'
5. Click on 'Create credentials' and then 'Service account key'
6. Next, select 'New service account', name it anything and select 'Project' and then 'Viewer' as the role from the dropdown list, finally, select JSON as the key type and click 'Create'
7. Upon clicking 'Create', a JSON file will be downloaded; pass the path to this file as the value of the `--google_service_account_file_path` flag
7. Click on 'Manage service accounts' (on the right-hand side), select your new Service Account, click on the three dots on the right-hand side, and then select Edit
8. Tick the box "Enable GSuite Domain-wide Delegation" and click Configure consent screen
9. Click on the options menu icon in the top left corner again, choose 'APIs & Services' and then 'Dashboard'
10. Click on the Drive API and then click the blue Enable button

### Give the Google service account access to files

For the Google service account to have access to index your files you must share the relevant folders with the service account's `client_email`.
