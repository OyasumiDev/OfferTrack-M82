#Requires -Version 5.1
# OfferTrack-M82 — Script de arranque completo
# Uso: .\start.ps1
# PowerShell 5.1 (Windows nativo)

$ErrorActionPreference = "Continue"

# ── Rutas ───────────────────────────────────────────────────────────────────
$ROOT         = "C:\Users\QUERTY\Documents\OfferTrack-M82"
$CORE_DIR     = "$ROOT\core"
$SCRAPER_DIR  = "$ROOT\scraper"
$BIN_DIR      = "$ROOT\bin"
$LOG_DIR      = "$ROOT\logs"
$LOG_FILE     = "$LOG_DIR\start.log"
$SCRAPER_LOG  = "$LOG_DIR\scraper.log"
$PID_FILE     = "$ROOT\.scraper.pid"
$BIN_PATH     = "$BIN_DIR\offertrack.exe"

# ── Crear directorios necesarios ─────────────────────────────────────────────
foreach ($dir in @($LOG_DIR, $BIN_DIR)) {
    if (-not (Test-Path $dir)) { New-Item -ItemType Directory -Path $dir | Out-Null }
}

# ── Helpers de salida ─────────────────────────────────────────────────────────
function Write-Log {
    param([string]$msg)
    "$(Get-Date -Format 'yyyy-MM-dd HH:mm:ss') $msg" | Add-Content -Path $LOG_FILE -Encoding UTF8
}

function OK  { param([string]$m) Write-Host "  [OK]   $m" -ForegroundColor Green;  Write-Log "[OK]   $m" }
function ERR { param([string]$m) Write-Host "  [ERR]  $m" -ForegroundColor Red;    Write-Log "[ERR]  $m" }
function WRN { param([string]$m) Write-Host "  [WARN] $m" -ForegroundColor Yellow; Write-Log "[WARN] $m" }
function INF { param([string]$m) Write-Host "  [INFO] $m" -ForegroundColor Cyan;   Write-Log "[INFO] $m" }

function Section {
    param([string]$title)
    $line = "─" * [Math]::Max(0, 56 - $title.Length)
    Write-Host ""
    Write-Host "  ── $title $line" -ForegroundColor Cyan
    Write-Log "=== $title ==="
}

# ── Helper HTTP (GET y POST con try/catch) ───────────────────────────────────
function Invoke-Http {
    param(
        [string]$Uri,
        [string]$Method  = "GET",
        [string]$Body    = $null,
        [int]   $Timeout = 10
    )
    try {
        $params = @{
            Uri             = $Uri
            Method          = $Method
            UseBasicParsing = $true
            TimeoutSec      = $Timeout
            ErrorAction     = "Stop"
        }
        if ($Body) {
            $params["Body"]        = $Body
            $params["ContentType"] = "application/json"
        }
        return Invoke-WebRequest @params
    }
    catch {
        return $null
    }
}

# ── Limpiar stale PID al inicio ───────────────────────────────────────────────
if (Test-Path $PID_FILE) {
    $savedPid = (Get-Content $PID_FILE -Raw).Trim()
    if ($savedPid -match "^\d+$") {
        if (-not (Get-Process -Id ([int]$savedPid) -ErrorAction SilentlyContinue)) {
            Remove-Item $PID_FILE -ErrorAction SilentlyContinue
        }
    }
}

# ── Banner ─────────────────────────────────────────────────────────────────────
Clear-Host
Write-Host ""
Write-Host "  ══════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "    OfferTrack-M82  —  Arranque del sistema"             -ForegroundColor White
Write-Host "  ══════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host "  $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"             -ForegroundColor DarkGray
Write-Host ""
Write-Log "=== START OfferTrack-M82 ==="

# ═════════════════════════════════════════════════════════════════════════════
# PASO 1 — VERIFICAR DEPENDENCIAS
# ═════════════════════════════════════════════════════════════════════════════
Section "PASO 1: Dependencias"

