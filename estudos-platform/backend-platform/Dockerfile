# API Go + CLI de migrations numa imagem pequena.
# No start, aplica as migrations (embutidas no binário) e sobe a API.
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api \
 && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/migrate ./cmd/migrate

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app
WORKDIR /app
COPY --from=build /out/api /out/migrate ./
USER app
ENV APP_ENV=production
EXPOSE 8080
CMD ["sh", "-c", "./migrate && exec ./api"]
