package dto

type MetadatoEntradaDTO struct {
	Campos map[string]string `json:"campos"`
}

type CrearAudioRespuestaDTO struct {
	ID            int32             `json:"id"`
	TipoID        int32             `json:"tipoId"`
	TipoNombre    string            `json:"tipoNombre"`
	Titulo        string            `json:"titulo"`
	RutaArchivo   string            `json:"rutaArchivo"`
	FechaRegistro string            `json:"fechaRegistro"`
	Metadatos     map[string]string `json:"metadatos"`
}

type TipoAudioDTO struct {
	ID     int32  `json:"id"`
	Nombre string `json:"nombre"`
}

type AudioResumenDTO struct {
	ID         int32             `json:"id"`
	TipoID     int32             `json:"tipoId"`
	TipoNombre string            `json:"tipoNombre"`
	Titulo     string            `json:"titulo"`
	Metadatos  map[string]string `json:"metadatos"`
}

type ErrorRespuestaDTO struct {
	Mensaje string `json:"mensaje"`
}
