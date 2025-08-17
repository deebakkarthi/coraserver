DROP DATABASE cora_db;
CREATE DATABASE cora_db;
USE cora_db;
CREATE TABLE IF NOT EXISTS slot (
    id INT,
    stime TIME NOT NULL, 
    etime TIME NOT NULL,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS subject (
    id CHAR(8),
    name VARCHAR(64) NOT NULL, 
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS faculty (
    id CHAR(254),
    name VARCHAR(64) NOT NULL,
    PRIMARY KEY (id)
);
CREATE TABLE IF NOT EXISTS static (
    class_id CHAR(4),
    day ENUM ("MON", "TUE", "WED", "THU", "FRI"), 
    slot_id INT, 
    faculty_id CHAR(254),
    subject_id CHAR(8),
    FOREIGN KEY (slot_id) REFERENCES slot (id), 
    FOREIGN KEY (faculty_id) REFERENCES faculty (id), 
    FOREIGN KEY (subject_id) REFERENCES subject (id), 
    PRIMARY KEY (class_id, day, slot_id)
);
CREATE TABLE IF NOT EXISTS dynamic (
    class_id CHAR(4),
    date DATE, 
    slot_id INT, 
    faculty_id CHAR(254) NOT NULL, 
    subject_id CHAR(8) NOT NULL, 
    FOREIGN KEY (faculty_id) REFERENCES faculty (id), 
    FOREIGN KEY (slot_id) REFERENCES slot (id), 
    FOREIGN KEY (subject_id) REFERENCES subject (id), 
    PRIMARY KEY (class_id, date, slot_id)
);
