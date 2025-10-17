# MKV Mender Installation Script for Windows
# Supports: Windows 10/11 (x64, ARM64)

$ErrorActionPreference = "Stop"

# Configuration
$GitHubRepo = "quentinsteinke/mkvmender"
$BinaryName = "mkvmender.exe"
$InstallDir = "$env:LOCALAPPDATA\mkvmender"

# Colors for output
function Write-Info {
    param([string]$Message)
    Write-Host "==> " -ForegroundColor Blue -NoNewline
    Write-Host $Message
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Error-Message {
    param([string]$Message)
    Write-Host "✗ " -ForegroundColor Red -NoNewline
    Write-Host $Message
}

function Write-Warn {
    param([string]$Message)
    Write-Host "! " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default {
            Write-Error-Message "Unsupported architecture: $arch"
            exit 1
        }
    }
}

# Get latest release version from GitHub
function Get-LatestVersion {
    Write-Info "Fetching latest version..."

    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$GitHubRepo/releases/latest"
        $version = $response.tag_name

        if ([string]::IsNullOrEmpty($version)) {
            throw "Failed to get version"
        }

        Write-Success "Latest version: $version"
        return $version
    }
    catch {
        Write-Error-Message "Failed to fetch latest version: $_"
        exit 1
    }
}

# Download binary
function Download-Binary {
    param(
        [string]$Architecture,
        [string]$Version
    )

    $filename = "mkvmender-windows-${Architecture}.exe"
    $downloadUrl = "https://github.com/$GitHubRepo/releases/download/$Version/$filename"
    $tmpFile = "$env:TEMP\$filename"

    Write-Info "Downloading mkvmender from $downloadUrl..."

    try {
        # Use TLS 1.2
        [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

        # Download with progress
        $ProgressPreference = 'SilentlyContinue'
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tmpFile -UseBasicParsing

        Write-Success "Downloaded successfully"
        return $tmpFile
    }
    catch {
        Write-Error-Message "Failed to download binary: $_"
        exit 1
    }
}

# Install binary
function Install-Binary {
    param([string]$TmpFile)

    $installPath = Join-Path $InstallDir $BinaryName

    Write-Info "Installing to $installPath..."

    try {
        # Create install directory if it doesn't exist
        if (-not (Test-Path $InstallDir)) {
            New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
        }

        # Move binary to install location
        Move-Item -Path $TmpFile -Destination $installPath -Force

        Write-Success "Installed to $installPath"
    }
    catch {
        Write-Error-Message "Failed to install binary: $_"
        exit 1
    }
}

# Add to PATH
function Add-ToPath {
    param([string]$Directory)

    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")

    if ($userPath -notlike "*$Directory*") {
        Write-Info "Adding $Directory to PATH..."

        try {
            $newPath = "$userPath;$Directory"
            [Environment]::SetEnvironmentVariable("Path", $newPath, "User")

            Write-Success "Added to PATH"
            Write-Warn "Please restart your terminal for PATH changes to take effect"
        }
        catch {
            Write-Warn "Could not add to PATH automatically"
            Write-Info "Add it manually: Settings > System > About > Advanced system settings > Environment Variables"
        }
    }
    else {
        Write-Success "$Directory is already in PATH"
    }
}

# Create config directory
function Initialize-Config {
    $configDir = Join-Path $env:USERPROFILE ".mkvmender"

    if (-not (Test-Path $configDir)) {
        Write-Info "Creating config directory at $configDir..."
        New-Item -ItemType Directory -Path $configDir -Force | Out-Null
    }
}

# Main installation
function Main {
    Write-Host ""
    Write-Info "MKV Mender Installer for Windows"
    Write-Host ""

    # Check if running with appropriate permissions
    $isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
    if ($isAdmin) {
        Write-Warn "Running as Administrator - installing to user directory anyway"
    }

    # Detect architecture
    $arch = Get-Architecture
    Write-Info "Detected: Windows $arch"

    # Get latest version
    $version = Get-LatestVersion

    # Download binary
    $tmpFile = Download-Binary -Architecture $arch -Version $version

    # Install binary
    Install-Binary -TmpFile $tmpFile

    # Add to PATH
    Add-ToPath -Directory $InstallDir

    # Initialize config
    Initialize-Config

    Write-Host ""
    Write-Success "MKV Mender installed successfully!"
    Write-Host ""
    Write-Info "Get started with:"
    Write-Host "    mkvmender register"
    Write-Host ""
    Write-Info "For more information:"
    Write-Host "    mkvmender --help"
    Write-Host "    https://github.com/$GitHubRepo"
    Write-Host ""
    Write-Warn "Please restart your terminal or PowerShell window to use mkvmender"
    Write-Host ""
}

# Run main installation
try {
    Main
}
catch {
    Write-Error-Message "Installation failed: $_"
    exit 1
}
