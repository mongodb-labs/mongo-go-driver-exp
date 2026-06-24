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

// TODO: implement TestAcos when $addFields is implemented

// TODO: implement TestAcosh when $addFields is implemented

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

// TODO: implement TestAsin when $addFields is implemented

// TODO: implement TestAsinh when $addFields is implemented

// TODO: implement TestAtan when $addFields is implemented

// TODO: implement TestAtan2 when $addFields is implemented

// TODO: implement TestAtanh when $addFields is implemented

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

// TODO: implement TestCos when $addFields is implemented

// TODO: implement TestCosh when $addFields is implemented

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

// TODO: implement TestDegreesToRadians when $addFields is implemented

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

// TODO: implement TestRadiansToDegrees when $addFields is implemented

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

// TODO: implement TestSin when $addFields is implemented

// TODO: implement TestSinh when $addFields is implemented

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

// TODO: implement TestTan when $addFields is implemented

// TODO: implement TestTanh when $addFields is implemented

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
