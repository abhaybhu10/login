FROM golang:1.16-alpine

COPY . /go/src/

WORKDIR /go/src/command/
RUN CGO_ENABLED=0 go build -o /bin/app


CMD ["app"]