CREATE TABLE IF NOT EXISTS users(
  id INTEGER PRIMARY KEY AUTOINCREMENT, 
  email TEXT,
  username TEXT,
  pass TEXT
);

INSERT INTO users(email,username,pass) VALUES("sam@ya.ru","sam","password");
INSERT INTO users(email,username,pass) VALUES("Dan@ya.ru","Danya","password");
INSERT INTO users(email,username,pass) VALUES("Sofia@ya.ru","Sofia","password");