-- セッション開始時に出題した問題を記録する（出題外の問題への回答を弾くため）
ALTER TABLE quiz_sessions ADD COLUMN IF NOT EXISTS question_ids JSONB;
