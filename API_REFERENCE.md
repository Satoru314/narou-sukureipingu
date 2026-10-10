# API リファレンス

## ベースURL

```
http://localhost:8080
```

---

## エンドポイント一覧

- [`GET /api/menu`](#get-apimenu) - 小説の目次取得
- [`GET /api/story`](#get-apistory) - 個別エピソード取得
- [`GET /api/novel/fetch-all`](#get-apinovelfetch-all) - 全エピソード一括取得
- [`GET /health`](#get-health) - ヘルスチェック

---

## GET /api/menu

小説の目次（エピソード一覧）を取得し、データベースに保存します。

### リクエストパラメータ

| パラメータ | 型     | 必須 | 説明                           | 例         |
|------------|--------|------|--------------------------------|------------|
| `ncode`    | string | ✓    | 小説のncode（英数字）          | `n9619hu`  |

### レスポンス

**成功時（200 OK）:**

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
    {
      "number": 2,
      "subtitle": "時々清掃人",
      "url": "/n9619hu/2/"
    }
    // ... 以下全エピソード
  ]
}
```

### レスポンスフィールド

| フィールド       | 型     | 説明                               |
|------------------|--------|------------------------------------|
| `ncode`          | string | 小説のncode                        |
| `title`          | string | 小説のタイトル                     |
| `author`         | string | 作者名                             |
| `author_url`     | string | 作者ページのURL                    |
| `description`    | string | 小説の説明文                       |
| `novel_url`      | string | 小説ページのURL                    |
| `episode_count`  | int    | エピソード総数                     |
| `episodes`       | array  | エピソード一覧                     |
| `episodes[].number` | int | エピソード番号                   |
| `episodes[].subtitle` | string | エピソードタイトル           |
| `episodes[].url` | string | エピソードの相対URL                |

### エラーレスポンス

**400 Bad Request:**
```
ncode query parameter is required
```

**500 Internal Server Error:**
```
Failed to scrape menu: <エラー詳細>
```

### 使用例

```bash
# curlで取得
curl "http://localhost:8080/api/menu?ncode=n9619hu"

# JSONを整形して表示
curl "http://localhost:8080/api/menu?ncode=n9619hu" | jq .

# エピソード数だけ取得
curl "http://localhost:8080/api/menu?ncode=n9619hu" | jq .episode_count
```

### 副作用

- 小説情報を`novels`テーブルに保存（既存の場合は更新）
- 全エピソード情報を`episodes`テーブルに保存（本文は保存されない）
- 目次が複数ページある場合、5秒間隔で全ページを自動取得

---

## GET /api/story

特定のエピソードの本文を取得し、データベースに保存します。

### リクエストパラメータ

| パラメータ        | 型     | 必須 | 説明                      | 例         |
|-------------------|--------|------|---------------------------|------------|
| `ncode`           | string | ✓    | 小説のncode               | `n9619hu`  |
| `episode_number`  | int    | ✓    | エピソード番号（1始まり） | `1`        |

### レスポンス

**成功時（200 OK）:**

```json
{
  "episode_number": 1,
  "subtitle": "中途半端な剣",
  "content": "ショートソードと比べると幾分長く、細かい取り回しに苦労する。\nロングソードと比較すればそのリーチはやや物足りず...",
  "posted_at": "2023-01-15T00:00:00Z"
}
```

### レスポンスフィールド

| フィールド        | 型       | 説明                                   |
|-------------------|----------|----------------------------------------|
| `episode_number`  | int      | エピソード番号                         |
| `subtitle`        | string   | エピソードタイトル                     |
| `content`         | string   | 本文（改行は`\n`）                     |
| `posted_at`       | string   | 投稿日時（ISO 8601形式、nullの場合あり）|

### エラーレスポンス

**400 Bad Request:**
```
ncode and episode_number query parameters are required
```
```
episode_number must be a valid integer
```

**500 Internal Server Error:**
```
Failed to scrape story: <エラー詳細>
```

### 使用例

```bash
# 第1話を取得
curl "http://localhost:8080/api/story?ncode=n9619hu&episode_number=1"

# 本文だけ抽出
curl "http://localhost:8080/api/story?ncode=n9619hu&episode_number=1" | jq -r .content

# 複数話を連続取得（5秒間隔）
for i in {1..5}; do
  curl "http://localhost:8080/api/story?ncode=n9619hu&episode_number=$i"
  sleep 5
done
```

### 副作用

- エピソードの本文と投稿日時を`episodes`テーブルに保存（更新）
- 事前に`/api/menu`でエピソード情報が登録されている必要がある

---

## GET /api/novel/fetch-all

小説の全エピソードを自動で一括取得します。

### リクエストパラメータ

| パラメータ | 型     | 必須 | 説明            | 例        |
|------------|--------|------|-----------------|-----------|
| `ncode`    | string | ✓    | 小説のncode     | `n9619hu` |

### レスポンス

**成功時（200 OK）:**

```json
{
  "ncode": "n9619hu",
  "title": "バスタード・ソードマン",
  "total_episodes": 406,
  "fetched_count": 396,
  "message": "Successfully fetched 396 new episodes, skipped 10 already fetched (total: 406)"
}
```

### レスポンスフィールド

| フィールド        | 型     | 説明                           |
|-------------------|--------|--------------------------------|
| `ncode`           | string | 小説のncode                    |
| `title`           | string | 小説のタイトル                 |
| `total_episodes`  | int    | 総エピソード数                 |
| `fetched_count`   | int    | 今回新規取得したエピソード数   |
| `message`         | string | 実行結果のメッセージ           |

### エラーレスポンス

**400 Bad Request:**
```
ncode query parameter is required
```

**500 Internal Server Error:**
```
Failed to scrape menu: <エラー詳細>
```

### 動作仕様

1. 小説の目次を取得してデータベースに保存
2. 全エピソードを順番に取得
3. **各リクエスト間に5秒の待機時間**（サーバー負荷軽減）
4. **既に本文が保存されているエピソードは自動スキップ**
5. 処理状況はサーバーログに出力

### 処理時間の目安

- 100話: 約8分（100話 × 5秒）
- 200話: 約17分（200話 × 5秒）
- 400話: 約33分（400話 × 5秒）

### 使用例

```bash
# 全エピソードを一括取得
curl "http://localhost:8080/api/novel/fetch-all?ncode=n9619hu"

# バックグラウンドで実行
curl "http://localhost:8080/api/novel/fetch-all?ncode=n9619hu" > result.json 2>&1 &

# 進捗をログで確認
docker compose logs web --tail=50 -f
```

### ログ出力例

```
2026/01/06 04:20:00 Fetching menu for n9619hu...
2026/01/06 04:20:21 Menu saved. Total episodes: 406
2026/01/06 04:20:21 Episode 1/406 already fetched, skipping: 中途半端な剣
2026/01/06 04:20:21 Episode 2/406 already fetched, skipping: 時々清掃人
2026/01/06 04:20:21 Fetching episode 3/406: 異世界の食事情
2026/01/06 04:20:21 Episode 3/406 saved successfully
2026/01/06 04:20:26 Fetching episode 4/406: 謎の発明家ケイオス卿
2026/01/06 04:20:26 Episode 4/406 saved successfully
```

### 注意事項

- リクエストがタイムアウトしても、サーバー側では処理が継続されます
- 途中で中断した場合、再実行すると未取得分のみを取得します
- 同じncodeで複数回実行しても、既取得分はスキップされるため安全です

### 副作用

- 小説情報と全エピソード情報をデータベースに保存
- 全エピソードの本文と投稿日時をデータベースに保存

---

## GET /health

サーバーのヘルスチェック用エンドポイント。

### リクエストパラメータ

なし

### レスポンス

**成功時（200 OK）:**

```
OK
```

### 使用例

```bash
curl "http://localhost:8080/health"
```

---

## エラーコード一覧

| HTTPステータス | 説明                                       |
|----------------|-------------------------------------------|
| 200            | 成功                                      |
| 400            | リクエストパラメータが不正                 |
| 405            | メソッドが許可されていない（GETを使用すること） |
| 500            | サーバー内部エラー（スクレイピング失敗など） |

---

## 共通ヘッダー

### レスポンスヘッダー

```
Content-Type: application/json
```

---

## レート制限

このAPIはスクレイピング対象サイトへの負荷を考慮して、以下の制限があります：

- `/api/story`: 推奨間隔 5秒
- `/api/novel/fetch-all`: 自動的に5秒間隔で実行
- `/api/menu`: 複数ページある場合、自動的に5秒間隔

**重要:** これらの制限を守らないと、対象サイトに過度な負荷をかける可能性があります。

---

## データベース連携

全エンドポイントはPostgreSQLデータベースと連携しています。

### データ永続化

- `/api/menu`: 小説情報とエピソード一覧を保存
- `/api/story`: エピソード本文を保存
- `/api/novel/fetch-all`: 全データを保存

### データ更新

- 同じncodeで再実行すると、小説情報は上書き更新されます
- エピソード情報は既存データとマージされます（`ON CONFLICT`使用）

---

## 使用例集

### 新しい小説を登録して全話取得

```bash
# 1. 目次を取得（オプション：確認用）
curl "http://localhost:8080/api/menu?ncode=n9619hu" | jq .episode_count

# 2. 全エピソードを一括取得
curl "http://localhost:8080/api/novel/fetch-all?ncode=n9619hu"
```

### 特定の話だけ取得

```bash
# 第1話だけ取得
curl "http://localhost:8080/api/story?ncode=n9619hu&episode_number=1" | jq .
```

### 更新チェックして差分取得

```bash
# 目次を再取得（エピソード数が増えているか確認）
curl "http://localhost:8080/api/menu?ncode=n9619hu" | jq .episode_count

# 全話一括取得（既取得分は自動スキップ）
curl "http://localhost:8080/api/novel/fetch-all?ncode=n9619hu"
```

### 複数の小説を順番に取得

```bash
#!/bin/bash
ncodes=("n9619hu" "n1234ab" "n5678cd")

for ncode in "${ncodes[@]}"; do
  echo "Fetching $ncode..."
  curl "http://localhost:8080/api/novel/fetch-all?ncode=$ncode"
  echo ""
done
```

---

## トラブルシューティング

### リクエストがタイムアウトする

`/api/novel/fetch-all`は大量のエピソードがある場合、完了まで時間がかかります。

**対策:**
- バックグラウンドで実行: `curl "..." > result.json 2>&1 &`
- ログで進捗確認: `docker compose logs web -f`

### エピソードが取得できない

1. 先に`/api/menu`で目次を取得していることを確認
2. ncodeとepisode_numberが正しいか確認
3. サーバーログでエラー詳細を確認: `docker compose logs web --tail=50`

### データベースエラー

```bash
# データベースを再初期化
make init-db
```

---

## サポート情報

- バグ報告: [GitHub Issues](https://github.com/your-repo/issues)
- ドキュメント: [README.md](./README.md)
