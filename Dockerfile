FROM golang:1.18.0 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY app/ app/
RUN CGO_ENABLED=0 go build -o mflight-api app/main.go

FROM scratch
COPY --from=builder /app/mflight-api /bin/
ENTRYPOINT ["/bin/mflight-api"]
