FROM golang:1.15.6 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
COPY config/ config/
COPY domain/ domain/
COPY handler/ handler/
COPY infrastructure/ infrastructure/
RUN CGO_ENABLED=0 go build -o mflight-exporter

FROM scratch
COPY --from=builder /app/mflight-exporter /bin/mflight-exporter
ENTRYPOINT ["/bin/mflight-exporter"]
