FROM golang:1.22.0 as builder

COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /client cmd/client/main.go

FROM scratch

COPY --from=builder /client /client
# чтения энвов нет, поэтому вот так
ENTRYPOINT ["/client", "-h", "server"]
