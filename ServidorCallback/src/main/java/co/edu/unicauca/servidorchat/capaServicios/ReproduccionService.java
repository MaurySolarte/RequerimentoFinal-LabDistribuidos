package co.edu.unicauca.servidorchat.capaServicios;

import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import co.edu.unicauca.servidorchat.capaModelo.CallbackMessage;

@Service
public class ReproduccionService {

    private final CallbackPublisherService callbackPublisherService;
    private final CallbackRoutingService callbackRoutingService;
    private final CanalAudioService canalAudioService;

    public ReproduccionService(CallbackPublisherService callbackPublisherService,
                CallbackRoutingService callbackRoutingService,
                CanalAudioService canalAudioService) {
        this.callbackPublisherService = callbackPublisherService;
        this.callbackRoutingService = callbackRoutingService;
        this.canalAudioService = canalAudioService;
    }

    public void iniciarReproduccion(CallbackMessage message, String sessionId) {
        if (esEventoNuevoOyente(message)) {
            canalAudioService.registrarNuevoOyente(message, sessionId);
            return;
        }

        canalAudioService.registrarReanudacion(message, sessionId);
    }

    public void detenerReproduccion(CallbackMessage message, String sessionId) {
        canalAudioService.registrarPausa(message, sessionId);
    }

    private boolean esEventoNuevoOyente(CallbackMessage message) {
        if (message == null) {
            return true;
        }
        if (StringUtils.hasText(message.tipo())) {
            return "nuevo_oyente".equalsIgnoreCase(message.tipo());
        }
        return !StringUtils.hasText(message.action()) || "start".equalsIgnoreCase(message.action());
    }
}