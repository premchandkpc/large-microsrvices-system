package com.platform.user.service;

import com.platform.user.dto.UpdateProfileRequest;
import com.platform.user.dto.UserProfileDto;
import com.platform.user.model.UserProfile;
import com.platform.user.repository.UserProfileRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.UUID;

@Service
@RequiredArgsConstructor
public class UserProfileService {

    private final UserProfileRepository profileRepository;

    public UserProfileDto getProfile(UUID userId) {
        UserProfile profile = profileRepository.findByUserId(userId)
            .orElseGet(() -> createDefaultProfile(userId));
        return toDto(profile);
    }

    @Transactional
    public UserProfileDto updateProfile(UUID userId, UpdateProfileRequest request) {
        UserProfile profile = profileRepository.findByUserId(userId)
            .orElseGet(() -> createDefaultProfile(userId));

        if (request.getAvatarUrl() != null) profile.setAvatarUrl(request.getAvatarUrl());
        if (request.getPhone() != null) profile.setPhone(request.getPhone());
        if (request.getBio() != null) profile.setBio(request.getBio());
        if (request.getDepartment() != null) profile.setDepartment(request.getDepartment());
        if (request.getJobTitle() != null) profile.setJobTitle(request.getJobTitle());
        if (request.getTimezone() != null) profile.setTimezone(request.getTimezone());
        if (request.getLocale() != null) profile.setLocale(request.getLocale());
        if (request.getDateOfBirth() != null) profile.setDateOfBirth(request.getDateOfBirth());
        if (request.getAddress() != null) profile.setAddress(request.getAddress());

        profile = profileRepository.save(profile);
        return toDto(profile);
    }

    private UserProfile createDefaultProfile(UUID userId) {
        UserProfile profile = UserProfile.builder()
            .userId(userId)
            .locale("en-US")
            .timezone("UTC")
            .build();
        return profileRepository.save(profile);
    }

    private UserProfileDto toDto(UserProfile profile) {
        return UserProfileDto.builder()
            .id(profile.getId())
            .userId(profile.getUserId())
            .avatarUrl(profile.getAvatarUrl())
            .phone(profile.getPhone())
            .bio(profile.getBio())
            .department(profile.getDepartment())
            .jobTitle(profile.getJobTitle())
            .timezone(profile.getTimezone())
            .locale(profile.getLocale())
            .emailVerified(profile.isEmailVerified())
            .twoFactorEnabled(profile.isTwoFactorEnabled())
            .createdAt(profile.getCreatedAt())
            .build();
    }
}
