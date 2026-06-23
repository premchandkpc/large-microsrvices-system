package com.platform.analytics.consumer;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.platform.analytics.model.AnalyticsEvent;
import com.platform.analytics.model.DocumentAnalytics;
import com.platform.analytics.repository.AnalyticsEventRepository;
import com.platform.analytics.repository.DocumentAnalyticsRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import java.util.Map;
import java.util.UUID;

@Component
@RequiredArgsConstructor
@Slf4j
public class EventConsumer {

    private final AnalyticsEventRepository eventRepository;
    private final DocumentAnalyticsRepository docAnalyticsRepository;
    private final ObjectMapper objectMapper;

    @KafkaListener(topics = "analytics-events", groupId = "analytics-service")
    public void consumeAnalyticsEvent(String message) {
        try {
            AnalyticsEvent event = objectMapper.readValue(message, AnalyticsEvent.class);
            eventRepository.save(event);
            log.debug("Analytics event saved: {}", event.getEventType());
        } catch (Exception e) {
            log.error("Failed to process analytics event", e);
        }
    }

    @KafkaListener(topics = "document-analytics", groupId = "analytics-service")
    public void consumeDocumentAnalytics(String message) {
        try {
            Map<String, Object> data = objectMapper.readValue(message, Map.class);
            DocumentAnalytics analytics = DocumentAnalytics.builder()
                .documentId(UUID.fromString((String) data.get("documentId")))
                .userId(data.get("userId") != null
                    ? UUID.fromString((String) data.get("userId")) : null)
                .action((String) data.get("action"))
                .fileSize(data.get("fileSize") != null
                    ? ((Number) data.get("fileSize")).longValue() : 0)
                .fileType((String) data.get("fileType"))
                .processingTimeMs(data.get("processingTimeMs") != null
                    ? ((Number) data.get("processingTimeMs")).longValue() : 0)
                .success(data.get("success") == null || (Boolean) data.get("success"))
                .build();
            docAnalyticsRepository.save(analytics);
            log.debug("Document analytics saved for doc: {}", analytics.getDocumentId());
        } catch (Exception e) {
            log.error("Failed to process document analytics", e);
        }
    }
}
