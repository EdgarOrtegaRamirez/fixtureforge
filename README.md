# FixtureForge

Generate realistic test data for development and testing. A fast, single-binary CLI tool written in Go.

## Features

- **Auto-detection** — Automatically determines column types from names (`email`, `first_name`, `created_at`, etc.)
- **28 column types** — strings, ints, floats, bools, names, emails, phones, addresses, cities, states, countries, zip codes, dates, datetimes, times, UUIDs, URLs, IPv4/IPv6, colors, lorem ipsum, enums, regex patterns, ranges, JSON objects, and more
- **7 output formats** — CSV, JSON, JSONL, YAML, Markdown table, HTML table, SQL INSERT
- **Schema files** — Define complex schemas in YAML or JSON for repeatable generation
- **Deterministic output** — Seed-based generation for reproducible results
- **Configurable** — Custom parameters per column (min/max, format, values, patterns)
- **Fast** — Single Go binary, no dependencies, instant startup

## Installation

```bash
# From source
go install github.com/EdgarOrtegaRamirez/fixtureforge/cmd/fixtureforge@latest

# Or build locally
git clone https://github.com/EdgarOrtegaRamirez/fixtureforge.git
cd fixtureforge
go build -o fixtureforge ./cmd/fixtureforge/
```

## Quick Start

```bash
# Generate 10 rows of CSV with auto-detected columns
fixtureforge generate -n 10 --columns name,email,age,city,active

# Generate JSON with a specific seed for reproducibility
fixtureforge generate -n 5 -f json --columns id,name,email --seed 42

# See what types would be detected
fixtureforge detect --columns id,first_name,last_name,email,phone,created_at,price,active

# Generate from a schema file
fixtureforge schema sample > schema.yaml
# Edit schema.yaml...
fixtureforge generate -s schema.yaml -f json
```

## Column Types

| Type | Description | Parameters |
|------|-------------|------------|
| `string` | Random lowercase string | `min_length`, `max_length`, `prefix` |
| `int` | Random integer | `min`, `max` |
| `float` | Random float with decimals | `min`, `max`, `decimals` |
| `bool` | Random boolean | — |
| `bool_weight` | Boolean with custom probability | `true_pct` |
| `name` | Full name (first + last) | — |
| `first_name` | First name | — |
| `last_name` | Last name | — |
| `email` | Realistic email address | — |
| `phone` | Phone number | `format` (us, intl, e164) |
| `address` | Street address | — |
| `city` | City name | — |
| `state` | US state code | — |
| `country` | Country name | — |
| `zipcode` | Zip/postal code | `format` (us5, us9, uk) |
| `date` | Date | `format` (iso, us, eu) |
| `datetime` | ISO 8601 datetime | — |
| `time` | Time (HH:MM:SS) | — |
| `uuid` | UUID v4 | — |
| `url` | Random URL | — |
| `ipv4` | IPv4 address | — |
| `ipv6` | IPv6 address | — |
| `color` | Color value | `format` (hex, rgb, hsl) |
| `lorem` | Lorem ipsum text | `words` |
| `enum` | Random from list | `values` (comma-separated) |
| `regex` | Pattern-based string | `pattern` |
| `range` | Number in range | `min`, `max`, `step`, `decimals` |
| `json` | JSON object | — |

## Schema Files

Define columns with specific types and parameters:

```yaml
columns:
  - name: id
    type: int
    params:
      min: "1"
      max: "1000"
  - name: name
    type: name
  - name: email
    type: email
  - name: age
    type: range
    params:
      min: "18"
      max: "80"
  - name: status
    type: enum
    params:
      values: "active,inactive,pending"
  - name: created_at
    type: datetime
count: 100
seed: 42
```

## Output Formats

```bash
# CSV (default)
fixtureforge generate -n 5 --columns id,name,email

# JSON array
fixtureforge generate -n 5 -f json --columns id,name,email

# JSON Lines (newline-delimited JSON)
fixtureforge generate -n 5 -f jsonl --columns id,name,email

# YAML
fixtureforge generate -n 5 -f yaml --columns id,name,email

# Markdown table
fixtureforge generate -n 5 -f markdown --columns id,name,email

# HTML table
fixtureforge generate -n 5 -f html --columns id,name,email

# SQL INSERT
fixtureforge generate -n 5 -f sql --columns id,name,active
```

## Auto-Detection Examples

```bash
fixtureforge detect --columns \
  id,first_name,last_name,email,phone, \
  address,city,state,zip_code,country, \
  created_at,updated_at,price,quantity, \
  is_active,description,url,uuid,ip_address

# Output:
# id                   → int
# first_name           → first_name
# last_name            → last_name
# email                → email
# phone                → phone
# address              → address
# city                 → city
# state                → state
# zip_code             → zipcode
# country              → country
# created_at           → datetime
# updated_at           → datetime
# price                → float
# quantity             → int
# is_active            → bool
# description          → lorem
# url                  → url
# uuid                 → uuid
# ip_address           → ipv4
```

## Building

```bash
go build -o fixtureforge ./cmd/fixtureforge/
```

## Testing

```bash
go test ./... -v
```

## Architecture

```
fixtureforge/
├── cmd/fixtureforge/     # CLI entry point (cobra)
├── pkg/fixtureforge/     # Public API (engine)
├── internal/
│   ├── types/            # Core data types
│   ├── generators/       # 28 column type generators
│   ├── detectors/        # Auto-detection from column names
│   └── formatters/       # Output formatters (CSV, JSON, etc.)
└── tests/                # Integration tests
```

## License

MIT
