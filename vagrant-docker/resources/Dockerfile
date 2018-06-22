FROM java:8u111-jdk

COPY spring-boot-sample-web-jsp-2.0.1.BUILD-SNAPSHOT.war /app/app.war
COPY entrypoint.sh /app/entrypoint.sh
RUN ["chmod", "+x", "/app/entrypoint.sh"]
EXPOSE 8080 9010
ENV SERVER_IP="localhost"
WORKDIR /app
ENTRYPOINT ["/app/entrypoint.sh"]