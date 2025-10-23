
INSERT INTO upload (title, file_url, created_at, updated_at)
VALUES 
    ('sofia_avatar','img/sofia_avatar.jpg','2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('pak_avatar','img/pak.jpg','2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('board_image','img/board1.jpg','2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('attachment_image','attachment/picture1.png','2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

INSERT INTO "user" (username, email, password, avatar_id, created_at, updated_at) 
VALUES 
    ('sofia', 'sofia@mail.ru', '$2a$12$5qc0NPEzjDdPrVq6jhCVfeMRPG/K6ZcWansBSFdH6Yra0Yd9.0vRe', 1, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('pak', 'pak@gmail.com', '$2a$12$EJGCrfUc0oO2lzTgSOpnA.HAyEy0DbpgDXJWo85xtZvUZPmPTujwW', 2, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    ('danila', 'danila@mail.ru', '$2a$11$fObdbVmIz6yatIyvZgIZ7.XNuF2yHP1Ro0aNh8TwsqhrmsrF5dsxm', NULL, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

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

INSERT INTO card (author_board_member_id, list_id, content, position, complete_before, created_at, updated_at)
VALUES 
    (1, 1, 'Сделать дз', 1, '2025-11-19 18:00:00+03', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (2, 2, 'Прочитать книгу', 1, '2025-11-01 15:00:00+03', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    (1, 2, 'Сходить в спортзал', 2, NULL, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),
    
    (4, 4, 'Записаться к врачу', 1, '2025-12-20 10:00:00+03', '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03'),

    (5, 5, 'Новая фича', 1, NULL, '2025-09-01 08:20:00+03', '2025-09-01 19:10:00+03');

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