CREATE TABLE IF NOT EXISTS `posts` (
  `id` char(16) NOT NULL,
  `modify_token` char(32) DEFAULT NULL,
  `text_appearance` char(4) NOT NULL DEFAULT 'norm',
  `last_modified` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;