$missing = $false
foreach ($cmd in @("docker", "node", "go")) {
    if (Get-Command $cmd -ErrorAction SilentlyContinue) {
        OK "$cmd encontrado"
    } else {
        ERR "$cmd no encontrado — instálalo antes de continuar"
        $missing = $true
    }
}

if ($missing) {
    ERR "Dependencias faltantes. Abortando."
    exit 1
}

if (Get-Command "ollama" -ErrorAction SilentlyContinue) {
    OK "ollama encontrado"
} else {
    WRN "ollama no encontrado — continuando sin él (se usará el proveedor de .env)"
}

# ═════════════════════════════════════════════════════════════════════════════
# PASO 2 — DOCKER DESKTOP
# ═════════════════════════════════════════════════════════════════════════════
Section "PASO 2: Docker Desktop"

docker info 2>&1 | Out-Null
if ($LASTEXITCODE -eq 0) {
    OK "Docker Desktop ya está corriendo"
} else {
    INF "Iniciando Docker Desktop..."
    $dockerExe = "C:\Program Files\Docker\Docker\Docker Desktop.exe"
    if (-not (Test-Path $dockerExe)) {
        ERR "No se encontró Docker Desktop en: $dockerExe"
        exit 1
    }
    Start-Process $dockerExe

    $waited = 0
    $ready  = $false
    while ($waited -lt 60) {
        Start-Sleep -Seconds 1
        $waited++
        Write-Host "`r  [....] Esperando Docker... ($waited/60)" -NoNewline -ForegroundColor Yellow
        docker info 2>&1 | Out-Null
        if ($LASTEXITCODE -eq 0) { $ready = $true; break }
    }
    Write-Host ""

    if ($ready) {
        OK "Docker Desktop listo"
    } else {
        ERR "Docker Desktop no respondió en 60 segundos."
        exit 1
    }
}

# ═════════════════════════════════════════════════════════════════════════════
# PASO 3 — QDRANT
# ═════════════════════════════════════════════════════════════════════════════
Section "PASO 3: Qdrant"

# Verificar health primero — si ya responde, no hay nada que hacer
$qdrantLive = Invoke-Http -Uri "http://localhost:6333/healthz" -Timeout 2
if ($qdrantLive -and $qdrantLive.StatusCode -eq 200) {
    OK "Qdrant ya corriendo en localhost:6333"
} else {
    # Inspeccionar contenedor
    $cStatus = docker inspect --format "{{.State.Status}}" qdrant-jobsearch 2>&1

    if ($LASTEXITCODE -ne 0) {
        INF "Creando contenedor qdrant-jobsearch..."
        docker run -d `
            --name qdrant-jobsearch `
            --restart unless-stopped `
            -p 6333:6333 `
            -p 6334:6334 `
            -v qdrant_storage:/qdrant/storage `
            qdrant/qdrant:v1.16.0 2>&1 | Out-Null
        if ($LASTEXITCODE -ne 0) {
            ERR "Error creando contenedor Qdrant."
            exit 1
        }
    } elseif ($cStatus.Trim() -eq "running") {
        INF "Contenedor corriendo pero sin responder — esperando..."
    } else {
        INF "Iniciando contenedor qdrant-jobsearch (estado: $($cStatus.Trim()))..."
        docker start qdrant-jobsearch 2>&1 | Out-Null
        if ($LASTEXITCODE -ne 0) {
            ERR "Error iniciando Qdrant."
            exit 1
        }
    }

    # Esperar respuesta
    $waited = 0
    $ready  = $false
    while ($waited -lt 30) {
        Start-Sleep -Seconds 1
        $waited++
        Write-Host "`r  [....] Esperando Qdrant... ($waited/30)" -NoNewline -ForegroundColor Yellow
        $r = Invoke-Http -Uri "http://localhost:6333/healthz" -Timeout 2
        if ($r -and $r.StatusCode -eq 200) { $ready = $true; break }
    }
    Write-Host ""

    if ($ready) {
        OK "Qdrant respondiendo en localhost:6333"
    } else {
        ERR "Qdrant no respondió en 30 segundos."
        exit 1
    }
}

