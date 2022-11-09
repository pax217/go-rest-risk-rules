#!/bin/sh

echo $BASE_URL
echo $REQUEST_HEADER_TOKEN

curl  --location --request POST ''$BASE_URL'/risk-rules/v1/merchants_score' --header 'Content-Type: application/json' --header 'X-Request-Risk: '$REQUEST_HEADER_TOKEN'' --data-raw '{}'