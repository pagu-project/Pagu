FROM golang:1.21-alpine as builder

RUN apk add --no-cache bash make git curl

RUN mkdir /app
WORKDIR /app

COPY . /app

# install tools
RUN make install-tools

# build
RUN make build-bot

FROM alpine:3.14
RUN apk --no-cache add ca-certificates tzdata git
RUN mkdir /app
RUN mkdir bot-data
COPY --from=builder /app/main /app
RUN chmod +x /app/main
CMD ["./app/main"]
