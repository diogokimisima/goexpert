# Sistema de Leilão com Concorrência

Sistema de leilão em tempo real desenvolvido em Go com processamento concorrente de lances (bids) utilizando goroutines e canais.

## 🏗️ Arquitetura

- **Backend**: Go 1.24.3 com Gin Framework
- **Banco de Dados**: MongoDB
- **Containerização**: Docker & Docker Compose

## 📋 Pré-requisitos

- Docker
- Docker Compose
- curl (para testes)

## 🚀 Como executar

### 1. Subir os serviços

```bash
docker-compose up --build -d
```

Isso irá iniciar:
- MongoDB na porta `27017`
- API REST na porta `8084`

### 2. Verificar se os serviços estão rodando

```bash
docker ps
```

Você deve ver dois containers:
- `auctionsDB` (MongoDB)
- `auction-api` (API Go)

## ⚙️ Configurações do Sistema

As variáveis de ambiente configuradas no `docker-compose.yml` controlam o comportamento do sistema:

### Variáveis de Ambiente

| Variável | Valor | Descrição |
|----------|-------|-----------|
| `MONGODB_URL` | `mongodb://mongodb:27017` | URL de conexão com o MongoDB |
| `MONGODB_DB` | `auctions` | Nome do banco de dados |
| `BATCH_INSERT_INTERVAL` | `10s` | Intervalo de tempo para processar lote de bids |
| `MAX_BATCH_SIZE` | `3` | Quantidade máxima de bids por lote |
| `AUCTION_INTERVAL` | `10m` | Tempo de duração de um leilão antes de ser marcado como completo |

### 📊 Como funciona o sistema de Bids em Batch

O sistema utiliza **processamento em lote (batch)** para otimizar a inserção de lances no banco de dados:

1. Quando um bid é criado via API, ele é enviado para um **canal (channel)** do Go
2. Os bids ficam acumulados em memória em um batch
3. O batch é processado e inserido no banco quando:
   - Atinge o `MAX_BATCH_SIZE` (3 bids), **OU**
   - Passa o tempo do `BATCH_INSERT_INTERVAL` (10 segundos)

**Exemplo prático:**
- Se você criar 3 bids rapidamente → eles são inseridos imediatamente
- Se você criar 1 ou 2 bids → eles serão inseridos após 10 segundos

### ⏱️ Ciclo de Vida do Leilão

Quando um leilão é criado:
1. Status inicial: `Active` (0)
2. Após `AUCTION_INTERVAL` (10 minutos): Status muda automaticamente para `Completed` (1)
3. **Importante**: Apenas leilões com status `Active` aceitam novos lances

## 🧪 Testando o Sistema

### Passo 1: Criar Usuários

Os usuários precisam ser criados manualmente no MongoDB:

```bash
docker exec -it auctionsDB mongosh

# Dentro do mongosh:
use auctions

db.users.insertMany([
  {
    "_id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "João Silva"
  },
  {
    "_id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "Maria Santos"
  }
])

exit
```

### Passo 2: Criar um Leilão

```bash
curl -X POST http://localhost:8084/auctions \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "iPhone 15 Pro",
    "category": "Smartphones",
    "description": "iPhone 15 Pro Max 256GB novo na caixa lacrada",
    "condition": 0
  }'
```

**Valores válidos para `condition`:**
- `0` = Novo (New)
- `1` = Usado (Used)
- `2` = Recondicionado (Refurbished)

### Passo 3: Listar Leilões

```bash
# Listar todos os leilões ativos
curl http://localhost:8084/auctions

# Listar com filtros
curl "http://localhost:8084/auctions?status=0&category=Smartphones"
```

**Valores de status:**
- `0` = Ativo (Active)
- `1` = Completo (Completed)

**Exemplo de resposta:**
```json
[
  {
    "id": "9d7b877f-8bf2-4aae-96bf-db56beb8e2c6",
    "product_name": "iPhone 15 Pro",
    "category": "Smartphones",
    "description": "iPhone 15 Pro Max 256GB novo na caixa lacrada",
    "condition": 0,
    "status": 0,
    "time_stamp": "2026-01-08T13:34:26Z"
  }
]
```

### Passo 4: Buscar Leilão por ID

```bash
# Substitua {auctionId} pelo ID real
curl http://localhost:8084/auctions/9d7b877f-8bf2-4aae-96bf-db56beb8e2c6
```

### Passo 5: Criar Lances (Bids)

**Importante**: Crie pelo menos 3 bids para atingir o `MAX_BATCH_SIZE` e ver o resultado imediatamente.

```bash
# Lance 1 - R$ 1.500
curl -X POST http://localhost:8084/bid \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "auction_id": "9d7b877f-8bf2-4aae-96bf-db56beb8e2c6",
    "amount": 1500.00
  }'

# Lance 2 - R$ 1.600
curl -X POST http://localhost:8084/bid \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440002",
    "auction_id": "9d7b877f-8bf2-4aae-96bf-db56beb8e2c6",
    "amount": 1600.00
  }'

# Lance 3 - R$ 1.700
curl -X POST http://localhost:8084/bid \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "auction_id": "9d7b877f-8bf2-4aae-96bf-db56beb8e2c6",
    "amount": 1700.00
  }'
```

