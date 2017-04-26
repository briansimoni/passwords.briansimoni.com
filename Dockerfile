# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

ADD . /go/src/github.com/briansimoni/simoni-password
WORKDIR /go/src/github.com/briansimoni/simoni-password
RUN go get
RUN go install github.com/briansimoni/simoni-password
ENV PORT 8888

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/simoni-password

# Document that the service listens on port 8888.
EXPOSE 8888