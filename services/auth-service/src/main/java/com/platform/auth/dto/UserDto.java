package com.platform.auth.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;

import java.util.Set;
import java.util.UUID;

@Data
@Builder
@AllArgsConstructor
public class UserDto {
    private UUID id;
    private String email;
    private String name;
    private Set<String> roles;
    private String tenantId;
}
