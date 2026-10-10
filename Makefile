# docker環境を起動
.PHONY: run-dev
run-dev:
	docker compose up -d

# Dockerコンテナを停止して、ボリュームも削除(データベースのデータも消える)
.PHONY: destroy
destroy:
	docker compose down -v

# データベースのマイグレーションだけを実行
# 注意⚠️: 事前にPostgreSQLをDocker環境で起動しておくこと
.PHONY: db-migrate
db-migrate:
	docker compose exec web go run ./cmd/migration/main.go

# データベースにデータを投入する
# 注意⚠️: データベースに対してマイグレーションを実行していること
.PHONY: insert-data
insert-data:
	docker compose exec web go run ./cmd/scripts/main.go
	
# データベースのコンテナとボリューム削除 & コンテナ立ち上げ & マイグレーション実行
.PHONY: init-db
init-db:
	$(MAKE) destroy
	$(MAKE) run-dev
	sleep 3 # PostgreSQL起動を待機中...
	$(MAKE) db-migrate
	$(MAKE) insert-data

# テストデータを投入するコマンド
seed:
	@echo "テストデータを投入します..."
	cat database/seeds/test_data.sql | docker compose exec -T db psql -U your_username -d your_database
	@echo "完了しました"