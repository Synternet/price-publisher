# Build container
FROM golang:1.21 AS build
ARG VERSION CI_JOB_TOKEN CI_SERVER_HOST
ARG NATS_URLS NATS_NKEY NATS_SUBJECT CMC_IDS CMC_API_KEY PUBLISH_INTERVAL_SEC

COPY . /home/src
WORKDIR /home/src

RUN CGO_ENABLED=0 go build -v -o . -installsuffix cgo -ldflags="-w -s" ./...

# Runtime container
FROM alpine

ENV TZ=UTC \
    PATH="/app:${PATH}"

RUN apk add --update --no-cache \
    tzdata \
    ca-certificates \
    bash \
    && \
    cp --remove-destination /usr/share/zoneinfo/${TZ} /etc/localtime && \
    echo "${TZ}" > /etc/timezone

WORKDIR /home/app
COPY --from=build /home/src/price-publisher .
COPY --from=build /home/src/docker/entrypoint.sh /home/app
RUN chmod +x /home/app/entrypoint.sh

ENTRYPOINT ["/home/app/entrypoint.sh"]
