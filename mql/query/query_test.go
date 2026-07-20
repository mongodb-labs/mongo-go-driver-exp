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

func TestAll_MatchValues(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("tags", query.All("appliance", "school", "book")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "tags", Value: bson.D{{Key: "$all", Value: bson.A{"appliance", "school", "book"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAll_WithElemMatch(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.All(
				query.ElemMatch(
					query.Field("size", query.Eq("M")),
					query.Field("num", query.Gt(50)),
				),
				query.ElemMatch(
					query.Field("num", query.Eq(100)),
					query.Field("color", query.Eq("green")),
				),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$all", Value: bson.A{
				bson.D{{Key: "$elemMatch", Value: bson.D{
					{Key: "size", Value: bson.D{{Key: "$eq", Value: "M"}}},
					{Key: "num", Value: bson.D{{Key: "$gt", Value: 50}}},
				}}},
				bson.D{{Key: "$elemMatch", Value: bson.D{
					{Key: "num", Value: bson.D{{Key: "$eq", Value: 100}}},
					{Key: "color", Value: bson.D{{Key: "$eq", Value: "green"}}},
				}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAnd_MultipleExpressionsSameField(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.And(
				query.Field("price", query.Ne(1.99)),
				query.Field("price", query.Exists(true)),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$and", Value: bson.A{
				bson.D{{Key: "price", Value: bson.D{{Key: "$ne", Value: 1.99}}}},
				bson.D{{Key: "price", Value: bson.D{{Key: "$exists", Value: true}}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAnd_MultipleExpressionsSameOperator(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.And(
				query.Or(
					query.Field("qty", query.Lt(10)),
					query.Field("qty", query.Gt(50)),
				),
				query.Or(
					query.Field("sale", query.Eq(true)),
					query.Field("price", query.Lt(5)),
				),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$and", Value: bson.A{
				bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "qty", Value: bson.D{{Key: "$lt", Value: 10}}}},
					bson.D{{Key: "qty", Value: bson.D{{Key: "$gt", Value: 50}}}},
				}}},
				bson.D{{Key: "$or", Value: bson.A{
					bson.D{{Key: "sale", Value: bson.D{{Key: "$eq", Value: true}}}},
					bson.D{{Key: "price", Value: bson.D{{Key: "$lt", Value: 5}}}},
				}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestElemMatch_ElementMatch(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("results", query.ElemMatch(query.Gte(80), query.Lt(85))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "results", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
				{Key: "$gte", Value: 80},
				{Key: "$lt", Value: 85},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestElemMatch_ArrayOfEmbeddedDocuments(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("results", query.ElemMatch(
				query.Field("product", query.Eq("xyz")),
				query.Field("score", query.Gte(8)),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "results", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
				{Key: "product", Value: bson.D{{Key: "$eq", Value: "xyz"}}},
				{Key: "score", Value: bson.D{{Key: "$gte", Value: 8}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestElemMatch_SingleQueryCondition(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("results", query.ElemMatch(
				query.Field("product", query.Ne("xyz")),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "results", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
				{Key: "product", Value: bson.D{{Key: "$ne", Value: "xyz"}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestElemMatch_WithOr(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("game", query.ElemMatch(
				query.Or(
					query.Field("score", query.Gt(10)),
					query.Field("score", query.Lt(5)),
				),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "game", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "score", Value: bson.D{{Key: "$gt", Value: 10}}}},
					bson.D{{Key: "score", Value: bson.D{{Key: "$lt", Value: 5}}}},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestElemMatch_SingleFieldOperator(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("results", query.ElemMatch(query.Gt(10))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "results", Value: bson.D{{Key: "$elemMatch", Value: bson.D{
				{Key: "$gt", Value: 10},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
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

func TestNor_TwoExpressions(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Nor(
				query.Field("price", query.Eq(1.99)),
				query.Field("sale", query.Eq(true)),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$nor", Value: bson.A{
				bson.D{{Key: "price", Value: bson.D{{Key: "$eq", Value: 1.99}}}},
				bson.D{{Key: "sale", Value: bson.D{{Key: "$eq", Value: true}}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNor_AdditionalComparisons(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Nor(
				query.Field("price", query.Eq(1.99)),
				query.Field("qty", query.Lt(20)),
				query.Field("sale", query.Eq(true)),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$nor", Value: bson.A{
				bson.D{{Key: "price", Value: bson.D{{Key: "$eq", Value: 1.99}}}},
				bson.D{{Key: "qty", Value: bson.D{{Key: "$lt", Value: 20}}}},
				bson.D{{Key: "sale", Value: bson.D{{Key: "$eq", Value: true}}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNor_AndExists(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Nor(
				query.Field("price", query.Eq(1.99)),
				query.Field("price", query.Exists(false)),
				query.Field("sale", query.Eq(true)),
				query.Field("sale", query.Exists(false)),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$nor", Value: bson.A{
				bson.D{{Key: "price", Value: bson.D{{Key: "$eq", Value: 1.99}}}},
				bson.D{{Key: "price", Value: bson.D{{Key: "$exists", Value: false}}}},
				bson.D{{Key: "sale", Value: bson.D{{Key: "$eq", Value: true}}}},
				bson.D{{Key: "sale", Value: bson.D{{Key: "$exists", Value: false}}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNot_Syntax(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("price", query.Not(query.Gt(1.99))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "price", Value: bson.D{{Key: "$not", Value: bson.D{{Key: "$gt", Value: 1.99}}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNot_RegularExpressions(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("price", query.Not(bson.Regex{Pattern: "^p.*"})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "price", Value: bson.D{{Key: "$not", Value: bson.Regex{Pattern: "^p.*"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestOr_Clauses(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Or(
				query.Field("quantity", query.Lt(20)),
				query.Field("price", query.Eq(10)),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$or", Value: bson.A{
				bson.D{{Key: "quantity", Value: bson.D{{Key: "$lt", Value: 20}}}},
				bson.D{{Key: "price", Value: bson.D{{Key: "$eq", Value: 10}}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestOr_ErrorHandling when $expr is implemented

func TestSize_QueryAnArrayByLength(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("tags", query.Size(3)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "tags", Value: bson.D{{Key: "$size", Value: 3}}},
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
