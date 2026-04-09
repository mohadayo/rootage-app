-- 問題に難易度カラムを追加
ALTER TABLE questions ADD COLUMN IF NOT EXISTS difficulty VARCHAR(20) NOT NULL DEFAULT 'beginner';

-- quiz_sessionsに難易度カラムを追加
ALTER TABLE quiz_sessions ADD COLUMN IF NOT EXISTS difficulty VARCHAR(20) NOT NULL DEFAULT 'beginner';
