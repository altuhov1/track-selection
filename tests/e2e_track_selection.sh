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
echo "   E2E TESTS - TRACK SELECTION"
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
        "email": "select_test@example.com",
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
    -d '{"email":"select_test@example.com","password":"testpass123"}')

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
# ПОЛУЧЕНИЕ СПИСКА ТРЕКОВ
# ============================================
echo -e "${BLUE}4. GET ALL TRACKS${NC}"
echo ""

echo "Test 2: Get all tracks"
RESPONSE=$(curl -s -X GET "${BASE_URL}/all-tracks" \
    -H "Authorization: Bearer $TOKEN")

TRACK_ID=$(echo "$RESPONSE" | jq -r '.[0].id' 2>/dev/null)
TRACK_ID_2=$(echo "$RESPONSE" | jq -r '.[1].id' 2>/dev/null)

if [ -n "$TRACK_ID" ] && [ "$TRACK_ID" != "null" ]; then
    echo -e "${GREEN}✓ PASSED (got track id: $TRACK_ID)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ВЫБОР ТРЕКА
# ============================================
echo -e "${BLUE}5. SELECT TRACK${NC}"
echo ""

echo "Test 3: Select a track"
RESPONSE=$(curl -s -X POST "${BASE_URL}/student/select-track" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"track_id\": \"$TRACK_ID\"}")

if echo "$RESPONSE" | grep -q "track selected"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ВЫБОР ВТОРОГО ТРЕКА
# ============================================
echo -e "${BLUE}6. SELECT SECOND TRACK${NC}"
echo ""

echo "Test 4: Select second track"
RESPONSE=$(curl -s -X POST "${BASE_URL}/student/select-track" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"track_id\": \"$TRACK_ID_2\"}")

if echo "$RESPONSE" | grep -q "track selected"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПОЛУЧЕНИЕ ВЫБРАННЫХ ТРЕКОВ
# ============================================
echo -e "${BLUE}7. GET SELECTED TRACKS${NC}"
echo ""

echo "Test 5: Get selected tracks"
RESPONSE=$(curl -s -X GET "${BASE_URL}/student/selected-tracks" \
    -H "Authorization: Bearer $TOKEN")

SELECTED_COUNT=$(echo "$RESPONSE" | jq '.tracks | length' 2>/dev/null)

if [ "$SELECTED_COUNT" -ge 2 ]; then
    echo -e "${GREEN}✓ PASSED (found $SELECTED_COUNT selected tracks)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПРОВЕРКА ЧТО ВЫБРАННЫЙ ТРЕК ЕСТЬ В СПИСКЕ
# ============================================
echo -e "${BLUE}8. VERIFY SELECTED TRACK${NC}"
echo ""

echo "Test 6: Verify selected track is in the list"
RESPONSE=$(curl -s -X GET "${BASE_URL}/student/selected-tracks" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RESPONSE" | grep -q "$TRACK_ID"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED (track not found)${NC}"
    test_result 1
fi
echo ""

# ============================================
# ОТМЕНА ВЫБОРА ТРЕКА
# ============================================
echo -e "${BLUE}9. UNSELECT TRACK${NC}"
echo ""

echo "Test 7: Unselect track"
RESPONSE=$(curl -s -X DELETE "${BASE_URL}/student/unselect-track/${TRACK_ID}" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RESPONSE" | grep -q "track unselected"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПРОВЕРКА ЧТО ТРЕК УДАЛЕН ИЗ ВЫБРАННЫХ
# ============================================
echo -e "${BLUE}10. VERIFY TRACK REMOVED${NC}"
echo ""

echo "Test 8: Verify track removed from selected"
RESPONSE=$(curl -s -X GET "${BASE_URL}/student/selected-tracks" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RESPONSE" | grep -q "$TRACK_ID"; then
    echo -e "${RED}✗ FAILED (track still selected)${NC}"
    test_result 1
else
    echo -e "${GREEN}✓ PASSED (track removed)${NC}"
    test_result 0
fi
echo ""

# ============================================
# ОТМЕНА ВЫБОРА НЕСУЩЕСТВУЮЩЕГО ТРЕКА
# ============================================
echo -e "${BLUE}11. UNSELECT NON-EXISTENT TRACK${NC}"
echo ""

echo "Test 9: Unselect non-existent track"
RESPONSE=$(curl -s -X DELETE "${BASE_URL}/student/unselect-track/00000000-0000-0000-0000-000000000000" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RESPONSE" | grep -q "NOT_FOUND"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ВЫБОР ТРЕКА БЕЗ ТОКЕНА
# ============================================
echo -e "${BLUE}12. UNAUTHORIZED ACCESS${NC}"
echo ""

echo "Test 10: Select track without token (should fail)"
RESPONSE=$(curl -s -X POST "${BASE_URL}/student/select-track" \
    -H "Content-Type: application/json" \
    -d "{\"track_id\": \"$TRACK_ID\"}")

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ВЫБОР ТРЕКА С НЕВЕРНЫМ ID
# ============================================
echo -e "${BLUE}13. SELECT INVALID TRACK${NC}"
echo ""

echo "Test 11: Select track with invalid id"
RESPONSE=$(curl -s -X POST "${BASE_URL}/student/select-track" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"track_id": "invalid-id"}')

if echo "$RESPONSE" | grep -q "NOT_FOUND"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПОВТОРНЫЙ ВЫБОР ТОГО ЖЕ ТРЕКА
# ============================================
echo -e "${BLUE}14. SELECT SAME TRACK AGAIN${NC}"
echo ""

echo "Test 12: Select same track again"
RESPONSE=$(curl -s -X POST "${BASE_URL}/student/select-track" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{\"track_id\": \"$TRACK_ID_2\"}")

# Должен вернуть успех (upsert)
if echo "$RESPONSE" | grep -q "track selected"; then
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
    echo -e "${GREEN}🎉 ALL TRACK SELECTION TESTS PASSED! 🎉${NC}"
    exit 0
else
    echo -e "${RED}❌ SOME TRACK SELECTION TESTS FAILED! ❌${NC}"
    exit 1
fi