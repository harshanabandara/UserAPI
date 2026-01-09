# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./docs ./docs

RUN CGO_ENABLED=0 go build -o /bin/server ./cmd/api-server

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/server /bin/server

USER 10001

EXPOSE 8080
ENTRYPOINT ["/bin/server"]


