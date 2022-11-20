--Группы
INSERT INTO groups (group_name)
VALUES
('122-1'),
('122-2'),
('122-2'),
('122-3');

--Студенты
INSERT INTO students (first_name, last_name, condition, average_mark, group_id)
VALUES
('Doroty', 'Smith', true, 98.9, 1),
('Patricia', 'Johnson', true, 50.3, 2),
('Austin', 'Cintron', false, 20.3, 2),
('Eduardo', 'Hiatt', true, 70.1, 2),
('Linda', 'Williams', true, 70.8, 1);

--Предметы
INSERT INTO courses (title, number_of_hours)
VALUES
('PostgreSQL', 100),
('Golang', 120),
('HTML/CSS', 80),
('Algorithms', 200);

--Оценки
INSERT INTO grades (grade, student_id, course_id)
VALUES
(4.5, 1, 1),
(4.2, 2, 1),
(3.5, 3, 1),
(5.0, 4, 1),
(3.1, 5, 1),
(4.5, 1, 2),
(2.2, 2, 2),
(4.5, 3, 2),
(3.4, 4, 2),
(3.5, 5, 2),
(3.5, 1, 3),
(4.9, 2, 3),
(4.0, 3, 3),
(3.1, 4, 3),
(3.7, 5, 3),
(3.5, 1, 4),
(4.1, 2, 4),
(4.2, 3, 4),
(4.1, 4, 4),
(4.7, 5, 4);
