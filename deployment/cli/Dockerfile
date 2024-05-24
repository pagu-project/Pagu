FROM golang:1.22.2-alpine as builder

RUN apk add --no-cache bash make git

RUN mkdir /pagu
WORKDIR /pagu
COPY ../.. .

RUN make build-cli

FROM alpine:3.14
RUN apk --no-cache add ca-certificates tzdata git

RUN mkdir /pagu-cli
COPY --from=builder /pagu/build /pagu-cli
COPY --from=builder /pagu/config/config.yml /pagu-cli/config/config.yml

RUN chmod +x /pagu-cli/pagu-cli

CMD ["./pagu-cli/pagu-cli", "-c", "./pagu-cli/config/config.yml", "run"]