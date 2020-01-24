FROM golang:1.13 as builder

WORKDIR /go/src/github.com/iamsayantan/messagerooms

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# it will take flags from the environment
RUN go build -o messagerooms cmd/main.go

# App
FROM scratch
WORKDIR /messagerooms
COPY --from=builder /go/src/github.com/iamsayantan/messagerooms/messagerooms .
ENTRYPOINT ["./messagerooms"]