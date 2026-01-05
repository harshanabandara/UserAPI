# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /src

COPY go.work ./
COPY go.work.sum ./
COPY ./app/go.mod ./app/go.sum ./app/
COPY ./db/go.mod ./db/go.sum ./db/

RUN cd /src/app && go mod tidy

COPY ./app ./app
COPY ./db ./db

WORKDIR /src/app

RUN CGO_ENABLED=0 go build -o /bin/server .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /bin/server /bin/server

USER 10001

EXPOSE 8080
ENTRYPOINT ["/bin/server"]


