
drop table dependents;
drop table employees;
drop table departments;
drop table locations;
drop table countries;
drop table jobs;
drop table regions;

CREATE TABLE regions (
 region_id NUMBER PRIMARY KEY,
 region_name VARCHAR2 (25) DEFAULT NULL
);

CREATE TABLE jobs (
 job_id NUMBER PRIMARY KEY,
 job_title VARCHAR2 (35) NOT NULL,
 min_salary NUMBER (8, 2) DEFAULT NULL,
 max_salary NUMBER (8, 2) DEFAULT NULL
);
 
CREATE TABLE countries (
 country_id CHAR (2) PRIMARY KEY,
 country_name VARCHAR2 (40) DEFAULT NULL,
 region_id NUMBER NOT NULL,
 FOREIGN KEY (region_id) REFERENCES regions (region_id) ON DELETE CASCADE 
);
 
CREATE TABLE locations (
 location_id NUMBER PRIMARY KEY,
 street_address VARCHAR2 (40) DEFAULT NULL,
 postal_code VARCHAR2 (12) DEFAULT NULL,
 city VARCHAR2 (30) NOT NULL,
 state_province VARCHAR2 (25) DEFAULT NULL,
 country_id CHAR (2) NOT NULL,
 FOREIGN KEY (country_id) REFERENCES countries (country_id) ON DELETE CASCADE 
);
 
CREATE TABLE departments (
 department_id NUMBER PRIMARY KEY,
 department_name VARCHAR2 (30) NOT NULL,
 location_id NUMBER DEFAULT NULL,
 FOREIGN KEY (location_id) REFERENCES locations (location_id) ON DELETE CASCADE 
);
 
CREATE TABLE employees (
 employee_id NUMBER PRIMARY KEY,
 first_name VARCHAR2 (20) DEFAULT NULL,
 last_name VARCHAR2 (25) NOT NULL,
 email VARCHAR2 (100) NOT NULL,
 phone_number VARCHAR2 (20) DEFAULT NULL,
 hire_date DATE NOT NULL,
 job_id NUMBER NOT NULL,
 salary NUMBER (8, 2) NOT NULL,
 manager_id NUMBER DEFAULT NULL,
 department_id NUMBER DEFAULT NULL,
 FOREIGN KEY (job_id) REFERENCES jobs (job_id) ON DELETE CASCADE,
 FOREIGN KEY (department_id) REFERENCES departments (department_id) ON DELETE CASCADE,
 FOREIGN KEY (manager_id) REFERENCES employees (employee_id)
);
 
CREATE TABLE dependents (
 dependent_id NUMBER PRIMARY KEY,
 first_name VARCHAR2 (50) NOT NULL,
 last_name VARCHAR2 (50) NOT NULL,
 relationship VARCHAR2 (25) NOT NULL,
 employee_id NUMBER NOT NULL,
 FOREIGN KEY (employee_id) REFERENCES employees (employee_id) ON DELETE CASCADE
);

-- Insert data into the regions table.
INSERT INTO regions (region_id, region_name) VALUES  (1, 'Europe');
INSERT INTO regions (region_id, region_name) VALUES  (2, 'Americas');
INSERT INTO regions (region_id, region_name) VALUES  (3, 'Asia');
INSERT INTO regions (region_id, region_name) VALUES  (4, 'Middle East and Africa');

-- Insert data into the jobs table.
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 1, 'Public Accountant', '4200.00', '9000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 2, 'Accounting Manager', '8200.00', '16000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 3, 'Administration Assistant', '3000.00', '6000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 4, 'President', '20000.00', '40000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 5, 'Administration Vice President', '15000.00', '30000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 6, 'Accountant', '4200.00', '9000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 7, 'Finance Manager', '8200.00', '16000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 8, 'Human Resources Representative', '4000.00', '9000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 9, 'Programmer', '4000.00', '10000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 10, 'Marketing Manager', '9000.00', '15000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 11, 'Marketing Representative', '4000.00', '9000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 12, 'Public Relations Representative', '4500.00', '10500.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 13, 'Purchasing Clerk', '2500.00', '5500.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 14, 'Purchasing Manager', '8000.00', '15000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 15, 'Sales Manager', '10000.00', '20000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 16, 'Sales Representative', '6000.00', '12000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 17, 'Shipping Clerk', '2500.00', '5500.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 18, 'Stock Clerk', '2000.00', '5000.00' );
INSERT INTO jobs ( job_id, job_title, min_salary, max_salary)
VALUES ( 19, 'Stock Manager', '5500.00', '8500.00' );

-- Insert data into the countries table.
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('AR', 'Argentina', 2);
INSERT INTO countries (country_id, country_name, region_id)
VALUES('AU', 'Australia', 3);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('BE', 'Belgium', 1);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('BR', 'Brazil', 2);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('CA', 'Canada', 2);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('CH', 'Switzerland', 1);
 INSERT INTO countries (country_id, country_name, region_id)
