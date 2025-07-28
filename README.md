# Omohide Map Backend

このプロジェクトは、Omohide Map のバックエンドを提供する Go 言語で書かれたアプリケーションです。

## 必要条件

- Go 1.24.5 以上
- PostgreSQL データベース

## 使用技術

- [Echo](https://echo.labstack.com/) - 高性能な Go の Web フレームワーク
- [GORM](https://gorm.io/) - Go のための ORM ライブラリ
- [PostgreSQL](https://www.postgresql.org/) - オープンソースのリレーショナルデータベース

## ディレクトリ構成

- `internal/handlers` - アプリケーションのハンドラを提供します。
- `internal/db` - データベース接続機能を提供します。
