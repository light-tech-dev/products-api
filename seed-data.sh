#!/bin/bash

# ═══════════════════════════════════════════════════════════════
#  Products API — Seed Data Script
#  Fills all models with realistic test data
# ═══════════════════════════════════════════════════════════════

# ─── Colors ───────────────────────────────────────────────────
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

# ─── Config ───────────────────────────────────────────────────
BASE_URL="${BASE_URL:-http://localhost:8080}"
API_URL="$BASE_URL/api/v1"

# ─── Counters ─────────────────────────────────────────────────
USERS_CREATED=0
PRODUCTS_CREATED=0
ITEMS_CREATED=0
ERRORS=0

# ─── Helpers ──────────────────────────────────────────────────

banner() {
    echo ""
    echo -e "${BOLD}${MAGENTA}╔═══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BOLD}${MAGENTA}║          🌱 Products API — Seed Data                          ║${NC}"
    echo -e "${BOLD}${MAGENTA}╚═══════════════════════════════════════════════════════════════╝${NC}"
}

section() {
    echo ""
    echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}${CYAN}  $1${NC}"
    echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════════${NC}"
}

step() {
    echo -e "  ${BLUE}▶${NC} $1"
}

ok() {
    echo -e "    ${GREEN}✅${NC} $1"
}

err() {
    ERRORS=$((ERRORS + 1))
    echo -e "    ${RED}❌${NC} $1"
}

info() {
    echo -e "    ${DIM}$1${NC}"
}

# create_user: ينشئ مستخدم
create_user() {
    local username="$1"
    local email="$2"
    local password="$3"

    local http_code=$(curl -s -o /tmp/seed_resp.json -w "%{http_code}" \
        -X POST "$API_URL/auth/register" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$username\",\"email\":\"$email\",\"password\":\"$password\"}")

    if [ "$http_code" = "201" ]; then
        USERS_CREATED=$((USERS_CREATED + 1))
        ok "User: $username ($email)"
    elif [ "$http_code" = "409" ]; then
        info "User already exists: $username"
    else
        err "Failed to create user: $username (HTTP $http_code)"
    fi
}

# login: يسجّل الدخول ويعيد token
login() {
    local username="$1"
    local password="$2"

    local token=$(curl -s -X POST "$API_URL/auth/login" \
        -H "Content-Type: application/json" \
        -d "{\"username\":\"$username\",\"password\":\"$password\"}" | jq -r .data.token 2>/dev/null)

    echo "$token"
}

# create_product: ينشئ منتج
create_product() {
    local token="$1"
    local name="$2"
    local sku="$3"
    local price="$4"
    local stock="$5"
    local category="$6"
    local description="$7"

    local http_code=$(curl -s -o /tmp/seed_resp.json -w "%{http_code}" \
        -X POST "$API_URL/products" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\":\"$name\",
            \"sku\":\"$sku\",
            \"price\":$price,
            \"stock\":$stock,
            \"category\":\"$category\",
            \"description\":\"$description\"
        }")

    if [ "$http_code" = "201" ]; then
        PRODUCTS_CREATED=$((PRODUCTS_CREATED + 1))
        local id=$(cat /tmp/seed_resp.json | jq -r .data.id)
        echo "$id"
    elif [ "$http_code" = "409" ]; then
        info "Product already exists: $sku"
        echo ""
    else
        err "Failed to create product: $name (HTTP $http_code)"
        echo ""
    fi
}

# create_item: ينشئ item
create_item() {
    local token="$1"
    local product_id="$2"
    local name="$3"
    local quantity="$4"
    local unit_price="$5"
    local notes="$6"

    local http_code=$(curl -s -o /tmp/seed_resp.json -w "%{http_code}" \
        -X POST "$API_URL/products/$product_id/items" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\":\"$name\",
            \"quantity\":$quantity,
            \"unit_price\":$unit_price,
            \"notes\":\"$notes\"
        }")

    if [ "$http_code" = "201" ]; then
        ITEMS_CREATED=$((ITEMS_CREATED + 1))
        ok "Item: $name (product #$product_id)"
    else
        err "Failed to create item: $name (HTTP $http_code)"
    fi
}

