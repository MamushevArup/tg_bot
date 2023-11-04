FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY utils/texts/*.txt /utils/texts/

RUN CGO_ENABLED=0 go build -o /krisha cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /krisha .

COPY --from=builder /utils/texts/ /app/utils/text/txt

ENTRYPOINT ["/app/krisha"]