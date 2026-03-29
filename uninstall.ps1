#requires -Version 5.1

param(
    [string] $InstallDir = ''
)

$ErrorActionPreference = 'Stop'

if (-not $InstallDir) {
    $InstallDir = Join-Path $env:LOCALAPPDATA 'Programs\clipper'
}

$exe = Join-Path $InstallDir 'clipper.exe'
if (Test-Path -LiteralPath $exe) {
    Remove-Item -LiteralPath $exe -Force
    Write-Host "Removed $exe"
} else {
    Write-Host "Not found (skipped): $exe"
}

$cfg = Join-Path $env:USERPROFILE '.config\clipper.json'
if (Test-Path -LiteralPath $cfg) {
    Remove-Item -LiteralPath $cfg -Force
    Write-Host "Removed $cfg"
}

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if (-not [string]::IsNullOrEmpty($userPath)) {
    $parts = $userPath -split ';' | Where-Object { $_ -and ($_ -ne $InstallDir) }
    $newPath = $parts -join ';'
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Host "Removed '$InstallDir' from your user PATH."
}

if (Test-Path -LiteralPath $InstallDir) {
    $left = @(Get-ChildItem -LiteralPath $InstallDir -Force -ErrorAction SilentlyContinue)
    if ($left.Count -eq 0) {
        Remove-Item -LiteralPath $InstallDir -Force -ErrorAction SilentlyContinue
    }
}

Write-Host ""
Write-Host 'Clipper uninstall finished. Open a new terminal so PATH changes apply.'
