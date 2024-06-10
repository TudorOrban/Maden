@echo off
SETLOCAL

:: Ensure Docker is running
docker info >nul 2>&1
IF %ERRORLEVEL% NEQ 0 (
    echo Docker is not running. Please start Docker and try again.
    exit /b 1
)

:: Navigate to the directory containing docker-compose.yml
cd /d %~dp0

:: Start the application using Docker Compose
echo Starting Maden API server and etcd...
docker-compose up -d

:: Allow some time for services to start
timeout /t 15

echo Maden API server and etcd are now running.
pause
ENDLOCAL