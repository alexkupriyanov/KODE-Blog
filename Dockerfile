FROM golang:latest
ARG GIT_COMMIT=unspecified
LABEL git_commit=$GIT_COMMIT
ARG BUILD_TIME=unspecified
LABEL build_time=$BUILD_TIME
RUN mkdir -p /go/src/KODE-Blog
WORKDIR /go/src/KODE-Blog
COPY . /go/src/KODE-Blog
RUN go get
RUN go install
CMD ["go", "run", "main.go"]
EXPOSE 8080