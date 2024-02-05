#!/bin/sh

CMD="./price-publisher"

if [ ! -z "$NATS_URLS" ]; then
  CMD="$CMD --nats-urls $NATS_URLS"
fi

if [ ! -z "$NATS_NKEY" ]; then
  CMD="$CMD --nats-nkey $NATS_NKEY"
fi

if [ ! -z "$NATS_SUBJECT" ]; then
  CMD="$CMD --nats-subject $NATS_SUBJECT"
fi

if [ ! -z "$CMC_IDS" ]; then
  CMD="$CMD --cmd-ids $CMC_IDS"
fi

if [ ! -z "$CMC_API_KEY" ]; then
  CMD="$CMD --cmc-api-key $CMC_API_KEY"
fi

if [ ! -z "$PUBLISH_INTERVAL_SEC" ]; then
  CMD="$CMD --publish-interval-sec $PUBLISH_INTERVAL_SEC"
fi

exec $CMD
