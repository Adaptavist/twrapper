FROM golang:1.16.1-alpine AS build
ARG DEST_DIR="/go/src/twrapper"
COPY . $DEST_DIR
WORKDIR $DEST_DIR
RUN apk update && apk add git openssh
RUN go build -o /usr/bin/twrapper ./cmd/twrapper

FROM golang:1.16.1-alpine AS tfenv
ARG DEST_DIR="/go/src/tfenv"
ARG TFENV_VERSION="v2.2.0"
RUN apk update && apk add git openssh && \
    git clone https://github.com/tfutils/tfenv.git ${DEST_DIR}
WORKDIR $DEST_DIR
RUN git checkout tags/${TFENV_VERSION} -b ${TFENV_VERSION}

FROM alpine:3
COPY --from=build /usr/bin/twrapper /bin/twrapper
COPY --from=tfenv /go/src/tfenv /opt/tfenv
RUN apk update && apk add bash curl git && \
    ln -s /opt/tfenv/bin/terraform /usr/bin && \
    ln -s /opt/tfenv/bin/tfenv /usr/bin
ENTRYPOINT [ "/bin/twrapper" ]