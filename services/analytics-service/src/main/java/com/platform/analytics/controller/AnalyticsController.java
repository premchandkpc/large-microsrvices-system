package com.platform.analytics.controller;

import com.platform.analytics.repository.AnalyticsEventRepository;
import com.platform.analytics.repository.DocumentAnalyticsRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.Instant;
import java.util.Map;

@RestController
@RequestMapping("/api/v1/analytics")
@RequiredArgsConstructor
public class AnalyticsController {

    private final AnalyticsEventRepository eventRepo;
    private final DocumentAnalyticsRepository docAnalyticsRepo;

    @GetMapping("/summary")
    public ResponseEntity<Map<String, Object>> getSummary() {
        long totalEvents = eventRepo.count();
        long totalDocuments = docAnalyticsRepo.count();
        return ResponseEntity.ok(Map.of(
            "totalEvents", totalEvents,
            "totalDocuments", totalDocuments,
            "timestamp", Instant.now()
        ));
    }

    @GetMapping("/documents")
    public ResponseEntity<Map<String, Object>> getDocumentAnalytics() {
        long processedCount = docAnalyticsRepo.count();
        return ResponseEntity.ok(Map.of(
            "totalProcessed", processedCount,
            "timestamp", Instant.now()
        ));
    }
}
