#!/bin/bash

# Script para testar o rate limiter

echo "====================================="
echo "Rate Limiter Test Script"
echo "====================================="
echo ""

BASE_URL="http://localhost:8080"

echo "1. Testing IP-based rate limiting (5 requests per second)..."
echo "Sending 7 requests rapidly..."
for i in {1..7}; do
  echo "Request $i:"
  curl -i -s "$BASE_URL/api/info" 2>&1 | grep -E "HTTP|message"
  echo ""
done

echo ""
echo "====================================="
echo ""

echo "2. Testing Token-based rate limiting with token 'abc123' (100 requests per second)..."
echo "Sending 3 requests with token..."
for i in {1..3}; do
  echo "Request $i with token:"
  curl -i -s -H "API_KEY: abc123" "$BASE_URL/api/info" 2>&1 | grep -E "HTTP|message"
  echo ""
done

echo ""
echo "====================================="
echo ""

echo "3. Testing block duration..."
echo "Waiting 2 seconds and trying again (should still be blocked)..."
sleep 2
echo "Request after 2 seconds:"
curl -i -s "$BASE_URL/api/info" 2>&1 | grep -E "HTTP|message"

echo ""
echo "====================================="
echo "Test completed!"
echo "====================================="
