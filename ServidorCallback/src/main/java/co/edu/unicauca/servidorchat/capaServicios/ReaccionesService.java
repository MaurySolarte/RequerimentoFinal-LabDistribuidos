package co.edu.unicauca.servidorchat.capaServicios;

import org.springframework.stereotype.Service;

import co.edu.unicauca.servidorchat.capaModelo.CallbackMessage;

@Service
public class ReaccionesService {

    private final CallbackPublisherService callbackPublisherService;
    private final CallbackRoutingService callbackRoutingService;
	private final CanalAudioService canalAudioService;

    public ReaccionesService(CallbackPublisherService callbackPublisherService,
	            CallbackRoutingService callbackRoutingService,
				CanalAudioService canalAudioService) {
        this.callbackPublisherService = callbackPublisherService;
        this.callbackRoutingService = callbackRoutingService;
		this.canalAudioService = canalAudioService;
    }

    public void procesarReaccion(CallbackMessage message, String sessionId) {
		canalAudioService.registrarReaccion(message, sessionId);
    }
}