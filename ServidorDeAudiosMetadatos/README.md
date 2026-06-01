# ServidorDeAudiosMetadatos

Servidor en Go que ahora unifica:

- almacenamiento de archivos de audio en disco;
- almacenamiento de metadatos en memoria;
- API REST para registro y consulta de audios;
- publicacion asincrona a RabbitMQ cuando se registra un nuevo audio;
- servicio gRPC de metadatos (compatibilidad temporal con el cliente actual).

## Endpoints REST

- `POST /audios`
  - `multipart/form-data`
  - campos:
    - `titulo` (string, obligatorio)
    - `tipo_id` (int, obligatorio)
    - `archivo` (file, obligatorio)
    - `metadatos` (string JSON opcional), ejemplo:
      - `[{"clave":"artista","valor":"Mon Laferte"}]`
- `GET /audios`
- `GET /audios/{id}`

## Puertos

- REST: `:8080`
- gRPC metadatos: `:50052`

## Variables de entorno RabbitMQ (opcionales)

- `RABBITMQ_HOST` (default `localhost`)
- `RABBITMQ_PORT` (default `5672`)
- `RABBITMQ_USER` (default `admin`)
- `RABBITMQ_PASSWORD` (default `1234`)

## Ejecutar

```bash
go mod tidy
go run ./main
```

## Ejemplo de registro (PowerShell)

```powershell
$metadatos='[{"clave":"artista","valor":"Demo"},{"clave":"genero","valor":"Rock"}]'
Invoke-RestMethod -Method Post -Uri "http://localhost:8080/audios" -Form @{
  titulo = "mi_audio_demo"
  tipo_id = "1"
  metadatos = $metadatos
  archivo = Get-Item "C:\ruta\a\audio.mp3"
}
```
