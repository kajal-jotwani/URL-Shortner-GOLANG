FROM golang:alpine as builder

RUN mkdir /builder

Add ./build/

WORKDIR /build

RUN go build -o main .

#stage 2

FROM alpine

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY . /app

COPY --from=builder /build/main /app/main

WORKDIR /app

Expose 3000

CMD ["./main"]