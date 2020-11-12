CREATE TABLE spaceships(
   Id INT(11) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
   Name VARCHAR(255) NOT NULL,
   Class VARCHAR(255) NOT NULL,
   Image VARCHAR(255) NOT NULL,
   Status VARCHAR(255) NOT NULL,
   Armaments LONGTEXT NOT NULL,
   Crew INT(11) UNSIGNED,
   Value BIGINT(11) UNSIGNED
);
