package com.platform.auth.repository;

import com.platform.auth.model.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.time.Instant;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface UserRepository extends JpaRepository<User, UUID> {

    Optional<User> findByEmail(String email);

    boolean existsByEmail(String email);

    @Modifying
    @Query("UPDATE User u SET u.failedLoginAttempts = :attempts WHERE u.id = :userId")
    void updateFailedAttempts(@Param("userId") UUID userId, @Param("attempts") int attempts);

    @Modifying
    @Query("UPDATE User u SET u.accountNonLocked = false, u.lockedUntil = :until WHERE u.id = :userId")
    void lockAccount(@Param("userId") UUID userId, @Param("until") Instant until);

    @Modifying
    @Query("UPDATE User u SET u.failedLoginAttempts = 0, u.lastLoginAt = :now WHERE u.id = :userId")
    void resetFailedAttempts(@Param("userId") UUID userId, @Param("now") Instant now);
}
