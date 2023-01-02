FROM golang:1.18 as builder

WORKDIR /workspace

COPY ./ /workspace

RUN go mod download && go mod vendor

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o covidstats cmd/main.go

FROM registry.access.redhat.com/ubi8/ubi:8.4

COPY web web

ENV COVID_STATS_API_KEY=""

COPY --from=builder /workspace/covidstats /usr/local/bin/covidstats

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/covidstats"]