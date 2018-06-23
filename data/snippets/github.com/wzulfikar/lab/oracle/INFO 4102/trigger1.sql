SET SERVEROUTPUT ON

-- CREATE TABLE test (test1 CHAR(2));

CREATE OR REPLACE TRIGGER test_var
	BEFORE INSERT OR UPDATE OR DELETE ON test -- test is table name
		FOR EACH ROW
		BEGIN
		DBMS_OUTPUT.PUT_LINE(':new.test: ' || :new.test1);
		DBMS_OUTPUT.PUT_LINE(':old.test: ' || :old.test1);

		END;
	/
