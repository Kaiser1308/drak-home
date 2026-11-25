# AGENTS.md

Hướng dẫn build và test cho agent và người làm trong drak-home. Mọi lệnh chạy từ gốc repo và tự đủ (không phụ thuộc thư mục build có sẵn), nên chạy được trong worktree mới sạch.

## Quy tắc chung

- Repo **public**: không commit secret (Wi-Fi, token Telegram, API key), không commit `device/sdkconfig`.
- `contracts/v1/` là hợp đồng chung giữa edge và device. Đổi contract là việc cần người duyệt; khi đổi phải chạy lại cả test edge lẫn device.
- `edge/go.mod`, `edge/go.sum`, `device/**/idf_component.yml`, `device/dependencies.lock` cũng cần người duyệt khi đổi.

## contracts/

- Nội dung: `schemas/*.json` (Draft 2020-12), `examples/*.json`, `conformance/scenarios.json` (test vector).
- Kiểm schema: test Go trong edge, `cd edge && go test ./internal/contract/ -run TestAcceptance -v`. Test biên dịch mọi schema, validate mọi ví dụ có schema tương ứng, và báo lỗi nếu có ví dụ mới chưa được khai báo trong `exampleSchema` hoặc `skippedExamples` (`edge/internal/contract/contract_test.go`).
- Thư mục contract được tìm bằng cách đi ngược lên từ thư mục hiện tại; đặt `DRAK_CONTRACTS_DIR` để chỉ định khác.

## edge/ (Go)

Yêu cầu: Go ≥ 1.27, `gotestsum` (`go install gotest.tools/gotestsum@v1.13.0`) nếu cần JUnit.

```sh
cd edge && go vet ./...                                   # L0
cd edge && go build ./...
cd edge && go test ./...                                  # L1
cd edge && gotestsum --junitfile "$PWD/edge-junit.xml" -- ./...
cd edge && go test ./... -run 'TestAcceptance'            # acceptance
```

Test tên `TestAcceptance*` là test nghiệm thu (contract và hành vi người dùng thấy).

## device/ (ESP-IDF C++)

Yêu cầu: ESP-IDF v5.5 (CI dùng image `espressif/idf:v5.5.6`), target `esp32` (ESP32-WROOM-32).

```sh
. "$IDF_PATH/export.sh"
cd device && idf.py set-target esp32 && idf.py build
```

Không có ESP-IDF trên máy thì chạy qua Docker:

```sh
docker run --rm --user "$(id -u):$(id -g)" -e HOME=/tmp -v "$PWD:/project" -w /project/device \
  espressif/idf:v5.5.6 bash -c 'idf.py set-target esp32 && idf.py build'
```

Code thuần C++ (không gọi API ESP-IDF) đặt trong `device/components/drak_core/` để test được trên máy tính.

### Host test (máy tính, không cần ESP-IDF)

Yêu cầu: CMake ≥ 3.21, g++ hỗ trợ C++17, mạng lần đầu (FetchContent tải GoogleTest 1.17.0, có kiểm SHA-256).

```sh
cmake -S device/test/host -B device/build-host
cmake --build device/build-host -j
ctest --test-dir device/build-host --output-on-failure --output-junit "$PWD/device-host-junit.xml"
```

`--output-junit` hiểu đường dẫn tương đối theo thư mục build, nên luôn truyền đường dẫn tuyệt đối.

Test trên phần cứng thật (phát IR, cảm biến) là kiểm tra thủ công, ghi kết quả vào PR.

## CI

- `.github/workflows/edge.yml`: gofmt, build, vet, test (JUnit artifact). Chạy khi đổi `edge/**` hoặc `contracts/**`.
- `.github/workflows/device.yml`: `idf.py build` trong `espressif/idf:v5.5.6` và host test (JUnit artifact). Chạy khi đổi `device/**` hoặc `contracts/**`.
