package dto

type MetadatoEntradaDTO struct {
	Titulo   string `json:"titulo"`
	Artista  string `json:"artista"`
	Genero   string `json:"genero"`
	Album    string `json:"album"`
	Duracion string `json:"duracion"`
}

type CrearAudioRespuestaDTO struct {
	ID            int32              `json:"id"`
	TipoID        int32              `json:"tipoId"`
	TipoNombre    string             `json:"tipoNombre"`
	Titulo        string             `json:"titulo"`
	RutaArchivo   string             `json:"rutaArchivo"`
	FechaRegistro string             `json:"fechaRegistro"`
	Metadatos     MetadatoEntradaDTO `json:"metadatos"`
}

type ErrorRespuestaDTO struct {
	Mensaje string `json:"mensaje"`
}
