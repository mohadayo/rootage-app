# rootage

SES企業の新人エンジニア向け育成アプリ。クイズ学習、実践ガイド、社内検索AIを統合したWebアプリケーションです。

## スクリーンショット

### ホーム画面
![ホーム画面](docs/home.png)

### クイズ学習
![クイズ学習](docs/quiz.png)

### 実践ガイド
![実践ガイド](docs/guides.png)

### 社内検索AI（RAG）
![社内検索AI](docs/rag.png)

## 機能

### クイズ学習
- 4カテゴリ（SES業界、SES面談、エンジニア基礎、現場マナー）x 3難易度（初級・中級・上級）
- カテゴリ別・難易度別のアチーブメント管理（80%以上でクリア）
- 間違えた問題の復習機能
- クイズ結果から次のアクションへの導線

### 実践ガイド
- カテゴリ別のナレッジ記事（キャリア構築、現場サバイバル、テンプレ・定型文、面談準備）
- Markdown形式の記事表示
- キーワード検索

### 社内検索AI（RAG）
- 社内文書をベースにしたAI検索（OpenAI API + pgvector）
- サンプル質問からワンタップで検索開始
- 参照ソースの表示
- 環境変数で機能の有効/無効を切替可能

### 管理機能
- 新人の学習状況ダッシュボード
- クイズ問題・カテゴリのCRUD（CSV一括インポート対応）
- 実践ガイド記事の作成・編集・公開
- 社内ナレッジ文書の登録・インデックス管理
- ユーザーのパスワードリセット

### 認証
- メールドメイン制限による登録制御
- JWT認証（24時間有効）
- メールによるパスワードリセット（Resend）
- 管理者/一般ユーザーのロール管理

## 技術スタック

| レイヤー | 技術 |
|---|---|
| バックエンド | Go (net/http) |
| フロントエンド | Vue.js 3 + TypeScript + Vite |
| スタイリング | Tailwind CSS 4 |
| 状態管理 | Pinia |
| データベース | PostgreSQL + pgvector |
| LLM/Embedding | OpenAI API |
| メール送信 | Resend |
| デプロイ | Render (Docker) |

## アーキテクチャ

```
frontend/          Vue 3 SPA
  src/
    api/           API クライアント
    views/         ページコンポーネント
    components/    共通コンポーネント
    stores/        Pinia ストア
    composables/   共通ロジック

backend/           Go API サーバー
  cmd/server/      エントリーポイント
  internal/
    handler/       HTTPハンドラー
    service/       ビジネスロジック
    repository/    データアクセス
    middleware/    認証・CORS・ロギング
    model/         データモデル
    dto/           リクエスト/レスポンス型
    pkg/           ユーティリティ（JWT, bcrypt, JSON）
    config/        環境変数管理
  migrations/      SQLマイグレーション
  seed/            シードデータ
```

## セットアップ

### 前提条件
- Go 1.22+
- Node.js 18+
- Docker / Docker Compose

### 1. 環境変数の設定

```bash
cp .env.example .env
# .env を編集
```

### 2. DB起動

```bash
docker compose up -d
```

### 3. バックエンド起動

```bash
cd backend
source ../.env
go run ./cmd/server/main.go
```

マイグレーションとシードデータは起動時に自動実行されます。

### 4. フロントエンド起動

```bash
cd frontend
npm install
npm run dev
```

`http://localhost:5173` でアクセスできます。

## 環境変数

| 変数 | 説明 | デフォルト |
|---|---|---|
| `APP_ENV` | 実行環境。`production` では `JWT_SECRET` 未設定だと起動を停止する | production |
| `DB_URL` | PostgreSQL接続文字列 | localhost:5432 |
| `JWT_SECRET` | JWTトークン署名キー。**本番では必須**（未設定だと `production` では起動しない。`development` のみ開発用の既定値でフォールバック） | -（開発時のみ dev-secret-key） |
| `OPENAI_API_KEY` | OpenAI APIキー | - |
| `RESEND_API_KEY` | Resend APIキー（パスワードリセット用） | - |
| `MAIL_FROM` | リセットメールの送信元。`onboarding@resend.dev` は Resend のサンドボックス用で本人以外に配送されないため、本番は検証済みドメインのアドレスを設定する | rootage \<onboarding@resend.dev\> |
| `BASE_URL` | アプリのベースURL | CORS_ORIGINの値 |
| `CORS_ORIGIN` | CORS許可オリジン | http://localhost:5173 |
| `ALLOWED_EMAIL_DOMAIN` | 登録許可メールドメイン（空なら制限なし） | - |
| `ADMIN_EMAIL` | 初回起動時に作成する管理者のメール（未設定なら作成しない） | - |
| `ADMIN_PASSWORD` | 同パスワード（8文字以上） | - |
| `SERVER_PORT` | サーバーポート | 8080 |
| `VITE_ENABLE_RAG` | 社内検索AI機能の有効化（フロント） | false |

> マイグレーションとシードは `schema_migrations` / `seed_applied` テーブルで適用状況を管理し、各ファイルを一度だけ実行します。マイグレーションが失敗した場合は握り潰さずに起動を停止します（壊れた状態のまま配信しないため）。マイグレーション `001` は pgvector 拡張に依存するため、DB に pgvector が入っている必要があります（`docker compose` の `pgvector/pgvector` イメージには同梱）。

## 管理者アカウント

管理者はシードデータには含まれません。`ADMIN_EMAIL` と `ADMIN_PASSWORD` を設定して起動すると、そのアカウントが存在しない場合のみ作成されます。

```bash
ADMIN_EMAIL=admin@your-domain.example ADMIN_PASSWORD='<8文字以上のパスワード>' go run ./cmd/server/main.go
```

同じメールアドレスのユーザーが既にいる場合は何もしません（運用中に変更したパスワードを上書きしないため）。作成後は環境変数を外して問題ありません。

一般ユーザーは新規登録画面から作成できます。`ALLOWED_EMAIL_DOMAIN` を設定するとドメイン制限が有効になります。

## テスト

```bash
docker compose up -d   # 統合テストは実際のDBに接続します
cd backend
go test ./... -v
```

ユニットテストと統合テストが実行されます。`internal/handler` のテストは DB に接続するため、`docker compose up -d` で PostgreSQL を起動しておく必要があります。別の接続先を使う場合は `TEST_DB_URL` で指定できます。

```bash
TEST_DB_URL="postgres://user:pass@localhost:5432/ses_quiz?sslmode=disable" go test ./...
```
