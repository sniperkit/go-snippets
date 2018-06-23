FROM alpine

RUN adduser -D -g '' jtobe

ADD app app
ADD seed.outputs.json seed.outputs.json
ADD seed.png.metadata.json seed.png.metadata.json

RUN chmod -R 777 app
RUN chmod 777 seed.outputs.json
RUN chmod 777 seed.png.metadata.json

ENTRYPOINT ["app/extract.sh"]
