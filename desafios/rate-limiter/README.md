# Rate Limiter em Go

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Um rate limiter robusto e configurável em Go que pode limitar requisições HTTP com base em endereço IP ou token de acesso, com armazenamento em Redis e suporte para diferentes estratégias de persistência.

## 📋 Índice

- [Características](#características)
- [Arquitetura](#arquitetura)
- [Pré-requisitos](#pré-requisitos)
- [Instalação](#instalação)
- [Configuração](#configuração)
- [Como Usar](#como-usar)
- [Endpoints](#endpoints)
- [Testes](#testes)
- [Estrutura do Projeto](#estrutura-do-projeto)
- [Exemplos](#exemplos)

## ✨ Características

- ✅ **Limitação por IP**: Controla o número de requisições por endereço IP
- ✅ **Limitação por Token**: Suporta tokens de API com limites customizados
- ✅ **Priorização de Token**: Limites de token sobrepõem limites de IP
- ✅ **Bloqueio Temporário**: Bloqueia IPs/tokens que excedem o limite por um período configurável
- ✅ **Redis Integration**: Usa Redis para armazenamento distribuído
- ✅ **Strategy Pattern**: Fácil troca de backend de armazenamento (Redis, Memory, etc.)
- ✅ **Middleware HTTP**: Integração simples com qualquer aplicação Go
- ✅ **Configuração via Environment**: Todas as configurações via variáveis de ambiente ou arquivo `.env`
- ✅ **Docker Ready**: Totalmente dockerizado com docker-compose
- ✅ **Testes Completos**: Suite de testes unitários e de integração
- ✅ **Produção Ready**: Separação de lógica, middleware e storage

## 🏗️ Arquitetura

O projeto segue princípios de Clean Architecture e SOLID:

```
┌─────────────────────────────────────────────────┐
│              HTTP Middleware Layer              │
│  (Extração de IP/Token, Resposta HTTP)         │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│            Rate Limiter Logic Layer             │
│  (Regras de negócio, Validação de limites)     │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│            Storage Strategy Layer               │
│  (Interface abstrata para persistência)         │
└────────────┬───────────────────┬────────────────┘
             │                   │
    ┌────────▼────────┐  ┌──────▼─────────┐
    │  Redis Storage  │  │ Memory Storage │
    └─────────────────┘  └────────────────┘
```

### Componentes Principais

1. **Config**: Gerenciamento de configurações via environment variables
2. **Storage Interface**: Abstração para diferentes backends de armazenamento
3. **Rate Limiter**: Lógica de negócio para verificação de limites
4. **Middleware**: Integração com HTTP handlers
5. **Server**: Aplicação HTTP principal

## 🔧 Pré-requisitos

- **Go 1.21+**
- **Docker** e **Docker Compose** (para execução com containers)
- **Redis** (se executar localmente sem Docker)
- **Make** (opcional, para usar o Makefile)

## 📦 Instalação

### Usando Docker (Recomendado)

1. Clone o repositório:
```bash
git clone https://github.com/diogokimisima/goexpert.git
cd goexpert/desafios/rate-limiter
```

2. Inicie os serviços:
```bash
docker-compose up -d
```

A aplicação estará disponível em `http://localhost:8080`

### Instalação Local

1. Clone o repositório:
```bash
git clone https://github.com/diogokimisima/goexpert.git
cd goexpert/desafios/rate-limiter
```

2. Instale as dependências:
```bash
go mod download
```

3. Configure o Redis (certifique-se de que está rodando)

4. Execute a aplicação:
```bash
go run cmd/server/main.go
```

## ⚙️ Configuração

### Arquivo .env

Crie um arquivo `.env` na raiz do projeto (ou use variáveis de ambiente):

```env
# Configuração Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Rate Limiter - Limitação por IP
RATE_LIMIT_IP_REQUESTS=5           # Máximo de requisições por período
RATE_LIMIT_IP_DURATION=1s          # Período de tempo para contagem
RATE_LIMIT_IP_BLOCK_DURATION=5m    # Tempo de bloqueio após exceder

# Rate Limiter - Limitação por Token (Padrão)
RATE_LIMIT_TOKEN_REQUESTS=10
RATE_LIMIT_TOKEN_DURATION=1s
RATE_LIMIT_TOKEN_BLOCK_DURATION=5m

# Tokens Customizados (formato: token:requests:duration:block_duration)
RATE_LIMIT_TOKENS=abc123:100:1s:10m,xyz789:50:1s:3m

# Servidor
SERVER_PORT=8080
```

### Parâmetros de Configuração

| Parâmetro | Descrição | Padrão |
|-----------|-----------|--------|
| `REDIS_HOST` | Host do Redis | `localhost` |
| `REDIS_PORT` | Porta do Redis | `6379` |
| `REDIS_PASSWORD` | Senha do Redis | `` |
| `REDIS_DB` | Database do Redis | `0` |
| `RATE_LIMIT_IP_REQUESTS` | Requisições permitidas por IP | `5` |
| `RATE_LIMIT_IP_DURATION` | Janela de tempo para IP | `1s` |
| `RATE_LIMIT_IP_BLOCK_DURATION` | Tempo de bloqueio do IP | `5m` |
| `RATE_LIMIT_TOKEN_REQUESTS` | Requisições permitidas por token | `10` |
| `RATE_LIMIT_TOKEN_DURATION` | Janela de tempo para token | `1s` |
| `RATE_LIMIT_TOKEN_BLOCK_DURATION` | Tempo de bloqueio do token | `5m` |
| `RATE_LIMIT_TOKENS` | Configuração de tokens específicos | `` |
| `SERVER_PORT` | Porta do servidor | `8080` |

### Formato de Duração

Os valores de duração seguem o formato do Go:
- `s` - segundos (ex: `1s`, `30s`)
- `m` - minutos (ex: `5m`, `30m`)
- `h` - horas (ex: `1h`, `24h`)

## 🚀 Como Usar

### Usando Docker Compose

```bash
# Iniciar serviços
docker-compose up -d

# Ver logs
docker-compose logs -f app

# Parar serviços
docker-compose down
```

### Usando Makefile

```bash
# Ver todos os comandos disponíveis
make help

# Build da aplicação
make build

# Executar localmente
make run

# Executar testes
make test

# Executar testes com coverage
make test-coverage

# Iniciar com Docker
make docker-up

# Parar Docker
make docker-down

# Rebuild completo
make docker-rebuild
```

### Requisições HTTP

#### Sem Token (Limitado por IP)

```bash
curl http://localhost:8080/api/info
```

#### Com Token (Limitado por Token)

```bash
curl -H "API_KEY: abc123" http://localhost:8080/api/info
```

#### Resposta quando o limite é excedido

```json
{
  "message": "you have reached the maximum number of requests or actions allowed within a certain time frame"
}
```

Status Code: `429 Too Many Requests`

## 🌐 Endpoints

### GET /

Retorna informações sobre a API e endpoints disponíveis.

**Resposta:**
```json
{
  "message": "Welcome to Rate Limiter API",
  "endpoints": [
    "GET /health - Health check",
    "GET /api/info - Get API information",
    "POST /api/data - Submit data"
  ],
  "usage": {
    "rate_limit_by_ip": "Requests are limited by IP address",
    "rate_limit_by_token": "Use API_KEY header to authenticate and get higher limits"
  }
}
```

### GET /health

Health check endpoint.

**Resposta:**
```json
{
  "status": "healthy"
}
```

### GET /api/info

Retorna informações sobre a API e autenticação.

**Headers:**
- `API_KEY` (opcional): Token de autenticação

**Resposta:**
```json
{
  "message": "API Information",
  "version": "1.0.0",
  "authenticated": true,
  "token": "abc123"
}
```

### POST /api/data

Endpoint para submeter dados.

**Headers:**
- `API_KEY` (opcional): Token de autenticação
- `Content-Type: application/json`

**Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Resposta:**
```json
{
  "message": "Data received successfully",
  "data": {
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

## 🧪 Testes

### Executar Todos os Testes

```bash
go test -v ./...
```

ou

```bash
make test
```

### Testes com Coverage

```bash
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out -o coverage.html
```

ou

```bash
make test-coverage
```

### Script de Teste de Integração

Execute o script de teste que simula múltiplas requisições:

```bash
# Primeiro, inicie a aplicação
docker-compose up -d

# Execute o script de teste
bash test.sh
```

O script testará:
1. Limitação por IP (5 req/s)
2. Limitação por token (100 req/s para token `abc123`)
3. Duração do bloqueio

### Estrutura de Testes

- `internal/config/config_test.go` - Testes de configuração
- `internal/storage/memory_test.go` - Testes do storage em memória
- `internal/limiter/limiter_test.go` - Testes da lógica de rate limiting
- `internal/middleware/ratelimiter_test.go` - Testes do middleware HTTP

## 📁 Estrutura do Projeto

```
rate-limiter/
├── cmd/
│   └── server/
│       └── main.go                 # Aplicação principal
├── internal/
│   ├── config/
│   │   ├── config.go              # Gerenciamento de configuração
│   │   └── config_test.go
│   ├── storage/
│   │   ├── storage.go             # Interface de Storage
│   │   ├── redis.go               # Implementação Redis
│   │   ├── memory.go              # Implementação Memory (testes)
│   │   └── memory_test.go
│   ├── limiter/
│   │   ├── limiter.go             # Lógica do Rate Limiter
│   │   └── limiter_test.go
│   └── middleware/
│       ├── ratelimiter.go         # Middleware HTTP
│       └── ratelimiter_test.go
├── .env                            # Configurações (não commitado)
├── .env.example                    # Exemplo de configuração
├── .gitignore
├── docker-compose.yml              # Docker Compose
├── Dockerfile                      # Dockerfile multi-stage
├── go.mod                          # Dependências Go
├── go.sum
├── Makefile                        # Comandos úteis
├── test.sh                         # Script de teste
└── README.md                       # Esta documentação
```

## 💡 Exemplos

### Exemplo 1: Limitação por IP

Configuração: 5 requisições por segundo

```bash
# Primeira requisição - OK (200)
curl http://localhost:8080/api/info

# Segunda requisição - OK (200)
curl http://localhost:8080/api/info

# ... até a quinta requisição - OK (200)

# Sexta requisição - BLOQUEADA (429)
curl http://localhost:8080/api/info
# Resposta: {"message":"you have reached the maximum number of requests or actions allowed within a certain time frame"}
```

### Exemplo 2: Limitação por Token

Token `abc123` configurado para 100 req/s

```bash
# Requisições 1-100 - OK (200)
for i in {1..100}; do
  curl -H "API_KEY: abc123" http://localhost:8080/api/info
done

# Requisição 101 - BLOQUEADA (429)
curl -H "API_KEY: abc123" http://localhost:8080/api/info
```

### Exemplo 3: Token Sobrepõe IP

```bash
# Excede o limite do IP (5 req/s) sem token
for i in {1..6}; do
  curl http://localhost:8080/api/info
done
# Última requisição retorna 429

# Mas com token ainda funciona (limite separado)
curl -H "API_KEY: abc123" http://localhost:8080/api/info
# Retorna 200 OK
```

### Exemplo 4: Teste de Carga com Apache Bench

```bash
# Teste com 100 requisições, 10 concorrentes
ab -n 100 -c 10 http://localhost:8080/api/info

# Teste com token
ab -n 100 -c 10 -H "API_KEY: abc123" http://localhost:8080/api/info
```

### Exemplo 5: Integração em Código Go

```go
package main

import (
    "github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/config"
    "github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/limiter"
    "github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/middleware"
    "github.com/diogokimisima/goexpert/desafios/rate-limiter/internal/storage"
    "net/http"
)

func main() {
    // Carrega configuração
    cfg, _ := config.Load()
    
    // Inicializa storage
    store, _ := storage.NewRedisStorage(
        cfg.Redis.Host,
        cfg.Redis.Port,
        cfg.Redis.Password,
        cfg.Redis.DB,
    )
    defer store.Close()
    
    // Cria rate limiter
    rl := limiter.New(store, cfg)
    
    // Seu handler
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
    
    // Aplica middleware
    http.Handle("/", middleware.RateLimiterMiddleware(rl)(handler))
    
    http.ListenAndServe(":8080", nil)
}
```

## 🔄 Strategy Pattern para Storage

O projeto utiliza o Strategy Pattern para permitir fácil troca do backend de armazenamento:

```go
// Interface de Storage
type Storage interface {
    Increment(ctx context.Context, key string, expiration time.Duration) (int64, error)
    Get(ctx context.Context, key string) (int64, error)
    SetBlock(ctx context.Context, key string, duration time.Duration) error
    IsBlocked(ctx context.Context, key string) (bool, error)
    Close() error
}
```

### Implementações Disponíveis

1. **RedisStorage**: Produção, distribuído
2. **MemoryStorage**: Desenvolvimento, testes

### Criar Nova Implementação

```go
type MyCustomStorage struct {
    // seus campos
}

func (m *MyCustomStorage) Increment(ctx context.Context, key string, expiration time.Duration) (int64, error) {
    // sua implementação
}

// Implementar outros métodos da interface...
```

## 🐛 Troubleshooting

### Erro: Failed to connect to Redis

**Solução**: Certifique-se de que o Redis está rodando:
```bash
docker-compose up redis -d
```

### Erro: Port 8080 already in use

**Solução**: Altere a porta no `.env`:
```env
SERVER_PORT=8081
```

### Testes falhando

**Solução**: Execute os testes individualmente:
```bash
go test -v ./internal/limiter
go test -v ./internal/middleware
go test -v ./internal/storage
```

## 📝 Licença

Este projeto está sob a licença MIT.

## 👤 Autor

Desenvolvido como parte do desafio Go Expert.

## 🤝 Contribuindo

Contribuições são bem-vindas! Sinta-se à vontade para abrir issues ou pull requests.

## 📚 Referências

- [Go Documentation](https://golang.org/doc/)
- [Redis Documentation](https://redis.io/documentation)
- [Chi Router](https://github.com/go-chi/chi)
- [Rate Limiting Patterns](https://en.wikipedia.org/wiki/Rate_limiting)
