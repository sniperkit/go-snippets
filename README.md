# docker-gradle-alpine
Dockerfile for Gradle

## Build

```sh
$ cd /path/to/gradle/project

$ docker run \
  -v $PWD:/app \
  -w /app \
  -u $UID:$GID \
  takitake/gradle-alpine \
  gradle wrapper && ./gradlew clean build
```
