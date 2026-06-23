package com.platform.auth.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;

import java.util.List;
import java.util.UUID;

@Data
@Builder
@AllArgsConstructor
public class ValidateTokenResponse {
    private boolean valid;
    private UUID userId;
    private List<String> roles;
    private String tenantId;
}
