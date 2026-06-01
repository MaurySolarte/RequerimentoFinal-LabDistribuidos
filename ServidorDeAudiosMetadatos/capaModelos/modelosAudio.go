package capamodelos

import "time"

type TipoAudio struct {
	ID     int32
	Nombre string
}

type Metadato struct {
	Titulo   string
	Artista  string
	Genero   string
	Album    string
	Duracion string
}

type Audio struct {
	ID            int32
	TipoID        int32
	TipoNombre    string
	Titulo        string
	ArchivoMP3    string
	FechaRegistro time.Time
	Metadatos     Metadato
}
