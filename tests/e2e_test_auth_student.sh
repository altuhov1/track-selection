#!/bin/bash

# Цвета
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=========================================="
echo "     E2E TESTS - REGISTER & LOGIN"
echo "=========================================="
echo ""

# Проверка что сервер запущен
if ! curl -s http://localhost:8080/login > /dev/null 2>&1; then
    echo -e "${RED}Server not running on port 8080${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Server running${NC}"
echo ""

# ============================================
# REGISTER TESTS
# ============================================
echo "========== REGISTER =========="
echo ""

# Тест 1: Успешная регистрация
echo "Test 1: Register student"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"pass123","role":"student"}')

if echo "$RESPONSE" | grep -q "user created"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 2: Регистрация с тем же email (ошибка)
echo "Test 2: Register duplicate email"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"pass123","role":"student"}')

if echo "$RESPONSE" | grep -q "EMAIL_EXISTS"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 3: Регистрация с коротким паролем (ошибка)
echo "Test 3: Register weak password (len<6)"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"weak@example.com","password":"123","role":"student"}')

if echo "$RESPONSE" | grep -q "WEAK_PASSWORD"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 4: Регистрация с плохим email (ошибка)
echo "Test 4: Register invalid email"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"not-an-email","password":"pass123","role":"student"}')

if echo "$RESPONSE" | grep -q "INVALID_EMAIL"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# ============================================
# LOGIN TESTS
# ============================================
echo "========== LOGIN =========="
echo ""

# Тест 5: Успешный логин
echo "Test 5: Login success"
RESPONSE=$(curl -s -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"pass123"}')

TOKEN=$(echo "$RESPONSE" | jq -r '.token' 2>/dev/null)

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ PASSED (got token)${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 6: Логин с неверным паролем
echo "Test 6: Login wrong password"
RESPONSE=$(curl -s -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"wrongpass"}')

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 7: Логин с несуществующим email
echo "Test 7: Login non-existent email"
RESPONSE=$(curl -s -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"email":"noexist@example.com","password":"pass123"}')

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

echo "=========================================="
echo "                DONE"
echo "=========================================="