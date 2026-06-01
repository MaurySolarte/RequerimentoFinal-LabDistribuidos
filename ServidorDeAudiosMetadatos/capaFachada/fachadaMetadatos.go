package fachada

import (
	"log"
	"time"

	capaAccesoDatos "servidoraudios.local/grpc-servidor-audios/capaAccesoDatos"
	modelos "servidoraudios.local/grpc-servidor-audios/capaModelos"
	componenteconexioncola "servidoraudios.local/grpc-servidor-audios/componenteConexionCola"
	pb "servidoraudios.local/grpc-servidor-audios/serviciosAudio"
)

var frasesMotivadoras = []string{
	"Cada reproduccion empieza con una buena organizacion.",
	"La musica adecuada puede cambiar el dia completo.",
	"Un nuevo audio, una nueva oportunidad de inspirar.",
}

func ObtenerTiposAudioDTO() *pb.AudioTypeList {
	log.Printf("GetAudioTypes llamado")
	tipos := capaAccesoDatos.ObtenerTiposAudio()
	respuesta := &pb.AudioTypeList{Items: []*pb.AudioType{}}

	for _, tipo := range tipos {
		respuesta.Items = append(respuesta.Items, &pb.AudioType{
			Id:   tipo.ID,
			Name: tipo.Nombre,
		})
	}

	return respuesta
}

func ObtenerAudiosPorTipoDTO(tipoID int32) *pb.AudioList {
	log.Printf("GetAudiosByType llamado con type_id=%d", tipoID)
	audios := capaAccesoDatos.ObtenerAudiosPorTipo(tipoID)
	respuesta := &pb.AudioList{Items: []*pb.AudioSummary{}}

	for _, audio := range audios {
		respuesta.Items = append(respuesta.Items, &pb.AudioSummary{
			Id:     audio.ID,
			TypeId: audio.TipoID,
			Title:  audio.Titulo,
		})
	}

	return respuesta
}

func ObtenerDetalleAudioDTO(audioID int32) (*pb.AudioDetails, error) {
	log.Printf("GetAudioDetails llamado con audio_id=%d", audioID)
	audio, err := capaAccesoDatos.ObtenerAudioPorID(audioID)
	if err != nil {
		return nil, err
	}

	respuesta := &pb.AudioDetails{
		Id:       audio.ID,
		TypeId:   audio.TipoID,
		TypeName: audio.TipoNombre,
		Title:    audio.Titulo,
		Filename: audio.ArchivoMP3,
		Metadata: []*pb.MetadataItem{},
	}

	respuesta.Metadata = append(respuesta.Metadata, &pb.MetadataItem{
		Key:   "Titulo",
		Value: audio.Metadatos.Titulo,
	})
	respuesta.Metadata = append(respuesta.Metadata, &pb.MetadataItem{
		Key:   "Artista",
		Value: audio.Metadatos.Artista,
	})
	respuesta.Metadata = append(respuesta.Metadata, &pb.MetadataItem{
		Key:   "Genero",
		Value: audio.Metadatos.Genero,
	})
	respuesta.Metadata = append(respuesta.Metadata, &pb.MetadataItem{
		Key:   "Album",
		Value: audio.Metadatos.Album,
	})
	respuesta.Metadata = append(respuesta.Metadata, &pb.MetadataItem{
		Key:   "Duracion",
		Value: audio.Metadatos.Duracion,
	})

	return respuesta, nil
}

func RegistrarAudio(tipoID int32, titulo string, data []byte, metadatos modelos.Metadato) (*modelos.Audio, error) {
	rutaArchivo, err := capaAccesoDatos.GuardarArchivoAudio(titulo, data)
	if err != nil {
		return nil, err
	}

	audio, err := capaAccesoDatos.AlmacenarAudio(tipoID, titulo, rutaArchivo, metadatos)
	if err != nil {
		return nil, err
	}

	publisher, err := componenteconexioncola.NewRabbitPublisher()
	if err != nil {
		log.Printf("[WARN] no se pudo inicializar publisher RabbitMQ: %v", err)
		return audio, nil
	}
	defer publisher.Close()

	if err := publisher.PublicarNotificacion(componenteconexioncola.NotificacionRegistroAudio{
		AudioID:         audio.ID,
		Titulo:          audio.Titulo,
		TipoAudio:       audio.TipoNombre,
		RutaArchivo:     audio.ArchivoMP3,
		Artista:         audio.Metadatos.Artista,
		Genero:          audio.Metadatos.Genero,
		Album:           audio.Metadatos.Album,
		Duracion:        audio.Metadatos.Duracion,
		FechaRegistro:   audio.FechaRegistro.Format(time.RFC3339),
		FraseMotivadora: frasesMotivadoras[int(audio.ID)%len(frasesMotivadoras)],
	}); err != nil {
		log.Printf("[WARN] no se pudo publicar notificacion en RabbitMQ: %v", err)
	}

	return audio, nil
}

func ObtenerAudiosModelo() []modelos.Audio {
	return capaAccesoDatos.ObtenerAudios()
}

func ObtenerAudioModelo(audioID int32) (*modelos.Audio, error) {
	return capaAccesoDatos.ObtenerAudioPorID(audioID)
}

func convertirMetadatosAMapa(metadatos modelos.Metadato) map[string]string {
	resultado := make(map[string]string)
	resultado["Titulo"] = metadatos.Titulo
	resultado["Artista"] = metadatos.Artista
	resultado["Genero"] = metadatos.Genero
	resultado["Album"] = metadatos.Album
	resultado["Duracion"] = metadatos.Duracion
	return resultado
}
