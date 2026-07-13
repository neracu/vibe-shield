# Refresh session PATH from registry (fixes stale Cursor/VS Code terminals)
$RepoRoot = Split-Path $PSScriptRoot -Parent

$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
            [System.Environment]::GetEnvironmentVariable("Path", "User")

$found = [bool](Get-Command vibe-shield -ErrorAction SilentlyContinue)

if ($found) {
    Write-Host "OK  PATH refreshed. vibe-shield is now available." -ForegroundColor Green
} else {
    Write-Host "!   PATH refreshed but vibe-shield still not found." -ForegroundColor Yellow
    Write-Host "    Run: powershell -ExecutionPolicy Bypass -File $RepoRoot\install.ps1" -ForegroundColor Yellow
    Write-Host "    Or from project dir: .\vibe-shield.exe python examples/fake_crash.py" -ForegroundColor Yellow
}
