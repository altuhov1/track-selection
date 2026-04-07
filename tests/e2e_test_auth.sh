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

# Тест 1: Успешная регистрация студента
echo "Test 1: Register student"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"student@example.com","password":"pass123","role":"student"}')

if echo "$RESPONSE" | grep -q "user created"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 2: Успешная регистрация админа
echo "Test 2: Register admin"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@example.com","password":"adminpass123","role":"admin"}')

if echo "$RESPONSE" | grep -q "user created"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 3: Регистрация с тем же email (ошибка)
echo "Test 3: Register duplicate email"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"student@example.com","password":"pass123","role":"student"}')

if echo "$RESPONSE" | grep -q "EMAIL_EXISTS"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 4: Регистрация с коротким паролем (ошибка)
echo "Test 4: Register weak password (len<6)"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"weak@example.com","password":"123","role":"student"}')

if echo "$RESPONSE" | grep -q "WEAK_PASSWORD"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 5: Регистрация с плохим email (ошибка)
echo "Test 5: Register invalid email"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"not-an-email","password":"pass123","role":"student"}')

if echo "$RESPONSE" | grep -q "INVALID_EMAIL"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 6: Регистрация с несуществующей ролью (ошибка)
echo "Test 6: Register invalid role"
RESPONSE=$(curl -s -X POST http://localhost:8080/register \
    -H "Content-Type: application/json" \
    -d '{"email":"invalid@example.com","password":"pass123","role":"superuser"}')

if echo "$RESPONSE" | grep -q "INVALID_ROLE"; then
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

# Тест 7: Успешный логин студента
echo "Test 7: Login student success"
RESPONSE=$(curl -s -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"email":"student@example.com","password":"pass123"}')

TOKEN=$(echo "$RESPONSE" | jq -r '.token' 2>/dev/null)

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ PASSED (got token)${NC}"
    STUDENT_TOKEN="$TOKEN"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 8: Успешный логин админа
echo "Test 8: Login admin success"
RESPONSE=$(curl -s -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@example.com","password":"adminpass123"}')

TOKEN=$(echo "$RESPONSE" | jq -r '.token' 2>/dev/null)

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ PASSED (got token)${NC}"
    ADMIN_TOKEN="$TOKEN"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 9: Логин с неверным паролем
echo "Test 9: Login wrong password"
RESPONSE=$(curl -s -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"email":"student@example.com","password":"wrongpass"}')

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 10: Логин с несуществующим email
echo "Test 10: Login non-existent email"
RESPONSE=$(curl -s -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{"email":"noexist@example.com","password":"pass123"}')

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# ============================================
# PROTECTED ENDPOINTS TESTS (опционально)
# ============================================
echo "========== PROTECTED ENDPOINTS =========="
echo ""

# Тест 11: Доступ к защищенному эндпоинту с токеном студента
if [ -n "$STUDENT_TOKEN" ]; then
    echo "Test 11: Access with student token"
    RESPONSE=$(curl -s -X GET http://localhost:8080/student/profile \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $STUDENT_TOKEN")
    
    if [ -n "$RESPONSE" ] && [ "$(echo "$RESPONSE" | jq -r '.error // empty')" = "" ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
    else
        echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    fi
else
    echo "Test 11: Skipped (no student token)"
fi
echo ""


# Тест 12: Доступ к админскому эндпоинту с токеном админа
if [ -n "$ADMIN_TOKEN" ]; then
    echo "Test 12: Admin endpoint with admin token"
    RESPONSE=$(curl -s -X GET http://localhost:8080/admin/tracks \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    
    # Может быть пустым списком или успешным ответом
    if [ -n "$RESPONSE" ] && [ "$(echo "$RESPONSE" | jq -r '.error // empty')" = "" ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
    else
        echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    fi
else
    echo "Test 13: Skipped (no admin token)"
fi
echo ""

echo "=========================================="
echo "                DONE"
echo "=========================================="