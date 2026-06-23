package com.platform.user.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;

import java.time.Instant;
import java.util.UUID;

@Data
@Builder
@AllArgsConstructor
public class UserProfileDto {
    private UUID id;
    private UUID userId;
    private String avatarUrl;
    private String phone;
    private String bio;
    private String department;
    private String jobTitle;
    private String timezone;
    private String locale;
    private boolean emailVerified;
    private boolean twoFactorEnabled;
    private Instant createdAt;
}
