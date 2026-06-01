package co.edu.unicauca.servidorchat.capaServicios;

import org.springframework.stereotype.Service;

import co.edu.unicauca.servidorchat.capaModelo.CallbackMessage;

@Service
public class ReaccionesService {

    private final CallbackPublisherService callbackPublisherService;
    private final CallbackRoutingService callbackRoutingService;

    public ReaccionesService(CallbackPublisherService callbackPublisherService,
            CallbackRoutingService callbackRoutingService) {
        this.callbackPublisherService = callbackPublisherService;
        this.callbackRoutingService = callbackRoutingService;
    }

    public void procesarReaccion(CallbackMessage message) {
        callbackPublisherService.publishJson("/brokerDeReacciones/reaccionesPorCancion", message);
        callbackPublisherService.publishJson(callbackRoutingService.reactionDestination(message), message);
    }
}