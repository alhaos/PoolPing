# pingpool

Практический пример паттерна **Worker Pool** на Go — параллельная проверка доступности хостов по TCP.

## Как это работает

```
hosts -> [jobs chan] -> [worker 1]  \
                    -> [worker 2]  --> [results chan] -> вывод
                    -> [worker N]  /
```

Хосты отправляются в канал `jobs`, пул воркеров параллельно делает TCP-коннект на порт 80 и отправляет результат в канал `results`. Отмена через `context` останавливает всё чисто.

## Структура

```
pingpool/
├── pingpool.go       # пакет: PingPool + Result
└── cmd/
    └── ping/
        └── main.go   # точка входа
```

## Запуск

```bash
go run ./cmd/ping/
```

Пример вывода:

```
Checking 10 hosts with 5 workers...

HOST                      LATENCY    STATUS
----                      -------    ------
cloudflare.com            12ms       OK
google.com                25ms       OK
github.com                31ms       OK
ya.ru                     55ms       OK
nonexistent.invalid       1ms        FAIL
192.168.1.1               1ms        FAIL

Done in 1.2s
```

## API

```go
func PingPool(ctx context.Context, numWorkers int, jobs <-chan string) <-chan Result
```

```go
type Result struct {
    Host    string
    Latency time.Duration
    Err     error        // nil если хост доступен
}
```

## Особенности

- TCP-коннект на порт 80 — не требует root и работает везде (в отличие от ICMP)
- `Ctrl+C` корректно останавливает пул через `signal.NotifyContext`
- Общее время ≈ время самого медленного хоста, а не сумма всех

## Связанные паттерны

- [Worker Pool (тесты)](../workerpool/) — базовая реализация с `chan int`
