FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copiar arquivos de dependências
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo o código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/auction

# Imagem final
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar o binário da etapa de build
COPY --from=builder /app/main .
COPY --from=builder /app/cmd/auction/.env ./cmd/auction/.env

EXPOSE 8080

CMD ["./main"]