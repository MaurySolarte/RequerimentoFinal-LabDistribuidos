package capaaccesodatos

import (
	"fmt"
	"sync"
	"time"

	modelos "servidoraudios.local/grpc-servidor-audios/capaModelos"
)

var (
	repoMu      sync.RWMutex
	audios            = []modelos.Audio{}
	siguienteID int32 = 1
	tiposAudios       = []modelos.TipoAudio{
		{ID: 1, Nombre: "Musica"},
		{ID: 2, Nombre: "Podcasts"},
		{ID: 3, Nombre: "Audiolibros"},
		{ID: 4, Nombre: "Ruido Blanco"},
	}
)

func ObtenerTiposAudio() []modelos.TipoAudio {
	repoMu.RLock()
	defer repoMu.RUnlock()

	tipos := make([]modelos.TipoAudio, len(tiposAudios))
	copy(tipos, tiposAudios)
	return tipos
}

func ObtenerAudios() []modelos.Audio {
	repoMu.RLock()
	defer repoMu.RUnlock()

	resultado := make([]modelos.Audio, len(audios))
	copy(resultado, audios)
	return resultado
}

func ObtenerAudiosPorTipo(tipoID int32) []modelos.Audio {
	repoMu.RLock()
	defer repoMu.RUnlock()

	audiosFiltrados := []modelos.Audio{}
	for _, audio := range audios {
		if audio.TipoID == tipoID {
			audiosFiltrados = append(audiosFiltrados, audio)
		}
	}
	return audiosFiltrados
}

func ObtenerAudioPorID(audioID int32) (*modelos.Audio, error) {
	repoMu.RLock()
	defer repoMu.RUnlock()

	for _, audio := range audios {
		if audio.ID == audioID {
			audioEncontrado := audio
			return &audioEncontrado, nil
		}
	}
	return nil, fmt.Errorf("audio con id %d no encontrado", audioID)
}

func AlmacenarAudio(tipoID int32, titulo string, rutaArchivo string, metadatos modelos.Metadato) (*modelos.Audio, error) {
	repoMu.Lock()
	defer repoMu.Unlock()

	tipoNombre, encontrado := buscarNombreTipo(tipoID)
	if !encontrado {
		return nil, fmt.Errorf("tipo de audio %d no existe", tipoID)
	}

	nuevoAudio := modelos.Audio{
		ID:            siguienteID,
		TipoID:        tipoID,
		TipoNombre:    tipoNombre,
		Titulo:        titulo,
		ArchivoMP3:    rutaArchivo,
		FechaRegistro: time.Now(),
		Metadatos:     metadatos,
	}

	audios = append(audios, nuevoAudio)
	siguienteID++

	copia := nuevoAudio
	return &copia, nil
}

func buscarNombreTipo(tipoID int32) (string, bool) {
	for _, tipo := range tiposAudios {
		if tipo.ID == tipoID {
			return tipo.Nombre, true
		}
	}
	return "", false
}
