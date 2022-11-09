FROM curlimages/curl:7.80.0

ENV BASE_URL $BASE_URL
ENV REQUEST_HEADER_TOKEN $REQUEST_HEADER_TOKEN

WORKDIR /app

COPY ./scripts/merchant-score.sh .

ENTRYPOINT [ "sh", "merchant-score.sh"]
