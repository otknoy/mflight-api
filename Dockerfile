FROM golang:1.16.0 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
COPY config/ config/
COPY domain/ domain/
COPY application/ application/
COPY interfaces/ interfaces/
COPY infrastructure/ infrastructure/
RUN CGO_ENABLED=0 go build -o mflight-api

FROM scratch
COPY --from=builder /app/mflight-api /bin/mflight-api
ENTRYPOINT ["/bin/mflight-api"]
