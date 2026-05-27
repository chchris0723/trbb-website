-- TRBB 鐵人三項運動社團 Database Schema
-- Encoding: UTF8MB4

SET NAMES utf8mb4;
SET time_zone = '+08:00';

-- ═══════════════════════════════════════════════════════════
-- USERS & AUTH
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `users` (
  `id`               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid`             VARCHAR(36)     NOT NULL UNIQUE,
  `username`         VARCHAR(50)     NOT NULL UNIQUE,
  `email`            VARCHAR(100)    NOT NULL UNIQUE,
  `phone`            VARCHAR(20)     DEFAULT NULL,
  `password_hash`    VARCHAR(255)    NOT NULL,
  `display_name`     VARCHAR(100)    DEFAULT NULL,
  `avatar_url`       VARCHAR(500)    DEFAULT NULL,
  `gender`           TINYINT(1)      DEFAULT NULL COMMENT '1=Male,2=Female,3=Other',
  `birthday`         DATE            DEFAULT NULL,
  `emergency_contact`VARCHAR(50)     DEFAULT NULL,
  `emergency_phone`  VARCHAR(20)     DEFAULT NULL,
  `role`             TINYINT(1)      NOT NULL DEFAULT 1 COMMENT '1=Member,2=Coach,9=Admin',
  `status`           TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '0=Pending,1=Active,2=Suspended,3=Rejected',
  `email_verified`   TINYINT(1)      NOT NULL DEFAULT 0,
  `phone_verified`   TINYINT(1)      NOT NULL DEFAULT 0,
  `garmin_connected` TINYINT(1)      NOT NULL DEFAULT 0,
  `garmin_user_id`   VARCHAR(100)    DEFAULT NULL,
  `oauth_provider`   VARCHAR(20)     DEFAULT NULL COMMENT 'google,facebook,line',
  `oauth_id`         VARCHAR(200)    DEFAULT NULL,
  `membership_type`  TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '0=None,1=Basic,2=Premium',
  `membership_start` DATE            DEFAULT NULL,
  `membership_end`   DATE            DEFAULT NULL,
  `created_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`       DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at`       DATETIME        DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_uuid` (`uuid`),
  INDEX `idx_status` (`status`),
  INDEX `idx_role` (`role`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user_tokens` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`      BIGINT UNSIGNED NOT NULL,
  `token`        VARCHAR(255)    NOT NULL UNIQUE,
  `token_type`   VARCHAR(20)     NOT NULL COMMENT 'access,refresh,email_verify,reset_pwd',
  `expires_at`   DATETIME        NOT NULL,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_token` (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `membership_applications` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`      BIGINT UNSIGNED NOT NULL,
  `type`         TINYINT(1)      NOT NULL DEFAULT 1,
  `real_name`    VARCHAR(50)     NOT NULL,
  `id_number`    VARCHAR(20)     DEFAULT NULL,
  `document_url` VARCHAR(500)    DEFAULT NULL,
  `status`       TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '0=Pending,1=Approved,2=Rejected',
  `reviewer_id`  BIGINT UNSIGNED DEFAULT NULL,
  `review_note`  TEXT            DEFAULT NULL,
  `reviewed_at`  DATETIME        DEFAULT NULL,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ═══════════════════════════════════════════════════════════
-- EVENTS & RACE REGISTRATION
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `events` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid`            VARCHAR(36)     NOT NULL UNIQUE,
  `title`           VARCHAR(200)    NOT NULL,
  `description`     LONGTEXT        DEFAULT NULL,
  `event_type`      TINYINT(1)      NOT NULL DEFAULT 1 COMMENT '1=Triathlon,2=Run,3=Swim,4=Bike,5=Training,6=Other',
  `location`        VARCHAR(300)    NOT NULL,
  `location_lat`    DECIMAL(10,8)   DEFAULT NULL,
  `location_lng`    DECIMAL(11,8)   DEFAULT NULL,
  `cover_url`       VARCHAR(500)    DEFAULT NULL,
  `start_at`        DATETIME        NOT NULL,
  `end_at`          DATETIME        NOT NULL,
  `reg_start_at`    DATETIME        NOT NULL,
  `reg_end_at`      DATETIME        NOT NULL,
  `max_participants`INT             DEFAULT NULL,
  `fee`             DECIMAL(10,2)   NOT NULL DEFAULT 0,
  `status`          TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '0=Draft,1=Published,2=Full,3=Ended,4=Cancelled',
  `is_group`        TINYINT(1)      NOT NULL DEFAULT 0,
  `group_min`       INT             DEFAULT NULL,
  `group_max`       INT             DEFAULT NULL,
  `transport_available` TINYINT(1)  NOT NULL DEFAULT 0,
  `creator_id`      BIGINT UNSIGNED NOT NULL,
  `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_uuid` (`uuid`),
  INDEX `idx_start_at` (`start_at`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `event_registrations` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid`         VARCHAR(36)     NOT NULL UNIQUE,
  `event_id`     BIGINT UNSIGNED NOT NULL,
  `user_id`      BIGINT UNSIGNED NOT NULL,
  `group_id`     BIGINT UNSIGNED DEFAULT NULL,
  `category`     VARCHAR(50)     DEFAULT NULL COMMENT 'e.g. Sprint, Olympic, Half, Full',
  `bib_number`   VARCHAR(20)     DEFAULT NULL,
  `status`       TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '0=Pending,1=Confirmed,2=Cancelled,3=Refunded',
  `payment_id`   BIGINT UNSIGNED DEFAULT NULL,
  `check_in_at`  DATETIME        DEFAULT NULL,
  `note`         TEXT            DEFAULT NULL,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_event_user` (`event_id`,`user_id`),
  INDEX `idx_event_id` (`event_id`),
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `group_registrations` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id`     BIGINT UNSIGNED NOT NULL,
  `group_name`   VARCHAR(100)    NOT NULL,
  `captain_id`   BIGINT UNSIGNED NOT NULL,
  `status`       TINYINT(1)      NOT NULL DEFAULT 0,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_event_id` (`event_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ═══════════════════════════════════════════════════════════
-- TRANSPORT
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `transport_options` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `event_id`     BIGINT UNSIGNED NOT NULL,
  `route_name`   VARCHAR(100)    NOT NULL,
  `departure_point` VARCHAR(200) NOT NULL,
  `departure_at` DATETIME        NOT NULL,
  `seats`        INT             NOT NULL DEFAULT 0,
  `fee`          DECIMAL(10,2)   NOT NULL DEFAULT 0,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_event_id` (`event_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `transport_bookings` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `transport_id` BIGINT UNSIGNED NOT NULL,
  `user_id`      BIGINT UNSIGNED NOT NULL,
  `seats`        INT             NOT NULL DEFAULT 1,
  `status`       TINYINT(1)      NOT NULL DEFAULT 0,
  `payment_id`   BIGINT UNSIGNED DEFAULT NULL,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_transport_id` (`transport_id`),
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ═══════════════════════════════════════════════════════════
-- PRODUCTS & SHOP
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `products` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid`         VARCHAR(36)     NOT NULL UNIQUE,
  `title`        VARCHAR(200)    NOT NULL,
  `description`  LONGTEXT        DEFAULT NULL,
  `category`     TINYINT(1)      NOT NULL DEFAULT 1 COMMENT '1=Apparel,2=Equipment,3=Nutrition,4=Accessories',
  `price`        DECIMAL(10,2)   NOT NULL,
  `stock`        INT             NOT NULL DEFAULT 0,
  `images`       JSON            DEFAULT NULL,
  `specs`        JSON            DEFAULT NULL,
  `status`       TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '0=Draft,1=Published,2=SoldOut',
  `creator_id`   BIGINT UNSIGNED NOT NULL,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_category` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `orders` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid`         VARCHAR(36)     NOT NULL UNIQUE,
  `user_id`      BIGINT UNSIGNED NOT NULL,
  `order_type`   TINYINT(1)      NOT NULL DEFAULT 1 COMMENT '1=Shop,2=Event,3=Transport,4=Membership',
  `total_amount` DECIMAL(10,2)   NOT NULL,
  `status`       TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '0=Pending,1=Paid,2=Shipped,3=Completed,4=Cancelled,5=Refunded',
  `shipping_name`VARCHAR(50)     DEFAULT NULL,
  `shipping_phone`VARCHAR(20)    DEFAULT NULL,
  `shipping_addr`VARCHAR(300)    DEFAULT NULL,
  `note`         TEXT            DEFAULT NULL,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `order_items` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `order_id`     BIGINT UNSIGNED NOT NULL,
  `ref_type`     VARCHAR(20)     NOT NULL COMMENT 'product,event,transport',
  `ref_id`       BIGINT UNSIGNED NOT NULL,
  `title`        VARCHAR(200)    NOT NULL,
  `price`        DECIMAL(10,2)   NOT NULL,
  `qty`          INT             NOT NULL DEFAULT 1,
  `spec`         JSON            DEFAULT NULL,
  PRIMARY KEY (`id`),
  INDEX `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ═══════════════════════════════════════════════════════════
-- PAYMENTS
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `payments` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid`            VARCHAR(36)     NOT NULL UNIQUE,
  `order_id`        BIGINT UNSIGNED NOT NULL,
  `user_id`         BIGINT UNSIGNED NOT NULL,
  `amount`          DECIMAL(10,2)   NOT NULL,
  `currency`        VARCHAR(5)      NOT NULL DEFAULT 'TWD',
  `provider`        VARCHAR(20)     NOT NULL COMMENT 'ecpay,stripe,linepay',
  `provider_txn_id` VARCHAR(200)    DEFAULT NULL,
  `status`          TINYINT(1)      NOT NULL DEFAULT 0 COMMENT '0=Pending,1=Success,2=Failed,3=Refunded',
  `payload`         JSON            DEFAULT NULL COMMENT 'raw callback payload',
  `paid_at`         DATETIME        DEFAULT NULL,
  `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_order_id` (`order_id`),
  INDEX `idx_provider_txn` (`provider_txn_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ═══════════════════════════════════════════════════════════
-- SECONDHAND EXCHANGE
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `secondhand_items` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid`         VARCHAR(36)     NOT NULL UNIQUE,
  `seller_id`    BIGINT UNSIGNED NOT NULL,
  `title`        VARCHAR(200)    NOT NULL,
  `description`  TEXT            DEFAULT NULL,
  `category`     TINYINT(1)      NOT NULL DEFAULT 1,
  `condition`    TINYINT(1)      NOT NULL DEFAULT 3 COMMENT '1=New,2=LikeNew,3=Good,4=Fair,5=Poor',
  `price`        DECIMAL(10,2)   NOT NULL,
  `trade_type`   TINYINT(1)      NOT NULL DEFAULT 1 COMMENT '1=Sell,2=Exchange,3=Gift',
  `images`       JSON            DEFAULT NULL,
  `status`       TINYINT(1)      NOT NULL DEFAULT 1 COMMENT '0=Removed,1=Available,2=Reserved,3=Sold',
  `buyer_id`     BIGINT UNSIGNED DEFAULT NULL,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_seller_id` (`seller_id`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ═══════════════════════════════════════════════════════════
-- TRAINING DIARY
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `training_logs` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `uuid`            VARCHAR(36)     NOT NULL UNIQUE,
  `user_id`         BIGINT UNSIGNED NOT NULL,
  `title`           VARCHAR(200)    NOT NULL,
  `sport_type`      TINYINT(1)      NOT NULL DEFAULT 1 COMMENT '1=Run,2=Swim,3=Bike,4=Brick,5=Gym,6=Other',
  `date`            DATE            NOT NULL,
  `duration_min`    INT             DEFAULT NULL COMMENT 'minutes',
  `distance_km`     DECIMAL(8,3)    DEFAULT NULL,
  `avg_heart_rate`  INT             DEFAULT NULL,
  `max_heart_rate`  INT             DEFAULT NULL,
  `calories`        INT             DEFAULT NULL,
  `elevation_m`     INT             DEFAULT NULL,
  `avg_pace`        VARCHAR(10)     DEFAULT NULL COMMENT 'mm:ss/km',
  `avg_speed_kph`   DECIMAL(5,2)    DEFAULT NULL,
  `power_avg`       INT             DEFAULT NULL COMMENT 'watts',
  `garmin_activity_id` VARCHAR(100) DEFAULT NULL,
  `garmin_data`     JSON            DEFAULT NULL,
  `note`            TEXT            DEFAULT NULL,
  `is_public`       TINYINT(1)      NOT NULL DEFAULT 0,
  `photos`          JSON            DEFAULT NULL,
  `created_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`      DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_date` (`date`),
  INDEX `idx_sport_type` (`sport_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `training_plans` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`      BIGINT UNSIGNED NOT NULL,
  `title`        VARCHAR(200)    NOT NULL,
  `description`  TEXT            DEFAULT NULL,
  `goal_event`   VARCHAR(100)    DEFAULT NULL,
  `start_date`   DATE            NOT NULL,
  `end_date`     DATE            NOT NULL,
  `weeks`        INT             DEFAULT NULL,
  `is_public`    TINYINT(1)      NOT NULL DEFAULT 0,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ═══════════════════════════════════════════════════════════
-- NOTIFICATIONS
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `notifications` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`      BIGINT UNSIGNED NOT NULL,
  `type`         VARCHAR(50)     NOT NULL,
  `title`        VARCHAR(200)    NOT NULL,
  `body`         TEXT            DEFAULT NULL,
  `ref_type`     VARCHAR(20)     DEFAULT NULL,
  `ref_id`       BIGINT UNSIGNED DEFAULT NULL,
  `is_read`      TINYINT(1)      NOT NULL DEFAULT 0,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_is_read` (`is_read`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ═══════════════════════════════════════════════════════════
-- GARMIN INTEGRATION
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `garmin_tokens` (
  `id`                  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id`             BIGINT UNSIGNED NOT NULL UNIQUE,
  `oauth_token`         VARCHAR(500)    NOT NULL,
  `oauth_token_secret`  VARCHAR(500)    NOT NULL,
  `garmin_user_id`      VARCHAR(100)    DEFAULT NULL,
  `scope`               VARCHAR(300)    DEFAULT NULL,
  `expires_at`          DATETIME        DEFAULT NULL,
  `created_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`          DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ═══════════════════════════════════════════════════════════
-- ANNOUNCEMENTS
-- ═══════════════════════════════════════════════════════════
CREATE TABLE IF NOT EXISTS `announcements` (
  `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `title`        VARCHAR(200)    NOT NULL,
  `content`      LONGTEXT        NOT NULL,
  `cover_url`    VARCHAR(500)    DEFAULT NULL,
  `is_pinned`    TINYINT(1)      NOT NULL DEFAULT 0,
  `status`       TINYINT(1)      NOT NULL DEFAULT 1,
  `creator_id`   BIGINT UNSIGNED NOT NULL,
  `published_at` DATETIME        DEFAULT NULL,
  `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- ═══════════════════════════════════════════════════════════
-- SEED: Super Admin
-- 預設帳號: admin@trbbtw.com
-- 預設密碼: Trbb@Super2024!  (首次登入後請立即修改)
-- ═══════════════════════════════════════════════════════════
INSERT INTO `users`
  (`uuid`,`username`,`email`,`phone`,`password_hash`,`display_name`,`role`,`status`,`email_verified`,`created_at`,`updated_at`)
VALUES
  (UUID(),'superadmin','admin@trbbtw.com','0000000000',
   '$2b$12$aDcf2hFro7SZ1XRH.s.l4eH.7pD87fyvHJFMVBN/uGwiytW4vbgVS',
   'TRBB 超級管理員',9,1,1,NOW(),NOW())
ON DUPLICATE KEY UPDATE `role`=9, `status`=1, `password_hash`='$2b$12$aDcf2hFro7SZ1XRH.s.l4eH.7pD87fyvHJFMVBN/uGwiytW4vbgVS';
