package capaaccesodatos

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var caracteresNoPermitidos = regexp.MustCompile(`[^a-zA-Z0-9\-_]+`)

func GuardarArchivoAudio(titulo string, data []byte) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("el archivo de audio esta vacio")
	}

	if err := os.MkdirAll("audios", os.ModePerm); err != nil {
		return "", fmt.Errorf("no se pudo crear la carpeta de audios: %w", err)
	}

	nombreBase := sanitizarNombreArchivo(titulo)
	nombreArchivo := fmt.Sprintf("%s.mp3", nombreBase)
	rutaArchivo := filepath.Join("audios", nombreArchivo)

	if err := os.WriteFile(rutaArchivo, data, 0644); err != nil {
		return "", fmt.Errorf("no se pudo guardar el archivo de audio: %w", err)
	}

	return rutaArchivo, nil
}

func sanitizarNombreArchivo(valor string) string {
	limpio := strings.TrimSpace(valor)
	if limpio == "" {
		return "audio_sin_titulo"
	}

	normalizado := caracteresNoPermitidos.ReplaceAllString(limpio, "_")
	normalizado = strings.Trim(normalizado, "_")
	if normalizado == "" {
		return "audio_sin_titulo"
	}
	return strings.ToLower(normalizado)
}
