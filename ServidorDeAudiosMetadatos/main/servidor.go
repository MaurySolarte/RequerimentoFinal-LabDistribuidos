package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
	capacontroladores "servidoraudios.local/grpc-servidor-audios/capaControladores"
	pb "servidoraudios.local/grpc-servidor-audios/serviciosAudio"
)

func main() {
	go iniciarServidorGRPC()
	iniciarServidorREST()
}

func iniciarServidorGRPC() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("no se pudo abrir el puerto gRPC: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMetadataServiceServer(grpcServer, &capacontroladores.ControladorMetadatos{})

	fmt.Println("Servidor de metadatos gRPC escuchando en :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("fallo servidor gRPC: %v", err)
	}
}

func iniciarServidorREST() {
	controladorRest := capacontroladores.NewControladorRestMetadatos()

	http.HandleFunc("/audios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controladorRest.RegistrarAudio(w, r)
		case http.MethodGet:
			controladorRest.ListarAudios(w, r)
		default:
			http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/audios/", controladorRest.ObtenerAudioPorID)
	http.HandleFunc("/audios/file/", controladorRest.DescargarArchivoPorID)

	fmt.Println("Servidor de audios/metadatos REST escuchando en :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("fallo servidor REST: %v", err)
	}
}
