package capafachadaservices

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func StreamAudioFile(tituloCancion string, funcionParaEnviarFragmento func([]byte) error) error {
	log.Printf("Canción solicitada: %s", tituloCancion)
	contenido, err := descargarAudioDesdeREST(tituloCancion)
	if err != nil {
		return err
	}

	buffer := make([]byte, 64*1024) // 64 KB se envian por fragmento
	fragmento := 0
	lector := bytes.NewReader(contenido)

	for {
		n, err := lector.Read(buffer)
		if err == io.EOF {
			log.Println("Canción enviada completamente desde la fachada.")
			break
		}
		if err != nil {
			return fmt.Errorf("error leyendo el archivo: %w", err)
		}

		fragmento++
		log.Printf("Fragmento #%d leído (%d bytes) y enviando", fragmento, n)
		//time.Sleep(1 * time.Second) // Simula un retardo en la lectura
		// ejecutamos la función para enviar el fragmento al cliente
		err = funcionParaEnviarFragmento(buffer[:n])
		if err != nil {
			return fmt.Errorf("error enviando fragmento #%d: %w", fragmento, err)
		}
	}

	return nil
}

func descargarAudioDesdeREST(referencia string) ([]byte, error) {
	ref := strings.TrimSpace(referencia)
	if ref == "" {
		return nil, fmt.Errorf("la referencia del audio esta vacia")
	}

	if _, err := strconv.ParseInt(ref, 10, 32); err != nil {
		ref = url.PathEscape(ref)
	}

	urlDescarga := "http://127.0.0.1:8081/audios/file/" + ref
	respuesta, err := http.Get(urlDescarga)
	if err != nil {
		return nil, fmt.Errorf("no se pudo descargar el audio desde REST: %w", err)
	}
	defer respuesta.Body.Close()

	if respuesta.StatusCode != http.StatusOK {
		cuerpo, _ := io.ReadAll(respuesta.Body)
		return nil, fmt.Errorf("error descargando audio desde REST (%s): %s", respuesta.Status, strings.TrimSpace(string(cuerpo)))
	}

	return io.ReadAll(respuesta.Body)
}
