package com.platform.auth.security;

import io.jsonwebtoken.*;
import io.jsonwebtoken.security.Keys;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.util.Date;
import java.util.List;
import java.util.UUID;

@Component
public class JwtTokenProvider {

    private final SecretKey accessKey;
    private final SecretKey refreshKey;
    private final long accessTokenExpiration;
    private final long refreshTokenExpiration;
    private final String issuer;

    public JwtTokenProvider(
            @Value("${jwt.secret}") String secret,
            @Value("${jwt.access-token-expiration}") long accessTokenExpiration,
            @Value("${jwt.refresh-token-expiration}") long refreshTokenExpiration,
            @Value("${jwt.issuer}") String issuer) {
        this.accessKey = Keys.hmacShaKeyFor(secret.getBytes(StandardCharsets.UTF_8));
        this.refreshKey = Keys.hmacShaKeyFor(
            (secret + "-refresh").getBytes(StandardCharsets.UTF_8));
        this.accessTokenExpiration = accessTokenExpiration;
        this.refreshTokenExpiration = refreshTokenExpiration;
        this.issuer = issuer;
    }

    public String generateAccessToken(UUID userId, String email, List<String> roles) {
        Date now = new Date();
        return Jwts.builder()
            .subject(userId.toString())
            .claim("email", email)
            .claim("roles", roles)
            .issuer(issuer)
            .issuedAt(now)
            .expiration(new Date(now.getTime() + accessTokenExpiration))
            .signWith(accessKey)
            .compact();
    }

    public String generateRefreshToken(UUID userId) {
        Date now = new Date();
        return Jwts.builder()
            .subject(userId.toString())
            .claim("type", "refresh")
            .issuer(issuer)
            .issuedAt(now)
            .expiration(new Date(now.getTime() + refreshTokenExpiration))
            .signWith(refreshKey)
            .compact();
    }

    public Claims validateAccessToken(String token) {
        return Jwts.parser()
            .verifyWith(accessKey)
            .requireIssuer(issuer)
            .build()
            .parseSignedClaims(token)
            .getPayload();
    }

    public Claims validateRefreshToken(String token) {
        return Jwts.parser()
            .verifyWith(refreshKey)
            .requireIssuer(issuer)
            .build()
            .parseSignedClaims(token)
            .getPayload();
    }

    public long getAccessTokenExpiration() {
        return accessTokenExpiration;
    }

    public long getRefreshTokenExpiration() {
        return refreshTokenExpiration;
    }
}
