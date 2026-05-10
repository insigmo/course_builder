# install.ps1 — installs the latest course-builder binary on Windows
# Usage (run in PowerShell as Administrator):
#   irm https://raw.githubusercontent.com/insigmo/course_builder/main/install.ps1 | iex

param(
  [string]$InstallDir = "$env:LOCALAPPDATA\Programs\course-builder"
)

$ErrorActionPreference = "Stop"
$Repo = "insigmo/course_builder"
$Binary = "course-builder"

# Detect architecture
$Arch = if ([System.Environment]::Is64BitOperatingSystem) {
  if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
  Write-Error "32-bit Windows is not supported."; exit 1
}

Write-Host "Fetching latest release..." -ForegroundColor Cyan

$Release = Invoke-RestMethod "https://api.github.com/repos/$Repo/releases/latest"
$Tag = $Release.tag_name

if (-not $Tag) {
  Write-Error "Could not determine latest release."
  exit 1
}

$Filename = "$Binary-windows-$Arch.exe"
$Url = "https://github.com/$Repo/releases/download/$Tag/$Filename"

Write-Host "Downloading $Filename ($Tag)..." -ForegroundColor Cyan

New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
$Dest = Join-Path $InstallDir "$Binary.exe"
Invoke-WebRequest -Uri $Url -OutFile $Dest

# Add to PATH if not already there
$UserPath = [System.Environment]::GetEnvironmentVariable("PATH", "User")
if ($UserPath -notlike "*$InstallDir*") {
  [System.Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
  Write-Host "Added $InstallDir to PATH (restart terminal to apply)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "✅ course-builder $Tag installed to $Dest" -ForegroundColor Green
Write-Host "   Run: course-builder --help"
