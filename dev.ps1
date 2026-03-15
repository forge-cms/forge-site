# dev.ps1 — start, stop, or restart the local forge-site server
# Usage:
#   .\dev.ps1 start          # start on default port 8080
#   .\dev.ps1 start 9090     # start on custom port
#   .\dev.ps1 stop           # kill any process on PORT (default 8080)
#   .\dev.ps1 restart        # stop then start

param(
    [Parameter(Mandatory)][ValidateSet("start","stop","restart")][string]$Action,
    [int]$Port = 8080
)

function Stop-Server {
    param([int]$P)
    $conns = Get-NetTCPConnection -LocalPort $P -ErrorAction SilentlyContinue
    if (-not $conns) {
        Write-Host "Nothing listening on port $P."
        return
    }
    $conns | ForEach-Object {
        Stop-Process -Id $_.OwningProcess -Force -ErrorAction SilentlyContinue
    }
    Write-Host "Stopped process(es) on port $P."
}

switch ($Action) {
    "start" {
        if (-not $env:SECRET) {
            $env:SECRET = "dev-secret"
            Write-Host "SECRET not set - using default: dev-secret"
        }
        Write-Host "Starting forge-site on :$Port ..."
        $env:PORT = "$Port"
        go run .
    }
    "stop" {
        Stop-Server -P $Port
    }
    "restart" {
        Stop-Server -P $Port
        if (-not $env:SECRET) {
            $env:SECRET = "dev-secret"
            Write-Host "SECRET not set - using default: dev-secret"
        }
        Write-Host "Starting forge-site on :$Port ..."
        $env:PORT = "$Port"
        go run .
    }
}
