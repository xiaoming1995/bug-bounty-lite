# Bug Bounty Lite - PowerShell Script
# Usage: .\run.ps1 [command]

param(
    [Parameter(Position=0)]
    [string]$Command = "help"
)

function Show-Help {
    Write-Host "Bug Bounty Lite - PowerShell Script" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "Usage: .\run.ps1 [command]" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Available Commands:" -ForegroundColor Green
    Write-Host "  run              - Run server"
    Write-Host "  migrate          - Run database migration"
    Write-Host "  init             - Initialize system data"
    Write-Host "  seed-projects    - Seed projects test data"
    Write-Host "  seed-users       - Seed users test data"
    Write-Host "  seed-reports     - Seed reports test data"
    Write-Host "  seed-all         - Seed all test data"
    Write-Host "  build            - Build project"
    Write-Host "  test             - Run tests"
    Write-Host "  clean            - Clean build artifacts"
    Write-Host "  help             - Show this help message"
    Write-Host ""
    Write-Host "Examples:" -ForegroundColor Yellow
    Write-Host "  .\run.ps1 migrate"
    Write-Host "  .\run.ps1 seed-all"
    Write-Host "  .\run.ps1 run"
}

function Run-Server {
    Write-Host "[Running Server]" -ForegroundColor Green
    go run cmd/server/main.go
}

function Run-Migration {
    Write-Host "[Running Database Migration]" -ForegroundColor Green
    go run cmd/migrate/main.go
}

function Run-Init {
    Write-Host "[Initializing System Data]" -ForegroundColor Green
    go run cmd/init/main.go
}

function Seed-Projects {
    Write-Host "[Seeding Projects Test Data]" -ForegroundColor Green
    go run cmd/seed-projects/main.go
}

function Seed-Users {
    Write-Host "[Seeding Users Test Data]" -ForegroundColor Green
    go run cmd/seed-users/main.go
}

function Seed-Reports {
    Write-Host "[Seeding Reports Test Data]" -ForegroundColor Green
    go run cmd/seed-reports/main.go
}

function Seed-All {
    Write-Host "[Seeding All Test Data]" -ForegroundColor Green
    go run cmd/seed-all/main.go
}

function Build-Project {
    Write-Host "[Building Project]" -ForegroundColor Green
    if (!(Test-Path "bin")) {
        New-Item -ItemType Directory -Path "bin" | Out-Null
    }
    go build -ldflags="-w -s" -o bin/server.exe ./cmd/server
    Write-Host "[OK] Build completed: bin/server.exe" -ForegroundColor Cyan
}

function Run-Tests {
    Write-Host "[Running Tests]" -ForegroundColor Green
    go test -v ./...
}

function Clean-Build {
    Write-Host "[Cleaning Build Artifacts]" -ForegroundColor Green
    if (Test-Path "bin") {
        Remove-Item -Recurse -Force "bin"
        Write-Host "[OK] Cleaned bin directory" -ForegroundColor Cyan
    }
    if (Test-Path "coverage.out") {
        Remove-Item "coverage.out"
    }
    if (Test-Path "coverage.html") {
        Remove-Item "coverage.html"
    }
}

# Main script logic
switch ($Command.ToLower()) {
    "run" { Run-Server }
    "migrate" { Run-Migration }
    "init" { Run-Init }
    "seed-projects" { Seed-Projects }
    "seed-users" { Seed-Users }
    "seed-reports" { Seed-Reports }
    "seed-all" { Seed-All }
    "build" { Build-Project }
    "test" { Run-Tests }
    "clean" { Clean-Build }
    "help" { Show-Help }
    default {
        Write-Host "Unknown command: $Command" -ForegroundColor Red
        Write-Host ""
        Show-Help
    }
}
