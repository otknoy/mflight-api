FROM golang:1.15.3 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
COPY mflight/ mflight/
RUN CGO_ENABLED=0 go build -o mflight-exporter

FROM scratch
COPY --from=builder /app/mflight-exporter /bin/mflight-exporter
EXPOSE 5000
ENTRYPOINT ["/bin/mflight-exporter"]
