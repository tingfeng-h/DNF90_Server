/*M!999999\- enable the sandbox mode */ 

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*M!100616 SET @OLD_NOTE_VERBOSITY=@@NOTE_VERBOSITY, NOTE_VERBOSITY=0 */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `account` varchar(64) NOT NULL,
  `password_hash` varchar(100) NOT NULL,
  `account_id` varchar(64) NOT NULL,
  `created_at` datetime NOT NULL,
  `last_login_at` datetime DEFAULT NULL,
  `last_login_ip` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_account` (`account`),
  UNIQUE KEY `uniq_account_id` (`account_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_auth_failures` (
  `subject` char(64) NOT NULL,
  `failures` int(11) NOT NULL DEFAULT 0,
  `blocked_until` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` bigint(20) NOT NULL,
  PRIMARY KEY (`subject`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_ban_receipts` (
  `request_id` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `payload_hash` char(64) NOT NULL,
  `result_json` text NOT NULL,
  `created_at` datetime(6) NOT NULL,
  PRIMARY KEY (`request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_bans` (
  `kind` varchar(16) NOT NULL,
  `target` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `label` varchar(128) NOT NULL,
  `reason` varchar(512) NOT NULL,
  `expires_at` datetime(6) DEFAULT NULL,
  `revoked` tinyint(1) NOT NULL,
  `revision` bigint(20) unsigned NOT NULL,
  `updated_at` datetime(6) NOT NULL,
  PRIMARY KEY (`kind`,`target`),
  KEY `ban_updated` (`updated_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_benefit_campaigns` (
  `id` char(32) NOT NULL,
  `created_at` datetime(6) NOT NULL,
  `payload` longtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_benefit_jobs` (
  `id` char(32) NOT NULL,
  `created_at` datetime(6) NOT NULL,
  `payload` longtext NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_benefit_policy` (
  `id` int(11) NOT NULL,
  `revision` bigint(20) unsigned NOT NULL,
  `cera` bigint(20) NOT NULL,
  `gold` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_benefit_signup` (
  `account_id` varchar(64) NOT NULL,
  `cera` bigint(20) NOT NULL,
  PRIMARY KEY (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_password_resets` (
  `request_id` varchar(96) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `account_id` varchar(64) NOT NULL,
  `fingerprint` char(64) NOT NULL,
  `created_at` bigint(20) NOT NULL,
  PRIMARY KEY (`request_id`),
  KEY `by_account` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_player_accounts` (
  `account_id` varchar(64) NOT NULL,
  `disabled` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_player_tokens` (
  `token` varchar(64) NOT NULL,
  `device` varchar(64) NOT NULL,
  PRIMARY KEY (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_qq_actions` (
  `id` char(64) NOT NULL,
  `account_id` varchar(64) NOT NULL,
  `qq` varchar(16) NOT NULL,
  `group_id` varchar(16) NOT NULL,
  `state` varchar(16) NOT NULL DEFAULT 'pending',
  `attempts` int(11) NOT NULL DEFAULT 0,
  `cause` varchar(128) NOT NULL DEFAULT '',
  `lease` char(32) NOT NULL DEFAULT '',
  `next_at` bigint(20) NOT NULL DEFAULT 0,
  `created_at` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `due` (`state`,`next_at`),
  KEY `by_account` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_qq_bindings` (
  `account_id` varchar(64) NOT NULL,
  `qq` varchar(16) NOT NULL,
  `group_id` varchar(16) NOT NULL,
  `bound_at` bigint(20) NOT NULL,
  `code_hash` char(64) NOT NULL,
  `state` varchar(16) NOT NULL DEFAULT 'active',
  `revision` bigint(20) unsigned NOT NULL DEFAULT 1,
  `updated_at` bigint(20) NOT NULL DEFAULT 0,
  `checked_at` bigint(20) NOT NULL DEFAULT 0,
  `check_cause` varchar(32) NOT NULL DEFAULT '',
  `next_check_at` bigint(20) NOT NULL DEFAULT 0,
  `check_lease` char(32) NOT NULL DEFAULT '',
  `check_until` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`account_id`),
  UNIQUE KEY `uniq_qq` (`qq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_qq_blacklist` (
  `qq` varchar(16) NOT NULL,
  `account_id` varchar(64) NOT NULL,
  `reason` varchar(512) NOT NULL,
  `created_at` bigint(20) NOT NULL,
  PRIMARY KEY (`qq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_qq_changes` (
  `request_id` varchar(80) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `fingerprint` char(64) NOT NULL,
  `result_json` text NOT NULL,
  `created_at` bigint(20) NOT NULL,
  PRIMARY KEY (`request_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_qq_code_reservations` (
  `code_hash` char(64) NOT NULL,
  `account_id` varchar(64) NOT NULL,
  `expires_at` bigint(20) NOT NULL,
  PRIMARY KEY (`code_hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_qq_codes` (
  `account_id` varchar(64) NOT NULL,
  `code_hash` char(64) NOT NULL,
  `expires_at` bigint(20) NOT NULL,
  `failures` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_qq_settings` (
  `id` int(11) NOT NULL,
  `enabled` tinyint(1) NOT NULL DEFAULT 0,
  `allow_binding` tinyint(1) NOT NULL DEFAULT 0,
  `groups_json` text NOT NULL,
  `revision` bigint(20) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_registration_ips` (
  `ip` varchar(45) NOT NULL,
  `used` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_registration_sources` (
  `account_id` varchar(64) NOT NULL,
  `ip` varchar(45) NOT NULL,
  PRIMARY KEY (`account_id`),
  KEY `ip` (`ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf90_web_accounts_tokens` (
  `token` varchar(64) NOT NULL,
  `account` varchar(64) NOT NULL,
  `account_id` varchar(64) NOT NULL,
  `created_at` datetime NOT NULL,
  `expires_at` datetime NOT NULL,
  `last_used_at` datetime DEFAULT NULL,
  `last_used_ip` varchar(45) DEFAULT NULL,
  `use_count` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`token`),
  KEY `idx_login_token_expires` (`expires_at`),
  KEY `idx_login_token_account_id` (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_account_inventories` (
  `account_id` varchar(128) NOT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_account_inventory_item_extra` (
  `account_id` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `extra_key` varchar(128) NOT NULL,
  `extra_value` longtext NOT NULL,
  PRIMARY KEY (`account_id`,`entry_key`,`extra_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_account_inventory_items` (
  `account_id` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `item_id` bigint(20) NOT NULL DEFAULT 0,
  `item_count` bigint(20) NOT NULL DEFAULT 0,
  `bind_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `expire_at` datetime(6) DEFAULT NULL,
  `raw_entry` longblob DEFAULT NULL,
  PRIMARY KEY (`account_id`,`entry_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_account_metadata` (
  `account_id` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `entry_value` longtext NOT NULL,
  PRIMARY KEY (`account_id`,`entry_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_accounts` (
  `account_id` varchar(128) NOT NULL,
  `state` varchar(32) NOT NULL DEFAULT '',
  `honor_exp` bigint(20) unsigned NOT NULL DEFAULT 0,
  `represent_account_name` varchar(64) DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`account_id`),
  UNIQUE KEY `uk_dnf_accounts_represent_account_name` (`represent_account_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_activity_counters` (
  `account_id` varchar(128) NOT NULL,
  `character_id` varchar(128) NOT NULL,
  `bucket` varchar(10) NOT NULL,
  `metric` varchar(64) NOT NULL,
  `total` bigint(20) NOT NULL,
  PRIMARY KEY (`account_id`,`bucket`,`character_id`,`metric`),
  KEY `idx_activity_ranking` (`metric`,`bucket`,`account_id`,`character_id`,`total`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_activity_events` (
  `account_id` varchar(128) NOT NULL,
  `event_id` varchar(128) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `character_id` varchar(128) NOT NULL,
  `character_name` varchar(64) NOT NULL,
  `metric` varchar(64) NOT NULL,
  `amount` bigint(20) NOT NULL,
  `occurred_at` datetime(6) NOT NULL,
  `service_day` char(10) NOT NULL,
  `dungeon_id` bigint(20) NOT NULL,
  `item_id` bigint(20) unsigned NOT NULL DEFAULT 0,
  `rarity` smallint(6) NOT NULL DEFAULT 0,
  `source` varchar(32) NOT NULL,
  PRIMARY KEY (`account_id`,`event_id`),
  KEY `idx_activity_account_time` (`account_id`,`occurred_at`,`event_id`),
  KEY `idx_activity_character_time` (`account_id`,`character_id`,`occurred_at`,`event_id`),
  KEY `idx_activity_day` (`account_id`,`service_day`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_activity_rank_rewards` (
  `account_id` varchar(128) NOT NULL,
  `board_key` varchar(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `rank_no` int(11) NOT NULL,
  `score` bigint(20) NOT NULL,
  `appearance_code` varchar(80) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `state` varchar(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  PRIMARY KEY (`account_id`,`board_key`),
  KEY `idx_rank_pending` (`account_id`,`state`,`board_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_activity_rank_weeks` (
  `board_key` varchar(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `snapshot_json` longtext NOT NULL,
  PRIMARY KEY (`board_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_character_locations` (
  `character_id` varchar(128) NOT NULL,
  `channel_id` int(11) NOT NULL DEFAULT 0,
  `town_id` bigint(20) NOT NULL DEFAULT 0,
  `dungeon_id` bigint(20) NOT NULL DEFAULT 0,
  `room_id` varchar(128) NOT NULL DEFAULT '',
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_character_roster_equipment` (
  `character_id` varchar(128) NOT NULL,
  `ordinal` int(11) NOT NULL,
  `slot` bigint(20) NOT NULL DEFAULT 0,
  `item_id_or_icon` bigint(20) NOT NULL DEFAULT 0,
  `raw_entry` longblob DEFAULT NULL,
  `packed_flags` bigint(20) NOT NULL DEFAULT 0,
  `optional_id_or_expire` bigint(20) NOT NULL DEFAULT 0,
  `aux_value` bigint(20) NOT NULL DEFAULT 0,
  `aux_flag` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`ordinal`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_character_roster_lists` (
  `character_id` varchar(128) NOT NULL,
  `list_name` varchar(32) NOT NULL,
  `ordinal` int(11) NOT NULL,
  `int_value` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`list_name`,`ordinal`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_character_rosters` (
  `character_id` varchar(128) NOT NULL,
  `header_unk_a` bigint(20) NOT NULL DEFAULT 0,
  `header_unk_b` bigint(20) NOT NULL DEFAULT 0,
  `header_total_or_slot_limit` bigint(20) NOT NULL DEFAULT 0,
  `header_used_or_remain` bigint(20) NOT NULL DEFAULT 0,
  `header_selected_or_page` bigint(20) NOT NULL DEFAULT 0,
  `header_roster_state` bigint(20) NOT NULL DEFAULT 0,
  `header_page_count` bigint(20) NOT NULL DEFAULT 0,
  `header_roster_flag` bigint(20) NOT NULL DEFAULT 0,
  `header_roster_value_a` bigint(20) NOT NULL DEFAULT 0,
  `header_roster_value_b` bigint(20) NOT NULL DEFAULT 0,
  `entry_byte_a` bigint(20) NOT NULL DEFAULT 0,
  `entry_packed_job_grow` bigint(20) NOT NULL DEFAULT 0,
  `entry_byte_c` bigint(20) NOT NULL DEFAULT 0,
  `entry_field_2cc` bigint(20) NOT NULL DEFAULT 0,
  `entry_state0` bigint(20) NOT NULL DEFAULT 0,
  `entry_time_a` bigint(20) NOT NULL DEFAULT 0,
  `entry_time_b` bigint(20) NOT NULL DEFAULT 0,
  `entry_value0` bigint(20) NOT NULL DEFAULT 0,
  `entry_value1` bigint(20) NOT NULL DEFAULT 0,
  `entry_value2` bigint(20) NOT NULL DEFAULT 0,
  `entry_reserved_a` bigint(20) NOT NULL DEFAULT 0,
  `entry_reserved_b` bigint(20) NOT NULL DEFAULT 0,
  `entry_value3` bigint(20) NOT NULL DEFAULT 0,
  `entry_object_id` bigint(20) NOT NULL DEFAULT 0,
  `entry_flag0_eq1` bigint(20) NOT NULL DEFAULT 0,
  `entry_special_status_flag` bigint(20) NOT NULL DEFAULT 0,
  `entry_value5` bigint(20) NOT NULL DEFAULT 0,
  `entry_display_flags` bigint(20) NOT NULL DEFAULT 0,
  `entry_reserved_c` bigint(20) NOT NULL DEFAULT 0,
  `entry_reserved_d` bigint(20) NOT NULL DEFAULT 0,
  `entry_value6` bigint(20) NOT NULL DEFAULT 0,
  `entry_flag1_nonzero` bigint(20) NOT NULL DEFAULT 0,
  `entry_bool_a_eq1` bigint(20) NOT NULL DEFAULT 0,
  `entry_bool_b_eq1` bigint(20) NOT NULL DEFAULT 0,
  `entry_bool_c_eq1` bigint(20) NOT NULL DEFAULT 0,
  `entry_flag2_nonzero` bigint(20) NOT NULL DEFAULT 0,
  `entry_flag3_nonzero` bigint(20) NOT NULL DEFAULT 0,
  `entry_flag4_nonzero` bigint(20) NOT NULL DEFAULT 0,
  `entry_flag5_nonzero` bigint(20) NOT NULL DEFAULT 0,
  `entry_value7` bigint(20) NOT NULL DEFAULT 0,
  `entry_flag6_eq1` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_character_stats` (
  `character_id` varchar(128) NOT NULL,
  `stat_key` varchar(128) NOT NULL,
  `stat_value` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`stat_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_character_stats_luke_backup_20261024_142606` (
  `character_id` varchar(128) NOT NULL,
  `stat_key` varchar(128) NOT NULL,
  `stat_value` bigint(20) NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_characters` (
  `character_id` varchar(128) NOT NULL,
  `account_id` varchar(128) NOT NULL,
  `slot` int(11) NOT NULL DEFAULT 0,
  `name` varchar(64) NOT NULL DEFAULT '',
  `job` varchar(64) NOT NULL DEFAULT '',
  `level` int(11) NOT NULL DEFAULT 0,
  `grow_type` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `exp` bigint(20) NOT NULL DEFAULT 0,
  `ex_equip_slot_stat` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `bonus_sp` int(11) NOT NULL DEFAULT 0,
  `bonus_tp` int(11) NOT NULL DEFAULT 0,
  `pvp_grade` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `pvp_rating_grade` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `user_state` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `tutorial_completed` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `tutorial_reward_progress_38` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `story_digest_last_level` int(10) unsigned NOT NULL DEFAULT 0,
  `story_digest_migration_version` smallint(5) unsigned NOT NULL DEFAULT 0,
  `gold` bigint(20) NOT NULL DEFAULT 0,
  `cera` bigint(20) NOT NULL DEFAULT 0,
  `coin` bigint(20) NOT NULL DEFAULT 0,
  `town_id` smallint(5) unsigned NOT NULL DEFAULT 38,
  `area_id` smallint(5) unsigned NOT NULL DEFAULT 1,
  `pos_x` smallint(6) NOT NULL DEFAULT 450,
  `pos_y` smallint(6) NOT NULL DEFAULT 234,
  `direction` tinyint(3) unsigned NOT NULL DEFAULT 5,
  `area_state` tinyint(3) unsigned NOT NULL DEFAULT 3,
  `delete_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `is_event_character` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `name_tag_item_id` bigint(20) NOT NULL DEFAULT 0,
  `stamina` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `fatigue` smallint(5) unsigned NOT NULL DEFAULT 156,
  `fatigue_limit` smallint(5) unsigned NOT NULL DEFAULT 156,
  `fatigue_penalty` bigint(20) NOT NULL DEFAULT 0,
  `pc_room_id` bigint(20) NOT NULL DEFAULT 65537,
  `is_private_store` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `is_premium_pc_room` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `server_group_id` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `black_count` bigint(20) NOT NULL DEFAULT 0,
  `guild_level` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `chaos_point` bigint(20) NOT NULL DEFAULT 0,
  `disguise_kind` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `is_disguised` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `expert_job_type` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `expert_job_exp` bigint(20) NOT NULL DEFAULT 0,
  `extra46` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `extra47` bigint(20) NOT NULL DEFAULT 0,
  `extra51` smallint(5) unsigned NOT NULL DEFAULT 0,
  `is_hardcore_mode` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `is_hardcore_dead` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `hardcore_death_count` smallint(5) unsigned NOT NULL DEFAULT 0,
  `user_state_bits` tinyint(3) unsigned NOT NULL DEFAULT 3,
  `chat_ban_end_time` bigint(20) NOT NULL DEFAULT 0,
  `fatigue_update` smallint(5) unsigned NOT NULL DEFAULT 0,
  `return_user_flag` tinyint(3) unsigned NOT NULL DEFAULT 1,
  `channel_display_mode` smallint(5) unsigned NOT NULL DEFAULT 0,
  `channel_type` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `channel_id` smallint(5) unsigned NOT NULL DEFAULT 2,
  `is_return_user` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `link_slot_enabled` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `link_type_a` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `link_type_b` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `emotion_index` smallint(5) unsigned NOT NULL DEFAULT 0,
  `action_byte` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `fatigue_display_update` smallint(5) unsigned NOT NULL DEFAULT 0,
  `costume_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `aura_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `premium_crystal_selection` tinyint(4) NOT NULL DEFAULT -1,
  `pet_display_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `title_display_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `pvp_stat_a` bigint(20) NOT NULL DEFAULT 0,
  `pvp_win_streak` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `pvp_lose_streak` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `pvp_rank_point` bigint(20) NOT NULL DEFAULT 0,
  `trailing_byte` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `stat_block_marker` bigint(20) NOT NULL DEFAULT 83,
  `stat_hp_max` bigint(20) NOT NULL DEFAULT 0,
  `stat_mp_max` bigint(20) NOT NULL DEFAULT 0,
  `stat_strength` int(11) NOT NULL DEFAULT 0,
  `stat_intelligence` int(11) NOT NULL DEFAULT 0,
  `stat_vitality` int(11) NOT NULL DEFAULT 0,
  `stat_spirit` int(11) NOT NULL DEFAULT 0,
  `stat_physical_attack` int(11) NOT NULL DEFAULT 0,
  `stat_physical_defense` int(11) NOT NULL DEFAULT 0,
  `stat_magical_attack` int(11) NOT NULL DEFAULT 0,
  `stat_magical_defense` int(11) NOT NULL DEFAULT 0,
  `stat_independent_attack` int(11) NOT NULL DEFAULT 0,
  `stat_fire_resistance` int(11) NOT NULL DEFAULT 0,
  `stat_water_resistance` int(11) NOT NULL DEFAULT 0,
  `stat_dark_resistance` int(11) NOT NULL DEFAULT 0,
  `stat_light_resistance` int(11) NOT NULL DEFAULT 0,
  `active_status_resistance_00` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_01` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_02` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_03` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_04` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_05` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_06` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_07` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_08` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_09` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_10` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_11` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_12` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_13` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_14` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_15` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_16` smallint(5) unsigned NOT NULL DEFAULT 0,
  `active_status_resistance_17` smallint(5) unsigned NOT NULL DEFAULT 0,
  `stat_inventory_limit` bigint(20) NOT NULL DEFAULT 0,
  `stat_hp_regen_speed` smallint(5) unsigned NOT NULL DEFAULT 0,
  `stat_mp_regen_speed` smallint(5) unsigned NOT NULL DEFAULT 0,
  `stat_move_speed` bigint(20) NOT NULL DEFAULT 0,
  `stat_attack_speed` smallint(5) unsigned NOT NULL DEFAULT 0,
  `stat_cast_speed` smallint(5) unsigned NOT NULL DEFAULT 0,
  `stat_hit_recovery` smallint(5) unsigned NOT NULL DEFAULT 0,
  `stat_jump_power` smallint(5) unsigned NOT NULL DEFAULT 0,
  `stat_weight` bigint(20) NOT NULL DEFAULT 0,
  `stat_level` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `name_tag_expire_time` bigint(20) NOT NULL DEFAULT 0,
  `skill_tree_index` smallint(6) NOT NULL DEFAULT -1,
  `equipped_creature_level` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `equip_list_trailing` bigint(20) NOT NULL DEFAULT 0,
  `manage_level` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `flag_byte` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `guild_power_war` bigint(20) NOT NULL DEFAULT 0,
  `server_timestamp` bigint(20) NOT NULL DEFAULT 0,
  `quest_shop_count` smallint(5) unsigned NOT NULL DEFAULT 0,
  `progress1` bigint(20) NOT NULL DEFAULT 0,
  `progress2` bigint(20) NOT NULL DEFAULT 0,
  `create_option_len` smallint(5) unsigned NOT NULL DEFAULT 0,
  `roster_state0` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_time_a` bigint(20) NOT NULL DEFAULT 0,
  `roster_time_b` bigint(20) NOT NULL DEFAULT 0,
  `roster_value0` bigint(20) NOT NULL DEFAULT 0,
  `roster_value1` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_value2` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_reserved_a` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_reserved_b` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_linked_id_00` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_linked_id_01` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_linked_id_02` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_linked_id_03` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_linked_id_04` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_linked_id_05` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_linked_id_06` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_linked_id_07` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_value3` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_object_id` bigint(20) NOT NULL DEFAULT 0,
  `roster_flag0_eq1` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_card_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_value5` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_display_flags` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_00` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_01` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_02` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_03` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_04` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_05` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_06` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_07` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_08` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_09` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_10` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_tail_11` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `roster_flag6_eq1` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_00` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_01` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_02` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_03` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_04` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_05` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_06` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_07` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_08` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_09` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_10` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_11` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_12` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_13` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_14` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_15` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_16` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_17` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_18` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_19` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_20` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_21` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_22` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_23` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_24` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_25` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_26` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_27` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_28` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_29` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_30` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_31` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_32` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_33` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_34` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_35` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_36` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_37` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_38` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_39` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_40` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_41` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_42` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_43` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_44` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_45` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_46` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_47` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_48` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_49` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_50` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_51` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_52` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_53` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_54` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_55` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_56` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_57` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_58` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_59` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_60` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_61` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_62` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `create_option_byte_63` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`),
  UNIQUE KEY `uk_dnf_characters_account_slot_active` (`account_id`,`slot`,`delete_flag`),
  KEY `idx_dnf_characters_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_clear_history` (
  `character_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `completion_hash` binary(32) NOT NULL,
  `account_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  `character_name` varchar(128) NOT NULL,
  `dungeon_id` bigint(20) NOT NULL,
  `dungeon_name` varchar(255) NOT NULL,
  `difficulty` int(11) NOT NULL,
  `started_at` datetime(6) NOT NULL,
  `finished_at` datetime(6) NOT NULL,
  `duration_ms` bigint(20) NOT NULL,
  `fatigue_spent` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`character_id`,`completion_hash`),
  KEY `clear_time` (`finished_at`),
  KEY `clear_dungeon` (`dungeon_id`,`finished_at`),
  KEY `clear_account` (`account_id`,`finished_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_uca1400_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_dungeon_permissions` (
  `character_id` varchar(128) NOT NULL,
  `dungeon_id` bigint(20) unsigned NOT NULL,
  `clear_state` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `sort_order` int(11) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`,`dungeon_id`),
  KEY `idx_dnf_dungeon_permissions_character_sort` (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_equipment_entries` (
  `character_id` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `slot_index` smallint(6) NOT NULL DEFAULT 0,
  `item_id` bigint(20) NOT NULL DEFAULT 0,
  `bind_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `expire_at` datetime(6) DEFAULT NULL,
  `raw_entry` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`,`entry_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_equipment_entry_extra` (
  `character_id` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `extra_key` varchar(128) NOT NULL,
  `extra_value` longtext NOT NULL,
  PRIMARY KEY (`character_id`,`entry_key`,`extra_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_equipments` (
  `character_id` varchar(128) NOT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_hell_party_permissions` (
  `character_id` varchar(128) NOT NULL,
  `area_id` bigint(20) unsigned NOT NULL,
  `quest_id` bigint(20) unsigned NOT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`,`area_id`),
  KEY `idx_dnf_hell_party_permissions_character_quest` (`character_id`,`quest_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_inventories` (
  `character_id` varchar(128) NOT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_inventory_item_extra` (
  `character_id` varchar(128) NOT NULL,
  `collection_name` varchar(32) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `extra_key` varchar(128) NOT NULL,
  `extra_value` longtext NOT NULL,
  PRIMARY KEY (`character_id`,`collection_name`,`entry_key`,`extra_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_inventory_items` (
  `character_id` varchar(128) NOT NULL,
  `collection_name` varchar(32) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `item_id` bigint(20) NOT NULL DEFAULT 0,
  `item_count` bigint(20) NOT NULL DEFAULT 0,
  `bind_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `expire_at` datetime(6) DEFAULT NULL,
  `raw_entry` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`,`collection_name`,`entry_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_account_cargo_state` (
  `account_id` bigint(20) NOT NULL,
  `selection_key` bigint(20) NOT NULL DEFAULT 0,
  `value32` bigint(20) NOT NULL DEFAULT 0,
  `item_count` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_account_character_entries` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `entry_index` bigint(20) NOT NULL,
  `slot_index` bigint(20) NOT NULL,
  `name` varchar(255) NOT NULL,
  `name_bytes` longblob DEFAULT NULL,
  `body_after_name` longblob NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `entry_index` (`entry_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_account_settings` (
  `account_id` bigint(20) NOT NULL,
  `main_game_option` longblob DEFAULT NULL,
  `quickchat_bank0` longblob DEFAULT NULL,
  `quickchat_bank1` longblob DEFAULT NULL,
  `hotkey_key_type` bigint(20) NOT NULL DEFAULT 0,
  `hotkey_slots` longblob DEFAULT NULL,
  PRIMARY KEY (`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_accounts` (
  `account_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `m_id` varchar(255) NOT NULL,
  `password_hash` varchar(255) NOT NULL DEFAULT '',
  `last_login_ip` varchar(255) NOT NULL DEFAULT '',
  `last_login_at` varchar(255) DEFAULT NULL,
  `created_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  `cera` bigint(20) NOT NULL DEFAULT 0,
  `token_cera` bigint(20) NOT NULL DEFAULT 0,
  `happy_token_cera` bigint(20) NOT NULL DEFAULT 0,
  `lucky_star` bigint(20) NOT NULL DEFAULT 0,
  `cube_black` bigint(20) NOT NULL DEFAULT 0,
  `cube_white` bigint(20) NOT NULL DEFAULT 0,
  `cube_red` bigint(20) NOT NULL DEFAULT 0,
  `cube_blue` bigint(20) NOT NULL DEFAULT 0,
  `cube_clear` bigint(20) NOT NULL DEFAULT 0,
  `cube_gold` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`account_id`),
  UNIQUE KEY `m_id` (`m_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_abuse_values` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `abuse_value` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_achievement_chunks` (
  `character_id` bigint(20) NOT NULL,
  `chunk_index` bigint(20) NOT NULL,
  `mode_byte` bigint(20) NOT NULL DEFAULT 0,
  `owner_id16` bigint(20) NOT NULL DEFAULT 0,
  `entries_blob` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`,`chunk_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_achievement_complete` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `achievement_id` bigint(20) NOT NULL,
  `p1` bigint(20) NOT NULL DEFAULT 0,
  `p2` bigint(20) NOT NULL DEFAULT 0,
  `p3` bigint(20) NOT NULL DEFAULT 0,
  `p4` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_container_state` (
  `character_id` bigint(20) NOT NULL,
  `list_type` bigint(20) NOT NULL,
  `list_param16` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`list_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_creatures` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `creature_key` bigint(20) NOT NULL,
  `field04` bigint(20) NOT NULL DEFAULT 0,
  `mode_flag` bigint(20) NOT NULL DEFAULT 0,
  `progress_value` bigint(20) NOT NULL DEFAULT 0,
  `mode1_field0a` bigint(20) NOT NULL DEFAULT 0,
  `mode1_field0b` bigint(20) NOT NULL DEFAULT 0,
  `field_after_value` bigint(20) NOT NULL DEFAULT 0,
  `creature_text` longblob DEFAULT NULL,
  `tail_flag` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_dimension_flags` (
  `character_id` bigint(20) NOT NULL,
  `flag1` bigint(20) NOT NULL DEFAULT 0,
  `flag2` bigint(20) NOT NULL DEFAULT 0,
  `flag3` bigint(20) NOT NULL DEFAULT 0,
  `flag4` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_dimensions` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `dim_key` bigint(20) NOT NULL,
  `val1` bigint(20) NOT NULL DEFAULT 0,
  `val2` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_dungeon_permissions` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `dungeon_id` bigint(20) NOT NULL,
  `clear_state` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_equipped_entries` (
  `character_id` bigint(20) NOT NULL,
  `slot` bigint(20) NOT NULL,
  `item_id` bigint(20) NOT NULL,
  `expire_time` bigint(20) NOT NULL DEFAULT 0,
  `raw_entry` longblob NOT NULL,
  PRIMARY KEY (`character_id`,`slot`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_event_info` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `repeat_event_index` bigint(20) NOT NULL,
  `event_data` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_growth_weapon_stages` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `stage_id` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_hotkey_slots` (
  `character_id` bigint(20) NOT NULL,
  `slot_index` bigint(20) NOT NULL,
  `hotkey_value` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`slot_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_init_bodies` (
  `character_id` bigint(20) NOT NULL,
  `noti_type` bigint(20) NOT NULL,
  `occurrence_index` bigint(20) NOT NULL DEFAULT 0,
  `body` longblob NOT NULL,
  PRIMARY KEY (`character_id`,`noti_type`,`occurrence_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_init_flags` (
  `character_id` bigint(20) NOT NULL,
  `shop_coin_event_flag` bigint(20) NOT NULL DEFAULT 0,
  `level60_ui_state` bigint(20) NOT NULL DEFAULT 0,
  `pc_room_state` bigint(20) NOT NULL DEFAULT 0,
  `expert_job_blob` longblob DEFAULT NULL,
  `champion_break_blob` longblob DEFAULT NULL,
  `boss_tower_placeholder` bigint(20) NOT NULL DEFAULT 0,
  `mailbox_loaded_count` bigint(20) NOT NULL DEFAULT 0,
  `mailbox_mode` bigint(20) NOT NULL DEFAULT 0,
  `mailbox_not_loaded_count` bigint(20) NOT NULL DEFAULT 0,
  `mailbox_unknown_count_c` bigint(20) NOT NULL DEFAULT 0,
  `event_info_tail_byte` bigint(20) NOT NULL DEFAULT 0,
  `hotkey_key_type` bigint(20) NOT NULL DEFAULT 0,
  `main_game_option_blob` longblob DEFAULT NULL,
  `quickchat_bank0` longblob DEFAULT NULL,
  `quickchat_bank1` longblob DEFAULT NULL,
  `charac_invisible_falgs_payload_len` bigint(20) NOT NULL DEFAULT 0,
  `racing_dungeon_current_enter_count` bigint(20) NOT NULL DEFAULT 0,
  `racing_dungeon_group_flags` longblob DEFAULT NULL,
  `ack_account_reg_time` bigint(20) NOT NULL DEFAULT 0,
  `ack_premium_blob` longblob DEFAULT NULL,
  `ack_quest_display_ids` longblob DEFAULT NULL,
  `ack_char_slot_index` bigint(20) NOT NULL DEFAULT 0,
  `ack_fatigue_battery` bigint(20) NOT NULL DEFAULT 0,
  `ack_fatigue_grownup_buff` bigint(20) NOT NULL DEFAULT 0,
  `ack_trade_punish_flag` bigint(20) NOT NULL DEFAULT 0,
  `ack_extra_field_86jp` bigint(20) NOT NULL DEFAULT 0,
  `ack_reserved_8b` longblob DEFAULT NULL,
  `ack_tutorial_skipable` bigint(20) NOT NULL DEFAULT 0,
  `ack_post_tutorial_u16` bigint(20) NOT NULL DEFAULT 0,
  `ack_unread_tail` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_invisible_falgs` (
  `character_id` bigint(20) NOT NULL,
  `slot_index` bigint(20) NOT NULL,
  `flag_value` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`slot_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_item_extra` (
  `item_uid` bigint(20) NOT NULL,
  `extra_key` varchar(128) NOT NULL,
  `extra_value` longtext NOT NULL,
  PRIMARY KEY (`item_uid`,`extra_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_item_locks` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `type_or_list` bigint(20) NOT NULL,
  `item_key_or_slot` bigint(20) NOT NULL,
  `state` bigint(20) NOT NULL,
  `extra_value` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_item_values` (
  `character_id` bigint(20) NOT NULL,
  `list_kind` varchar(255) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `item_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`list_kind`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_items` (
  `item_uid` bigint(20) NOT NULL AUTO_INCREMENT,
  `owner_scope` varchar(255) NOT NULL,
  `owner_id` bigint(20) NOT NULL,
  `character_id` bigint(20) DEFAULT NULL,
  `list_type` bigint(20) NOT NULL,
  `slot_index` bigint(20) NOT NULL,
  `item_template_id` bigint(20) NOT NULL,
  `item_kind` varchar(255) NOT NULL DEFAULT 'unknown',
  `stack_count` bigint(20) NOT NULL DEFAULT 0,
  `instance_value` bigint(20) NOT NULL DEFAULT 0,
  `durability` bigint(20) NOT NULL DEFAULT 0,
  `seal_flag` bigint(20) NOT NULL DEFAULT 0,
  `option_value` bigint(20) NOT NULL DEFAULT 0,
  `expire_time` bigint(20) NOT NULL DEFAULT 0,
  `marker_16` bigint(20) NOT NULL DEFAULT 0,
  `pet_serial_or_handle` bigint(20) NOT NULL DEFAULT 0,
  `created_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`item_uid`),
  UNIQUE KEY `owner_scope` (`owner_scope`,`owner_id`,`list_type`,`slot_index`,`item_kind`),
  KEY `idx_character_items_owner_container` (`owner_scope`,`owner_id`,`list_type`,`slot_index`),
  KEY `idx_character_items_template` (`item_template_id`),
  KEY `idx_character_items_character` (`character_id`,`list_type`,`slot_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_pvp_missions` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `mission_id` bigint(20) NOT NULL,
  `progress_value` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_pvp_results` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `value_u32` bigint(20) NOT NULL DEFAULT 0,
  `value_u16a` bigint(20) NOT NULL DEFAULT 0,
  `value_u16b` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_racing_dungeon_entries` (
  `character_id` bigint(20) NOT NULL,
  `group_index` bigint(20) NOT NULL,
  `entry_index` bigint(20) NOT NULL,
  `track_like_id` bigint(20) NOT NULL,
  `value_a` bigint(20) NOT NULL,
  `value_b` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`group_index`,`entry_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_racing_dungeon_groups` (
  `character_id` bigint(20) NOT NULL,
  `group_index` bigint(20) NOT NULL,
  `group_id` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`group_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_racing_dungeon_tail_ids` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `id_value` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_show_effects` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `effect_index` bigint(20) NOT NULL,
  `duration_seconds` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_skill_points` (
  `character_id` bigint(20) NOT NULL,
  `total_sp` bigint(20) NOT NULL DEFAULT 0,
  `remaining_sp` bigint(20) NOT NULL DEFAULT 0,
  `total_sfp` bigint(20) NOT NULL DEFAULT 0,
  `remaining_sfp` bigint(20) NOT NULL DEFAULT 0,
  `total_tp` bigint(20) NOT NULL DEFAULT 0,
  `remaining_tp` bigint(20) NOT NULL DEFAULT 0,
  `synced_level` bigint(20) NOT NULL DEFAULT 1,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_skill_tail` (
  `character_id` bigint(20) NOT NULL,
  `tail0` bigint(20) NOT NULL DEFAULT 0,
  `tail1` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_skills` (
  `character_id` bigint(20) NOT NULL,
  `page_index` bigint(20) NOT NULL,
  `page_header` bigint(20) NOT NULL DEFAULT 0,
  `slot` bigint(20) NOT NULL,
  `skill_id` bigint(20) NOT NULL,
  `level` bigint(20) NOT NULL DEFAULT 0,
  `extra_values` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`,`page_index`,`slot`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_sort_item_locks` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `list_type` bigint(20) NOT NULL,
  `slot_index` bigint(20) NOT NULL,
  `state` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`sort_order`),
  UNIQUE KEY `character_id` (`character_id`,`list_type`,`slot_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_subtype0_fields` (
  `character_id` bigint(20) NOT NULL,
  `name_tag_item_id` bigint(20) NOT NULL DEFAULT 0,
  `creature_field1` bigint(20) NOT NULL DEFAULT 0,
  `creature_field2` bigint(20) NOT NULL DEFAULT 0,
  `creature_field3` bigint(20) NOT NULL DEFAULT 0,
  `creature_field4` bigint(20) NOT NULL DEFAULT 0,
  `creature_buffer` longblob DEFAULT NULL,
  `stamina` bigint(20) NOT NULL DEFAULT 0,
  `fatigue_penalty` bigint(20) NOT NULL DEFAULT 0,
  `is_event_character` bigint(20) NOT NULL DEFAULT 0,
  `pc_room_id` bigint(20) NOT NULL DEFAULT 65537,
  `is_private_store` bigint(20) NOT NULL DEFAULT 0,
  `is_premium_pc_room` bigint(20) NOT NULL DEFAULT 0,
  `server_group_id` bigint(20) NOT NULL DEFAULT 0,
  `black_count` bigint(20) NOT NULL DEFAULT 0,
  `guild_level` bigint(20) NOT NULL DEFAULT 0,
  `chaos_point` bigint(20) NOT NULL DEFAULT 0,
  `disguise_kind` bigint(20) NOT NULL DEFAULT 0,
  `is_disguised` bigint(20) NOT NULL DEFAULT 0,
  `expert_job_type` bigint(20) NOT NULL DEFAULT 0,
  `expert_job_exp` bigint(20) NOT NULL DEFAULT 0,
  `extra46` bigint(20) NOT NULL DEFAULT 0,
  `extra47` bigint(20) NOT NULL DEFAULT 0,
  `extra51` bigint(20) NOT NULL DEFAULT 0,
  `is_hardcore_mode` bigint(20) NOT NULL DEFAULT 0,
  `is_hardcore_dead` bigint(20) NOT NULL DEFAULT 0,
  `hardcore_death_count` bigint(20) NOT NULL DEFAULT 0,
  `user_state_bits` bigint(20) NOT NULL DEFAULT 3,
  `chat_ban_end_time` bigint(20) NOT NULL DEFAULT 0,
  `fatigue_update` bigint(20) NOT NULL DEFAULT 0,
  `return_user_flag` bigint(20) NOT NULL DEFAULT 1,
  `channel_display_mode` bigint(20) NOT NULL DEFAULT 0,
  `channel_type` bigint(20) NOT NULL DEFAULT 0,
  `channel_id` bigint(20) NOT NULL DEFAULT 2,
  `is_return_user` bigint(20) NOT NULL DEFAULT 0,
  `link_slot_enabled` bigint(20) NOT NULL DEFAULT 0,
  `link_type_a` bigint(20) NOT NULL DEFAULT 0,
  `link_type_b` bigint(20) NOT NULL DEFAULT 0,
  `emotion_index` bigint(20) NOT NULL DEFAULT 0,
  `action_byte` bigint(20) NOT NULL DEFAULT 0,
  `fatigue_display_update` bigint(20) NOT NULL DEFAULT 0,
  `costume_flag` bigint(20) NOT NULL DEFAULT 0,
  `aura_flag` bigint(20) NOT NULL DEFAULT 0,
  `pet_display_flag` bigint(20) NOT NULL DEFAULT 0,
  `title_display_flag` bigint(20) NOT NULL DEFAULT 0,
  `pvp_stat_a` bigint(20) NOT NULL DEFAULT 0,
  `pvp_win_streak` bigint(20) NOT NULL DEFAULT 0,
  `pvp_lose_streak` bigint(20) NOT NULL DEFAULT 0,
  `pvp_rank_point` bigint(20) NOT NULL DEFAULT 0,
  `trailing_byte` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_subtype1_fields` (
  `character_id` bigint(20) NOT NULL,
  `stat_hp_max` bigint(20) NOT NULL DEFAULT 0,
  `stat_mp_max` bigint(20) NOT NULL DEFAULT 0,
  `stat_strength` bigint(20) NOT NULL DEFAULT 0,
  `stat_intelligence` bigint(20) NOT NULL DEFAULT 0,
  `stat_vitality` bigint(20) NOT NULL DEFAULT 0,
  `stat_spirit` bigint(20) NOT NULL DEFAULT 0,
  `stat_physical_attack` bigint(20) NOT NULL DEFAULT 0,
  `stat_physical_defense` bigint(20) NOT NULL DEFAULT 0,
  `stat_magical_attack` bigint(20) NOT NULL DEFAULT 0,
  `stat_magical_defense` bigint(20) NOT NULL DEFAULT 0,
  `stat_independent_attack` bigint(20) NOT NULL DEFAULT 0,
  `stat_fire_resistance` bigint(20) NOT NULL DEFAULT 0,
  `stat_water_resistance` bigint(20) NOT NULL DEFAULT 0,
  `stat_dark_resistance` bigint(20) NOT NULL DEFAULT 0,
  `stat_light_resistance` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_00` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_01` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_02` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_03` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_04` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_05` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_06` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_07` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_08` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_09` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_10` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_11` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_12` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_13` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_14` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_15` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_16` bigint(20) NOT NULL DEFAULT 0,
  `active_status_resistance_17` bigint(20) NOT NULL DEFAULT 0,
  `stat_inventory_limit` bigint(20) NOT NULL DEFAULT 0,
  `stat_hp_regen_speed` bigint(20) NOT NULL DEFAULT 0,
  `stat_mp_regen_speed` bigint(20) NOT NULL DEFAULT 0,
  `stat_move_speed` bigint(20) NOT NULL DEFAULT 0,
  `stat_attack_speed` bigint(20) NOT NULL DEFAULT 0,
  `stat_cast_speed` bigint(20) NOT NULL DEFAULT 0,
  `stat_hit_recovery` bigint(20) NOT NULL DEFAULT 0,
  `stat_jump_power` bigint(20) NOT NULL DEFAULT 0,
  `stat_weight` bigint(20) NOT NULL DEFAULT 0,
  `stat_level` bigint(20) NOT NULL DEFAULT 0,
  `name_tag_item_id` bigint(20) NOT NULL DEFAULT 0,
  `name_tag_expire_time` bigint(20) NOT NULL DEFAULT 0,
  `skill_tree_index` bigint(20) NOT NULL DEFAULT 0,
  `equipped_creature_level` bigint(20) NOT NULL DEFAULT 0,
  `equip_list_trailing` bigint(20) NOT NULL DEFAULT 0,
  `manage_level` bigint(20) NOT NULL DEFAULT 0,
  `flag_byte` bigint(20) NOT NULL DEFAULT 0,
  `guild_power_war` bigint(20) NOT NULL DEFAULT 0,
  `server_timestamp` bigint(20) NOT NULL DEFAULT 0,
  `quest_shop_count` bigint(20) NOT NULL DEFAULT 0,
  `progress1` bigint(20) NOT NULL DEFAULT 0,
  `progress2` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_unknown725` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `param_a` bigint(20) NOT NULL,
  `mode_or_state` bigint(20) NOT NULL,
  `content_id` bigint(20) NOT NULL,
  `param_b` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_unknown730` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL,
  `entry_id` bigint(20) NOT NULL,
  `sentinel_or_value` bigint(20) NOT NULL,
  `flag` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo109_state` (
  `character_id` bigint(20) NOT NULL,
  `value0` bigint(20) NOT NULL DEFAULT 0,
  `value1` bigint(20) NOT NULL DEFAULT 0,
  `value2` bigint(20) NOT NULL DEFAULT 0,
  `value3` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo10c_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo117_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo118_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key` bigint(20) NOT NULL DEFAULT 0,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `value_d` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo11d_state` (
  `character_id` bigint(20) NOT NULL,
  `mode` bigint(20) NOT NULL DEFAULT 0,
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `byte_b` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo126_state` (
  `character_id` bigint(20) NOT NULL,
  `text` varchar(255) NOT NULL DEFAULT '',
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo12a_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo154_state` (
  `character_id` bigint(20) NOT NULL,
  `key_a` bigint(20) NOT NULL DEFAULT 0,
  `slot_or_value_a` bigint(20) NOT NULL DEFAULT 0,
  `key_b` bigint(20) NOT NULL DEFAULT 0,
  `delta` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo159_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `slot` bigint(20) NOT NULL DEFAULT 255,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo17c_control` (
  `character_id` bigint(20) NOT NULL,
  `header_flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo17c_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `selector` bigint(20) NOT NULL DEFAULT 0,
  `item_or_key` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo182_control` (
  `character_id` bigint(20) NOT NULL,
  `actor_key` bigint(20) NOT NULL DEFAULT 0,
  `header_flag` bigint(20) NOT NULL DEFAULT 0,
  `outer_count` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo182_first_rows` (
  `character_id` bigint(20) NOT NULL,
  `group_index` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `row_state` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_index`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo182_first_values` (
  `character_id` bigint(20) NOT NULL,
  `group_index` bigint(20) NOT NULL DEFAULT 0,
  `row_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value_index` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_index`,`row_sort_order`,`value_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo182_groups` (
  `character_id` bigint(20) NOT NULL,
  `phase` bigint(20) NOT NULL DEFAULT 0,
  `group_index` bigint(20) NOT NULL DEFAULT 0,
  `group_flag` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`phase`,`group_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo182_second_values` (
  `character_id` bigint(20) NOT NULL,
  `group_index` bigint(20) NOT NULL DEFAULT 0,
  `value_index` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_index`,`value_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo183_control` (
  `character_id` bigint(20) NOT NULL,
  `header_a` bigint(20) NOT NULL DEFAULT 0,
  `header_b` bigint(20) NOT NULL DEFAULT 0,
  `header_value` bigint(20) NOT NULL DEFAULT 0,
  `global_flag` bigint(20) NOT NULL DEFAULT 0,
  `mode` bigint(20) NOT NULL DEFAULT 0,
  `extra_value` bigint(20) NOT NULL DEFAULT 0,
  `tail_flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo183_primary_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `word0` bigint(20) NOT NULL DEFAULT 0,
  `key_or_value` bigint(20) NOT NULL DEFAULT 0,
  `word1` bigint(20) NOT NULL DEFAULT 0,
  `word2` bigint(20) NOT NULL DEFAULT 0,
  `flag0` bigint(20) NOT NULL DEFAULT 0,
  `flag1` bigint(20) NOT NULL DEFAULT 0,
  `bool_flag` bigint(20) NOT NULL DEFAULT 0,
  `flag2` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo183_secondary_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `row_type` bigint(20) NOT NULL DEFAULT 0,
  `value0` bigint(20) NOT NULL DEFAULT 0,
  `value1` bigint(20) NOT NULL DEFAULT 0,
  `value2` bigint(20) NOT NULL DEFAULT 0,
  `word0` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `word1` bigint(20) NOT NULL DEFAULT 0,
  `word2` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo184_control` (
  `character_id` bigint(20) NOT NULL,
  `header` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo184_first_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo184_second_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo184_third_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo1bf_control` (
  `character_id` bigint(20) NOT NULL,
  `refresh_flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo1bf_groups` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `context_state` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo1bf_selectors` (
  `character_id` bigint(20) NOT NULL,
  `group_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `selector` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_sort_order`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo1bf_values` (
  `character_id` bigint(20) NOT NULL,
  `group_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `selector_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_sort_order`,`selector_sort_order`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo22d_state` (
  `character_id` bigint(20) NOT NULL,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `word_c` bigint(20) NOT NULL DEFAULT 0,
  `word_d` bigint(20) NOT NULL DEFAULT 0,
  `mode` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo22e_pairs` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo22e_state` (
  `character_id` bigint(20) NOT NULL,
  `object_key` bigint(20) NOT NULL DEFAULT 0,
  `mode` bigint(20) NOT NULL DEFAULT 0,
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `byte_b` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo237_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo238_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo23_fixed_values` (
  `character_id` bigint(20) NOT NULL,
  `section` varchar(255) NOT NULL,
  `slot_index` bigint(20) NOT NULL,
  `source_order` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`section`,`slot_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo23_object_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `object_or_slot_id` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `value_d` bigint(20) NOT NULL DEFAULT 0,
  `value_e` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo23_pair_entries` (
  `character_id` bigint(20) NOT NULL,
  `section` varchar(255) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `entry_key` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`section`,`sort_order`),
  KEY `idx_character_userinfo23_pair_section` (`character_id`,`section`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo23_scalar_values` (
  `character_id` bigint(20) NOT NULL,
  `family` varchar(255) NOT NULL,
  `scalar_index` bigint(20) NOT NULL,
  `source_order` bigint(20) NOT NULL DEFAULT 0,
  `byte_offset` bigint(20) NOT NULL DEFAULT -1,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `width_bits` bigint(20) NOT NULL DEFAULT 32,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`family`,`scalar_index`),
  KEY `idx_character_userinfo23_scalar_order` (`character_id`,`source_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo253_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo253_state` (
  `character_id` bigint(20) NOT NULL,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `tail_word_a` bigint(20) NOT NULL DEFAULT 0,
  `tail_word_b` bigint(20) NOT NULL DEFAULT 0,
  `tail_flag` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo254_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo254_state` (
  `character_id` bigint(20) NOT NULL,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo255_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo25b_state` (
  `character_id` bigint(20) NOT NULL,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo26e_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo274_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo275_state` (
  `character_id` bigint(20) NOT NULL,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `word_c` bigint(20) NOT NULL DEFAULT 0,
  `word_d` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo276_state` (
  `character_id` bigint(20) NOT NULL,
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `byte_b` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo287_extras` (
  `character_id` bigint(20) NOT NULL,
  `row_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `extra_index` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`row_sort_order`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo287_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `name` varchar(255) NOT NULL DEFAULT '',
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `byte_b` bigint(20) NOT NULL DEFAULT 0,
  `packed_flag` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `text` varchar(255) NOT NULL DEFAULT '',
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `value_d` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `byte_c` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo287_state` (
  `character_id` bigint(20) NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo28a_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `row_type` bigint(20) NOT NULL DEFAULT 0,
  `text` varchar(255) NOT NULL DEFAULT '',
  `word` bigint(20) NOT NULL DEFAULT 0,
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `packed_flag` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo28a_state` (
  `character_id` bigint(20) NOT NULL,
  `tail_value_a` bigint(20) NOT NULL DEFAULT 0,
  `tail_value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo28b_state` (
  `character_id` bigint(20) NOT NULL,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo29f_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `category` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo29f_state` (
  `character_id` bigint(20) NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2a9_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2a9_state` (
  `character_id` bigint(20) NOT NULL,
  `header_word` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2a9_values` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2aa_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2b0_state` (
  `character_id` bigint(20) NOT NULL,
  `key` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2bc_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2c1_state` (
  `character_id` bigint(20) NOT NULL,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2d2_groups` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `group_key` bigint(20) NOT NULL DEFAULT 0,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2d2_pairs` (
  `character_id` bigint(20) NOT NULL,
  `group_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `row_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `pair_kind` varchar(255) NOT NULL DEFAULT 'first',
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_sort_order`,`row_sort_order`,`pair_kind`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2d2_rows` (
  `character_id` bigint(20) NOT NULL,
  `group_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_sort_order`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2d2_state` (
  `character_id` bigint(20) NOT NULL,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2d3_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2d8_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo2ef_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo31d_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo324_state` (
  `character_id` bigint(20) NOT NULL,
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `text` varchar(255) NOT NULL DEFAULT '',
  `byte_b` bigint(20) NOT NULL DEFAULT 0,
  `byte_c` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo327_blobs` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `blob_key` bigint(20) NOT NULL DEFAULT 0,
  `payload` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo329_groups` (
  `character_id` bigint(20) NOT NULL,
  `target_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `context_state` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`target_sort_order`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo329_selectors` (
  `character_id` bigint(20) NOT NULL,
  `target_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `group_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `selector` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`target_sort_order`,`group_sort_order`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo329_targets` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `target_key` bigint(20) NOT NULL DEFAULT 0,
  `refresh_flag` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo329_values` (
  `character_id` bigint(20) NOT NULL,
  `target_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `group_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `selector_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`target_sort_order`,`group_sort_order`,`selector_sort_order`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo336_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `text` varchar(255) NOT NULL DEFAULT '',
  `byte_b` bigint(20) NOT NULL DEFAULT 0,
  `byte_c` bigint(20) NOT NULL DEFAULT 0,
  `packed_flag` bigint(20) NOT NULL DEFAULT 0,
  `byte_d` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo336_state` (
  `character_id` bigint(20) NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34b_state` (
  `character_id` bigint(20) NOT NULL,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `state` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34c_control` (
  `character_id` bigint(20) NOT NULL,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34c_first_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key` bigint(20) NOT NULL DEFAULT 0,
  `word` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `flag_c` bigint(20) NOT NULL DEFAULT 0,
  `flag_d` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34c_second_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `row_type` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `word_c` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34c_text_state` (
  `character_id` bigint(20) NOT NULL,
  `category` bigint(20) NOT NULL DEFAULT 0,
  `text` varchar(255) NOT NULL DEFAULT '',
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34d_control` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34d_rows` (
  `character_id` bigint(20) NOT NULL,
  `group_index` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_index`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34d_value_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34e_byte_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo34e_state` (
  `character_id` bigint(20) NOT NULL,
  `state` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo352_state` (
  `character_id` bigint(20) NOT NULL,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo354_state` (
  `character_id` bigint(20) NOT NULL,
  `word0` bigint(20) NOT NULL DEFAULT 0,
  `word1` bigint(20) NOT NULL DEFAULT 0,
  `word2` bigint(20) NOT NULL DEFAULT 0,
  `word3` bigint(20) NOT NULL DEFAULT 0,
  `word4` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo355_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo359_groups` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `group_key` bigint(20) NOT NULL DEFAULT 0,
  `raw0` longblob NOT NULL,
  `raw1` longblob NOT NULL,
  `raw2` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo359_state` (
  `character_id` bigint(20) NOT NULL,
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `byte_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo36b_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `kind` bigint(20) NOT NULL DEFAULT 0,
  `byte_a` bigint(20) NOT NULL DEFAULT 0,
  `byte_b` bigint(20) NOT NULL DEFAULT 0,
  `byte_c` bigint(20) NOT NULL DEFAULT 0,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `byte_d` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo36b_state` (
  `character_id` bigint(20) NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo374_rows` (
  `character_id` bigint(20) NOT NULL,
  `group_kind` varchar(255) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `payload` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_kind`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo374_state` (
  `character_id` bigint(20) NOT NULL,
  `header` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo379_rows` (
  `character_id` bigint(20) NOT NULL,
  `group_kind` varchar(255) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `payload` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_kind`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo379_state` (
  `character_id` bigint(20) NOT NULL,
  `header` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo37a_rows` (
  `character_id` bigint(20) NOT NULL,
  `group_kind` varchar(255) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `payload` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_kind`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo37a_state` (
  `character_id` bigint(20) NOT NULL,
  `header1` longblob NOT NULL,
  `header33` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo37b_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `raw124` longblob NOT NULL,
  `raw64` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo393_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo393_state` (
  `character_id` bigint(20) NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo3cd_state` (
  `character_id` bigint(20) NOT NULL,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo3e6_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo57_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `object_key` bigint(20) NOT NULL DEFAULT 0,
  `field_a` bigint(20) NOT NULL DEFAULT 0,
  `route_or_index` bigint(20) NOT NULL DEFAULT 0,
  `field_c` bigint(20) NOT NULL DEFAULT 0,
  `state` bigint(20) NOT NULL DEFAULT 0,
  `value32` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo57_slots` (
  `character_id` bigint(20) NOT NULL,
  `row_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `slot_index` bigint(20) NOT NULL,
  `mode` bigint(20) NOT NULL DEFAULT 255,
  `value` bigint(20) NOT NULL DEFAULT 65535,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`row_sort_order`,`slot_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo58_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `object_key` bigint(20) NOT NULL DEFAULT 0,
  `state` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo59_control` (
  `character_id` bigint(20) NOT NULL,
  `object_key` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo59_slots` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `category` bigint(20) NOT NULL,
  `mode` bigint(20) NOT NULL DEFAULT 255,
  `value` bigint(20) NOT NULL DEFAULT 65535,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo5b_control` (
  `character_id` bigint(20) NOT NULL,
  `header_flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo5b_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `value_d` bigint(20) NOT NULL DEFAULT 0,
  `value_e` bigint(20) NOT NULL DEFAULT 0,
  `value_f` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo5c_state` (
  `character_id` bigint(20) NOT NULL,
  `header_a` bigint(20) NOT NULL DEFAULT 0,
  `header_b` bigint(20) NOT NULL DEFAULT 0,
  `state0` bigint(20) NOT NULL DEFAULT 0,
  `state1` bigint(20) NOT NULL DEFAULT 0,
  `state2` bigint(20) NOT NULL DEFAULT 0,
  `state3` bigint(20) NOT NULL DEFAULT 0,
  `state4` bigint(20) NOT NULL DEFAULT 0,
  `state5` bigint(20) NOT NULL DEFAULT 0,
  `state6` bigint(20) NOT NULL DEFAULT 0,
  `state7` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo5f_state` (
  `character_id` bigint(20) NOT NULL,
  `category` bigint(20) NOT NULL DEFAULT 0,
  `mode_or_apply_flag` bigint(20) NOT NULL DEFAULT 0,
  `scale_or_visual_flag` bigint(20) NOT NULL DEFAULT 0,
  `delta_value` bigint(20) NOT NULL DEFAULT 0,
  `existing_slot_value` bigint(20) NOT NULL DEFAULT 65535,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo60_pairs` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo60_wide_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key` bigint(20) NOT NULL DEFAULT 0,
  `value0` bigint(20) NOT NULL DEFAULT 0,
  `value1` bigint(20) NOT NULL DEFAULT 0,
  `value2` bigint(20) NOT NULL DEFAULT 0,
  `value3` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo64_state` (
  `character_id` bigint(20) NOT NULL,
  `object_key` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo67_state` (
  `character_id` bigint(20) NOT NULL,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo6a_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `object_key` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo6b_state` (
  `character_id` bigint(20) NOT NULL,
  `object_key` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo73_state` (
  `character_id` bigint(20) NOT NULL,
  `mode` bigint(20) NOT NULL DEFAULT 0,
  `state` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 65535,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo7a_control` (
  `character_id` bigint(20) NOT NULL,
  `mode` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo7a_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key_a` bigint(20) NOT NULL DEFAULT 0,
  `key_b` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo80_control` (
  `character_id` bigint(20) NOT NULL,
  `actor_key` bigint(20) NOT NULL DEFAULT 0,
  `profile_a` bigint(20) NOT NULL DEFAULT 0,
  `profile_b` bigint(20) NOT NULL DEFAULT 0,
  `route` bigint(20) NOT NULL DEFAULT 0,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo80_extra_words` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo80_slots` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `category` bigint(20) NOT NULL,
  `mode` bigint(20) NOT NULL DEFAULT 255,
  `value` bigint(20) NOT NULL DEFAULT 65535,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`),
  KEY `idx_character_userinfo80_slots_category` (`character_id`,`category`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo81_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo83_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo85_state` (
  `character_id` bigint(20) NOT NULL,
  `object_key` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo86_children` (
  `character_id` bigint(20) NOT NULL,
  `row_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `child_key` bigint(20) NOT NULL DEFAULT 0,
  `state` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`row_sort_order`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo86_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key_value` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo87_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo88_state` (
  `character_id` bigint(20) NOT NULL,
  `object_key` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo89_state` (
  `character_id` bigint(20) NOT NULL,
  `state` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo8f_control` (
  `character_id` bigint(20) NOT NULL,
  `context_key` bigint(20) NOT NULL DEFAULT 0,
  `root_value` bigint(20) NOT NULL DEFAULT 0,
  `header_value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo8f_list_a` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key_value` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `bool_flag` bigint(20) NOT NULL DEFAULT 0,
  `state` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo8f_list_b` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `row_type` bigint(20) NOT NULL DEFAULT 0,
  `key_a` bigint(20) NOT NULL DEFAULT 0,
  `key_b` bigint(20) NOT NULL DEFAULT 0,
  `key_c` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo90_control` (
  `character_id` bigint(20) NOT NULL,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 1,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `include_primary_block` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo90_summaries` (
  `character_id` bigint(20) NOT NULL,
  `is_primary` bigint(20) NOT NULL DEFAULT 0,
  `group_index` bigint(20) NOT NULL DEFAULT -1,
  `summary_word` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`is_primary`,`group_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo90_text_rows` (
  `character_id` bigint(20) NOT NULL,
  `is_primary` bigint(20) NOT NULL DEFAULT 0,
  `group_index` bigint(20) NOT NULL DEFAULT -1,
  `slot_index` bigint(20) NOT NULL DEFAULT 0,
  `text_value` varchar(255) NOT NULL DEFAULT '',
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`is_primary`,`group_index`,`slot_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo91_control` (
  `character_id` bigint(20) NOT NULL,
  `header_value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo91_rows` (
  `character_id` bigint(20) NOT NULL,
  `group_index` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key_value` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_index`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo92_state` (
  `character_id` bigint(20) NOT NULL,
  `mode_flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo98_state` (
  `character_id` bigint(20) NOT NULL,
  `header` bigint(20) NOT NULL DEFAULT 0,
  `word_a` bigint(20) NOT NULL DEFAULT 0,
  `word_b` bigint(20) NOT NULL DEFAULT 0,
  `state0` bigint(20) NOT NULL DEFAULT 0,
  `state1` bigint(20) NOT NULL DEFAULT 0,
  `state2` bigint(20) NOT NULL DEFAULT 0,
  `state3` bigint(20) NOT NULL DEFAULT 0,
  `apply_flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo9b_state` (
  `character_id` bigint(20) NOT NULL,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo_blobs` (
  `character_id` bigint(20) NOT NULL,
  `blob_kind` varchar(255) NOT NULL,
  `subtype` bigint(20) NOT NULL,
  `user_info_type` bigint(20) NOT NULL DEFAULT 0,
  `gate_or_count` bigint(20) NOT NULL DEFAULT 0,
  `user_id` bigint(20) NOT NULL DEFAULT 0,
  `name_bytes` longblob DEFAULT NULL,
  `remaining_bytes` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`,`blob_kind`,`subtype`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo_byte_count_raw_rows` (
  `character_id` bigint(20) NOT NULL,
  `noti_type` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `payload` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`noti_type`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo_fixed_raws` (
  `character_id` bigint(20) NOT NULL,
  `noti_type` bigint(20) NOT NULL,
  `payload` longblob NOT NULL,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`noti_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo_slot_category_state` (
  `character_id` bigint(20) NOT NULL,
  `category` bigint(20) NOT NULL,
  `state_a` bigint(20) NOT NULL DEFAULT 0,
  `state_b` bigint(20) NOT NULL DEFAULT 0,
  `enable_flag` bigint(20) NOT NULL DEFAULT 0,
  `group12_active_flag` bigint(20) NOT NULL DEFAULT 0,
  `empty_marker` bigint(20) NOT NULL DEFAULT -1,
  `last_selection_tick` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo_slot_control` (
  `character_id` bigint(20) NOT NULL,
  `visible_category_count` bigint(20) NOT NULL DEFAULT 4,
  `group0_entry_total` bigint(20) NOT NULL DEFAULT 0,
  `group0_tail_refresh_value` bigint(20) NOT NULL DEFAULT 0,
  `refresh_value_a` bigint(20) NOT NULL DEFAULT 0,
  `refresh_value_b` bigint(20) NOT NULL DEFAULT 0,
  `selected_category_a` bigint(20) NOT NULL DEFAULT -1,
  `selected_category_b` bigint(20) NOT NULL DEFAULT -1,
  `msg72_value_a` bigint(20) NOT NULL DEFAULT 0,
  `msg72_state` bigint(20) NOT NULL DEFAULT 0,
  `render_flag_1949` bigint(20) NOT NULL DEFAULT 0,
  `render_flag_1950` bigint(20) NOT NULL DEFAULT 0,
  `special_group_flag_1951` bigint(20) NOT NULL DEFAULT 0,
  `special_group_mode_1952` bigint(20) NOT NULL DEFAULT 0,
  `mode_bits` bigint(20) NOT NULL DEFAULT 0,
  `tail_value_a` bigint(20) NOT NULL DEFAULT 0,
  `tail_flag_a` bigint(20) NOT NULL DEFAULT 0,
  `tail_bool_b` bigint(20) NOT NULL DEFAULT 0,
  `tail_value_b` bigint(20) NOT NULL DEFAULT 0,
  `tail_pair_a` bigint(20) NOT NULL DEFAULT 0,
  `tail_pair_b` bigint(20) NOT NULL DEFAULT 0,
  `tail_index6_value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfo_slot_group_entries` (
  `character_id` bigint(20) NOT NULL,
  `group_id` bigint(20) NOT NULL,
  `category` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key_or_item_id` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `value_width_bits` bigint(20) NOT NULL DEFAULT 32,
  `source_noti_type` bigint(20) NOT NULL DEFAULT 35,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_id`,`category`,`sort_order`),
  KEY `idx_character_userinfo_slot_group_lookup` (`character_id`,`category`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoa0_control` (
  `character_id` bigint(20) NOT NULL,
  `header_flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoa0_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `selector` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoa1_state` (
  `character_id` bigint(20) NOT NULL,
  `value0` bigint(20) NOT NULL DEFAULT 0,
  `value1` bigint(20) NOT NULL DEFAULT 0,
  `value2` bigint(20) NOT NULL DEFAULT 0,
  `value3` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoa2_state` (
  `character_id` bigint(20) NOT NULL,
  `mode` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoa3_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoaa_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfob0_state` (
  `character_id` bigint(20) NOT NULL,
  `enabled` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfob6_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `text_a` varchar(255) NOT NULL DEFAULT '',
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `flag_c` bigint(20) NOT NULL DEFAULT 0,
  `text_b` varchar(255) NOT NULL DEFAULT '',
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfob6_values` (
  `character_id` bigint(20) NOT NULL,
  `row_sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value_index` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`row_sort_order`,`value_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfobc_state` (
  `character_id` bigint(20) NOT NULL,
  `state` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoc8_state` (
  `character_id` bigint(20) NOT NULL,
  `delta` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoc9_state` (
  `character_id` bigint(20) NOT NULL,
  `value0` bigint(20) NOT NULL DEFAULT 0,
  `value1` bigint(20) NOT NULL DEFAULT 0,
  `value2` bigint(20) NOT NULL DEFAULT 0,
  `value3` bigint(20) NOT NULL DEFAULT 0,
  `value4` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfocf_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod0_state` (
  `character_id` bigint(20) NOT NULL,
  `value0` bigint(20) NOT NULL DEFAULT 0,
  `value1` bigint(20) NOT NULL DEFAULT 0,
  `value2` bigint(20) NOT NULL DEFAULT 0,
  `value3` bigint(20) NOT NULL DEFAULT 0,
  `value4` bigint(20) NOT NULL DEFAULT 0,
  `value5` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod1_control` (
  `character_id` bigint(20) NOT NULL,
  `header_a` bigint(20) NOT NULL DEFAULT 0,
  `header_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod1_rows` (
  `character_id` bigint(20) NOT NULL,
  `group_index` bigint(20) NOT NULL DEFAULT 0,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `key` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`group_index`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod2_control` (
  `character_id` bigint(20) NOT NULL,
  `tail_value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod2_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `row_type` bigint(20) NOT NULL DEFAULT 0,
  `context_key` bigint(20) NOT NULL DEFAULT 0,
  `key` bigint(20) NOT NULL DEFAULT 0,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `value_c` bigint(20) NOT NULL DEFAULT 0,
  `value_d` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod3_control` (
  `character_id` bigint(20) NOT NULL,
  `header_a` bigint(20) NOT NULL DEFAULT 0,
  `header_b` bigint(20) NOT NULL DEFAULT 0,
  `header_value` bigint(20) NOT NULL DEFAULT 0,
  `global_flag` bigint(20) NOT NULL DEFAULT 0,
  `mode` bigint(20) NOT NULL DEFAULT 0,
  `extra_value` bigint(20) NOT NULL DEFAULT 0,
  `tail_flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod3_primary_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `word0` bigint(20) NOT NULL DEFAULT 0,
  `byte0` bigint(20) NOT NULL DEFAULT 0,
  `byte1` bigint(20) NOT NULL DEFAULT 0,
  `word1` bigint(20) NOT NULL DEFAULT 0,
  `word2` bigint(20) NOT NULL DEFAULT 0,
  `word3` bigint(20) NOT NULL DEFAULT 0,
  `byte2` bigint(20) NOT NULL DEFAULT 0,
  `word4` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `byte3` bigint(20) NOT NULL DEFAULT 0,
  `byte4` bigint(20) NOT NULL DEFAULT 0,
  `bool_flag` bigint(20) NOT NULL DEFAULT 0,
  `byte5` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod3_secondary_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `row_type` bigint(20) NOT NULL DEFAULT 0,
  `value0` bigint(20) NOT NULL DEFAULT 0,
  `value1` bigint(20) NOT NULL DEFAULT 0,
  `value2` bigint(20) NOT NULL DEFAULT 0,
  `word0` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `word1` bigint(20) NOT NULL DEFAULT 0,
  `word2` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod5_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod6_state` (
  `character_id` bigint(20) NOT NULL,
  `flag_a` bigint(20) NOT NULL DEFAULT 0,
  `flag_b` bigint(20) NOT NULL DEFAULT 0,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod7_state` (
  `character_id` bigint(20) NOT NULL,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod8_control` (
  `character_id` bigint(20) NOT NULL,
  `mode` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfod8_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfodc_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfodd_state` (
  `character_id` bigint(20) NOT NULL,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfodf_state` (
  `character_id` bigint(20) NOT NULL,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `value0` bigint(20) NOT NULL DEFAULT 0,
  `value1` bigint(20) NOT NULL DEFAULT 0,
  `value2` bigint(20) NOT NULL DEFAULT 0,
  `value3` bigint(20) NOT NULL DEFAULT 0,
  `value4` bigint(20) NOT NULL DEFAULT 0,
  `value5` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoe0_control` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoe0_rows` (
  `character_id` bigint(20) NOT NULL,
  `sort_order` bigint(20) NOT NULL DEFAULT 0,
  `text` varchar(255) NOT NULL DEFAULT '',
  `note` varchar(255) NOT NULL DEFAULT '',
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`,`sort_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoe6_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoeb_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfofe_state` (
  `character_id` bigint(20) NOT NULL,
  `value_a` bigint(20) NOT NULL DEFAULT 0,
  `value_b` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_character_userinfoff_state` (
  `character_id` bigint(20) NOT NULL,
  `value` bigint(20) NOT NULL DEFAULT 0,
  `flag` bigint(20) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_characters` (
  `character_id` bigint(20) NOT NULL,
  `account_id` bigint(20) NOT NULL,
  `name` varchar(255) NOT NULL,
  `job` bigint(20) NOT NULL DEFAULT 0,
  `grow_type` bigint(20) NOT NULL DEFAULT 0,
  `level` bigint(20) NOT NULL DEFAULT 1,
  `exp` bigint(20) NOT NULL DEFAULT 0,
  `ex_equip_slot_stat` bigint(20) NOT NULL DEFAULT 0,
  `bonus_sp` bigint(20) NOT NULL DEFAULT 0,
  `bonus_tp` bigint(20) NOT NULL DEFAULT 0,
  `pvp_grade` bigint(20) NOT NULL DEFAULT 0,
  `pvp_rating_grade` bigint(20) NOT NULL DEFAULT 0,
  `user_state` bigint(20) NOT NULL DEFAULT 0,
  `gold` bigint(20) NOT NULL DEFAULT 0,
  `coin` bigint(20) NOT NULL DEFAULT 0,
  `town_id` bigint(20) NOT NULL DEFAULT 0,
  `area_id` bigint(20) NOT NULL DEFAULT 0,
  `pos_x` bigint(20) NOT NULL DEFAULT 0,
  `pos_y` bigint(20) NOT NULL DEFAULT 0,
  `direction` bigint(20) NOT NULL DEFAULT 5,
  `area_state` bigint(20) NOT NULL DEFAULT 3,
  `name_bytes` longblob DEFAULT NULL,
  `appearance_blob` longblob DEFAULT NULL,
  `delete_flag` bigint(20) NOT NULL DEFAULT 0,
  `created_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`character_id`),
  UNIQUE KEY `idx_characters_name_unique` (`name`),
  KEY `idx_characters_account` (`account_id`,`delete_flag`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_equipped_items` (
  `character_id` bigint(20) NOT NULL,
  `equip_list_blob` longblob NOT NULL,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_get_userinfo_template` (
  `id` bigint(20) NOT NULL DEFAULT 1,
  `seed_character_id` bigint(20) NOT NULL DEFAULT 1000,
  `response_blob` longblob DEFAULT NULL,
  `pkt0_routing_byte7` bigint(20) NOT NULL DEFAULT 0,
  `gate_or_count1` bigint(20) NOT NULL DEFAULT 32,
  `gate_or_count2` bigint(20) NOT NULL DEFAULT 32,
  `flag_or_manage` bigint(20) NOT NULL DEFAULT 2,
  `key_or_point` bigint(20) NOT NULL DEFAULT 0,
  `unknown16` bigint(20) NOT NULL DEFAULT 0,
  `unknown32` bigint(20) NOT NULL DEFAULT 0,
  `pkt2_result_code` bigint(20) NOT NULL DEFAULT 1,
  `pkt2_character_key` bigint(20) NOT NULL DEFAULT 0,
  `pkt2_slot_flag1` bigint(20) NOT NULL DEFAULT 0,
  `pkt2_slot_flag2` bigint(20) NOT NULL DEFAULT 1,
  `pkt2_state_flag` bigint(20) NOT NULL DEFAULT 255,
  `pkt2_flag3` bigint(20) NOT NULL DEFAULT 1,
  `pkt2_reserved` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_getuserinfo_extra_packets` (
  `seq` bigint(20) NOT NULL,
  `command` bigint(20) NOT NULL,
  `noti_type` bigint(20) NOT NULL,
  `body` longblob NOT NULL,
  PRIMARY KEY (`seq`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_global_raw_packets` (
  `noti_type` bigint(20) NOT NULL,
  `packet_body` longblob NOT NULL,
  PRIMARY KEY (`noti_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_global_server_event_phase` (
  `id` bigint(20) NOT NULL,
  `event_phase_bitmap` longblob NOT NULL,
  `updated_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_item_audit_log` (
  `audit_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `owner_scope` varchar(255) NOT NULL,
  `owner_id` bigint(20) NOT NULL,
  `character_id` bigint(20) DEFAULT NULL,
  `action_name` varchar(255) NOT NULL,
  `list_type` bigint(20) DEFAULT NULL,
  `slot_index` bigint(20) DEFAULT NULL,
  `item_uid` bigint(20) DEFAULT NULL,
  `item_template_id` bigint(20) DEFAULT NULL,
  `delta_stack_count` bigint(20) NOT NULL DEFAULT 0,
  `created_at` datetime(6) NOT NULL DEFAULT current_timestamp(6),
  PRIMARY KEY (`audit_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_item_audit_payload_values` (
  `audit_id` bigint(20) NOT NULL,
  `value_path` varchar(512) NOT NULL,
  `value_type` varchar(16) NOT NULL,
  `string_value` longtext DEFAULT NULL,
  `number_value` varchar(128) DEFAULT NULL,
  `bool_value` tinyint(3) unsigned DEFAULT NULL,
  PRIMARY KEY (`audit_id`,`value_path`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_packet_sequence` (
  `character_id` bigint(20) NOT NULL,
  `seq_index` bigint(20) NOT NULL,
  `command` bigint(20) NOT NULL,
  `noti_type` bigint(20) NOT NULL,
  `kind` bigint(20) NOT NULL DEFAULT 0,
  `item_list_type` bigint(20) NOT NULL DEFAULT -1,
  `occurrence_index` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`seq_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_packet_templates` (
  `character_id` bigint(20) NOT NULL,
  `command` bigint(20) NOT NULL,
  `noti_type` bigint(20) NOT NULL,
  `occurrence_index` bigint(20) NOT NULL DEFAULT 0,
  `body` longblob NOT NULL,
  `body_length` bigint(20) NOT NULL,
  PRIMARY KEY (`character_id`,`command`,`noti_type`,`occurrence_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_legacy_unequipped_entries` (
  `character_id` bigint(20) NOT NULL,
  `item_template_id` bigint(20) NOT NULL,
  `raw_entry` longblob NOT NULL,
  PRIMARY KEY (`character_id`,`item_template_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_mail_attachment_extra` (
  `character_id` varchar(128) NOT NULL,
  `mail_id` varchar(128) NOT NULL,
  `attachment_index` int(11) NOT NULL,
  `extra_key` varchar(128) NOT NULL,
  `extra_value` longtext NOT NULL,
  PRIMARY KEY (`character_id`,`mail_id`,`attachment_index`,`extra_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_mail_attachments` (
  `character_id` varchar(128) NOT NULL,
  `mail_id` varchar(128) NOT NULL,
  `attachment_index` int(11) NOT NULL,
  `item_id` bigint(20) NOT NULL DEFAULT 0,
  `item_count` bigint(20) NOT NULL DEFAULT 0,
  `bind_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `expire_at` datetime(6) DEFAULT NULL,
  `raw_entry` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`,`mail_id`,`attachment_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_mail_metadata` (
  `character_id` varchar(128) NOT NULL,
  `mail_id` varchar(128) NOT NULL,
  `metadata_key` varchar(128) NOT NULL,
  `metadata_value` longtext NOT NULL,
  PRIMARY KEY (`character_id`,`mail_id`,`metadata_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_mailboxes` (
  `character_id` varchar(128) NOT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_mails` (
  `character_id` varchar(128) NOT NULL,
  `mail_id` varchar(128) NOT NULL,
  `sender_character_id` varchar(128) NOT NULL DEFAULT '',
  `sender_name` varchar(128) NOT NULL DEFAULT '',
  `recipient_character_id` varchar(128) NOT NULL DEFAULT '',
  `recipient_name` varchar(128) NOT NULL DEFAULT '',
  `title` varchar(255) NOT NULL DEFAULT '',
  `body` longtext NOT NULL,
  `gold` bigint(20) NOT NULL DEFAULT 0,
  `read_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `claimed_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `deleted_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `created_at` datetime(6) DEFAULT NULL,
  `expire_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`,`mail_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_packet_template_metadata` (
  `template_id` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `entry_value` longtext NOT NULL,
  PRIMARY KEY (`template_id`,`entry_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_packet_templates` (
  `template_id` varchar(128) NOT NULL,
  `name` varchar(128) NOT NULL DEFAULT '',
  `body` longblob DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`template_id`),
  KEY `idx_dnf_packet_templates_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_pet_artifact_extra` (
  `character_id` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `extra_key` varchar(128) NOT NULL,
  `extra_value` longtext NOT NULL,
  PRIMARY KEY (`character_id`,`entry_key`,`extra_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_pet_artifacts` (
  `character_id` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `item_id` bigint(20) NOT NULL DEFAULT 0,
  `item_count` bigint(20) NOT NULL DEFAULT 0,
  `bind_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `expire_at` datetime(6) DEFAULT NULL,
  `raw_entry` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`,`entry_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_pet_clear_tokens` (
  `character_id` varchar(128) NOT NULL,
  `pet_key` varchar(128) NOT NULL,
  `token_order` int(11) NOT NULL,
  `token` varchar(255) NOT NULL,
  `applied` tinyint(3) unsigned NOT NULL DEFAULT 1,
  PRIMARY KEY (`character_id`,`pet_key`,`token`),
  UNIQUE KEY `uk_dnf_pet_clear_token_order` (`character_id`,`pet_key`,`token_order`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_pet_entries` (
  `character_id` varchar(128) NOT NULL,
  `pet_key` varchar(128) NOT NULL,
  `creature_key` bigint(20) unsigned NOT NULL DEFAULT 0,
  `item_id` bigint(20) NOT NULL DEFAULT 0,
  `source_list_type` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `source_slot_index` smallint(6) NOT NULL DEFAULT 0,
  `pet_name` varchar(128) NOT NULL DEFAULT '',
  `name_raw` longblob DEFAULT NULL,
  `satiety` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `satiety_micros` bigint(20) NOT NULL DEFAULT 0,
  `mode_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `mode1_field_0a` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `mode1_field_0b` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `pet_level` bigint(20) NOT NULL DEFAULT 0,
  `pet_exp` bigint(20) NOT NULL DEFAULT 0,
  `tail_flag` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `raw_entry` longblob DEFAULT NULL,
  PRIMARY KEY (`character_id`,`pet_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_pet_entry_extra` (
  `character_id` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `extra_key` varchar(128) NOT NULL,
  `extra_value` longtext NOT NULL,
  PRIMARY KEY (`character_id`,`entry_key`,`extra_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_pets` (
  `character_id` varchar(128) NOT NULL,
  `equipped_key` varchar(128) NOT NULL DEFAULT '',
  `town_display` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_quest_state_extra` (
  `character_id` varchar(128) NOT NULL,
  `state_group` varchar(16) NOT NULL,
  `quest_id` bigint(20) NOT NULL,
  `extra_key` varchar(128) NOT NULL,
  `extra_value` longtext NOT NULL,
  PRIMARY KEY (`character_id`,`state_group`,`quest_id`,`extra_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_quest_states` (
  `character_id` varchar(128) NOT NULL,
  `state_group` varchar(16) NOT NULL,
  `quest_id` bigint(20) NOT NULL,
  `status` varchar(64) NOT NULL DEFAULT '',
  `trigger_type` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `progress_value` bigint(20) NOT NULL DEFAULT 0,
  `reward_select_index` bigint(20) NOT NULL DEFAULT 0,
  `multiplier` bigint(20) NOT NULL DEFAULT 0,
  `state_updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`,`state_group`,`quest_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_quests` (
  `character_id` varchar(128) NOT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_setting_values` (
  `scope` varchar(128) NOT NULL,
  `entry_key` varchar(128) NOT NULL,
  `entry_value` longtext NOT NULL,
  PRIMARY KEY (`scope`,`entry_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_settings` (
  `scope` varchar(128) NOT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`scope`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_skill_cooldowns` (
  `character_id` varchar(128) NOT NULL,
  `skill_id` bigint(20) NOT NULL,
  `expires_at` datetime(6) NOT NULL,
  PRIMARY KEY (`character_id`,`skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_skill_layouts` (
  `character_id` varchar(128) NOT NULL,
  `tree_id` int(11) NOT NULL,
  `slot_index` int(11) NOT NULL,
  `skill_id` int(10) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`tree_id`,`slot_index`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_skill_states` (
  `character_id` varchar(128) NOT NULL,
  `skill_id` bigint(20) NOT NULL,
  `skill_level` int(11) NOT NULL DEFAULT 0,
  `enabled` tinyint(3) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_skill_tree_points` (
  `character_id` varchar(128) NOT NULL,
  `tree_id` int(11) NOT NULL,
  `total_sp` int(11) NOT NULL DEFAULT 0,
  `remaining_sp` int(11) NOT NULL DEFAULT 0,
  `total_tp` int(11) NOT NULL DEFAULT 0,
  `remaining_tp` int(11) NOT NULL DEFAULT 0,
  `synced_level` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`tree_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_skill_tree_profile_states` (
  `character_id` varchar(128) NOT NULL,
  `skill_id` bigint(20) NOT NULL,
  `skill_level` int(11) NOT NULL DEFAULT 0,
  `enabled` tinyint(3) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_skill_tree_profiles` (
  `character_id` varchar(128) NOT NULL,
  `expanded` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `active_tree` tinyint(3) unsigned NOT NULL DEFAULT 0,
  `total_sp` int(11) NOT NULL DEFAULT 0,
  `remaining_sp` int(11) NOT NULL DEFAULT 0,
  `total_tp` int(11) NOT NULL DEFAULT 0,
  `remaining_tp` int(11) NOT NULL DEFAULT 0,
  `synced_level` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_skill_tree_states` (
  `character_id` varchar(128) NOT NULL,
  `tree_id` int(11) NOT NULL,
  `skill_id` bigint(20) NOT NULL,
  `skill_level` int(11) NOT NULL DEFAULT 0,
  `enabled` tinyint(3) unsigned NOT NULL DEFAULT 0,
  PRIMARY KEY (`character_id`,`tree_id`,`skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_skills` (
  `character_id` varchar(128) NOT NULL,
  `total_sp` int(11) NOT NULL DEFAULT 0,
  `remaining_sp` int(11) NOT NULL DEFAULT 0,
  `total_tp` int(11) NOT NULL DEFAULT 0,
  `remaining_tp` int(11) NOT NULL DEFAULT 0,
  `synced_level` int(11) NOT NULL DEFAULT 0,
  `updated_at` datetime(6) DEFAULT NULL,
  PRIMARY KEY (`character_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8mb4 */;
CREATE TABLE `dnf_system_mail_deliveries` (
  `character_id` varchar(128) NOT NULL,
  `completion_key` varbinary(512) NOT NULL,
  `mail_id` varchar(128) NOT NULL,
  `payload_sha256` char(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  PRIMARY KEY (`character_id`,`completion_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*M!100616 SET NOTE_VERBOSITY=@OLD_NOTE_VERBOSITY */;

