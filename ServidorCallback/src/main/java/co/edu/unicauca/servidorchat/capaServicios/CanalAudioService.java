package co.edu.unicauca.servidorchat.capaServicios;

import java.util.HashMap;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;

import org.springframework.context.event.EventListener;
import org.springframework.messaging.simp.stomp.StompHeaderAccessor;
import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;
import org.springframework.web.socket.messaging.SessionDisconnectEvent;

import co.edu.unicauca.servidorchat.capaModelo.CallbackMessage;

@Service
public class CanalAudioService {

	private final CallbackPublisherService callbackPublisherService;
	private final CallbackRoutingService callbackRoutingService;
	private final ConcurrentHashMap<String, Set<String>> sesionesPorAudio = new ConcurrentHashMap<>();
	private final ConcurrentHashMap<String, String> audioPorSesion = new ConcurrentHashMap<>();
	private final ConcurrentHashMap<String, String> nicknamePorSesion = new ConcurrentHashMap<>();

	public CanalAudioService(CallbackPublisherService callbackPublisherService,
			CallbackRoutingService callbackRoutingService) {
		this.callbackPublisherService = callbackPublisherService;
		this.callbackRoutingService = callbackRoutingService;
	}

	public void registrarNuevoOyente(CallbackMessage message, String sessionId) {
		String audioId = resolverAudioId(message);
		registrarSesion(sessionId, audioId, resolverNickname(message));
		callbackPublisherService.publishObject(callbackRoutingService.listenersDestination(message), crearPayload("nuevo_oyente", message, "comenzo a reproducir"));
	}

	public void registrarPausa(CallbackMessage message, String sessionId) {
		String audioId = resolverAudioId(message);
		removerSesion(sessionId);
		callbackPublisherService.publishObject(callbackRoutingService.statesDestination(message), crearPayload("pausa", message, "pausó"));
	}

	public void registrarReanudacion(CallbackMessage message, String sessionId) {
		String audioId = resolverAudioId(message);
		registrarSesion(sessionId, audioId, resolverNickname(message));
		callbackPublisherService.publishObject(callbackRoutingService.statesDestination(message), crearPayload("reanuda", message, "reanudó"));
	}

	public void registrarReaccion(CallbackMessage message, String sessionId) {
		if (StringUtils.hasText(sessionId) && StringUtils.hasText(resolverAudioId(message))) {
			registrarSesion(sessionId, resolverAudioId(message), resolverNickname(message));
		}
		HashMap<String, Object> payload = new HashMap<>();
		payload.put("tipo", "reaccion");
		payload.put("nickname", resolverNickname(message));
		payload.put("audioId", resolverAudioId(message));
		payload.put("emoji", resolverEmoji(message));
		callbackPublisherService.publishObject(callbackRoutingService.reactionsDestination(message), payload);
	}

	@EventListener
	public void alDesconectar(SessionDisconnectEvent event) {
		StompHeaderAccessor accessor = StompHeaderAccessor.wrap(event.getMessage());
		String sessionId = accessor.getSessionId();
		removerSesion(sessionId);
	}

	private void registrarSesion(String sessionId, String audioId, String nickname) {
		if (!StringUtils.hasText(sessionId) || !StringUtils.hasText(audioId)) {
			return;
		}
		audioPorSesion.put(sessionId, audioId);
		nicknamePorSesion.put(sessionId, nickname);
		sesionesPorAudio.computeIfAbsent(audioId, clave -> ConcurrentHashMap.newKeySet()).add(sessionId);
	}

	private void removerSesion(String sessionId) {
		if (!StringUtils.hasText(sessionId)) {
			return;
		}
		String audioId = audioPorSesion.remove(sessionId);
		nicknamePorSesion.remove(sessionId);
		if (!StringUtils.hasText(audioId)) {
			return;
		}
		Set<String> sesiones = sesionesPorAudio.get(audioId);
		if (sesiones != null) {
			sesiones.remove(sessionId);
			if (sesiones.isEmpty()) {
				sesionesPorAudio.remove(audioId);
			}
		}
	}

	private HashMap<String, Object> crearPayload(String tipo, CallbackMessage message, String estado) {
		HashMap<String, Object> payload = new HashMap<>();
		payload.put("tipo", tipo);
		payload.put("nickname", resolverNickname(message));
		payload.put("audioId", resolverAudioId(message));
		if (StringUtils.hasText(estado)) {
			payload.put("estado", estado);
		}
		return payload;
	}

	private String resolverAudioId(CallbackMessage message) {
		if (message == null) {
			return "general";
		}
		if (StringUtils.hasText(message.audioId())) {
			return message.audioId().trim();
		}
		if (StringUtils.hasText(message.songKey())) {
			return message.songKey().trim();
		}
		if (StringUtils.hasText(message.songTitle())) {
			return message.songTitle().trim();
		}
		return "general";
	}

	private String resolverNickname(CallbackMessage message) {
		if (message == null || !StringUtils.hasText(message.nickname())) {
			return "Anónimo";
		}
		return message.nickname().trim();
	}

	private String resolverEmoji(CallbackMessage message) {
		if (message == null) {
			return "";
		}
		if (StringUtils.hasText(message.emoji())) {
			return message.emoji().trim();
		}
		if (StringUtils.hasText(message.reaction())) {
			return message.reaction().trim();
		}
		return "";
	}
}