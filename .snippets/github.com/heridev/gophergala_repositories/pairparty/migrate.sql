DROP TABLE IF EXISTS user;

CREATE TABLE "user" ("id" INTEGER PRIMARY KEY  ,"username" varchar(255) NOT NULL ,"name" varchar(255)  ,"email" varchar(255) NOT NULL ,"password" varchar(255) NOT NULL ,"created_at" datetime  ,"updated_at" datetime  ,"deleted_at" datetime  );

insert into user (username, name, email, password) values('foo','Foo Bar', 'foo@bar.com', '12345');