# ═══════════════════════════════════════════════════════════════
#  START
# ═══════════════════════════════════════════════════════════════

clear
banner

echo ""
echo -e "${DIM}Base URL: $BASE_URL${NC}"
echo -e "${DIM}Started:  $(date '+%Y-%m-%d %H:%M:%S')${NC}"

# ─── Check server ─────────────────────────────────────────────
echo ""
echo -e "${YELLOW}⏳ Checking server...${NC}"

if ! curl -s "$BASE_URL/health" > /dev/null 2>&1; then
    echo -e "${RED}❌ Server not responding at $BASE_URL${NC}"
    echo -e "${YELLOW}   Run: go run ./cmd/server${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Server is up${NC}"

# ═══════════════════════════════════════════════════════════════
#  1. CREATE USERS
# ═══════════════════════════════════════════════════════════════

section "1. Creating Users"

step "Admin user"
create_user "admin" "admin@example.com" "admin123456"

step "Regular users"
create_user "john_doe" "john@example.com" "john123456"
create_user "jane_smith" "jane@example.com" "jane123456"
create_user "ali_ahmed" "ali@example.com" "ali123456"
create_user "sara_mohamed" "sara@example.com" "sara123456"

# ═══════════════════════════════════════════════════════════════
#  2. LOGIN AS ADMIN
# ═══════════════════════════════════════════════════════════════

section "2. Logging in as Admin"

step "Login admin"
TOKEN=$(login "admin" "admin123456")

if [ -n "$TOKEN" ] && [ "$TOKEN" != "null" ]; then
    ok "Token obtained: ${TOKEN:0:30}..."
else
    err "Failed to login"
    exit 1
fi

# ═══════════════════════════════════════════════════════════════
#  3. CREATE PRODUCTS — ELECTRONICS
# ═══════════════════════════════════════════════════════════════

section "3. Creating Electronics Products"

step "iPhones"
ID_IPHONE_15=$(create_product "$TOKEN" "iPhone 15 Pro" "APPLE-IP15PRO-001" 1299.99 50 "Electronics" "Latest Apple flagship with A17 Pro chip")
ID_IPHONE_14=$(create_product "$TOKEN" "iPhone 14" "APPLE-IP14-002" 899.99 100 "Electronics" "Previous generation iPhone")
ID_IPHONE_SE=$(create_product "$TOKEN" "iPhone SE" "APPLE-IPSE-003" 429.99 200 "Electronics" "Budget iPhone with A15 chip")

step "Samsung"
ID_GALAXY_S24=$(create_product "$TOKEN" "Samsung Galaxy S24 Ultra" "SAMSUNG-S24U-004" 1199.99 75 "Electronics" "Samsung flagship with S Pen")
ID_GALAXY_A54=$(create_product "$TOKEN" "Samsung Galaxy A54" "SAMSUNG-A54-005" 449.99 150 "Electronics" "Mid-range Samsung phone")

step "Google Pixel"
ID_PIXEL_8=$(create_product "$TOKEN" "Google Pixel 8 Pro" "GOOGLE-P8P-006" 999.99 60 "Electronics" "Google flagship with AI features")

# ═══════════════════════════════════════════════════════════════
#  4. CREATE PRODUCTS — COMPUTERS
# ═══════════════════════════════════════════════════════════════

section "4. Creating Computer Products"

step "Laptops"
ID_MACBOOK_PRO=$(create_product "$TOKEN" "MacBook Pro 16 M3" "APPLE-MBP16-007" 2499.99 15 "Computers" "Powerful MacBook with M3 Max chip")
ID_MACBOOK_AIR=$(create_product "$TOKEN" "MacBook Air 15" "APPLE-MBA15-008" 1299.99 40 "Computers" "Thin and light MacBook")
ID_DELL_XPS=$(create_product "$TOKEN" "Dell XPS 15" "DELL-XPS15-009" 1799.99 25 "Computers" "Premium Windows laptop")
ID_LENOVO_X1=$(create_product "$TOKEN" "Lenovo ThinkPad X1 Carbon" "LENOVO-X1C-010" 1699.99 30 "Computers" "Business-class ThinkPad")

