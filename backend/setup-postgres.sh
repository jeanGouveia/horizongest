#!/bin/bash

# PostgreSQL Setup Script for HorizonGest
# This script creates the user and database for the application

set -e

# Configuration
DB_USER="horizongest_user"
DB_PASSWORD="horizongest_secure_password"
DB_NAME="horizongest"
DB_TEST_NAME="horizongest_test"
DB_HOST="localhost"
DB_PORT="5432"

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${YELLOW}PostgreSQL Setup Script${NC}"
echo "=============================="
echo ""

# Check if PostgreSQL is running
echo -e "${YELLOW}Checking if PostgreSQL is running...${NC}"
if ! pg_isready -h $DB_HOST -p $DB_PORT > /dev/null 2>&1; then
    echo -e "${RED}ERROR: PostgreSQL is not running on ${DB_HOST}:${DB_PORT}${NC}"
    echo "Please start PostgreSQL and try again."
    echo "On Ubuntu/Debian: sudo service postgresql start"
    echo "On macOS: brew services start postgresql"
    exit 1
fi
echo -e "${GREEN}✓ PostgreSQL is running${NC}"
echo ""

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo -e "${RED}ERROR: psql command not found${NC}"
    echo "Please install PostgreSQL client tools."
    exit 1
fi

# Create user if it doesn't exist
echo -e "${YELLOW}Creating user '${DB_USER}' if it doesn't exist...${NC}"
USER_EXISTS=$(sudo -u postgres psql -tAc "SELECT 1 FROM pg_user WHERE usename = '${DB_USER}'" 2>/dev/null || echo "")

if [ "$USER_EXISTS" = "1" ]; then
    echo -e "${GREEN}✓ User '${DB_USER}' already exists${NC}"
else
    sudo -u postgres psql -c "CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASSWORD}';" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ User '${DB_USER}' created${NC}"
    else
        echo -e "${RED}ERROR: Failed to create user '${DB_USER}'${NC}"
        exit 1
    fi
fi

# Grant CREATEDB permission
echo -e "${YELLOW}Granting CREATEDB permission to user '${DB_USER}'...${NC}"
sudo -u postgres psql -c "ALTER USER ${DB_USER} CREATEDB;" > /dev/null 2>&1
echo -e "${GREEN}✓ CREATEDB permission granted${NC}"
echo ""

# Create database if it doesn't exist
echo -e "${YELLOW}Creating database '${DB_NAME}' if it doesn't exist...${NC}"
DB_EXISTS=$(sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname = '${DB_NAME}'" 2>/dev/null || echo "")

if [ "$DB_EXISTS" = "1" ]; then
    echo -e "${GREEN}✓ Database '${DB_NAME}' already exists${NC}"
else
    sudo -u postgres psql -c "CREATE DATABASE ${DB_NAME} OWNER ${DB_USER};" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Database '${DB_NAME}' created${NC}"
    else
        echo -e "${RED}ERROR: Failed to create database '${DB_NAME}'${NC}"
        exit 1
    fi
fi

# Create test database if it doesn't exist
echo -e "${YELLOW}Creating test database '${DB_TEST_NAME}' if it doesn't exist...${NC}"
DB_TEST_EXISTS=$(sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname = '${DB_TEST_NAME}'" 2>/dev/null || echo "")

if [ "$DB_TEST_EXISTS" = "1" ]; then
    echo -e "${GREEN}✓ Test database '${DB_TEST_NAME}' already exists${NC}"
else
    sudo -u postgres psql -c "CREATE DATABASE ${DB_TEST_NAME} OWNER ${DB_USER};" > /dev/null 2>&1
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ Test database '${DB_TEST_NAME}' created${NC}"
    else
        echo -e "${RED}ERROR: Failed to create test database '${DB_TEST_NAME}'${NC}"
        exit 1
    fi
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}PostgreSQL setup completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "Database connection details:"
echo "  Host: ${DB_HOST}"
echo "  Port: ${DB_PORT}"
echo "  Database: ${DB_NAME}"
echo "  Test Database: ${DB_TEST_NAME}"
echo "  User: ${DB_USER}"
echo ""
echo "You can now run the application with:"
echo "  cd backend && go run cmd/server/main.go"
echo ""
echo "To run tests:"
echo "  cd backend && go test ./..."
