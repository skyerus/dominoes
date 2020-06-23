# Compile stage
FROM golang:1.14.4-alpine3.12 AS build-env
ENV CGO_ENABLED 0

RUN apk add --no-cache git
RUN apk add --update nodejs npm

ADD . /
WORKDIR /
RUN cd ./web; npm run build
WORKDIR /

RUN go build -gcflags "all=-N -l" -o /dominoes

# Final stage
FROM alpine:3.12.0

ENV CGO_ENABLED 0
RUN apk add --no-cache libc6-compat
RUN apk add --no-cache bash

WORKDIR /
COPY --from=build-env /dominoes /
COPY --from=build-env /web/dist /web/dist
CMD ["/dominoes"]