# ── Leer proveedor y modelo del .env antes de Paso 4 ────────────────────────
$aiProvider = "gemini"
$aiModel    = "gemini-2.0-flash"
if (Test-Path "$ROOT\.env") {
    Get-Content "$ROOT\.env" | ForEach-Object {
        if ($_ -match "^AI_PROVIDER=(.+)$") { $aiProvider = $Matches[1].Trim() }
        if ($_ -match "^AI_MODEL=(.+)$")    { $aiModel    = $Matches[1].Trim() }
    }
}

# ═════════════════════════════════════════════════════════════════════════════
# PASO 4 — OLLAMA
# ═════════════════════════════════════════════════════════════════════════════
Section "PASO 4: Ollama"

$ollamaResp = Invoke-Http -Uri "http://localhost:11434/api/tags" -Timeout 3

if ($ollamaResp -and $ollamaResp.StatusCode -eq 200) {
    OK "Ollama ya corriendo"
} else {
    if (Get-Command "ollama" -ErrorAction SilentlyContinue) {
        INF "Iniciando Ollama..."
        Start-Process "ollama" -ArgumentList "serve" -WindowStyle Hidden

        $waited = 0
        $ready  = $false
        while ($waited -lt 20) {
            Start-Sleep -Seconds 1
            $waited++
            Write-Host "`r  [....] Esperando Ollama... ($waited/20)" -NoNewline -ForegroundColor Yellow
            $r = Invoke-Http -Uri "http://localhost:11434/api/tags" -Timeout 2
            if ($r -and $r.StatusCode -eq 200) { $ollamaResp = $r; $ready = $true; break }
        }
        Write-Host ""

        if ($ready) {
            OK "Ollama iniciado correctamente"
        } else {
            WRN "Ollama no respondió en 20 segundos — continuando sin él"
        }
    } else {
        WRN "Ollama no instalado — continuando sin él"
    }
}

# Verificar el modelo configurado en .env
if ($ollamaResp -and $ollamaResp.StatusCode -eq 200 -and $aiProvider -eq "ollama") {
    try {
        $tags     = $ollamaResp.Content | ConvertFrom-Json
        $hasModel = $tags.models | Where-Object { $_.name -eq $aiModel -or $_.name -like "$aiModel*" }
        if ($hasModel) {
            OK "Modelo '$aiModel' disponible en Ollama"
        } else {
            $available = ($tags.models | ForEach-Object { $_.name }) -join ", "
            WRN "Modelo '$aiModel' no encontrado en Ollama."
            WRN "Para instalarlo: ollama pull $aiModel"
            WRN "Modelos disponibles: $available"
        }
    } catch {
        WRN "No se pudo verificar modelos de Ollama"
    }
}

# ═════════════════════════════════════════════════════════════════════════════
# PASO 5 — COMPILAR BINARIO GO
# ═════════════════════════════════════════════════════════════════════════════
Section "PASO 5: Compilar binario Go"

INF "Ejecutando go build..."
Push-Location $CORE_DIR
$buildOutput = & go build -o "$BIN_PATH" ".\cmd\offertrack\" 2>&1
$buildExit   = $LASTEXITCODE
Pop-Location

if ($buildExit -ne 0) {
    ERR "Error de compilación Go:"
    $buildOutput | ForEach-Object { Write-Host "    $_" -ForegroundColor Red }
    Write-Log "[ERR] go build: $buildOutput"
    exit 1
}
OK "Binario compilado: bin\offertrack.exe"

# ═════════════════════════════════════════════════════════════════════════════
# PASO 6 — SCRAPER NODE.JS
# ═════════════════════════════════════════════════════════════════════════════
Section "PASO 6: Scraper Node.js"

