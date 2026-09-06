-- 難易度は beginner / intermediate / advanced の3値のみ許可する。
-- CSV 一括投入で '中級' のような不正値が保存され、出題クエリの完全一致から
-- 漏れて誰にも出題されない事故を防ぐ。
-- 既存の不正値は beginner に寄せてから CHECK を張る（ALTER 失敗を避ける）。
UPDATE questions SET difficulty = 'beginner'
    WHERE difficulty IS NULL OR difficulty NOT IN ('beginner', 'intermediate', 'advanced');

ALTER TABLE questions DROP CONSTRAINT IF EXISTS questions_difficulty_check;
ALTER TABLE questions ADD CONSTRAINT questions_difficulty_check
    CHECK (difficulty IN ('beginner', 'intermediate', 'advanced'));
