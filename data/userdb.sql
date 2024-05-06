\! echo 'Creating Cloudtacts user database and table...'

CREATE DATABASE IF NOT EXISTS cloudtacts;

CREATE TABLE IF NOT EXISTS cloudtacts.user
(
	ctuser	VARCHAR(20) NOT NULL,
	ctpass	CHAR(66) NOT NULL,
	ctprof	VARCHAR(20) NOT NULL,
	uemail	VARCHAR(50) NOT NULL,
	ctppic	VARCHAR(52),
	atoken	CHAR(20),
	llogin	DATETIME,
	uvalid	DATETIME,
	CONSTRAINT PRIMARY KEY (ctuser, ctprof, uemail)
) ENGINE=InnoDB;

COMMIT;
