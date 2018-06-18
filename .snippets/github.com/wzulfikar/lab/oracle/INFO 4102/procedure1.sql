CREATE OR REPLACE PROCEDURE Emp_Proc
	(peno IN employee.empno%TYPE
	 pename OUT employee.ename%TYPE
	 pjob OUT employee.job%TYPE
	 pbranch_name OUT employee.branch_name%TYPE
	 pbranch_location OUT employee.branch_location%TYPE)
AS
	SELECT ename, job, branch_name, branch_location
	INTO pename, pjob, pbranch_name, pbranch_location
	FROM employee NATURAL JOIN branch
	WHERE empno = peno;
BEGIN   -- Find employee details

END Emp_Proc;

ACCEPT eno PROMPT 'Enter the employee number: ' 

-- ANONYMOUS BLOCK
DECLARE
BEGIN
v_name v_job v_bname v_blocation
employee.ename%TYPE; employee.job%TYPE; branch.branch_name%TYPE; branch.branch_location%TYPE;
END;

Emp_Proc(&eno, v_name, v_job, v_bname, v_blocation);
DBMS_OUTPUT.put_line('Name: ' || v_name);
DBMS_OUTPUT.put_line('Job: ' || v_job);
DBMS_OUTPUT.put_line('Branch Name: ' || v_bname);
DBMS_OUTPUT.put_line('Branch Location: ' || v_blocation);
END;
/