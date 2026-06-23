package com.platform.user.dto;

import jakarta.validation.constraints.Size;
import lombok.Data;

import java.time.Instant;

@Data
public class UpdateProfileRequest {
    @Size(max = 500)
    private String avatarUrl;

    @Size(max = 50)
    private String phone;

    @Size(max = 1000)
    private String bio;

    @Size(max = 255)
    private String department;

    @Size(max = 255)
    private String jobTitle;

    @Size(max = 100)
    private String timezone;

    @Size(max = 10)
    private String locale;

    private Instant dateOfBirth;

    @Size(max = 500)
    private String address;
}