# Matar proceso viejo en :3001 para garantizar que corre el dist actualizado
$oldPid = $null
$netstatLine = netstat -ano 2>$null | Select-String ":3001\s.*LISTENING"
if ($netstatLine) {
    $oldPid = $netstatLine.ToString().Trim().Split()[-1]
    if ($oldPid -match "^\d+$") {
        taskkill /F /PID $oldPid 2>&1 | Out-Null
        INF "Proceso anterior en :3001 (PID $oldPid) detenido — reiniciando con dist actualizado"
        Start-Sleep -Seconds 1
    }
}

INF "Iniciando scraper Node.js..."

# Limpiar log anterior
if (Test-Path $SCRAPER_LOG) { Clear-Content $SCRAPER_LOG -ErrorAction SilentlyContinue }

$cmdArgs = "/c title OfferTrack Scraper && node dist\server.js >> `"$SCRAPER_LOG`" 2>&1"
$scraperProc = Start-Process `
    -FilePath    "cmd.exe" `
    -ArgumentList $cmdArgs `
    -WorkingDirectory $SCRAPER_DIR `
    -WindowStyle Minimized `
    -PassThru

if ($scraperProc) {
    "$($scraperProc.Id)" | Out-File -FilePath $PID_FILE -Encoding ASCII -NoNewline
    INF "Scraper PID $($scraperProc.Id) guardado en .scraper.pid"
}

# Esperar health — hasta 60 seg (carga del modelo de embeddings)
$waited = 0
$ready  = $false
while ($waited -lt 60) {
    Start-Sleep -Seconds 3
    $waited += 3
    Write-Host "`r  [....] Esperando scraper... (modelo cargando, $waited/60 seg)" -NoNewline -ForegroundColor Yellow
    $r = Invoke-Http -Uri "http://localhost:3001/health" -Timeout 3
    if ($r -and $r.StatusCode -eq 200) { $ready = $true; break }
}
Write-Host ""

if ($ready) {
    OK "Scraper listo en localhost:3001"
} else {
    WRN "Scraper no respondió en 60 seg — revisa la ventana minimizada o logs\scraper.log"
}

# ═════════════════════════════════════════════════════════════════════════════
# PASO 7 — PRUEBAS DE INTEGRACIÓN
# ═════════════════════════════════════════════════════════════════════════════
Section "PASO 7: Pruebas de integración"

# Helper: muestra los últimos N chars de un archivo temporal (debug)
function Show-ResponsePreview {
    param([string]$File, [int]$Chars = 300)
    if (Test-Path $File) {
        $content = Get-Content $File -Raw -ErrorAction SilentlyContinue
        if ($content) {
            Write-Host "    >> Respuesta ($($content.Length) bytes): $($content.Substring(0, [Math]::Min($Chars,$content.Length)))" -ForegroundColor DarkGray
        } else {
            Write-Host "    >> Archivo de respuesta vacío" -ForegroundColor DarkGray
        }
    } else {
        Write-Host "    >> No se creó archivo de respuesta" -ForegroundColor DarkGray
    }
}

# ── TEST 1 — Qdrant health ───────────────────────────────────────────────────
Write-Host "  [TEST 1] Qdrant health..." -NoNewline -ForegroundColor Cyan
$r = Invoke-Http -Uri "http://localhost:6333/healthz" -Timeout 5
if ($r -and $r.StatusCode -eq 200) {
    $ver = ""
    try { $ver = " (v$(($r.Content | ConvertFrom-Json).version))" } catch {}
    Write-Host " PASS$ver" -ForegroundColor Green
    Write-Log "[TEST1] PASS$ver"
} else {
    Write-Host " FAIL" -ForegroundColor Red
    Write-Log "[TEST1] FAIL"
    ERR "TEST 1 falló — Qdrant no responde. Abortando."
    exit 1
}

