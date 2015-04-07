FROM golang:1.3

COPY . /go/src/github.com/icecrime/octostats
WORKDIR /go/src/github.com/icecrime/octostats
RUN GOPATH=$GOPATH:/go/src/github.com/icecrime/octostats/Godeps/_workspace go install ./...
ENTRYPOINT ["octostats"]
