package componenteconexioncola

import (
	"encoding/json"
	"fmt"
	"os"

	amqp "github.com/streadway/amqp"
)

type RabbitPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type NotificacionRegistroAudio struct {
	AudioID         int32  `json:"audioId"`
	Titulo          string `json:"titulo"`
	TipoAudio       string `json:"tipoAudio"`
	RutaArchivo     string `json:"rutaArchivo"`
	Artista         string `json:"artista"`
	Genero          string `json:"genero"`
	Album           string `json:"album"`
	Duracion        string `json:"duracion"`
	FechaRegistro   string `json:"fechaRegistro"`
	FraseMotivadora string `json:"fraseMotivadora"`
}

func NewRabbitPublisher() (*RabbitPublisher, error) {
	host := obtenerVariable("RABBITMQ_HOST", "192.168.1.48")
	puerto := obtenerVariable("RABBITMQ_PORT", "5672")
	usuario := obtenerVariable("RABBITMQ_USER", "admin")
	password := obtenerVariable("RABBITMQ_PASSWORD", "1234")
	urlConexion := fmt.Sprintf("amqp://%s:%s@%s:%s/", usuario, password, host, puerto)

	conn, err := amqp.Dial(urlConexion)
	if err != nil {
		return nil, fmt.Errorf("error conectando a RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("error creando canal de RabbitMQ: %w", err)
	}

	queue, err := channel.QueueDeclare(
		"notificaciones_canciones",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = channel.Close()
		_ = conn.Close()
		return nil, fmt.Errorf("error declarando la cola: %w", err)
	}

	return &RabbitPublisher{conn: conn, channel: channel, queue: queue}, nil
}

func (p *RabbitPublisher) PublicarNotificacion(msg NotificacionRegistroAudio) error {
	if p == nil || p.channel == nil {
		return fmt.Errorf("publisher no inicializado")
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error serializando notificacion: %w", err)
	}

	err = p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("error publicando notificacion: %w", err)
	}

	fmt.Printf("[RabbitMQ] Notificacion publicada: %s\n", string(body))
	return nil
}

func (p *RabbitPublisher) Close() {
	if p == nil {
		return
	}
	if p.channel != nil {
		_ = p.channel.Close()
	}
	if p.conn != nil {
		_ = p.conn.Close()
	}
}

func obtenerVariable(clave string, porDefecto string) string {
	if valor := os.Getenv(clave); valor != "" {
		return valor
	}
	return porDefecto
}
