package co.edu.unicauca.servidorchat.capaControladores;

import org.springframework.messaging.handler.annotation.MessageMapping;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Controller;

import co.edu.unicauca.servidorchat.capaModelo.CallbackMessage;
import co.edu.unicauca.servidorchat.capaServicios.ReproduccionService;

@Controller
public class ReproduccionController {

    private final ReproduccionService reproduccionService;

    public ReproduccionController(ReproduccionService reproduccionService) {
        this.reproduccionService = reproduccionService;
    }

    @MessageMapping("/iniciarReproduccion")
    public void iniciarReproduccion(@Payload CallbackMessage message) {
        reproduccionService.iniciarReproduccion(message);
    }

    @MessageMapping("/detenerReproduccion")
    public void detenerReproduccion(@Payload CallbackMessage message) {
        reproduccionService.detenerReproduccion(message);
    }
}