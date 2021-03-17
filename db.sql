create table post(
  id serial not null unique, 
  authorname int not null,
  title varchar(100) not null, 
  article text not null, 
  published timestamp not null
)

CREATE TABLE users (
  id SERIAL NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  username VARCHAR(255) NOT NULL,
  password CHAR(60) NOT NULL
);