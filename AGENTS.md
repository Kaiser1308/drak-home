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

<!-- apos:begin (managed by `apos init`; edits inside are replaced) -->
This project uses APOS, the output-assurance layer: a change counts as done when the gate says so, not when an agent says so.

- Every code change belongs to a Goal and follows the roles: PM writes the Goal, Lead plans and submits, Tester writes the acceptance tests, Worker writes the code, Reviewer blocks with a red test. Use the matching skill: `apos-pm`, `apos-lead`, `apos-tester`, `apos-worker`, `apos-reviewer`.
- Investigating a bug (reading logs, reproducing, experiments on a scratch branch) needs no Goal. Fixing it does: open a small Goal with `apos bug open` (the Human approves it with `apos bug approve`) whose acceptance test reproduces the bug — red on the current code, locked, then green with the fix.
- Restructuring code that has no tests needs a safety net first: a `characterize` Goal captures what the code does today (`apos goal create --kind characterize --covers <glob>`), then a `refactor` Goal (`--kind refactor --paths <glob>`) restructures only what is captured. Accepting a repo with no tests (`apos setup approve --baseline-missing`) is the Human's decision.
- Upgrading a dependency or patching a security finding is a `maintenance` Goal: `apos maintenance open` (`--security <finding>` to name the finding the patch must remove, `--from-pr <n>` for a Dependabot pull request); the Human approves it with `apos maintenance approve`. It needs no criterion: the gate holds it to the base (no new failure, no new finding, only manifests, lockfiles and the declared paths).
- When production is on fire the Human lets a fix through now with `apos hotfix "<what>" --commit <sha>` (it merges nothing: the Human merges). APOS records it, opens a bugfix Goal that owes the evidence by the deadline (a test red on the commit before the hotfix, green with it) and shows `hotfix chưa có bằng chứng: n` in `apos inbox` until that Goal closes. Plan it with `apos lead propose --bugfix`; the delivery must contain the hotfix commit.
- Work is done only when `apos package submit` (or the integration gate) is green. "Done" in an agent message does not count.
- Do not edit `.apos/`, locked tests or the test infrastructure.
- Do not run Human commands: `goal approve|appetite|select`, `bug approve`, `maintenance approve`, `hotfix`, `outcome close|revise|stop`, `gate confirm`, `gate run --base`, `blocker resolve`, `setup approve`, `exception approve --human`. Do not pass `--on-behalf` unless the Human has clearly said yes somewhere else.
- `apos inbox` shows what is waiting for the Human. The setup is `.apos/project.yaml`; see APOS `docs/setup.md`.
<!-- apos:end -->
