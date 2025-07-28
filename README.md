# Omohide Map Backend

このプロジェクトは、Omohide Mapのバックエンドを提供するGo言語で書かれたアプリケーションです。

## 必要条件

- Go 1.24.5以上
- PostgreSQL データベース

## セットアップ

1. リポジトリをクローンします。

   ```bash
   git clone https://github.com/yourusername/omohide_map_backend.git
   cd omohide_map_backend
   ```

2. 必要なGoモジュールをインストールします。

   ```bash
   go mod tidy
   ```

3. `.env`ファイルを作成し、データベースのURLを設定します。

   ```
   DATABASE_URL=your_database_url
   ```

4. アプリケーションを起動します。

   ```bash
   go run main.go
   ```

## 使用技術

- [Echo](https://echo.labstack.com/) - 高性能なGoのWebフレームワーク
- [GORM](https://gorm.io/) - GoのためのORMライブラリ
- [PostgreSQL](https://www.postgresql.org/) - オープンソースのリレーショナルデータベース

## ディレクトリ構成

- `internal/handlers` - アプリケーションのハンドラを提供します。
- `internal/db` - データベース接続機能を提供します。

## 貢献

貢献を歓迎します。バグ報告やプルリクエストをお待ちしています。

## ライセンス

このプロジェクトはMITライセンスの下でライセンスされています。
