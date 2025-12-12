@echo off
chcp 65001 >nul
REM Bug Bounty Lite - Windows Batch Script
REM Usage: run.bat [command]

if "%1"=="" goto help
if "%1"=="help" goto help
if "%1"=="run" goto run
if "%1"=="migrate" goto migrate
if "%1"=="init" goto init
if "%1"=="seed-projects" goto seed-projects
if "%1"=="seed-users" goto seed-users
if "%1"=="seed-reports" goto seed-reports
if "%1"=="seed-all" goto seed-all
if "%1"=="build" goto build
if "%1"=="test" goto test
goto help

:run
echo [Running Server]
go run cmd/server/main.go
goto end

:migrate
echo [Running Database Migration]
go run cmd/migrate/main.go
goto end

:init
echo [Initializing System Data]
go run cmd/init/main.go
goto end

:seed-projects
echo [Seeding Projects Test Data]
go run cmd/seed-projects/main.go
goto end

:seed-users
echo [Seeding Users Test Data]
go run cmd/seed-users/main.go
goto end

:seed-reports
echo [Seeding Reports Test Data]
go run cmd/seed-reports/main.go
goto end

:seed-all
echo [Seeding All Test Data]
go run cmd/seed-all/main.go
goto end

:build
echo [Building Project]
if not exist bin mkdir bin
go build -ldflags="-w -s" -o bin/server.exe ./cmd/server
echo [OK] Build completed: bin/server.exe
goto end

:test
echo [Running Tests]
go test -v ./...
goto end

:help
echo Bug Bounty Lite - Windows Batch Script
echo.
echo Usage: run.bat [command]
echo.
echo Available Commands:
echo   run              - Run server
echo   migrate          - Run database migration
echo   init             - Initialize system data
echo   seed-projects    - Seed projects test data
echo   seed-users       - Seed users test data
echo   seed-reports     - Seed reports test data
echo   seed-all         - Seed all test data
echo   build            - Build project
echo   test             - Run tests
echo   help             - Show this help message
echo.
echo Examples:
echo   run.bat migrate
echo   run.bat seed-all
echo   run.bat run
goto end

:end
