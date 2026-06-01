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
	mux := http.NewServeMux()

	mux.HandleFunc("/audios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controladorRest.RegistrarAudio(w, r)
		case http.MethodGet:
			controladorRest.ListarAudios(w, r)
		default:
			http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/tipos-audio", controladorRest.ListarTiposAudio)

	mux.HandleFunc("/audios/", controladorRest.ObtenerAudioPorID)
	mux.HandleFunc("/audios/file/", controladorRest.DescargarArchivoPorID)

	handlerConCORS := permitirCORS(mux)
	fmt.Println("Servidor de audios/metadatos REST escuchando en :8081")
	if err := http.ListenAndServe(":8081", handlerConCORS); err != nil {
		log.Fatalf("fallo servidor REST: %v", err)
	}
}

func permitirCORS(siguiente http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		siguiente.ServeHTTP(w, r)
	})
}
