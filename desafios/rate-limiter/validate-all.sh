#!/bin/bash

echo "============================================"
echo "  VALIDAÇÃO COMPLETA DO RATE LIMITER"
echo "============================================"
echo ""

FAILED=0
PASSED=0

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para testar
test_endpoint() {
  local name="$1"
  local expected_code="$2"
  local header="$3"
  
  echo -n "[$name] "
  
  if [ -z "$header" ]; then
    CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:8080/api/info)
  else
    CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "$header" http://localhost:8080/api/info)
  fi
  
  if [ "$CODE" == "$expected_code" ]; then
    echo -e "${GREEN}✅ PASS${NC} (HTTP $CODE)"
    ((PASSED++))
  else
    echo -e "${RED}❌ FAIL${NC} (Expected $expected_code, got $CODE)"
    ((FAILED++))
  fi
}

# Verificar se o servidor está rodando
echo "Verificando se o servidor está rodando..."
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
  echo -e "${RED}❌ Servidor não está rodando!${NC}"
  echo "Execute: docker-compose up -d"
  exit 1
fi
echo -e "${GREEN}✅ Servidor está rodando${NC}"
echo ""

# Limpar Redis
echo "Limpando Redis..."
docker-compose exec -T redis redis-cli FLUSHALL > /dev/null 2>&1
sleep 1

echo ""
echo "=== Testes de Limitação por IP (5 req/s) ==="
test_endpoint "IP Request 1" "200"
test_endpoint "IP Request 2" "200"
test_endpoint "IP Request 3" "200"
test_endpoint "IP Request 4" "200"
test_endpoint "IP Request 5" "200"
test_endpoint "IP Request 6 (bloqueado)" "429"
test_endpoint "IP Request 7 (ainda bloqueado)" "429"

# Limpar Redis
docker-compose exec -T redis redis-cli FLUSHALL > /dev/null 2>&1
sleep 1

echo ""
echo "=== Testes de Limitação por Token (10 req/s) ==="
for i in {1..10}; do
  test_endpoint "Token Request $i" "200" "API_KEY: test-token"
done
test_endpoint "Token Request 11 (bloqueado)" "429" "API_KEY: test-token"

# Limpar Redis
docker-compose exec -T redis redis-cli FLUSHALL > /dev/null 2>&1
sleep 1

echo ""
echo "=== Teste de Prioridade Token sobre IP ==="
# Bloquear IP
echo "Bloqueando IP primeiro..."
for i in {1..6}; do
  curl -s -o /dev/null http://localhost:8080/api/info
done
test_endpoint "IP bloqueado" "429"
test_endpoint "Token funciona com IP bloqueado" "200" "API_KEY: abc123"

echo ""
echo "=== Teste de Token VIP abc123 (100 req/s) ==="
# Limpar Redis
docker-compose exec -T redis redis-cli FLUSHALL > /dev/null 2>&1
sleep 1

# Fazer 100 requisições com token VIP
echo "Enviando 100 requisições com token VIP..."
SUCCESS=0
for i in {1..100}; do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" -H "API_KEY: abc123" http://localhost:8080/api/info)
  if [ "$CODE" == "200" ]; then
    ((SUCCESS++))
  fi
  # Mostrar progresso a cada 20 requisições
  if [ $((i % 20)) -eq 0 ]; then
    echo "  Progresso: $i/100..."
  fi
done

echo -n "[Token VIP - 100 requisições] "
if [ $SUCCESS -eq 100 ]; then
  echo -e "${GREEN}✅ PASS${NC} ($SUCCESS/100)"
  ((PASSED++))
else
  echo -e "${RED}❌ FAIL${NC} ($SUCCESS/100)"
  ((FAILED++))
fi

test_endpoint "Token VIP Request 101 (bloqueado)" "429" "API_KEY: abc123"

# Limpar Redis
docker-compose exec -T redis redis-cli FLUSHALL > /dev/null 2>&1
sleep 1

echo ""
echo "=== Teste de Mensagem de Erro ==="
MESSAGE=$(curl -s http://localhost:8080/api/info)
for i in {1..6}; do
  curl -s -o /dev/null http://localhost:8080/api/info
done
MESSAGE=$(curl -s http://localhost:8080/api/info)
echo -n "[Mensagem de erro 429] "
if echo "$MESSAGE" | grep -q "you have reached the maximum number of requests or actions allowed within a certain time frame"; then
  echo -e "${GREEN}✅ PASS${NC}"
  ((PASSED++))
else
  echo -e "${RED}❌ FAIL${NC} - Mensagem incorreta"
  echo "Recebido: $MESSAGE"
  ((FAILED++))
fi

echo ""
echo "=== Verificando Redis ==="
echo -n "[Redis está armazenando dados] "
KEYS=$(docker-compose exec -T redis redis-cli KEYS "*" 2>/dev/null | wc -l)
if [ $KEYS -gt 0 ]; then
  echo -e "${GREEN}✅ PASS${NC} ($KEYS chaves encontradas)"
  ((PASSED++))
else
  echo -e "${RED}❌ FAIL${NC} (Nenhuma chave encontrada)"
  ((FAILED++))
fi

echo ""
echo "============================================"
echo "  RESULTADO FINAL"
echo "============================================"
echo -e "${GREEN}✅ Testes Passados: $PASSED${NC}"
if [ $FAILED -gt 0 ]; then
  echo -e "${RED}❌ Testes Falhados: $FAILED${NC}"
else
  echo -e "${GREEN}❌ Testes Falhados: $FAILED${NC}"
fi
echo ""

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}🎉 TODOS OS TESTES PASSARAM!${NC}"
  echo ""
  echo "Seu Rate Limiter está funcionando perfeitamente! ✨"
  exit 0
else
  echo -e "${YELLOW}⚠️  ALGUNS TESTES FALHARAM${NC}"
  echo ""
  echo "Verifique os logs com: docker-compose logs app"
  exit 1
fi
