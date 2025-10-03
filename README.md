# KidsPOS Server (Go Version)

子供向けPOSシステムのGo実装版。Spring Boot版と比較して大幅なメモリ削減とパフォーマンス向上を実現。

## 特徴

- **超軽量**: メモリ使用量 30-50MB（Spring Boot版の1/10）
- **高速起動**: 2-3秒で起動（Spring Boot版は30-45秒）
- **Raspberry Pi対応**: Pi Zero 2Wでも快適に動作
- **SQLite内蔵**: 外部データベース不要
- **クロスプラットフォーム**: Linux/macOS/Windows/ARM対応

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
make test
```

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

### REST API

- `GET /api/items` - 商品一覧取得
- `POST /api/items` - 商品作成
- `PUT /api/items/:id` - 商品更新
- `DELETE /api/items/:id` - 商品削除
- `POST /api/sales` - 販売登録
- `GET /api/reports/sales` - 売上データ取得

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