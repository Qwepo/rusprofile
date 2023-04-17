FROM golang:latest


ENV GOPATH /go
WORKDIR /rusprof

COPY go.mod go.sum ./
RUN go mod download


COPY cmd cmd
COPY pkg pkg
COPY internal internal
COPY gen gen
COPY config.yaml ./

RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/app/main.go

CMD [ "./main" ]