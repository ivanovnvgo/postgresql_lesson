-- Группы
DROP TABLE IF EXISTS groups CASCADE;
CREATE TABLE groups (
id INT GENERATED ALWAYS AS IDENTITY,
group_name VARCHAR(200) NOT NULL,
CONSTRAINT groups_id_pkey PRIMARY KEY (id)
);

-- Студенты
DROP TABLE IF EXISTS students CASCADE;
CREATE TABLE students (
id INT GENERATED ALWAYS AS IDENTITY,
first_name VARCHAR(200) NOT NULL,
last_name VARCHAR(200) NOT NULL,
condition BOOLEAN NOT NULL,
average_mark REAL default NULL,
group_id INT NOT NULL,
CONSTRAINT student_id_pkey PRIMARY KEY (id),
CONSTRAINT students_fk_group_id FOREIGN KEY (group_id) REFERENCES groups (id)
);

-- Предметы
DROP TABLE IF EXISTS courses CASCADE;
CREATE TABLE courses (
id INT GENERATED ALWAYS AS IDENTITY,
title VARCHAR(200) NOT NULL,
number_of_hours INT NOT NULL,
CONSTRAINT courses_id_pkey PRIMARY KEY (id)
);

-- Оценки
DROP TABLE IF EXISTS grades CASCADE;
CREATE TABLE grades (
grade REAL default NULL,
student_id INT NOT NULL,
course_id INT NOT NULL,
CONSTRAINT grades_fk_student_id FOREIGN KEY (student_id) REFERENCES students (id),
CONSTRAINT grades_fk_course_id FOREIGN KEY (course_id) REFERENCES courses (id),
CONSTRAINT grade_check CHECK (grade >= 0 AND grade <= 5)
);

create index full_name_student on students (first_name, last_name);