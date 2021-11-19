FROM golang:1.16-alpine AS build

WORKDIR /app

COPY . .

RUN apk add --update gcc g++
RUN go mod download
RUN go build -o /eth-pools-metrics

FROM golang:1.16-alpine

WORKDIR /

COPY --from=build /eth-pools-metrics /eth-pools-metrics

ENTRYPOINT ["/eth-pools-metrics"]
