@echo off
echo ========================================
echo    Building RPA DFS Engine
echo ========================================
echo.

REM Check for Go installation
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Go not found in system
    echo Install Go from https://golang.org/
    pause
    exit /b 1
)

echo Go found: 
go version

echo.
echo Initializing module and downloading dependencies...
go mod tidy

echo.
echo Building application for Windows...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o dist/rpa-dfs-engine.exe ./cmd/rpa-dfs-engine

echo.
echo Build completed!
echo Executable file: dist/rpa-dfs-engine.exe
echo.
pause 