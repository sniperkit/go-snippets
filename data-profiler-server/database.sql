CREATE TABLE IF NOT EXISTS `datasets` (
		`id` INTEGER PRIMARY KEY AUTOINCREMENT,
		`path` VARCHAR(200),
		`name` VARCHAR(200),
		`description` VARCHAR(2000)
);


CREATE TABLE IF NOT EXISTS `matrices` (
		`id` INTEGER PRIMARY KEY AUTOINCREMENT,
		`path` VARCHAR(500),
		`filename` VARCHAR(500),
		`configuration` VARCHAR(2000),
		`datasetid` INTEGER,
		`estimatorpath` VARCHAR(500),
		FOREIGN KEY(datasetid) REFERENCES datasets(id)
);

CREATE TABLE IF NOT EXISTS `coordinates` (
		`id` INTEGER PRIMARY KEY AUTOINCREMENT,
		`path` VARCHAR(500),
		`filename` VARCHAR(500),
		`k` VARCHAR(500),
		`gof` VARCHAR(2000),
		`stress` VARCHAR(2000),
		`matrixid` INTEGER,
		FOREIGN KEY(matrixid) REFERENCES matrices(id)
);

CREATE TABLE IF NOT EXISTS `operators` (
		`id` INTEGER PRIMARY KEY AUTOINCREMENT,
		`name` VARCHAR(500),
		`description` VARCHAR(500),
		`path` VARCHAR(500),
		`datasetid` INTEGER,
		`scoresfile` VARCHAR(500),
		FOREIGN KEY(datasetid) REFERENCES datasets(id)
);

CREATE TABLE IF NOT EXISTS `models` (
		`id` INTEGER PRIMARY KEY AUTOINCREMENT,
		`coordinatesid` INTEGER,
		`operatorid` INTEGER,
		`datasetid` INTEGER,
		`samplingrate` DECIMAL,
		`configuration` VARCHAR(2000),
		`errors` VARCHAR(2000),
		`samplespath` VARCHAR(500), 
		`appxvaluespath` VARCHAR(500), 
		FOREIGN KEY(datasetid) REFERENCES datasets(id),
		FOREIGN KEY(coordinatesid) REFERENCES coordinates(id),
		FOREIGN KEY(operatorid) REFERENCES operators(id)
);
