FROM golang:1.16.1-alpine AS build
ARG DEST_DIR="/go/src/twrapper"
COPY . $DEST_DIR
WORKDIR $DEST_DIR
RUN apk update && apk add git openssh
RUN go build -o /usr/bin/twrapper ./cmd/twrapper

FROM alpine:3
COPY --from=build /usr/bin/twrapper /bin/twrapper
ENTRYPOINT [ "/bin/twrapper" ]