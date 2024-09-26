CREATE TABLE whiteboards (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255),
    owner_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE

);

CREATE TABLE strokes (
    id INT PRIMARY KEY AUTO_INCREMENT,
    whiteboard_id INT,                           -- Foreign key linking to the whiteboard
    owner_id INT,                                -- Who created this stroke
    x_start FLOAT,                               -- Starting x coordinate of the stroke
    y_start FLOAT,                               -- Starting y coordinate of the stroke
    x_end FLOAT,                                 -- Ending x coordinate of the stroke
    y_end FLOAT,                                 -- Ending y coordinate of the stroke
    color VARCHAR(7),                            -- Stroke color (e.g., #000000 for black)
    width INT,                                   -- Stroke width
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (whiteboard_id) REFERENCES whiteboards(id) ON DELETE CASCADE,  -- Cascade delete if whiteboard is removed
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL  -- If a user is deleted, the stroke remains but owner_id is null);
);

CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,  -- User name
    email VARCHAR(255) NOT NULL UNIQUE,  -- Email must be unique
    role ENUM('Admin', 'Editor', 'Viewer') DEFAULT 'Viewer',  -- User roles
    -- created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- When the user was created
    -- updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP  -- When the user was last updated
);