-- カテゴリ・問題の削除で学習履歴（quiz_sessions / quiz_answers）が
-- CASCADE で巻き込まれて消える事故を防ぐ。
-- 削除は論理削除（is_active=false）に切り替え、履歴側の FK は RESTRICT にして
-- 直接 DELETE が走っても履歴を壊せないようにする。

ALTER TABLE categories ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;
ALTER TABLE questions  ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE;

-- category / question を参照する履歴の FK を CASCADE -> RESTRICT に張り替える。
-- （user_id / session_id 側は CASCADE のまま：ユーザーやセッションの削除で
--   その配下の履歴が消えるのは正しい挙動なので変更しない。）
ALTER TABLE quiz_sessions DROP CONSTRAINT IF EXISTS quiz_sessions_category_id_fkey;
ALTER TABLE quiz_sessions ADD CONSTRAINT quiz_sessions_category_id_fkey
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT;

ALTER TABLE quiz_answers DROP CONSTRAINT IF EXISTS quiz_answers_question_id_fkey;
ALTER TABLE quiz_answers ADD CONSTRAINT quiz_answers_question_id_fkey
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE RESTRICT;

CREATE INDEX IF NOT EXISTS idx_categories_active ON categories(is_active);
CREATE INDEX IF NOT EXISTS idx_questions_active  ON questions(is_active);
