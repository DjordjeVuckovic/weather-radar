ARG GO_IMAGE=golang:1.23.2
ARG DISTROLESS_IMAGE=gcr.io/distroless/base-debian12
FROM ${GO_IMAGE} AS builder

WORKDIR /app
COPY go.mod .
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/weather-radar main.go

FROM ${DISTROLESS_IMAGE} AS publisher

WORKDIR /
COPY --from=builder /app/weather-radar /weather-radar

EXPOSE 80
ENTRYPOINT ["/weather-radar"]
