package agg_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/mongodb-labs/mongo-go-driver-exp/mql/agg"
	"github.com/mongodb-labs/mongo-go-driver-exp/mql/query"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// assertPipelineEqual marshals pipeline and a bson.A of expected stages and
// fails the test if they differ.
func assertPipelineEqual(t *testing.T, pipeline agg.Pipeline, wantStages bson.A) {
	t.Helper()
	got, err := bson.Marshal(pipeline)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	want, err := bson.Marshal(bson.D{{Key: "pipeline", Value: wantStages}})
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if !bytes.Equal(want, got) {
		t.Errorf("Pipelines don't match.\nWant: %s\nGot:  %s",
			bson.Raw(want).String(),
			bson.Raw(got).String())
	}
}

// --- $addFields ---

func TestAddFieldsStage_UsingTwoStages(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("totalHomework", agg.Sum("$homework")),
			agg.Assign("totalQuiz", agg.Sum("$quiz")),
		),
		agg.AddFieldsStage(
			agg.Assign("totalScore", agg.Add("$totalHomework", "$totalQuiz", "$extraCredit")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "totalHomework", Value: bson.D{
				{Key: "$sum", Value: bson.A{"$homework"}},
			}},
			{Key: "totalQuiz", Value: bson.D{
				{Key: "$sum", Value: bson.A{"$quiz"}},
			}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "totalScore", Value: bson.D{
				{Key: "$add", Value: bson.A{"$totalHomework", "$totalQuiz", "$extraCredit"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAddFieldsStage_AddingFieldsToEmbeddedDoc(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("specs.fuel_type", "unleaded"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "specs.fuel_type", Value: "unleaded"}}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAddFieldsStage_OverwritingExistingField(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("cats", 20),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "cats", Value: 20}}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAddFieldsStage_AddElementToArray(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("_id", query.Eq(1)),
		),
		agg.AddFieldsStage(
			agg.Assign("homework", agg.ConcatArrays("$homework", []int{7})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "$eq", Value: 1},
			}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "homework", Value: bson.D{
				{Key: "$concatArrays", Value: bson.A{"$homework", bson.A{7}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $bucket ---

func TestBucketStage_BucketByYearAndFilterByBucketResults(t *testing.T) {
	got := agg.Pipeline{
		agg.BucketStage(
			"$year_born",
			[]any{1840, 1850, 1860, 1870, 1880},
			agg.WithBucketDefault("Other"),
			agg.WithBucketOutput(
				agg.Accumulate("count", agg.SumAccumulator(1)),
				agg.Accumulate("artists", agg.PushAccumulator(bson.D{
					{Key: "name", Value: agg.Concat("$first_name", " ", "$last_name")},
					{Key: "year_born", Value: "$year_born"},
				})),
			),
		),
		agg.MatchStage(
			query.Field("count", query.Gt(3)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$bucket", Value: bson.D{
			{Key: "groupBy", Value: "$year_born"},
			{Key: "boundaries", Value: bson.A{1840, 1850, 1860, 1870, 1880}},
			{Key: "default", Value: "Other"},
			{Key: "output", Value: bson.D{
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "artists", Value: bson.D{{Key: "$push", Value: bson.D{
					{Key: "name", Value: bson.D{{Key: "$concat", Value: bson.A{"$first_name", " ", "$last_name"}}}},
					{Key: "year_born", Value: "$year_born"},
				}}}},
			}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "count", Value: bson.D{{Key: "$gt", Value: 3}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBucketStage_UseBucketWithFacetToBucketByMultipleFields(t *testing.T) {
	got := agg.Pipeline{
		agg.FacetStage(
			agg.Facet("price",
				agg.BucketStage(
					"$price",
					[]any{0, 200, 400},
					agg.WithBucketDefault("Other"),
					agg.WithBucketOutput(
						agg.Accumulate("count", agg.SumAccumulator(1)),
						agg.Accumulate("artwork", agg.PushAccumulator(bson.D{
							{Key: "title", Value: "$title"},
							{Key: "price", Value: "$price"},
						})),
						agg.Accumulate("averagePrice", agg.AvgAccumulator("$price")),
					),
				),
			),
			agg.Facet("year",
				agg.BucketStage(
					"$year",
					[]any{1890, 1910, 1920, 1940},
					agg.WithBucketDefault("Unknown"),
					agg.WithBucketOutput(
						agg.Accumulate("count", agg.SumAccumulator(1)),
						agg.Accumulate("artwork", agg.PushAccumulator(bson.D{
							{Key: "title", Value: "$title"},
							{Key: "year", Value: "$year"},
						})),
					),
				),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$facet", Value: bson.D{
			{Key: "price", Value: bson.A{
				bson.D{{Key: "$bucket", Value: bson.D{
					{Key: "groupBy", Value: "$price"},
					{Key: "boundaries", Value: bson.A{0, 200, 400}},
					{Key: "default", Value: "Other"},
					{Key: "output", Value: bson.D{
						{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
						{Key: "artwork", Value: bson.D{{Key: "$push", Value: bson.D{
							{Key: "title", Value: "$title"},
							{Key: "price", Value: "$price"},
						}}}},
						{Key: "averagePrice", Value: bson.D{{Key: "$avg", Value: "$price"}}},
					}},
				}}},
			}},
			{Key: "year", Value: bson.A{
				bson.D{{Key: "$bucket", Value: bson.D{
					{Key: "groupBy", Value: "$year"},
					{Key: "boundaries", Value: bson.A{1890, 1910, 1920, 1940}},
					{Key: "default", Value: "Unknown"},
					{Key: "output", Value: bson.D{
						{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
						{Key: "artwork", Value: bson.D{{Key: "$push", Value: bson.D{
							{Key: "title", Value: "$title"},
							{Key: "year", Value: "$year"},
						}}}},
					}},
				}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $bucketAuto ---

func TestBucketAutoStage_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.BucketAutoStage(
			"$price",
			4,
		),
	}
	want := bson.A{
		bson.D{{Key: "$bucketAuto", Value: bson.D{
			{Key: "groupBy", Value: "$price"},
			{Key: "buckets", Value: 4},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBucketAutoStage_WithGranularity(t *testing.T) {
	got := agg.Pipeline{
		agg.BucketAutoStage(
			"$price",
			4,
			agg.WithBucketAutoGranularity("R5"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$bucketAuto", Value: bson.D{
			{Key: "groupBy", Value: "$price"},
			{Key: "buckets", Value: 4},
			{Key: "granularity", Value: "R5"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBucketAutoStage_WithOutputAndGranularity(t *testing.T) {
	got := agg.Pipeline{
		agg.BucketAutoStage(
			"$price",
			4,
			agg.WithBucketAutoGranularity("POWERSOF2"),
			agg.WithBucketAutoOutput(
				agg.Accumulate("count", agg.SumAccumulator(1)),
				agg.Accumulate("averagePrice", agg.AvgAccumulator("$price")),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$bucketAuto", Value: bson.D{
			{Key: "groupBy", Value: "$price"},
			{Key: "buckets", Value: 4},
			{Key: "output", Value: bson.D{
				{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				{Key: "averagePrice", Value: bson.D{{Key: "$avg", Value: "$price"}}},
			}},
			{Key: "granularity", Value: "POWERSOF2"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $changeStream ---

func TestChangeStream_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.ChangeStreamStage(),
	}
	want := bson.A{
		bson.D{{Key: "$changeStream", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestChangeStream_WithFullDocument(t *testing.T) {
	got := agg.Pipeline{
		agg.ChangeStreamStage(
			agg.WithChangeStreamFullDocument("updateLookup"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$changeStream", Value: bson.D{
			{Key: "fullDocument", Value: "updateLookup"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestChangeStream_WithFullDocumentBeforeChange(t *testing.T) {
	got := agg.Pipeline{
		agg.ChangeStreamStage(
			agg.WithChangeStreamFullDocumentBeforeChange("whenAvailable"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$changeStream", Value: bson.D{
			{Key: "fullDocumentBeforeChange", Value: "whenAvailable"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestChangeStream_WithResumeAfter(t *testing.T) {
	token := bson.D{{Key: "_data", Value: "8264A0F1C0000000012B"}}
	got := agg.Pipeline{
		agg.ChangeStreamStage(
			agg.WithChangeStreamResumeAfter(token),
		),
	}
	want := bson.A{
		bson.D{{Key: "$changeStream", Value: bson.D{
			{Key: "resumeAfter", Value: bson.D{{Key: "_data", Value: "8264A0F1C0000000012B"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestChangeStream_WithStartAtOperationTime(t *testing.T) {
	got := agg.Pipeline{
		agg.ChangeStreamStage(
			agg.WithChangeStreamStartAtOperationTime(bson.Timestamp{T: 1688674200, I: 1}),
		),
	}
	want := bson.A{
		bson.D{{Key: "$changeStream", Value: bson.D{
			{Key: "startAtOperationTime", Value: bson.Timestamp{T: 1688674200, I: 1}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestChangeStream_WithMultipleOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ChangeStreamStage(
			agg.WithChangeStreamShowExpandedEvents(true),
			agg.WithChangeStreamFullDocument("updateLookup"),
			agg.WithChangeStreamAllChangesForCluster(true),
		),
	}
	want := bson.A{
		bson.D{{Key: "$changeStream", Value: bson.D{
			{Key: "allChangesForCluster", Value: true},
			{Key: "fullDocument", Value: "updateLookup"},
			{Key: "showExpandedEvents", Value: true},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $changeStreamSplitLargeEvent ---

func TestChangeStreamSplitLargeEvent_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.ChangeStreamSplitLargeEventStage(),
	}
	want := bson.A{
		bson.D{{Key: "$changeStreamSplitLargeEvent", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestChangeStreamSplitLargeEvent_AsFinalChangeStreamStage(t *testing.T) {
	got := agg.Pipeline{
		agg.ChangeStreamStage(),
		agg.ChangeStreamSplitLargeEventStage(),
	}
	want := bson.A{
		bson.D{{Key: "$changeStream", Value: bson.D{}}},
		bson.D{{Key: "$changeStreamSplitLargeEvent", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestChangeStreamSplitLargeEvent_AfterChangeStreamWithOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ChangeStreamStage(
			agg.WithChangeStreamFullDocument("required"),
			agg.WithChangeStreamFullDocumentBeforeChange("required"),
		),
		agg.ChangeStreamSplitLargeEventStage(),
	}
	want := bson.A{
		bson.D{{Key: "$changeStream", Value: bson.D{
			{Key: "fullDocument", Value: "required"},
			{Key: "fullDocumentBeforeChange", Value: "required"},
		}}},
		bson.D{{Key: "$changeStreamSplitLargeEvent", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $collStats ---

func TestCollStatsStage_LatencyStatsDocument(t *testing.T) {
	got := agg.Pipeline{
		agg.CollStatsStage(
			agg.WithCollStatsLatencyStats(true),
		),
	}
	want := bson.A{
		bson.D{{Key: "$collStats", Value: bson.D{
			{Key: "latencyStats", Value: bson.D{{Key: "histograms", Value: true}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCollStatsStage_StorageStatsDocument(t *testing.T) {
	got := agg.Pipeline{
		agg.CollStatsStage(
			agg.WithCollStatsStorageStats(),
		),
	}
	want := bson.A{
		bson.D{{Key: "$collStats", Value: bson.D{
			{Key: "storageStats", Value: bson.D{}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCollStatsStage_Count(t *testing.T) {
	got := agg.Pipeline{
		agg.CollStatsStage(
			agg.WithCollStatsCount(),
		),
	}
	want := bson.A{
		bson.D{{Key: "$collStats", Value: bson.D{
			{Key: "count", Value: bson.D{}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCollStatsStage_QueryExecStats(t *testing.T) {
	got := agg.Pipeline{
		agg.CollStatsStage(
			agg.WithCollStatsQueryExecStats(),
		),
	}
	want := bson.A{
		bson.D{{Key: "$collStats", Value: bson.D{
			{Key: "queryExecStats", Value: bson.D{}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCollStatsStage_MultipleOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.CollStatsStage(
			agg.WithCollStatsQueryExecStats(),
			agg.WithCollStatsCount(),
			agg.WithCollStatsLatencyStats(true),
			agg.WithCollStatsStorageStats(),
		),
	}
	want := bson.A{
		bson.D{{Key: "$collStats", Value: bson.D{
			{Key: "latencyStats", Value: bson.D{{Key: "histograms", Value: true}}},
			{Key: "storageStats", Value: bson.D{}},
			{Key: "count", Value: bson.D{}},
			{Key: "queryExecStats", Value: bson.D{}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $count ---

func TestCountStage(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("score", query.Gt(80))),
		agg.CountStage("passing_scores"),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "score", Value: bson.D{{Key: "$gt", Value: 80}}}}}},
		bson.D{{Key: "$count", Value: "passing_scores"}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $currentOp ---

func TestCurrentOpStage_InactiveSessions(t *testing.T) {
	got := agg.Pipeline{
		agg.CurrentOpStage(
			agg.WithCurrentOpAllUsers(true),
			agg.WithCurrentOpIdleSessions(true),
		),
		agg.MatchStage(
			query.Field("active", query.Eq(false)),
			query.Field("transaction", query.Exists(true)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$currentOp", Value: bson.D{
			{Key: "allUsers", Value: true},
			{Key: "idleSessions", Value: true},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "active", Value: bson.D{{Key: "$eq", Value: false}}},
			{Key: "transaction", Value: bson.D{{Key: "$exists", Value: true}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCurrentOpStage_IdleConnections(t *testing.T) {
	got := agg.Pipeline{
		agg.CurrentOpStage(
			agg.WithCurrentOpIdleConnections(true),
		),
	}
	want := bson.A{
		bson.D{{Key: "$currentOp", Value: bson.D{
			{Key: "idleConnections", Value: true},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCurrentOpStage_LocalOps(t *testing.T) {
	got := agg.Pipeline{
		agg.CurrentOpStage(
			agg.WithCurrentOpLocalOps(true),
		),
	}
	want := bson.A{
		bson.D{{Key: "$currentOp", Value: bson.D{
			{Key: "localOps", Value: true},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCurrentOpStage_MultipleOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.CurrentOpStage(
			agg.WithCurrentOpLocalOps(true),
			agg.WithCurrentOpIdleCursors(true),
			agg.WithCurrentOpAllUsers(true),
			agg.WithCurrentOpIdleSessions(false),
			agg.WithCurrentOpIdleConnections(true),
		),
	}
	want := bson.A{
		bson.D{{Key: "$currentOp", Value: bson.D{
			{Key: "allUsers", Value: true},
			{Key: "idleConnections", Value: true},
			{Key: "idleCursors", Value: true},
			{Key: "idleSessions", Value: false},
			{Key: "localOps", Value: true},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $densify ---

func TestDensifyStage_DensifyTimeSeriesData(t *testing.T) {
	lower := time.Date(2021, 5, 18, 0, 0, 0, 0, time.UTC)
	upper := time.Date(2021, 5, 18, 8, 0, 0, 0, time.UTC)
	got := agg.Pipeline{
		agg.DensifyStage(
			"timestamp",
			1,
			agg.DensifyBoundsValues(lower, upper),
			agg.WithDensifyUnit("hour"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$densify", Value: bson.D{
			{Key: "field", Value: "timestamp"},
			{Key: "range", Value: bson.D{
				{Key: "bounds", Value: bson.A{lower, upper}},
				{Key: "step", Value: 1},
				{Key: "unit", Value: "hour"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDensifyStage_DensificationWithPartitions(t *testing.T) {
	got := agg.Pipeline{
		agg.DensifyStage(
			"altitude",
			200,
			agg.DensifyBoundsFull(),
			agg.WithDensifyPartitionByFields("variety"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$densify", Value: bson.D{
			{Key: "field", Value: "altitude"},
			{Key: "partitionByFields", Value: bson.A{"variety"}},
			{Key: "range", Value: bson.D{
				{Key: "bounds", Value: "full"},
				{Key: "step", Value: 200},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDensifyStage_NumericBounds(t *testing.T) {
	got := agg.Pipeline{
		agg.DensifyStage(
			"value",
			25,
			agg.DensifyBoundsValues(0, 100),
		),
	}
	want := bson.A{
		bson.D{{Key: "$densify", Value: bson.D{
			{Key: "field", Value: "value"},
			{Key: "range", Value: bson.D{
				{Key: "bounds", Value: bson.A{0, 100}},
				{Key: "step", Value: 25},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDensifyStage_PartitionBounds(t *testing.T) {
	got := agg.Pipeline{
		agg.DensifyStage(
			"timestamp",
			1,
			agg.DensifyBoundsPartition(),
			agg.WithDensifyUnit("hour"),
			agg.WithDensifyPartitionByFields("sensorId"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$densify", Value: bson.D{
			{Key: "field", Value: "timestamp"},
			{Key: "partitionByFields", Value: bson.A{"sensorId"}},
			{Key: "range", Value: bson.D{
				{Key: "bounds", Value: "partition"},
				{Key: "step", Value: 1},
				{Key: "unit", Value: "hour"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $documents ---

func TestDocumentsStage_TestAPipelineStage(t *testing.T) {
	got := agg.Pipeline{
		agg.DocumentsStage(bson.A{
			bson.D{{Key: "x", Value: 10}},
			bson.D{{Key: "x", Value: 2}},
			bson.D{{Key: "x", Value: 5}},
		}),
		agg.BucketAutoStage("$x", 4),
	}
	want := bson.A{
		bson.D{{Key: "$documents", Value: bson.A{
			bson.D{{Key: "x", Value: 10}},
			bson.D{{Key: "x", Value: 2}},
			bson.D{{Key: "x", Value: 5}},
		}}},
		bson.D{{Key: "$bucketAuto", Value: bson.D{
			{Key: "groupBy", Value: "$x"},
			{Key: "buckets", Value: 4},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $facet ---

func TestFacetStage(t *testing.T) {
	got := agg.Pipeline{
		agg.FacetStage(
			agg.Facet("categorizedByTags",
				agg.UnwindStage("$tags"),
				agg.SortByCountStage("$tags"),
			),
			agg.Facet("categorizedByPrice",
				agg.MatchStage(query.Field("price", query.Exists(true))),
				agg.BucketStage(
					"$price",
					[]any{0, 150, 200, 300, 400},
					agg.WithBucketDefault("Other"),
					agg.WithBucketOutput(
						agg.Accumulate("count", agg.SumAccumulator(1)),
						agg.Accumulate("titles", agg.PushAccumulator("$title")),
					),
				),
			),
			agg.Facet("categorizedByYears(Auto)",
				agg.BucketAutoStage("$year", 4),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$facet", Value: bson.D{
			{Key: "categorizedByTags", Value: bson.A{
				bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$tags"}}}},
				bson.D{{Key: "$sortByCount", Value: "$tags"}},
			}},
			{Key: "categorizedByPrice", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "price", Value: bson.D{{Key: "$exists", Value: true}}}}}},
				bson.D{{Key: "$bucket", Value: bson.D{
					{Key: "groupBy", Value: "$price"},
					{Key: "boundaries", Value: bson.A{0, 150, 200, 300, 400}},
					{Key: "default", Value: "Other"},
					{Key: "output", Value: bson.D{
						{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
						{Key: "titles", Value: bson.D{{Key: "$push", Value: "$title"}}},
					}},
				}}},
			}},
			{Key: "categorizedByYears(Auto)", Value: bson.A{
				bson.D{{Key: "$bucketAuto", Value: bson.D{
					{Key: "groupBy", Value: "$year"},
					{Key: "buckets", Value: 4},
				}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $fill ---

func TestFillStage_ConstantValue(t *testing.T) {
	got := agg.Pipeline{
		agg.FillStage([]agg.FillOutput{
			agg.FillWithValue("bootsSold", 0),
			agg.FillWithValue("sandalsSold", 0),
			agg.FillWithValue("sneakersSold", 0),
		}),
	}
	want := bson.A{
		bson.D{{Key: "$fill", Value: bson.D{
			{Key: "output", Value: bson.D{
				{Key: "bootsSold", Value: bson.D{{Key: "value", Value: 0}}},
				{Key: "sandalsSold", Value: bson.D{{Key: "value", Value: 0}}},
				{Key: "sneakersSold", Value: bson.D{{Key: "value", Value: 0}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFillStage_LinearInterpolation(t *testing.T) {
	got := agg.Pipeline{
		agg.FillStage(
			[]agg.FillOutput{agg.FillWithMethod("price", "linear")},
			agg.WithFillSortBy(agg.Sort("time", agg.Asc)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$fill", Value: bson.D{
			{Key: "sortBy", Value: bson.D{{Key: "time", Value: 1}}},
			{Key: "output", Value: bson.D{
				{Key: "price", Value: bson.D{{Key: "method", Value: "linear"}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFillStage_LastObservedValue(t *testing.T) {
	got := agg.Pipeline{
		agg.FillStage(
			[]agg.FillOutput{agg.FillWithMethod("score", "locf")},
			agg.WithFillSortBy(agg.Sort("date", agg.Asc)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$fill", Value: bson.D{
			{Key: "sortBy", Value: bson.D{{Key: "date", Value: 1}}},
			{Key: "output", Value: bson.D{
				{Key: "score", Value: bson.D{{Key: "method", Value: "locf"}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFillStage_DistinctPartitions(t *testing.T) {
	got := agg.Pipeline{
		agg.FillStage(
			[]agg.FillOutput{agg.FillWithMethod("score", "locf")},
			agg.WithFillSortBy(agg.Sort("date", agg.Asc)),
			agg.WithFillPartitionBy(bson.D{{Key: "restaurant", Value: "$restaurant"}}),
		),
	}
	want := bson.A{
		bson.D{{Key: "$fill", Value: bson.D{
			{Key: "sortBy", Value: bson.D{{Key: "date", Value: 1}}},
			{Key: "partitionBy", Value: bson.D{{Key: "restaurant", Value: "$restaurant"}}},
			{Key: "output", Value: bson.D{
				{Key: "score", Value: bson.D{{Key: "method", Value: "locf"}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $geoNear ---

func TestGeoNearStage_MaximumDistance(t *testing.T) {
	got := agg.Pipeline{
		agg.GeoNearStage(
			bson.D{
				{Key: "type", Value: "Point"},
				{Key: "coordinates", Value: bson.A{-73.99279, 40.719296}},
			},
			agg.WithGeoNearDistanceField("dist.calculated"),
			agg.WithGeoNearMaxDistance(2),
			agg.WithGeoNearQuery(query.Field("category", query.Eq("Parks"))),
			agg.WithGeoNearIncludeLocs("dist.location"),
			agg.WithGeoNearSpherical(true),
		),
	}
	want := bson.A{
		bson.D{{Key: "$geoNear", Value: bson.D{
			{Key: "near", Value: bson.D{
				{Key: "type", Value: "Point"},
				{Key: "coordinates", Value: bson.A{-73.99279, 40.719296}},
			}},
			{Key: "distanceField", Value: "dist.calculated"},
			{Key: "maxDistance", Value: 2},
			{Key: "query", Value: bson.D{{Key: "category", Value: bson.D{{Key: "$eq", Value: "Parks"}}}}},
			{Key: "includeLocs", Value: "dist.location"},
			{Key: "spherical", Value: true},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGeoNearStage_MinimumDistance(t *testing.T) {
	got := agg.Pipeline{
		agg.GeoNearStage(
			bson.D{
				{Key: "type", Value: "Point"},
				{Key: "coordinates", Value: bson.A{-73.99279, 40.719296}},
			},
			agg.WithGeoNearDistanceField("dist.calculated"),
			agg.WithGeoNearMinDistance(2),
			agg.WithGeoNearQuery(query.Field("category", query.Eq("Parks"))),
			agg.WithGeoNearIncludeLocs("dist.location"),
			agg.WithGeoNearSpherical(true),
		),
	}
	want := bson.A{
		bson.D{{Key: "$geoNear", Value: bson.D{
			{Key: "near", Value: bson.D{
				{Key: "type", Value: "Point"},
				{Key: "coordinates", Value: bson.A{-73.99279, 40.719296}},
			}},
			{Key: "distanceField", Value: "dist.calculated"},
			{Key: "minDistance", Value: 2},
			{Key: "query", Value: bson.D{{Key: "category", Value: bson.D{{Key: "$eq", Value: "Parks"}}}}},
			{Key: "includeLocs", Value: "dist.location"},
			{Key: "spherical", Value: true},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $graphLookup ---

func TestGraphLookupStage_WithinASingleCollection(t *testing.T) {
	got := agg.Pipeline{
		agg.GraphLookupStage("employees", "$reportsTo", "reportsTo", "name", "reportingHierarchy"),
	}
	want := bson.A{
		bson.D{{Key: "$graphLookup", Value: bson.D{
			{Key: "from", Value: "employees"},
			{Key: "startWith", Value: "$reportsTo"},
			{Key: "connectFromField", Value: "reportsTo"},
			{Key: "connectToField", Value: "name"},
			{Key: "as", Value: "reportingHierarchy"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGraphLookupStage_AcrossMultipleCollections(t *testing.T) {
	got := agg.Pipeline{
		agg.GraphLookupStage("airports", "$nearestAirport", "connects", "airport", "destinations",
			agg.WithGraphLookupMaxDepth(2),
			agg.WithGraphLookupDepthField("numConnections"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$graphLookup", Value: bson.D{
			{Key: "from", Value: "airports"},
			{Key: "startWith", Value: "$nearestAirport"},
			{Key: "connectFromField", Value: "connects"},
			{Key: "connectToField", Value: "airport"},
			{Key: "as", Value: "destinations"},
			{Key: "maxDepth", Value: 2},
			{Key: "depthField", Value: "numConnections"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGraphLookupStage_WithAQueryFilter(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("name", query.Eq("Tanya Jordan"))),
		agg.GraphLookupStage("people", "$friends", "friends", "name", "golfers",
			agg.WithGraphLookupRestrictSearchWithMatch(query.Field("hobbies", query.Eq("golf"))),
		),
		agg.ProjectStage(
			agg.Include("name"),
			agg.Include("friends"),
			agg.Compute("connections who play golf", "$golfers.name"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "name", Value: bson.D{{Key: "$eq", Value: "Tanya Jordan"}}}}}},
		bson.D{{Key: "$graphLookup", Value: bson.D{
			{Key: "from", Value: "people"},
			{Key: "startWith", Value: "$friends"},
			{Key: "connectFromField", Value: "friends"},
			{Key: "connectToField", Value: "name"},
			{Key: "as", Value: "golfers"},
			{Key: "restrictSearchWithMatch", Value: bson.D{{Key: "hobbies", Value: bson.D{{Key: "$eq", Value: "golf"}}}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "friends", Value: 1},
			{Key: "connections who play golf", Value: "$golfers.name"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $group ---

func TestGroupStage_CountNumDocsInColl(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage(
			agg.Null,
			agg.Accumulate("count", agg.CountAccumulator()),
		),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "count", Value: bson.D{{Key: "$count", Value: bson.D{}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGroupStage_RetrieveDistinctVals(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage(
			"$rated",
			agg.Accumulate("totalRuntime", agg.SumAccumulator("$runtime")),
		),
		agg.MatchStage(
			query.Field("totalRuntime", query.Gte(100000)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$rated"},
			{Key: "totalRuntime", Value: bson.D{{Key: "$sum", Value: "$runtime"}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "totalRuntime", Value: bson.D{{Key: "$gte", Value: 100000}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGroupStage_CalculateCountSumAvg(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("year", query.Lt(1910)),
		),
		agg.GroupStage(
			"$year",
			agg.Accumulate("totalRuntime", agg.SumAccumulator("$runtime")),
			agg.Accumulate("averageRuntime", agg.AvgAccumulator("$runtime")),
			agg.Accumulate("count", agg.SumAccumulator(1)),
		),
		agg.SortStage(
			agg.Sort("totalRuntime", agg.Desc),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$lt", Value: 1910}}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$year"},
			{Key: "totalRuntime", Value: bson.D{{Key: "$sum", Value: "$runtime"}}},
			{Key: "averageRuntime", Value: bson.D{{Key: "$avg", Value: "$runtime"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "totalRuntime", Value: -1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGroupStage_GroupByNull(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("year", query.Lt(1910)),
		),
		agg.GroupStage(
			"$year",
			agg.Accumulate("totalRuntime", agg.SumAccumulator("$runtime")),
			agg.Accumulate("averageRuntime", agg.AvgAccumulator("$runtime")),
			agg.Accumulate("count", agg.SumAccumulator(1)),
		),
		agg.SortStage(
			agg.Sort("totalRuntime", agg.Desc),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$lt", Value: 1910}}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$year"},
			{Key: "totalRuntime", Value: bson.D{{Key: "$sum", Value: "$runtime"}}},
			{Key: "averageRuntime", Value: bson.D{{Key: "$avg", Value: "$runtime"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "totalRuntime", Value: -1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGroupStage_GroupTitlesByYear(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("year", query.Lt(1910)),
		),
		agg.GroupStage(
			"$year",
			agg.Accumulate("titles", agg.PushAccumulator("$title")),
		),
		agg.SortStage(
			agg.Sort("_id", agg.Asc),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$lt", Value: 1910}}},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$year"},
			{Key: "titles", Value: bson.D{{Key: "$push", Value: "$title"}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "_id", Value: 1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestGroupStage_RetrieveDistinctValues
// spec test is $group with _id only and no accumulator fields; ready to implement

// TODO: implement TestGroupStage_GroupByItemHaving
// ready to implement; uses SumAccumulator(Multiply(...)) followed by MatchStage

// TODO: implement TestGroupStage_CalculateCountSumAvgWithDates
// after $dateToString expression operator is implemented

// TODO: implement TestGroupStage_GroupByNullSumMultiply
// ready to implement; existing TestGroupStage_GroupByNull incorrectly duplicates TestGroupStage_CalculateCountSumAvg

func TestGroupStage_GroupDocumentsByAuthor(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage(
			"$author",
			agg.Accumulate("books", agg.PushAccumulator(agg.RootObject())),
		),
		agg.AddFieldsStage(
			agg.Assign("totalCopies", agg.Sum("$books.copies")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$author"},
			{Key: "books", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "totalCopies", Value: bson.D{{Key: "$sum", Value: bson.A{"$books.copies"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $indexStats ---

func TestIndexStatsStage(t *testing.T) {
	got := agg.Pipeline{
		agg.IndexStatsStage(),
	}
	want := bson.A{
		bson.D{{Key: "$indexStats", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $limit ---

func TestLimitStage(t *testing.T) {
	got := agg.Pipeline{
		agg.LimitStage(5),
	}
	want := bson.A{
		bson.D{{Key: "$limit", Value: 5}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $listLocalSessions ---

func TestListLocalSessionsStage_AllUsers(t *testing.T) {
	got := agg.Pipeline{
		agg.ListLocalSessionsStage(agg.WithListLocalSessionsAllUsers(true)),
	}
	want := bson.A{
		bson.D{{Key: "$listLocalSessions", Value: bson.D{{Key: "allUsers", Value: true}}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestListLocalSessionsStage_SpecifiedUsers(t *testing.T) {
	got := agg.Pipeline{
		agg.ListLocalSessionsStage(agg.WithListLocalSessionsUsers(
			agg.SessionUser{User: "myAppReader", DB: "test"},
		)),
	}
	want := bson.A{
		bson.D{{Key: "$listLocalSessions", Value: bson.D{
			{Key: "users", Value: bson.A{
				bson.D{{Key: "user", Value: "myAppReader"}, {Key: "db", Value: "test"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestListLocalSessionsStage_CurrentUser(t *testing.T) {
	got := agg.Pipeline{
		agg.ListLocalSessionsStage(),
	}
	want := bson.A{
		bson.D{{Key: "$listLocalSessions", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $listSampledQueries ---

func TestListSampledQueriesStage_AllCollections(t *testing.T) {
	got := agg.Pipeline{
		agg.ListSampledQueriesStage(),
	}
	want := bson.A{
		bson.D{{Key: "$listSampledQueries", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestListSampledQueriesStage_SpecificCollection(t *testing.T) {
	got := agg.Pipeline{
		agg.ListSampledQueriesStage(agg.WithListSampledQueriesNamespace("social.post")),
	}
	want := bson.A{
		bson.D{{Key: "$listSampledQueries", Value: bson.D{{Key: "namespace", Value: "social.post"}}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $listSearchIndexes ---

func TestListSearchIndexesStage_AllIndexes(t *testing.T) {
	got := agg.Pipeline{
		agg.ListSearchIndexesStage(),
	}
	want := bson.A{
		bson.D{{Key: "$listSearchIndexes", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestListSearchIndexesStage_ByName(t *testing.T) {
	got := agg.Pipeline{
		agg.ListSearchIndexesStage(agg.WithListSearchIndexesName("synonym-mappings")),
	}
	want := bson.A{
		bson.D{{Key: "$listSearchIndexes", Value: bson.D{{Key: "name", Value: "synonym-mappings"}}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestListSearchIndexesStage_ById(t *testing.T) {
	got := agg.Pipeline{
		agg.ListSearchIndexesStage(agg.WithListSearchIndexesID("6524096020da840844a4c4a7")),
	}
	want := bson.A{
		bson.D{{Key: "$listSearchIndexes", Value: bson.D{{Key: "id", Value: "6524096020da840844a4c4a7"}}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $listSessions ---

func TestListSessionsStage_AllUsers(t *testing.T) {
	got := agg.Pipeline{
		agg.ListSessionsStage(agg.WithListSessionsAllUsers(true)),
	}
	want := bson.A{
		bson.D{{Key: "$listSessions", Value: bson.D{{Key: "allUsers", Value: true}}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestListSessionsStage_SpecifiedUsers(t *testing.T) {
	got := agg.Pipeline{
		agg.ListSessionsStage(agg.WithListSessionsUsers(
			agg.SessionUser{User: "myAppReader", DB: "test"},
		)),
	}
	want := bson.A{
		bson.D{{Key: "$listSessions", Value: bson.D{
			{Key: "users", Value: bson.A{
				bson.D{{Key: "user", Value: "myAppReader"}, {Key: "db", Value: "test"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestListSessionsStage_CurrentUser(t *testing.T) {
	got := agg.Pipeline{
		agg.ListSessionsStage(),
	}
	want := bson.A{
		bson.D{{Key: "$listSessions", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $lookup ---

func TestLookupStage_SingleEqualityJoin(t *testing.T) {
	got := agg.Pipeline{
		agg.LookupStage("inventory_docs",
			agg.WithLookupFrom("inventory"),
			agg.WithLookupLocalField("item"),
			agg.WithLookupForeignField("sku"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "inventory"},
			{Key: "localField", Value: "item"},
			{Key: "foreignField", Value: "sku"},
			{Key: "as", Value: "inventory_docs"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLookupStage_WithMergeObjects(t *testing.T) {
	got := agg.Pipeline{
		agg.LookupStage("fromItems",
			agg.WithLookupFrom("items"),
			agg.WithLookupLocalField("item"),
			agg.WithLookupForeignField("item"),
		),
		agg.ReplaceRootStage(agg.MergeObjects(
			agg.ArrayElemAt("$fromItems", 0),
			agg.RootObject(),
		)),
		agg.ProjectStage(agg.Exclude("fromItems")),
	}
	want := bson.A{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "items"},
			{Key: "localField", Value: "item"},
			{Key: "foreignField", Value: "item"},
			{Key: "as", Value: "fromItems"},
		}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{
			{Key: "newRoot", Value: bson.D{{Key: "$mergeObjects", Value: bson.A{
				bson.D{{Key: "$arrayElemAt", Value: bson.A{"$fromItems", 0}}},
				"$$ROOT",
			}}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{{Key: "fromItems", Value: 0}}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLookupStage_UncorrelatedSubquery(t *testing.T) {
	got := agg.Pipeline{
		agg.LookupStage("holidays",
			agg.WithLookupFrom("holidays"),
			agg.WithLookupPipeline(
				agg.MatchStage(query.Field("year", query.Eq(2018))),
				agg.ProjectStage(
					agg.Compute("date", bson.D{
						{Key: "name", Value: "$name"},
						{Key: "date", Value: "$date"},
					}),
				),
				agg.ReplaceRootStage("$date"),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "holidays"},
			{Key: "pipeline", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "year", Value: bson.D{{Key: "$eq", Value: 2018}}}}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "name", Value: "$name"},
						{Key: "date", Value: "$date"},
					}},
				}}},
				bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$date"}}}},
			}},
			{Key: "as", Value: "holidays"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestLookupStage_CorrelatedSubquery (uses WithLookupLet plus a
// $match with the $expr query operator, which is not implemented yet)

// --- $match ---

func TestMatchStage_EqualityMatch(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("rated", query.Eq("TV-PG")),
			query.Field("runtime", query.Gt(1000)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "rated", Value: bson.D{{Key: "$eq", Value: "TV-PG"}}},
			{Key: "runtime", Value: bson.D{{Key: "$gt", Value: 1000}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestMatchStage_PerformCount
// after multi-condition field queries are supported in the query package

func TestMatchStage_EmptyFieldCondition(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("x"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "x", Value: bson.D{}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $merge ---

func TestMergeStage_OnDemandMaterializedView(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage(
			bson.D{{Key: "fiscal_year", Value: "$fiscal_year"}, {Key: "dept", Value: "$dept"}},
			agg.Accumulate("salaries", agg.SumAccumulator("$salary")),
		),
		agg.MergeStage("budgets",
			agg.WithMergeIntoDB("reporting"),
			agg.WithMergeOn("_id"),
			agg.WithMergeWhenMatched("replace"),
			agg.WithMergeWhenNotMatched("insert"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "fiscal_year", Value: "$fiscal_year"}, {Key: "dept", Value: "$dept"}}},
			{Key: "salaries", Value: bson.D{{Key: "$sum", Value: "$salary"}}},
		}}},
		bson.D{{Key: "$merge", Value: bson.D{
			{Key: "into", Value: bson.D{{Key: "db", Value: "reporting"}, {Key: "coll", Value: "budgets"}}},
			{Key: "on", Value: "_id"},
			{Key: "whenMatched", Value: "replace"},
			{Key: "whenNotMatched", Value: "insert"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMergeStage_OnlyInsertNewData(t *testing.T) {
	got := agg.Pipeline{
		agg.MergeStage("orgArchive",
			agg.WithMergeIntoDB("reporting"),
			agg.WithMergeOn("dept", "fiscal_year"),
			agg.WithMergeWhenMatched("fail"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$merge", Value: bson.D{
			{Key: "into", Value: bson.D{{Key: "db", Value: "reporting"}, {Key: "coll", Value: "orgArchive"}}},
			{Key: "on", Value: bson.A{"dept", "fiscal_year"}},
			{Key: "whenMatched", Value: "fail"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMergeStage_MergeResultsFromMultipleCollections(t *testing.T) {
	got := agg.Pipeline{
		agg.MergeStage("quarterlyreport",
			agg.WithMergeOn("_id"),
			agg.WithMergeWhenMatched("merge"),
			agg.WithMergeWhenNotMatched("insert"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$merge", Value: bson.D{
			{Key: "into", Value: "quarterlyreport"},
			{Key: "on", Value: "_id"},
			{Key: "whenMatched", Value: "merge"},
			{Key: "whenNotMatched", Value: "insert"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $out ---

func TestOutStage_SameDatabase(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage("$author", agg.Accumulate("books", agg.PushAccumulator("$title"))),
		agg.OutStage("authors"),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$author"},
			{Key: "books", Value: bson.D{{Key: "$push", Value: "$title"}}},
		}}},
		bson.D{{Key: "$out", Value: bson.D{{Key: "coll", Value: "authors"}}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestOutStage_DifferentDatabase(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage("$author", agg.Accumulate("books", agg.PushAccumulator("$title"))),
		agg.OutStage("authors", agg.WithOutDB("reporting")),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$author"},
			{Key: "books", Value: bson.D{{Key: "$push", Value: "$title"}}},
		}}},
		bson.D{{Key: "$out", Value: bson.D{
			{Key: "db", Value: "reporting"},
			{Key: "coll", Value: "authors"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestOutStage_TimeSeriesCollection(t *testing.T) {
	got := agg.Pipeline{
		agg.OutStage("sensorData",
			agg.WithOutDB("reporting"),
			agg.WithOutTimeseries(agg.NewTimeseries("timestamp",
				agg.WithTimeseriesMetaField("sensorId"),
				agg.WithTimeseriesGranularity("hours"),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$out", Value: bson.D{
			{Key: "db", Value: "reporting"},
			{Key: "coll", Value: "sensorData"},
			{Key: "timeseries", Value: bson.D{
				{Key: "timeField", Value: "timestamp"},
				{Key: "metaField", Value: "sensorId"},
				{Key: "granularity", Value: "hours"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $planCacheStats ---

func TestPlanCacheStatsStage_ReturnInformationForAllEntries(t *testing.T) {
	got := agg.Pipeline{
		agg.PlanCacheStatsStage(),
	}
	want := bson.A{
		bson.D{{Key: "$planCacheStats", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestPlanCacheStatsStage_FindCacheEntryDetailsForQueryHash(t *testing.T) {
	got := agg.Pipeline{
		agg.PlanCacheStatsStage(),
		agg.MatchStage(query.Field("planCacheKey", query.Eq("B1435201"))),
	}
	want := bson.A{
		bson.D{{Key: "$planCacheStats", Value: bson.D{}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "planCacheKey", Value: bson.D{{Key: "$eq", Value: "B1435201"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $project ---

func TestProjectStage_IncludeSpecificFieldsInOutputDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Include("title"),
			agg.Include("rated"),
		),
		agg.LimitStage(1),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: bson.D{{Key: "$eq", Value: "The Great Train Robbery"}}}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "rated", Value: 1},
		}}},
		bson.D{{Key: "$limit", Value: 1}},
	}
	assertPipelineEqual(t, got, want)
}

func TestProjectStage_SuppressIdFieldInOutputDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("title"),
			agg.Include("rated"),
		),
		agg.LimitStage(1),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: bson.D{{Key: "$eq", Value: "The Great Train Robbery"}}}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "title", Value: 1},
			{Key: "rated", Value: 1},
		}}},
		bson.D{{Key: "$limit", Value: 1}},
	}
	assertPipelineEqual(t, got, want)
}

func TestProjectStage_ExcludeFieldsFromOutputDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(agg.Exclude("rated")),
		agg.LimitStage(1),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: bson.D{{Key: "$eq", Value: "The Great Train Robbery"}}}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "rated", Value: 0},
		}}},
		bson.D{{Key: "$limit", Value: 1}},
	}
	assertPipelineEqual(t, got, want)
}

func TestProjectStage_ExcludeFieldsFromEmbeddedDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Exclude("imdb.id"),
			agg.Exclude("type"),
		),
		agg.LimitStage(1),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: bson.D{{Key: "$eq", Value: "The Great Train Robbery"}}}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "imdb.id", Value: 0},
			{Key: "type", Value: 0},
		}}},
		bson.D{{Key: "$limit", Value: 1}},
	}
	assertPipelineEqual(t, got, want)
}

func TestProjectStage_IncludeSpecificFieldsFromEmbeddedDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(agg.Include("imdb.rating")),
		agg.LimitStage(1),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: bson.D{{Key: "$eq", Value: "The Great Train Robbery"}}}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "imdb.rating", Value: 1},
		}}},
		bson.D{{Key: "$limit", Value: 1}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: Add computed field leadActor: { $arrayElemAt: [ "$cast", 0 ] }
// when arrayElemAt operator is implemented
func TestProjectStage_IncludeComputedFields(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("The Great Train Robbery"))),
		agg.ProjectStage(
			agg.Include("title"),
			agg.Compute("releaseYear", "$year"),
		),
		agg.LimitStage(1),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: bson.D{{Key: "$eq", Value: "The Great Train Robbery"}}}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "releaseYear", Value: "$year"},
		}}},
		bson.D{{Key: "$limit", Value: 1}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestProjectStage_ConditionallyExcludeFields
// after $cond operator is implemented

// TODO: implement TestProjectStage_IncludeComputedFieldsSubstr
// after $substr operator is implemented

// TODO: implement TestProjectStage_ProjectNewArrayFields
// after array-literal ProjectionField variant is implemented

// --- $redact ---

// TODO: implement TestRedactStage_EvaluateAccessAtEveryDocumentLevel
// after the $cond expression operator is implemented (the spec example uses
// $cond with $$DESCEND / $$PRUNE system variables)

// TODO: implement TestRedactStage_ExcludeAllFieldsAtAGivenLevel
// after the $cond expression operator is implemented

// --- $replaceRoot ---

func TestReplaceRootStage_WithAnEmbeddedDocumentField(t *testing.T) {
	got := agg.Pipeline{
		agg.ReplaceRootStage(agg.MergeObjects(
			bson.D{
				{Key: "dogs", Value: 0},
				{Key: "cats", Value: 0},
				{Key: "birds", Value: 0},
				{Key: "fish", Value: 0},
			},
			"$pets",
		)),
	}
	want := bson.A{
		bson.D{{Key: "$replaceRoot", Value: bson.D{
			{Key: "newRoot", Value: bson.D{{Key: "$mergeObjects", Value: bson.A{
				bson.D{
					{Key: "dogs", Value: 0},
					{Key: "cats", Value: 0},
					{Key: "birds", Value: 0},
					{Key: "fish", Value: 0},
				},
				"$pets",
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReplaceRootStage_WithADocumentNestedInAnArray(t *testing.T) {
	got := agg.Pipeline{
		agg.UnwindStage("$grades"),
		agg.MatchStage(query.Field("grades.grade", query.Gte(90))),
		agg.ReplaceRootStage("$grades"),
	}
	want := bson.A{
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$grades"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "grades.grade", Value: bson.D{{Key: "$gte", Value: 90}}}}}},
		bson.D{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$grades"}}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $replaceWith ---

func TestReplaceWithStage_AnEmbeddedDocumentField(t *testing.T) {
	got := agg.Pipeline{
		agg.ReplaceWithStage(agg.MergeObjects(
			bson.D{
				{Key: "dogs", Value: 0},
				{Key: "cats", Value: 0},
				{Key: "birds", Value: 0},
				{Key: "fish", Value: 0},
			},
			"$pets",
		)),
	}
	want := bson.A{
		bson.D{{Key: "$replaceWith", Value: bson.D{{Key: "$mergeObjects", Value: bson.A{
			bson.D{
				{Key: "dogs", Value: 0},
				{Key: "cats", Value: 0},
				{Key: "birds", Value: 0},
				{Key: "fish", Value: 0},
			},
			"$pets",
		}}}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReplaceWithStage_ADocumentNestedInAnArray(t *testing.T) {
	got := agg.Pipeline{
		agg.UnwindStage("$grades"),
		agg.MatchStage(query.Field("grades.grade", query.Gte(90))),
		agg.ReplaceWithStage("$grades"),
	}
	want := bson.A{
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$grades"}}}},
		bson.D{{Key: "$match", Value: bson.D{{Key: "grades.grade", Value: bson.D{{Key: "$gte", Value: 90}}}}}},
		bson.D{{Key: "$replaceWith", Value: "$grades"}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $sample ---

func TestSampleStagee(t *testing.T) {
	got := agg.Pipeline{
		agg.SampleStage(3),
	}
	want := bson.A{
		bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: 3}}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $set ---

func TestSetStage_AddingFieldsToEmbeddedDoc(t *testing.T) {
	got := agg.Pipeline{
		agg.SetStage(agg.Assign("specs.fuel_type", "unleaded")),
	}
	want := bson.A{
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "specs.fuel_type", Value: "unleaded"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSetStage_OverwriteExistingField(t *testing.T) {
	got := agg.Pipeline{
		agg.SetStage(agg.Assign("cats", 20)),
	}
	want := bson.A{
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "cats", Value: 20},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestSetStage_UsingTwoSetStages
// after $sum expression operator is implemented

// TODO: implement TestSetStage_AddElementToArray
// after $concatArrays expression operator is implemented

// TODO: implement TestSetStage_CreatingNewFieldWithExistingFields
// after $avg expression operator is implemented

// --- $setWindowFields ---

func TestSetWindowFieldsStage_CumulativeQuantityForState(t *testing.T) {
	got := agg.Pipeline{
		agg.SetWindowFieldsStage(
			[]agg.WindowField{
				agg.WindowOutput("cumulativeQuantityForState", agg.SumAccumulator("$quantity"),
					agg.WithWindowDocuments(agg.WindowUnbounded, agg.WindowCurrent)),
			},
			agg.WithSetWindowFieldsPartitionBy("$state"),
			agg.WithSetWindowFieldsSortBy(agg.Sort("orderDate", agg.Asc)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$setWindowFields", Value: bson.D{
			{Key: "partitionBy", Value: "$state"},
			{Key: "sortBy", Value: bson.D{{Key: "orderDate", Value: 1}}},
			{Key: "output", Value: bson.D{
				{Key: "cumulativeQuantityForState", Value: bson.D{
					{Key: "$sum", Value: "$quantity"},
					{Key: "window", Value: bson.D{{Key: "documents", Value: bson.A{"unbounded", "current"}}}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSetWindowFieldsStage_CumulativeQuantityForYear(t *testing.T) {
	got := agg.Pipeline{
		agg.SetWindowFieldsStage(
			[]agg.WindowField{
				agg.WindowOutput("cumulativeQuantityForYear", agg.SumAccumulator("$quantity"),
					agg.WithWindowDocuments(agg.WindowUnbounded, agg.WindowCurrent)),
			},
			agg.WithSetWindowFieldsPartitionBy(agg.Year("$orderDate")),
			agg.WithSetWindowFieldsSortBy(agg.Sort("orderDate", agg.Asc)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$setWindowFields", Value: bson.D{
			{Key: "partitionBy", Value: bson.D{{Key: "$year", Value: bson.D{{Key: "date", Value: "$orderDate"}}}}},
			{Key: "sortBy", Value: bson.D{{Key: "orderDate", Value: 1}}},
			{Key: "output", Value: bson.D{
				{Key: "cumulativeQuantityForYear", Value: bson.D{
					{Key: "$sum", Value: "$quantity"},
					{Key: "window", Value: bson.D{{Key: "documents", Value: bson.A{"unbounded", "current"}}}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSetWindowFieldsStage_MovingAverageQuantityForYear(t *testing.T) {
	got := agg.Pipeline{
		agg.SetWindowFieldsStage(
			[]agg.WindowField{
				agg.WindowOutput("averageQuantity", agg.AvgAccumulator("$quantity"),
					agg.WithWindowDocuments(agg.WindowOffset(-1), agg.WindowOffset(0))),
			},
			agg.WithSetWindowFieldsPartitionBy(agg.Year("$orderDate")),
			agg.WithSetWindowFieldsSortBy(agg.Sort("orderDate", agg.Asc)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$setWindowFields", Value: bson.D{
			{Key: "partitionBy", Value: bson.D{{Key: "$year", Value: bson.D{{Key: "date", Value: "$orderDate"}}}}},
			{Key: "sortBy", Value: bson.D{{Key: "orderDate", Value: 1}}},
			{Key: "output", Value: bson.D{
				{Key: "averageQuantity", Value: bson.D{
					{Key: "$avg", Value: "$quantity"},
					{Key: "window", Value: bson.D{{Key: "documents", Value: bson.A{-1, 0}}}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $shardedDataDistribution ---

func TestShardedDataDistributionStage(t *testing.T) {
	got := agg.Pipeline{
		agg.ShardedDataDistributionStage(),
	}
	want := bson.A{
		bson.D{{Key: "$shardedDataDistribution", Value: bson.D{}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $skip ---

func TestSkipStage(t *testing.T) {
	got := agg.Pipeline{
		agg.SkipStage(5),
	}
	want := bson.A{
		bson.D{{Key: "$skip", Value: 5}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $sort ---

func TestSortStage_AscendingDescendingSort(t *testing.T) {
	got := agg.Pipeline{
		agg.SortStage(
			agg.Sort("age", agg.Desc),
			agg.Sort("posts", agg.Asc),
		),
	}
	want := bson.A{
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "age", Value: -1},
			{Key: "posts", Value: 1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestSortStage_TextScoreMetadataSort
// after $meta sort modifier and $text query operator are implemented

// --- $sortByCount ---

func TestSortByCountStage(t *testing.T) {
	got := agg.Pipeline{
		agg.UnwindStage("$tags"),
		agg.SortByCountStage("$tags"),
	}
	want := bson.A{
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$tags"}}}},
		bson.D{{Key: "$sortByCount", Value: "$tags"}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $unset ---

func TestUnsetStage_RemoveASingleField(t *testing.T) {
	got := agg.Pipeline{
		agg.UnsetStage("copies"),
	}
	want := bson.A{
		bson.D{{Key: "$unset", Value: bson.A{"copies"}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestUnsetStage_RemoveTopLevelFields(t *testing.T) {
	got := agg.Pipeline{
		agg.UnsetStage("isbn", "copies"),
	}
	want := bson.A{
		bson.D{{Key: "$unset", Value: bson.A{"isbn", "copies"}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestUnsetStage_RemoveEmbeddedFields(t *testing.T) {
	got := agg.Pipeline{
		agg.UnsetStage("isbn", "author.first", "copies.warehouse"),
	}
	want := bson.A{
		bson.D{{Key: "$unset", Value: bson.A{"isbn", "author.first", "copies.warehouse"}}},
	}
	assertPipelineEqual(t, got, want)
}

// --- $unwind ---

func TestUnwindStage_UnwindArray(t *testing.T) {
	got := agg.Pipeline{
		agg.UnwindStage("$sizes"),
	}
	want := bson.A{
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$sizes"}}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestUnwindStage_PreserveNullAndEmptyArrays(t *testing.T) {
	got := agg.Pipeline{
		agg.UnwindStage("$sizes", agg.WithUnwindPreserveNullAndEmptyArrays(true)),
	}
	want := bson.A{
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$sizes"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestUnwindStage_IncludeArrayIndex(t *testing.T) {
	got := agg.Pipeline{
		agg.UnwindStage("$sizes", agg.WithUnwindIncludeArrayIndex("arrayIndex")),
	}
	want := bson.A{
		bson.D{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$sizes"},
			{Key: "includeArrayIndex", Value: "arrayIndex"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestUnwindStage_UnwindEmbeddedArrays(t *testing.T) {
	got := agg.Pipeline{
		agg.UnwindStage("$items"),
		agg.UnwindStage("$items.tags"),
		agg.GroupStage("$items.tags",
			agg.Accumulate("totalSalesAmount", agg.SumAccumulator(agg.Multiply("$items.price", "$items.quantity"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$items"}}}},
		bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$items.tags"}}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$items.tags"},
			{Key: "totalSalesAmount", Value: bson.D{{Key: "$sum", Value: bson.D{
				{Key: "$multiply", Value: bson.A{"$items.price", "$items.quantity"}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}
