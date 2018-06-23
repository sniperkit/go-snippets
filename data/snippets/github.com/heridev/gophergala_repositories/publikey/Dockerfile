FROM ubuntu
ADD publikey_linux_amd64 /
RUN mkdir /data
VOLUME /data
EXPOSE 8080
ENTRYPOINT /publikey_linux_amd64 server -p 8080 --data-file /data/data.db