### Passo 6: Buscar Lances de um Leilão

```bash
# Aguardar 2 segundos (se não criou 3 bids)
sleep 2

# Buscar todos os lances
curl http://localhost:8084/bid/9d7b877f-8bf2-4aae-96bf-db56beb8e2c6
```

**Exemplo de resposta:**
```json
[
  {
    "id": "a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "auction_id": "9d7b877f-8bf2-4aae-96bf-db56beb8e2c6",
    "amount": 1500,
    "timestamp": "2026-01-08 14:30:45"
  },
  {
    "id": "b2c3d4e5-f6a7-4b5c-9d0e-2f3a4b5c6d7e",
    "user_id": "550e8400-e29b-41d4-a716-446655440002",
    "auction_id": "9d7b877f-8bf2-4aae-96bf-db56beb8e2c6",
    "amount": 1600,
    "timestamp": "2026-01-08 14:30:46"
  },
  {
    "id": "c3d4e5f6-a7b8-4c5d-0e1f-3a4b5c6d7e8f",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "auction_id": "9d7b877f-8bf2-4aae-96bf-db56beb8e2c6",
    "amount": 1700,
    "timestamp": "2026-01-08 14:30:47"
  }
]
```

### Passo 7: Buscar o Vencedor do Leilão

```bash
curl http://localhost:8084/auction/winner/9d7b877f-8bf2-4aae-96bf-db56beb8e2c6
```

**Exemplo de resposta:**
```json
{
  "auction": {
    "id": "9d7b877f-8bf2-4aae-96bf-db56beb8e2c6",
    "product_name": "iPhone 15 Pro",
    "category": "Smartphones",
    "description": "iPhone 15 Pro Max 256GB novo na caixa lacrada",
    "condition": 0,
    "status": 0,
    "time_stamp": "2026-01-08T13:34:26Z"
  },
  "bid": {
    "id": "c3d4e5f6-a7b8-4c5d-0e1f-3a4b5c6d7e8f",
    "user_id": "550e8400-e29b-41d4-a716-446655440001",
    "auction_id": "9d7b877f-8bf2-4aae-96bf-db56beb8e2c6",
    "amount": 1700,
    "timestamp": "2026-01-08 14:30:47"
  }
}
```

### Passo 8: Buscar Usuário por ID

```bash
curl http://localhost:8084/users/550e8400-e29b-41d4-a716-446655440001
```

## 📝 Endpoints da API

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| `POST` | `/auctions` | Criar novo leilão |
| `GET` | `/auctions` | Listar leilões (com filtros opcionais) |
| `GET` | `/auctions/:auctionId` | Buscar leilão por ID |
| `POST` | `/bid` | Criar novo lance |
| `GET` | `/bid/:auctionId` | Buscar lances de um leilão |
| `GET` | `/auction/winner/:auctionId` | Buscar lance vencedor |
| `GET` | `/users/:userId` | Buscar usuário por ID |

## 🔍 Troubleshooting

### Bids retornam null

**Problema**: Os bids não aparecem ao buscar.

**Possíveis causas:**

1. **Leilão expirado**: Verifique se o leilão está com `status: 0` (Active)
   ```bash
   curl http://localhost:8084/auctions
   ```
   Se o status for `1` (Completed), crie um novo leilão.

2. **Batch não processado**: Aguarde 10 segundos ou crie mais bids para atingir o `MAX_BATCH_SIZE` de 3.

3. **Usuário não existe**: Verifique se o usuário foi criado no MongoDB.

### Ver logs da aplicação

```bash
# Ver logs em tempo real
docker logs -f auction-api

# Ver últimas 50 linhas
docker logs --tail 50 auction-api
```

### Verificar dados no MongoDB

```bash
docker exec -it auctionsDB mongosh

use auctions
db.auctions.find().pretty()
db.bids.find().pretty()
db.users.find().pretty()
exit
```

### Limpar o banco de dados

```bash
docker exec -it auctionsDB mongosh

use auctions
db.auctions.deleteMany({})
db.bids.deleteMany({})
db.users.deleteMany({})
exit
```

## 🛑 Parar os serviços

```bash
# Parar containers
docker-compose down

# Parar e remover volumes (apaga todos os dados)
docker-compose down -v
```

## 🔧 Ajustes para Produção

Para ambiente de produção, considere ajustar as seguintes variáveis no `docker-compose.yml`:

```yaml
environment:
  BATCH_INSERT_INTERVAL: 5m    # 5 minutos
  MAX_BATCH_SIZE: 50           # 50 bids por lote
  AUCTION_INTERVAL: 24h        # 24 horas de duração
```

## 📚 Conceitos de Concorrência Utilizados

- **Goroutines**: Processamento assíncrono de bids e atualização de status do leilão
- **Channels**: Comunicação entre goroutines para processamento em batch
- **WaitGroups**: Sincronização de múltiplas goroutines na inserção de bids
- **Timers**: Controle de intervalo para processamento de batches

## 🏆 Características do Sistema

- ✅ Processamento concorrente de múltiplos lances
- ✅ Otimização com batch processing
- ✅ Validação de UUIDs
- ✅ Leilões com tempo de expiração automático
- ✅ Determinação automática do vencedor (maior lance)
- ✅ API RESTful com validações
