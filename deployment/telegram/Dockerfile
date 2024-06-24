FROM golang:1.22.2-alpine as builder

RUN apk add --no-cache bash make git

RUN mkdir /pagu
WORKDIR /pagu
COPY ../.. .

RUN make build-telegram

FROM alpine:3.14
RUN apk --no-cache add ca-certificates tzdata git

RUN mkdir /pagu-telegram
COPY --from=builder /pagu/build /pagu-telegram

RUN chmod +x /pagu-telegram/pagu-telegram

CMD ["./pagu-telegram/pagu-telegram", "-c", "./config/config.yml", "run"]
