# multi-stage build
# build image
FROM golang:1.14-alpine AS build-go

# Get bash and git
RUN apk update && apk add --no-cache git

COPY . /go/src/sitesearch
WORKDIR /go/src/sitesearch

RUN go build -o bin/sitesearch

# deployment image
FROM alpine

COPY --from=build-go /go/src/sitesearch/bin/sitesearch /sitesearch


EXPOSE 8080
ENTRYPOINT ["/sitesearch"]