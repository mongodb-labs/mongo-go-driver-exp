---
name: bson-test-writing
description: Writes BSON pipeline tests for a stage, operator, or accumulator. Reads test examples from the MQL spec and generates assertPipelineEqual tests matching the style of agg/agg_test.go.
disable-model-invocation: false
---

You are writing tests for the `agg` package. Follow every step in order.

## Step 1: Locate the spec

The operator name is: **$ARGUMENTS** (strip the leading `$` if present for the filename lookup).

Search for the spec YAML across all subdirectories of `~/Projects/mql-specifications/definitions/`:

```
find ~/Projects/mql-specifications/definitions -name "<operatorname>.yaml"
```

Read all matching files. The `tests` section of each file is your source of truth for test cases.

## Step 2: Read the existing test file

Read `agg/agg_test.go` in full before writing anything. Understand:
- The `assertPipelineEqual(t, pipeline, wantStages)` helper signature
- The naming convention: `TestStageName_DescriptionInPascalCase`
- The section comment style: `// --- $stagename ---`
- How `bson.D`, `bson.A`, `bson.E` are constructed for the expected output

## Step 3: Extract test cases from the spec

From the `tests` section of the spec YAML, select tests that:
- Exercise the Go API for the operator being tested
- Are representable with currently-implemented stages and operators (skip tests that require stages not yet in `stage.go` or operators not yet in `operator.go`/`accumulator.go`)
- Cover distinct behaviors (don't duplicate tests that check the same encoding)

For each test, note:
- `name` — becomes the `_DescriptionInPascalCase` suffix of the function name (strip articles, convert to PascalCase)
- `pipeline` — the MongoDB pipeline JSON; this defines both the Go builder call and the expected BSON

## Step 4: Translate each test case to Go

For every selected test case, produce one `func TestXxx_Yyy(t *testing.T)` function:

```go
func TestGroupStage_CountNumDocsInColl(t *testing.T) {
    assertPipelineEqual(t,
        agg.Pipeline{
            agg.GroupStage(
                agg.Null,
                agg.Accumulate("count", agg.CountAccumulator()),
            ),
        },
        bson.A{
            bson.D{{Key: "$group", Value: bson.D{
                {Key: "_id", Value: nil},
                {Key: "count", Value: bson.D{{Key: "$count", Value: bson.D{}}}},
            }}},
        },
    )
}
```

Rules:
- The first argument to `assertPipelineEqual` is the **Go builder call** (use `agg.XxxStage(...)`, `agg.Accumulate(...)`, etc.)
- The second argument is the **expected BSON** (`bson.A` of `bson.D` stages) — transcribe directly from the spec's pipeline JSON
- Field references like `"$fieldname"` or `"$$ROOT"` are plain Go string literals
- `null` / `~` in YAML → `nil` in Go BSON
- Numeric literals: use untyped Go literals (`1`, `100`) — let the compiler infer
- If a spec test includes stages that aren't yet implemented, add a `// TODO: requires $xxx` comment and skip that stage in the pipeline (or skip the whole test if it's the primary stage under test)

## Step 5: Add a TODO comment for unsupported operators

If the spec test uses an operator that isn't yet implemented in the `agg` package, add a `// TODO: Add <description> when <stage/operator> is implemented` comment above the test function, matching the style of existing TODO comments in `agg_test.go`.

## Step 6: Insert into agg_test.go

- If a `// --- $operatorname ---` section already exists in `agg_test.go`, append the new test functions after the last existing test in that section.
- If no section exists yet, add the section comment and all test functions at the end of the file.
- Do not reorder or reformat existing tests.

## Step 7: Verify

Run `go test -v work` from the repo root. All new tests must pass. Fix any failures before declaring done.
