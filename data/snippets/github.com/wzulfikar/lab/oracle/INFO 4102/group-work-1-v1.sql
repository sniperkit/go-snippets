SET SERVEROUTPUT ON
SET VERIFY OFF

/**
 * Get detail of employee 
 * by its employee number.
 *
 * Raise exception if input is not a valid number
 */

ACCEPT input PROMPT 'Enter employee number: '

CREATE OR REPLACE FUNCTION getMemberCount(emp_no IN INTEGER) 
   RETURN NUMBER 
   IS memberCount NUMBER(2);
   BEGIN 
      SELECT COUNT(EMPNO)
      INTO memberCount 
      FROM EMPLOYEE
      WHERE MGR = emp_no;
      
      RETURN(memberCount); 
    END;
/

CREATE OR REPLACE PROCEDURE displayMemberInfo(employee_name EMPLOYEE.ENAME%TYPE, member_count integer) AS
    BEGIN
    if member_count = 0 THEN
        DBMS_OUTPUT.PUT_LINE('- ' || employee_name || ' is not manager');
    ELSIF member_count = 1 THEN
        DBMS_OUTPUT.PUT_LINE('- ' || employee_name || ' is managing 1 person');
    ELSE
        DBMS_OUTPUT.PUT_LINE('- ' || employee_name || ' is managing ' || member_count || ' person');
    END IF;
    END;
/

CREATE OR REPLACE PROCEDURE displayManagerInfo(employee_name EMPLOYEE.ENAME%TYPE, managerId EMPLOYEE.MGR%TYPE) AS
    BEGIN
    if managerId IS NULL THEN
        DBMS_OUTPUT.PUT_LINE('- ' || employee_name || ' is not being managed by anyone');
    ELSE
        DBMS_OUTPUT.PUT_LINE('- ' || employee_name || ' is under supervision of employee #' || managerId);
    END IF;
    END;
/

DECLARE
	CURSOR employee_cursor IS
	SELECT 
        E.EMPNO EMPNO, 
        E.ENAME ENAME, 
        E.JOB JOB, 
        E.SAL SAL, 
        B.BRANCH_NAME BRANCH_NAME, 
        B.BRANCH_LOCATION BRANCH_LOCATION,
        E.HIREDATE,
        E.MGR MGR
	FROM EMPLOYEE E
    JOIN BRANCH B
    ON B.BRANCH_NUM = E.BRANCH_NUM
	WHERE EMPNO = '&input';

	employee employee_cursor%ROWTYPE;
	EMPLOYEE_NOT_FOUND_EXCEPTION EXCEPTION;
    
	liner           varchar2(30)  := '----------------------------';    
	rpad_size       integer       := 13;
	months_worked   integer;
	years_worked    decimal;
	user_input      varchar2(20)  := '&input';
	
BEGIN
	OPEN employee_cursor;
    
	LOOP
	FETCH employee_cursor INTO employee;
    
    -- display descriptive message if no records found
    IF employee_cursor%NOTFOUND AND employee_cursor%ROWCOUNT = 0 THEN
        raise EMPLOYEE_NOT_FOUND_EXCEPTION;
    END IF;
    
	EXIT WHEN employee_cursor%NOTFOUND; -- use explicit cursor `employee_cursor`
    
    -- calculate months and years worked
    months_worked := ROUND(MONTHS_BETWEEN(CURRENT_DATE, employee.HIREDATE));
    years_worked  := ROUND(months_worked/12);
    
    -- craft employee info
    DBMS_OUTPUT.PUT_LINE(liner);
    DBMS_OUTPUT.PUT_LINE('Employee Info');
    DBMS_OUTPUT.PUT_LINE(liner);
    DBMS_OUTPUT.PUT_LINE(RPAD('Emp No.', rpad_size) || ' : ' || employee.EMPNO);
	DBMS_OUTPUT.PUT_LINE(RPAD('Name', rpad_size) || ' : ' || employee.ENAME);
    DBMS_OUTPUT.PUT_LINE(RPAD('Job', rpad_size) || ' : ' || employee.JOB);
    DBMS_OUTPUT.PUT_LINE(RPAD('Salary', rpad_size) || ' : ' || employee.SAL);
    DBMS_OUTPUT.PUT_LINE(RPAD('Hiredate', rpad_size) || ' : ' || employee.HIREDATE);
    DBMS_OUTPUT.PUT_LINE(RPAD('Months Worked', rpad_size) || ' : ' || months_worked || ' (' || years_worked || 'yrs)');
    DBMS_OUTPUT.PUT_LINE(RPAD('Branch', rpad_size) || ' : ' || employee.BRANCH_NAME);
	DBMS_OUTPUT.PUT_LINE(RPAD('Branch Loc.', rpad_size) || ' : ' || employee.BRANCH_LOCATION);

    -- craft header for management info
    DBMS_OUTPUT.PUT_LINE('');
    DBMS_OUTPUT.PUT_LINE(liner);
    DBMS_OUTPUT.PUT_LINE('Management Info');
    DBMS_OUTPUT.PUT_LINE(liner);
    
    -- call procedure to display info related to management
    displayMemberInfo(employee.ENAME, getMemberCount(employee.EMPNO));
    displayManagerInfo(employee.ENAME, employee.MGR);
    
	END LOOP;

	CLOSE employee_cursor;
    
    -- exception handler begins
    EXCEPTION
       WHEN INVALID_NUMBER THEN
       DBMS_OUTPUT.PUT_LINE('Oops! Something went wrong..');
       DBMS_OUTPUT.PUT_LINE('- "' || user_input || '" is not valid employee number');
       DBMS_OUTPUT.PUT_LINE('- ' || 'Employee number must be integer');
       WHEN EMPLOYEE_NOT_FOUND_EXCEPTION THEN
       DBMS_OUTPUT.PUT_LINE('Not found: employee with employee number ' || user_input || ' does not exist.');
END;
/
