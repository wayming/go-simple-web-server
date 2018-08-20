FROM golang:1.10 as builder

WORKDIR /go/src/app
COPY cmds/simple_web_server/main.go .

RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o SimpleServer .


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /go/src/app
EXPOSE 8080/tcp

COPY --from=builder /go/src/app .
CMD ["./SimpleServer"]