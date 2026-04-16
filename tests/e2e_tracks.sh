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
echo "        E2E TESTS - TRACKS"
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
# СОЗДАНИЕ ТЕСТОВЫХ ПОЛЬЗОВАТЕЛЕЙ
# ============================================
echo -e "${BLUE}2. CREATE TEST USERS${NC}"
echo ""

echo "Creating admin user..."
curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "admin@example.com",
        "password": "adminpass123",
        "role": "admin",
        "first_name": "Admin",
        "last_name": "Adminov"
    }' > /dev/null

echo "Creating student user..."
curl -s -X POST "${BASE_URL}/register" \
    -H "Content-Type: application/json" \
    -d '{
        "email": "student@example.com",
        "password": "pass123",
        "role": "student",
        "first_name": "Student",
        "last_name": "Studentov"
    }' > /dev/null

echo -e "${GREEN}✓ Test users created${NC}"
echo ""

# ============================================
# GET /all-tracks (публичный)
# ============================================
echo -e "${BLUE}3. GET /api/all-tracks (public)${NC}"
echo ""

echo "Test 1: Get all tracks without auth"
RESPONSE=$(curl -s -X GET "${BASE_URL}/all-tracks")
COUNT=$(echo "$RESPONSE" | jq '. | length' 2>/dev/null)

if [ "$COUNT" -ge 3 ]; then
    echo -e "${GREEN}✓ PASSED (found $COUNT tracks)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED (expected at least 3 tracks, got $COUNT)${NC}"
    test_result 1
fi
echo ""

# ============================================
# ЛОГИН КАК АДМИН
# ============================================
echo -e "${BLUE}4. ADMIN LOGIN${NC}"
echo ""

