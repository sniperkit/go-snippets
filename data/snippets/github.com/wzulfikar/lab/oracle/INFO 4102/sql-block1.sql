SET SERVEROUTPUT ON
SET VERIFY OFF

ACCEPT bcode prompt 'Enter the book code: '

DECLARE
    v_title book.title%TYPE;
BEGIN
    SELECT title INTO v_title
    FROM book
    WHERE book_code = '&bcode';
    DBMS_OUTPUT.PUT_LINE
    ('The title of the books is ' || v_title);
END;
/

SET SERVEROUTPUT OFF
SET VERIFY ONmove