VALUES ('CN', 'China', 3);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('DE', 'Germany', 1);
 INSERT INTO countries (country_id, country_name, region_id)
VALUES ('DK', 'Denmark', 1);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('EG', 'Egypt', 4);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('FR', 'France', 1);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('HK', 'HongKong', 3);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('IL', 'Israel', 4);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('IN', 'India', 3);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('IT', 'Italy', 1);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('JP', 'Japan', 3);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('KW', 'Kuwait', 4);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('MX', 'Mexico', 2);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('NG', 'Nigeria', 4);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('NL', 'Netherlands', 1);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('SG', 'Singapore', 3);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('UK', 'United Kingdom', 1);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ( 'US', 'United States of America', 2 );
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('ZM', 'Zambia', 4);
INSERT INTO countries (country_id, country_name, region_id)
VALUES ('ZW', 'Zimbabwe', 4);

-- Insert data into the locations table.
INSERT INTO locations ( location_id, street_address, postal_code, city, state_province,
 country_id)
VALUES ( 1400, '2014 Jabberwocky Rd', '26192', 'Southlake', 'Texas', 'US' );
 INSERT INTO locations ( location_id, street_address, postal_code, city, state_province,
 country_id)
VALUES (1500, '2011 Interiors Blvd', '99236', 'South San Francisco', 'California',
 'US' );
INSERT INTO locations ( location_id, street_address, postal_code, city, state_province,
 country_id)
VALUES  ( 1700, '2004 Charade Rd', '98199', 'Seattle', 'Washington', 'US' );
 INSERT INTO locations ( location_id, street_address, postal_code, city, state_province,
 country_id)
VALUES ( 1800, '147 Spadina Ave', 'M5V 2L7', 'Toronto', 'Ontario', 'CA');
 INSERT INTO locations ( location_id, street_address, postal_code, city, state_province,
 country_id)
VALUES ( 2400, '8204 Arthur St', NULL, 'London', NULL, 'UK' );
 INSERT INTO locations ( location_id, street_address, postal_code, city, state_province,
 country_id)
VALUES ( 2500, 'Magdalen Centre, The Oxford Science Park', 'OX9 9ZB', 'Oxford', 'Oxford',
 'UK' );
INSERT INTO locations ( location_id, street_address, postal_code, city, state_province,
 country_id)
VALUES  ( 2700, 'Schwanthalerstr. 7031', '80925', 'Munich', 'Bavaria', 'DE' );

-- Insert departments table. 
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (1, 'Administration', 1700);
 INSERT INTO departments ( department_id, department_name, location_id)
VALUES (2, 'Marketing', 1800);
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (3, 'Purchasing', 1700);
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (4, 'Human Resources', 2400);
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (5, 'Shipping', 1500);
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (6, 'IT', 1400);
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (7, 'Public Relations', 2700);
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (8, 'Sales', 2500);
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (9, 'Executive', 1700);
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (10, 'Finance', 1700);
INSERT INTO departments ( department_id, department_name, location_id)
VALUES (11, 'Accounting', 1700); 