# ── TEST 2 — Embeddings 384 dims ─────────────────────────────────────────────
Write-Host "  [TEST 2] Embeddings 384 dims..." -NoNewline -ForegroundColor Cyan

# Esperar 1.5s para evitar rate-limit si el health check fue reciente
Start-Sleep -Milliseconds 1500

$tmpEmbed = [System.IO.Path]::GetTempFileName()
$tmpBody2 = [System.IO.Path]::GetTempFileName()
'{"texts":["backend developer"]}' | Out-File -FilePath $tmpBody2 -Encoding utf8 -NoNewline

$curlStatus = & curl.exe -s -o $tmpEmbed -w "%{http_code}" `
    -X POST "http://localhost:3001/embed/texts" `
    -H "Content-Type: application/json" `
    --data-binary "@$tmpBody2" `
    --max-time 180 2>$null

Remove-Item $tmpBody2 -ErrorAction SilentlyContinue
$httpCode2 = $curlStatus.Trim()

if ($httpCode2 -eq "200") {
    try {
        $json = Get-Content $tmpEmbed -Raw | ConvertFrom-Json
        if ($json.dimensions -eq 384) {
            Write-Host "  [TEST 2] PASS (384 dims)" -ForegroundColor Green
            Write-Log "[TEST2] PASS — 384 dims"
        } else {
            Write-Host "  [TEST 2] FAIL (dims=$($json.dimensions), esperado 384)" -ForegroundColor Red
            Write-Log "[TEST2] FAIL — dims=$($json.dimensions)"
            ERR "TEST 2 falló — dimensiones incorrectas. Abortando."
            Remove-Item $tmpEmbed -ErrorAction SilentlyContinue
            exit 1
        }
    } catch {
        Write-Host "  [TEST 2] FAIL (no se pudo parsear JSON)" -ForegroundColor Red
        Write-Log "[TEST2] FAIL — parse error: $_"
        ERR "TEST 2 falló — error parseando respuesta. Abortando."
        Remove-Item $tmpEmbed -ErrorAction SilentlyContinue
        exit 1
    }
} else {
    Write-Host "  [TEST 2] FAIL (HTTP $httpCode2)" -ForegroundColor Red
    Write-Log "[TEST2] FAIL — HTTP $httpCode2"
    ERR "TEST 2 falló — /embed/texts devolvió HTTP $httpCode2. Abortando."
    Remove-Item $tmpEmbed -ErrorAction SilentlyContinue
    exit 1
}
Remove-Item $tmpEmbed -ErrorAction SilentlyContinue

# ── TEST 3 — Scraper /scrape responde (sin lanzar Playwright) ────────────────
Write-Host "  [TEST 3] Scraper endpoint /scrape..." -NoNewline -ForegroundColor Cyan
$tmpBody3 = [System.IO.Path]::GetTempFileName()
'{"role":"","portals":["occ"]}' | Out-File -FilePath $tmpBody3 -Encoding utf8 -NoNewline

$tmpScrape3 = [System.IO.Path]::GetTempFileName()
$httpCode3 = & curl.exe -s -o $tmpScrape3 -w "%{http_code}" `
    -X POST "http://localhost:3001/scrape" `
    -H "Content-Type: application/json" `
    --data-binary "@$tmpBody3" `
    --max-time 8 2>$null
Remove-Item $tmpBody3  -ErrorAction SilentlyContinue
Remove-Item $tmpScrape3 -ErrorAction SilentlyContinue

if ($httpCode3 -match "^(200|400)$") {
    Write-Host " PASS (HTTP $httpCode3)" -ForegroundColor Green
    Write-Log "[TEST3] PASS — HTTP $httpCode3"
} else {
    Write-Host " WARN (HTTP $httpCode3 — el scraper tardó o no respondió)" -ForegroundColor Yellow
    Write-Log "[TEST3] WARN — HTTP $httpCode3"
}

