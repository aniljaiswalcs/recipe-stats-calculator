FROM golang:1.20-alpine as builder
WORKDIR /aggregator
ENV GO111MODULE=on

COPY . .
RUN go mod download
RUN go build -o recipe-stats-calculator main.go


FROM alpine:latest
WORKDIR /recipe/
COPY --from=builder /aggregator/recipe-stats-calculator .
COPY --from=builder /aggregator/hf_test_calculation_fixtures.json .
COPY --from=builder /aggregator/config.yml .
CMD ["./recipe-stats-calculator"]