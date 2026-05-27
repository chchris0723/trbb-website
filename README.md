# TRBB 鐵人三項運動社團

> 挑戰極限，超越自我。

一站式鐵人三項社團管理平台，涵蓋會員系統、賽事報名、商品銷售、二手交換、訓練日記、團體報名、賽事交通，並整合金流與 Garmin Connect。

---

## 技術架構

| 層級 | 技術 |
|------|------|
| 前端 Web | Vue 3 + Vite + Pinia |
| 前端 App (Mobile) | Vue 3 + Vite + Pinia |
| 後台 Admin | Vue 3 + Vite + Pinia |
| 後端 | Go (Gin) |
| 資料庫 | MySQL 8.0 |
| 快取 | Redis 7 |
| 物件儲存 | MinIO |
| 反向代理 | Nginx 1.25 |
| 容器化 | Docker + Docker Compose |

---

## 快速開始

### 前置需求
- Docker & Docker Compose
- Node.js 18+
- Go 1.22+（本地開發）

### 初始化（首次）

```bash
# 1. 複製環境設定
cp .env.example .env            # 修改密碼等設定
cp backend/.env.example backend/.env   # 修改第三方金鑰

# 2. 賦予執行權限
chmod +x service.sh

# 3. 初始化並啟動
./service.sh init
```

### 常用指令

```bash
./service.sh init       # 首次初始化：建置前端 + 啟動所有服務
./service.sh build      # 建置前端靜態檔案 + 後端 Docker image
./service.sh start      # 啟動所有服務
./service.sh stop       # 停止所有服務
./service.sh restart    # 重啟所有服務
./service.sh rebuild    # 完整重建：停止 → 建置前端 → 重建後端 → 啟動
./service.sh status     # 查看容器狀態
./service.sh logs               # 查看所有日誌
./service.sh logs backend       # 查看後端日誌
./service.sh logs nginx         # 查看 Nginx 日誌
```

---

## 服務網址

| 服務 | 網址 |
|------|------|
| 前台 Web | http://trbbtw.com |
| 手機前台 | http://m.trbbtw.com |
| 後台管理 | http://admin.trbbtw.com |
| 圖片 CDN | http://images.trbbtw.com |
| MinIO Console | http://localhost:9001 |
| Backend Health | http://localhost:8080/health |

> 本地開發請將以上域名加入 `/etc/hosts` 指向 `127.0.0.1`

---

## API 路由規劃

| 前綴 | 說明 |
|------|------|
| `/v1/api/*` | 前台 API（需會員 JWT） |
| `/v1/admin/*` | 後台 API（需管理員 JWT，role=9） |
| `/v1/third/*` | 第三方回呼（金流、Garmin、OAuth）|

---

## 主要功能模組

- **會員系統** — 註冊、登入、Email 驗證、OAuth（Google / Facebook / LINE）
- **管理員審核** — 會員申請審核、分級制度（一般 / 進階 / 管理員）
- **賽事報名** — 個人 / 團體報名、確認、退款、號碼布管理
- **賽事交通** — 包車預訂、座位管理
- **商品銷售** — 商品管理、購物車、訂單流程
- **二手交換** — 刊登、洽詢、交易確認
- **訓練日記** — 手動記錄、Garmin 自動同步、數據分析
- **金流整合** — ECPay / Stripe / LINE Pay（預留接口）
- **Garmin 串接** — OAuth 1.0a 授權、活動推播、健康數據
- **通知系統** — 站內通知、Email（SendGrid）、SMS（Every8D）

---

## 第三方整合說明

所有第三方設定均在 `backend/.env` 中配置：

### 金流
```
PAYMENT_PROVIDER=ecpay        # ecpay | stripe | linepay
PAYMENT_MERCHANT_ID=...
PAYMENT_HASH_KEY=...
PAYMENT_HASH_IV=...
```

### Garmin Connect API
```
GARMIN_CLIENT_ID=...
GARMIN_CLIENT_SECRET=...
GARMIN_REDIRECT_URI=https://trbbtw.com/v1/third/garmin/callback
```

### OAuth
```
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
FB_APP_ID=...
FB_APP_SECRET=...
LINE_CHANNEL_ID=...
LINE_CHANNEL_SECRET=...
```

---

## 資料庫

- 引擎：MySQL 8.0，字符集 `utf8mb4`
- Schema 初始化：`mysql/init/01_schema.sql`（Docker 啟動時自動執行）
- 時區：`Asia/Taipei`

預設管理員帳號（首次啟動後請立即修改密碼）：
```
Email: admin@trbbtw.com
```

---

## 日誌

- 後端：JSON 格式，每日輪轉，存於 `logs/` volume（`trbb-YYYY-MM-DD.log`）
- Nginx：JSON 格式，存於 `nginx/logs/`

---

## 目錄結構

```
trbb/
├── backend/                 # Go 後端
│   ├── cmd/server/          # 程式入口
│   ├── internal/
│   │   ├── api/             # /v1/api 路由
│   │   ├── admin/           # /v1/admin 路由
│   │   ├── third/           # /v1/third 路由
│   │   ├── middleware/      # JWT、CORS、Logger
│   │   └── config/          # 環境設定
│   └── pkg/
│       ├── database/        # MySQL
│       ├── cache/           # Redis
│       ├── storage/         # MinIO
│       └── logger/          # Zap JSON logger
├── frontend/
│   ├── web/                 # 前台 Web (Vue 3)
│   ├── app/                 # 手機前台 (Vue 3)
│   └── admin/               # 後台管理 (Vue 3)
├── mysql/
│   └── init/                # SQL Schema
├── nginx/
│   ├── nginx.conf           # 主設定
│   └── conf.d/              # 虛擬主機設定
├── docker-compose.yml
├── service.sh               # 管理腳本
└── .gitignore
```

---

## License

© 2024 TRBB 鐵人三項運動社團. All rights reserved.
