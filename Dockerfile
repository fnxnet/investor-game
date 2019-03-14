FROM golang:alpine
ADD . /go/src/app
WORKDIR /go/src/app
RUN apk add --no-cache git mercurial && go get -u github.com/google/uuid && go get -u github.com/gorilla/websocket && go get -u github.com/go-redis/redis
ENV PORT=3001
CMD ["go", "run", "*.go --front"]