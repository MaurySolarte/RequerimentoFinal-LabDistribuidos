package co.edu.unicauca.servidorchat.capaControladores;

import org.springframework.messaging.handler.annotation.MessageMapping;
import org.springframework.messaging.handler.annotation.SendTo;
import org.springframework.stereotype.Controller;

@Controller
public class ReaccionesController {
  
  @MessageMapping("/enviarReaccion")
  @SendTo("/brokerDeReacciones/reaccionesPorCancion")
  public String processMessage(String message) {
      return message;
  }
}

