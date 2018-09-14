FROM golang:1.10-alpine3.7
MAINTAINER YangJunhai <yangjunhai@xindong.com>
COPY . /go/src/github.com/pirateXD/registrator
RUN go build -ldflags "-X main.Version=$(cat /go/src/github.com/pirateXD/registrator/VERSION)" -o /bin/registrator github.com/pirateXD/registrator
RUN rm -rf /go
ENTRYPOINT ["/bin/registrator"]
CMD ["--help"]
