
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
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