echo "Test 2: Login as admin"
RESPONSE=$(curl -s -X POST "${BASE_URL}/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"admin@example.com","password":"adminpass123"}')

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
# ПРОПУСКАЕМ ОСТАЛЬНЫЕ ТЕСТЫ, ЕСЛИ НЕТ ТОКЕНА
# ============================================
if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo -e "${RED}No admin token available. Skipping remaining tests.${NC}"
    echo ""
    echo "=========================================="
    echo -e "${BLUE}TEST RESULTS${NC}"
    echo "=========================================="
    echo -e "${GREEN}PASSED: $PASSED${NC}"
    echo -e "${RED}FAILED: $FAILED${NC}"
    exit 1
fi

# ============================================
# СОЗДАНИЕ ТРЕКА (ADMIN) - С НОВОЙ СТРУКТУРОЙ
# ============================================
echo -e "${BLUE}5. POST /api/new-track (admin)${NC}"
echo ""

echo "Test 3: Create new track with semester-based curriculum"
RESPONSE=$(curl -s -X POST "${BASE_URL}/new-track" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
        "name": "DevOps-инжиниринг",
        "description": "Изучи CI/CD, Docker, Kubernetes и облачные технологии",
        "curriculum": {
            "semesters": [
                {
                    "number": 1,
                    "courses": [
                        {"name": "Основы Linux", "description": "Администрирование Linux", "is_elective": false},
                        {"name": "Сети", "description": "TCP/IP, DNS, HTTP", "is_elective": false}
                    ]
                },
                {
                    "number": 2,
                    "courses": [
                        {"name": "Docker", "description": "Контейнеризация", "is_elective": false},
                        {"name": "Kubernetes", "description": "Оркестрация", "is_elective": true, "options": ["K8s", "Nomad"]}
                    ]
                }
            ]
        },
        "requirements": [
            {"subject": "programming", "min_grade": 4},
            {"subject": "os_knowledge", "min_grade": 3}
        ],
        "teachers": ["Анна Иванова", "Петр Сидоров"],
        "difficulty": 3,
        "type": 1,
        "employment_prospects": 9,
        "alumni_reviews": 8,
        "web_link": "https://example.com/devops",
        "has_certificates": 1,
        "learning_style": 2,
        "desired_tech_skills": 8,
        "desired_math_skills": 5,
        "desired_soft_skills": 6,
        "professional_goals": [1, 2]
    }')

TRACK_ID=$(echo "$RESPONSE" | jq -r '.id' 2>/dev/null)

if [ -n "$TRACK_ID" ] && [ "$TRACK_ID" != "null" ] && [ ${#TRACK_ID} -gt 10 ]; then
    echo -e "${GREEN}✓ PASSED (created track with id: $TRACK_ID)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПРОВЕРКА, ЧТО ТРЕК ПОЯВИЛСЯ В СПИСКЕ
# ============================================
echo -e "${BLUE}6. VERIFY TRACK IN LIST${NC}"
echo ""

echo "Test 4: Verify new track appears in list"
RESPONSE=$(curl -s -X GET "${BASE_URL}/all-tracks")

if echo "$RESPONSE" | grep -q "DevOps-инжиниринг"; then
    echo -e "${GREEN}✓ PASSED (track found in list)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED (track not found)${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПРОВЕРКА СТРУКТУРЫ CURRICULUM
# ============================================
echo -e "${BLUE}7. VERIFY CURRICULUM STRUCTURE${NC}"
echo ""

echo "Test 5: Verify curriculum has semesters"
RESPONSE=$(curl -s -X GET "${BASE_URL}/all-tracks")
SEMESTER_COUNT=$(echo "$RESPONSE" | jq ".[] | select(.name==\"DevOps-инжиниринг\") | .curriculum.semesters | length" 2>/dev/null)

if [ "$SEMESTER_COUNT" -ge 2 ]; then
    echo -e "${GREEN}✓ PASSED (found $SEMESTER_COUNT semesters)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED (expected at least 2 semesters, got $SEMESTER_COUNT)${NC}"
    test_result 1
fi
echo ""

# ============================================
# ОБНОВЛЕНИЕ ТРЕКА (ADMIN)
# ============================================
echo -e "${BLUE}8. PUT /api/edit-track/{id} (admin)${NC}"
echo ""

echo "Test 6: Update track name, difficulty and type"
RESPONSE=$(curl -s -X PUT "${BASE_URL}/edit-track/${TRACK_ID}" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
        "name": "DevOps инжиниринг расширенный",
        "difficulty": 4,
        "type": 2,
        "employment_prospects": 10
    }')

if echo "$RESPONSE" | grep -q "track updated"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПРОВЕРКА ОБНОВЛЕНИЙ
# ============================================
echo -e "${BLUE}9. VERIFY UPDATES${NC}"
echo ""

echo "Test 7: Verify track was updated"
RESPONSE=$(curl -s -X GET "${BASE_URL}/all-tracks")

if echo "$RESPONSE" | grep -q "DevOps инжиниринг расширенный"; then
    echo -e "${GREEN}✓ PASSED (updated name found)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED (updated name not found)${NC}"
    test_result 1
fi
echo ""

# ============================================
# ПРОВЕРКА ПОЛЕЙ DIFFICULTY И TYPE
# ============================================
echo -e "${BLUE}10. VERIFY DIFFICULTY AND TYPE FIELDS${NC}"
echo ""

echo "Test 8: Verify difficulty and type were updated"
RESPONSE=$(curl -s -X GET "${BASE_URL}/all-tracks")
DIFFICULTY=$(echo "$RESPONSE" | jq ".[] | select(.name==\"DevOps инжиниринг расширенный\") | .difficulty" 2>/dev/null)
TYPE=$(echo "$RESPONSE" | jq ".[] | select(.name==\"DevOps инжиниринг расширенный\") | .type" 2>/dev/null)

if [ "$DIFFICULTY" = "4" ] && [ "$TYPE" = "2" ]; then
    echo -e "${GREEN}✓ PASSED (difficulty=$DIFFICULTY, type=$TYPE)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED (difficulty=$DIFFICULTY, type=$TYPE)${NC}"
    test_result 1
fi
echo ""

# ============================================
# ДОСТУП БЕЗ ТОКЕНА
# ============================================
echo -e "${BLUE}11. UNAUTHORIZED ACCESS${NC}"
echo ""

echo "Test 9: Create track without token (should fail)"
RESPONSE=$(curl -s -X POST "${BASE_URL}/new-track" \
    -H "Content-Type: application/json" \
    -d '{"name": "test"}')

if echo "$RESPONSE" | grep -q "UNAUTHORIZED"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ДОСТУП С ТОКЕНОМ СТУДЕНТА
# ============================================
echo -e "${BLUE}12. STUDENT LOGIN${NC}"
echo ""

echo "Test 10: Login as student"
RESPONSE=$(curl -s -X POST "${BASE_URL}/login" \
    -H "Content-Type: application/json" \
    -d '{"email":"student@example.com","password":"pass123"}')

STUDENT_TOKEN=$(echo "$RESPONSE" | jq -r '.token' 2>/dev/null)

if [ -n "$STUDENT_TOKEN" ] && [ "$STUDENT_TOKEN" != "null" ]; then
    echo -e "${GREEN}✓ PASSED (got student token)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

echo "Test 11: Create track with student token (should be forbidden)"
RESPONSE=$(curl -s -X POST "${BASE_URL}/new-track" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $STUDENT_TOKEN" \
    -d '{"name": "test"}')

if echo "$RESPONSE" | grep -q "FORBIDDEN"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# ОБНОВЛЕНИЕ НЕСУЩЕСТВУЮЩЕГО ТРЕКА
# ============================================
echo -e "${BLUE}13. UPDATE NON-EXISTENT TRACK${NC}"
echo ""

echo "Test 12: Update non-existent track"
RESPONSE=$(curl -s -X PUT "${BASE_URL}/edit-track/00000000-0000-0000-0000-000000000000" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{"name": "test"}')

if echo "$RESPONSE" | grep -q "NOT_FOUND"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

# ============================================
# УДАЛЕНИЕ ТРЕКА (ADMIN)
# ============================================
echo -e "${BLUE}14. DELETE TRACK (admin)${NC}"
echo ""

echo "Test 13: Delete track"
RESPONSE=$(curl -s -X DELETE "${BASE_URL}/delete-track/${TRACK_ID}" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN")

if echo "$RESPONSE" | grep -q "track deleted"; then
    echo -e "${GREEN}✓ PASSED${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED: $RESPONSE${NC}"
    test_result 1
fi
echo ""

echo "Test 14: Verify track no longer exists"
RESPONSE=$(curl -s -X GET "${BASE_URL}/all-tracks")

if echo "$RESPONSE" | grep -q "DevOps инжиниринг расширенный"; then
    echo -e "${RED}✗ FAILED (track still exists)${NC}"
    test_result 1
else
    echo -e "${GREEN}✓ PASSED (track deleted)${NC}"
    test_result 0
fi
echo ""

# ============================================
# ДЕФОЛТНЫЕ ТРЕКИ
# ============================================
echo -e "${BLUE}15. DEFAULT TRACKS${NC}"
echo ""

echo "Test 15: Check default tracks exist"
RESPONSE=$(curl -s -X GET "${BASE_URL}/all-tracks")

if echo "$RESPONSE" | grep -q "Backend-разработка" && \
   echo "$RESPONSE" | grep -q "Data Science" && \
   echo "$RESPONSE" | grep -q "Frontend-разработка"; then
    echo -e "${GREEN}✓ PASSED (all default tracks exist)${NC}"
    test_result 0
else
    echo -e "${RED}✗ FAILED (missing some default tracks)${NC}"
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
    echo -e "${GREEN}🎉 ALL TRACKS TESTS PASSED! 🎉${NC}"
    exit 0
else
    echo -e "${RED}❌ SOME TRACKS TESTS FAILED! ❌${NC}"
    exit 1
fi