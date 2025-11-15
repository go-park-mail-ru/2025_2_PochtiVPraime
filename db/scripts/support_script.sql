CREATE TABLE IF NOT EXISTS support_form (
    "id" BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    helper_id BIGINT,
    form_type TEXT NOT NULL CHECK (form_type IN ('Баг', 'Предложение', 'Продуктовая жалоба')),
    form_status TEXT NOT NULL CHECK (form_status IN ('Открыто', 'В работе', 'Закрыто')),
    text TEXT NOT NULL,
    contact_email TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES "user"(id) ON DELETE CASCADE
);