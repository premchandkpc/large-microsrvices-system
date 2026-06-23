package com.platform.user.model;

import jakarta.persistence.*;
import lombok.*;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

import java.time.Instant;
import java.util.UUID;

@Entity
@Table(name = "user_profiles")
@Data
@NoArgsConstructor
@AllArgsConstructor
@Builder
@EntityListeners(AuditingEntityListener.class)
public class UserProfile {

    @Id
    private UUID id;

    @Column(nullable = false, unique = true)
    private UUID userId;

    @Column(length = 500)
    private String avatarUrl;

    @Column(length = 50)
    private String phone;

    @Column(columnDefinition = "TEXT")
    private String bio;

    @Column(length = 255)
    private String department;

    @Column(length = 255)
    private String jobTitle;

    @Column(length = 100)
    private String timezone;

    @Column(length = 10)
    private String locale;

    @Column
    private Instant dateOfBirth;

    @Column(length = 500)
    private String address;

    @Builder.Default
    private boolean emailVerified = false;

    @Builder.Default
    private boolean twoFactorEnabled = false;

    @Column(length = 255)
    private String twoFactorSecret;

    @CreatedDate
    private Instant createdAt;

    @LastModifiedDate
    private Instant updatedAt;
}
