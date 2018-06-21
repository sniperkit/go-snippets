# login as root;

mysql -u root -p;

# use database mysql;

use mysql;

# create use scott with password tiger;

create user 'scott'@'localhost' identified by 'tiger';

grant select,insert,update,delete,create,drop,execute,references on *.* to 'scott'@'localhost';

# create database 

create database databasename;

# remove database

drop database databasename;

# create table

create table Course (
    courseId char(5), 
    subjectId char(4) not null,
    courseNumber integer,
    title varchar(50) not null,
    numOfCredits integer,
    primary key (courseId)
);

create table Student (
    ssn char(9),
    firstName varchar(25),
    mi char(1),
    lastName varchar(25),
    birthDate date,
    street varchar(25),
    phone char(11),
    zipCode char(5),
    deptId char(4),
    primary key (ssn)
);

create table Enrollment (
    ssn char(9),
    courseId char(5),
    dateRegistered date,
    grade char(1),
    primary key (ssn, courseId),
    foreign key (ssn) references Student (ssn),
    foreigh key (courseId) references Course (courseId)
);

# remove table

drop table tablename

# insert into table

insert into tableName [(column1, column2, ..., column)] values (value1, value2, ..., valuen);
#insert into Course (courseId,subjectId,courseNumber,title,numOfCredits) values ('11113', 'CSCI', '3720', 'Database System', 3);

# update table

update tableName set column1 = newValue1[, columu2 = newValue2,...] [where condition];
# update Course set numOfCredits = 4 where title = 'Database System';

# delete from table

delete [from] tableName [where condition];
# delete Course where title = 'Database System';
# delete Course;

#Simple Query

select column-list from table-list [where condition];
# select firstName, mi, lastName from Student where deptId = 'CS';
# select firstName, mi, lastName from Student where deptId = 'CS' and zipCode = '31411';
# select * from Student where deptId = 'CS' and zipCode = '31411';

# The like, between-and, and is null Operators
# SQL has a like operator that can be used for pattern matching. The syntax to check wheter a string s has a pattern p is : 
    's like p' or 's not like p'
    # You can use the wild-card characters %(percent symbol) and _ (underline symbol) in the pattern p. % matches zero or more charactors, and _ matches any single characer in s. For example, lastName like '_mi%' matches any string whose second and third letters are m and i. lastName not like '_mi%' excludes any string whose second and third letters are m and i.

# The between-and operator checks wheter a value v is between two other values, v1 and v2, using the following syntax:
    'v between v1 and v2' or 'v not between v1 and v2'

# The is null operator checks whether a value v is null using the folloing syntax:
    'v is null' or 'v is not null'

# select ssn from Enrollment where grade between 'C' and 'A'

# select lastName as "Last Name", zipCode as "Zip Code" from Student where deptId='CS';

# select title, 50 * numOfCredits as "Lecture Minutes Per Week" from Course where subjectId='CSCI';

# select subjectId as "Subject ID" from Course;

# select distinct subjectId as "Subject ID" from Course;

# Display by order

select column-list from table-list [where condition] [order by columns-to-be-sorted];

# select lastName, firstName, depId from Student where deptId='CS' order by lastName desc, firstName asc;

# select distinct lastName, firstName, courseId from Student, Enrollment where Student.ssn = Enrollment.ssn and lastName = 'Smith' and firstName = 'Jacob'
