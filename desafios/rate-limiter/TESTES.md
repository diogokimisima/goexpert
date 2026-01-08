# Guia Completo de Testes - Rate Limiter

Este documento contém todos os testes necessários para validar que o rate limiter atende todos os requisitos do projeto.

## 📋 Checklist de Requisitos

- [ ] Limitação por IP
- [ ] Limitação por Token (header API_KEY)
- [ ] Token sobrepõe configuração de IP
- [ ] Middleware injetável
- [ ] Configuração por variáveis de ambiente/.env
- [ ] Bloqueio temporário configurável
- [ ] Resposta HTTP 429 com mensagem correta
- [ ] Redis como storage
- [ ] Strategy pattern implementado
- [ ] Lógica separada do middleware
- [ ] Docker/Docker Compose funcional
- [ ] Servidor na porta 8080

---

## 🚀 Passo 1: Preparar o Ambiente

### 1.1 Iniciar os Serviços

```bash
# Parar tudo que estiver rodando
docker-compose down -v

# Build e iniciar os containers
docker-compose up --build -d

# Verificar se os serviços estão rodando
docker-compose ps
```

**Resultado esperado:**
```
NAME                    STATUS              PORTS
rate-limiter-app        Up                  0.0.0.0:8080->8080/tcp
rate-limiter-redis      Up (healthy)        0.0.0.0:6379->6379/tcp
```

### 1.2 Verificar Logs

```bash
# Ver logs da aplicação
docker-compose logs app

# Verificar se conectou no Redis
docker-compose logs app | grep -i redis
```

**Deve mostrar:** `Connected to Redis successfully`

---

## 🧪 Passo 2: Testes Unitários

### 2.1 Executar Todos os Testes

```bash
go test -v ./...
```

**Resultado esperado:** Todos os testes devem passar (PASS)

### 2.2 Testes com Coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func coverage.out
```

**Resultado esperado:** Coverage > 70%

---

## 🌐 Passo 3: Testes de Integração

### 3.1 Teste Básico - Servidor Respondendo

```bash
curl http://localhost:8080/
```

**Resultado esperado:**
- Status: `200 OK`
- JSON com informações da API

### 3.2 Health Check

```bash
curl http://localhost:8080/health
```

**Resultado esperado:**
```json
{"status":"healthy"}
```

---

## 🔒 Passo 4: Testar Limitação por IP

### 4.1 Configuração Atual

Verifique no `.env`:
```
RATE_LIMIT_IP_REQUESTS=5
RATE_LIMIT_IP_DURATION=1s
RATE_LIMIT_IP_BLOCK_DURATION=5m
```

### 4.2 Teste Manual - Enviar 7 Requisições

```bash
echo "=== Teste de Rate Limit por IP ==="
for i in {1..7}; do
  echo -n "Request $i: "
  curl -s -o /dev/null -w "HTTP %{http_code}\n" http://localhost:8080/api/info
  sleep 0.1
done
```

**Resultado esperado:**
```
Request 1: HTTP 200
Request 2: HTTP 200
Request 3: HTTP 200
Request 4: HTTP 200
Request 5: HTTP 200
Request 6: HTTP 429  ✅ BLOQUEADO
Request 7: HTTP 429  ✅ CONTINUA BLOQUEADO
```

### 4.3 Verificar Mensagem de Erro

```bash
curl -i http://localhost:8080/api/info
```

**Deve retornar:**
```
HTTP/1.1 429 Too Many Requests
Content-Type: application/json

