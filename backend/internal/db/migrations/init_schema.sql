CREATE TABLE whiteboards (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- Use AUTO_INCREMENT for primary key in MySQL
    name VARCHAR(255),  -- Name of the whiteboard, could be optional
    owner_id INT,  -- Reference to the user who created the whiteboard
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- When the whiteboard was created
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- Last updated time
    current_state JSON  -- Storing the latest state of the whiteboard as a JSON blob for quick recovery
    FOREIGN KEY (owner_id) REFERENCES users(id)  -- Reference to the user table
);


CREATE TABLE strokes (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- Use AUTO_INCREMENT for MySQL
    whiteboard_id INT,  -- Foreign key to the whiteboard this stroke belongs to
    user_id INT,  -- Foreign key to the user who made this action
    stroke_data JSON,  -- Store the detailed data about the stroke, stored in JSON
    -- example of stroke_data
    -- {
    --     "type": "stroke",
    --     "color": "#ff0000",
    --     "width": 5,
    --     "path": [
    --         { "x": 10, "y": 15 },
    --         { "x": 12, "y": 18 },
    --         { "x": 15, "y": 20 }
    --     ]
    -- }
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- When the action happened
);


CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,  -- Automatically increments the primary key
    name VARCHAR(255),  -- Name of the user
    email VARCHAR(255),  -- Email address of the user
    role ENUM('Admin', 'Editor', 'Viewer')  -- User role: Admin, Editor, or Viewer
);

