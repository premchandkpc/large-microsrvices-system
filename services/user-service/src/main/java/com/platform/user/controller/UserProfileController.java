package com.platform.user.controller;

import com.platform.user.dto.UpdateProfileRequest;
import com.platform.user.dto.UserProfileDto;
import com.platform.user.service.UserProfileService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.security.Principal;
import java.util.UUID;

@RestController
@RequestMapping("/api/v1/users")
@RequiredArgsConstructor
public class UserProfileController {

    private final UserProfileService profileService;

    @GetMapping("/{userId}/profile")
    public ResponseEntity<UserProfileDto> getProfile(@PathVariable UUID userId) {
        return ResponseEntity.ok(profileService.getProfile(userId));
    }

    @PutMapping("/{userId}/profile")
    public ResponseEntity<UserProfileDto> updateProfile(
            @PathVariable UUID userId,
            @Valid @RequestBody UpdateProfileRequest request) {
        return ResponseEntity.ok(profileService.updateProfile(userId, request));
    }

    @GetMapping("/me")
    public ResponseEntity<UserProfileDto> getMyProfile(Principal principal) {
        UUID userId = UUID.fromString(principal.getName());
        return ResponseEntity.ok(profileService.getProfile(userId));
    }

    @PutMapping("/me")
    public ResponseEntity<UserProfileDto> updateMyProfile(
            Principal principal,
            @Valid @RequestBody UpdateProfileRequest request) {
        UUID userId = UUID.fromString(principal.getName());
        return ResponseEntity.ok(profileService.updateProfile(userId, request));
    }
}
