package co.edu.unicauca.servidorchat.capaServicios;

import org.springframework.stereotype.Service;
import org.springframework.util.StringUtils;

import co.edu.unicauca.servidorchat.capaModelo.CallbackMessage;

@Service
public class CallbackRoutingService {

    private static final String BASE_DESTINATION = "/brokerDeReacciones/canciones/";

    public String reactionDestination(CallbackMessage message) {
        return BASE_DESTINATION + resolveSongKey(message) + "/reacciones";
    }

    public String playbackStartDestination(CallbackMessage message) {
        return BASE_DESTINATION + resolveSongKey(message) + "/reproducciones/inicio";
    }

    public String playbackStopDestination(CallbackMessage message) {
        return BASE_DESTINATION + resolveSongKey(message) + "/reproducciones/fin";
    }

    public String listenersDestination(CallbackMessage message) {
        return BASE_DESTINATION + resolveAudioId(message) + "/oyentes";
    }

    public String statesDestination(CallbackMessage message) {
        return BASE_DESTINATION + resolveAudioId(message) + "/estados";
    }

    public String reactionsDestination(CallbackMessage message) {
        return BASE_DESTINATION + resolveAudioId(message) + "/reacciones";
    }

    public String resolveSongKey(CallbackMessage message) {
        String rawKey = firstText(message.songKey(), message.songTitle(), "general");
        String normalized = rawKey.toLowerCase().replaceAll("[^a-z0-9]+", "-");
        return normalized.replaceAll("(^-|-$)", "");
    }

    public String resolveAudioId(CallbackMessage message) {
        String rawAudioId = firstText(message.audioId(), message.songKey(), resolveSongKey(message));
        return rawAudioId == null || rawAudioId.isBlank() ? "general" : rawAudioId.trim();
    }

    private String firstText(String firstChoice, String secondChoice, String fallback) {
        if (StringUtils.hasText(firstChoice)) {
            return firstChoice;
        }

        if (StringUtils.hasText(secondChoice)) {
            return secondChoice;
        }

        return fallback;
    }
}