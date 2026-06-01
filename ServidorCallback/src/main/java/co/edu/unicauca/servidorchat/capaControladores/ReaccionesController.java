package co.edu.unicauca.servidorchat.capaControladores;

import org.springframework.messaging.handler.annotation.MessageMapping;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Controller;
import org.springframework.messaging.simp.stomp.StompHeaderAccessor;

import co.edu.unicauca.servidorchat.capaModelo.CallbackMessage;
import co.edu.unicauca.servidorchat.capaServicios.ReaccionesService;

@Controller
public class ReaccionesController {

  private final ReaccionesService reaccionesService;

  public ReaccionesController(ReaccionesService reaccionesService) {
    this.reaccionesService = reaccionesService;
  }

  @MessageMapping("/enviarReaccion")
  public void enviarReaccion(@Payload CallbackMessage message, StompHeaderAccessor headerAccessor) {
	    reaccionesService.procesarReaccion(message, headerAccessor.getSessionId());
  }
}
