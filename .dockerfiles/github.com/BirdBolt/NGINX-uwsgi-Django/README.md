# alpine-nginx-uwsgi-django
A multi container setup for a nginx-uwsgi-django stack

####Will eventually be adding more details/instructions concerning use of this

### How to deploy your application
Copy the root of your Django project into the app folder, make sure to include a `requirements.txt` file. 
You need to modify the following files:

- [uwsgi.ini]
...
- [nginx-app.conf]
...


Then:

```
$ docker-compose build
$ docker-compose up
```

If that succeds, then `curl localhost` (or navigating to localhost in a browser)
should return the output of the present project, be it the default django page, 
or your own django project you copied in. 
