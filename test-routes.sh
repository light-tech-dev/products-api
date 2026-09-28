#!/bin/bash

# ═══════════════════════════════════════════════════════════════
#  Products API — Full Routes Test
# ═══════════════════════════════════════════════════════════════

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_URL="$BASE_URL/api/v1"

PASSED=0
FAILED=0
TOTAL=0

section() {
    echo ""
    echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}${CYAN}  $1${NC}"
    echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════════${NC}"
}

test_case() {
    TOTAL=$((TOTAL + 1))
    echo ""
    echo -e "${BOLD}${BLUE}▶ Test #$TOTAL:${NC} $1"
}

pass() {
    PASSED=$((PASSED + 1))
    echo -e "  ${GREEN}✅ PASS${NC} — $1"
}

fail() {
    FAILED=$((FAILED + 1))
    echo -e "  ${RED}❌ FAIL${NC} — $1"
}

expect() {
    local actual="$1"
    local expected="$2"
    local msg="$3"
    if [ "$actual" = "$expected" ]; then
        pass "$msg (got: $actual)"
    else
        fail "$msg (expected: $expected, got: $actual)"
    fi
}

expect_status() {
    local actual="$1"
    local expected="$2"
    local msg="$3"
    if [ "$actual" = "$expected" ]; then
        pass "$msg (HTTP $actual)"
    else
        fail "$msg (expected HTTP $expected, got $actual)"
    fi
}

json_get() {
    echo "$1" | jq -r "$2" 2>/dev/null
}

clear
echo ""
echo -e "${BOLD}${MAGENTA}╔═══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BOLD}${MAGENTA}║          🧪 Products API — Full Routes Test                   ║${NC}"
echo -e "${BOLD}${MAGENTA}╚═══════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "${DIM}Base URL: $BASE_URL${NC}"
echo -e "${DIM}Started:  $(date '+%Y-%m-%d %H:%M:%S')${NC}"

if ! command -v jq &> /dev/null; then
    echo -e "${RED}❌ jq is required${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}⏳ Waiting for server...${NC}"
for i in {1..10}; do
    if curl -s "$BASE_URL/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Server is up${NC}"
        break
    fi
    if [ $i -eq 10 ]; then
        echo -e "${RED}❌ Server not responding${NC}"
        exit 1
    fi
    sleep 1
done

# ─── SECTION 1: API INFO ───
section "1. API Info Endpoints"

test_case "GET /health"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" "$BASE_URL/health")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "Health returns 200"
STATUS=$(json_get "$BODY" ".status")
expect "$STATUS" "healthy" "Health status"

test_case "GET /api/v1/"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" "$API_URL/")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "API Info returns 200"
NAME=$(json_get "$BODY" ".name")
expect "$NAME" "Products API" "API name"

test_case "GET /api/v1/endpoints"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" "$API_URL/endpoints")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "Endpoints returns 200"
TOTAL_EPS=$(json_get "$BODY" ".total")
if [ "$TOTAL_EPS" -gt 0 ]; then
    pass "Total endpoints: $TOTAL_EPS"
else
    fail "No endpoints"
fi

test_case "GET /api/v1/schema"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" "$API_URL/schema")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "Schema returns 200"

# ─── SECTION 2: AUTH ───
section "2. Authentication"

TEST_USER="testuser_$(date +%s)"
TEST_EMAIL="test_${TEST_USER}@test.com"

test_case "POST /auth/register"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USER\",\"email\":\"$TEST_EMAIL\",\"password\":\"pass123\"}")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "201" "Register returns 201"
USER_ID=$(json_get "$BODY" ".data.id")
if [ -n "$USER_ID" ] && [ "$USER_ID" != "null" ]; then
    pass "User created: ID=$USER_ID"
else
    fail "User ID not returned"
fi

test_case "POST /auth/register — duplicate"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USER\",\"email\":\"other@test.com\",\"password\":\"pass123\"}")
expect_status "$HTTP_CODE" "409" "Duplicate returns 409"

test_case "POST /auth/register — validation"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X POST "$API_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d '{"username":"x","email":"bad","password":"123"}')
expect_status "$HTTP_CODE" "400" "Validation returns 400"

test_case "POST /auth/login"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USER\",\"password\":\"pass123\"}")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "Login returns 200"
TOKEN=$(json_get "$BODY" ".data.token")
if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    pass "Token: ${TOKEN:0:30}..."
else
    fail "No token"
fi

test_case "POST /auth/login — wrong password"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X POST "$API_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"$TEST_USER\",\"password\":\"wrong\"}")
expect_status "$HTTP_CODE" "401" "Wrong password returns 401"

test_case "GET /auth/me"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/auth/me")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "Me returns 200"

test_case "GET /auth/me — without token"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" "$API_URL/auth/me")
expect_status "$HTTP_CODE" "401" "No token returns 401"

# ─── SECTION 3: PRODUCTS ───
section "3. Products CRUD"

