FROM golang:1.21.2-alpine3.17 as builder
WORKDIR /app
COPY . .
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/cloud-resources /app/cmd/main/main.go

FROM alpine:3.17
COPY --from=builder /app/cloud-resources /cloud-resources
ENTRYPOINT [ "/cloud-resources" ]