{"message":"you have reached the maximum number of requests or actions allowed within a certain time frame"}
```

### 4.4 Verificar Bloqueio Persistente

```bash
# Aguardar 2 segundos e tentar novamente
sleep 2
curl -s -o /dev/null -w "HTTP %{http_code}\n" http://localhost:8080/api/info
```

**Resultado esperado:** `HTTP 429` (ainda bloqueado)

### 4.5 Verificar no Redis

```bash
docker-compose exec redis redis-cli KEYS "*"
```

**Deve mostrar:** Chaves como `ip:192.168.65.1` e `ip:192.168.65.1:blocked`

---

## 🎫 Passo 5: Testar Limitação por Token

### 5.1 Configuração de Tokens

Verifique no `.env`:
```
RATE_LIMIT_TOKEN_REQUESTS=10          # Token padrão
RATE_LIMIT_TOKENS=abc123:100:1s:10m   # Token abc123 com 100 req/s
```

### 5.2 Teste com Token Padrão (10 requisições)

```bash
echo "=== Teste com Token (limite padrão 10 req/s) ==="
for i in {1..12}; do
  echo -n "Request $i: "
  curl -s -o /dev/null -w "HTTP %{http_code}\n" \
    -H "API_KEY: token-padrao" \
    http://localhost:8080/api/info
  sleep 0.1
done
```

**Resultado esperado:**
```
Request 1-10: HTTP 200
Request 11: HTTP 429  ✅ BLOQUEADO
Request 12: HTTP 429
```

### 5.3 Teste com Token VIP (100 requisições)

```bash
echo "=== Teste com Token VIP abc123 (100 req/s) ==="
for i in {1..101}; do
  if [ $i -le 5 ] || [ $i -eq 100 ] || [ $i -eq 101 ]; then
    echo -n "Request $i: "
    curl -s -o /dev/null -w "HTTP %{http_code}\n" \
      -H "API_KEY: abc123" \
      http://localhost:8080/api/info
  fi
  sleep 0.01
done
```

**Resultado esperado:**
```
Request 1-5: HTTP 200
...
Request 100: HTTP 200
Request 101: HTTP 429  ✅ BLOQUEADO após 100
```

---

## 🎯 Passo 6: Testar Prioridade de Token sobre IP

### 6.1 Bloquear o IP Primeiro

```bash
echo "=== Bloqueando IP primeiro ==="
# Fazer 6 requisições sem token para bloquear o IP
for i in {1..6}; do
  curl -s -o /dev/null http://localhost:8080/api/info
done

# Verificar que IP está bloqueado
echo -n "IP bloqueado? "
curl -s -o /dev/null -w "HTTP %{http_code}\n" http://localhost:8080/api/info
```

**Resultado esperado:** `HTTP 429`

### 6.2 Testar com Token no Mesmo IP

```bash
echo "=== Testando com token no IP bloqueado ==="
curl -s -o /dev/null -w "HTTP %{http_code}\n" \
  -H "API_KEY: abc123" \
  http://localhost:8080/api/info
```

**Resultado esperado:** `HTTP 200` ✅ **Token funciona mesmo com IP bloqueado!**

---

## ⏱️ Passo 7: Testar Tempo de Bloqueio

### 7.1 Configurar Bloqueio Curto para Teste

Edite `.env` temporariamente:
```env
RATE_LIMIT_IP_BLOCK_DURATION=10s  # 10 segundos para teste rápido
```

```bash
# Reiniciar aplicação
docker-compose restart app
sleep 3
```

### 7.2 Executar Teste de Bloqueio

```bash
echo "=== Teste de Duração do Bloqueio ==="

# Bloquear o IP
for i in {1..6}; do
  curl -s -o /dev/null http://localhost:8080/api/info
done

# Verificar bloqueio
echo -n "Bloqueado: "
curl -s -o /dev/null -w "HTTP %{http_code}\n" http://localhost:8080/api/info

# Aguardar 11 segundos
echo "Aguardando 11 segundos..."
sleep 11

# Tentar novamente
echo -n "Após 11 segundos: "
curl -s -o /dev/null -w "HTTP %{http_code}\n" http://localhost:8080/api/info
```

**Resultado esperado:**
```
Bloqueado: HTTP 429
Aguardando 11 segundos...
Após 11 segundos: HTTP 200  ✅ DESBLOQUEADO
```

### 7.3 Restaurar Configuração

```bash
# Restaurar .env para 5 minutos
sed -i '' 's/RATE_LIMIT_IP_BLOCK_DURATION=10s/RATE_LIMIT_IP_BLOCK_DURATION=5m/' .env
docker-compose restart app
```

---

## 🐳 Passo 8: Testar Redis e Strategy Pattern

### 8.1 Verificar Dados no Redis

```bash
# Listar todas as chaves
docker-compose exec redis redis-cli KEYS "*"