-- Insert data into employees table
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 100, 'Steven', 'King', 'steven.king@sqltutorial.org', '515.123.4567', '17-JUN-87', 4, '24000.00', NULL, 9 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 101, 'Neena', 'Kochhar', 'neena.kochhar@sqltutorial.org', '515.123.4568', '21-SEP-89', 5, '17000.00', 100, 9 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 102, 'Lex', 'De Haan', 'lex.de haan@sqltutorial.org', '515.123.4569', '13-JAN-93', 5, '17000.00', 100, 9 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 103, 'Alexander', 'Hunold', 'alexander.hunold@sqltutorial.org', '590.423.4567', '3-JAN-90', 9, '9000.00', 102, 6 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 104, 'Bruce', 'Ernst', 'bruce.ernst@sqltutorial.org', '590.423.4568', '21-MAY-91', 9, '6000.00', 103, 6 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 105, 'David', 'Austin', 'david.austin@sqltutorial.org', '590.423.4569', '25-JUN-97', 9, '4800.00', 103, 6 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 106, 'Valli', 'Pataballa', 'valli.pataballa@sqltutorial.org', '590.423.4560', '05-FEB-98', 9, '4800.00', 103, 6 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 107, 'Diana', 'Lorentz', 'diana.lorentz@sqltutorial.org', '590.423.5567', '7-FEB-99', 9, '4200.00', 103, 6 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 108, 'Nancy', 'Greenberg', 'nancy.greenberg@sqltutorial.org', '515.124.4569', '17-AUG-94', 7, '12000.00', 101, 10 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 109, 'Daniel', 'Faviet', 'daniel.faviet@sqltutorial.org', '515.124.4169', '16-AUG-94', 6, '9000.00', 108, 10 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 110, 'John', 'Chen', 'john.chen@sqltutorial.org', '515.124.4269', '28-SEP-97', 6, '8200.00', 108, 10 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 111, 'Ismael', 'Sciarra', 'ismael.sciarra@sqltutorial.org', '515.124.4369', '30-SEP-97', 6, '7700.00', 108, 10 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 112, 'Jose Manuel', 'Urman', 'jose manuel.urman@sqltutorial.org', '515.124.4469', '7-MAR-98', 6, '7800.00', 108, 10 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 113, 'Luis', 'Popp', 'luis.popp@sqltutorial.org', '515.124.4567', '7-DEC-99', 6, '6900.00', 108, 10 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 114, 'Den', 'Raphaely', 'den.raphaely@sqltutorial.org', '515.127.4561', '7-DEC-94', 14, '11000.00', 100, 3 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 115, 'Alexander', 'Khoo', 'alexander.khoo@sqltutorial.org', '515.127.4562', '18-MAY-95', 13, '3100.00', 114, 3 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 116, 'Shelli', 'Baida', 'shelli.baida@sqltutorial.org', '515.127.4563', '24-DEC-97', 13, '2900.00', 114, 3 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 117, 'Sigal', 'Tobias', 'sigal.tobias@sqltutorial.org', '515.127.4564', '24-JUL-97', 13, '2800.00', 114, 3 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 118, 'Guy', 'Himuro', 'guy.himuro@sqltutorial.org', '515.127.4565', '15-NOV-98', 13, '2600.00', 114, 3 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 119, 'Karen', 'Colmenares', 'karen.colmenares@sqltutorial.org', '515.127.4566', '10-AUG-99', 13, '2500.00', 114, 3 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 120, 'Matthew', 'Weiss', 'matthew.weiss@sqltutorial.org', '650.123.1234', '18-JUL-96', 19, '8000.00', 100, 5 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 121, 'Adam', 'Fripp', 'adam.fripp@sqltutorial.org', '650.123.2234', '10-APR-97', 19, '8200.00', 100, 5 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 122, 'Payam', 'Kaufling', 'payam.kaufling@sqltutorial.org', '650.123.3234', '1-MAY-95', 19, '7900.00', 100, 5 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 123, 'Shanta', 'Vollman', 'shanta.vollman@sqltutorial.org', '650.123.4234', '10-OCT-97', 19, '6500.00', 100, 5 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 126, 'Irene', 'Mikkilineni', 'irene.mikkilineni@sqltutorial.org', '650.124.1224', '28-SEP-98', 18, '2700.00', 120, 5 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 145, 'John', 'Russell', 'john.russell@sqltutorial.org', NULL, '01-OCT-96', 15, '14000.00', 100, 8 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 146, 'Karen', 'Partners', 'karen.partners@sqltutorial.org', NULL, '5-JAN-97', 15, '13500.00', 100, 8 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 176, 'Jonathon', 'Taylor', 'jonathon.taylor@sqltutorial.org', NULL, '24-MAR-98', 16, '8600.00', 100, 8 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 177, 'Jack', 'Livingston', 'jack.livingston@sqltutorial.org', NULL, '23-APR-98', 16, '8400.00', 100, 8 );
INSERT INTO employees ( employee_id, first_name, last_name, email, phone_number, hire_date, job_id, salary, manager_id, department_id)
VALUES ( 178, 'Kimberely', 'Grant', 'kimberely.grant@sqltutorial.org', NULL, '24-MAY-99', 16, '7000.00', 100, 8 );

--  Insert data into the dependents table
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 4, 'Jennifer', 'King', 'Child', 100);
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 5, 'Johnny', 'Kochhar', 'Child', 101 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 6, 'Bette', 'De Haan', 'Child', 102 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 7, 'Grace', 'Faviet', 'Child', 109 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 8, 'Matthew', 'Chen', 'Child', 110 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 9, 'Joe', 'Sciarra', 'Child', 111 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 10, 'Christian', 'Urman', 'Child', 112 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 11, 'Zero', 'Popp', 'Child', 113 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 12, 'Karl', 'Greenberg', 'Child', 108 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 14, 'Vivien', 'Hunold', 'Child', 103 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 15, 'Cuba', 'Ernst', 'Child', 104 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 16, 'Fred', 'Austin', 'Child', 105 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 17, 'Helen', 'Pataballa', 'Child', 106 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 18, 'Dan', 'Lorentz', 'Child', 107 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 22, 'Elvis', 'Khoo', 'Child', 115 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 23, 'Sandra', 'Baida', 'Child', 116 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 24, 'Cameron', 'Tobias', 'Child', 117 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 25, 'Kevin', 'Himuro', 'Child', 118 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 26, 'Rip', 'Colmenares', 'Child', 119 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 27, 'Julia', 'Raphaely', 'Child', 114 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 28, 'Woody', 'Russell', 'Child', 145 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 29, 'Alec', 'Partners', 'Child', 146 );
INSERT INTO dependents ( dependent_id, first_name, last_name, relationship, employee_id)
VALUES ( 30, 'Sandra', 'Taylor', 'Child', 176 );