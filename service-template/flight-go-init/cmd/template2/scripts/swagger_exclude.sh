#!/bin/bash

# Script to auto-detect which common lib controllers are registered in RegisterRoutes
# and generate include/exclude lists for swagger generation
#
# Output format:
#   Line 1: Comma-separated list of packages to EXCLUDE
#   Line 2: Comma-separated list of directories to INCLUDE (for --dir flag)
#
# How to add new common lib controller:
#   Add entry to COMMON_LIB_CONTROLLERS with format:
#   ControllerName:module_path:subdir

set -e

CONTROLLER_DI="internal/controller/di.go"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# ==========================================
# COMMON LIB CONTROLLERS REGISTRY
# Format: "ControllerName:module_path:subdir"
# ==========================================
COMMON_LIB_CONTROLLERS="
SystemParamController:github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO:system_param
CredentialController:github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO:credential
CredentialMasterDataController:github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO:credential_master_data
PromotionController:github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO:promotion
IntegratorErrorMappingController:github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO:errormapper
"
# Add more controllers here as needed:
# ExampleController:github.com/tiket/TIX-FLIGHT-COMMON-LIB-GO:example_package

# ==========================================
# Functions
# ==========================================

is_controller_registered() {
    local controller_name=$1
    # Check if impl.ControllerName. exists and is NOT commented out
    local found=$(grep "impl\.${controller_name}\." "$CONTROLLER_DI" 2>/dev/null | grep -v "^[[:space:]]*//")
    [ -n "$found" ]
}

get_module_path() {
    local module=$1
    go list -m -f '{{.Dir}}' "$module" 2>/dev/null || echo ""
}

# ==========================================
# Main
# ==========================================

if [ ! -f "$CONTROLLER_DI" ]; then
    echo "Error: ${CONTROLLER_DI} not found" >&2
    exit 1
fi

exclude_packages=""
include_dirs=""

echo -e "${GREEN}Checking RegisterRoutes in ${CONTROLLER_DI}...${NC}" >&2

while IFS=: read -r controller module subdir; do
    [ -z "$controller" ] && continue
    
    full_package="${module}/${subdir}"
    
    if is_controller_registered "$controller"; then
        echo -e "  ${GREEN}✓${NC} ${controller} → included" >&2
        
        module_path=$(get_module_path "$module")
        if [ -n "$module_path" ]; then
            dir="${module_path}/${subdir}"
            [ -z "$include_dirs" ] && include_dirs="$dir" || include_dirs="${include_dirs},${dir}"
        fi
    else
        echo -e "  ${YELLOW}✗${NC} ${controller} → excluded" >&2
        [ -z "$exclude_packages" ] && exclude_packages="$full_package" || exclude_packages="${exclude_packages},${full_package}"
    fi
done <<< "$COMMON_LIB_CONTROLLERS"

# Output: line 1 = exclude, line 2 = include
echo "$exclude_packages"
echo "$include_dirs"

