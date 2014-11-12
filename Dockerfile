FROM golang:1.3

COPY . /go/src/github.com/icecrime/octostats
WORKDIR /go/src/github.com/icecrime/octostats
RUN go get -d ./... && go install ./...
ENTRYPOINT ["octostats"]
