-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS upload (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    file_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "user" (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    username TEXT UNIQUE NOT NULL CHECK (char_length(username) > 2),
    email TEXT UNIQUE NOT NULL CHECK (email LIKE '%_@_%._%' and char_length(email) BETWEEN 5 AND 254),
    password BYTEA NOT NULL,
    avatar_id BIGINT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (avatar_id) REFERENCES upload(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS board (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    owner_user_id BIGINT NOT NULL,
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    image_id BIGINT,
    archived BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    visibility TEXT DEFAULT 'private' CHECK (visibility IN ('private', 'link', 'public')),
    FOREIGN KEY (owner_user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (image_id) REFERENCES upload(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS board_member (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    board_id BIGINT NOT NULL,
    member_role TEXT NOT NULL CHECK (member_role IN ('admin', 'member', 'observer')),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE,
    FOREIGN KEY (board_id) REFERENCES board(id) ON DELETE CASCADE,
    UNIQUE(user_id, board_id)
);

CREATE TABLE IF NOT EXISTS list (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    board_id BIGINT NOT NULL,
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    position INTEGER CHECK (position > 0),
    UNIQUE (board_id, position),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (board_id) REFERENCES board(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS card (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    author_board_member_id BIGINT NOT NULL,
    list_id BIGINT NOT NULL,
    content TEXT CHECK (char_length(content) < 5000),
    position INTEGER CHECK (position > 0),
    UNIQUE (list_id, position),
    completed  BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    complete_before TIMESTAMPTZ,
    FOREIGN KEY (author_board_member_id) REFERENCES board_member(id) ON DELETE CASCADE,
    FOREIGN KEY (list_id) REFERENCES list(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS card_member (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    board_member_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE,
    FOREIGN KEY (board_member_id) REFERENCES board_member(id) ON DELETE CASCADE,
    UNIQUE(card_id, board_member_id)
);

CREATE TABLE IF NOT EXISTS comment (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    board_member_owner_id BIGINT NOT NULL,
    content TEXT NOT NULL CHECK (char_length(content) < 5000),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE,
    FOREIGN KEY (board_member_owner_id) REFERENCES board_member(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS attachment (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    file_id BIGINT NOT NULL,
    position INTEGER CHECK (position > 0),
    UNIQUE (card_id, position),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES upload(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS checklist (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    card_id BIGINT NOT NULL,
    UNIQUE (card_id),
    title TEXT NOT NULL CHECK (char_length(title) BETWEEN 1 AND 100),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (card_id) REFERENCES card(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS checklist_point (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    checklist_id BIGINT NOT NULL,
    content TEXT NOT NULL CHECK (char_length(content) < 5000),
    checked BOOLEAN DEFAULT FALSE,
    position INTEGER CHECK (position > 0), 
    UNIQUE (checklist_id, position),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (checklist_id) REFERENCES checklist(id) ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION set_updated_at()  
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_upload_updated_at
BEFORE UPDATE ON upload
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_user_updated_at
BEFORE UPDATE ON "user"
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_board_updated_at
BEFORE UPDATE ON board
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_board_member_updated_at
BEFORE UPDATE ON board_member
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_list_updated_at
BEFORE UPDATE ON list
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_card_updated_at
BEFORE UPDATE ON card
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_card_member_updated_at
BEFORE UPDATE ON card_member
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_comment_updated_at
BEFORE UPDATE ON comment
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_attachment_updated_at
BEFORE UPDATE ON attachment
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_checklist_updated_at
BEFORE UPDATE ON checklist
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_checklist_point_updated_at
BEFORE UPDATE ON checklist_point
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

INSERT INTO upload (title, file_url, created_at, updated_at)
VALUES 
    ('sofia_avatar','img/sofia_avatar.jpg','2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('pak_avatar','img/pak.jpg','2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('board_image','img/board1.jpg','2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('attachment_image','attachment/picture1.png','2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO "user" (username, email, password, avatar_id, created_at, updated_at) 
VALUES 
    ('sofia', 'sofia@mail.ru', '$2a$12$5qc0NPEzjDdPrVq6jhCVfeMRPG/K6ZcWansBSFdH6Yra0Yd9.0vRe', 1, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('pak', 'pak@gmail.com', '$2a$12$EJGCrfUc0oO2lzTgSOpnA.HAyEy0DbpgDXJWo85xtZvUZPmPTujwW', 1, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('danila', 'danila@mail.ru', '$2a$11$fObdbVmIz6yatIyvZgIZ7.XNuF2yHP1Ro0aNh8TwsqhrmsrF5dsxm', 1, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO board (owner_user_id, title, image_id, visibility, created_at, updated_at)
VALUES 
    (1, 'Задачи на неделю', 3, 'private', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (1, 'Личные задачи', NULL, 'private', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (2, 'Проект на семестр', NULL, 'public', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO board_member (user_id, board_id, member_role, created_at, updated_at)
VALUES 
    (1, 1, 'admin', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (2, 1, 'member', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (3, 1, 'observer', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),

    (1, 2, 'admin', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),  

    (2, 3, 'admin', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (1, 3, 'member', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO list (board_id, title, position, created_at, updated_at)
VALUES 
    (1, 'Понедельник', 1, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (1, 'Вторник', 2, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (1, 'Готово', 3, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (2, 'TODO', 1, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (3, 'Идеи', 1, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (3, 'Наработки', 2, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO card (author_board_member_id, list_id, content, position, complete_before,completed, created_at, updated_at)
VALUES 
    (1, 1, 'Сделать дз', 1, '2025-11-19 18:00:00+03', false, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (2, 2, 'Прочитать книгу', 1, '2025-11-01 15:00:00+03', false, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (1, 2, 'Сходить в спортзал', 2, NULL, false, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    
    (4, 4, 'Записаться к врачу', 1, '2025-12-20 10:00:00+03', false,'2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),

    (5, 5, 'Новая фича', 1, NULL, false, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO card_member (card_id, board_member_id, created_at, updated_at)
VALUES 
    (5, 6, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO comment (card_id, board_member_owner_id, content, created_at, updated_at)
VALUES 
    (1, 2, 'Сделай это до среды', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (1, 1, 'Проверь информацию', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (2, 3, 'Когда будем реализовывать?', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO attachment (card_id, title, file_id, position, created_at, updated_at)
VALUES 
    (1, 'Фотка', 4, 1, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO checklist (card_id, title, created_at, updated_at)
VALUES 
    (2, 'Шаги разработки', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO checklist_point (checklist_id, content, position, checked, created_at, updated_at)
VALUES 
    (1, 'Сделать схему', 1, true, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (1, 'Разделить задачи', 2, false, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (1, 'Проверить работу', 3, false, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
