FROM golang:1.16-alpine AS build

WORKDIR /app

COPY . .

RUN apk add --update gcc g++
RUN go mod download

# DO NOT MERGE THIS: DIRTY TRICK TO AVOID
# prysmaticlabs/prysm#10153
RUN sed -i.bak -e '4033d' /go/pkg/mod/github.com/prysmaticlabs/prysm@v0.0.0-20220128215931-aba628b56bc2/proto/prysm/v1alpha1/generated.ssz.go
RUN go build -o /eth-pools-metrics

FROM golang:1.16-alpine

WORKDIR /

COPY --from=build /eth-pools-metrics /eth-pools-metrics

ENTRYPOINT ["/eth-pools-metrics"]
