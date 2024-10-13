CREATE TABLE `admins` (
  `userid` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
ALTER TABLE `admins`
  ADD PRIMARY KEY (`userid`);
COMMIT;
