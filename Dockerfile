FROM golang:alpine3.19 as builder
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /workspace
COPY go.mod .
COPY go.sum .
COPY . .
RUN go mod download && go build ./cmd/magneticod && go build ./cmd/magneticow

FROM alpine:3.19
WORKDIR /tmp
COPY --from=builder /workspace/magneticod /usr/bin/
COPY --from=builder /workspace/magneticow /usr/bin/
ENTRYPOINT ["/usr/bin/magneticod", "--help"]
LABEL org.opencontainers.image.source=https://github.com/tgragnato/magnetico
