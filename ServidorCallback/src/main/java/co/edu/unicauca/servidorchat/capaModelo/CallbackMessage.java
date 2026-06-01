package co.edu.unicauca.servidorchat.capaModelo;

import java.util.Map;

public record CallbackMessage(
        String nickname,
        String songKey,
        String songTitle,
        String songArtist,
        String reaction,
        String action) {

    public static CallbackMessage fromRawPayload(String payload) {
        if (payload == null) {
            return new CallbackMessage(null, null, null, null, null, null);
        }

        String trimmedPayload = payload.trim();
        if (trimmedPayload.startsWith("{")) {
            return fromJson(trimmedPayload);
        }

        return new CallbackMessage(null, null, null, null, trimmedPayload, null);
    }

    private static CallbackMessage fromJson(String payload) {
        try {
            var objectMapper = new com.fasterxml.jackson.databind.ObjectMapper();
            Map<String, String> data = objectMapper.readValue(payload, Map.class);
            return new CallbackMessage(
                    data.get("nickname"),
                    data.get("songKey"),
                    data.get("songTitle"),
                    data.get("songArtist"),
                    data.get("reaction"),
                    data.get("action"));
        } catch (Exception exception) {
            return new CallbackMessage(null, null, null, null, payload, null);
        }
    }
}