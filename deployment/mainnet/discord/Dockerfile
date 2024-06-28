FROM golang:1.22.2-alpine as builder

RUN apk add --no-cache bash make git

RUN mkdir /pagu
WORKDIR /pagu
COPY ../../.. .

RUN go mod tidy
RUN go mod vendor
RUN make build-discord-mainnet

FROM alpine:3.14
RUN apk --no-cache add ca-certificates tzdata git

RUN mkdir /pagu-discord-mainnet
COPY --from=builder /pagu/build /pagu-discord-mainnet

RUN chmod +x /pagu-discord/pagu-discord-mainnet

CMD ["./pagu-discord/pagu-discord-mainnet", "-c", "./config/config.yml", "run"]