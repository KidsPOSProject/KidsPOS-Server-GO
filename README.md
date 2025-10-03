# KidsPOS Server (Go Version)

子供向けPOSシステムのGo実装版。Spring Boot版と比較して大幅なメモリ削減とパフォーマンス向上を実現。

## 特徴

- **超軽量**: メモリ使用量 30-50MB（Spring Boot版の1/10）
- **高速起動**: 2-3秒で起動（Spring Boot版は30-45秒）
- **Raspberry Pi対応**: Pi Zero 2Wでも快適に動作
- **SQLite内蔵**: 外部データベース不要
- **クロスプラットフォーム**: Linux/macOS/Windows/ARM対応
- **APK管理機能**: Androidアプリのバージョン管理とOTA配信
- **完全なCRUD**: 全エンティティで作成・読取・更新・削除をサポート

## 比較表

| 項目                   | Spring Boot版 | Go版         |
|----------------------|--------------|-------------|
| メモリ使用量               | 512MB+       | **30-50MB** |
| 起動時間                 | 30-45秒       | **2-3秒**    |
| バイナリサイズ              | 50MB (JAR)   | **15MB**    |
| Raspberry Pi Zero 2W | ❌ 動作困難       | ✅ 快適動作      |
| Raspberry Pi 3       | ⚠️ 動作するが重い   | ✅ 快適動作      |
| Raspberry Pi 4       | ✅ 動作可能       | ✅ 超快適       |

## クイックスタート

### 必要環境

- Go 1.21以上
- Make（オプション）

### インストール・実行

```bash
# 依存関係のダウンロード
make deps

# ビルド
make build

# 実行
make run

# または直接実行
go run cmd/server/main.go
```

アプリケーションは http://localhost:8080 で起動します。

## 新機能（v2.0）

### 1. Staff/Store の削除・更新機能
- 物理削除を実装（論理削除ではなくデータベースから完全削除）
- 外部キー制約の自動チェック（販売履歴がある場合は削除を拒否）
- Web UIおよびREST API両方に対応

### 2. APKバージョン管理システム
- **ファイルアップロード**: multipart/form-data形式でAPKファイルをアップロード
- **バージョン管理**: セマンティックバージョニングとバージョンコードの管理
- **OTA配信**: Androidアプリから最新版の確認とダウンロード
- **リリースノート**: 各バージョンの変更内容を記録
- **ファイルストレージ**: ローカルファイルシステム（`./uploads/apk/`）
- **バージョン無効化**: 古いバージョンを配信停止（物理削除せず無効化）
- **Web UI**: アップロード、一覧表示、ダウンロードをブラウザから操作可能

### 3. 包括的なテストスイート
- **リポジトリ層**: データベース操作の全テスト
- **サービス層**: ビジネスロジックとファイル処理のテスト
- **外部キー制約**: 参照整合性の検証テスト
- テストフレームワーク: `testify` を使用

## Raspberry Pi向けビルド

```bash
# 全Raspberry Pi向けビルド
make build-pi

# 個別ビルド
# Pi 4/5 (64-bit)
GOOS=linux GOARCH=arm64 go build -o kidspos-arm64 cmd/server/main.go

# Pi 3 (32-bit ARMv7)
GOOS=linux GOARCH=arm GOARM=7 go build -o kidspos-armv7 cmd/server/main.go

# Pi Zero 2W (32-bit ARMv6)
GOOS=linux GOARCH=arm GOARM=6 go build -o kidspos-armv6 cmd/server/main.go
```

## Raspberry Piへのデプロイ

### 1. バイナリを転送

```bash
# Pi 4/5の場合
scp dist/kidspos-arm64 pi@raspberrypi.local:/home/pi/

# Pi 3の場合
scp dist/kidspos-armv7 pi@raspberrypi.local:/home/pi/

# 実行権限を付与
ssh pi@raspberrypi.local "chmod +x /home/pi/kidspos-arm*"
```

### 2. systemdサービス設定

```bash
# サービスファイル作成
sudo nano /etc/systemd/system/kidspos.service
```

```ini
[Unit]
Description=KidsPOS Go Server
After=network.target

[Service]
Type=simple
User=pi
WorkingDirectory=/home/pi
ExecStart=/home/pi/kidspos-arm64
Restart=always
Environment="PORT=8080"

[Install]
WantedBy=multi-user.target
```

```bash
# サービス有効化と起動
sudo systemctl daemon-reload
sudo systemctl enable kidspos
sudo systemctl start kidspos
```

## Docker実行

```bash
# Dockerイメージビルド
make docker-build

# 実行
make docker-run
```

## 開発

### 開発モード（ホットリロード）

```bash
# airをインストールして開発モード起動
make dev
```

### テスト実行

