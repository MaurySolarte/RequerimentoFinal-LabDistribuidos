package fachada

import (
	"fmt"
	"log"

	capaAccesoDatos "servidoraudios.local/grpc-servidor-audios/capaAccesoDatos"
	modelos "servidoraudios.local/grpc-servidor-audios/capaModelos"
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

	for clave, valor := range audio.Metadatos.Campos {
		respuesta.Metadata = append(respuesta.Metadata, &pb.MetadataItem{
			Key:   clave,
			Value: valor,
		})
	}

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

	return audio, nil
}

func ObtenerAudiosModelo() []modelos.Audio {
	return capaAccesoDatos.ObtenerAudios()
}

func ObtenerAudioModelo(audioID int32) (*modelos.Audio, error) {
	return capaAccesoDatos.ObtenerAudioPorID(audioID)
}

func ObtenerTiposAudioModelo() []modelos.TipoAudio {
	return capaAccesoDatos.ObtenerTiposAudio()
}

func ObtenerAudioPorTituloModelo(titulo string) (*modelos.Audio, error) {
	return capaAccesoDatos.ObtenerAudioPorTitulo(titulo)
}

func ObtenerAudioPorReferencia(referencia string) (*modelos.Audio, error) {
	if referencia == "" {
		return nil, fmt.Errorf("referencia vacia")
	}

	if audio, err := capaAccesoDatos.ObtenerAudioPorID(parsearEntero(referencia)); err == nil {
		return audio, nil
	}

	return capaAccesoDatos.ObtenerAudioPorTitulo(referencia)
}

func parsearEntero(valor string) int32 {
	var resultado int32
	_, _ = fmt.Sscanf(valor, "%d", &resultado)
	return resultado
}

func convertirMetadatosAMapa(metadatos modelos.Metadato) map[string]string {
	resultado := make(map[string]string)
	for clave, valor := range metadatos.Campos {
		resultado[clave] = valor
	}
	return resultado
}
