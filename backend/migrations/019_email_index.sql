-- メールアドレスの大文字小文字を区別しない照合を高速化する。
-- 既存データの lower 化と UNIQUE 化は、既存の大小違い重複が存在すると失敗しうるため
-- ここでは行わない（PR の「マージ前チェックリスト」で手当てののち UNIQUE INDEX を張る）。
CREATE INDEX IF NOT EXISTS idx_users_email_lower ON users (lower(email));
