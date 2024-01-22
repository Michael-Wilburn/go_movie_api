-- primero creamos la extensión uuid-ossp
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- creamos la tabla de usuarios
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username VARCHAR(255) UNIQUE,
    email VARCHAR(255) UNIQUE,
    password_hash VARCHAR(255),
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
    updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE FUNCTION update_updated_date()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_date = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_updated_date
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_date();

-- creamos la tabla de visualizations
CREATE TABLE visualizations (
    movie_id INT PRIMARY KEY,
    views INT DEFAULT 0
);

-- creamos la tbala de comentarios
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    user_id UUID,
    movie_id INT,
    comment_text VARCHAR(1000),
    created_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (movie_id) REFERENCES visualizations(movie_id)
);

-- Crear un trigger que actualice updated_date antes de una actualización en la fila
CREATE OR REPLACE FUNCTION update_comments_updated_date()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_date = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Vincular el trigger a la tabla de comentarios
CREATE TRIGGER update_comments_trigger
BEFORE UPDATE ON comments
FOR EACH ROW EXECUTE FUNCTION update_comments_updated_date();

