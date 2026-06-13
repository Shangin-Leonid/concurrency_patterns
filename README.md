<!---
![Golang](https://img.shields.io/badge/%20-navy?style=flat&logo=Go&logoSize=auto)
![Golang](https://img.shields.io/badge/%20-00ADD8?style=flat&logo=Go&logoSize=auto&logoColor=white)
![Golang](https://img.shields.io/badge/%20-white?style=flat&logo=Go&logoSize=auto&logoColor=00ADD8)
-->

# ![Golang](https://img.shields.io/badge/%20-00ADD8?style=flat&logo=Go&logoSize=auto&logoColor=E7FEFB) Go Concurrency Patterns

![Go version](https://img.shields.io/github/go-mod/go-version/Shangin-Leonid/concurrency_patterns?style=flat&logo=Go&logoColor=E7FEFB&logoSize=auto&labelColor=00ADD8&color=00ADD8)
![Makefile](https://img.shields.io/static/v1?label=&message=Makefile&style=flat&logo=CMake&color=FE7A16)
![Linux](https://img.shields.io/static/v1?label=&message=Linux&style=flat&logo=linux&color=2B822F)
![Windows](https://img.shields.io/static/v1?label=&message=Windows&style=flat&logo=windows&color=4169e1)

This repository contains practical implementations of key concurrency and multithreading patterns in Golang. These design patterns help efficiently manage goroutines, safely transfer data via channels, and prevent memory leaks.

## :white_check_mark: Implemented patterns and tips

Some patterns have been taken from  "Concurrency in Go" by Katherine Cox-Buday.

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

## :rocket: Getting Started

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

### 3. Delete all rubish

```bash
make clean
```

---
