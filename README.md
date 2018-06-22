# alpine-tomcat

1. base on alpine-linux
2. with java8 and tomcat8
3. timezone is UTC+8
4. it is very small

##Use case
```
docker run -d -p 8080:8080 -v app.war:/app/webapps/app.war lioncui/alpine-tomcat
```

> if u want to write it into dockerfile


```
FROM lioncui/alpine-tomcat    

ADD app.war /app/webapps/
```

Ths all
