-- Down Migration: Remove 10 main categories
-- Migration name: 20250430_add_categories
-- Удаление 10 основных категорий
delete from categories
where
    name in (
        'Концерты',
        'Театр',
        'Кино',
        'Выставки',
        'Фестивали',
        'Мастер-классы',
        'Спорт',
        'Детям',
        'Экскурсии',
        'Вечеринки'
    )
;

