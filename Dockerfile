# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/briansimoni/simoni-password

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
WORKDIR /go/src/github.com/briansimoni/simoni-password
RUN go get

RUN go install github.com/briansimoni/simoni-password

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/simoni-password

# Document that the service listens on port 8888.
EXPOSE 8888