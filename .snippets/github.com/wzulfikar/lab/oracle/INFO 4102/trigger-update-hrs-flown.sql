SET SERVEROUTPUT ON

CREATE OR REPLACE TRIGGER hrs_flown
	BEFORE INSERT OR UPDATE OR DELETE ON charter FOR EACH ROW
		BEGIN 
		IF DELETING THEN
		UPDATE aircraft SET ac_hours_flown = ac_hours_flown - :old.char_hours_flown WHERE ac_number = :old.ac_number;
		DBMS_OUTPUT.PUT_LINE('Charter trip deleted and aircraft hours flown updated.');
		ELSE
		UPDATE aircraft SET ac_hours_flown = ac_hours_flown + :new.char_hours_flown WHERE ac_number = :new.ac_number;
		DBMS_OUTPUT.PUT_LINE('Charter trip added/updated and aircraft hours flown updated.');
		END IF;
		END;
	/
