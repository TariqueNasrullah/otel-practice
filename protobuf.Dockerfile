FROM golang:1.23.3

RUN apt-get update
RUN apt-get install unzip
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v28.3/protoc-28.3-linux-x86_64.zip
RUN unzip protoc-28.3-linux-x86_64.zip -d /usr/local/

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest


WORKDIR /proto/src

ENTRYPOINT ["make"]

CMD ["all"]