step "Desktops"
ID_IMAC=$(create_product "$TOKEN" "iMac 24 M3" "APPLE-IMAC24-011" 1499.99 20 "Computers" "All-in-one desktop")

# ═══════════════════════════════════════════════════════════════
#  5. CREATE PRODUCTS — ACCESSORIES
# ═══════════════════════════════════════════════════════════════

section "5. Creating Accessories"

step "Audio"
ID_AIRPODS_PRO=$(create_product "$TOKEN" "AirPods Pro 2" "APPLE-APP2-012" 249.99 200 "Accessories" "Noise cancelling earbuds")
ID_AIRPODS_MAX=$(create_product "$TOKEN" "AirPods Max" "APPLE-APMAX-013" 549.99 50 "Accessories" "Over-ear headphones")
ID_SONY_WH=$(create_product "$TOKEN" "Sony WH-1000XM5" "SONY-WH1000XM5-014" 399.99 80 "Accessories" "Premium noise cancelling headphones")

step "Watches"
ID_APPLE_WATCH=$(create_product "$TOKEN" "Apple Watch Ultra 2" "APPLE-AWU2-015" 799.99 60 "Accessories" "Rugged Apple Watch")
ID_GALAXY_WATCH=$(create_product "$TOKEN" "Samsung Galaxy Watch 6" "SAMSUNG-GW6-016" 349.99 90 "Accessories" "Android smartwatch")

step "Chargers"
ID_MAGSAFE=$(create_product "$TOKEN" "MagSafe Charger" "APPLE-MSCHG-017" 39.99 500 "Accessories" "Wireless magnetic charger")
ID_USBC_CABLE=$(create_product "$TOKEN" "USB-C to Lightning Cable" "APPLE-USBC-018" 19.99 1000 "Accessories" "1m certified cable")

# ═══════════════════════════════════════════════════════════════
#  6. CREATE PRODUCTS — FURNITURE
# ═══════════════════════════════════════════════════════════════

section "6. Creating Furniture Products"

step "Office"
ID_DESK=$(create_product "$TOKEN" "Standing Desk Pro" "FURN-DESK-019" 599.99 25 "Furniture" "Electric adjustable standing desk")
ID_CHAIR=$(create_product "$TOKEN" "Ergonomic Office Chair" "FURN-CHAIR-020" 449.99 40 "Furniture" "Herman Miller style chair")

# ═══════════════════════════════════════════════════════════════
#  7. CREATE ITEMS FOR PRODUCTS
# ═══════════════════════════════════════════════════════════════

section "7. Creating Items for Products"

if [ -n "$ID_IPHONE_15" ]; then
    step "Items for iPhone 15 Pro (#$ID_IPHONE_15)"
    create_item "$TOKEN" "$ID_IPHONE_15" "Protective Case" 1 49.99 "Black leather case"
    create_item "$TOKEN" "$ID_IPHONE_15" "Screen Protector" 2 19.99 "Tempered glass"
    create_item "$TOKEN" "$ID_IPHONE_15" "Fast Charger 30W" 1 39.99 "USB-C PD charger"
    create_item "$TOKEN" "$ID_IPHONE_15" "Lightning Cable 2m" 2 29.99 "Certified cable"
fi

if [ -n "$ID_MACBOOK_PRO" ]; then
    step "Items for MacBook Pro (#$ID_MACBOOK_PRO)"
    create_item "$TOKEN" "$ID_MACBOOK_PRO" "Protective Sleeve" 1 79.99 "Leather sleeve"
    create_item "$TOKEN" "$ID_MACBOOK_PRO" "USB-C Hub 7-in-1" 1 89.99 "Multi-port adapter"
    create_item "$TOKEN" "$ID_MACBOOK_PRO" "AppleCare+ 3 Years" 1 399.99 "Extended warranty"
fi

if [ -n "$ID_GALAXY_S24" ]; then
    step "Items for Galaxy S24 Ultra (#$ID_GALAXY_S24)"
    create_item "$TOKEN" "$ID_GALAXY_S24" "S Pen Replacement" 1 49.99 "Original S Pen"
    create_item "$TOKEN" "$ID_GALAXY_S24" "Silicone Case" 1 34.99 "Blue silicone"
    create_item "$TOKEN" "$ID_GALAXY_S24" "45W Super Fast Charger" 1 59.99 "Samsung charger"
