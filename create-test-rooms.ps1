# Script para crear habitaciones de prueba
Write-Host "=== CREANDO HABITACIONES DE PRUEBA ===" -ForegroundColor Cyan
Write-Host ""

# Login y obtener token
Write-Host "1. Obteniendo token de administrador..." -ForegroundColor Yellow
$loginBody = @{
    username_or_email = 'admintest'
    password = 'admin123'
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Method POST -Uri http://localhost:8080/login -ContentType 'application/json' -Body $loginBody
    $token = $loginResponse.login.token
    Write-Host "   ✅ Token obtenido exitosamente" -ForegroundColor Green
} catch {
    Write-Host "   ❌ Error al hacer login: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "   Asegúrate de que el usuario 'admintest' existe con password 'admin123'" -ForegroundColor Yellow
    exit
}

Write-Host ""
Write-Host "2. Creando habitaciones..." -ForegroundColor Yellow

# Definir habitaciones
$rooms = @(
    @{ number="101"; type="single"; price=100; description="Single room - Floor 1"; capacity=1; floor=1; has_minibar=$false },
    @{ number="102"; type="single"; price=100; description="Single room - Floor 1"; capacity=1; floor=1; has_minibar=$false },
    @{ number="103"; type="single"; price=110; description="Single room with balcony - Floor 1"; capacity=1; floor=1; has_minibar=$false },
    @{ number="201"; type="double"; price=150; description="Double room - Floor 2"; capacity=2; floor=2; has_minibar=$false },
    @{ number="202"; type="double"; price=150; description="Double room - Floor 2"; capacity=2; floor=2; has_minibar=$false },
    @{ number="203"; type="double"; price=160; description="Double room with ocean view - Floor 2"; capacity=2; floor=2; has_minibar=$true },
    @{ number="301"; type="suite"; price=350; description="Luxury suite with panoramic view - Floor 3"; capacity=4; floor=3; has_minibar=$true },
    @{ number="302"; type="suite"; price=350; description="Luxury suite with jacuzzi - Floor 3"; capacity=4; floor=3; has_minibar=$true },
    @{ number="401"; type="deluxe"; price=500; description="Deluxe suite - Presidential - Floor 4"; capacity=6; floor=4; has_minibar=$true },
    @{ number="402"; type="deluxe"; price=500; description="Deluxe suite - Penthouse - Floor 4"; capacity=6; floor=4; has_minibar=$true }
)

$successCount = 0
$errorCount = 0

foreach ($room in $rooms) {
    $roomBody = @{
        number = $room.number
        type = $room.type
        price = $room.price
        description = $room.description
        capacity = $room.capacity
        floor = $room.floor
        has_wifi = $true
        has_ac = $true
        has_tv = $true
        has_minibar = $room.has_minibar
    } | ConvertTo-Json

    try {
        $response = Invoke-RestMethod -Method POST -Uri http://localhost:8081/api/v1/admin/rooms `
            -Headers @{Authorization="Bearer $token"} `
            -ContentType 'application/json' -Body $roomBody
        
        Write-Host "   ✅ Habitación $($room.number) ($($room.type)) creada" -ForegroundColor Green
        $successCount++
    } catch {
        Write-Host "   ❌ Error creando habitación $($room.number): $($_.Exception.Message)" -ForegroundColor Red
        $errorCount++
    }
}

Write-Host ""
Write-Host "=== RESUMEN ===" -ForegroundColor Cyan
Write-Host "Habitaciones creadas: $successCount" -ForegroundColor Green
Write-Host "Errores: $errorCount" -ForegroundColor $(if($errorCount -eq 0){"Green"}else{"Red"})
Write-Host ""
Write-Host "✨ Proceso completado!" -ForegroundColor Cyan
Write-Host "Puedes ver las habitaciones en: http://localhost:3000/admin/rooms" -ForegroundColor Yellow
