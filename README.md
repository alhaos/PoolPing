# PoolPing

Практический пример паттерна **Worker Pool** на Go — параллельная проверка
доступности хостов через ICMP.

## Как это работает

```text
hosts -> [jobs chan] -> [worker 1]  \
                     -> [worker 2]  --> [results chan] -> вывод
                     -> [worker N]  /
```

Хосты отправляются в канал `jobs`, пул воркеров параллельно отправляет ICMP echo
request и пишет результат в канал `results`. Отмена через `context` останавливает
всё чисто.

## Структура

```text

PoolPing/
├── cmd/
│   └── PoolPing/
│       └── main.go           # точка входа
├── config/
│   └── config.yml            # список хостов и настройки
├── internal/
│   ├── config/
│   │   └── config.go         # загрузка конфига через cleanenv
│   └── pingpool/
│       └── pingpool.go       # PingPool + Result
├── go.mod
└── go.sum
```

## Запуск

> Требуются права администратора (Windows) или root (Linux) для ICMP сокетов.

```bash
# с дефолтным конфигом
go run ./cmd/PoolPing/

# с кастомным конфигом
go run ./cmd/PoolPing/ -config path/to/config.yml
```

Пример вывода:

```text

Results:
  - gitlab.com                                         48ms
  - github.com                                         89ms
  - ya.ru                                              51ms
  - amazon.com                                         no reply
  - nonexistent.invalid                                lookup nonexistent.invalid: no such host
Total host processed: 42
  - Alive: 21
  - Dead: 21

```

## Конфиг

```yaml
workers: 10       # количество воркеров
timeoutMs: 3000   # таймаут ICMP в миллисекундах

hosts:
  - google.com
  - github.com
  - 8.8.8.8
  # ...
```

## API

```go

func PingPool(ctx context.Context, numWorkers int, timeout time.Duration, jobs <-chan string) <-chan Result

```

```go

type Result struct {
    Host    string
    Latency time.Duration
    Err     error        // nil если хост доступен
}

```

## Особенности

- Настоящий ICMP ping через [pro-bing](https://github.com/prometheus-community/pro-bing)
- `Ctrl+C` корректно останавливает пул через `signal.NotifyContext`
- Результаты выводятся сразу по мере готовности — порядок вывода = порядок завершения воркеров
- Общее время ≈ время самого медленного хоста, а не сумма всех
