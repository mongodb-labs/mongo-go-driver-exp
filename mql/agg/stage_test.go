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

// --- $project ---

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_IncludeSpecificFieldsInOutputDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", "The Great Train Robbery")),
		agg.ProjectStage(
			agg.Include("title"),
			agg.Include("rated"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: "The Great Train Robbery"}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "rated", Value: 1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_SuppressIdFieldInOutputDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", "The Great Train Robbery")),
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("title"),
			agg.Include("rated"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: "The Great Train Robbery"}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "title", Value: 1},
			{Key: "rated", Value: 1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_ExcludeFieldsFromOutputDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", "The Great Train Robbery")),
		agg.ProjectStage(agg.Exclude("rated")),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: "The Great Train Robbery"}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "rated", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_ExcludeFieldsFromEmbeddedDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", "The Great Train Robbery")),
		agg.ProjectStage(
			agg.Exclude("imdb.id"),
			agg.Exclude("type"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: "The Great Train Robbery"}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "imdb.id", Value: 0},
			{Key: "type", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: Add limit of 1 when limit stage is implemented
func TestProjectStage_IncludeSpecificFieldsFromEmbeddedDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", "The Great Train Robbery")),
		agg.ProjectStage(agg.Include("imdb.rating")),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: "The Great Train Robbery"}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "imdb.rating", Value: 1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: Add limit of 1 when limit stage is implemented
// TODO: Add computed field leadActor: { $arrayElemAt: [ "$cast", 0 ] }
// when arrayElemAt operator is implemented
func TestProjectStage_IncludeComputedFields(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", "The Great Train Robbery")),
		agg.ProjectStage(
			agg.Include("title"),
			agg.Compute("releaseYear", "$year"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{{Key: "title", Value: "The Great Train Robbery"}}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "releaseYear", Value: "$year"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestProjectStage_ConditionallyExcludeFields
// after $cond operator is implemented

// TODO: implement TestProjectStage_IncludeComputedFieldsSubstr
// after $substr operator is implemented

// TODO: implement TestProjectStage_ProjectNewArrayFields
// after array-literal ProjectionField variant is implemented

// --- $match ---

func TestMatchStage_EqualityMatch(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("rated", "TV-PG"),
			query.Field("runtime", query.Gt(1000)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "rated", Value: "TV-PG"},
			{Key: "runtime", Value: bson.D{{Key: "$gt", Value: 1000}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestMatchStage_PerformCount
// after multi-condition field queries are supported in the query package

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

// TODO: implement TestGroupStage_GroupDocumentsByAuthor
// after $addFields stage is implemented
