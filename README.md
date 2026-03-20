# Fehu ᚠ

A CLI-based personal accounting system written in Go. Fehu lets you manage accounts, record transactions, tag entries, and run custom calculations — all from an interactive command-line interface backed by a local SQLite database.
this program is licensed under MIT license

---

## Summary

Fehu is a CLI accounting system written in Go.

Key technical highlights:

- Custom DFA-based command parser
- AST expression engine for financial queries
- SQLite persistence layer
- Modular architecture (cli / core / ast)

The goal of the project was to design a flexible financial analysis tool entirely from scratch.

---

## Features

### Account Management
Organize your finances in a hierarchical account tree using colon-separated names (e.g., `income:salary`, `expense:transport`). Accounts support optional descriptions and can be searched by name, description, or parent prefix.

### Transaction Recording
Log financial transactions as a set of account entries. Each entry specifies a direction (`<` for inflow, `>` for outflow) and an amount. Fehu enforces that every transaction is **balanced** — the sum of all entries must equal zero.

```
new txn income:salary<500000;expense:food>200000;asset:bank<-300000
```

Transactions carry an optional description and timestamp, and can be queried by ID, description, or time range.

### 🏷️ Auto-tagging
Write `#tagname` anywhere in an account or transaction description and Fehu automatically creates the tag and links it — no extra command needed.

```
new acc -d="monthly rent #housing #fixed" expense:rent
new txn -d="grocery run #food" expense:food<30000;asset:bank>30000
```

Tags can also be created and managed manually, and used to filter transaction tables in the calculator.

### 🧮 Expression Calculator
An AST-based expression engine lets you query and aggregate saved transaction tables interactively.

**Save a query as a named table:**
```
get txn -save=monthly
```

**Aggregate over an account name pattern** (supports `*` and `?` wildcards):
```
calc sum(acc(monthly, "expense*"))
calc avg(acc(monthly, "income:salary"))
```

**Filter by time range:**
```
calc between(monthly, "2024-01-01;00:00:00", "2024-12-31;23:59:59")
```

**Set operations on tables:**
```
calc union(tableA, tableB)
calc intersect(tableA, tableB)
calc xor(tableA, tableB)
```

**Filter by tag:**
```
calc atag(monthly, "food")     # amounts from accounts tagged #food
calc ttag(monthly, "holiday")  # transactions tagged #holiday
```

**Built-in functions:**

| Function | Description |
|----------|-------------|
| `sum(nums...)` | Sum of numbers |
| `avg(nums...)` | Average |
| `min(nums...)` / `max(nums...)` | Min / Max |
| `count(list\|table)` | Number of items or transactions |
| `acc(table, pattern)` | Record amounts for accounts matching the wildcard pattern |
| `atag(table, tag)` | Record amounts from accounts with a given tag |
| `ttag(table, tag)` | Transactions with a given tag |
| `between(table, start, end)` | Filter transactions by time period |
| `union(tables...)` | Set union of tables |
| `intersect(tables...)` | Set intersection of tables |
| `xor(tables...)` | Symmetric difference of tables |

**Operators:** `+ - * /`, `== != < > <= >=`, `&& || !`

**Define reusable variables:**
```
def salary acc(monthly, "income:salary")
calc sum(salary)
```

The built-in variable `__all__` always refers to every transaction in the database.

### ⚡ Benchmark Mode
Run with `-b` to print execution time after every command — useful for profiling.

---

## Installation Guide

