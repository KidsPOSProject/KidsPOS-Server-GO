# KidsPos Go版 プロジェクトガイドライン

## プロジェクト概要

KidsPos Go版は、Spring Boot版を軽量化したPOSシステムで、Gin + Go + SQLiteで構築されています。

## 重要な実行指針

### 必須ルール

- **rules ディレクトリ内の全ルールを実行前に必ず読み込み、絶対に守ること**
- 各ルールファイルは汎用的で恒久的に利用可能な形式で記述されている
- **メモリ効率最優先**: 低メモリ環境（Raspberry Pi）での動作を常に考慮

## アーキテクチャ

### 技術スタック

- **Backend**: Go 1.21 + Gin Framework
- **Database**: SQLite (純Go実装)
- **Frontend**: HTML Templates + Bootstrap + jQuery
- **Build**: Make

### ディレクトリ構造

```
KidsPOS-Server-GO/
├── cmd/
│   └── server/         # メインエントリーポイント
├── internal/
│   ├── config/        # 設定管理
│   ├── handlers/      # HTTPハンドラー（コントローラー層）
│   │   ├── handlers.go # Web UI ハンドラー
│   │   └── api.go     # REST API ハンドラー
│   ├── models/        # データモデル（エンティティ）
│   ├── repository/    # データアクセス層
│   └── service/       # ビジネスロジック層
├── web/
│   ├── templates/     # HTMLテンプレート
│   └── static/        # 静的ファイル（CSS/JS）
└── migrations/        # DBマイグレーション
```

## コーディング規約

### Go

- Go標準のフォーマッティング（`go fmt`）を使用
- エラーハンドリングは明示的に行う
- インターフェースを活用した疎結合設計
- コンテキスト（context）を適切に伝播
- ゴルーチンリークに注意

### データベース

- SQLite を使用（外部依存なし）
- マイグレーションは `internal/repository/db.go` で管理
- プリペアドステートメントを使用してSQLインジェクション対策

## API 設計

- REST API は `/api/` プレフィックスを使用
- 標準的な HTTP メソッドとステータスコードを使用
- レスポンスは JSON 形式
- Spring Boot版との互換性を維持

## フロントエンド

- Go標準の `html/template` を使用
- Bootstrap 5.3.0 でスタイリング
- 静的ファイルは埋め込み可能（embed）

## ビルドとデプロイ

### 開発

- `make dev` - ホットリロード付き開発モード
- `make run` - 通常実行
- `make test` - テスト実行

### ビルド

- `make build` - ローカルビルド
- `make build-pi` - Raspberry Pi向けビルド（ARM）
- `make docker-build` - Dockerイメージビルド

### デプロイ

- バイナリサイズ: 約15MB
- メモリ使用量: 30-50MB
- 起動時間: 2-3秒

## パフォーマンス目標

- **メモリ使用量**: 50MB以下
- **起動時間**: 5秒以内
- **レスポンス時間**: 100ms以内
- **同時接続数**: 100以上

## Raspberry Pi対応

- Pi Zero 2W: 快適動作（推奨）
- Pi 3: 快適動作
- Pi 4: 超快適動作
- ARMv6/v7/ARM64対応

## 今後の改善点

- テストカバレッジ80%以上
- CI/CD パイプライン（GitHub Actions）
- プロファイリングツールの統合
- メトリクス収集（Prometheus）