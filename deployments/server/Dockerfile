FROM golang:1.22.0 as builder

COPY go.mod .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server cmd/server/main.go

FROM scratch

COPY --from=builder /server /server
EXPOSE 7771
# чтения энвов нет, поэтому вот так
ENTRYPOINT ["/server", "-h", "0.0.0.0"]
