package co.edu.unicauca.servidorchat.capaServicios;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;

import org.springframework.messaging.simp.SimpMessagingTemplate;
import org.springframework.stereotype.Service;

import co.edu.unicauca.servidorchat.capaModelo.CallbackMessage;

@Service
public class CallbackPublisherService {

    private final ObjectMapper objectMapper;
    private final SimpMessagingTemplate messagingTemplate;

    public CallbackPublisherService(SimpMessagingTemplate messagingTemplate, ObjectMapper objectMapper) {
        this.messagingTemplate = messagingTemplate;
        this.objectMapper = objectMapper;
    }

    @SuppressWarnings("null")
    public void publishJson(String destination, CallbackMessage message) {
        try {
            messagingTemplate.convertAndSend(destination, objectMapper.writeValueAsString(message));
        } catch (JsonProcessingException exception) {
            throw new IllegalStateException("No fue posible serializar el callback para " + destination, exception);
        }
    }

    @SuppressWarnings("null")
    public void publishText(String destination, String payload) {
        messagingTemplate.convertAndSend(destination, payload);
    }
}