CREATE TABLE `osu!`.`scores` (`scoreid` INT NOT NULL , `mapchecksum` TEXT NOT NULL , `username` TEXT NOT NULL , `OnlineScoreChecksum` TEXT NOT NULL , `Count300` INT NOT NULL , `Count100` INT NOT NULL , `Count50` INT NOT NULL , `CountGeki` INT NOT NULL , `CountKatu` INT NOT NULL , `CountMiss` INT NOT NULL , `TotalScore` INT NOT NULL , `MaxCombo` INT NOT NULL , `Perfect` TEXT NOT NULL , `Ranking` TEXT NOT NULL , `EnabledMods` TEXT NOT NULL , `Pass` TEXT NOT NULL, `Accuracy` FLOAT NOT NULL ) ENGINE = InnoDB; 