### Prerequisites
- [Go](https://go.dev/) 1.18+

| Platform | Requirement |
|----------|-------------|
| Windows | GCC (via MSYS2 or TDM-GCC) |
| macOS | Xcode Command Line Tools (`xcode-select --install`) |
| Linux | `gcc` (`apt install gcc` / `yum install gcc`) |

### Option 1 — Build from source

```bash
git clone https://github.com/pilboy97/fehu.git
cd fehu
go build -o fehu ./main/
```

Windows:
```powershell
go build -o fehu.exe .\main\
```

### Option 2 — Download prebuilt binary

Download the binary for your platform from GitHub Releases:
```
https://github.com/pilboy97/fehu/releases
```

---

## Usage

### Start the REPL

```bash
./fehu
```

### Open or create a database

```
open mybudget
```

This creates (or opens) `mybudget.db` in the current directory. The schema is created automatically on first open.

### Command-line flags

| Flag | Description |
|------|-------------|
| `-d <n>` | Open a database on startup |
| `-c "<command>"` | Execute a single command and exit |
| `-b` | Print elapsed time after each command |
| `-CODE <currency>` | Set currency code (default: `KRW`) |

```bash
./fehu -d mybudget -CODE USD
./fehu -c "get acc" -d mybudget
```

## MCP Server

Fehu can run as an [MCP (Model Context Protocol)](https://modelcontextprotocol.io) server, allowing AI assistants like Claude to manage your ledger through natural language.

### Start the MCP Server

```bash
./fehu -mcp
# or open a database on startup
./fehu -mcp -d mybudget
```

### Connecting to Claude Desktop

Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "fehu": {
      "command": "/path/to/fehu",
      "args": ["-mcp"]
    }
  }
}
```

### Available MCP Tools

| Tool | Description |
|------|-------------|
| `open_db` | Open or create a database |
| `get_accounts` | List accounts with balances (filterable by name/desc) |
| `create_account` | Create a new account |
| `delete_account` | Delete an account |
| `get_transactions` | List transactions (filterable by desc, id, time range) |
| `create_transaction` | Create a new transaction |
| `batch_create_transactions` | Create multiple transactions at once using a JSON array |
| `update_transaction` | Update transaction description or time |
| `delete_transaction` | Delete a transaction |
| `get_tags` | List all tags |
| `delete_tag` | Delete a tag |
| `calc` | Evaluate a calc expression |
| `get_readme` | Retrieve this documentation |

### Notes

- **String literals** in `calc` expressions must use **single quotes**: `sum(acc(__all__, 'expense*'))`
- `open_db` with a relative path resolves from the directory of the `fehu` executable. Use absolute paths to avoid ambiguity.
- The MCP server communicates over **stdio** — it is designed to be launched as a subprocess by an MCP client, not run interactively.

---

## Command Reference

### Accounts

| Command | Description |
|---------|-------------|
| `new acc [-d=<desc>] <name>` | Create an account |
| `get acc` | List all accounts with balances |
| `get acc name <name>` | Find account by name |
| `get acc desc <desc>` | Find accounts by description |
| `get acc child <name>` | List child accounts (prefix match) |
| `alt acc [-d=<desc>] <name>` | Update account description |
| `alt acc rename <old> <new>` | Rename an account |
| `del acc <name>` | Delete an account |

Account names use colon-separated hierarchy (e.g., `expense:food:dining`). Prefix a name with `~` in a transaction record to invert the direction.

### Transactions

| Command | Description |
|---------|-------------|
| `new txn [-t=<time>] [-d=<desc>] <record>` | Create a transaction |
| `get txn [-save=<name>]` | List all transactions |
| `get txn id <id> [-save=<name>]` | Find transaction by ID |
| `get txn time <from>~<to> [-save=<name>]` | Find transactions within a time range |
| `get txn desc <desc> [-save=<name>]` | Find transactions by description |
| `alt txn [-t=<time>] [-d=<desc>] <id>` | Update transaction metadata |
| `alt txn record <id> <record>` | Replace all records of a transaction |
| `del txn <id>` | Delete a transaction |

**Time format:** `YYYY-MM-DD;HH:MM:SS`

**Record format:** `account<amount` (inflow) or `account>amount` (outflow), entries separated by `;`. All entries must sum to zero.

**Open-ended time range:** Leave either side of `~` empty to query without a bound.
```
get txn time 2024-01-01;00:00:00~
get txn time ~2024-12-31;23:59:59
```

### Tags

| Command | Description |
|---------|-------------|
| `new tag [-d=<desc>] <name>` | Create a tag manually |
| `get tag` | List all tags |
| `get tag name <name>` | Find tag by name |
| `get tag desc <desc>` | Find tags by description |
| `alt tag [-d=<desc>] <name>` | Update tag description |
| `alt tag rename <old> <new>` | Rename a tag |
| `del tag <name>` | Delete a tag |

### Calculator

```
calc <expression>
def <name> <expression>
```

`def` stores the result of an expression as a named variable for reuse. `calc` evaluates and prints an expression immediately.

### Other

| Command | Description |
|---------|-------------|
| `open <name>` | Open (or create) a `.db` file |
| `close` | Close the current database |
| `quit` | Exit Fehu |

---

## Best Practice Guide

### 1. Default Account Structure

#### Assets — what you own
```
asset:bank            # Bank account
asset:cash            # Cash on hand
asset:savings         # Savings / time deposits
asset:investment      # Stocks, funds, crypto
asset:realestate      # Real estate
asset:car             # Vehicles
```

#### Liabilities — what you owe
```
liability:creditcard  # Credit card balance
liability:loan        # Personal loans
liability:mortgage    # Mortgage
```

#### Equity — net worth tracking
```
equity:opening        # Opening balance (use when starting a new ledger)
```

#### Income — money coming in
```
income:salary         # Salary
income:freelance      # Freelance / side income
income:dividend       # Dividends
income:interest       # Interest income
income:rental         # Rental income
income:capital_gain   # Realized capital gains
```

#### Expenses — money going out
```
expense:food          # Food & dining
expense:transport     # Transportation
expense:housing       # Rent / maintenance
expense:utility       # Utilities
expense:health        # Healthcare
expense:education     # Education
expense:culture       # Entertainment & leisure
expense:clothing      # Clothing
expense:subscription  # Subscriptions
```

---

### 2. Recommended `def` Variables

Add these to your workflow to track net worth and profitability at a glance.

```
# Net worth (Assets - Liabilities)
def net_worth sum(acc(__all__, 'asset*')) - sum(acc(__all__, 'liability*'))

