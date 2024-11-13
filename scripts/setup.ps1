# Requires -Version 5.1
# Setup script for the Go project
Write-Host "Setting up environment..." -ForegroundColor Cyan

# Change to the project directory if not already there
$ProjectPath = Split-Path -Path $MyInvocation.MyCommand.Path -Parent
Set-Location -Path $ProjectPath

# Check if Go is installed and accessible
if (!(Get-Command "go" -ErrorAction SilentlyContinue)) {
    Write-Error "Go is not installed or not in the PATH."
    Exit
}

# Install dependencies using Go
try {
    Write-Host "Installing Go dependencies..."
    go mod tidy
    Write-Host "Dependencies installed successfully." -ForegroundColor Green
} catch {
    Write-Error "Failed to install Go dependencies. Error: $_"
    Exit
}

# Set up environment variables
Write-Host "Setting up environment variables..."
$env:GEMINI_API_KEY = "your-api-key-here"
[Environment]::SetEnvironmentVariable("GEMINI_API_KEY", "your-api-key-here", [EnvironmentVariableTarget]::User)

Write-Host "Environment setup completed." -ForegroundColor Green

# Optionally, open the project directory in File Explorer
# Start-Process explorer.exe $ProjectPath



