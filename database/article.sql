CREATE TABLE article (
  id SERIAL PRIMARY KEY, 
  title VARCHAR (50) UNIQUE NOT NULL, 
  url VARCHAR (50) NOT NULL, 
  content TEXT NOT NULL, 
  type VARCHAR (50) NOT NULL, 
  updated_at TIMESTAMP , 
  created_at TIMESTAMP NOT NULL 
);

