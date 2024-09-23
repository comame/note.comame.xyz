
CREATE TABLE `nt_post` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `url_key` varchar(32) NOT NULL,
  `created_datetime` datetime NOT NULL,
  `updated_datetime` datetime NOT NULL,
  `title` text NOT NULL,
  `text` text NOT NULL,
  `visibility` int NOT NULL COMMENT '0=非公開, 1=限定公開, 2=全体公開',
  PRIMARY KEY (`id`),
  UNIQUE KEY `url_key` (`url_key`),
  KEY `visibility` (`visibility`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE `nt_post_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `post_id` int unsigned NOT NULL COMMENT 'nt_post.id',
  `url_key` varchar(32) NOT NULL,
  `created_datetime` datetime NOT NULL,
  `updated_datetime` datetime NOT NULL,
  `text` text NOT NULL,
  `visibility` int NOT NULL,
  PRIMARY KEY (`id`),
  KEY `post_id` (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='nt_postのログテーブル';

