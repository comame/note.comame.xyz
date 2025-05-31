
CREATE TABLE `nt_post` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `url_key` varchar(32) NOT NULL,
  `created_datetime` datetime NOT NULL,
  `updated_datetime` datetime NOT NULL,
  `title` text NOT NULL,
  `text` text NOT NULL,
  `permission` varchar(16) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `url_key` (`url_key`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `nt_post_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `post_id` int unsigned NOT NULL COMMENT 'nt_post.id',
  `url_key` varchar(32) NOT NULL,
  `created_datetime` datetime NOT NULL,
  `updated_datetime` datetime NOT NULL,
  `text` text NOT NULL,
  `permission` varchar(16) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `post_id` (`post_id`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='nt_postのログテーブル';

