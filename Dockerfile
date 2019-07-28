FROM golang:latest

RUN mkdir -p /go/src/meter-panel

WORKDIR /go/src/meter-panel/

COPY . /go/src/meter-panel/

RUN go build /go/src/meter-panel/main.go

EXPOSE 12300

CMD /go/src/meter-panel/main
