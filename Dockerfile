FROM golang:1.19.3 as builder
WORKDIR /usr/local/bin
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/app ./cmd

FROM alpine:latest
EXPOSE 8888
WORKDIR /usr/local/bin
RUN apk --no-cache add ca-certificates
COPY --from=builder /usr/local/bin/bin .
COPY --from=builder /usr/local/bin/config/config.yml ./config/
COPY --from=builder /usr/local/bin/migrations/ ./migrations/
CMD [ "app" ]