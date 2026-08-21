# Crypto Berserker — Multi-Exchange Cryptocurrency Trading System in Go

**Go · REST APIs · WebSocket · Concurrency · MySQL · Automated Trading · Exchange Integration**

Developed and maintained from **2019 to 2024** as a personal quantitative trading project — designed, built, and operated live on AWS by a single developer. ~13,500 lines of Go across 569 commits.

## What it does

Crypto Berserker connects to **10 cryptocurrency exchanges** through a unified abstraction layer and runs multiple automated trading strategies on top of it:

- **Cross-exchange arbitrage** (`strategy/exchanger`) — maintains real-time futures order books for Binance and Bybit over WebSocket and captures price divergence between venues
- **Triangular arbitrage** (`strategy/arbitrage/Triangular`) — exploits pricing inconsistencies across three pairs within a single exchange
- **Funding-rate arbitrage** (`strategy/arbitrage/FundingRate`) — hedged spot/perpetual positions to harvest funding payments
- **New-listing sniper** (`strategy/buyer-first-online*`) — low-latency buyer for the first seconds of new token listings
- **Market maker** (`strategy/maker`) and **automated margin lending** (`strategy/Lender`)

## Exchange integrations

Each integration is hand-written against the exchange's official API — no third-party SDK:

| Exchange | Scope |
|---|---|
| Binance | Spot REST + WebSocket |
| Binance Futures (USD-M) | REST + WebSocket user/market streams |
| Bybit | Spot/derivatives REST + WebSocket |
| Bybit Linear (v5) | USDT perpetuals REST + WebSocket |
| OKEx | Swap REST + WebSocket |
| FTX / FTX OTC | REST + WebSocket *(historical — exchange defunct since Nov 2022)* |
| MaiCoin MAX | REST + WebSocket order book |
| AscendEX (BitMax) | REST |
| CoinEx | REST |

Common concerns are factored out: HMAC-SHA256 request signing, timestamp/`recvWindow` handling, order placement with retry, and order-book stream normalization.

## Architecture

```
strategy/           Trading strategies (venue-agnostic)
exchange/           Core abstraction: Exchange interface, order-booker engine,
                    wallets, fees, HMAC signing helpers
exchange-list/      Per-exchange implementations (REST + WebSocket)
ksql/ + sql/        MySQL persistence for trades and quotes
message_tool/       Telegram bot + e-mail alerting
jmath / jtime / log Shared utilities
```

Design principles:

- **One `Exchange` interface, many venues** — strategies never touch exchange-specific code; adding a venue never touches strategy code
- **Local order books from WebSocket diff streams** — add/update/remove deltas applied concurrently per market, with automatic reconnection and resync on drop
- **Goroutine-per-connection concurrency** — each market stream and each strategy loop runs independently; shared wallet state is lock-protected
- **INI-driven configuration** — global strategy parameters with per-symbol overrides, no recompile needed for tuning

## Engineering notes from production

Highlights from the development log (`feature.md`):

- Profiled the order path end-to-end: isolated HTTP POST latency (~50 ms, verified identical between `net/http` and `fasthttp`), confirmed in-process overhead under 100 ns, then cut Binance Futures WebSocket update interval from 250 ms to 100 ms
- Diagnosed WebSocket "instability" that turned out to be AWS t3a.micro CPU-credit throttling — the disconnects were the symptom, not the cause; resolved by instance resizing
- Survived the FTX collapse (Nov 2022): migrated live strategies from FTX to Binance Futures + Bybit linear perpetuals
- Hardened order handling: retry with exchange-specific rate-limit backoff, cancel-all safety on shutdown, millisecond-timestamped fill reporting

## Notes

- All API credentials have been scrubbed from the repository and its history (`***REMOVED***` placeholders); the referenced private `setting` package was never committed
- Commit messages are mostly in Traditional Chinese — this was a personal project
- No longer actively maintained; published as a portfolio piece. Nothing here is financial advice.
