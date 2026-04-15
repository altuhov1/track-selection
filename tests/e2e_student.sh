#!/bin/bash

# Цвета
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Базовый URL
BASE_URL="http://localhost:8080/api"

# Счетчики тестов
PASSED=0
FAILED=0

# Функция для вывода результата
test_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAILED${NC}"
        ((FAILED++))
    fi
}

# Функция для проверки наличия поля в JSON
check_json_field() {
    echo "$1" | jq -e ".$2" > /dev/null 2>&1
    return $?
}

echo "=========================================="
echo "     COMPLETE E2E TESTS"
echo "=========================================="
echo ""

# ============================================
# 1. ПРОВЕРКА СЕРВЕРА
# ============================================
echo -e "${BLUE}1. SERVER CHECK${NC}"
echo ""

if ! curl -s "${BASE_URL}/login" > /dev/null 2>&1; then
    echo -e "${RED}Server not running on port 8080${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Server running${NC}"
echo ""

# ============================================
# 2. РЕГИСТРАЦИЯ
# ============================================
echo -e "${BLUE}2. REGISTRATION${NC}"
echo ""

echo "Test 1: Register student"
RESPONSE=$(curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "e2e_test@example.com",
        "password": "testpass123",
        "role": "student",
        "first_name": "Тест",
        "last_name": "Тестов"
    }')

if echo "$RESPONSE" | grep -q "user created"; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

# ============================================
# 3. ЛОГИН
# ============================================
echo -e "${BLUE}3. LOGIN${NC}"
echo ""

