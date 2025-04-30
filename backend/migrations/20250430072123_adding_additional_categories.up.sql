-- Up Migration: Add 10 main categories
-- Migration name: 20250430_add_categories
INSERT INTO categories (name) VALUES
    ('Концерты'),
    ('Театр'),
    ('Кино'),
    ('Выставки'),
    ('Фестивали'),
    ('Мастер-классы'),
    ('Спорт'),
    ('Детям'),
    ('Экскурсии'),
    ('Вечеринки');

