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

	capaAccesoDatos "servidoraudios.local/grpc-servidor-audios/capaAccesoDatos"
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

	tipoFiltro := strings.TrimSpace(r.URL.Query().Get("tipo_id"))
	audios := fachada.ObtenerAudiosModelo()
	if tipoFiltro != "" {
		tipoID, err := strconv.ParseInt(tipoFiltro, 10, 32)
		if err != nil {
			responderError(w, http.StatusBadRequest, "tipo_id invalido")
			return
		}
		filtrados := make([]modelos.Audio, 0)
		for _, audio := range audios {
			if audio.TipoID == int32(tipoID) {
				filtrados = append(filtrados, audio)
			}
		}
		audios = filtrados
	}
	respuesta := make([]dto.AudioResumenDTO, 0, len(audios))
	for _, audio := range audios {
		respuesta = append(respuesta, dto.AudioResumenDTO{
			ID:         audio.ID,
			TipoID:     audio.TipoID,
			TipoNombre: audio.TipoNombre,
			Titulo:     audio.Titulo,
			Metadatos:  convertirMetadatosDTO(audio.Metadatos),
		})
	}

	responderJSON(w, http.StatusOK, respuesta)
}

func (c *ControladorRestMetadatos) ObtenerAudioPorID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	referencia := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/audios/"))
	referencia = strings.Trim(referencia, "/")
	if referencia == "" {
		responderError(w, http.StatusBadRequest, "id o titulo de audio obligatorio")
		return
	}

	audio, err := obtenerAudioPorReferencia(referencia)
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

	referencia := strings.TrimSpace(strings.TrimPrefix(r.URL.Path, "/audios/file/"))
	referencia = strings.Trim(referencia, "/")
	if referencia == "" {
		responderError(w, http.StatusBadRequest, "id o titulo de audio obligatorio")
		return
	}

	audio, err := obtenerAudioPorReferencia(referencia)
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

func (c *ControladorRestMetadatos) ListarTiposAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, http.StatusMethodNotAllowed, "metodo no permitido")
		return
	}

	tipos := fachada.ObtenerTiposAudioModelo()
	respuesta := make([]dto.TipoAudioDTO, 0, len(tipos))
	for _, tipo := range tipos {
		respuesta = append(respuesta, dto.TipoAudioDTO{ID: tipo.ID, Nombre: tipo.Nombre})
	}

	responderJSON(w, http.StatusOK, respuesta)
}

func obtenerAudioPorReferencia(referencia string) (*modelos.Audio, error) {
	if audioID, err := strconv.ParseInt(referencia, 10, 32); err == nil {
		return fachada.ObtenerAudioModelo(int32(audioID))
	}

	return capaAccesoDatos.ObtenerAudioPorTitulo(referencia)
}

func parsearMetadatos(raw string) (modelos.Metadato, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return modelos.Metadato{Campos: map[string]string{}}, nil
	}

	var entrada map[string]string
	if err := json.Unmarshal([]byte(raw), &entrada); err != nil {
		return modelos.Metadato{}, err
	}

	resultado := make(map[string]string, len(entrada))
	for clave, valor := range entrada {
		resultado[strings.TrimSpace(clave)] = strings.TrimSpace(valor)
	}

	return modelos.Metadato{Campos: resultado}, nil
}

func convertirMetadatosDTO(metadato modelos.Metadato) map[string]string {
	resultado := make(map[string]string, len(metadato.Campos))
	for clave, valor := range metadato.Campos {
		resultado[clave] = valor
	}
	return resultado
}

func responderJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func responderError(w http.ResponseWriter, status int, mensaje string) {
	responderJSON(w, status, dto.ErrorRespuestaDTO{Mensaje: mensaje})
}
