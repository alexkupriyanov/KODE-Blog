FROM golang:latest as build
RUN mkdir -p /go/src/KODE-Blog
WORKDIR /go/src/KODE-Blog
COPY . /go/src/KODE-Blog
RUN go get
RUN go install
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM golang:latest as release
ARG GIT_COMMIT=unspecified
LABEL git_commit=$GIT_COMMIT
ARG BUILD_TIME=unspecified
LABEL build_time=$BUILD_TIME
WORKDIR /app
EXPOSE 8080
COPY --from=build /go/src/KODE-Blog/app .
RUN mkdir files
ENTRYPOINT ["./app"]
