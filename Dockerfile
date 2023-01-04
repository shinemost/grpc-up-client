FROM golang:alpine AS build

LABEL maintainer="hjfu"

ENV GO11MODULE=on \
  CGO_ENABLE=on \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app

COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go get -d ./...
RUN go install ./...

RUN go build -mod=mod -o grpc-client .


FROM alpine
WORKDIR /app

COPY --from=build /app/grpc-client ./grpc-client
COPY ./certs ./certs

ENTRYPOINT ["/app/grpc-client"]