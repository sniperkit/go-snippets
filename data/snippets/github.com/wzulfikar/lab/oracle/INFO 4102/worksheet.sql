-- Class #4: Wed, 20 Sep 2017 at 14:23:36 MYT

# exercise 1
SET LINESIZE 75
SET PAGESIZE 50

TTITLE 'LIST OF BOOKS || TYPE: PSY' -- `|` means new line
COLUMN bcode HEADING 'BOOK|CODE' FORMAT A10
COLUMN title HEADING 'BOOK|TITLE' FORMAT A10
COLUMN type HEADING 'BOOK|TYPE' FORMAT A10
COLUMN price HEADING 'PRICE' FORMAT $999.99

SELECT book_code bcode , title, type, price
FROM book
WHERE type = 'PSY';

SET LINESIZE 80
TTITLE OFF
CLEAR COLUMNS

# exercise 2
SET LINESIZE 75
SET PAGESIZE 50

TTITLE 'LIST OF BOOKS' -- `|` means new line
BREAK ON type SKIP 2
COLUMN bcode HEADING 'BOOK|CODE' FORMAT A10
COLUMN title HEADING 'BOOK|TITLE' FORMAT A10
COLUMN type HEADING 'BOOK|TYPE' FORMAT A10
COLUMN price HEADING 'PRICE' FORMAT $999.99

SELECT type, book_code bcode, title, price
FROM book
WHERE type IN ('POE', 'ART');

SET LINESIZE 80
TTITLE OFF
CLEAR COLUMNS

# exercise 3
SET LINESIZE 100
SET PAGESIZE 50

TTITLE 'LIST OF BOOKS' -- `|` means new line
BREAK ON city SKIP 2 ON publisher_code SKIP 2
COLUMN city HEADING 'CITY' FORMAT A20
COLUMN publisher_code HEADING 'PUB' FORMAT A5
COLUMN title HEADING 'TITLE' FORMAT A40

select p.city city, p.publisher_code, b.title
from book b 
join publisher p 
on b.publisher_code = p.publisher_code
where p.city IN ('Boston', 'Boulder CO', 'Westport CT')
order by p.city;

SET LINESIZE 80
TTITLE OFF
CLEAR COLUMNS

-- Class #5: Mon, 25 Sep 2017 at 14:10:36 MYT
-- see `sql-block1.sql` & `sql-block2.sql`

-- Class #6: Wed, 27 Sep 2017 at 14:33:02 MYT
-- see `cursor1.sql`