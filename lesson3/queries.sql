SELECT students.first_name, students.last_name, groups.group_name
FROM students
INNER JOIN groups
ON students.group_id = groups.id AND groups.group_name = '122-2';

SELECT students.first_name, students.last_name, groups.group_name, grades.grade
FROM students
INNER JOIN groups
ON students.group_id = groups.id AND groups.group_name = '122-2'
INNER JOIN grades
ON students.id = grades.student_id
WHERE grades.grade > 4
ORDER BY students.first_name;

SELECT students.first_name, students.last_name, sum(courses.number_of_hours)
FROM students
INNER JOIN grades
ON students.id = grades.student_id
INNER JOIN courses
ON grades.course_id = courses.id
GROUP BY students.id;

SELECT students.first_name, students.last_name, courses.title
FROM students
INNER JOIN grades
ON students.id = grades.student_id
INNER JOIN courses
ON grades.course_id = courses.id
WHERE grades.grade = 5;

SELECT groups.group_name, count(students.group_id) number_of_students
FROM groups
INNER JOIN students
ON groups.id = students.group_id
GROUP BY groups.group_name
ORDER BY number_of_students DESC;