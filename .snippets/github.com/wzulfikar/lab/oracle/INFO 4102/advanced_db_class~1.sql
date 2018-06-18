select s.student_id, c.course_no from STUDENT s 
join enrollment e on s.student_id = e.student_id 
join course_section c on e.c_sec_id = c.c_sec_id
where c.course_no = 'MIS 101' and c.course_no != 'MIS 451';

select databases from dual;

-- Practice 2.1
-- insert classes conducted in BUS and CR to respective tables `Class_Bus` and `Class_CR` with ONE INSERT STATEMENT

    -- delete table
    DROP TABLE Class_Bus;DROP TABLE Class_CR;
    
    -- insert 
    INSERT ALL