# TEST 4 — Binario Go
Write-Host "  [TEST 4] Binario Go (--help)..." -NoNewline -ForegroundColor Cyan
& "$BIN_PATH" --help 2>&1 | Out-Null
if ($LASTEXITCODE -eq 0) {
    Write-Host " PASS" -ForegroundColor Green
    Write-Log "[TEST4] PASS"
} else {
    Write-Host " FAIL (exit $LASTEXITCODE)" -ForegroundColor Red
    Write-Log "[TEST4] FAIL"
    ERR "TEST 4 falló — el binario no responde. Abortando."
    exit 1
}

# ═════════════════════════════════════════════════════════════════════════════
# PASO 8 — BÚSQUEDA AUTOMÁTICA
# ═════════════════════════════════════════════════════════════════════════════
Section "PASO 8: Buscar vacantes"

Write-Host ""
Write-Host "  Sistema listo. Iniciando búsqueda de vacantes." -ForegroundColor White
Write-Host ""
$autoRole = Read-Host "  Puesto a buscar (ej: backend developer — Enter para saltar al menú)"

if ($autoRole) {
    $autoLoc = Read-Host "  Ciudad (ej: monterrey — Enter para cualquiera)"
    Write-Host ""
    INF "Ejecutando búsqueda con re-ranking semántico..."
    Write-Host ""
    if ($autoLoc) {
        & "$BIN_PATH" search --role "$autoRole" --location "$autoLoc"
    } else {
        & "$BIN_PATH" search --role "$autoRole"
    }
    Write-Host ""
    Write-Log "[PASO8] Búsqueda: role=$autoRole loc=$autoLoc"
    Read-Host "  Presiona Enter para continuar al menú"
} else {
    INF "Búsqueda omitida — continuando al menú"
}

# ═════════════════════════════════════════════════════════════════════════════
# PASO 9 — RESUMEN Y MENÚ
# ═════════════════════════════════════════════════════════════════════════════

function Show-Status {
    $q = if ((& curl.exe -s -o NUL -w "%{http_code}" "http://localhost:6333/healthz" --max-time 2 2>$null) -eq "200") { "[OK]" } else { "[--]" }
    $s = if ((& curl.exe -s -o NUL -w "%{http_code}" "http://localhost:3001/health"  --max-time 2 2>$null) -eq "200") { "[OK]" } else { "[--]" }
    $o = if ((& curl.exe -s -o NUL -w "%{http_code}" "http://localhost:11434/api/tags" --max-time 2 2>$null) -eq "200") { "[OK]" } else { "[NA]" }
    $b = if (Test-Path $BIN_PATH) { "[OK]" } else { "[--]" }

    Write-Host ""
    Write-Host "  ══════════════════════════════════════════════════════" -ForegroundColor Cyan
    Write-Host "    OfferTrack-M82  —  Sistema listo" -ForegroundColor White
    Write-Host "  ══════════════════════════════════════════════════════" -ForegroundColor Cyan

    $qColor = if ($q -eq "[OK]") { "Green" } else { "Red" }
    $sColor = if ($s -eq "[OK]") { "Green" } else { "Red" }
    $oColor = if ($o -eq "[OK]") { "Green" } else { "Yellow" }
    $bColor = if ($b -eq "[OK]") { "Green" } else { "Red" }

    Write-Host "    $q Qdrant v1.16.0    -> http://localhost:6333/dashboard" -ForegroundColor $qColor
    Write-Host "    $s Scraper Node.js   -> http://localhost:3001/health"    -ForegroundColor $sColor
    Write-Host "    $o Ollama            -> http://localhost:11434"           -ForegroundColor $oColor
    Write-Host "    $b Binario Go        -> bin\offertrack.exe"              -ForegroundColor $bColor
    Write-Host ""
    Write-Host "    Proveedor IA activo: $aiProvider / $aiModel" -ForegroundColor Cyan
    Write-Host "    (edita .env para cambiar de proveedor)"      -ForegroundColor DarkGray
    Write-Host "  ══════════════════════════════════════════════════════" -ForegroundColor Cyan
    Write-Host ""
    Write-Host "    [1] Buscar vacantes (CLI)"          -ForegroundColor Gray
    Write-Host "    [2] Listar vacantes guardadas"     -ForegroundColor Gray
    Write-Host "    [3] Abrir TUI interactiva"         -ForegroundColor Gray
    Write-Host "    [4] Abrir dashboard Qdrant"        -ForegroundColor Gray
    Write-Host "    [5] Ver logs del scraper"          -ForegroundColor Gray
    Write-Host "    [6] Detener todo y salir"          -ForegroundColor Gray
    Write-Host "    [0] Solo salir (servicios siguen)" -ForegroundColor Gray
    Write-Host ""
}

