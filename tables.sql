
CREATE TABLE `nt_post` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `parent` int unsigned COMMENT '0の場合、最上位とする',
  `url_key` varchar(32) NOT NULL,
  `created_datetime` datetime NOT NULL,
  `updated_datetime` datetime NOT NULL,
  `title` text NOT NULL,
  `text` text NOT NULL,
  `permission` varchar(16) NOT NULL,
  `permission_inherited` tinyint(1),
  PRIMARY KEY (`id`),
  UNIQUE KEY `url_key` (`url_key`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `nt_post_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `parent` int unsigned COMMENT '0の場合、最上位とする',
  `post_id` int unsigned NOT NULL COMMENT 'nt_post.id',
  `url_key` varchar(32) NOT NULL,
  `created_datetime` datetime NOT NULL,
  `updated_datetime` datetime NOT NULL,
  `text` text NOT NULL,
  `permission` varchar(16) NOT NULL,
  `permission_inherited` tinyint(1),
  PRIMARY KEY (`id`),
  KEY `post_id` (`post_id`)
) ENGINE=InnoDB AUTO_INCREMENT=32 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='nt_postのログテーブル';

