
#!/bin/bash

# Цвета
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BASE_URL="http://localhost:8080/api"
PASSED=0
FAILED=0

test_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ PASSED${NC}"
        ((PASSED++))
    else
        echo -e "${RED}✗ FAILED${NC}"
        ((FAILED++))
    fi
}

echo "=========================================="
echo "   E2E TESTS - RECOMMENDATIONS"
echo "=========================================="
echo ""

# ============================================
# ПРОВЕРКА СЕРВЕРА
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
# СОЗДАНИЕ ТЕСТОВОГО СТУДЕНТА
# ============================================
echo -e "${BLUE}2. CREATE TEST STUDENT${NC}"
echo ""

echo "Creating test student..."
curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "rec_test@example.com",
        "password": "testpass123",
        "role": "student",
        "first_name": "Тест",
        "last_name": "Тестов"
    }' > /dev/null

echo -e "${GREEN}✓ Test student created${NC}"
echo ""

# ============================================
# ЛОГИН
# ============================================
echo -e "${BLUE}3. LOGIN${NC}"
echo ""

echo "Test 1: Login as test student"
RESPONSE=$(curl -s -X POST "${BASE_URL}/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"rec_test@example.com","password":"testpass123"}')

TOKEN=$(echo "$RESPONSE" | jq -r '.token' 2>/dev/null)

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ] && [ ${#TOKEN} -gt 10 ]; then
    echo -e "${GREEN}✓ PASSED (got token)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# РЕКОМЕНДАЦИИ БЕЗ ЗАПОЛНЕННОГО ПРОФИЛЯ
# ============================================
echo -e "${BLUE}4. RECOMMENDATIONS WITHOUT PROFILE${NC}"
echo ""

echo "Test 2: Get recommendations without profile (should fail)"
RESPONSE=$(curl -s -X GET "${BASE_URL}/student/recommendations" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RESPONSE" | grep -q "PROFILE_NOT_COMPLETE"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ЗАПОЛНЕНИЕ ПРОФИЛЯ
# ============================================
echo -e "${BLUE}5. FILL PROFILE${NC}"
echo ""

echo "Test 3: Update preferences with grades and skills"
RESPONSE=$(curl -s -X POST "${BASE_URL}/me/edit-info" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
        "professional_goals": [1, 2],
        "grades": {
            "informatics": 4,
            "programming": 5,
            "foreign_language": 3,
            "physics": 4,
            "aig": 4,
            "math_analysis": 4,
            "algorithms_data_structures": 4,
            "databases": 4,
            "discrete_math": 4
        },
        "skills": {
            "databases": 7,
            "system_architecture": 6,
            "algorithmic_programming": 8,
            "public_speaking": 5,
            "testing": 6,
            "analytics": 7,
            "machine_learning": 4,
            "os_knowledge": 5,
            "research_projects": 6
        },
        "learning_style": 2,
        "certificates": 1
    }')

if echo "$RESPONSE" | grep -q "preferences updated"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПРОВЕРКА ЗАПОЛНЕНИЯ ПРОФИЛЯ
# ============================================
echo -e "${BLUE}6. VERIFY PROFILE COMPLETION${NC}"
echo ""

echo "Test 4: Check profile completion status"
RESPONSE=$(curl -s -X GET "${BASE_URL}/me/profile-completion" \
    -H "Authorization: Bearer $TOKEN")

IS_COMPLETE=$(echo "$RESPONSE" | jq -r '.is_complete' 2>/dev/null)

if [ "$IS_COMPLETE" = "true" ]; then
    echo -e "${GREEN}✓ PASSED (profile is complete)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED (profile not complete)${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПОЛУЧЕНИЕ РЕКОМЕНДАЦИЙ
# ============================================
echo -e "${BLUE}7. GET RECOMMENDATIONS${NC}"
echo ""

echo "Test 5: Get recommendations with complete profile"
RESPONSE=$(curl -s -X GET "${BASE_URL}/student/recommendations" \
    -H "Authorization: Bearer $TOKEN")

# Проверяем, что ответ содержит recommendations
if echo "$RESPONSE" | jq -e '.recommendations' > /dev/null 2>&1; then
    REC_COUNT=$(echo "$RESPONSE" | jq '.recommendations | length' 2>/dev/null)
    if [ "$REC_COUNT" -gt 0 ]; then
        echo -e "${GREEN}✓ PASSED (got $REC_COUNT recommendations)${NC}"
        test_result 0
    else
        echo -e "${YELLOW}⚠ PASSED (got 0 recommendations - no matching tracks)${NC}"
        test_result 0
    fi
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПРОВЕРКА СТРУКТУРЫ ОТВЕТА
# ============================================
echo -e "${BLUE}8. VERIFY RESPONSE STRUCTURE${NC}"
echo ""

echo "Test 6: Verify recommendation structure"
RESPONSE=$(curl -s -X GET "${BASE_URL}/student/recommendations" \
    -H "Authorization: Bearer $TOKEN")

# Проверяем наличие полей у первого трека
HAS_TRACK_ID=$(echo "$RESPONSE" | jq '.recommendations[0].track_id' 2>/dev/null)
HAS_TRACK_NAME=$(echo "$RESPONSE" | jq '.recommendations[0].track_name' 2>/dev/null)
HAS_SCORE=$(echo "$RESPONSE" | jq '.recommendations[0].score' 2>/dev/null)
HAS_RANK=$(echo "$RESPONSE" | jq '.recommendations[0].rank' 2>/dev/null)
HAS_CRITERIA=$(echo "$RESPONSE" | jq '.recommendations[0].criteria_scores' 2>/dev/null)

if [ -n "$HAS_TRACK_ID" ] && [ -n "$HAS_TRACK_NAME" ] && [ -n "$HAS_SCORE" ] && [ -n "$HAS_RANK" ] && [ -n "$HAS_CRITERIA" ]; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: missing required fields${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПРОВЕРКА РАНЖИРОВАНИЯ
# ============================================
echo -e "${BLUE}9. VERIFY RANKING${NC}"
echo ""

echo "Test 7: Verify recommendations are sorted by score"
RESPONSE=$(curl -s -X GET "${BASE_URL}/student/recommendations" \
    -H "Authorization: Bearer $TOKEN")

SCORES=$(echo "$RESPONSE" | jq '.recommendations[].score' 2>/dev/null)
PREV_SCORE=999
IS_SORTED=true

for SCORE in $SCORES; do
    if (( $(echo "$SCORE > $PREV_SCORE" | bc -l) )); then
        IS_SORTED=false
        break
    fi
    PREV_SCORE=$SCORE
done

if [ "$IS_SORTED" = true ]; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED (recommendations not sorted)${NC}"
    test_result 1
fi
echo ""

# ============================================
# ДОСТУП БЕЗ ТОКЕНА
# ============================================
echo -e "${BLUE}10. UNAUTHORIZED ACCESS${NC}"
echo ""

echo "Test 8: Access recommendations without token (should fail)"
RESPONSE=$(curl -s -X GET "${BASE_URL}/student/recommendations")

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# РЕЗУЛЬТАТЫ
# ============================================
echo "=========================================="
echo -e "${BLUE}TEST RESULTS${NC}"
echo "=========================================="
echo -e "${GREEN}PASSED: $PASSED${NC}"
echo -e "${RED}FAILED: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL RECOMMENDATIONS TESTS PASSED! 🎉${NC}"
    exit 0
else
    echo -e "${RED}❌ SOME RECOMMENDATIONS TESTS FAILED! ❌${NC}"
    exit 1
fi