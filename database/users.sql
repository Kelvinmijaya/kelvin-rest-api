CREATE TABLE users (
  id SERIAL PRIMARY KEY, 
  name VARCHAR (50) UNIQUE NOT NULL, 
  password VARCHAR (255) NOT NULL, 
  email VARCHAR (255) UNIQUE NOT NULL, 
  created_at TIMESTAMP NOT NULL, 
  last_login TIMESTAMP
);

INSERT INTO users VALUES (1, 'kelvin', '$2a$14$orzPWHEbg/61ePxETQYRCO0azwy7f/aDIHPeZ1HaNTLsGRmMEZkbm', 'kelvinmijaya@gmail.com', 2024-01-25 08:14:24.99874, '');
