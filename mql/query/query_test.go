package query_test

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

func TestEq_EqualsSpecificValue(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Eq(20)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$eq", Value: 20}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestEq_FieldInEmbeddedDocEqualsValue(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("item.name", query.Eq("ab")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "item.name", Value: bson.D{{Key: "$eq", Value: "ab"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestEq_EqualsArrayValue(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("tags", query.Eq([]string{"A", "B"})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "tags", Value: bson.D{{Key: "$eq", Value: bson.A{"A", "B"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestEq_RegexMatchBehavior(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("company", query.Eq("MongoDB"))),
		agg.MatchStage(query.Field("company", query.Eq(bson.Regex{Pattern: "^MongoDB"}))),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "company", Value: bson.D{{Key: "$eq", Value: "MongoDB"}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "company", Value: bson.D{{Key: "$eq", Value: bson.Regex{Pattern: "^MongoDB"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestExists_AndNotEqualTo(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Exists(true), query.Nin(5, 15)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{
				{Key: "$exists", Value: true},
				{Key: "$nin", Value: bson.A{5, 15}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestExists_NullValues(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Exists(true)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$exists", Value: true}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestExists_MissingField(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Exists(false)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$exists", Value: false}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGt(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Gt(20)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$gt", Value: 20}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGte(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Gte(20)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$gte", Value: 20}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIn_MatchValuesInArray(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("tags", query.In("home", "school")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "tags", Value: bson.D{{Key: "$in", Value: bson.A{"home", "school"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIn_RegularExpression(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("tags", query.In(bson.Regex{Pattern: "^be"}, bson.Regex{Pattern: "^st"})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "tags", Value: bson.D{{Key: "$in", Value: bson.A{
				bson.Regex{Pattern: "^be"},
				bson.Regex{Pattern: "^st"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLt(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Lt(20)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$lt", Value: 20}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLte(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Lte(20)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$lte", Value: 20}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNe(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Ne(20)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$ne", Value: 20}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNin_SelectOnUnmatchingDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("quantity", query.Nin(5, 15)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "quantity", Value: bson.D{{Key: "$nin", Value: bson.A{5, 15}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNin_SelectOnElementsNotInArray(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("tags", query.Nin("school")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "tags", Value: bson.D{{Key: "$nin", Value: bson.A{"school"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestType_QueryingByDataType(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("zipCode", query.Type(2))),
		agg.MatchStage(query.Field("zipCode", query.Type("string"))),
		agg.MatchStage(query.Field("zipCode", query.Type(1))),
		agg.MatchStage(query.Field("zipCode", query.Type("double"))),
		agg.MatchStage(query.Field("zipCode", query.Type("number"))),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{2}}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{"string"}}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{1}}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{"double"}}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{"number"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestType_QueryingByMultipleDataType(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("zipCode", query.Type(2, 1))),
		agg.MatchStage(query.Field("zipCode", query.Type("string", "double"))),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{2, 1}}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{"string", "double"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestType_QueryingByMinKeyAndMaxKey(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("zipCode", query.Type("minKey"))),
		agg.MatchStage(query.Field("zipCode", query.Type("maxKey"))),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{"minKey"}}}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{"maxKey"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestType_QueryingByArrayType(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("zipCode", query.Type("array"))),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "zipCode", Value: bson.D{{Key: "$type", Value: bson.A{"array"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}
