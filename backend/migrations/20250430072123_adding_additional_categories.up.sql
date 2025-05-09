-- Up Migration: Add 10 main categories
-- Migration name: 20250430_add_categories
INSERT INTO categories (name) VALUES
    ('Концерты'),
    ('Театр'),
    ('Кино'),
    ('Выставки'),
    ('Фестивали'),
    ('Спорт'),
    ('Детям'),
    ('Экскурсии'),
    ('Вечеринки'),
    ('Клубы');

