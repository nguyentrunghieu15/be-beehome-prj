#!/bin/bash

# Wait until Elasticsearch is up
until curl -u "elastic:${ELASTIC_PASSWORD}" -X POST "http://localhost:9200/_security/user/kibana_system/_password" -d "{\"password\":\"${KIBANA_PASSWORD}\"}" -H "Content-Type: application/json" | grep -q "^{}"; do
  echo "Waiting for Elasticsearch..."
  sleep 10
done

echo "Kibana system user password set."