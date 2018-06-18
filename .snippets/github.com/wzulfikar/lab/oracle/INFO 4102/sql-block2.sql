-- Mon, 25 Sep 2017 at 14:59:32 MYT
SET SERVEROUTPUT ON
SET VERIFY OFF

ACCEPT bcode prompt 'Enter the book code: '

DECLARE
    v_title book.title%TYPE;
    v_price book.price%TYPE;
    v_price_discounted book.price%TYPE;
BEGIN
    SELECT title, price, 
    	CASE WHEN TYPE = 'ART' THEN
    		ROUND(price / (1-0.10), 2)
    	WHEN TYPE = 'FIC' THEN
    		ROUND(price * (1-0.15), 2)
    	ELSE
    		price
    	END
    INTO v_title, v_price, v_price_discounted
    FROM book
    WHERE book_code = '&bcode';

    DBMS_OUTPUT.PUT_LINE
    	('The title of the books is ' || v_title);
    DBMS_OUTPUT.PUT_LINE
    	('Original Price: RM ' || v_price);
    DBMS_OUTPUT.PUT_LINE
    	('The price after discount: RM ' || v_price_discounted);
END;
/

SET SERVEROUTPUT OFF
SET VERIFY ON
