package com.platform.auth.service;

import com.platform.auth.dto.*;
import com.platform.auth.model.RefreshToken;
import com.platform.auth.model.User;
import com.platform.auth.repository.RefreshTokenRepository;
import com.platform.auth.repository.UserRepository;
import com.platform.auth.security.JwtTokenProvider;
import io.jsonwebtoken.Claims;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.*;

@Service
@RequiredArgsConstructor
@Slf4j
public class AuthService {

    private final UserRepository userRepository;
    private final RefreshTokenRepository refreshTokenRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtTokenProvider tokenProvider;

    @Value("${app.max-login-attempts}")
    private int maxLoginAttempts;

    @Value("${app.lockout-duration-minutes}")
    private int lockoutDurationMinutes;

    @Transactional
    public AuthResponse login(LoginRequest request) {
        User user = userRepository.findByEmail(request.getEmail())
            .orElseThrow(() -> new BadCredentialsException("Invalid credentials"));

        if (!user.isEnabled()) {
            throw new BadCredentialsException("Account is disabled");
        }

        if (!user.isAccountNonLocked()) {
            if (user.getLockedUntil() != null && user.getLockedUntil().isAfter(Instant.now())) {
                throw new BadCredentialsException("Account is locked. Try again later.");
            }
            user.setAccountNonLocked(true);
            user.setFailedLoginAttempts(0);
        }

        if (!passwordEncoder.matches(request.getPassword(), user.getPasswordHash())) {
            handleFailedLogin(user);
            throw new BadCredentialsException("Invalid credentials");
        }

        userRepository.resetFailedAttempts(user.getId(), Instant.now());

        List<String> roles = new ArrayList<>(user.getRoles());
        String accessToken = tokenProvider.generateAccessToken(
            user.getId(), user.getEmail(), roles);
        String refreshTokenStr = tokenProvider.generateRefreshToken(user.getId());

        saveRefreshToken(user.getId(), refreshTokenStr);

        return AuthResponse.builder()
            .accessToken(accessToken)
            .refreshToken(refreshTokenStr)
            .expiresIn(tokenProvider.getAccessTokenExpiration() / 1000)
            .tokenType("Bearer")
            .user(UserDto.builder()
                .id(user.getId())
                .email(user.getEmail())
                .name(user.getName())
                .roles(user.getRoles())
                .tenantId(user.getTenantId())
                .build())
            .build();
    }

    @Transactional
    public AuthResponse register(RegisterRequest request) {
        if (userRepository.existsByEmail(request.getEmail())) {
            throw new IllegalArgumentException("Email already registered");
        }

        User user = User.builder()
            .email(request.getEmail())
            .passwordHash(passwordEncoder.encode(request.getPassword()))
            .name(request.getName())
            .roles(new HashSet<>(Set.of("ROLE_USER")))
            .build();
        user = userRepository.save(user);

        List<String> roles = new ArrayList<>(user.getRoles());
        String accessToken = tokenProvider.generateAccessToken(
            user.getId(), user.getEmail(), roles);
        String refreshTokenStr = tokenProvider.generateRefreshToken(user.getId());

        saveRefreshToken(user.getId(), refreshTokenStr);

        return AuthResponse.builder()
            .accessToken(accessToken)
            .refreshToken(refreshTokenStr)
            .expiresIn(tokenProvider.getAccessTokenExpiration() / 1000)
            .tokenType("Bearer")
            .user(UserDto.builder()
                .id(user.getId())
                .email(user.getEmail())
                .name(user.getName())
                .roles(user.getRoles())
                .build())
            .build();
    }

    @Transactional
    public AuthResponse refreshToken(RefreshTokenRequest request) {
        Claims claims = tokenProvider.validateRefreshToken(request.getRefreshToken());
        UUID userId = UUID.fromString(claims.getSubject());

        RefreshToken storedToken = refreshTokenRepository
            .findByToken(request.getRefreshToken())
            .orElseThrow(() -> new IllegalArgumentException("Invalid refresh token"));

        if (storedToken.isRevoked()) {
            refreshTokenRepository.deleteByUserId(userId);
            throw new IllegalArgumentException("Refresh token has been revoked");
        }

        if (storedToken.getExpiresAt().isBefore(Instant.now())) {
            refreshTokenRepository.delete(storedToken);
            throw new IllegalArgumentException("Refresh token has expired");
        }

        storedToken.setRevoked(true);
        refreshTokenRepository.save(storedToken);

        User user = userRepository.findById(userId)
            .orElseThrow(() -> new IllegalArgumentException("User not found"));

        List<String> roles = new ArrayList<>(user.getRoles());
        String newAccessToken = tokenProvider.generateAccessToken(
            userId, user.getEmail(), roles);
        String newRefreshToken = tokenProvider.generateRefreshToken(userId);

        saveRefreshToken(userId, newRefreshToken);

        return AuthResponse.builder()
            .accessToken(newAccessToken)
            .refreshToken(newRefreshToken)
            .expiresIn(tokenProvider.getAccessTokenExpiration() / 1000)
            .tokenType("Bearer")
            .user(UserDto.builder()
                .id(user.getId())
                .email(user.getEmail())
                .name(user.getName())
                .roles(user.getRoles())
                .build())
            .build();
    }

    @Transactional
    public void logout(UUID userId, String refreshToken) {
        if (refreshToken != null) {
            refreshTokenRepository.findByToken(refreshToken)
                .ifPresent(token -> {
                    token.setRevoked(true);
                    refreshTokenRepository.save(token);
                });
        }
        refreshTokenRepository.deleteByUserId(userId);
    }

    public ValidateTokenResponse validateToken(String token) {
        Claims claims = tokenProvider.validateAccessToken(token);
        UUID userId = UUID.fromString(claims.getSubject());
        @SuppressWarnings("unchecked")
        List<String> roles = claims.get("roles", List.class);
        String tenantId = claims.get("tenant_id", String.class);

        return ValidateTokenResponse.builder()
            .userId(userId)
            .valid(true)
            .roles(roles != null ? roles : List.of())
            .tenantId(tenantId)
            .build();
    }

    private void saveRefreshToken(UUID userId, String tokenStr) {
        RefreshToken refreshToken = RefreshToken.builder()
            .token(tokenStr)
            .userId(userId)
            .expiresAt(Instant.now()
                .plusMillis(tokenProvider.getRefreshTokenExpiration()))
            .build();
        refreshTokenRepository.save(refreshToken);
    }

    private void handleFailedLogin(User user) {
        int attempts = user.getFailedLoginAttempts() + 1;
        user.setFailedLoginAttempts(attempts);

        if (attempts >= maxLoginAttempts) {
            user.setAccountNonLocked(false);
            user.setLockedUntil(
                Instant.now().plusSeconds(lockoutDurationMinutes * 60));
            log.warn("Account locked due to failed attempts: {}", user.getEmail());
        }

        userRepository.save(user);
    }
}
