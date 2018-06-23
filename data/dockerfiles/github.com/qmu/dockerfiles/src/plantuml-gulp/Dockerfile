FROM qmu1/plantuml:latest

MAINTAINER TAMURA Yoshiya <a@qmu.jp>

WORKDIR /workspace

RUN \
    apk add --no-cache nodejs \
    && npm install -g gulp \
    && npm install --save-dev gulp gulp-coffee gulp-plantuml gulp-cached gulp-plumber

RUN mkdir -p /workspace/uml \
    && mkdir -p /workspace/img

ADD plantuml.sh /bin/plantuml
ADD gulpfile.coffee /workspace/gulpfile.coffee

CMD []
ENTRYPOINT []
