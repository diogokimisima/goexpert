# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copia os arquivos de dependências
COPY go.mod go.sum ./

# Download das dependências
RUN go mod download

# Copia o código fonte
COPY . .

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copia o binário do stage de build
COPY --from=builder /app/server .

# Expõe a porta 8080
EXPOSE 8080

# Comando para executar a aplicação
CMD ["./server"]