# Ver valor de uma chave de contador
docker-compose exec redis redis-cli GET "ip:192.168.65.1"

# Ver TTL (tempo de expiração)
docker-compose exec redis redis-cli TTL "ip:192.168.65.1"
```

### 8.2 Limpar Redis

```bash
docker-compose exec redis redis-cli FLUSHALL
echo "Redis limpo! Pode testar novamente."
```

### 8.3 Verificar Strategy Pattern

```bash
# Verificar que existe interface Storage
grep -n "type Storage interface" internal/storage/storage.go

# Verificar implementações
ls internal/storage/*.go
```

**Deve mostrar:** `storage.go`, `redis.go`, `memory.go`

---

## 📊 Passo 9: Testes de Carga

### 9.1 Teste com Apache Bench (se instalado)

```bash
# Instalar ab (macOS)
# brew install httpd

# Teste de carga
ab -n 100 -c 10 http://localhost:8080/api/info
```

### 9.2 Teste com Loop Rápido

```bash
echo "=== Teste de Carga - 50 requisições ==="
for i in {1..50}; do
  curl -s -o /dev/null -w "%{http_code} " http://localhost:8080/api/info &
done
wait
echo -e "\n\nVerifique quantas retornaram 200 vs 429"
```

### 9.3 Teste Concorrente com Múltiplos Tokens

```bash
#!/bin/bash
echo "=== Teste Concorrente ==="

# Token 1 em background
for i in {1..10}; do
  curl -s -o /dev/null -H "API_KEY: abc123" http://localhost:8080/api/info &
done

# Token 2 em background
for i in {1..10}; do
  curl -s -o /dev/null -H "API_KEY: xyz789" http://localhost:8080/api/info &
done

wait
echo "Teste concorrente finalizado"
```

---

## 🎬 Passo 10: Script de Teste Automatizado

### 10.1 Executar Script Completo

```bash
chmod +x test.sh
./test.sh
```

### 10.2 Criar Script de Validação Completo

Crie o arquivo `validate-all.sh`:

```bash
#!/bin/bash

echo "============================================"
echo "  VALIDAÇÃO COMPLETA DO RATE LIMITER"
echo "============================================"
echo ""

FAILED=0
PASSED=0

# Função para testar
test_endpoint() {
  local name="$1"
  local expected_code="$2"
  shift 2
  local curl_args="$@"
  
  echo -n "[$name] "
  CODE=$(curl -s -o /dev/null -w "%{http_code}" $curl_args http://localhost:8080/api/info)
  
  if [ "$CODE" == "$expected_code" ]; then
    echo "✅ PASS (HTTP $CODE)"
    ((PASSED++))
  else
    echo "❌ FAIL (Expected $expected_code, got $CODE)"
    ((FAILED++))
  fi
}

# Limpar Redis
echo "Limpando Redis..."
docker-compose exec -T redis redis-cli FLUSHALL > /dev/null
sleep 1

echo ""
echo "=== Testes de Limitação por IP ==="
test_endpoint "IP Request 1" "200"
test_endpoint "IP Request 2" "200"
test_endpoint "IP Request 3" "200"
test_endpoint "IP Request 4" "200"
test_endpoint "IP Request 5" "200"
test_endpoint "IP Request 6 (bloqueado)" "429"
test_endpoint "IP Request 7 (ainda bloqueado)" "429"

# Limpar Redis
docker-compose exec -T redis redis-cli FLUSHALL > /dev/null
sleep 1

echo ""
echo "=== Testes de Limitação por Token ==="
for i in {1..10}; do
  test_endpoint "Token Request $i" "200" "-H 'API_KEY: test-token'"
done
test_endpoint "Token Request 11 (bloqueado)" "429" "-H 'API_KEY: test-token'"

# Limpar Redis
docker-compose exec -T redis redis-cli FLUSHALL > /dev/null
sleep 1

echo ""
echo "=== Teste de Prioridade Token sobre IP ==="
# Bloquear IP
for i in {1..6}; do
  curl -s -o /dev/null http://localhost:8080/api/info
done
test_endpoint "IP bloqueado" "429"
test_endpoint "Token funciona com IP bloqueado" "200" "-H 'API_KEY: abc123'"

echo ""
echo "=== Teste de Token VIP (100 req/s) ==="
# Limpar Redis
docker-compose exec -T redis redis-cli FLUSHALL > /dev/null
sleep 1

# Fazer 100 requisições com token VIP (mais rápido, só testa algumas)
SUCCESS=0
for i in {1..100}; do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: abc123" http://localhost:8080/api/info)
  if [ "$CODE" == "200" ]; then
    ((SUCCESS++))
  fi
done
echo -n "[Token VIP - 100 requisições] "
if [ $SUCCESS -eq 100 ]; then
  echo "✅ PASS ($SUCCESS/100)"
  ((PASSED++))
else
  echo "❌ FAIL ($SUCCESS/100)"
  ((FAILED++))
fi

test_endpoint "Token VIP Request 101 (bloqueado)" "429" "-H 'API_KEY: abc123'"

echo ""
echo "============================================"
echo "  RESULTADO FINAL"
echo "============================================"
echo "✅ Testes Passados: $PASSED"
echo "❌ Testes Falhados: $FAILED"
echo ""

if [ $FAILED -eq 0 ]; then
  echo "🎉 TODOS OS TESTES PASSARAM!"
  exit 0
else
  echo "⚠️  ALGUNS TESTES FALHARAM"
  exit 1
fi
```

### 10.3 Executar Validação Completa

```bash
chmod +x validate-all.sh
./validate-all.sh
```

---

## ✅ Passo 11: Checklist Final de Validação

Execute cada item e marque como concluído:

### Funcionalidades

- [ ] Limitação por IP funciona (5 req/s bloqueio na 6ª)
- [ ] Limitação por Token funciona (10 req/s para token padrão)
- [ ] Token VIP funciona (100 req/s para abc123)
- [ ] Token sobrepõe limite de IP
- [ ] Bloqueio temporário funciona (5 minutos)
- [ ] Resposta 429 com mensagem correta
- [ ] Header API_KEY é reconhecido

### Arquitetura

- [ ] Middleware está separado da lógica
- [ ] Storage usa interface Strategy
- [ ] Redis está funcionando
- [ ] Dados são salvos no Redis
- [ ] Configuração via .env funciona

### Docker

- [ ] docker-compose up funciona
- [ ] Aplicação roda na porta 8080
- [ ] Redis está acessível
- [ ] Health checks funcionam
- [ ] Logs estão disponíveis

### Testes

- [ ] Testes unitários passam
- [ ] Testes de integração passam
- [ ] Script de teste funciona
- [ ] Teste de carga funciona

---

## 🐛 Troubleshooting

### Problema: "connection refused"

```bash
docker-compose ps
docker-compose logs app
```

### Problema: Redis não conecta

```bash
docker-compose logs redis
docker-compose exec redis redis-cli ping
```

### Problema: Sempre retorna 200

```bash
# Verificar se o middleware está aplicado
docker-compose logs app | grep -i "rate"

# Verificar .env
cat .env | grep RATE_LIMIT
```

### Resetar Tudo

```bash
docker-compose down -v
rm -rf coverage.out
docker-compose up --build -d
sleep 5
./validate-all.sh
```

---

## 📝 Conclusão

Se todos os testes acima passarem, seu Rate Limiter está **100% funcional** e atende todos os requisitos do projeto! 🎉
