#!/bin/bash

# Script to generate Swagger documentation per distribution
# Usage: ./scripts/swagger_gen.sh <distribution_type>
# Example: ./scripts/swagger_gen.sh {{PROVIDER}}
#
# Features:
# - Uses @basePath from main.go (controllers use relative paths)
# - Smart include/exclude based on RegisterRoutes
# - Only registered controllers appear in swagger

set -e

DISTRIBUTION_TYPE=${1:-"{{PROVIDER}}"}
BASE_PATH="tix-flight-${DISTRIBUTION_TYPE}-integrator"
DOCS_DIR="docs/${DISTRIBUTION_TYPE}"
MAIN_PATH="cmd/${DISTRIBUTION_TYPE}/main.go"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}==================================================="
echo "Swagger Generation (Smart Mode)"
echo "===================================================${NC}"
echo "Distribution: ${DISTRIBUTION_TYPE}"
echo "Base Path: /${BASE_PATH}"
echo "Output: ${DOCS_DIR}"
echo -e "${GREEN}===================================================${NC}"

# Check swag
if ! command -v swag &> /dev/null; then
    echo -e "${YELLOW}swag not installed. Run: go install github.com/swaggo/swag/cmd/swag@latest${NC}"
    exit 1
fi

# Check main.go
if [ ! -f "$MAIN_PATH" ]; then
    echo -e "${YELLOW}Error: ${MAIN_PATH} not found${NC}"
    exit 1
fi

# Check @BasePath
if ! grep -q "@BasePath" "$MAIN_PATH"; then
    echo -e "${YELLOW}Warning: Add @BasePath to ${MAIN_PATH}: // @BasePath /${BASE_PATH}${NC}"
fi

mkdir -p "${DOCS_DIR}"

# Step 1: Detect registered controllers
echo ""
echo "Step 1: Detecting registered controllers..."
chmod +x "${SCRIPT_DIR}/swagger_exclude.sh"

RESULT=$("${SCRIPT_DIR}/swagger_exclude.sh")
EXCLUDE_PACKAGES=$(echo "$RESULT" | sed -n '1p')
INCLUDE_DIRS=$(echo "$RESULT" | sed -n '2p')

# Step 2: Build swag command
echo ""
echo "Step 2: Generating Swagger..."

SWAG_CMD="swag init -g ${MAIN_PATH} -o ${DOCS_DIR} --parseDependency --parseInternal --parseDepth 2"

if [ -n "$INCLUDE_DIRS" ]; then
    echo "  Including: ${INCLUDE_DIRS}"
    SWAG_CMD="${SWAG_CMD} --dir .,${INCLUDE_DIRS}"
fi

if [ -n "$EXCLUDE_PACKAGES" ]; then
    echo "  Excluding: ${EXCLUDE_PACKAGES}"
    SWAG_CMD="${SWAG_CMD} --exclude ${EXCLUDE_PACKAGES}"
fi

[ -z "$INCLUDE_DIRS" ] && [ -z "$EXCLUDE_PACKAGES" ] && echo "  No common lib controllers configured"

eval "$SWAG_CMD"

# Step 3: Update package name
echo ""
echo "Step 3: Updating package name..."
DOCS_GO="${DOCS_DIR}/docs.go"
if [ -f "$DOCS_GO" ]; then
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' "s|package docs|package ${DISTRIBUTION_TYPE}|g" "$DOCS_GO"
    else
        sed -i "s|package docs|package ${DISTRIBUTION_TYPE}|g" "$DOCS_GO"
    fi
fi

echo ""
echo -e "${GREEN}==================================================="
echo "Done! Output: ${DOCS_DIR}"
echo ""
echo "Swagger only includes registered controllers."
echo "To add new common lib controller, edit:"
echo "  scripts/swagger_exclude.sh → COMMON_LIB_CONTROLLERS"
echo -e "===================================================${NC}"
