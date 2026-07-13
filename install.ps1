# Vibe-Shield Windows installer
# Usage: powershell -ExecutionPolicy Bypass -File install.ps1

$ErrorActionPreference = "Stop"
$RepoRoot = $PSScriptRoot
$BinaryName = "vibe-shield.exe"

function Log($msg)  { Write-Host "==> $msg" -ForegroundColor Cyan }
function Ok($msg)   { Write-Host "OK  $msg" -ForegroundColor Green }
function Warn($msg)  { Write-Host "!   $msg" -ForegroundColor Yellow }
function Fail($msg)  { Write-Host "X   $msg" -ForegroundColor Red; exit 1 }

function Refresh-SessionPath {
    $env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
                [System.Environment]::GetEnvironmentVariable("Path", "User")
}

Log "Vibe-Shield Windows Installer"
Log ""

# Refresh PATH from registry before any command lookups (fixes stale terminal sessions)
Refresh-SessionPath

$defaultGoPath = Join-Path $env:USERPROFILE "go"
$installDir = Join-Path $defaultGoPath "bin"
$installedBin = Join-Path $installDir $BinaryName
$localBin = Join-Path $RepoRoot $BinaryName

New-Item -ItemType Directory -Force -Path $installDir | Out-Null

if (Get-Command go -ErrorAction SilentlyContinue) {
    Ok "Go found: $((go version 2>&1 | Out-String).Trim())"
    Log "Building and installing vibe-shield..."
    Push-Location $RepoRoot
    try {
        go install ./cmd/vibe-shield
        if ($LASTEXITCODE -ne 0) { Fail "go install failed" }
    } finally {
        Pop-Location
    }
    $installDir = Join-Path (go env GOPATH).Trim() "bin"
    $installedBin = Join-Path $installDir $BinaryName
} elseif (Test-Path $installedBin) {
    Ok "Using existing binary at $installedBin"
} elseif (Test-Path $localBin) {
    Log "Go not in session PATH - copying local build to $installedBin"
    Copy-Item -Path $localBin -Destination $installedBin -Force
    Ok "Copied local binary to $installedBin"
} else {
    Fail "Go is not available and no prebuilt binary found. Install Go (winget install GoLang.Go) or build first: go build -o vibe-shield.exe ./cmd/vibe-shield"
}

if (-not (Test-Path $installedBin)) {
    Fail "Expected binary not found at $installedBin"
}

Ok "Installed to $installedBin"

# Ensure install dir is on User PATH
$userPath = [System.Environment]::GetEnvironmentVariable("Path", "User")
$pathEntries = @($userPath -split ';' | Where-Object { $_ -ne '' })
$normalizedInstall = $installDir.TrimEnd('\')

if ($pathEntries -notcontains $normalizedInstall) {
    Log "Adding $normalizedInstall to User PATH..."
    $newPath = if ($userPath) { "$userPath;$normalizedInstall" } else { $normalizedInstall }
    [System.Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Ok "Added to User PATH"
} else {
    Ok "Already on User PATH"
}

Refresh-SessionPath

if (Get-Command vibe-shield -ErrorAction SilentlyContinue) {
    Ok "vibe-shield is available in this session"
} else {
    Warn "Restart Cursor completely, or refresh PATH in your current terminal:"
    Log '. .\scripts\refresh-path.ps1'
}

Log ""
Log 'Usage: vibe-shield [your-command] [args...]'
Log 'Example: vibe-shield python examples/fake_crash.py'
Log ""
Log 'Stale terminal? Run: . .\scripts\refresh-path.ps1'
