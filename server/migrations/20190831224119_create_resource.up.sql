CREATE TABLE resources(
  id          	INTEGER PRIMARY KEY,
	title       	TEXT NOT NULL,
	description 	TEXT,
	url         	TEXT NOT NULL,
	referrer    	TEXT,
	credit      	TEXT,
	published_at	datetime,
	type        	TEXT NOT NULL,
	-- tags        	TEXT,
	created_at  	datetime DEFAULT CURRENT_TIMESTAMP,
	updated_at  	datetime DEFAULT CURRENT_TIMESTAMP
);
