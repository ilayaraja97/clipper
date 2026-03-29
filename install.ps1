#requires -Version 5.1

param(
    [string] $InstallDir = ''
)

$ErrorActionPreference = 'Stop'

if (-not $InstallDir) {
    $InstallDir = Join-Path $env:LOCALAPPDATA 'Programs\clipper'
}

$RepoOwner = 'ilayaraja97'
$RepoName = 'clipper'

function Get-GoArch {
    switch ($env:PROCESSOR_ARCHITECTURE) {
        'AMD64' { return 'amd64' }
        'x86' { return '386' }
        'ARM64' { return 'arm64' }
        default {
            throw "Unsupported processor architecture: $($env:PROCESSOR_ARCHITECTURE). Clipper provides Windows builds for amd64 and 386."
        }
    }
}

$headers = @{
    'User-Agent' = 'clipper-install-script'
    'Accept'     = 'application/vnd.github+json'
}

$release = Invoke-RestMethod -Uri "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest" -Headers $headers
$tag = $release.tag_name
$goarch = Get-GoArch

$assetName = "clipper_{0}_windows_{1}.tar.gz" -f $tag, $goarch
$asset = $release.assets | Where-Object { $_.name -eq $assetName } | Select-Object -First 1
if (-not $asset) {
    $names = ($release.assets | ForEach-Object { $_.name }) -join ', '
    throw "No GitHub release asset named '$assetName'. Available assets: $names"
}

$url = $asset.browser_download_url
Write-Host "Downloading Clipper $tag for Windows ($goarch)..."
Write-Host "  $url"
Write-Host ""

$tmp = Join-Path ([System.IO.Path]::GetTempPath()) ('clipper-install-' + [Guid]::NewGuid().ToString('n'))
New-Item -ItemType Directory -Path $tmp | Out-Null

try {
    $archive = Join-Path $tmp 'dist.tar.gz'
    Invoke-WebRequest -Uri $url -OutFile $archive -UseBasicParsing
    Push-Location $tmp
    try {
        & tar.exe -xzf $archive
    } finally {
        Pop-Location
    }

    $exeSrc = Join-Path $tmp 'clipper.exe'
    if (-not (Test-Path -LiteralPath $exeSrc)) {
        throw 'clipper.exe not found after extracting the release archive.'
    }

    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $exeDest = Join-Path $InstallDir 'clipper.exe'
    Copy-Item -LiteralPath $exeSrc -Destination $exeDest -Force
} finally {
    Remove-Item -LiteralPath $tmp -Recurse -Force -ErrorAction SilentlyContinue
}

$userPath = [Environment]::GetEnvironmentVariable('Path', 'User')
if ($userPath -notlike "*$InstallDir*") {
    $newPath = if ([string]::IsNullOrEmpty($userPath)) { $InstallDir } else { "$InstallDir;$userPath" }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Host "Added to your user PATH: $InstallDir"
}

$env:Path = "$InstallDir;$env:Path"

Write-Host ""
Write-Host "Installation of version $tag complete."
Write-Host "  clipper.exe -> $exeDest"
Write-Host ""
Write-Host "Open a new terminal, then run: clipper"
