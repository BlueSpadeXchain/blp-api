#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "🚀 Running TypeScript tests..."

# Run TypeScript compilation
echo -n "Compiling TypeScript... "
npx tsc
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC}"
    
    # Run the compiled JavaScript
    echo -e "\n📝 Test output:"
    echo "----------------------------------------"
    node dist/test.js
    echo "----------------------------------------"
    
    echo -e "\n${GREEN}✓ Tests completed${NC}"
else
    echo -e "${RED}⨯ Compilation failed${NC}"
    exit 1
fi