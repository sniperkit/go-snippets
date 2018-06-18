SET SERVEROUTPUT ON
SET VERIFY OFF

/**
 * Function to calculate age. No return value.
 */
CREATE OR REPLACE FUNCTION calc_age(birthdate IN STUDENT.S_DOB%TYPE) 
   RETURN NUMBER 
   IS 
    age DECIMAL(4,1);
   BEGIN 
        age := ROUND(MONTHS_BETWEEN(SYSDATE, birthdate)/12, 1);
      RETURN(age); 
    END;
/

/**
 * Procedure with IN and OUT args
 */
CREATE OR REPLACE PROCEDURE count_courses_completed(s_id IN STUDENT.STUDENT_ID%TYPE, courses_completed OUT integer) AS
    BEGIN
    SELECT count(grade) INTO courses_completed FROM enrollment where student_id = s_id;
    END;
/

/**
 * Procedure to display enrolled courses of given student id.
 */
CREATE OR REPLACE PROCEDURE display_enrolled_courses(input ENROLLMENT.STUDENT_ID%TYPE) AS
    -- cursor declaration
	CURSOR student_cursor IS
	select 
		en.grade GRADE,
        cs.course_no COURSE_NO,
        cs.c_sec_id SECTION
    from enrollment en
	join course_section cs on en.c_sec_id = cs.c_sec_id
	where en.student_id = upper(input);
    
    student student_cursor%ROWTYPE;
    grade   varchar2(10);
    STUDENT_NOT_FOUND_EXCEPTION EXCEPTION;
    
    BEGIN
    
    DBMS_OUTPUT.PUT_LINE(RPAD('COURSE', 10) || RPAD('SECTION', 10) || RPAD('GRADE', 10));
    DBMS_OUTPUT.PUT_LINE('----------------------------------------------');

    OPEN student_cursor;
    
	LOOP
	FETCH student_cursor INTO student;
    
    -- display descriptive message if no records found
    IF student_cursor%NOTFOUND AND student_cursor%ROWCOUNT = 0 THEN
        raise STUDENT_NOT_FOUND_EXCEPTION;
    END IF;
    
    EXIT WHEN student_cursor%NOTFOUND; -- use explicit cursor `employee_cursor`
    
    grade := student.GRADE;
    
    IF (grade IS NULL) THEN
        grade := 'INCOMPLETE';
    END IF;

    DBMS_OUTPUT.PUT_LINE(RPAD(student.COURSE_NO, 10) || RPAD(student.SECTION, 10) || RPAD(grade, 10));
    
    END LOOP;
    CLOSE student_cursor;
    
    -- start exception handler
    EXCEPTION
       -- handle user-defined exception
       WHEN STUDENT_NOT_FOUND_EXCEPTION THEN
	       DBMS_OUTPUT.PUT_LINE('Not found: student id "' || input || '" has not enrolled to any courses');
END display_enrolled_courses;
/


-- ANONYMOUS BLOCK STARTS HERE --

ACCEPT input PROMPT 'Enter student id: '

DECLARE
	-- cursor declaration
	CURSOR student_cursor IS
	select 
		s.STUDENT_ID STUDENT_ID, 
        s.S_FIRST S_FIRST, 
        s.S_LAST S_LAST, 
        s.S_PHONE S_PHONE,
        s.S_DOB S_DOB,
        s.S_ADDRESS S_ADDRESS
    from student s
	where s.STUDENT_ID = upper('&input'); -- support uppercase and lowercase input

	student student_cursor%ROWTYPE;

	-- declare user-defined exception
	STUDENT_NOT_FOUND_EXCEPTION EXCEPTION;
    
	liner           varchar2(30)  := '----------------------------';    
	name           	varchar2(50)  := '----------------------------';    
	rpad_size       integer       := 17;
	user_input      varchar2(20)  := '&input';
	courses_completed integer(3);
BEGIN
	OPEN student_cursor;
    
	LOOP
	FETCH student_cursor INTO student;
    
    -- display descriptive message if no records found
    IF student_cursor%NOTFOUND AND student_cursor%ROWCOUNT = 0 THEN
        raise STUDENT_NOT_FOUND_EXCEPTION;
    END IF;
    
	EXIT WHEN student_cursor%NOTFOUND; -- use explicit cursor `employee_cursor`
    
    -- call procedure
    count_courses_completed(student.STUDENT_ID, courses_completed);    
    
    -- craft employee info
    DBMS_OUTPUT.PUT_LINE(liner);
    DBMS_OUTPUT.PUT_LINE('Student Info');
    DBMS_OUTPUT.PUT_LINE(liner);
    DBMS_OUTPUT.PUT_LINE(RPAD('Student Id.', rpad_size) || ' : ' || student.STUDENT_ID);
	DBMS_OUTPUT.PUT_LINE(RPAD('Name', rpad_size) || ' : ' || student.S_FIRST || ' ' || student.S_LAST);
    DBMS_OUTPUT.PUT_LINE(RPAD('Address', rpad_size) || ' : ' || student.S_ADDRESS);
    DBMS_OUTPUT.PUT_LINE(RPAD('Phone', rpad_size) || ' : ' || student.S_PHONE);
    DBMS_OUTPUT.PUT_LINE(RPAD('Birthdate', rpad_size) || ' : ' || student.S_DOB);
    DBMS_OUTPUT.PUT_LINE(RPAD('Age', rpad_size) || ' : ' || calc_age(student.S_DOB) || ' year-old');
    DBMS_OUTPUT.PUT_LINE(RPAD('Courses Completed', rpad_size) || ' : ' || courses_completed);

	END LOOP;
    
	CLOSE student_cursor;
    
    DBMS_OUTPUT.NEW_LINE;
    DBMS_OUTPUT.NEW_LINE;
    DBMS_OUTPUT.PUT_LINE('####### ENROLLMENT DETAIL #######');
    DBMS_OUTPUT.NEW_LINE;
    DBMS_OUTPUT.NEW_LINE;    
    
    display_enrolled_courses(student.STUDENT_ID);
    
    -- start exception handler
    EXCEPTION
       -- handle user-defined exception
       WHEN STUDENT_NOT_FOUND_EXCEPTION THEN
	       DBMS_OUTPUT.PUT_LINE('Not found: student with student number "' || user_input || '" does not exist.');
END;
/
-- END OF ANONYMOUS BLOCK
