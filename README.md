# 小説家になろう スクレイピング API

「小説家になろう」から小説データを取得し、PostgreSQLデータベースに保存するWeb APIサーバーです。

## 機能

- 小説の目次（エピソード一覧）を取得
- 個別エピソードの本文を取得
- 小説の全エピソードを一括取得（2秒間隔、既取得スキップ機能付き）

## セットアップ

### 前提条件

- Docker & Docker Compose
- Make（オプション）

### 環境構築

1. リポジトリをクローン
2. `.env`ファイルを作成（既に存在する場合はそのまま）
3. データベース初期化とサーバー起動

```bash
# データベース初期化（初回のみ、または再構築時）
make init-db

# 開発サーバー起動
make run-dev
```

サーバーは `http://localhost:8080` で起動します。

## API エンドポイント

### 1. 目次取得 - `/api/menu`

小説の目次（エピソード一覧）を取得します。

**メソッド:** `GET`

**クエリパラメータ:**
- `ncode` (必須): 小説のncode（例: `n9619hu`）

**使用例:**
```bash
curl "http://localhost:8080/api/menu?ncode=n9619hu"
```

**レスポンス例:**
```json
{
  "ncode": "n9619hu",
  "title": "バスタード・ソードマン",
  "author": "ジェームズ・リッチマン",
  "author_url": "https://mypage.syosetu.com/375252/",
  "description": "それなりに強力なギフトを持って異世界に転生したものの...",
  "novel_url": "https://ncode.syosetu.com/n9619hu",
  "episode_count": 406,
  "episodes": [
    {
      "number": 1,
      "subtitle": "中途半端な剣",
      "url": "/n9619hu/1/"
    },
    ...
  ]
}
```

**効果:**
- 小説情報をデータベースに保存
- 全エピソードの情報（タイトル、URL）をデータベースに保存
- 複数ページの目次がある場合、自動で全ページ取得（5秒間隔）

---

### 2. エピソード本文取得 - `/api/story`

特定のエピソードの本文を取得します。

**メソッド:** `GET`

**クエリパラメータ:**
- `ncode` (必須): 小説のncode（例: `n9619hu`）
- `episode_number` (必須): エピソード番号（例: `1`）

**使用例:**
```bash
curl "http://localhost:8080/api/story?ncode=n9619hu&episode_number=1"
```

**レスポンス例:**
```json
{
  "episode_number": 1,
  "subtitle": "中途半端な剣",
  "content": "ショートソードと比べると幾分長く、細かい取り回しに苦労する...",
  "posted_at": "2023-01-15T00:00:00Z"
}
```

**効果:**
- エピソードの本文と投稿日時をデータベースに保存
- 既に目次取得済みのエピソードに本文を追加

---

### 3. 全エピソード一括取得 - `/api/novel/fetch-all`

小説の全エピソードを自動で取得します。

**メソッド:** `GET`

**クエリパラメータ:**
- `ncode` (必須): 小説のncode（例: `n9619hu`）

**使用例:**
```bash
curl "http://localhost:8080/api/novel/fetch-all?ncode=n9619hu"
```

**レスポンス例:**
```json
{
  "ncode": "n9619hu",
  "title": "バスタード・ソードマン",
  "total_episodes": 406,
  "fetched_count": 396,
  "message": "Successfully fetched 396 new episodes, skipped 10 already fetched (total: 406)"
}
```

**効果:**
1. 小説の目次を取得してデータベースに保存
2. 全エピソードを順番に取得
3. **各リクエスト間に5秒の待機時間**（サーバー負荷軽減）
4. **既に取得済みのエピソードは自動スキップ**（効率的な再取得）
5. 進捗状況がサーバーログに出力

**注意事項:**
- 大量のエピソードがある場合、完了まで時間がかかります
  - 例: 400話の小説 → 約33分（400話 × 5秒）
- 処理中にリクエストがタイムアウトしても、サーバー側では処理が継続されます
- 途中で中断した場合、再度実行すると未取得のエピソードのみを取得します

---

## Makeコマンド

```bash
# データベース初期化（DBを削除して再作成）
make init-db

# 開発サーバー起動
make run-dev

# サーバー停止
make stop

# ログ表示
make logs

# 全コンテナ削除（データも削除）
make destroy
```

## データベーススキーマ

### novels テーブル
- `ncode` (PRIMARY KEY): 小説のncode
- `title`: 小説タイトル
- `author`: 作者名
- `author_url`: 作者ページURL
- `description`: 小説の説明
- `novel_url`: 小説のURL
- `episode_count`: エピソード数
- `created_at`, `updated_at`: 作成/更新日時

### episodes テーブル
- `ncode` (PRIMARY KEY): 小説のncode
- `episode_number` (PRIMARY KEY): エピソード番号
- `subtitle`: エピソードタイトル
- `url`: エピソードURL
- `content`: 本文（NULL = 未取得）
- `posted_at`: 投稿日時
- `created_at`: 作成日時

## 技術スタック

- **言語**: Go 1.25+
- **データベース**: PostgreSQL 18
- **スクレイピング**: goquery
- **コンテナ**: Docker & Docker Compose
- **開発ツール**: Air (ホットリロード)

## ライセンス

このプロジェクトは教育目的で作成されています。スクレイピング対象サイトの利用規約を遵守してください。
