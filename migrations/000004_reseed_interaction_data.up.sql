-- 000004_reseed_interaction_data.up.sql
-- Полная очистка и повторный засев тестовых взаимодействий

-- Удаляем тестовые seed-данные
DELETE FROM reports   WHERE id IN (1);
DELETE FROM likes     WHERE id <= 100;
DELETE FROM bookmarks WHERE id <= 100;

-- Заново добавляем жалобу для проверки панели AI-аудита
INSERT INTO reports (id, reporter_id, target_id, target_author_id, target_type, reason, status, room_id) VALUES
(1, 4, 1001, 3, 'post', 'Токсичное обсуждение и нарушение правил сообщества', 'OPEN', 101)
ON CONFLICT (id) DO UPDATE SET
    status = 'OPEN',
    reason = EXCLUDED.reason;

-- Сбрасываем счетчики
SELECT setval(pg_get_serial_sequence('reports',   'id'), GREATEST(coalesce(max(id), 1), 100)) FROM reports;
SELECT setval(pg_get_serial_sequence('likes',     'id'), GREATEST(coalesce(max(id), 1), 100)) FROM likes;
SELECT setval(pg_get_serial_sequence('bookmarks', 'id'), GREATEST(coalesce(max(id), 1), 100)) FROM bookmarks;