TEST_SKU="SKU-$(date +%s)"

test_case "POST /products"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X POST "$API_URL/products" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Test Product\",\"sku\":\"$TEST_SKU\",\"price\":99.99,\"stock\":10,\"category\":\"Test\"}")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "201" "Create returns 201"
PRODUCT_ID=$(json_get "$BODY" ".data.id")
if [ -n "$PRODUCT_ID" ] && [ "$PRODUCT_ID" != "null" ]; then
    pass "Product created: ID=$PRODUCT_ID"
else
    fail "No product ID"
fi

test_case "POST /products — duplicate SKU"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X POST "$API_URL/products" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Dup\",\"sku\":\"$TEST_SKU\",\"price\":50}")
expect_status "$HTTP_CODE" "409" "Duplicate SKU returns 409"

test_case "GET /products/:id"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products/$PRODUCT_ID")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "Get by ID returns 200"

test_case "GET /products/:id — not found"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products/99999")
expect_status "$HTTP_CODE" "404" "Not found returns 404"

test_case "GET /products/:id — invalid"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products/abc")
expect_status "$HTTP_CODE" "400" "Invalid ID returns 400"

test_case "GET /products"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "List returns 200"
TOTAL_P=$(json_get "$BODY" ".total")
if [ "$TOTAL_P" -ge 1 ]; then
    pass "Found $TOTAL_P products"
else
    fail "No products"
fi

test_case "GET /products — pagination"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products?page=1&per_page=5")
expect_status "$HTTP_CODE" "200" "Pagination returns 200"

test_case "PUT /products/:id"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X PUT "$API_URL/products/$PRODUCT_ID" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"price":89.99}')
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "Update returns 200"
PRICE=$(json_get "$BODY" ".data.price")
expect "$PRICE" "89.99" "Price updated"

# ─── SECTION 4: FILTERS ───
section "4. Filters & Search"

test_case "GET /products?category=Test"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products?category=Test")
expect_status "$HTTP_CODE" "200" "Category filter returns 200"

test_case "GET /products?min_price=50"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products?min_price=50")
expect_status "$HTTP_CODE" "200" "Price filter returns 200"

test_case "GET /products?q=Test"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products?q=Test")
expect_status "$HTTP_CODE" "200" "Search returns 200"

# ─── SECTION 5: ITEMS ───
section "5. Items (Nested)"

test_case "POST /products/:id/items"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X POST "$API_URL/products/$PRODUCT_ID/items" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"Item 1","quantity":2,"unit_price":10.50}')
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "201" "Create item returns 201"
ITEM_ID=$(json_get "$BODY" ".data.id")
if [ -n "$ITEM_ID" ] && [ "$ITEM_ID" != "null" ]; then
    pass "Item created: ID=$ITEM_ID"
else
    fail "No item ID"
fi

test_case "GET /products/:id/items"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products/$PRODUCT_ID/items")
BODY=$(cat /tmp/resp.json)
expect_status "$HTTP_CODE" "200" "List items returns 200"
COUNT=$(json_get "$BODY" ".count")
if [ "$COUNT" -ge 1 ]; then
    pass "Found $COUNT items"
else
    fail "No items"
fi

test_case "DELETE /items/:id"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X DELETE "$API_URL/items/$ITEM_ID" \
    -H "Authorization: Bearer $TOKEN")
expect_status "$HTTP_CODE" "204" "Delete item returns 204"

# ─── SECTION 6: DELETE PRODUCT ───
section "6. Delete Products"

test_case "DELETE /products/:id"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -X DELETE "$API_URL/products/$PRODUCT_ID" \
    -H "Authorization: Bearer $TOKEN")
expect_status "$HTTP_CODE" "204" "Delete returns 204"

test_case "GET /products/:id — verify deleted"
HTTP_CODE=$(curl -s -o /tmp/resp.json -w "%{http_code}" \
    -H "Authorization: Bearer $TOKEN" \
    "$API_URL/products/$PRODUCT_ID")
expect_status "$HTTP_CODE" "404" "Deleted returns 404"

# ─── RESULTS ───
echo ""
echo ""
echo -e "${BOLD}${MAGENTA}╔═══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BOLD}${MAGENTA}║                        TEST RESULTS                            ║${NC}"
echo -e "${BOLD}${MAGENTA}╚═══════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "  ${BOLD}Total:${NC}    $TOTAL"
echo -e "  ${GREEN}${BOLD}Passed:${NC}   $PASSED"
echo -e "  ${RED}${BOLD}Failed:${NC}   $FAILED"

if [ $TOTAL -gt 0 ]; then
    PERCENT=$((PASSED * 100 / TOTAL))
else
    PERCENT=0
fi

echo ""
echo -e "  ${BOLD}Success Rate:${NC} ${PERCENT}%"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}${BOLD}  🎉 ALL TESTS PASSED!${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}${BOLD}  ❌ $FAILED TESTS FAILED${NC}"
    echo ""
    exit 1
fi
