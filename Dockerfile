FROM golang:latest AS builder
LABEL stage=intermediate
WORKDIR /
COPY . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o /hillibot-discord .

FROM alpine:latest
LABEL maintainer = "hilli@github.com"
LABEL org.opencontainers.image.source https://github.com/hilli/hillibot-discord
RUN apk --no-cache add ca-certificates
WORKDIR /hillibot
COPY --from=builder /hillibot-discord ./
RUN chmod +x ./hillibot-discord
ENTRYPOINT [ "./hillibot-discord" ]