$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $MyInvocation.MyCommand.Path
chcp 65001 | Out-Null
[Console]::OutputEncoding = [Text.UTF8Encoding]::new()
$OutputEncoding = [Text.UTF8Encoding]::new()

function Write-Section {
    param([string]$title)
    $line = "-" * 60
    Write-Host $line
    Write-Host $title
    Write-Host $line
}

function Start-GoServer {
    param(
        [string]$dir,
        [string]$file = "main.go"
    )
    Start-Job -WorkingDirectory $dir -ScriptBlock {
        param($file)
        chcp 65001 | Out-Null
        [Console]::OutputEncoding = [Text.UTF8Encoding]::new()
        $OutputEncoding = [Text.UTF8Encoding]::new()
        & go run $file
    } -ArgumentList $file
}

function Stop-GoServer {
    param([System.Management.Automation.Job]$job)
    if ($null -ne $job) {
        Stop-Job $job -ErrorAction SilentlyContinue
        Receive-Job $job -ErrorAction SilentlyContinue | ForEach-Object { "Сервер: $_" } | Out-Host
        Remove-Job $job -ErrorAction SilentlyContinue
        Start-Sleep -Milliseconds 300
    }
}

function Run-InDir {
    param([string]$dir, [scriptblock]$script)
    Push-Location $dir
    try {
        & $script
    } finally {
        Pop-Location
    }
}

function Invoke-Client {
    param([scriptblock]$script)
    & $script 2>&1 | ForEach-Object { "Клиент: $_" } | Out-Host
}

function Read-TcpText {
    param([string]$hostName, [int]$port)
    $client = [System.Net.Sockets.TcpClient]::new()
    $client.Connect($hostName, $port)
    try {
        $stream = $client.GetStream()
        $buffer = New-Object byte[] 4096
        $bytes = $stream.Read($buffer, 0, $buffer.Length)
        if ($bytes -gt 0) {
            [Text.Encoding]::UTF8.GetString($buffer, 0, $bytes)
        }
    } finally {
        $client.Close()
    }
}

Write-Section "Задание 1: TCP эхо-сервер и статистика"
$job = Start-GoServer "$root\\1_tcp\\server"
Start-Sleep -Seconds 1
Run-InDir "$root\\1_tcp\\client" { "hello" | & go run main.go }
$stats = Read-TcpText "localhost" 8082
if ($stats) { Write-Host $stats }
Stop-GoServer $job

Write-Section "Задание 2: UDP эхо-сервер"
$job = Start-GoServer "$root\\2_udp\\server"
Start-Sleep -Seconds 1
Run-InDir "$root\\2_udp\\client" { @("hello", "exit") | & go run main.go }
Stop-GoServer $job

Write-Section "Задание 3: HTTP клиент/сервер"
$job = Start-GoServer "$root\\3_http\\server"
Start-Sleep -Seconds 1
Run-InDir "$root\\3_http\\client" { & go run main.go }
Stop-GoServer $job

Write-Section "Задание 4: REST API (GET/POST/PUT/DELETE)"
$job = Start-GoServer "$root\\4_rest\\server"
Start-Sleep -Seconds 1
Invoke-Client { & curl.exe -s -u admin:secret -X POST -H "Content-Type: application/json" -d '{"name":"Charlie"}' http://localhost:8080/users }
Invoke-Client { & curl.exe -s -u admin:secret http://localhost:8080/users }
Invoke-Client { & curl.exe -s -u admin:secret http://localhost:8080/users/1 }
Invoke-Client { & curl.exe -s -u admin:secret -X PUT -H "Content-Type: application/json" -d '{"name":"Updated"}' http://localhost:8080/users/1 }
Invoke-Client { & curl.exe -s -u admin:secret -X DELETE http://localhost:8080/users/1 }
Invoke-Client { & curl.exe -s -u admin:secret http://localhost:8080/users/1 }
Stop-GoServer $job

Write-Section "Задание 5: Basic Auth роли"
$job = Start-GoServer "$root\\basic_auth"
Start-Sleep -Seconds 1
Invoke-Client { & curl.exe -s -u admin:secret http://localhost:8080/admin }
Invoke-Client { & curl.exe -s -u editor:edit123 http://localhost:8080/admin }
Invoke-Client { & curl.exe -s -u viewer:view123 http://localhost:8080/admin }
Invoke-Client { & curl.exe -s -u bad:creds http://localhost:8080/admin }
Invoke-Client { & curl.exe -s http://localhost:8080/admin }
Stop-GoServer $job

Write-Section "Задание 6: Cookies"
$job = Start-GoServer "$root\\cookies"
Start-Sleep -Seconds 1
Run-InDir "$root\\cookies" {
    Invoke-Client { & curl.exe -s -c cookies.txt http://localhost:8080/ }
    Invoke-Client { & curl.exe -s -b cookies.txt http://localhost:8080/ }
    Remove-Item cookies.txt -ErrorAction SilentlyContinue
}
Stop-GoServer $job

Write-Section "Задание 7: JWT (login/protected)"
$job = Start-GoServer "$root\\jwt" "main_check.go"
Start-Sleep -Seconds 1
$login = & curl.exe -s -X POST -H "Content-Type: application/json" -d '{"username":"user","password":"user123"}' http://localhost:8080/login
$token = ($login | ConvertFrom-Json).token
Write-Host "Token: $token"
Invoke-Client { & curl.exe -s -H "Authorization: Bearer $token" http://localhost:8080/protected }
Invoke-Client { & curl.exe -s http://localhost:8080/protected }
Invoke-Client { & curl.exe -s -H "Authorization: Bearer invalid" http://localhost:8080/protected }
Invoke-Client { & curl.exe -s -X POST -H "Content-Type: application/json" -d '{"username":"user","password":"wrong"}' http://localhost:8080/login }
Invoke-Client { & curl.exe -s -X POST -H "Content-Type: application/json" -d '{"username":"nope","password":"nope"}' http://localhost:8080/login }
Invoke-Client { & curl.exe -s -X POST -H "Content-Type: application/json" -d '{bad json' http://localhost:8080/login }
Stop-GoServer $job

Write-Section "Задание 8: Sessions"
$job = Start-GoServer "$root\\session"
Start-Sleep -Seconds 1
Run-InDir "$root\\session" {
    Invoke-Client { & curl.exe -s -b cookies.txt http://localhost:8080/protected }
    Invoke-Client { & curl.exe -s -c cookies.txt -X POST -H "Content-Type: application/json" -d '{"username":"user","password":"user123"}' http://localhost:8080/login }
    Invoke-Client { & curl.exe -s -b cookies.txt http://localhost:8080/protected }
    Invoke-Client { & curl.exe -s -b cookies.txt http://localhost:8080/logout }
    Invoke-Client { & curl.exe -s -b cookies.txt http://localhost:8080/protected }
    Remove-Item cookies.txt -ErrorAction SilentlyContinue
}
Stop-GoServer $job

Write-Section "Задание 9: Swagger REST API"
$job = Start-GoServer "$root\\swaggo_example"
Start-Sleep -Seconds 1
Invoke-Client { & curl.exe -s http://localhost:8080/users }
Invoke-Client { & curl.exe -s http://localhost:8080/users/1 }
Invoke-Client { & curl.exe -s http://localhost:8080/users/999 }
Invoke-Client { & curl.exe -s -X POST -H "Content-Type: application/json" -d '{"name":"New"}' http://localhost:8080/users }
Invoke-Client { & curl.exe -s -X DELETE http://localhost:8080/users/1 }
Invoke-Client { & curl.exe -s -X DELETE http://localhost:8080/users/999 }
Stop-GoServer $job
