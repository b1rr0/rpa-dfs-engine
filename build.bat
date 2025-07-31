@echo off
echo ========================================
echo    Building Facebook Auto Login
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
cd src
go mod tidy

echo.
echo Building application for Windows...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w" -o ../dist/facebook-login.exe .
cd ..

echo.
echo Build completed!
echo Executable file: dist/facebook-login.exe
echo.
pause 