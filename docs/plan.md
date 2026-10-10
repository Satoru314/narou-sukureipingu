# 小説家になろうスクレイピングAPI 実装プラン

## 概要
小説家になろうから小説データを取得し、PostgreSQLに保存するHTTP APIサーバーを構築する。

## 要件

### 機能1: 目次取得API
- **エンドポイント**: `POST /api/menu`
- **リクエスト**: `{"url": "https://ncode.syosetu.com/n9619hu"}`
- **レスポンス**:
```json
{
  "title": "小説タイトル",
  "author": "著者名",
  "episode_count": 405,
  "episodes": [
    {"number": 1, "subtitle": "サブタイトル1", "url": "/n9619hu/1/"},
    {"number": 2, "subtitle": "サブタイトル2", "url": "/n9619hu/2/"}
  ]
}
```
- **処理**:
  1. 目次1ページ目を取得
  2. 最終ページ番号を確認（`.c-pager__item--last`から）
  3. 全ページを順次取得してサブタイトル一覧を構築
  4. PostgreSQLに保存

### 機能2: 話取得API
- **エンドポイント**: `POST /api/story`
- **リクエスト**: `{"url": "https://ncode.syosetu.com/n9619hu/1"}`
- **レスポンス**:
```json
{
  "episode_number": 1,
  "subtitle": "サブタイトル",
  "content": "本文..."
}
```
- **処理**:
  1. 指定URLから話を取得
  2. タイトル、本文を抽出
  3. PostgreSQLに保存

## アーキテクチャ

### ディレクトリ構成
```
narou-sukureipingu/
├── cmd/
│   └── api/
│       └── main.go          # APIサーバーのエントリーポイント
├── internal/
│   ├── handler/             # HTTPハンドラー
│   │   ├── menu.go          # 目次取得ハンドラー
│   │   └── story.go         # 話取得ハンドラー
│   ├── scraper/             # スクレイピングロジック
│   │   ├── menu.go          # 目次スクレイパー
│   │   └── story.go         # 話スクレイパー
│   ├── db/                  # データベース操作
│   │   ├── db.go            # DB接続
│   │   ├── novel.go         # 小説テーブル操作
│   │   └── episode.go       # エピソードテーブル操作
│   └── models/              # データモデル
│       └── models.go        # 構造体定義
├── docker-compose.yml       # PostgreSQL + API
├── .env                     # 環境変数
└── go.mod
```

### データベーススキーマ

**novels テーブル**
```sql
CREATE TABLE novels (
    id SERIAL PRIMARY KEY,
    ncode VARCHAR(20) UNIQUE NOT NULL,  -- n9619huなど
    title TEXT NOT NULL,
    author TEXT NOT NULL,
    episode_count INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**episodes テーブル**
```sql
CREATE TABLE episodes (
    id SERIAL PRIMARY KEY,
    novel_id INT REFERENCES novels(id) ON DELETE CASCADE,
    episode_number INT NOT NULL,
    subtitle TEXT NOT NULL,
    url TEXT NOT NULL,
    content TEXT,  -- 後から取得
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(novel_id, episode_number)
);
```

## 実装ステップ

### Phase 1: 基本セットアップ
1. Dockerでpostgresqlを起動
2. .envファイル作成（DB接続情報）
3. データベース接続とテーブル作成

### Phase 2: スクレイピング機能実装
1. 既存のcmd/scraping_menu/main.goを参考に目次スクレイパー実装
   - 複数ページ対応
   - サブタイトルを配列で取得
   - 各話のURLも取得
2. 既存のcmd/scraping_story/main.goを参考に話スクレイパー実装

### Phase 3: データベース操作実装
1. PostgreSQL接続（lib/pq使用）
2. novels/episodesテーブルへのCRUD操作

### Phase 4: HTTPサーバー実装
1. net/httpまたはginでAPIサーバー構築
2. `/api/menu` エンドポイント実装
3. `/api/story` エンドポイント実装

### Phase 5: Docker統合
1. docker-compose.ymlでpostgresql + apiを起動
2. フロントエンドから利用可能な状態に

## 技術スタック
- **言語**: Go 1.21+
- **Webフレームワーク**: net/http (シンプル) or gin (高機能)
- **DB**: PostgreSQL 15
- **スクレイピング**: goquery + bluemonday
- **コンテナ**: Docker + docker-compose

## 注意事項
- User-Agent設定必須（403エラー回避）
- サーバーへの負荷軽減のため、リクエスト間に1秒待機
- エラーハンドリングは最低限でOK（シンプル重視）

## 次のアクション
1. docker-compose.yml作成
2. データベーススキーマ作成
3. スクレイピングロジック実装
4. APIサーバー実装
