$confirm = Read-Host "Esto eliminara todos los datos de Qdrant. Continuar? (s/N)"
if ($confirm -eq "s" -or $confirm -eq "S") {
    docker compose down -v
    docker compose up qdrant -d
    Write-Host "Qdrant reiniciado." -ForegroundColor Green
} else {
    Write-Host "Cancelado."
}
