FROM linterhub/image-ruby
LABEL maintainer="linterhub@repometric.com"
RUN apk add --no-cache libcurl && \
    apk add --no-cache --virtual .build-deps \
    make \
    ruby-dev \
    gcc \
    musl-dev \
    libffi-dev && \
    gem install csvlint && \
    apk del .build-deps
