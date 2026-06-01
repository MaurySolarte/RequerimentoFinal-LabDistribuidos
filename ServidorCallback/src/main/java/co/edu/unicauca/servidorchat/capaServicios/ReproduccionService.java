package co.edu.unicauca.servidorchat.capaServicios;

import org.springframework.stereotype.Service;

import co.edu.unicauca.servidorchat.capaModelo.CallbackMessage;

@Service
public class ReproduccionService {

    private final CallbackPublisherService callbackPublisherService;
    private final CallbackRoutingService callbackRoutingService;

    public ReproduccionService(CallbackPublisherService callbackPublisherService,
            CallbackRoutingService callbackRoutingService) {
        this.callbackPublisherService = callbackPublisherService;
        this.callbackRoutingService = callbackRoutingService;
    }

    public void iniciarReproduccion(CallbackMessage message) {
        callbackPublisherService.publishJson(callbackRoutingService.playbackStartDestination(message), message);
    }

    public void detenerReproduccion(CallbackMessage message) {
        callbackPublisherService.publishJson(callbackRoutingService.playbackStopDestination(message), message);
    }
}