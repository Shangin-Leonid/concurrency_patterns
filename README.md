# ⚙️ Go Concurrency Patterns

This repository contains practical implementations of key concurrency and multithreading patterns in Go. These design patterns help efficiently manage goroutines, safely transfer data via channels, and prevent memory leaks.

## 🛠️ Implemented patterns and tips

The project features the following patterns and mechanisms (each implemented in a dedicated file):

🔹 **Lexical Confinement** (`lexical_confinement.go`) — restricting channel access to a specific lexical scope to ensure compile-time thread safety.

🔹 **Generator** (`generator.go`) — generating an output channel and writting goroutine.

🔹 **Repeat-take** (`repeat_take.go`) — generating an infinite source of data and taking a finite sequence of it.

🔹 **Or-Channel** (`or_channel.go`) — combining multiple read channels into one that closes as soon as any of its component channels close.

🔹 **Or-Done-Loop** (`or_done_loop.go`) — a clean wrapper pattern for safely reading from channel while respecting cancellation via `context`.

🔹 **Select Priority** (`select_priority.go`) — enforcing execution priority among multiple cases within a `select` block.

🔹 **Pipeline** (`pipeline.go`) — pipelined processing of data streams.

🔹 **Channel Filter** (`channel_filter.go`) — filtering data flowing through channels based on a predicate.

🔹 **Fan-Out / Fan-In** (`fan_out.go`, `fan_in.go`, `fan_out_fan_in.go`) — distributing resource-intensive tasks across multiple goroutines and multiplexing the results into a single channel.

🔹 **Tee-Channel** (`tee_channel.go`) — splitting a single input channel into two independent output streams (similar to the `tee` command-line utility).

🔹 **Bridge-Channel** (`bridge_channel.go`) — flattening a sequence of channels into a single stream of values.

🔹 **Context deadline** (`context_deadline.go`) — evident example of usage contexts with timeouts and deadlines.

---

## 🚀 Getting Started

To run these examples locally, you will need [Go](https://go.dev) installed and optionally the `make` utility.

### 1. Clone the Repository

```bash
git clone https://github.com/Shangin-Leonid/concurrency_patterns
cd concurrency_patterns
```

### 2. Running the Code

A `makefile` is provided in the root directory to quickly execute the demonstration code inside `main.go`:

```bash
make run ARG=<pattern to be tested>
```

If you do not have `make` installed, use the standard Go command:

```bash
go run main.go <pattern to be tested>
```

---
