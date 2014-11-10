FROM google/golang

RUN apt-get install -y -q netcat

WORKDIR /gopath/src/app
ADD . /gopath/src/app/
RUN go get app

CMD []
ENTRYPOINT ["/gopath/src/app/run.sh"]
