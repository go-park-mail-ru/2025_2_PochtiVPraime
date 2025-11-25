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
    member_role TEXT NOT NULL CHECK (member_role IN ('owner','admin', 'member', 'observer')),
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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
