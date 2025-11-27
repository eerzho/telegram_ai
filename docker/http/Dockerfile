FROM golang:1.25-alpine3.22 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./main ./cmd/http/main.go

FROM gcr.io/distroless/static-debian12:nonroot AS prod

WORKDIR /app

COPY --from=build /app/main /app/main
COPY --from=build /app/prompts /app/prompts

EXPOSE 80

ENTRYPOINT ["/app/main"]
