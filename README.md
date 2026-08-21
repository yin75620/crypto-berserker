# Crypto Berserker — Multi-Exchange Cryptocurrency Trading System in Go

**Go · REST APIs · WebSocket · Concurrency · MySQL · AWS · Automated Trading · Exchange Integration**

Developed and maintained from **2019 to 2024** as an independent automated trading system — designed, built, and operated live on AWS with real capital. Approximately **13,500 lines of Go across 550+ commits**.

## What it does

Crypto Berserker supports **nine exchange integrations** through a unified abstraction layer, providing a common foundation for several automated trading strategies:

* **Cross-exchange arbitrage** (`strategy/exchanger`) — maintains real-time futures order books for Binance and Bybit over WebSocket and identifies price divergence between venues
* **Triangular arbitrage** (`strategy/arbitrage/Triangular`) — evaluates pricing inconsistencies across three trading pairs within a single exchange
* **Funding-rate arbitrage** (`strategy/arbitrage/FundingRate`) — manages hedged spot/perpetual positions designed to capture funding-rate opportunities
* **New-listing execution strategy** (`strategy/buyer-first-online*`) — low-latency order execution during the first seconds of newly listed assets
* **Market making** (`strategy/maker`) and **automated margin lending** (`strategy/Lender`)

## Exchange integrations

Each integration was implemented directly against the exchange's official REST and WebSocket APIs without relying on third-party exchange SDKs.

| Exchange / Integration  | Scope                                                             |
| ----------------------- | ----------------------------------------------------------------- |
| Binance                 | Spot REST + WebSocket                                             |
| Binance Futures (USD-M) | REST + WebSocket user/market streams                              |
| Bybit                   | Spot/derivatives REST + WebSocket                                 |
| Bybit Linear (v5)       | USDT perpetuals REST + WebSocket                                  |
| OKEx                    | Swap REST + WebSocket                                             |
| FTX / FTX OTC           | REST + WebSocket *(historical — exchange defunct since Nov 2022)* |
| MaiCoin MAX             | REST + WebSocket order book                                       |
| AscendEX (BitMax)       | REST                                                              |
| CoinEx                  | REST                                                              |

Common exchange concerns are factored into reusable components, including:

* HMAC-SHA256 request signing
* Timestamp and `recvWindow` handling
* Order placement and retry logic
* Exchange-specific rate-limit handling
* WebSocket reconnection
* Order-book stream normalization
* Wallet and account state management

## Architecture

```text
strategy/           Trading strategies and execution logic
exchange/           Core abstractions: Exchange interfaces, order-book engine,
                    wallets, fees, authentication helpers
exchange-list/      Exchange-specific REST and WebSocket implementations
ksql/ + sql/        MySQL persistence for trades and market data
message_tool/       Telegram and e-mail notifications
jmath / jtime / log Shared utilities
```

### Design principles

* **Unified exchange abstraction** — trading strategies are largely isolated from exchange-specific implementations, making it easier to add or migrate venues
* **Local order books from WebSocket streams** — market-data deltas are processed into local books with reconnection and synchronization handling
* **Concurrent execution with goroutines** — market streams and strategy loops run independently, while shared wallet state is protected against concurrent access
* **Configuration-driven strategies** — global strategy parameters and per-symbol overrides can be adjusted without recompiling the application
* **Operational fault handling** — retries, exchange-specific errors, disconnections, and shutdown behavior are handled explicitly because the system operated with live capital

## Engineering notes from production

Highlights from the development log (`feature.md`):

* **Profiled the end-to-end order path** — measured HTTP request latency at approximately 50 ms and determined that local processing overhead was negligible relative to network latency; subsequently reduced Binance Futures WebSocket update intervals from 250 ms to 100 ms
* **Diagnosed an infrastructure-level WebSocket issue** — intermittent disconnects were traced to AWS `t3a.micro` CPU-credit throttling rather than the WebSocket implementation itself, and were resolved by resizing the instance
* **Responded to the FTX collapse in November 2022** — migrated live trading strategies away from FTX to Binance Futures and Bybit linear perpetuals
* **Hardened order execution** — added retry behavior, exchange-specific rate-limit backoff, cancel-all safety mechanisms on shutdown, and millisecond-timestamped trade reporting
* **Maintained compatibility with evolving exchange APIs** — including migration to newer Bybit API versions and adjustments to exchange-specific order and market-data behavior

## Operational context

This was not a simulated trading demo. The system was deployed on AWS and operated with **real capital** over multiple years, so failures in order execution, exchange connectivity, or state management had real financial consequences. Many of the reliability, recovery, and operational mechanisms described above were introduced in response to issues encountered during live operation.

## Repository notes

* All API credentials have been scrubbed from the repository and its Git history (`***REMOVED***` placeholders)
* The private `setting` package referenced by the application was never committed
* Most historical commit messages are written in Traditional Chinese
* The project is no longer actively maintained and is published as a portfolio piece
* Nothing in this repository constitutes financial advice
