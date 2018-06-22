DROP TABLE IF EXISTS Lookup;
DROP TABLE IF EXISTS User;
DROP TABLE IF EXISTS Post;
DROP TABLE IF EXISTS Comment;
DROP TABLE IF EXISTS Tag;

CREATE TABLE Lookup ( id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name VARCHAR(128) NOT NULL, code INTEGER NOT NULL, type VARCHAR(128) NOT NULL, position INTEGER NOT NULL );
CREATE TABLE User ( id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, username VARCHAR(128) NOT NULL, password VARCHAR(128) NOT NULL, salt VARCHAR(128) NOT NULL, email VARCHAR(128) NOT NULL, profile TEXT );
CREATE TABLE Post ( id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, title VARCHAR(128) NOT NULL, content TEXT NOT NULL, tags TEXT, status INTEGER NOT NULL, create_time INTEGER, update_time INTEGER, author_id INTEGER NOT NULL, CONSTRAINT FK_post_author FOREIGN KEY (author_id) REFERENCES User (id) ON DELETE CASCADE ON UPDATE RESTRICT );
CREATE TABLE Comment ( id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, content TEXT NOT NULL, status INTEGER NOT NULL, create_time INTEGER, author VARCHAR(128) NOT NULL, email VARCHAR(128) NOT NULL, url VARCHAR(128), post_id INTEGER NOT NULL, CONSTRAINT FK_comment_post FOREIGN KEY (post_id) REFERENCES Post (id) ON DELETE CASCADE ON UPDATE RESTRICT );
CREATE TABLE Tag ( id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name VARCHAR(128) NOT NULL, frequency INTEGER DEFAULT 1 );

INSERT INTO Lookup (name, type, code, position) VALUES ('Draft', 'PostStatus', 1, 1);
INSERT INTO Lookup (name, type, code, position) VALUES ('Published', 'PostStatus', 2, 2);
INSERT INTO Lookup (name, type, code, position) VALUES ('Archived', 'PostStatus', 3, 3);
INSERT INTO Lookup (name, type, code, position) VALUES ('Pending Approval', 'CommentStatus', 1, 1);
INSERT INTO Lookup (name, type, code, position) VALUES ('Approved', 'CommentStatus', 2, 2);

INSERT INTO User (username, password, salt, email) VALUES ('demo','2e5c7db760a33498023813489cfadc0b','28b206548469ce62182048fd9cf91760','webmaster@example.com');
INSERT INTO Post (title, content, status, create_time, update_time, author_id, tags) VALUES ('Welcome!','<p>Rerum velit quos est <strong>similique</strong>. Consectetur tempora eos ullam velit nobis sit debitis. Magni explicabo omnis delectus labore vel recusandae.</p><p>Aut a minus laboriosam harum placeat quas minima fuga. Quos nulla fuga quam officia tempore. Rerum occaecati ut eum et tempore. Nam ab repudiandae et nemo praesentium.</p><p>Cumque corporis officia occaecati ducimus sequi laborum omnis ut. Nam aspernatur veniam fugit. Nihil eum libero ea dolorum ducimus impedit sed. Quidem inventore porro corporis debitis eum in. Nesciunt unde est est qui nulla. Esse sunt placeat molestiae molestiae sed quia. Sunt qui quidem quos velit reprehenderit quos blanditiis ducimus. Sint et molestiae maxime ut consequatur minima. Quaerat rem voluptates voluptatem quos. Corporis perferendis in provident iure. Commodi odit exercitationem excepturi et deserunt qui.</p><p>Optio iste necessitatibus velit non. Neque sed occaecati culpa porro culpa. Quia quam in molestias ratione et necessitatibus consequatur. Est est tempora consequatur voluptatem vel. Mollitia tenetur non quis omnis perspiciatis deserunt sed necessitatibus. Ad rerum reiciendis sunt aspernatur.</p><p>Est ullam ut magni aspernatur. Eum et sed tempore modi.</p><p>Earum aperiam sit neque quo laborum suscipit unde. Expedita nostrum itaque non non adipisci. Ut delectus quis delectus est at sint. Iste hic qui ea eaque eaque sed id. Hic placeat rerum numquam id velit deleniti voluptatem. Illum adipisci voluptas adipisci ut alias. Earum exercitationem iste quidem eveniet aliquid hic reiciendis. Exercitationem est sunt in minima consequuntur. Aut quaerat libero dolorem.</p>',2,1230952187,1230952187,1,'announce, blog');
INSERT INTO Post (title, content, status, create_time, update_time, author_id, tags) VALUES ('A Test Post', '<p>Aliquam magni <em>tempora</em> quas enim. Perspiciatis libero corporis sunt eum nam. Molestias est sunt molestiae natus.</p><p>Blanditiis dignissimos autem culpa itaque. Explicabo perferendis ullam officia ut quia nemo. Eaque perspiciatis perspiciatis est hic non ullam et. Expedita exercitationem enim sit ut dolore.</p><h2>Sed in sunt officia blanditiis ipsam maiores perspiciatis amet</h2><p>Vero fugiat facere officiis aut quis rerum velit. Autem eius sint ullam. Nemo sunt molestiae nulla accusantium est voluptatem voluptas sed. In blanditiis neque libero voluptatem praesentium occaecati nulla libero. Perspiciatis eos voluptatem facere voluptatibus. Explicabo quo eveniet nihil culpa. Qui eos officia consequuntur sed esse praesentium dolorum. Eius perferendis qui quia autem nostrum sed. Illum in ex excepturi voluptas. Qui veniam sit alias delectus nihil. Impedit est ut alias illum repellendus qui.</p><p>Veniam est aperiam quisquam soluta. Magni blanditiis praesentium sed similique velit ipsam consequatur. Porro omnis magni sunt incidunt aspernatur ut.</p>', 2,1230952187,1230952187,1,'test');

INSERT INTO Comment (content, status, create_time, author, email, post_id) VALUES ('This is a test comment.', 2, 1230952187, 'Tester', 'tester@example.com', 2);

INSERT INTO Tag (name) VALUES ('announce');
INSERT INTO Tag (name) VALUES ('blog');
INSERT INTO Tag (name) VALUES ('test');
