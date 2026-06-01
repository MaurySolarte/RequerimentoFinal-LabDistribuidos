package capacontroladores

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	fachada "servidoraudios.local/grpc-servidor-audios/capaFachada"
	modelos "servidoraudios.local/grpc-servidor-audios/capaModelos"
	dto "servidoraudios.local/grpc-servidor-audios/dto"
)

type ControladorRestMetadatos struct{}

func NewControladorRestMetadatos() *ControladorRestMetadatos {
	return &ControladorRestMetadatos{}
}

func (c *ControladorRestMetadatos) RegistrarAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	if err := r.ParseMultipartForm(50 << 20); err != nil {
		responderError(w, http.StatusBadRequest, fmt.Sprintf("formato multipart invalido: %v", err))
		return
	}

	titulo := strings.TrimSpace(r.FormValue("titulo"))
	if titulo == "" {
		responderError(w, http.StatusBadRequest, "el titulo es obligatorio")
		return
	}

	tipoID, err := strconv.ParseInt(strings.TrimSpace(r.FormValue("tipo_id")), 10, 32)
	if err != nil {
		responderError(w, http.StatusBadRequest, "tipo_id invalido")
		return
	}

	file, _, err := r.FormFile("archivo")
	if err != nil {
		responderError(w, http.StatusBadRequest, "no se recibio el archivo de audio")
		return
	}
	defer file.Close()

	contenidoArchivo, err := io.ReadAll(file)
	if err != nil {
		responderError(w, http.StatusBadRequest, "no se pudo leer el archivo de audio")
		return
	}

	metadatos, err := parsearMetadatos(r.FormValue("metadatos"))
	if err != nil {
		responderError(w, http.StatusBadRequest, fmt.Sprintf("metadatos invalidos: %v", err))
		return
	}

	audio, err := fachada.RegistrarAudio(int32(tipoID), titulo, contenidoArchivo, metadatos)
	if err != nil {
		responderError(w, http.StatusBadRequest, err.Error())
		return
	}

	respuesta := dto.CrearAudioRespuestaDTO{
		ID:            audio.ID,
		TipoID:        audio.TipoID,
		TipoNombre:    audio.TipoNombre,
		Titulo:        audio.Titulo,
		RutaArchivo:   audio.ArchivoMP3,
		FechaRegistro: audio.FechaRegistro.Format(time.RFC3339),
		Metadatos:     convertirMetadatosDTO(audio.Metadatos),
	}

	responderJSON(w, http.StatusCreated, respuesta)
}

func (c *ControladorRestMetadatos) ListarAudios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	audios := fachada.ObtenerAudiosModelo()
	respuesta := make([]dto.CrearAudioRespuestaDTO, 0, len(audios))
	for _, audio := range audios {
		respuesta = append(respuesta, dto.CrearAudioRespuestaDTO{
			ID:            audio.ID,
			TipoID:        audio.TipoID,
			TipoNombre:    audio.TipoNombre,
			Titulo:        audio.Titulo,
			RutaArchivo:   audio.ArchivoMP3,
			FechaRegistro: audio.FechaRegistro.Format(time.RFC3339),
			Metadatos:     convertirMetadatosDTO(audio.Metadatos),
		})
	}

	responderJSON(w, http.StatusOK, respuesta)
}

func (c *ControladorRestMetadatos) ObtenerAudioPorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	prefijo := "/audios/"
	if !strings.HasPrefix(r.URL.Path, prefijo) {
		responderError(w, http.StatusBadRequest, "ruta invalida")
		return
	}

	idTexto := strings.TrimPrefix(r.URL.Path, prefijo)
	idTexto = strings.TrimSpace(idTexto)
	if idTexto == "" {
		responderError(w, http.StatusBadRequest, "id de audio obligatorio")
		return
	}

	audioID, err := strconv.ParseInt(idTexto, 10, 32)
	if err != nil {
		responderError(w, http.StatusBadRequest, "id de audio invalido")
		return
	}

	audio, err := fachada.ObtenerAudioModelo(int32(audioID))
	if err != nil {
		responderError(w, http.StatusNotFound, err.Error())
		return
	}

	responderJSON(w, http.StatusOK, dto.CrearAudioRespuestaDTO{
		ID:            audio.ID,
		TipoID:        audio.TipoID,
		TipoNombre:    audio.TipoNombre,
		Titulo:        audio.Titulo,
		RutaArchivo:   audio.ArchivoMP3,
		FechaRegistro: audio.FechaRegistro.Format(time.RFC3339),
		Metadatos:     convertirMetadatosDTO(audio.Metadatos),
	})
}

func (c *ControladorRestMetadatos) DescargarArchivoPorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	prefijo := "/audios/file/"
	if !strings.HasPrefix(r.URL.Path, prefijo) {
		responderError(w, http.StatusBadRequest, "ruta invalida")
		return
	}

	idTexto := strings.TrimPrefix(r.URL.Path, prefijo)
	idTexto = strings.TrimSpace(idTexto)
	if idTexto == "" {
		responderError(w, http.StatusBadRequest, "id de audio obligatorio")
		return
	}

	audioID, err := strconv.ParseInt(idTexto, 10, 32)
	if err != nil {
		responderError(w, http.StatusBadRequest, "id de audio invalido")
		return
	}

	audio, err := fachada.ObtenerAudioModelo(int32(audioID))
	if err != nil {
		responderError(w, http.StatusNotFound, err.Error())
		return
	}

	// Abrir archivo y stream hacia el cliente
	file, err := os.Open(audio.ArchivoMP3)
	if err != nil {
		responderError(w, http.StatusNotFound, "archivo no encontrado: "+err.Error())
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "audio/mpeg")
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, file)
}

func parsearMetadatos(raw string) (modelos.Metadato, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return modelos.Metadato{}, nil
	}

	var entrada dto.MetadatoEntradaDTO
	if err := json.Unmarshal([]byte(raw), &entrada); err != nil {
		return modelos.Metadato{}, err
	}

	return modelos.Metadato{
		Titulo:   strings.TrimSpace(entrada.Titulo),
		Artista:  strings.TrimSpace(entrada.Artista),
		Genero:   strings.TrimSpace(entrada.Genero),
		Album:    strings.TrimSpace(entrada.Album),
		Duracion: strings.TrimSpace(entrada.Duracion),
	}, nil
}

func convertirMetadatosDTO(metadato modelos.Metadato) dto.MetadatoEntradaDTO {
	return dto.MetadatoEntradaDTO{
		Titulo:   metadato.Titulo,
		Artista:  metadato.Artista,
		Genero:   metadato.Genero,
		Album:    metadato.Album,
		Duracion: metadato.Duracion,
	}
}

func responderJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func responderError(w http.ResponseWriter, status int, mensaje string) {
	responderJSON(w, status, dto.ErrorRespuestaDTO{Mensaje: mensaje})
}
