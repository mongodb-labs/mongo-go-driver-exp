package agg_test

import (
	"testing"

	"github.com/mongodb-labs/mongo-go-driver-exp/mql/agg"
	"github.com/mongodb-labs/mongo-go-driver-exp/mql/query"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestAbs(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("delta", agg.Abs(agg.Subtract("$startTemp", "$endTemp"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "delta", Value: bson.D{
				{Key: "$abs", Value: bson.D{
					{Key: "$subtract", Value: bson.A{"$startTemp", "$endTemp"}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAcos(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("angle_a", agg.RadiansToDegrees(agg.Acos(agg.Divide("$side_b", "$hypotenuse")))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "angle_a", Value: bson.D{
				{Key: "$radiansToDegrees", Value: bson.D{
					{Key: "$acos", Value: bson.D{
						{Key: "$divide", Value: bson.A{"$side_b", "$hypotenuse"}},
					}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAcosh(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("y-coordinate", agg.RadiansToDegrees(agg.Acosh("$x-coordinate"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "y-coordinate", Value: bson.D{
				{Key: "$radiansToDegrees", Value: bson.D{
					{Key: "$acosh", Value: "$x-coordinate"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAdd_AddNumbers(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("total", agg.Add("$price", "$fee")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "total", Value: bson.D{
				{Key: "$add", Value: bson.A{"$price", "$fee"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestAdd_PerformAdditionOnDate when date functionality is implemented

func TestAllElementsTrue(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("responses"),
			agg.Compute("isAllTrue", agg.AllElementsTrue("$responses")),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "responses", Value: 1},
			{Key: "isAllTrue", Value: bson.D{{Key: "$allElementsTrue", Value: bson.A{"$responses"}}}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAnd(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Include("qty"),
			agg.Compute("result", agg.And(agg.Gt("$qty", 100), agg.Lt("$qty", 250))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "qty", Value: 1},
			{Key: "result", Value: bson.D{
				{Key: "$and", Value: bson.A{
					bson.D{{Key: "$gt", Value: bson.A{"$qty", 100}}},
					bson.D{{Key: "$lt", Value: bson.A{"$qty", 250}}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAnyElementTrue(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("responses"),
			agg.Compute("isAllTrue", agg.AnyElementTrue("$responses")),
			agg.Include("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "responses", Value: 1},
			{Key: "isAllTrue", Value: bson.D{{Key: "$anyElementTrue", Value: bson.A{"$responses"}}}},
			{Key: "_id", Value: 1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestArrayElemAt(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("first", agg.ArrayElemAt("$favorites", 0)),
			agg.Compute("last", agg.ArrayElemAt("$favorites", -1)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "first", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$favorites", 0}}}},
			{Key: "last", Value: bson.D{{Key: "$arrayElemAt", Value: bson.A{"$favorites", -1}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestArrayToObject_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("dimensions", agg.ArrayToObject("$dimensions")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "dimensions", Value: bson.D{
				{Key: "$arrayToObject", Value: "$dimensions"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestArrayToObject_ObjectToArray when ObjectToArray
// operator is implemented

func TestAsin(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("angle_a", agg.RadiansToDegrees(agg.Asin(agg.Divide("$side_a", "$hypotenuse")))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "angle_a", Value: bson.D{
				{Key: "$radiansToDegrees", Value: bson.D{
					{Key: "$asin", Value: bson.D{
						{Key: "$divide", Value: bson.A{"$side_a", "$hypotenuse"}},
					}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAsinh(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("y-coordinate", agg.RadiansToDegrees(agg.Asinh("$x-coordinate"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "y-coordinate", Value: bson.D{
				{Key: "$radiansToDegrees", Value: bson.D{
					{Key: "$asinh", Value: "$x-coordinate"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAtan(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("angle_a", agg.RadiansToDegrees(agg.Atan(agg.Divide("$side_b", "$side_a")))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "angle_a", Value: bson.D{
				{Key: "$radiansToDegrees", Value: bson.D{
					{Key: "$atan", Value: bson.D{
						{Key: "$divide", Value: bson.A{"$side_b", "$side_a"}},
					}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAtan2(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("angle_a", agg.RadiansToDegrees(agg.Atan2("$side_b", "$side_a"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "angle_a", Value: bson.D{
				{Key: "$radiansToDegrees", Value: bson.D{
					{Key: "$atan2", Value: bson.A{"$side_b", "$side_a"}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAtanh(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("y-coordinate", agg.RadiansToDegrees(agg.Atanh("$x-coordinate"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "y-coordinate", Value: bson.D{
				{Key: "$radiansToDegrees", Value: bson.D{
					{Key: "$atanh", Value: "$x-coordinate"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestAvg(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("quizAvg", agg.Avg("$quizzes")),
			agg.Compute("labAvg", agg.Avg("$labs")),
			agg.Compute("examAvg", agg.Avg("$final", "$midterm")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "quizAvg", Value: bson.D{{Key: "$avg", Value: bson.A{"$quizzes"}}}},
			{Key: "labAvg", Value: bson.D{{Key: "$avg", Value: bson.A{"$labs"}}}},
			{Key: "examAvg", Value: bson.D{{Key: "$avg", Value: bson.A{"$final", "$midterm"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitAnd_TwoIntegers(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("result", agg.BitAnd("$a", "$b")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "result", Value: bson.D{
				{Key: "$bitAnd", Value: bson.A{"$a", "$b"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitAnd_LongAndInteger(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("result", agg.BitAnd("$a", int64(63))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "result", Value: bson.D{
				{Key: "$bitAnd", Value: bson.A{"$a", int64(63)}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitNot(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("result", agg.BitNot("$a")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "result", Value: bson.D{
				{Key: "$bitNot", Value: "$a"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitOr_TwoIntegers(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("result", agg.BitOr("$a", "$b")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "result", Value: bson.D{
				{Key: "$bitOr", Value: bson.A{"$a", "$b"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitOr_LongAndInteger(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("result", agg.BitOr("$a", int64(63))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "result", Value: bson.D{
				{Key: "$bitOr", Value: bson.A{"$a", int64(63)}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBitXor(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("result", agg.BitXor("$a", "$b")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "result", Value: bson.D{
				{Key: "$bitXor", Value: bson.A{"$a", "$b"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBottom(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("bottomScore", agg.Bottom(
				"$results",
				[]any{"$playerId", "$score"},
				agg.Sort("score", agg.Desc),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "bottomScore", Value: bson.D{{Key: "$bottom", Value: bson.D{
				{Key: "sortBy", Value: bson.D{{Key: "score", Value: -1}}},
				{Key: "output", Value: bson.A{"$playerId", "$score"}},
				{Key: "input", Value: "$results"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBottomN(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("bottomScores", agg.BottomN(
				3,
				"$results",
				[]any{"$playerId", "$score"},
				agg.Sort("score", agg.Desc),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "bottomScores", Value: bson.D{{Key: "$bottomN", Value: bson.D{
				{Key: "n", Value: 3},
				{Key: "sortBy", Value: bson.D{{Key: "score", Value: -1}}},
				{Key: "output", Value: []any{"$playerId", "$score"}},
				{Key: "input", Value: "$results"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCeil(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("value"),
			agg.Compute("ceilingValue", agg.Ceil("$value")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "value", Value: 1},
			{Key: "ceilingValue", Value: bson.D{
				{Key: "$ceil", Value: "$value"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCmp(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Include("qty"),
			agg.Compute("cmpTo250", agg.Cmp("$qty", 250)),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "qty", Value: 1},
			{Key: "cmpTo250", Value: bson.D{
				{Key: "$cmp", Value: bson.A{"$qty", 250}},
			}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestConcat(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("itemDescription", agg.Concat("$item", " - ", "$description")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "itemDescription", Value: bson.D{
				{Key: "$concat", Value: bson.A{"$item", " - ", "$description"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestConcatArrays(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("items", agg.ConcatArrays("$instock", "$ordered")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "items", Value: bson.D{{Key: "$concatArrays", Value: bson.A{"$instock", "$ordered"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCos(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("side_a", agg.Multiply(agg.Cos(agg.DegreesToRadians("$angle_a")), "$hypotenuse")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "side_a", Value: bson.D{
				{Key: "$multiply", Value: bson.A{
					bson.D{{Key: "$cos", Value: bson.D{
						{Key: "$degreesToRadians", Value: "$angle_a"},
					}}},
					"$hypotenuse",
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestCosh(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("cosh_output", agg.Cosh(agg.DegreesToRadians("$angle"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "cosh_output", Value: bson.D{
				{Key: "$cosh", Value: bson.D{
					{Key: "$degreesToRadians", Value: "$angle"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDegreesToRadians(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("angle_a_rad", agg.DegreesToRadians("$angle_a")),
			agg.Assign("angle_b_rad", agg.DegreesToRadians("$angle_b")),
			agg.Assign("angle_c_rad", agg.DegreesToRadians("$angle_c")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "angle_a_rad", Value: bson.D{{Key: "$degreesToRadians", Value: "$angle_a"}}},
			{Key: "angle_b_rad", Value: bson.D{{Key: "$degreesToRadians", Value: "$angle_b"}}},
			{Key: "angle_c_rad", Value: bson.D{{Key: "$degreesToRadians", Value: "$angle_c"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDivide(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("city"),
			agg.Compute("workdays", agg.Divide("$hours", 8)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "city", Value: 1},
			{Key: "workdays", Value: bson.D{
				{Key: "$divide", Value: bson.A{"$hours", 8}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestEq(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Include("qty"),
			agg.Compute("qtyEq250", agg.Eq("$qty", 250)),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "qty", Value: 1},
			{Key: "qtyEq250", Value: bson.D{
				{Key: "$eq", Value: bson.A{"$qty", 250}},
			}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestExp(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("effectiveRate", agg.Subtract(agg.Exp("$interestRate"), 1)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "effectiveRate", Value: bson.D{
				{Key: "$subtract", Value: bson.A{
					bson.D{{Key: "$exp", Value: "$interestRate"}},
					1,
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFilterArray_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("items", agg.FilterArray("$items", "item", agg.Gte("$$item.price", 100))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "items", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$items"},
					{Key: "as", Value: "item"},
					{Key: "cond", Value: bson.D{
						{Key: "$gte", Value: bson.A{"$$item.price", 100}},
					}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFilterArray_UsingLimitField(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("items", agg.FilterArray("$items", "item", agg.Gte("$$item.price", 100), 1)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "items", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$items"},
					{Key: "as", Value: "item"},
					{Key: "cond", Value: bson.D{
						{Key: "$gte", Value: bson.A{"$$item.price", 100}},
					}},
					{Key: "limit", Value: 1},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFilterArray_LimitGreaterThanPossibleMatches(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("items", agg.FilterArray("$items", "item", agg.Gte("$$item.price", 100), 5)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "items", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$items"},
					{Key: "as", Value: "item"},
					{Key: "cond", Value: bson.D{
						{Key: "$gte", Value: bson.A{"$$item.price", 100}},
					}},
					{Key: "limit", Value: 5},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFirst(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("firstItem", agg.First("$items")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "firstItem", Value: bson.D{{Key: "$first", Value: "$items"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFirstN(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("firstScores", agg.FirstN(3, "$score")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "firstScores", Value: bson.D{{Key: "$firstN", Value: bson.D{
				{Key: "n", Value: 3},
				{Key: "input", Value: "$score"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFloor(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("value"),
			agg.Compute("floorValue", agg.Floor("$value")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "value", Value: 1},
			{Key: "floorValue", Value: bson.D{
				{Key: "$floor", Value: "$value"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGt(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Include("qty"),
			agg.Compute("qtyGt250", agg.Gt("$qty", 250)),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "qty", Value: 1},
			{Key: "qtyGt250", Value: bson.D{
				{Key: "$gt", Value: bson.A{"$qty", 250}},
			}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestGte(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Include("qty"),
			agg.Compute("qtyGte250", agg.Gte("$qty", 250)),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "qty", Value: 1},
			{Key: "qtyGte250", Value: bson.D{
				{Key: "$gte", Value: bson.A{"$qty", 250}},
			}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIfNull_SingleInputExpr(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("year", query.Lt(1910)),
		),
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("title"),
			agg.Compute("rated", agg.IfNull("$rated", "Not Rated")),
		),
		agg.SortStage(
			agg.Sort("title", agg.Asc),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$lt", Value: 1910}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "title", Value: 1},
			{Key: "rated", Value: bson.D{
				{Key: "$ifNull", Value: bson.A{"$rated", "Not Rated"}},
			}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "title", Value: 1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIfNull_MultipleInputExpr(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("year", query.Lt(1910)),
		),
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("title"),
			agg.Compute("rating", agg.IfNull("$tomatoes.critic.rating", "$tomatoes.viewer.rating", 0)),
		),
		agg.SortStage(
			agg.Sort("title", agg.Asc),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$lt", Value: 1910}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "title", Value: 1},
			{Key: "rating", Value: bson.D{
				{Key: "$ifNull", Value: bson.A{"$tomatoes.critic.rating", "$tomatoes.viewer.rating", 0}},
			}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "title", Value: 1},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIn(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("store location", "$location"),
			agg.Compute("has bananas", agg.In("bananas", "$in_stock")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "store location", Value: "$location"},
			{Key: "has bananas", Value: bson.D{
				{Key: "$in", Value: bson.A{"bananas", "$in_stock"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfArray_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("index", agg.IndexOfArray("$items", 2, nil, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "index", Value: bson.D{{Key: "$indexOfArray", Value: bson.A{"$items", 2}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfArray_StartOption(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("index", agg.IndexOfArray("$items", 2, 1, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "index", Value: bson.D{{Key: "$indexOfArray", Value: bson.A{"$items", 2, 1}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfArray_EndOption(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("index", agg.IndexOfArray("$items", 2, nil, 4)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "index", Value: bson.D{{Key: "$indexOfArray", Value: bson.A{"$items", 2, 0, 4}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfArray_StartAndEndOption(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("index", agg.IndexOfArray("$items", 2, 1, 4)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "index", Value: bson.D{{Key: "$indexOfArray", Value: bson.A{"$items", 2, 1, 4}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfBytes_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("byteLocation", agg.IndexOfBytes("$item", "foo", nil, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "byteLocation", Value: bson.D{
				{Key: "$indexOfBytes", Value: bson.A{"$item", "foo"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfBytes_StartOption(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("byteLocation", agg.IndexOfBytes("$item", "foo", 1, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "byteLocation", Value: bson.D{
				{Key: "$indexOfBytes", Value: bson.A{"$item", "foo", 1}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfBytes_EndOption(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("byteLocation", agg.IndexOfBytes("$item", "foo", nil, 4)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "byteLocation", Value: bson.D{
				{Key: "$indexOfBytes", Value: bson.A{"$item", "foo", 0, 4}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfBytes_StartAndEndOption(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("byteLocation", agg.IndexOfBytes("$item", "foo", 1, 4)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "byteLocation", Value: bson.D{
				{Key: "$indexOfBytes", Value: bson.A{"$item", "foo", 1, 4}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfCP_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("cpLocation", agg.IndexOfCP("$item", "foo", nil, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "cpLocation", Value: bson.D{{Key: "$indexOfCP", Value: bson.A{"$item", "foo"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfCP_StartOption(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("cpLocation", agg.IndexOfCP("$item", "foo", 1, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "cpLocation", Value: bson.D{{Key: "$indexOfCP", Value: bson.A{"$item", "foo", 1}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfCP_EndOption(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("cpLocation", agg.IndexOfCP("$item", "foo", nil, 4)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "cpLocation", Value: bson.D{{Key: "$indexOfCP", Value: bson.A{"$item", "foo", 0, 4}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfCP_StartAndEndOption(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("cpLocation", agg.IndexOfCP("$item", "foo", 1, 4)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "cpLocation", Value: bson.D{{Key: "$indexOfCP", Value: bson.A{"$item", "foo", 1, 4}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLast(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("lastItem", agg.Last("$items")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "lastItem", Value: bson.D{{Key: "$last", Value: "$items"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLastN(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("lastScores", agg.LastN(3, "$score")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "lastScores", Value: bson.D{{Key: "$lastN", Value: bson.D{
				{Key: "n", Value: 3},
				{Key: "input", Value: "$score"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLn(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("x", "$year"),
			agg.Compute("y", agg.Ln("$sales")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "x", Value: "$year"},
			{Key: "y", Value: bson.D{
				{Key: "$ln", Value: "$sales"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLog(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("bitsNeeded", agg.Floor(agg.Add(1, agg.Log("$int", 2)))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "bitsNeeded", Value: bson.D{
				{Key: "$floor", Value: bson.D{
					{Key: "$add", Value: bson.A{1, bson.D{{Key: "$log", Value: bson.A{"$int", 2}}}}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLog10(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("pH", agg.Multiply(-1, agg.Log10("$H3O"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "pH", Value: bson.D{
				{Key: "$multiply", Value: bson.A{-1, bson.D{{Key: "$log10", Value: "$H3O"}}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLt(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Include("qty"),
			agg.Compute("qtyLt250", agg.Lt("$qty", 250)),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "qty", Value: 1},
			{Key: "qtyLt250", Value: bson.D{
				{Key: "$lt", Value: bson.A{"$qty", 250}},
			}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLte(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Include("qty"),
			agg.Compute("qtyLte250", agg.Lte("$qty", 250)),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "qty", Value: 1},
			{Key: "qtyLte250", Value: bson.D{
				{Key: "$lte", Value: bson.A{"$qty", 250}},
			}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLtrim(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("description", agg.Ltrim("$description", nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "description", Value: bson.D{{Key: "$ltrim", Value: bson.D{
				{Key: "input", Value: "$description"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMax(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("quizMax", agg.Max("$quizzes")),
			agg.Compute("labMax", agg.Max("$labs")),
			agg.Compute("examMax", agg.Max("$final", "$midterm")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "quizMax", Value: bson.D{{Key: "$max", Value: bson.A{"$quizzes"}}}},
			{Key: "labMax", Value: bson.D{{Key: "$max", Value: bson.A{"$labs"}}}},
			{Key: "examMax", Value: bson.D{{Key: "$max", Value: bson.A{"$final", "$midterm"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMaxN(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("maxScores", agg.MaxN(2, "$score")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "maxScores", Value: bson.D{{Key: "$maxN", Value: bson.D{
				{Key: "n", Value: 2},
				{Key: "input", Value: "$score"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMin(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("quizMin", agg.Min("$quizzes")),
			agg.Compute("labMin", agg.Min("$labs")),
			agg.Compute("examMin", agg.Min("$final", "$midterm")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "quizMin", Value: bson.D{{Key: "$min", Value: bson.A{"$quizzes"}}}},
			{Key: "labMin", Value: bson.D{{Key: "$min", Value: bson.A{"$labs"}}}},
			{Key: "examMin", Value: bson.D{{Key: "$min", Value: bson.A{"$final", "$midterm"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMinN(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("minScores", agg.MinN(2, "$score")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "minScores", Value: bson.D{{Key: "$minN", Value: bson.D{
				{Key: "n", Value: 2},
				{Key: "input", Value: "$score"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMod(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("remainder", agg.Mod("$hours", "$tasks")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "remainder", Value: bson.D{
				{Key: "$mod", Value: bson.A{"$hours", "$tasks"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMultiply(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("date"),
			agg.Include("item"),
			agg.Compute("total", agg.Multiply("$price", "$quantity")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "date", Value: 1},
			{Key: "item", Value: 1},
			{Key: "total", Value: bson.D{
				{Key: "$multiply", Value: bson.A{"$price", "$quantity"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNe(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Include("qty"),
			agg.Compute("qtyNe250", agg.Ne("$qty", 250)),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "qty", Value: 1},
			{Key: "qtyNe250", Value: bson.D{
				{Key: "$ne", Value: bson.A{"$qty", 250}},
			}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestNot(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("result", agg.Not(agg.Gt("$qty", 250))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "result", Value: bson.D{
				{Key: "$not", Value: bson.A{
					bson.D{{Key: "$gt", Value: bson.A{"$qty", 250}}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestOr(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("result", agg.Or(agg.Gt("$qty", 250), agg.Lt("$qty", 200))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "result", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "$gt", Value: bson.A{"$qty", 250}}},
					bson.D{{Key: "$lt", Value: bson.A{"$qty", 200}}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestPow(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("variance", agg.Pow(agg.StdDevPop("$scores.score"), 2)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "variance", Value: bson.D{
				{Key: "$pow", Value: bson.A{
					bson.D{{Key: "$stdDevPop", Value: bson.A{"$scores.score"}}},
					2,
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRadiansToDegrees(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("angle_a_deg", agg.RadiansToDegrees("$angle_a")),
			agg.Assign("angle_b_deg", agg.RadiansToDegrees("$angle_b")),
			agg.Assign("angle_c_deg", agg.RadiansToDegrees("$angle_c")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "angle_a_deg", Value: bson.D{{Key: "$radiansToDegrees", Value: "$angle_a"}}},
			{Key: "angle_b_deg", Value: bson.D{{Key: "$radiansToDegrees", Value: "$angle_b"}}},
			{Key: "angle_c_deg", Value: bson.D{{Key: "$radiansToDegrees", Value: "$angle_c"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRange(t *testing.T) {
	step := 25
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("city"),
			agg.Compute("Rest stops", agg.Range(0, "$distance", &step)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "city", Value: 1},
			{Key: "Rest stops", Value: bson.D{{Key: "$range", Value: bson.A{0, "$distance", 25}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRegexFind_AndItsOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFind("$description", bson.Regex{Pattern: "line"}, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "returnObject", Value: bson.D{{Key: "$regexFind", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: bson.Regex{Pattern: "line"}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRegexFind_IOption(t *testing.T) {
	opts := "i"
	got := agg.Pipeline{
		// Specify i as part of the Regex type
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFind("$description", bson.Regex{Pattern: "line", Options: "i"}, nil)),
		),
		// Specify i in the options field
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFind("$description", "line", &opts)),
		),
		// Mix Regex type with options field
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFind("$description", bson.Regex{Pattern: "line"}, &opts)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "returnObject", Value: bson.D{{Key: "$regexFind", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: bson.Regex{Pattern: "line", Options: "i"}},
			}}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "returnObject", Value: bson.D{{Key: "$regexFind", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: "line"},
				{Key: "options", Value: "i"},
			}}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "returnObject", Value: bson.D{{Key: "$regexFind", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: bson.Regex{Pattern: "line"}},
				{Key: "options", Value: "i"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRegexFindAll_AndItsOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFindAll("$description", bson.Regex{Pattern: "line"}, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "returnObject", Value: bson.D{{Key: "$regexFindAll", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: bson.Regex{Pattern: "line"}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRegexFindAll_IOption(t *testing.T) {
	opts := "i"
	got := agg.Pipeline{
		// Specify i as part of the regex type
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFindAll("$description", bson.Regex{Pattern: "line", Options: "i"}, nil)),
		),
		// Specify i in the options field
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFindAll("$description", "line", &opts)),
		),
		// Mix Regex type with options field
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFindAll("$description", bson.Regex{Pattern: "line"}, &opts)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "returnObject", Value: bson.D{{Key: "$regexFindAll", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: bson.Regex{Pattern: "line", Options: "i"}},
			}}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "returnObject", Value: bson.D{{Key: "$regexFindAll", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: "line"},
				{Key: "options", Value: "i"},
			}}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "returnObject", Value: bson.D{{Key: "$regexFindAll", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: bson.Regex{Pattern: "line"}},
				{Key: "options", Value: "i"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRegexFindAll_ParseEmailFromString(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("email", agg.RegexFindAll("$comment", bson.Regex{
				Pattern: `[a-z0-9_.+-]+@[a-z0-9_.+-]+\.[a-z0-9_.+-]+`, Options: "i"}, nil)),
		),
		agg.SetStage(agg.Assign("email", "$email.match")),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "email", Value: bson.D{{Key: "$regexFindAll", Value: bson.D{
				{Key: "input", Value: "$comment"},
				{Key: "regex", Value: bson.Regex{Pattern: `[a-z0-9_.+-]+@[a-z0-9_.+-]+\.[a-z0-9_.+-]+`, Options: "i"}},
			}}}},
		}}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "email", Value: "$email.match"},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestRegexFindAll_UseCapturedGroupingsToParseUserName when $reduce is implemented

func TestRegexMatch_AndItsOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("result", agg.RegexMatch("$description", bson.Regex{Pattern: "line"}, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "result", Value: bson.D{{Key: "$regexMatch", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: bson.Regex{Pattern: "line"}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRegexMatch_IOption(t *testing.T) {
	opts := "i"
	got := agg.Pipeline{
		// Specify i as part of the Regex type
		agg.AddFieldsStage(
			agg.Assign("result", agg.RegexMatch("$description", bson.Regex{Pattern: "line", Options: "i"}, nil)),
		),
		// Specify i in the options field
		agg.AddFieldsStage(
			agg.Assign("result", agg.RegexMatch("$description", "line", &opts)),
		),
		// Mix Regex type with options field
		agg.AddFieldsStage(
			agg.Assign("result", agg.RegexMatch("$description", bson.Regex{Pattern: "line"}, &opts)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "result", Value: bson.D{{Key: "$regexMatch", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: bson.Regex{Pattern: "line", Options: "i"}},
			}}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "result", Value: bson.D{{Key: "$regexMatch", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: "line"},
				{Key: "options", Value: "i"},
			}}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "result", Value: bson.D{{Key: "$regexMatch", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "regex", Value: bson.Regex{Pattern: "line"}},
				{Key: "options", Value: "i"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestRegexMatch_CheckEmailAddress when $cond is implemented

func TestReplaceAll_ReplaceUsingString(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("item", agg.ReplaceAll("$item", "blue paint", "red paint")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: bson.D{{Key: "$replaceAll", Value: bson.D{
				{Key: "input", Value: "$item"},
				{Key: "find", Value: "blue paint"},
				{Key: "replacement", Value: "red paint"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReplaceAll_ReplaceUsingRegex(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("item", agg.ReplaceAll("$item", bson.Regex{Pattern: `\bblue paint\b`}, "red paint")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: bson.D{{Key: "$replaceAll", Value: bson.D{
				{Key: "input", Value: "$item"},
				{Key: "find", Value: bson.Regex{Pattern: `\bblue paint\b`}},
				{Key: "replacement", Value: "red paint"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReplaceOne_ReplaceUsingString(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("item", agg.ReplaceOne("$item", "blue paint", "red paint")),
    ),
  }
  want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: bson.D{{Key: "$replaceOne", Value: bson.D{
				{Key: "input", Value: "$item"},
				{Key: "find", Value: "blue paint"},
				{Key: "replacement", Value: "red paint"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReplaceOne_ReplaceUsingRegex(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("item", agg.ReplaceOne("$item", bson.Regex{Pattern: `\bblue paint\b`}, "red paint")),
  	),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
      	{Key: "item", Value: bson.D{{Key: "$replaceOne", Value: bson.D{
				{Key: "input", Value: "$item"},
				{Key: "find", Value: bson.Regex{Pattern: `\bblue paint\b`}},
				{Key: "replacement", Value: "red paint"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReverseArray(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("reverseFavorites", agg.ReverseArray("$favorites")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "reverseFavorites", Value: bson.D{{Key: "$reverseArray", Value: "$favorites"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRound(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("roundedValue", agg.Round("$value", 1)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "roundedValue", Value: bson.D{
				{Key: "$round", Value: bson.A{"$value", 1}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRtrim(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("description", agg.Rtrim("$description", nil)),
    ),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
        {Key: "item", Value: 1},
			  {Key: "description", Value: bson.D{{Key: "$rtrim", Value: bson.D{
			  {Key: "input", Value: "$description"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSetDifference(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("flowerFieldA"),
			agg.Include("flowerFieldB"),
			agg.Compute("inBOnly", agg.SetDifference("$flowerFieldB", "$flowerFieldA")),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "flowerFieldA", Value: 1},
			{Key: "flowerFieldB", Value: 1},
			{Key: "inBOnly", Value: bson.D{{Key: "$setDifference", Value: bson.A{"$flowerFieldB", "$flowerFieldA"}}}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSetEquals(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("cakes"),
			agg.Include("cupcakes"),
			agg.Compute("sameFlavors", agg.SetEquals("$cakes", "$cupcakes")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "cakes", Value: 1},
			{Key: "cupcakes", Value: 1},
			{Key: "sameFlavors", Value: bson.D{{Key: "$setEquals", Value: bson.A{"$cakes", "$cupcakes"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSetIntersection_ElementsArrayExample(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("flowerFieldA"),
			agg.Include("flowerFieldB"),
			agg.Compute("commonToBoth", agg.SetIntersection("$flowerFieldA", "$flowerFieldB")),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "flowerFieldA", Value: 1},
			{Key: "flowerFieldB", Value: 1},
			{Key: "commonToBoth", Value: bson.D{{Key: "$setIntersection", Value: bson.A{"$flowerFieldA", "$flowerFieldB"}}}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestSetIntersection_RetrieveDocsForRolesGrantedToCurrentUser when $not and $expr are implemented

func TestSetIsSubset(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("flowerFieldA"),
			agg.Include("flowerFieldB"),
			agg.Compute("AisSubset", agg.SetIsSubset("$flowerFieldA", "$flowerFieldB")),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "flowerFieldA", Value: 1},
			{Key: "flowerFieldB", Value: 1},
			{Key: "AisSubset", Value: bson.D{{Key: "$setIsSubset", Value: bson.A{"$flowerFieldA", "$flowerFieldB"}}}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSetUnion(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("flowerFieldA"),
			agg.Include("flowerFieldB"),
			agg.Compute("allValues", agg.SetUnion("$flowerFieldA", "$flowerFieldB")),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "flowerFieldA", Value: 1},
			{Key: "flowerFieldB", Value: 1},
			{Key: "allValues", Value: bson.D{{Key: "$setUnion", Value: bson.A{"$flowerFieldA", "$flowerFieldB"}}}},
			{Key: "_id", Value: 0},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSigmoid(t *testing.T) {
	got := agg.Pipeline{
		agg.SetStage(agg.Assign("scaled", agg.Sigmoid("$score"))),
	}
	want := bson.A{
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "scaled", Value: bson.D{
				{Key: "$sigmoid", Value: "$score"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSin(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("side_b", agg.Multiply(agg.Sin(agg.DegreesToRadians("$angle_a")), "$hypotenuse")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "side_b", Value: bson.D{
				{Key: "$multiply", Value: bson.A{
					bson.D{{Key: "$sin", Value: bson.D{
						{Key: "$degreesToRadians", Value: "$angle_a"},
					}}},
					"$hypotenuse",
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSinh(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("sinh_output", agg.Sinh(agg.DegreesToRadians("$angle"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "sinh_output", Value: bson.D{
				{Key: "$sinh", Value: bson.D{
					{Key: "$degreesToRadians", Value: "$angle"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestSize when $cond is implemented          
           
// TODO: implement TestSplit after $unwind stage is implemented

func TestSlice(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("threeFavorites", agg.Slice("$favorites", 3, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "threeFavorites", Value: bson.D{{Key: "$slice", Value: bson.A{"$favorites", 3}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSortArray_SortOnAField(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("result", agg.SortArray("$team", agg.Sort("name", agg.Asc))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "result", Value: bson.D{{Key: "$sortArray", Value: bson.D{
				{Key: "input", Value: "$team"},
				{Key: "sortBy", Value: bson.D{{Key: "name", Value: 1}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSortArray_SortOnASubfield(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("result", agg.SortArray("$team", agg.Sort("address.city", agg.Desc))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "result", Value: bson.D{{Key: "$sortArray", Value: bson.D{
				{Key: "input", Value: "$team"},
				{Key: "sortBy", Value: bson.D{{Key: "address.city", Value: -1}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSortArray_SortOnMultipleFields(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("result", agg.SortArray("$team", agg.Sort("age", agg.Desc), agg.Sort("name", agg.Asc))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "result", Value: bson.D{{Key: "$sortArray", Value: bson.D{
				{Key: "input", Value: "$team"},
				{Key: "sortBy", Value: bson.D{
					{Key: "age", Value: -1},
					{Key: "name", Value: 1},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSortArrayByValue_SortArrayOfIntegers(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("result", agg.SortArrayByValue([]int{1, 4, 1, 6, 12, 5}, agg.Asc)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "result", Value: bson.D{{Key: "$sortArray", Value: bson.D{
				{Key: "input", Value: []int{1, 4, 1, 6, 12, 5}},
				{Key: "sortBy", Value: 1},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSortArrayByValue_SortOnMixedTypeFields(t *testing.T) {
	d, err := bson.ParseDecimal128("10.23")
	if err != nil {
		t.Fatal(err)
	}
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("result", agg.SortArrayByValue([]any{
				20, 4, bson.D{{Key: "a", Value: "Free"}}, 6, 21, 5, "Gratis",
				bson.D{{Key: "a", Value: nil}},
				bson.D{{Key: "a", Value: bson.D{{Key: "sale", Value: true}, {Key: "price", Value: 19}}}},
				d,
				bson.D{{Key: "a", Value: "On sale"}},
			}, agg.Asc)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "result", Value: bson.D{{Key: "$sortArray", Value: bson.D{
				{Key: "input", Value: []any{
					20, 4, bson.D{{Key: "a", Value: "Free"}}, 6, 21, 5, "Gratis",
					bson.D{{Key: "a", Value: nil}},
					bson.D{{Key: "a", Value: bson.D{{Key: "sale", Value: true}, {Key: "price", Value: 19}}}},
					d,
					bson.D{{Key: "a", Value: "On sale"}},
				}},
				{Key: "sortBy", Value: 1},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSqrt(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("distance", agg.Sqrt(agg.Add(
				agg.Pow(agg.Subtract("$p2.y", "$p1.y"), 2),
				agg.Pow(agg.Subtract("$p2.x", "$p1.x"), 2),
			))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "distance", Value: bson.D{
				{Key: "$sqrt", Value: bson.D{
					{Key: "$add", Value: bson.A{
						bson.D{{Key: "$pow", Value: bson.A{
							bson.D{{Key: "$subtract", Value: bson.A{"$p2.y", "$p1.y"}}},
							2,
						}}},
						bson.D{{Key: "$pow", Value: bson.A{
							bson.D{{Key: "$subtract", Value: bson.A{"$p2.x", "$p1.x"}}},
							2,
						}}},
					}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestStdDevPop(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("stdDev", agg.StdDevPop("$scores.score")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "stdDev", Value: bson.D{
				{Key: "$stdDevPop", Value: bson.A{"$scores.score"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestStrcasecmp(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("comparisonResult", agg.Strcasecmp("$quarter", "13q4")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "comparisonResult", Value: bson.D{{Key: "$strcasecmp", Value: bson.A{"$quarter", "13q4"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestStrLenBytes(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("length", agg.StrLenBytes("$name")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "length", Value: bson.D{{Key: "$strLenBytes", Value: "$name"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestStrLenCP(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("length", agg.StrLenCP("$name")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "length", Value: bson.D{{Key: "$strLenCP", Value: "$name"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSubstr(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("yearSubstring", agg.Substr("$quarter", 0, 2)),
			agg.Compute("quarterSubstring", agg.Substr("$quarter", 2, -1)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "yearSubstring", Value: bson.D{{Key: "$substr", Value: bson.A{"$quarter", 0, 2}}}},
			{Key: "quarterSubstring", Value: bson.D{{Key: "$substr", Value: bson.A{"$quarter", 2, -1}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSubstrBytes_SingleByteCharacterSet(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("yearSubstring", agg.SubstrBytes("$quarter", 0, 2)),
			agg.Compute("quarterSubstring", agg.SubstrBytes("$quarter", 2, agg.Subtract(agg.StrLenBytes("$quarter"), 2))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "yearSubstring", Value: bson.D{{Key: "$substrBytes", Value: bson.A{"$quarter", 0, 2}}}},
			{Key: "quarterSubstring", Value: bson.D{{Key: "$substrBytes", Value: bson.A{
				"$quarter",
				2,
				bson.D{{Key: "$subtract", Value: bson.A{
					bson.D{{Key: "$strLenBytes", Value: "$quarter"}},
					2,
				}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSubstrBytes_SingleByteAndMultibyteCharacterSet(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("menuCode", agg.SubstrBytes("$name", 0, 3)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "menuCode", Value: bson.D{{Key: "$substrBytes", Value: bson.A{"$name", 0, 3}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSubstrCP_SingleByteCharacterSet(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("yearSubstring", agg.SubstrCP("$quarter", 0, 2)),
			agg.Compute("quarterSubstring", agg.SubstrCP("$quarter", 2, agg.Subtract(agg.StrLenCP("$quarter"), 2))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "yearSubstring", Value: bson.D{{Key: "$substrCP", Value: bson.A{"$quarter", 0, 2}}}},
			{Key: "quarterSubstring", Value: bson.D{{Key: "$substrCP", Value: bson.A{
				"$quarter",
				2,
				bson.D{{Key: "$subtract", Value: bson.A{
					bson.D{{Key: "$strLenCP", Value: "$quarter"}},
					2,
				}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSubstrCP_SingleByteAndMultibyteCharacterSet(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("menuCode", agg.SubstrCP("$name", 0, 3)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "menuCode", Value: bson.D{{Key: "$substrCP", Value: bson.A{"$name", 0, 3}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSubtract_SubtractNumbers(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("total", agg.Subtract(agg.Add("$price", "$fee"), "$discount")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "total", Value: bson.D{
				{Key: "$subtract", Value: bson.A{
					bson.D{{Key: "$add", Value: bson.A{"$price", "$fee"}}},
					"$discount",
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestSubtract_SubtractTwoDates when date functionality is implemented

// TODO: implement TestSubtract_SubtractMillisecondsFromDate when date functionality is implemented

func TestSum(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("quizTotal", agg.Sum("$quizzes")),
			agg.Compute("labTotal", agg.Sum("$labs")),
			agg.Compute("examTotal", agg.Sum("$final", "$midterm")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "quizTotal", Value: bson.D{{Key: "$sum", Value: bson.A{"$quizzes"}}}},
			{Key: "labTotal", Value: bson.D{{Key: "$sum", Value: bson.A{"$labs"}}}},
			{Key: "examTotal", Value: bson.D{{Key: "$sum", Value: bson.A{"$final", "$midterm"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestTan(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("side_b", agg.Multiply(agg.Tan(agg.DegreesToRadians("$angle_a")), "$side_a")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "side_b", Value: bson.D{
				{Key: "$multiply", Value: bson.A{
					bson.D{{Key: "$tan", Value: bson.D{
						{Key: "$degreesToRadians", Value: "$angle_a"},
					}}},
					"$side_a",
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestTanh(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("tanh_output", agg.Tanh(agg.DegreesToRadians("$angle"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "tanh_output", Value: bson.D{
				{Key: "$tanh", Value: bson.D{
					{Key: "$degreesToRadians", Value: "$angle"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestToLower(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("item", agg.ToLower("$item")),
			agg.Compute("description", agg.ToLower("$description")),
    ),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
      {Key: "item", Value: bson.D{{Key: "$toLower", Value: "$item"}}},
			{Key: "description", Value: bson.D{{Key: "$toLower", Value: "$description"}}},
		}}},
	}
  assertPipelineEqual(t, got, want)
}
           
func TestToUpper(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("item", agg.ToUpper("$item")),
			agg.Compute("description", agg.ToUpper("$description")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: bson.D{{Key: "$toUpper", Value: "$item"}}},
			{Key: "description", Value: bson.D{{Key: "$toUpper", Value: "$description"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}
           
func TestTop(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("topScore", agg.Top("$results", []any{"$playerId", "$score"}, agg.Sort("score", agg.Desc))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
      			{Key: "topScore", Value: bson.D{{Key: "$top", Value: bson.D{
				{Key: "sortBy", Value: bson.D{{Key: "score", Value: -1}}},
				{Key: "output", Value: bson.A{"$playerId", "$score"}},
				{Key: "input", Value: "$results"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestTopN(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("topScores", agg.TopN(3, "$results", []any{"$playerId", "$score"}, agg.Sort("score", agg.Desc))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "topScores", Value: bson.D{{Key: "$topN", Value: bson.D{
				{Key: "n", Value: 3},
				{Key: "sortBy", Value: bson.D{{Key: "score", Value: -1}}},
				{Key: "output", Value: bson.A{"$playerId", "$score"}},
				{Key: "input", Value: "$results"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestTrim(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("description", agg.Trim("$description", nil)),
    ),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
      {Key: "item", Value: 1},
			{Key: "description", Value: bson.D{{Key: "$trim", Value: bson.D{
			{Key: "input", Value: "$description"},
    }}}},
    }}},
	}
	assertPipelineEqual(t, got, want)
}

func TestTrunc(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("truncatedValue", agg.Trunc("$value", 1)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "truncatedValue", Value: bson.D{
				{Key: "$trunc", Value: bson.A{"$value", 1}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestZip_MatrixTransposition(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("transposed", agg.Zip([]any{
				agg.ArrayElemAt("$matrix", 0),
				agg.ArrayElemAt("$matrix", 1),
				agg.ArrayElemAt("$matrix", 2),
			}, false)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "transposed", Value: bson.D{{Key: "$zip", Value: bson.D{
				{Key: "inputs", Value: bson.A{
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$matrix", 0}}},
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$matrix", 1}}},
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$matrix", 2}}},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestZip_FilteringAndPreservingIndexes when $let is implemented
