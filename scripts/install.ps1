param(
  [string]$Version = "",
  [string]$Repo = "Senasphy/unisocket",
  [string]$BinName = "unisocket",
  [string]$InstallDir = "$env:LOCALAPPDATA\Programs\unisocket\bin"
)

$ErrorActionPreference = "Stop"

function Resolve-Arch {
  switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString().ToLower()) {
    "x64" { "amd64" }
    "arm64" { "arm64" }
    default { throw "Unsupported architecture" }
  }
}

function Resolve-LatestVersion {
  (Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" `
    -Headers @{ "User-Agent" = "unisocket-installer" }).tag_name
}

if (-not $Version) { $Version = $env:UNISOCKET_VERSION }
if (-not $Version) { $Version = Resolve-LatestVersion }
if (-not $Version) { throw "Failed to resolve version" }

$Version = $Version.TrimStart("v")

$arch = Resolve-Arch
$artifact = "${BinName}_${Version}_windows_${arch}.zip"
$url = "https://github.com/$Repo/releases/download/v$Version/$artifact"

$tmp = Join-Path $env:TEMP ("unisocket-" + [guid]::NewGuid())
New-Item -ItemType Directory -Path $tmp | Out-Null

try {
  $zip = Join-Path $tmp $artifact
  Write-Host "Downloading $url ..."
  Invoke-WebRequest -Uri $url -OutFile $zip

  Expand-Archive $zip -DestinationPath $tmp -Force

  $src = Join-Path $tmp "$BinName.exe"
  $dest = Join-Path $InstallDir "$BinName.exe"

  New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
  Copy-Item $src $dest -Force

  $path = [Environment]::GetEnvironmentVariable("Path", "User")
  if ($path -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$path;$InstallDir", "User")
    Write-Host "Added to PATH. Restart terminal."
  }

  Write-Host "Installed $BinName $Version → $dest"
}
finally {
  Remove-Item $tmp -Recurse -Force -ErrorAction SilentlyContinue
}
