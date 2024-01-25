FROM golang:1.21 as builder

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

WORKDIR /app

RUN mkdir -p /app/docs

COPY --from=builder /go/src/app/main .

COPY --from=builder /go/src/app/docs /app/docs

EXPOSE 4000

CMD ["./main"]