# Net income (Income - Expenses)
def net_income sum(acc(__all__, 'income*')) - sum(acc(__all__, 'expense*'))

# Total assets
def total_assets sum(acc(__all__, 'asset*'))

# Total liabilities
def total_liabilities sum(acc(__all__, 'liability*'))

# Total expenses
def total_expenses sum(acc(__all__, 'expense*'))

# Total income
def total_income sum(acc(__all__, 'income*'))
```

Usage:
```
calc net_worth
calc net_income
```

---

### 3. Tag Conventions

#### Future / Pending transactions
| Tag | Meaning |
|-----|---------|
| `#NOTYET` | Transaction recorded but not yet settled (e.g. future installments) |
| `#PENDING` | Awaiting confirmation (e.g. pending transfers) |

#### Asset classification
| Tag | Meaning |
|-----|---------|
| `#cash` | Cash-equivalent assets (bank, cash on hand) |
| `#illiquid` | Hard-to-liquidate assets (real estate, car) |
| `#investment` | Investment assets subject to market value changes |

#### Unrealized gains & losses
| Tag | Meaning |
|-----|---------|
| `#unrealized` | Value change not yet realized (e.g. stock price increase) |
| `#realized` | Gain or loss actually settled |

---

### 4. Tracking Patterns

#### Cash-equivalent assets
```
# See total liquid assets
calc sum(atag(__all__, 'cash'))
```

#### Unrealized gains/losses
```
# Record a stock price increase (not yet sold)
new acc -d="Unrealized gain on stocks #unrealized #investment" income:capital_gain
new txn -d="AAPL price increase #unrealized" asset:investment<500000;income:capital_gain>500000

# When sold, mark as realized
alt txn -d="AAPL sold #realized" <txn_id>
```

#### Exclude future transactions from current net worth
```
# Net worth excluding #NOTYET transactions
def settled_txns xor(__all__, ttag(__all__, 'NOTYET'))
def real_net_worth sum(acc(settled_txns, 'asset*')) - sum(acc(settled_txns, 'liability*'))
calc real_net_worth
```

---

### 5. Example: Starting a New Ledger

```
# 1. Open your database
open mybudget

# 2. Create accounts
new acc asset:bank
new acc asset:cash
new acc liability:creditcard
new acc equity:opening
new acc income:salary
new acc expense:food
new acc expense:transport

# 3. Record opening balances
new txn -d="Opening balance" asset:bank<3000000;asset:cash<500000;equity:opening>3500000

# 4. Set up def variables
def net_worth sum(acc(__all__, 'asset*')) - sum(acc(__all__, 'liability*'))
def net_income sum(acc(__all__, 'income*')) - sum(acc(__all__, 'expense*'))

# 5. Check your position
calc net_worth
```

---

### 6. Tips

