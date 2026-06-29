---
name: spec-to-go
description: Implements a new MongoDB aggregation stage, expression operator, or accumulator. Reads the MQL spec first, then generates Go code following the project's sealed builder + generics pattern.
disable-model-invocation: false
---

You are implementing a new operator for the `agg` package. Follow every step in order.

## Step 1: Locate the spec

The operator name is: **$ARGUMENTS** (strip the leading `$` if present for the filename lookup).

Search for the spec YAML in `~/Projects/mql-specifications/definitions/` across these subdirectories:
- `stage/` — pipeline stages
- `expression/` — expression operators
- `accumulator/` — accumulator operators
- `query/` — query filter operators

An operator may appear in more than one subdirectory (e.g. `$avg` is in both `accumulator/` and `expression/`). Read all matching files.

Use `find ~/Projects/mql-specifications/definitions -name "<operatorname>.yaml"` to locate them.

## Step 2: Parse the spec

From each YAML file, extract:
- `name` — the MongoDB operator name (e.g. `$limit`)
- `type` — what kind of operator it is (stage, accumulator, resolvesToNumber, etc.)
- `encode` — how arguments are encoded (`single`, `object`, `array`)
- `description` — what it does (use for the Go doc comment)
- `arguments` — list of args; for each: `name`, `type[]`, `variadic`, `variadicMin`, `description`

## Step 3: Map spec types to Go

Use the type constraint table from CLAUDE.md:

| Spec type            | Go constraint / type |
|----------------------|----------------------|
| `resolvesToNumber`   | `NumberResolver`        |
| `resolvesToArray`    | `ArrayResolver`         |
| `resolvesToString`   | `StringResolver`        |
| `resolvesToBool`     | `BoolResolver`          |
| `expression`         | `Expr`               |
| `accumulator`        | `Accumulator`        |
| plain numeric literal | `Number`            |

If an argument accepts multiple spec types (e.g. `resolvesToNumber | resolvesToDate`), use the widest matching Go constraint, or `Expr` if no single constraint covers all types.

## Step 4: Determine the Go encoding

| Spec `encode` | BSON shape | Go implementation |
|---------------|------------|-------------------|
| `single`      | `{$op: value}` | Single generic arg |
| `object`      | `{$op: {field1: v1, field2: v2}}` | Named params, build `bson.D` |
| `array`       | `{$op: [v1, v2, ...]}` | Variadic args marshaled as slice |

Variadic args (`variadic: array` in spec) become `...T` in Go. If `variadicMin: 0`, the variadic is optional.

## Step 5: Choose the right file and pattern

| Operator type   | File            | Pattern to follow |
|-----------------|-----------------|-------------------|
| Pipeline stage  | `agg/stage.go`  | Sealed interface (see $set, $sort, $group) |
| Accumulator     | `agg/accumulator.go` | `Accumulator{doc: bson.D{...}}` |
| Expression op   | `agg/operator.go` | Read existing operators for the pattern |

Read the relevant existing file before writing any code so you match the exact style.

## Step 6: Implement

For a **pipeline stage**:
1. Define the public `XxxField` interface and private `xxxField` struct (if the stage takes structured field args beyond simple expressions).
2. Write public constructors.
3. Write the `XxxStage(...)` function that builds `Stage{{Key: "$xxx", Value: ...}}`.
4. Add a `// --- $xxx ---` comment section to separate it from adjacent stages.

For an **accumulator**:
1. Write a constructor function `XxxAccumulator[T Constraint](args...) Accumulator`.
2. Return `Accumulator{doc: bson.D{{Key: "$xxx", Value: ...}}}`.
3. For `encode: object`, build the inner `bson.D` from named arguments.

For an **expression operator**:
1. Follow the pattern of existing operators in `operator.go`.
2. Return the appropriate typed expression struct (`NumberExpr`, `ArrayExpr`, etc.) based on what the spec's `type` field says the operator resolves to.

Always write a one-line doc comment using the spec's `description`. Keep it concise.

## Step 7: Verify

Run `go build -v work` from the repo root. Fix any compile errors before declaring done.