```bash
# 全テスト実行
make test

# または
go test ./...

# カバレッジ付きテスト
go test ./... -cover

# 詳細表示
go test ./... -v
```

実装済みのテスト:
- リポジトリ層テスト（Staff/Store/APK の CRUD操作）
- サービス層テスト（ビジネスロジックとファイル操作）
- 外部キー制約の検証テスト

### コードフォーマット

```bash
make fmt
```

## プロジェクト構造

```
KidsPOS-Server-GO/
├── cmd/
│   └── server/
│       └── main.go         # エントリーポイント
├── internal/
│   ├── config/            # 設定
│   ├── handlers/          # HTTPハンドラー
│   ├── models/            # データモデル
│   ├── repository/        # データアクセス層
│   └── service/           # ビジネスロジック
├── web/
│   ├── templates/         # HTMLテンプレート
│   └── static/            # 静的ファイル
├── migrations/            # DBマイグレーション
├── Makefile              # ビルドスクリプト
├── go.mod                # Go依存関係
└── README.md
```

## 環境変数

```bash
PORT=8080                           # サーバーポート
DATABASE_PATH=./kidspos.db         # SQLiteファイルパス
RECEIPT_PRINTER_HOST=localhost     # レシートプリンタホスト
RECEIPT_PRINTER_PORT=9100          # レシートプリンタポート
QR_CODE_SIZE=200                   # QRコードサイズ
ALLOWED_IP_PREFIX=192.168.         # 許可IPプレフィックス
```

## API エンドポイント

### Web UI

- `GET /` - ホーム
- `GET /items` - 商品一覧
- `GET /sales` - 販売一覧
- `GET /stores` - 店舗一覧
- `GET /staffs` - スタッフ一覧
- `GET /settings` - 設定
- `GET /reports/sales` - 売上レポート
- `GET /apk` - APKバージョン一覧
- `GET /apk/upload` - APKアップロードページ
- `POST /apk/upload` - APKアップロード処理

### REST API

#### 商品 (Items)
- `GET /api/items` - 商品一覧取得
- `GET /api/items/:id` - 商品詳細取得
- `POST /api/items` - 商品作成
- `PUT /api/items/:id` - 商品更新
- `DELETE /api/items/:id` - 商品削除

#### 販売 (Sales)
- `GET /api/sales` - 販売一覧取得
- `GET /api/sales/:id` - 販売詳細取得
- `POST /api/sales` - 販売登録

#### 店舗 (Stores)
- `GET /api/stores` - 店舗一覧取得
- `GET /api/stores/:id` - 店舗詳細取得
- `POST /api/stores` - 店舗作成
- `PUT /api/stores/:id` - 店舗更新
- `DELETE /api/stores/:id` - 店舗削除（物理削除、販売履歴がある場合はエラー）

#### スタッフ (Staffs)
- `GET /api/staffs` - スタッフ一覧取得
- `GET /api/staffs/:id` - スタッフ詳細取得
- `POST /api/staffs` - スタッフ作成
- `PUT /api/staffs/:id` - スタッフ更新
- `DELETE /api/staffs/:id` - スタッフ削除（物理削除、販売履歴がある場合はエラー）

#### 設定 (Settings)
- `GET /api/settings` - 設定一覧取得
- `PUT /api/settings/:key` - 設定更新

#### レポート (Reports)
- `GET /api/reports/sales` - 売上データ取得
- `GET /api/reports/sales/excel` - 売上データExcelダウンロード

#### APKバージョン管理 (APK Versions)
- `GET /api/apk/version/latest` - 最新APKバージョン取得
- `GET /api/apk/version/check?currentVersionCode=X` - アップデート確認
- `GET /api/apk/version/all` - 全APKバージョン取得
- `GET /api/apk/download/:id` - APKファイルダウンロード（ID指定）
- `GET /api/apk/download/latest` - 最新APKファイルダウンロード
- `POST /api/apk/upload` - APKファイルアップロード
- `DELETE /api/apk/version/:id` - APKバージョン削除（物理削除）
- `PUT /api/apk/version/:id/deactivate` - APKバージョン無効化

## トラブルシューティング

### Raspberry Piでメモリ不足の場合

```bash
# スワップを無効化（SDカード寿命延長）
sudo dphys-swapfile swapoff
sudo systemctl disable dphys-swapfile

# 不要なサービスを停止
sudo systemctl disable bluetooth
sudo systemctl disable avahi-daemon
```

### ポート使用中エラー

```bash
# 8080ポートを使用中のプロセスを確認
lsof -i:8080

# プロセスを停止
kill -9 [PID]
```

## ライセンス

MIT License

## 貢献

プルリクエスト歓迎です！

## サポート

問題がある場合は、GitHubのIssuesでお知らせください。