FROM golang:latest

COPY ./server .

RUN apt-get -y update && apt-get -y install libvips-dev \
    && rm -rf /var/lib/apt/lists/*

RUN unset GOPATH && export GO111MODULE=on && go mod download
RUN unset GOPATH && go build -o server *.go

# CMD ["sleep","infinity"]
CMD ["./server"]