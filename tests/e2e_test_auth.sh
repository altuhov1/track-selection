#!/bin/bash

# Цвета
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=========================================="
echo "     E2E TESTS - REGISTER & LOGIN & ME"
echo "=========================================="
echo ""

# Базовый URL
BASE_URL="http://localhost:3000/api"

# Проверка что сервер запущен
if ! curl -s "${BASE_URL}/login" > /dev/null 2>&1; then
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
RESPONSE=$(curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "student@example.com",
        "password": "pass123",
        "role": "student",
        "first_name": "Иван",
        "last_name": "Петров"
    }')

if echo "$RESPONSE" | grep -q "user created"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 2: Успешная регистрация админа
echo "Test 2: Register admin"
RESPONSE=$(curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "admin@example.com",
        "password": "adminpass123",
        "role": "admin",
        "first_name": "Петр",
        "last_name": "Сидоров"
    }')

if echo "$RESPONSE" | grep -q "user created"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 3: Регистрация с тем же email (ошибка)
echo "Test 3: Register duplicate email"
RESPONSE=$(curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "student@example.com",
        "password": "pass123",
        "role": "student",
        "first_name": "Иван",
        "last_name": "Петров"
    }')

if echo "$RESPONSE" | grep -q "EMAIL_EXISTS"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 4: Регистрация с коротким паролем (ошибка)
echo "Test 4: Register weak password (len<6)"
RESPONSE=$(curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "weak@example.com",
        "password": "123",
        "role": "student",
        "first_name": "Тест",
        "last_name": "Тестов"
    }')

if echo "$RESPONSE" | grep -q "WEAK_PASSWORD"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 5: Регистрация с плохим email (ошибка)
echo "Test 5: Register invalid email"
RESPONSE=$(curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "not-an-email",
        "password": "pass123",
        "role": "student",
        "first_name": "Тест",
        "last_name": "Тестов"
    }')

if echo "$RESPONSE" | grep -q "INVALID_EMAIL"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 6: Регистрация с несуществующей ролью (ошибка)
echo "Test 6: Register invalid role"
RESPONSE=$(curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "invalid@example.com",
        "password": "pass123",
        "role": "superuser",
        "first_name": "Тест",
        "last_name": "Тестов"
    }')

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
RESPONSE=$(curl -s -X POST "${BASE_URL}/login" \
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
RESPONSE=$(curl -s -X POST "${BASE_URL}/login" \
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
RESPONSE=$(curl -s -X POST "${BASE_URL}/login" \
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
RESPONSE=$(curl -s -X POST "${BASE_URL}/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"noexist@example.com","password":"pass123"}')

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# ============================================
# GET /ME TESTS
# ============================================
echo "========== GET /api/me =========="
echo ""

# Тест 11: Получение информации о студенте
if [ -n "$STUDENT_TOKEN" ]; then
    echo "Test 11: Get student info from /me"
    RESPONSE=$(curl -s -X GET "${BASE_URL}/me" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $STUDENT_TOKEN")
    
    EMAIL=$(echo "$RESPONSE" | jq -r '.email' 2>/dev/null)
    ROLE=$(echo "$RESPONSE" | jq -r '.role' 2>/dev/null)
    FIRST_NAME=$(echo "$RESPONSE" | jq -r '.first_name' 2>/dev/null)
    LAST_NAME=$(echo "$RESPONSE" | jq -r '.last_name' 2>/dev/null)
    
    if [ "$EMAIL" = "student@example.com" ] && [ "$ROLE" = "student" ] && [ "$FIRST_NAME" = "Иван" ] && [ "$LAST_NAME" = "Петров" ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
    else
        echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    fi
else
    echo "Test 11: Skipped (no student token)"
fi
echo ""

# Тест 12: Получение информации об админе
if [ -n "$ADMIN_TOKEN" ]; then
    echo "Test 12: Get admin info from /me"
    RESPONSE=$(curl -s -X GET "${BASE_URL}/me" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ADMIN_TOKEN")
    
    EMAIL=$(echo "$RESPONSE" | jq -r '.email' 2>/dev/null)
    ROLE=$(echo "$RESPONSE" | jq -r '.role' 2>/dev/null)
    FIRST_NAME=$(echo "$RESPONSE" | jq -r '.first_name' 2>/dev/null)
    LAST_NAME=$(echo "$RESPONSE" | jq -r '.last_name' 2>/dev/null)
    
    if [ "$EMAIL" = "admin@example.com" ] && [ "$ROLE" = "admin" ] && [ "$FIRST_NAME" = "Петр" ] && [ "$LAST_NAME" = "Сидоров" ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
    else
        echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    fi
else
    echo "Test 12: Skipped (no admin token)"
fi
echo ""

# Тест 13: Доступ к /me без токена (ошибка)
echo "Test 13: Access /me without token (should fail)"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me" \
    -H "Content-Type: application/json")

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

# Тест 14: Доступ к /me с невалидным токеном (ошибка)
echo "Test 14: Access /me with invalid token (should fail)"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer invalid.token.here")

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
fi
echo ""

echo "=========================================="
echo "                DONE"
echo "=========================================="