- **Always balance**: every transaction must sum to zero across all entries.
- **Use `#NOTYET`** for future installments or scheduled payments so they don't distort your current net worth.
- **Use `#unrealized`** for paper gains/losses on investments — convert to `#realized` when actually settled.
- **Use `atag` vs `ttag`**: `atag` filters by account description tags, `ttag` filters by transaction description tags.
- **String literals** in `calc` expressions use **single quotes** (`'expense*'`), not double quotes.

---

## Architecture

Fehu is split into three modules:

### `cmd/fehu` — Application layer
Entry point and command handlers. Wires together the `cli` and `core` packages.

```
cmd/fehu/
├── main.go       # Entry point, REPL loop, flags
├── proc.go       # Command dispatcher (switch on State)
├── states.go     # DFA state graph & parser initialization
├── accfunc.go    # Account command handlers
├── txnfunc.go    # Transaction command handlers
├── tagfunc.go    # Tag command handlers
└── calc.go       # Calculator command handlers
```

### `cli` — DFA-based command parser
A reusable Go library that drives the interactive REPL. Commands are parsed using a **Deterministic Finite Automaton (DFA)** — each keyword transitions the parser to a new state, enabling unambiguous hierarchical sub-commands like `get txn time`.

```
cli/
├── cli.go        # REPL loop (Run / Exec)
├── parse.go      # DFA walker, Cmd & FlagVar types
├── state.go      # State & Flag node definitions
└── strings.go    # Unicode-aware tokenizer with quote support
```

| Type | Role |
|------|------|
| `State` | A DFA node: regex pattern, child states, and valid flags |
| `Flag` | A named option (`-f=value` or `--flag=value`) with regex-validated value |
| `Cmd` | Parse result: resolved `State`, matched `FlagVar`s, positional params |
| `Parser` | Walks the DFA tree token by token and returns a `Cmd` |
| `CLI` | Reads from an `io.Reader`, prints the prompt, dispatches each line |

### `core` — Business logic & persistence
All SQLite queries, domain structs, the expression calculator, and the `Table` type used by the calc engine.

```
core/
├── db.go          # DB open/close, schema creation
├── struct.go      # Acc, Txn, Record, Tag structs
├── account.go     # Account CRUD & balance calculation
├── transaction.go # Transaction CRUD & pretty-printing
├── record.go      # Record CRUD
├── tag.go         # Tag CRUD, auto-tag parsing from #hashtags
├── table.go       # Table type: set ops, acc/tag/period filtering
├── stmt.go        # AST evaluator & built-in function dispatch
├── func.go        # SureName, SureID, wildcard Search, ParseTime
└── variable.go    # Global state: currency Code, DBPath, Vars map
```

The `Table` type is a set of transaction IDs. All `get txn` results can be saved as a `Table` and composed with set operations (`union`, `intersect`, `xor`) before aggregation.

### `ast` — Expression parser & AST
A standalone package that tokenizes and parses calc expressions into an Abstract Syntax Tree, which `core` then evaluates.

```
ast/
├── tokenize.go   # Lexer: converts expression strings into Token slices
├── tokens.go     # Token constants, operator precedence & arity tables
├── ast.go        # AST builder: sorts operators by precedence, wires child nodes
└── types.go      # Value types: Num, Str, Bool, List, Variable, Sym, Computed, Void
```

Expressions are parsed in three steps: **tokenize** → **sort by precedence** → **build tree bottom-up**. The `Computed` type allows lazy values (e.g. `__all__`) that re-evaluate on every access.

---

## Database Schema

| Table | Columns | Description |
|-------|---------|-------------|
| `acc` | `id`, `name`, `desc` | Accounts |
| `txn` | `id`, `desc`, `time` | Transactions |
| `record` | `id`, `tid`, `aid`, `amount` | Individual entries within a transaction |
| `Tag` | `id`, `name`, `desc` | Tags |
| `tagacc` | `tagid`, `aid` | Account–tag associations |
| `tagtxn` | `tagid`, `tid` | Transaction–tag associations |

Foreign keys are enforced with `ON DELETE CASCADE`, so deleting an account or transaction automatically cleans up its records and tag links.

---

## Dependencies

| Package | Purpose |
|---------|---------|
| [`mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3) | SQLite driver (CGO) |
| [`Rhymond/go-money`](https://github.com/Rhymond/go-money) | Currency-aware money arithmetic |
| [`pkg/errors`](https://github.com/pkg/errors) | Structured error wrapping |

---

## License

MIT
