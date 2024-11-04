ARG GO_IMAGE=golang:1.23.2
ARG DISTROLESS_IMAGE=gcr.io/distroless/base-debian12
FROM ${GO_IMAGE} AS builder

WORKDIR /app
COPY go.mod go.sum ./ ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/weather-radar ./cmd/main.go

FROM ${DISTROLESS_IMAGE} AS publisher

WORKDIR /
COPY --from=builder /app/weather-radar /weather-radar

EXPOSE 80
ENTRYPOINT ["/weather-radar"]
