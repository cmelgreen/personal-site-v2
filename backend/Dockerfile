FROM 010629071893.dkr.ecr.us-east-1.amazonaws.com/golang:latest

RUN apt-get -y update && apt-get -y install libvips-dev \
    && rm -rf /var/lib/apt/lists/*

COPY ./server .

RUN unset GOPATH && export GO111MODULE=on && go mod download
RUN unset GOPATH && go build -o server *.go

CMD ["./server"]