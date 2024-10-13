CREATE TABLE `users` (
  `userid` int(11) NOT NULL,
  `username` text NOT NULL,
  `password` text NOT NULL,
  `ranked_score` bigint(20) NOT NULL,
  `accuracy` float NOT NULL,
  `playcount` int(11) NOT NULL,
  `total_score` bigint(20) NOT NULL,
  `rank` int(11) NOT NULL,
  `lastonline` text NOT NULL,
  `joindate` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
INSERT INTO `users`(`userid`, `username`, `password`, `ranked_score`, `accuracy`, `playcount`, `total_score`, `rank`, `lastonline`) VALUES ('1','BanchoBot','nononoononononnoon','0','0','0','0','0','0','even server does not now')