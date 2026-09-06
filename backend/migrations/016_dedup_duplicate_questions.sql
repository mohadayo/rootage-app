-- 起動のたびに seed_java.sql が重複挿入していた問題を統合する（一回性のデータ移行）。
-- 同一 text の問題は最古の1件（canonical）を残し、残りに紐づく回答を canonical へ寄せる。
-- 新規インストールでは questions が空のため何も起きない。
-- （このマイグレーションは FK を RESTRICT に張り替える 017 より前に実行すること。
--   ここではまだ questions -> quiz_answers が CASCADE なので、寄せきれない重複回答は
--   重複問題の削除に伴って安全に消える。）

-- 1. 重複問題に紐づく回答を canonical へ付け替える。
--    ただし同一セッションで canonical に既に回答済みの場合は
--    (session_id, question_id) 一意制約に反するため付け替えない（後段の削除で消える）。
UPDATE quiz_answers qa
SET question_id = canon.id
FROM (
    SELECT id, text FROM (
        SELECT id, text, ROW_NUMBER() OVER (PARTITION BY text ORDER BY created_at, id) AS rn
        FROM questions
    ) t WHERE rn = 1
) canon
JOIN (
    SELECT id, text FROM (
        SELECT id, text, ROW_NUMBER() OVER (PARTITION BY text ORDER BY created_at, id) AS rn
        FROM questions
    ) t WHERE rn > 1
) dup ON dup.text = canon.text
WHERE qa.question_id = dup.id
  AND NOT EXISTS (
      SELECT 1 FROM quiz_answers x
      WHERE x.session_id = qa.session_id AND x.question_id = canon.id
  );

-- 2. 重複問題の行を削除する（残った付け替え不能な回答は CASCADE で消える）。
DELETE FROM questions WHERE id IN (
    SELECT id FROM (
        SELECT id, ROW_NUMBER() OVER (PARTITION BY text ORDER BY created_at, id) AS rn
        FROM questions
    ) t WHERE rn > 1
);
