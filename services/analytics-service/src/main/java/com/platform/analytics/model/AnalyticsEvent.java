package com.platform.analytics.model;

import jakarta.persistence.*;
import lombok.*;

import java.time.Instant;
import java.util.UUID;

@Entity
@Table(name = "analytics_events", indexes = {
    @Index(name = "idx_event_type", columnList = "eventType"),
    @Index(name = "idx_event_created", columnList = "createdAt"),
    @Index(name = "idx_event_user", columnList = "userId")
})
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class AnalyticsEvent {

    @Id
    private UUID id;

    @Column(nullable = false, length = 100)
    private String eventType;

    private UUID userId;

    private String tenantId;

    @Column(columnDefinition = "JSONB")
    private String properties;

    @Column(length = 45)
    private String ipAddress;

    @Column(length = 500)
    private String userAgent;

    @Column(nullable = false)
    private Instant createdAt;

    @PrePersist
    public void prePersist() {
        if (id == null) id = UUID.randomUUID();
        if (createdAt == null) createdAt = Instant.now();
    }
}
