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

func TestBitsAllClear_BitPositionArray(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAllClear([]int{1, 5})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAllClear", Value: bson.A{1, 5}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAllClear_IntegerBitmask(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAllClear(35)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAllClear", Value: 35}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAllClear_BinDataBitmask(t *testing.T) {
	bitmask := bson.Binary{Subtype: 0x00, Data: []byte{0x20}}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAllClear(bitmask)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAllClear", Value: bitmask}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAllSet_BitPositionArray(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAllSet([]int{1, 5})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAllSet", Value: bson.A{1, 5}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAllSet_IntegerBitmask(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAllSet(50)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAllSet", Value: 50}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAllSet_BinDataBitmask(t *testing.T) {
	bitmask := bson.Binary{Subtype: 0x00, Data: []byte{0x20}}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAllSet(bitmask)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAllSet", Value: bitmask}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAnyClear_BitPositionArray(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAnyClear([]int{1, 5})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAnyClear", Value: bson.A{1, 5}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAnyClear_IntegerBitmask(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAnyClear(35)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAnyClear", Value: 35}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAnyClear_BinDataBitmask(t *testing.T) {
	bitmask := bson.Binary{Subtype: 0x00, Data: []byte{0x20}}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAnyClear(bitmask)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAnyClear", Value: bitmask}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAnySet_BitPositionArray(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAnySet([]int{1, 5})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAnySet", Value: bson.A{1, 5}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAnySet_IntegerBitmask(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAnySet(35)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAnySet", Value: 35}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitsAnySet_BinDataBitmask(t *testing.T) {
	bitmask := bson.Binary{Subtype: 0x00, Data: []byte{0x20}}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("a", query.BitsAnySet(bitmask)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "a", Value: bson.D{{Key: "$bitsAnySet", Value: bitmask}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBox(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoWithin(
				query.Box([]float64{0, 0}, []float64{3, 6}),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoWithin", Value: bson.D{
				{Key: "$box", Value: bson.A{bson.A{0.0, 0.0}, bson.A{3.0, 6.0}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCenter(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoWithin(
				query.Center([]float64{-74, 40.74}, 10),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoWithin", Value: bson.D{
				{Key: "$center", Value: bson.A{bson.A{-74.0, 40.74}, 10.0}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCenterSphere(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoWithin(
				query.CenterSphere([]float64{-88, 30}, 0.01),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoWithin", Value: bson.D{
				{Key: "$centerSphere", Value: bson.A{bson.A{-88.0, 30.0}, 0.01}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestComment(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("x", query.Gt(0)),
			query.Comment("Don't allow negative inputs."),
		),
		agg.GroupStage(
			agg.Mod("$x", 2),
			agg.Accumulate("total", agg.SumAccumulator("$x")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "x", Value: bson.D{{Key: "$gt", Value: 0}}},
			{Key: "$comment", Value: "Don't allow negative inputs."},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$mod", Value: bson.A{"$x", 2}}}},
			{Key: "total", Value: bson.D{{Key: "$sum", Value: "$x"}}},
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

func TestExpr_CompareTwoFieldsFromSingleDoc(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Expr(agg.Gt("$spent", "$budget")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$expr", Value: bson.D{{Key: "$gt", Value: bson.A{"$spent", "$budget"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestExpr_ConditionalStatements when $cond is implemented

func TestGeoIntersects_BigPolygon(t *testing.T) {
	coords := [][][]float64{{{-100, 60}, {-100, 0}, {-100, -60}, {100, -60}, {100, 60}, {-100, 60}}}
	crs := bson.D{
		{Key: "type", Value: "name"},
		{Key: "properties", Value: bson.D{
			{Key: "name", Value: "urn:x-mongodb:crs:strictwinding:EPSG:4326"},
		}},
	}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoIntersects(
				query.GeoJSON("Polygon", coords, query.WithGeoJSONCRS(crs)),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoIntersects", Value: bson.D{
				{Key: "$geometry", Value: bson.D{
					{Key: "type", Value: "Polygon"},
					{Key: "coordinates", Value: bson.A{bson.A{
						bson.A{-100.0, 60.0}, bson.A{-100.0, 0.0}, bson.A{-100.0, -60.0},
						bson.A{100.0, -60.0}, bson.A{100.0, 60.0}, bson.A{-100.0, 60.0},
					}}},
					{Key: "crs", Value: crs},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGeoIntersects_Polygon(t *testing.T) {
	coords := [][][]float64{{{0, 0}, {3, 6}, {6, 1}, {0, 0}}}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoIntersects(query.GeoJSON("Polygon", coords))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoIntersects", Value: bson.D{
				{Key: "$geometry", Value: bson.D{
					{Key: "type", Value: "Polygon"},
					{Key: "coordinates", Value: bson.A{bson.A{
						bson.A{0.0, 0.0}, bson.A{3.0, 6.0}, bson.A{6.0, 1.0}, bson.A{0.0, 0.0},
					}}},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGeoJSON_Point(t *testing.T) {
	coords := []float64{-73.9667, 40.78}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoIntersects(query.GeoJSON("Point", coords))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoIntersects", Value: bson.D{
				{Key: "$geometry", Value: bson.D{
					{Key: "type", Value: "Point"},
					{Key: "coordinates", Value: bson.A{-73.9667, 40.78}},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGeoJSON_Polygon(t *testing.T) {
	coords := [][][]float64{{{0, 0}, {3, 6}, {6, 1}, {0, 0}}}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoWithin(query.GeoJSON("Polygon", coords))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoWithin", Value: bson.D{
				{Key: "$geometry", Value: bson.D{
					{Key: "type", Value: "Polygon"},
					{Key: "coordinates", Value: bson.A{bson.A{
						bson.A{0.0, 0.0}, bson.A{3.0, 6.0}, bson.A{6.0, 1.0}, bson.A{0.0, 0.0},
					}}},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGeoJSON_WithCRS(t *testing.T) {
	coords := [][][]float64{{{-100, 60}, {-100, -60}, {100, -60}, {100, 60}, {-100, 60}}}
	crs := bson.D{
		{Key: "type", Value: "name"},
		{Key: "properties", Value: bson.D{
			{Key: "name", Value: "urn:x-mongodb:crs:strictwinding:EPSG:4326"},
		}},
	}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoWithin(
				query.GeoJSON("Polygon", coords, query.WithGeoJSONCRS(crs)),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoWithin", Value: bson.D{
				{Key: "$geometry", Value: bson.D{
					{Key: "type", Value: "Polygon"},
					{Key: "coordinates", Value: bson.A{bson.A{
						bson.A{-100.0, 60.0}, bson.A{-100.0, -60.0}, bson.A{100.0, -60.0},
						bson.A{100.0, 60.0}, bson.A{-100.0, 60.0},
					}}},
					{Key: "crs", Value: crs},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGeoWithin_Polygon(t *testing.T) {
	coords := [][][]float64{{{0, 0}, {3, 6}, {6, 1}, {0, 0}}}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoWithin(query.GeoJSON("Polygon", coords))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoWithin", Value: bson.D{
				{Key: "$geometry", Value: bson.D{
					{Key: "type", Value: "Polygon"},
					{Key: "coordinates", Value: bson.A{bson.A{
						bson.A{0.0, 0.0}, bson.A{3.0, 6.0}, bson.A{6.0, 1.0}, bson.A{0.0, 0.0},
					}}},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGeoWithin_BigPolygon(t *testing.T) {
	coords := [][][]float64{{{-100, 60}, {-100, 0}, {-100, -60}, {100, -60}, {100, 60}, {-100, 60}}}
	crs := bson.D{
		{Key: "type", Value: "name"},
		{Key: "properties", Value: bson.D{
			{Key: "name", Value: "urn:x-mongodb:crs:strictwinding:EPSG:4326"},
		}},
	}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoWithin(
				query.GeoJSON("Polygon", coords, query.WithGeoJSONCRS(crs)),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoWithin", Value: bson.D{
				{Key: "$geometry", Value: bson.D{
					{Key: "type", Value: "Polygon"},
					{Key: "coordinates", Value: bson.A{bson.A{
						bson.A{-100.0, 60.0}, bson.A{-100.0, 0.0}, bson.A{-100.0, -60.0},
						bson.A{100.0, -60.0}, bson.A{100.0, 60.0}, bson.A{-100.0, 60.0},
					}}},
					{Key: "crs", Value: crs},
				}},
			}}}},
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

func TestJSONSchema(t *testing.T) {
	schema := bson.D{
		{Key: "required", Value: bson.A{"name", "major", "gpa", "address"}},
		{Key: "properties", Value: bson.D{
			{Key: "name", Value: bson.D{
				{Key: "bsonType", Value: "string"},
				{Key: "description", Value: "must be a string and is required"},
			}},
			{Key: "address", Value: bson.D{
				{Key: "bsonType", Value: "object"},
				{Key: "required", Value: bson.A{"zipcode"}},
				{Key: "properties", Value: bson.D{
					{Key: "street", Value: bson.D{{Key: "bsonType", Value: "string"}}},
					{Key: "zipcode", Value: bson.D{{Key: "bsonType", Value: "string"}}},
				}},
			}},
		}},
	}
	got := agg.Pipeline{
		agg.MatchStage(
			query.JSONSchema(schema),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$jsonSchema", Value: schema},
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

func TestMaxDistance(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("loc", query.MaxDistance(5000))),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$maxDistance", Value: 5000}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMinDistance(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("loc", query.MinDistance((1000)))),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$minDistance", Value: 1000}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMod_SelectDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("qty", query.Mod(4, 0)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "qty", Value: bson.D{{Key: "$mod", Value: bson.A{4, 0}}}},
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

func TestNear(t *testing.T) {
	coords := []float64{-73.9667, 40.78}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("location", query.Near(
				query.GeoJSON("Point", coords),
				query.WithNearMinDistance(1000),
				query.WithNearMaxDistance(5000),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "location", Value: bson.D{{Key: "$near", Value: bson.D{
				{Key: "$geometry", Value: bson.D{
					{Key: "type", Value: "Point"},
					{Key: "coordinates", Value: bson.A{-73.9667, 40.78}},
				}},
				{Key: "$minDistance", Value: 1000},
				{Key: "$maxDistance", Value: 5000},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNearSphere(t *testing.T) {
	coords := []float64{-73.9667, 40.78}
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("location", query.NearSphere(
				query.GeoJSON("Point", coords),
				query.WithNearMinDistance(1000),
				query.WithNearMaxDistance(5000),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "location", Value: bson.D{{Key: "$nearSphere", Value: bson.D{
				{Key: "$geometry", Value: bson.D{
					{Key: "type", Value: "Point"},
					{Key: "coordinates", Value: bson.A{-73.9667, 40.78}},
				}},
				{Key: "$minDistance", Value: 1000},
				{Key: "$maxDistance", Value: 5000},
			}}}},
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

func TestPolygon(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("loc", query.GeoWithin(
				query.Polygon([]float64{0, 0}, []float64{3, 6}, []float64{6, 0}),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "loc", Value: bson.D{{Key: "$geoWithin", Value: bson.D{
				{Key: "$polygon", Value: bson.A{
					bson.A{0.0, 0.0}, bson.A{3.0, 6.0}, bson.A{6.0, 0.0},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRegex_PerformLikeMatch(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("sku", query.Regex(bson.Regex{Pattern: "789$"})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "sku", Value: bson.D{{Key: "$regex", Value: bson.Regex{Pattern: "789$"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRegex_PerformCaseInsensitiveMatch(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("sku", query.Regex(bson.Regex{Pattern: "^ABC", Options: "i"})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "sku", Value: bson.D{{Key: "$regex", Value: bson.Regex{Pattern: "^ABC", Options: "i"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement $sampleRate tests when $count stage is implemented

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

func TestText_SearchForSingleWord(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Text("coffee"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$text", Value: bson.D{{Key: "$search", Value: "coffee"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestText_MatchAnyOfSearchTerms(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Text("bake coffee cake"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$text", Value: bson.D{{Key: "$search", Value: "bake coffee cake"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestText_SearchDifferentLanguage(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Text("leche", query.WithTextLanguage("es")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$text", Value: bson.D{
				{Key: "$search", Value: "leche"},
				{Key: "$language", Value: "es"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestText_CaseAndDiacriticInsensitiveSearch(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Text("сы́рники CAFÉS"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$text", Value: bson.D{
				{Key: "$search", Value: "сы́рники CAFÉS"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestText_CaseSensitiveSearch(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Text("Coffee", query.WithTextCaseSensitive(true)),
		),
		agg.MatchStage(
			query.Text(`"Café Con Leche"`, query.WithTextCaseSensitive(true)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$text", Value: bson.D{
				{Key: "$search", Value: "Coffee"},
				{Key: "$caseSensitive", Value: true},
			}},
		}}},
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$text", Value: bson.D{
				{Key: "$search", Value: `"Café Con Leche"`},
				{Key: "$caseSensitive", Value: true},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestText_DiacriticSensitiveSearch(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Text("CAFÉ", query.WithTextDiacriticSensitive(true)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "$text", Value: bson.D{
				{Key: "$search", Value: "CAFÉ"},
				{Key: "$diacriticSensitive", Value: true},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: test TestText_TextSearchScoreExamples when $meta is implemented

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
