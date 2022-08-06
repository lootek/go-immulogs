FROM golang:1.18-alpine AS builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o immulogsd ./cmd/immulogsd/

FROM scratch

COPY --from=builder /build/immulogsd .
ENTRYPOINT ["./immulogsd"]
