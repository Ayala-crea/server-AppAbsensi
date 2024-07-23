CREATE TABLE users (
    user_id INT PRIMARY KEY AUTO_INCREMENT,
    full_name VARCHAR(100) NOT NULL,
    nickname VARCHAR(50),
    email VARCHAR(100) NOT NULL UNIQUE,
    phone_number VARCHAR(20),
    password VARCHAR(255) NOT NULL, -- hashed password
    role_id INT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

CREATE TABLE roles (
    role_id INT PRIMARY KEY AUTO_INCREMENT,
    role_name VARCHAR(50) NOT NULL,
    description TEXT
);

CREATE TABLE profiles (
    profile_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    image VARCHAR(255),
    address TEXT,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);


CREATE TABLE students_employees (
    id INT PRIMARY KEY AUTO_INCREMENT,
    admin_id INT NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    status ENUM('student', 'employee') NOT NULL,
    class VARCHAR(50),
    npk_or_npm VARCHAR(50),
    phone_number VARCHAR(20),
    FOREIGN KEY (admin_id) REFERENCES users(user_id)
);

CREATE TABLE locations (
    location_id INT PRIMARY KEY AUTO_INCREMENT,
    admin_id INT NOT NULL,
    location_name VARCHAR(100) NOT NULL,
    latitude DECIMAL(10, 7) NOT NULL,
    longitude DECIMAL(10, 7) NOT NULL,
    radius INT NOT NULL, -- radius dalam meter
    FOREIGN KEY (admin_id) REFERENCES users(user_id)
);

CREATE TABLE wifi_fingerprints (
    fingerprint_id INT PRIMARY KEY AUTO_INCREMENT,
    location_id INT NOT NULL,
    ap_name VARCHAR(100) NOT NULL,
    rssi INT NOT NULL,
    FOREIGN KEY (location_id) REFERENCES locations(location_id)
);


CREATE TABLE attendance_records (
    record_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    location_id INT NOT NULL,
    check_in_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('present', 'absent') NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (location_id) REFERENCES locations(location_id)
);

CREATE TABLE admin_actions (
    action_id INT PRIMARY KEY AUTO_INCREMENT,
    admin_id INT NOT NULL,
    action_type ENUM('start_attendance', 'end_attendance') NOT NULL,
    action_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id) REFERENCES users(user_id)
);

