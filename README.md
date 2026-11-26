# drak-home

DRAK: IoT tại nhà 3 tầng (Device ESP-IDF / Edge Go / Central), luật chạy offline trên ESP32. Việc đầu tiên: máy lạnh TCL tự tắt đúng giờ, kể cả khi mất mạng.

Trạng thái: **S0 — khung repo**. Chưa có tính năng dùng được.

## Cấu trúc

| Thư mục | Nội dung |
|---|---|
| `contracts/v1/` | Contract v1 (JSON Schema 2020-12, ví dụ, test vector `conformance/scenarios.json`), chép nguyên từ blueprint v0.1. Edge và device đọc chung, không ghim hash. |
| `edge/` | Edge runtime, module Go (`go.mod` ở `edge/`). |
| `device/` | Firmware ESP-IDF cho ESP32-WROOM-32 (target `esp32`), cùng host test ở `device/test/host/`. |
| `.github/workflows/` | `edge.yml` (đổi `edge/**`, `contracts/**`), `device.yml` (đổi `device/**`, `contracts/**`). |

## Build và test nhanh

```sh
# Edge (Go 1.27+)
cd edge && go build ./... && go vet ./... && go test ./...

# Device firmware (ESP-IDF v5.5, đã source export.sh)
cd device && idf.py set-target esp32 && idf.py build

# Device host test (CMake ≥ 3.21, g++)
cmake -S device/test/host -B device/build-host && cmake --build device/build-host -j
ctest --test-dir device/build-host --output-junit "$PWD/device-host-junit.xml"
```

Chi tiết cho người và agent: [AGENTS.md](AGENTS.md).

## Bảo mật

Repo public. Không commit Wi-Fi, token, API key hay `sdkconfig` có thông tin riêng; cấu hình riêng để trong file bị `.gitignore`.

## Giấy phép

MIT, xem [LICENSE](LICENSE).
