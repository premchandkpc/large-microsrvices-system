package com.platform.analytics.model;

import jakarta.persistence.*;
import lombok.*;
import java.time.Instant;
import java.util.UUID;

@Entity
@Table(name = "document_analytics")
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class DocumentAnalytics {

    @Id
    private UUID id;

    private UUID documentId;

    private UUID userId;

    @Column(length = 50)
    private String action;

    private long fileSize;

    @Column(length = 100)
    private String fileType;

    private long processingTimeMs;

    private boolean success;

    private Instant createdAt;

    @PrePersist
    public void prePersist() {
        if (id == null) id = UUID.randomUUID();
        if (createdAt == null) createdAt = Instant.now();
    }
}
