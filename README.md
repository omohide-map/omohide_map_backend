# Omohide Map Backend

このプロジェクトは、Omohide Map のバックエンドを提供する Go 言語で書かれたアプリケーションです。

## 概要

- フレームワーク: Echo
- 認証: Firebase Authentication（ID トークン検証）
- データストア: Firestore
- ストレージ: Amazon S3（画像アップロード）
- バリデーション: go-playground/validator
- エラーハンドリング: 共通 JSON 形式

## 必要条件

- Go 1.24.5 以上
- [Firestore](https://cloud.google.com/firestore)
- [Amazon S3](https://aws.amazon.com/s3/)

## セットアップ

1. 依存関係の取得

   ```bash
   go mod download
   ```

2. 環境変数（`.env` 推奨）

   ```env
   # Firebase サービスアカウント JSON への絶対パス（必須）
   GOOGLE_APPLICATION_CREDENTIALS=/abs/path/to/service-account.json

   # Web サーバ設定（任意）
   PORT=8080

   # S3 設定（必須）
   AWS_S3_BUCKET=your-bucket-name
   AWS_REGION=ap-northeast-1
   # 認証は AWS のクレデンシャルプロバイダチェーンを使用
   # 必要に応じて以下を利用
   # AWS_ACCESS_KEY_ID=...
   # AWS_SECRET_ACCESS_KEY=...
   ```

3. ローカル起動

   ```bash
   go run main.go
   ```

   - ルート: `GET /` → 稼働確認（トークン不要）
   - API グループ: `/api/*` は Firebase ID トークン（Bearer）必須

## 使用技術

- [Echo](https://echo.labstack.com/)
- [GORM](https://gorm.io/)

## ディレクトリ構成

```text
omohide_map_backend/
├─ internal/
│  ├─ di/            # 依存注入
│  ├─ handler/       # ハンドラー
│  ├─ middleware/    # ミドルウェア
│  ├─ models/        # モデル
│  ├─ repository/    # リポジトリ
│  ├─ service/       # サービス
│  └─ storage/       # ストレージ
├─ pkg/              # 共通パッケージ
├─ routes/
│  └─ main_router.go
├─ .env
├─ main.go
├─ go.mod
├─ go.sum
└─ README.md
```
