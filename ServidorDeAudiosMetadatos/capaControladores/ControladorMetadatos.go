package capacontroladores

import (
	"context"
	"fmt"

	capaFachada "servidoraudios.local/grpc-servidor-audios/capaFachada"
	pb "servidoraudios.local/grpc-servidor-audios/serviciosAudio"
)

type ControladorMetadatos struct {
	pb.UnimplementedMetadataServiceServer
}

func (s *ControladorMetadatos) ListAudioTypes(ctx context.Context, req *pb.Empty) (*pb.AudioTypeList, error) {
	fmt.Printf("[RPC] ListAudioTypes invocado\n")
	return capaFachada.ObtenerTiposAudioDTO(), nil
}

func (s *ControladorMetadatos) ListAudiosByType(ctx context.Context, req *pb.AudioTypeRequest) (*pb.AudioList, error) {
	fmt.Printf("[RPC] ListAudiosByType invocado con type_id=%d\n", req.TypeId)
	return capaFachada.ObtenerAudiosPorTipoDTO(req.TypeId), nil
}

func (s *ControladorMetadatos) GetAudioDetails(ctx context.Context, req *pb.AudioByIdRequest) (*pb.AudioDetails, error) {
	fmt.Printf("[RPC] GetAudioDetails invocado con audio_id=%d\n", req.AudioId)
	return capaFachada.ObtenerDetalleAudioDTO(req.AudioId)
}
