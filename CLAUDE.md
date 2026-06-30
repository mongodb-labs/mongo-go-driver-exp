# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

Experimental Go library (`github.com/mongodb-labs/mongo-go-driver-exp/agg`) that provides a type-safe builder API for MongoDB aggregation pipelines. All APIs are unstable and subject to breaking changes. The code lives in `agg/` inside a Go workspace at the repo root.

## Build and test

From the repo root (as CI does):

```
go build -v work
go test -v work
```

Tests are pure unit tests — no MongoDB instance is needed. They marshal the pipeline to BSON and compare bytes.

## Core pattern: sealed builder interface

Every stage and field type follows this pattern:

1. **Public interface** with a single unexported method that returns the private struct:
   ```go
   type SetField interface{ setField() setField }
   ```
2. **Private struct** that holds the data:
   ```go
   type setField struct { name string; expr Expr }
   ```
3. **Public constructors** that return the interface:
   ```go
   func Assign(field string, expr Expr) SetField { ... }
   ```
4. **Stage function** that unwraps via the interface method and builds `bson.D`:
   ```go
   func SetStage(fields ...SetField) Stage { ... }
   ```

Never expose the private struct directly. Always unwrap through the interface method inside the stage function.

## Type constraints (expr.go)

Map MQL spec argument types to these Go constraints:

| Spec type            | Go constraint / type |
|----------------------|----------------------|
| `resolvesToNumber`   | `NumberResolver`        |
| `resolvesToArray`    | `ArrayResolver`         |
| `resolvesToString`   | `StringResolver`        |
| `resolvesToBool`     | `BoolResolver`          |
| `expression`         | `Expr`               |
| `accumulator`        | `Accumulator`        |
| plain Go number      | `Number`             |

The `Expr` alias is `any`. Use the narrower typed constraints (`NumberResolver`, etc.) when the spec restricts the argument.

## BSON encoding

Spec `encode` field maps to:
- `single` — one field, scalar value: `bson.D{{Key: "$op", Value: expr}}`
- `object` — multiple named sub-fields: `bson.D{{Key: "$op", Value: bson.D{{Key: "input", Value: ...}, ...}}}`
- `array` — variadic args marshaled as a slice

Accumulators live in `accumulator.go` and hold a `bson.D` marshaled via `MarshalBSON`. Stages live in `stage.go`.

## MQL specifications

Full operator specs (YAML) are at:
```
~/Projects/mql-specifications/definitions/
  stage/          ← pipeline stages ($group, $match, $sort, …)
  expression/     ← expression operators ($add, $abs, …)
  accumulator/    ← accumulator operators ($sum, $avg, $push, …)
  query/          ← query filter operators
```

Each YAML file has `name`, `type`, `encode`, `description`, and `arguments` (with `type`, `variadic`, `variadicMin`) plus `tests`.
