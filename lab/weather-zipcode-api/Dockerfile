FROM golang:1.24 AS build

WORKDIR /app

# Copiar apenas os arquivos de módulo primeiro para aproveitar o cache do Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar o resto do código-fonte
COPY . .

# Compilar um binário estático otimizado
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-api ./cmd/server/main.go

# Usar scratch como imagem final (completamente vazia e mínima)
FROM scratch

# Copiar certificados CA para permitir chamadas HTTPS
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copiar o binário compilado
COPY --from=build /app/weather-api /app/weather-api

# Definir o executável como ponto de entrada
ENTRYPOINT ["/app/weather-api"]