-- パスワード変更・ロール変更・退職時に既存 JWT を失効させるためのトークン世代。
-- 発行時に claims へ埋め、検証時に DB の現在値と突合する。
-- DEFAULT 0 なので、この列を持たない既存トークン（claim なし=0）は引き続き有効。
ALTER TABLE users ADD COLUMN IF NOT EXISTS token_version INTEGER NOT NULL DEFAULT 0;
