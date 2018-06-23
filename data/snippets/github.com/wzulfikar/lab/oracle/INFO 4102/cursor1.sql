SET SERVEROUTPUT ON
SET VERIFY OFF

ACCEPT pcode PROMPT 'Enter the publisher code: '

DECLARE
	CURSOR c1 IS
	SELECT BOOK_CODE, TITLE, PRICE
	FROM book
	WHERE PUBLISHER_CODE = '&pcode';

	book_rec c1%ROWTYPE;

BEGIN
	OPEN c1;

	LOOP
	FETCH c1 INTO book_rec;
	EXIT WHEN c1%NOTFOUND; -- use explicit cursor `c1`

	DBMS_OUTPUT.PUT_LINE
	(book_rec.book_code || ' ' || book_rec.title || ' ' || book_rec.price);
	END LOOP;

	CLOSE c1;
END;
/