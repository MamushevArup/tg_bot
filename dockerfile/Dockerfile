FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /krisha cmd/main.go

FROM scratch

WORKDIR /

COPY --from=builder /krisha /krisha

ENTRYPOINT ["/krisha"]