# Bucle del menú
while ($true) {
    Show-Status
    $choice = Read-Host "  Opcion"

    switch ($choice.Trim()) {
        "1" {
            Write-Host ""
            $role = Read-Host "  Puesto (ej: backend developer)"
            if (-not $role) { WRN "El puesto no puede estar vacío."; continue }
            $loc  = Read-Host "  Ciudad (ej: monterrey — Enter para cualquiera)"
            Write-Host ""
            if ($loc) {
                & "$BIN_PATH" search --role "$role" --location "$loc"
            } else {
                & "$BIN_PATH" search --role "$role"
            }
            Write-Host ""
            Read-Host "  Presiona Enter para volver al menu"
        }
        "2" {
            Write-Host ""
            & "$BIN_PATH" list
            Write-Host ""
            Read-Host "  Presiona Enter para volver al menu"
        }
        "3" {
            Write-Host ""
            INF "Abriendo TUI (q para salir)..."
            Write-Host ""
            & "$BIN_PATH" tui
            Write-Host ""
        }
        "4" {
            Start-Process "http://localhost:6333/dashboard"
            INF "Dashboard abierto en el navegador"
            Start-Sleep -Seconds 1
        }
        "5" {
            Write-Host ""
            if (Test-Path $SCRAPER_LOG) {
                INF "Ultimas 25 lineas de logs\scraper.log:"
                Write-Host ""
                Get-Content $SCRAPER_LOG -Tail 25 | ForEach-Object {
                    Write-Host "    $_" -ForegroundColor Gray
                }
            } else {
                WRN "No se encontró logs\scraper.log aún."
            }
            if (Test-Path $PID_FILE) {
                $scraperPid = (Get-Content $PID_FILE -Raw).Trim()
                INF "PID del proceso scraper: $scraperPid"
            }
            Write-Host ""
            Read-Host "  Presiona Enter para volver al menu"
        }
        "6" {
            Write-Host ""
            INF "Deteniendo servicios..."

            if (Test-Path $PID_FILE) {
                $scraperPid = (Get-Content $PID_FILE -Raw).Trim()
                if ($scraperPid -match "^\d+$") {
                    & taskkill /F /T /PID $scraperPid 2>&1 | Out-Null
                    Remove-Item $PID_FILE -ErrorAction SilentlyContinue
                    OK "Scraper detenido (PID $scraperPid)"
                }
            } else {
                WRN "No se encontró .scraper.pid — el scraper se detuvo manualmente o no inició"
            }

            docker stop qdrant-jobsearch 2>&1 | Out-Null
            if ($LASTEXITCODE -eq 0) { OK "Qdrant detenido" } else { WRN "Qdrant ya estaba detenido" }

            Write-Host ""
            OK "Sistema detenido. Hasta luego."
            Write-Log "=== STOP (menu opcion 6) ==="
            exit 0
        }
        "0" {
            Write-Host ""
            OK "Saliendo. Qdrant y scraper siguen corriendo en segundo plano."
            Write-Log "=== EXIT (menu opcion 0) ==="
            exit 0
        }
        default {
            WRN "Opcion no valida. Elige entre 0 y 6."
            Start-Sleep -Seconds 1
        }
    }
}