echo "Test 2: Login"
RESPONSE=$(curl -s -X POST "${BASE_URL}/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"e2e_test@example.com","password":"testpass123"}')

TOKEN=$(echo "$RESPONSE" | jq -r '.token' 2>/dev/null)

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ] && [ ${#TOKEN} -gt 10 ]; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

# ============================================
# 4. GET /me
# ============================================
echo -e "${BLUE}4. GET /api/me${NC}"
echo ""

echo "Test 3: Get current user info"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me" \
    -H "Authorization: Bearer $TOKEN")

USER_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)
EMAIL=$(echo "$RESPONSE" | jq -r '.email' 2>/dev/null)
FIRST_NAME=$(echo "$RESPONSE" | jq -r '.first_name' 2>/dev/null)
LAST_NAME=$(echo "$RESPONSE" | jq -r '.last_name' 2>/dev/null)

if [ -n "$USER_ID" ] && [ "$EMAIL" = "e2e_test@example.com" ] && [ "$FIRST_NAME" = "Тест" ] && [ "$LAST_NAME" = "Тестов" ]; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

# ============================================
# 5. GET /me/info (пустые предпочтения)
# ============================================
echo -e "${BLUE}5. GET /api/me/info (initial)${NC}"
echo ""

echo "Test 4: Get preferences (should have random grades)"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me/info" \
    -H "Authorization: Bearer $TOKEN")

# Проверяем наличие полей
check_json_field "$RESPONSE" "professional_goals"
HAS_GOALS=$?
check_json_field "$RESPONSE" "grades"
HAS_GRADES=$?
check_json_field "$RESPONSE" "skills"
HAS_SKILLS=$?
check_json_field "$RESPONSE" "learning_style"
HAS_STYLE=$?

if [ $HAS_GOALS -eq 0 ] && [ $HAS_GRADES -eq 0 ] && [ $HAS_SKILLS -eq 0 ] && [ $HAS_STYLE -eq 0 ]; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response missing fields: $RESPONSE${NC}"
fi
echo ""

# ============================================
# 6. POST /me/edit-info (частичное обновление)
# ============================================
echo -e "${BLUE}6. POST /api/me/edit-info (partial)${NC}"
echo ""

echo "Test 5: Update learning style"
RESPONSE=$(curl -s -X POST "${BASE_URL}/me/edit-info" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"learning_style": 2}')

if echo "$RESPONSE" | grep -q "preferences updated"; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

# ============================================
# 7. Проверка обновления
# ============================================
echo -e "${BLUE}7. Verify partial update${NC}"
echo ""

echo "Test 6: Verify learning_style changed"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me/info" \
    -H "Authorization: Bearer $TOKEN")

LEARNING_STYLE=$(echo "$RESPONSE" | jq -r '.learning_style' 2>/dev/null)

if [ "$LEARNING_STYLE" = "2" ]; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Expected learning_style=2, got $LEARNING_STYLE${NC}"
fi
echo ""

# ============================================
# 8. POST /me/edit-info (полное обновление)
# ============================================
echo -e "${BLUE}8. POST /api/me/edit-info (full)${NC}"
echo ""

echo "Test 7: Full update all fields"
RESPONSE=$(curl -s -X POST "${BASE_URL}/me/edit-info" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
        "professional_goals": [1, 2, 3, 4],
        "skills": {
            "databases": 9,
            "system_architecture": 8,
            "algorithmic_programming": 7,
            "public_speaking": 6,
            "testing": 8,
            "analytics": 9,
            "machine_learning": 5,
            "os_knowledge": 7,
            "research_projects": 8
        },
        "learning_style": 3,
        "certificates": 1
    }')

if echo "$RESPONSE" | grep -q "preferences updated"; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

# ============================================
# 9. Проверка полного обновления
# ============================================
echo -e "${BLUE}9. Verify full update${NC}"
echo ""

echo "Test 8: Verify all fields updated"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me/info" \
    -H "Authorization: Bearer $TOKEN")

PROF_GOALS_COUNT=$(echo "$RESPONSE" | jq -r '.professional_goals | length' 2>/dev/null)
LEARNING_STYLE=$(echo "$RESPONSE" | jq -r '.learning_style' 2>/dev/null)
CERTIFICATES=$(echo "$RESPONSE" | jq -r '.certificates' 2>/dev/null)
SKILLS_DB=$(echo "$RESPONSE" | jq -r '.skills.databases' 2>/dev/null)

if [ "$PROF_GOALS_COUNT" = "4" ] && [ "$LEARNING_STYLE" = "3" ] && [ "$CERTIFICATES" = "1" ] && [ "$SKILLS_DB" = "9" ]; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Goals count: $PROF_GOALS_COUNT, style: $LEARNING_STYLE, cert: $CERTIFICATES, db: $SKILLS_DB${NC}"
fi
echo ""

# ============================================
# 10. GET /me/profile-completion
# ============================================
echo -e "${BLUE}10. GET /api/me/profile-completion${NC}"
echo ""

echo "Test 9: Get profile completion status"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me/profile-completion" \
    -H "Authorization: Bearer $TOKEN")

IS_COMPLETE=$(echo "$RESPONSE" | jq -r '.is_complete' 2>/dev/null)

if [ "$IS_COMPLETE" = "true" ]; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Expected is_complete=true, got $IS_COMPLETE${NC}"
fi
echo ""

# ============================================
# 11. ТЕСТЫ ОШИБОК
# ============================================
echo -e "${BLUE}11. ERROR HANDLING${NC}"
echo ""

echo "Test 10: Access /me without token"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me")

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

echo "Test 11: Invalid learning_style (should fail)"
RESPONSE=$(curl -s -X POST "${BASE_URL}/me/edit-info" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"learning_style": 99}')

if echo "$RESPONSE" | grep -q "INVALID_LEARNING_STYLE"; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

echo "Test 12: Invalid skill value (should fail)"
RESPONSE=$(curl -s -X POST "${BASE_URL}/me/edit-info" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"skills": {"databases": 99}}')

if echo "$RESPONSE" | grep -q "INVALID_SKILL_VALUE"; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

echo "Test 13: Invalid certificates (should fail)"
RESPONSE=$(curl -s -X POST "${BASE_URL}/me/edit-info" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"certificates": 5}')

if echo "$RESPONSE" | grep -q "INVALID_CERTIFICATE"; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

# ============================================
# 12. ОЧИСТКА ДАННЫХ
# ============================================
echo -e "${BLUE}12. DATA CLEANUP${NC}"
echo ""

echo "Test 14: Clear some data (remove professional_goals)"
RESPONSE=$(curl -s -X POST "${BASE_URL}/me/edit-info" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"professional_goals": []}')

if echo "$RESPONSE" | grep -q "preferences updated"; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Response: $RESPONSE${NC}"
fi
echo ""

echo "Test 15: Verify profile completion becomes false"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me/profile-completion" \
    -H "Authorization: Bearer $TOKEN")

IS_COMPLETE=$(echo "$RESPONSE" | jq -r '.is_complete' 2>/dev/null)

if [ "$IS_COMPLETE" = "false" ]; then
    test_result 0
else
    test_result 1
    echo -e "${RED}Expected is_complete=false, got $IS_COMPLETE${NC}"
fi
echo ""

# ============================================
# 13. РЕЗУЛЬТАТЫ
# ============================================
echo "=========================================="
echo -e "${BLUE}TEST RESULTS${NC}"
echo "=========================================="
echo -e "${GREEN}PASSED: $PASSED${NC}"
echo -e "${RED}FAILED: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! 🎉${NC}"
    exit 0
else
    echo -e "${RED}❌ SOME TESTS FAILED! ❌${NC}"
    exit 1
fi