fi

if [ -n "$ID_AIRPODS_PRO" ]; then
    step "Items for AirPods Pro (#$ID_AIRPODS_PRO)"
    create_item "$TOKEN" "$ID_AIRPODS_PRO" "Replacement Ear Tips" 1 9.99 "4 sizes included"
    create_item "$TOKEN" "$ID_AIRPODS_PRO" "Wireless Charging Case" 1 99.99 "MagSafe compatible"
fi

if [ -n "$ID_DELL_XPS" ]; then
    step "Items for Dell XPS (#$ID_DELL_XPS)"
    create_item "$TOKEN" "$ID_DELL_XPS" "Extended Warranty 4 Years" 1 299.99 "Dell ProSupport"
    create_item "$TOKEN" "$ID_DELL_XPS" "Backpack" 1 89.99 "Dell branded"
fi

if [ -n "$ID_MACBOOK_AIR" ]; then
    step "Items for MacBook Air (#$ID_MACBOOK_AIR)"
    create_item "$TOKEN" "$ID_MACBOOK_AIR" "USB-C to HDMI Adapter" 1 69.99 "4K support"
    create_item "$TOKEN" "$ID_MACBOOK_AIR" "Sleeve 13 inch" 1 39.99 "Felt material"
fi

# ═══════════════════════════════════════════════════════════════
#  8. FINAL STATS
# ═══════════════════════════════════════════════════════════════

section "8. Fetching Final Statistics"

step "Total users"
USERS=$(curl -s -H "Authorization: Bearer $TOKEN" "$API_URL/auth/me" | jq -r .data.username 2>/dev/null)
ok "Current user: $USERS"

step "Total products"
TOTAL_PRODUCTS=$(curl -s -H "Authorization: Bearer $TOKEN" "$API_URL/products?per_page=100" | jq -r .total 2>/dev/null)
ok "Total products: $TOTAL_PRODUCTS"

step "Products by category"
echo -e "    ${DIM}Electronics:$(curl -s -H "Authorization: Bearer $TOKEN" "$API_URL/products?category=Electronics" | jq -r .total 2>/dev/null)${NC}"
echo -e "    ${DIM}Computers:  $(curl -s -H "Authorization: Bearer $TOKEN" "$API_URL/products?category=Computers" | jq -r .total 2>/dev/null)${NC}"
echo -e "    ${DIM}Accessories:$(curl -s -H "Authorization: Bearer $TOKEN" "$API_URL/products?category=Accessories" | jq -r .total 2>/dev/null)${NC}"
echo -e "    ${DIM}Furniture:  $(curl -s -H "Authorization: Bearer $TOKEN" "$API_URL/products?category=Furniture" | jq -r .total 2>/dev/null)${NC}"

# ═══════════════════════════════════════════════════════════════
#  RESULTS
# ═══════════════════════════════════════════════════════════════

echo ""
echo ""
echo -e "${BOLD}${MAGENTA}╔═══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${BOLD}${MAGENTA}║                     SEED DATA SUMMARY                          ║${NC}"
echo -e "${BOLD}${MAGENTA}╚═══════════════════════════════════════════════════════════════╝${NC}"
echo ""
echo -e "  ${BOLD}Users created:${NC}     $USERS_CREATED"
echo -e "  ${BOLD}Products created:${NC}  $PRODUCTS_CREATED"
echo -e "  ${BOLD}Items created:${NC}     $ITEMS_CREATED"
echo -e "  ${BOLD}Errors:${NC}            $ERRORS"
echo ""

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}${BOLD}  🎉 SEED DATA COMPLETE!${NC}"
    echo ""
    echo -e "${DIM}  Login credentials:${NC}"
    echo -e "${DIM}    Username: admin${NC}"
    echo -e "${DIM}    Password: admin123456${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}${BOLD}  ⚠️  COMPLETED WITH $ERRORS ERRORS${NC}"
    echo ""
    exit 1
fi
