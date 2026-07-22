package agg_test

import (
	"bytes"
	"testing"

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

func TestBucketAutoStage(t *testing.T) {
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

// --- $documents ---

func TestDocumentsStage_TestAPipelineStage(t *testing.T) {
	got := agg.Pipeline{
		agg.DocumentsStage(
			bson.D{{Key: "x", Value: 10}},
			bson.D{{Key: "x", Value: 2}},
			bson.D{{Key: "x", Value: 5}},
		),
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
