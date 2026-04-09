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
| `DB_URL` | PostgreSQL接続文字列 | localhost:5432 |
| `JWT_SECRET` | JWTトークン署名キー | dev-secret-key |
| `OPENAI_API_KEY` | OpenAI APIキー | - |
| `RESEND_API_KEY` | Resend APIキー（パスワードリセット用） | - |
| `BASE_URL` | アプリのベースURL | CORS_ORIGINの値 |
| `CORS_ORIGIN` | CORS許可オリジン | http://localhost:5173 |
| `ALLOWED_EMAIL_DOMAIN` | 登録許可メールドメイン | - |
| `SERVER_PORT` | サーバーポート | 8080 |
| `VITE_ENABLE_RAG` | 社内検索AI機能の有効化（フロント） | false |

## デモアカウント

| ユーザー | メール | パスワード | 権限 |
|---|---|---|---|
| 管理者 | admin@example.com | admin1234 | admin |

一般ユーザーは新規登録画面から作成できます。`ALLOWED_EMAIL_DOMAIN` を設定するとドメイン制限が有効になります。

## テスト

```bash
cd backend
go test ./... -v
```

48テスト（ユニットテスト + 統合テスト）が実行されます。
