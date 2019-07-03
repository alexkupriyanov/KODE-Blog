FROM golang:latest

RUN mkdir -p /go/src/app

WORKDIR /go/src/app

COPY . /go/src/app

RUN go get

RUN go install

CMD ["go", "run", "main.go"]

EXPOSE 8080