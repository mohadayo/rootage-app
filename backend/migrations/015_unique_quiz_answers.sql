-- 同一セッション・同一問題への重複回答を防ぐ。
-- 一意インデックスを張る前に、既存の重複行を最初の1件だけ残して削除する。
DELETE FROM quiz_answers a
USING quiz_answers b
WHERE a.session_id = b.session_id
  AND a.question_id = b.question_id
  AND a.ctid > b.ctid;

CREATE UNIQUE INDEX IF NOT EXISTS idx_quiz_answers_session_question
  ON quiz_answers (session_id, question_id);
