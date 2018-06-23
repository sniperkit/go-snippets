FROM golang:1.6.2-alpine

ENV PATCHELF_VER=0.9 \
    PATCHELF_URL=https://github.com/strothj/alpine-patchelf/releases/download/
ENV PROJECT_ROOT=/go/src/github.com/strothj/alpine-debtool
    
RUN apk add --no-cache bash ca-certificates libstdc++ && \
    wget ${PATCHELF_URL}/${PATCHELF_VER}/${PATCHELF_VER}.tar.gz && \
    tar -C / -xzf ${PATCHELF_VER}.tar.gz && \
    rm ${PATCHELF_VER}.tar.gz
ADD cmd $PROJECT_ROOT/cmd
ADD commands $PROJECT_ROOT/commands
ADD debtool $PROJECT_ROOT/debtool
ADD vendor $PROJECT_ROOT/vendor
WORKDIR $PROJECT_ROOT/cmd/debtool
RUN go install

CMD ["/bin/bash"]