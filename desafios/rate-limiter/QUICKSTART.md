# Quick Start Guide

## Iniciar o Projeto

```bash
# Com Docker (Recomendado)
docker-compose up -d

# Verificar logs
docker-compose logs -f app

# Parar
docker-compose down
```

## Testar

```bash
# Testes Unitários
go test -v ./...

# Script de Teste de Integração (precisa da app rodando)
docker-compose up -d
bash test.sh
```

## Exemplos de Uso

### 1. Teste Simples
```bash
curl http://localhost:8080/
```

### 2. Teste de Rate Limit por IP (5 req/s)
```bash
# Enviar 7 requisições rapidamente
for i in {1..7}; do
  curl -i http://localhost:8080/api/info
done
```

### 3. Teste com Token (100 req/s para token abc123)
```bash
curl -H "API_KEY: abc123" http://localhost:8080/api/info
```

### 4. Teste POST
```bash
curl -X POST http://localhost:8080/api/data \
  -H "Content-Type: application/json" \
  -d '{"name":"John","email":"john@example.com"}'
```

## Configuração

Edite o arquivo `.env` para ajustar os limites:

```env
# Limites por IP
RATE_LIMIT_IP_REQUESTS=5
RATE_LIMIT_IP_DURATION=1s
RATE_LIMIT_IP_BLOCK_DURATION=5m

# Limites por Token
RATE_LIMIT_TOKEN_REQUESTS=10
RATE_LIMIT_TOKEN_DURATION=1s
RATE_LIMIT_TOKEN_BLOCK_DURATION=5m

# Tokens Customizados
RATE_LIMIT_TOKENS=abc123:100:1s:10m,xyz789:50:1s:3m
```

## Endpoints

- `GET /` - Home com informações da API
- `GET /health` - Health check
- `GET /api/info` - Informações da API
- `POST /api/data` - Submeter dados

## Resposta quando Bloqueado

Status: `429 Too Many Requests`
```json
{
  "message": "you have reached the maximum number of requests or actions allowed within a certain time frame"
}
```

## Comandos Make

```bash
make help          # Mostra todos os comandos
make build         # Build da aplicação
make run           # Executar localmente
make test          # Executar testes
make docker-up     # Iniciar com Docker
make docker-down   # Parar Docker
```

## Estrutura

```
rate-limiter/
├── cmd/server/main.go           # Aplicação principal
├── internal/
│   ├── config/                  # Configurações
│   ├── storage/                 # Strategy pattern (Redis, Memory)
│   ├── limiter/                 # Lógica do rate limiter
│   └── middleware/              # Middleware HTTP
├── docker-compose.yml           # Docker Compose
├── Dockerfile                   # Multi-stage build
└── README.md                    # Documentação completa
```

## Troubleshooting

### Redis não conecta
```bash
docker-compose up redis -d
docker-compose logs redis
```

### Porta 8080 em uso
```bash
# Edite .env e altere SERVER_PORT=8081
docker-compose down
docker-compose up -d
```

### Rebuild completo
```bash
docker-compose down -v
docker-compose build --no-cache
docker-compose up -d
```
