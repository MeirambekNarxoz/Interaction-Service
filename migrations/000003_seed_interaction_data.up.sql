-- 000003_seed_interaction_data.up.sql

-- Наполняем жалобы (reports) на контент для проверки панели модератора
INSERT INTO reports (id, reporter_id, target_id, target_author_id, target_type, reason, status, room_id) VALUES
(1, 4, 1001, 3, 'post', 'Токсичное обсуждение и нарушение правил сообщества', 'OPEN', 101)
ON CONFLICT (id) DO NOTHING;

-- Сбрасываем счетчики sequence
SELECT setval(pg_get_serial_sequence('reports', 'id'), coalesce(max(id), 1)) FROM reports;
SELECT setval(pg_get_serial_sequence('likes', 'id'), coalesce(max(id), 1)) FROM likes;
SELECT setval(pg_get_serial_sequence('bookmarks', 'id'), coalesce(max(id), 1)) FROM bookmarks;
