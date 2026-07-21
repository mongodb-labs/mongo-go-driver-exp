package agg_test

import (
	"testing"
	"time"

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

func TestBinarySize(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("name", "$name"),
			agg.Compute("imageSize", agg.BinarySize("$binary")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: "$name"},
			{Key: "imageSize", Value: bson.D{{Key: "$binarySize", Value: "$binary"}}},
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

func TestBsonSize_ReturnSizesOfDocs(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("object_size", agg.BsonSize("$$ROOT")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "object_size", Value: bson.D{{Key: "$bsonSize", Value: "$$ROOT"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBsonSize_ReturnCombinedSizeOfAllDocsInCollection(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage(
			agg.Null,
			agg.Accumulate("combined_object_size", agg.SumAccumulator(agg.BsonSize("$$ROOT"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "combined_object_size", Value: bson.D{
				{Key: "$sum", Value: bson.D{{Key: "$bsonSize", Value: "$$ROOT"}}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestBsonSize_ReturnDocWithLargestSpecifiedField(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("name", "$name"),
			agg.Compute("task_object_size", agg.BsonSize("$current_task")),
		),
		agg.SortStage(agg.Sort("task_object_size", agg.Desc)),
		agg.LimitStage(1),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: "$name"},
			{Key: "task_object_size", Value: bson.D{{Key: "$bsonSize", Value: "$current_task"}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "task_object_size", Value: -1},
		}}},
		bson.D{{Key: "$limit", Value: 1}},
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

func TestCond(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("discount", agg.Cond(agg.Gte("$qty", 250), 30, 20)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "discount", Value: bson.D{
				{Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{
						{Key: "$gte", Value: bson.A{"$qty", 250}},
					}},
					{Key: "then", Value: 30},
					{Key: "else", Value: 20},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestConvert_Example when $switch is implemented

func TestConvert_ConvertHexadecimalStringToInteger(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("decimalValue", agg.Convert("$hexString", "int", agg.WithConvertBase(16))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "decimalValue", Value: bson.D{{Key: "$convert", Value: bson.D{
				{Key: "input", Value: "$hexString"},
				{Key: "to", Value: "int"},
				{Key: "base", Value: int32(16)},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestConvert_ConvertIntegerToBinaryString(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("binaryString", agg.Convert("$value", "string", agg.WithConvertBase(2))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "binaryString", Value: bson.D{{Key: "$convert", Value: bson.D{
				{Key: "input", Value: "$value"},
				{Key: "to", Value: "string"},
				{Key: "base", Value: int32(2)},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestConvert_OnErrorReturnsFallbackString(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("convertedQty", agg.Convert("$qty", "int",
				agg.WithConvertOnError(agg.ToString("$qty")))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "convertedQty", Value: bson.D{{Key: "$convert", Value: bson.D{
				{Key: "input", Value: "$qty"},
				{Key: "to", Value: "int"},
				{Key: "onError", Value: bson.D{{Key: "$toString", Value: "$qty"}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestConvert_OnNullReturnsFallbackDecimal(t *testing.T) {
	onNull, err := bson.ParseDecimal128("0")
	if err != nil {
		t.Fatal(err)
	}
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("convertedPrice", agg.Convert("$price", "decimal", agg.WithConvertOnError("Error"), agg.WithConvertOnNull(onNull))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "convertedPrice", Value: bson.D{{Key: "$convert", Value: bson.D{
				{Key: "input", Value: "$price"},
				{Key: "to", Value: "decimal"},
				{Key: "onError", Value: "Error"},
				{Key: "onNull", Value: onNull},
			}}}},
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

func TestCreateObjectId(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("objectId", agg.CreateObjectId()),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "objectId", Value: bson.D{
				{Key: "$createObjectId", Value: bson.D{}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: add $merge stage when implemented
func TestDateAdd_AddFutureDate(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("expectedDeliveryDate", agg.DateAdd("$purchaseDate", "day", 3)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "expectedDeliveryDate", Value: bson.D{
				{Key: "$dateAdd", Value: bson.D{
					{Key: "startDate", Value: "$purchaseDate"},
					{Key: "unit", Value: "day"},
					{Key: "amount", Value: 3},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestDateAdd_FilterOnDateRange when $expr is implemented

func TestDateAdd_AdjustForDaylightSavingsTime(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("location"),
			agg.Compute("start", agg.DateToString("$login", agg.WithDateToStringFormat("%Y-%m-%d %H:%M"))),
			agg.Compute("days", agg.DateToString(
				agg.DateAdd("$login", "day", 1, agg.WithDateAddTimezone("$location")),
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"))),
			agg.Compute("hours", agg.DateToString(
				agg.DateAdd("$login", "hour", 24, agg.WithDateAddTimezone("$location")),
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"))),
			agg.Compute("startTZInfo", agg.DateToString("$login",
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"),
				agg.WithDateToStringTimezone("$location"))),
			agg.Compute("daysTZInfo", agg.DateToString(
				agg.DateAdd("$login", "day", 1, agg.WithDateAddTimezone("$location")),
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"),
				agg.WithDateToStringTimezone("$location"))),
			agg.Compute("hoursTZInfo", agg.DateToString(
				agg.DateAdd("$login", "hour", 24, agg.WithDateAddTimezone("$location")),
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"),
				agg.WithDateToStringTimezone("$location"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: int32(0)},
			{Key: "location", Value: int32(1)},
			{Key: "start", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$login"},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
				}},
			}},
			{Key: "days", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "$dateAdd", Value: bson.D{
							{Key: "startDate", Value: "$login"},
							{Key: "unit", Value: "day"},
							{Key: "amount", Value: 1},
							{Key: "timezone", Value: "$location"},
						}},
					}},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
				}},
			}},
			{Key: "hours", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "$dateAdd", Value: bson.D{
							{Key: "startDate", Value: "$login"},
							{Key: "unit", Value: "hour"},
							{Key: "amount", Value: 24},
							{Key: "timezone", Value: "$location"},
						}},
					}},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
				}},
			}},
			{Key: "startTZInfo", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$login"},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
					{Key: "timezone", Value: "$location"},
				}},
			}},
			{Key: "daysTZInfo", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "$dateAdd", Value: bson.D{
							{Key: "startDate", Value: "$login"},
							{Key: "unit", Value: "day"},
							{Key: "amount", Value: 1},
							{Key: "timezone", Value: "$location"},
						}},
					}},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
					{Key: "timezone", Value: "$location"},
				}},
			}},
			{Key: "hoursTZInfo", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "$dateAdd", Value: bson.D{
							{Key: "startDate", Value: "$login"},
							{Key: "unit", Value: "hour"},
							{Key: "amount", Value: 24},
							{Key: "timezone", Value: "$location"},
						}},
					}},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
					{Key: "timezone", Value: "$location"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateDiff_ElapsedTime(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage(
			agg.Null,
			agg.Accumulate("averageTime", agg.AvgAccumulator(agg.DateDiff("$purchased", "$delivered", "day"))),
		),
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("numDays", agg.Trunc("$averageTime", agg.WithTruncPlace(1))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "averageTime", Value: bson.D{
				{Key: "$avg", Value: bson.D{
					{Key: "$dateDiff", Value: bson.D{
						{Key: "startDate", Value: "$purchased"},
						{Key: "endDate", Value: "$delivered"},
						{Key: "unit", Value: "day"},
					}},
				}},
			}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: int32(0)},
			{Key: "numDays", Value: bson.D{
				{Key: "$trunc", Value: bson.A{"$averageTime", 1}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateDiff_ResultPrecision(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("Start", "$start"),
			agg.Compute("End", "$end"),
			agg.Compute("years", agg.DateDiff("$start", "$end", "year")),
			agg.Compute("months", agg.DateDiff("$start", "$end", "month")),
			agg.Compute("days", agg.DateDiff("$start", "$end", "day")),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "Start", Value: "$start"},
			{Key: "End", Value: "$end"},
			{Key: "years", Value: bson.D{
				{Key: "$dateDiff", Value: bson.D{
					{Key: "startDate", Value: "$start"},
					{Key: "endDate", Value: "$end"},
					{Key: "unit", Value: "year"},
				}},
			}},
			{Key: "months", Value: bson.D{
				{Key: "$dateDiff", Value: bson.D{
					{Key: "startDate", Value: "$start"},
					{Key: "endDate", Value: "$end"},
					{Key: "unit", Value: "month"},
				}},
			}},
			{Key: "days", Value: bson.D{
				{Key: "$dateDiff", Value: bson.D{
					{Key: "startDate", Value: "$start"},
					{Key: "endDate", Value: "$end"},
					{Key: "unit", Value: "day"},
				}},
			}},
			{Key: "_id", Value: int32(0)},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateDiff_WeeksPerMonth(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("wks_default", agg.DateDiff("$start", "$end", "week")),
			agg.Compute("wks_monday", agg.DateDiff("$start", "$end", "week", agg.WithDateDiffStartOfWeek("Monday"))),
			agg.Compute("wks_friday", agg.DateDiff("$start", "$end", "week", agg.WithDateDiffStartOfWeek("fri"))),
			agg.Exclude("_id"),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "wks_default", Value: bson.D{
				{Key: "$dateDiff", Value: bson.D{
					{Key: "startDate", Value: "$start"},
					{Key: "endDate", Value: "$end"},
					{Key: "unit", Value: "week"},
				}},
			}},
			{Key: "wks_monday", Value: bson.D{
				{Key: "$dateDiff", Value: bson.D{
					{Key: "startDate", Value: "$start"},
					{Key: "endDate", Value: "$end"},
					{Key: "unit", Value: "week"},
					{Key: "startOfWeek", Value: "Monday"},
				}},
			}},
			{Key: "wks_friday", Value: bson.D{
				{Key: "$dateDiff", Value: bson.D{
					{Key: "startDate", Value: "$start"},
					{Key: "endDate", Value: "$end"},
					{Key: "unit", Value: "week"},
					{Key: "startOfWeek", Value: "fri"},
				}},
			}},
			{Key: "_id", Value: int32(0)},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateFromParts(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("date", agg.DateFromParts(
				agg.WithDateFromPartsYear(2017),
				agg.WithDateFromPartsMonth(2),
				agg.WithDateFromPartsDay(8),
				agg.WithDateFromPartsHour(12),
			)),
			agg.Compute("date_iso", agg.DateFromParts(
				agg.WithDateFromPartsIsoWeekYear(2017),
				agg.WithDateFromPartsIsoWeek(6),
				agg.WithDateFromPartsIsoDayOfWeek(3),
				agg.WithDateFromPartsHour(12),
			)),
			agg.Compute("date_timezone", agg.DateFromParts(
				agg.WithDateFromPartsYear(2016),
				agg.WithDateFromPartsMonth(12),
				agg.WithDateFromPartsDay(31),
				agg.WithDateFromPartsHour(23),
				agg.WithDateFromPartsMinute(46),
				agg.WithDateFromPartsSecond(12),
				agg.WithDateFromPartsTimezone("America/New_York"),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "date", Value: bson.D{
				{Key: "$dateFromParts", Value: bson.D{
					{Key: "year", Value: 2017},
					{Key: "month", Value: 2},
					{Key: "day", Value: 8},
					{Key: "hour", Value: 12},
				}},
			}},
			{Key: "date_iso", Value: bson.D{
				{Key: "$dateFromParts", Value: bson.D{
					{Key: "isoWeekYear", Value: 2017},
					{Key: "isoWeek", Value: 6},
					{Key: "isoDayOfWeek", Value: 3},
					{Key: "hour", Value: 12},
				}},
			}},
			{Key: "date_timezone", Value: bson.D{
				{Key: "$dateFromParts", Value: bson.D{
					{Key: "year", Value: 2016},
					{Key: "month", Value: 12},
					{Key: "day", Value: 31},
					{Key: "hour", Value: 23},
					{Key: "minute", Value: 46},
					{Key: "second", Value: 12},
					{Key: "timezone", Value: "America/New_York"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateFromString_ConvertingDates(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("date", agg.DateFromString("$date", agg.WithDateFromStringTimezone("America/New_York"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "date", Value: bson.D{
				{Key: "$dateFromString", Value: bson.D{
					{Key: "dateString", Value: "$date"},
					{Key: "timezone", Value: "America/New_York"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateFromString_OnError(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute(
				"date",
				agg.DateFromString(
					"$date",
					agg.WithDateFromStringTimezone("$timezone"),
					agg.WithDateFromStringOnError("$date"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "date", Value: bson.D{
				{Key: "$dateFromString", Value: bson.D{
					{Key: "dateString", Value: "$date"},
					{Key: "timezone", Value: "$timezone"},
					{Key: "onError", Value: "$date"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateFromString_OnNull(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute(
				"date",
				agg.DateFromString(
					"$date",
					agg.WithDateFromStringTimezone("$timezone"),
					agg.WithDateFromStringOnNull(time.UnixMilli(0).UTC()))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "date", Value: bson.D{
				{Key: "$dateFromString", Value: bson.D{
					{Key: "dateString", Value: "$date"},
					{Key: "timezone", Value: "$timezone"},
					{Key: "onNull", Value: bson.DateTime(0)},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestDateSubtract_SubtractFixedAmount and TestDateSubtract_FilterByRelativeDates
// when $expr is implemented

func TestDateSubtract_AdjustForDaylightSavingsTime(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("location"),
			agg.Compute("start", agg.DateToString("$login", agg.WithDateToStringFormat("%Y-%m-%d %H:%M"))),
			agg.Compute("days", agg.DateToString(
				agg.DateSubtract("$login", "day", 1, agg.WithDateSubtractTimezone("$location")),
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"))),
			agg.Compute("hours", agg.DateToString(
				agg.DateSubtract("$login", "hour", 24, agg.WithDateSubtractTimezone("$location")),
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"))),
			agg.Compute("startTZInfo", agg.DateToString("$login",
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"),
				agg.WithDateToStringTimezone("$location"))),
			agg.Compute("daysTZInfo", agg.DateToString(
				agg.DateSubtract("$login", "day", 1, agg.WithDateSubtractTimezone("$location")),
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"),
				agg.WithDateToStringTimezone("$location"))),
			agg.Compute("hoursTZInfo", agg.DateToString(
				agg.DateSubtract("$login", "hour", 24, agg.WithDateSubtractTimezone("$location")),
				agg.WithDateToStringFormat("%Y-%m-%d %H:%M"),
				agg.WithDateToStringTimezone("$location"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: int32(0)},
			{Key: "location", Value: int32(1)},
			{Key: "start", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$login"},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
				}},
			}},
			{Key: "days", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "$dateSubtract", Value: bson.D{
							{Key: "startDate", Value: "$login"},
							{Key: "unit", Value: "day"},
							{Key: "amount", Value: 1},
							{Key: "timezone", Value: "$location"},
						}},
					}},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
				}},
			}},
			{Key: "hours", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "$dateSubtract", Value: bson.D{
							{Key: "startDate", Value: "$login"},
							{Key: "unit", Value: "hour"},
							{Key: "amount", Value: 24},
							{Key: "timezone", Value: "$location"},
						}},
					}},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
				}},
			}},
			{Key: "startTZInfo", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$login"},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
					{Key: "timezone", Value: "$location"},
				}},
			}},
			{Key: "daysTZInfo", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "$dateSubtract", Value: bson.D{
							{Key: "startDate", Value: "$login"},
							{Key: "unit", Value: "day"},
							{Key: "amount", Value: 1},
							{Key: "timezone", Value: "$location"},
						}},
					}},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
					{Key: "timezone", Value: "$location"},
				}},
			}},
			{Key: "hoursTZInfo", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: bson.D{
						{Key: "$dateSubtract", Value: bson.D{
							{Key: "startDate", Value: "$login"},
							{Key: "unit", Value: "hour"},
							{Key: "amount", Value: 24},
							{Key: "timezone", Value: "$location"},
						}},
					}},
					{Key: "format", Value: "%Y-%m-%d %H:%M"},
					{Key: "timezone", Value: "$location"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateToParts(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("date", agg.DateToParts("$date")),
			agg.Compute("date_iso", agg.DateToParts("$date", agg.WithDateToPartsIso8601(true))),
			agg.Compute("date_timezone", agg.DateToParts("$date", agg.WithDateToPartsTimezone("America/New_York"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "date", Value: bson.D{
				{Key: "$dateToParts", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
			{Key: "date_iso", Value: bson.D{
				{Key: "$dateToParts", Value: bson.D{
					{Key: "date", Value: "$date"},
					{Key: "iso8601", Value: true},
				}},
			}},
			{Key: "date_timezone", Value: bson.D{
				{Key: "$dateToParts", Value: bson.D{
					{Key: "date", Value: "$date"},
					{Key: "timezone", Value: "America/New_York"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateToString(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("yearMonthDayUTC", agg.DateToString("$date",
				agg.WithDateToStringFormat("%Y-%m-%d"))),
			agg.Compute("timewithOffsetNY", agg.DateToString("$date",
				agg.WithDateToStringFormat("%H:%M:%S:%L%z"),
				agg.WithDateToStringTimezone("America/New_York"))),
			agg.Compute("timewithOffset430", agg.DateToString("$date",
				agg.WithDateToStringFormat("%H:%M:%S:%L%z"),
				agg.WithDateToStringTimezone("+04:30"))),
			agg.Compute("minutesOffsetNY", agg.DateToString("$date",
				agg.WithDateToStringFormat("%Z"),
				agg.WithDateToStringTimezone("America/New_York"))),
			agg.Compute("minutesOffset430", agg.DateToString("$date",
				agg.WithDateToStringFormat("%Z"),
				agg.WithDateToStringTimezone("+04:30"))),
			agg.Compute("abbreviated_month", agg.DateToString("$date",
				agg.WithDateToStringFormat("%b"),
				agg.WithDateToStringTimezone("+04:30"))),
			agg.Compute("full_month", agg.DateToString("$date",
				agg.WithDateToStringFormat("%B"),
				agg.WithDateToStringTimezone("+04:30"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "yearMonthDayUTC", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$date"},
					{Key: "format", Value: "%Y-%m-%d"},
				}},
			}},
			{Key: "timewithOffsetNY", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$date"},
					{Key: "format", Value: "%H:%M:%S:%L%z"},
					{Key: "timezone", Value: "America/New_York"},
				}},
			}},
			{Key: "timewithOffset430", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$date"},
					{Key: "format", Value: "%H:%M:%S:%L%z"},
					{Key: "timezone", Value: "+04:30"},
				}},
			}},
			{Key: "minutesOffsetNY", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$date"},
					{Key: "format", Value: "%Z"},
					{Key: "timezone", Value: "America/New_York"},
				}},
			}},
			{Key: "minutesOffset430", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$date"},
					{Key: "format", Value: "%Z"},
					{Key: "timezone", Value: "+04:30"},
				}},
			}},
			{Key: "abbreviated_month", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$date"},
					{Key: "format", Value: "%b"},
					{Key: "timezone", Value: "+04:30"},
				}},
			}},
			{Key: "full_month", Value: bson.D{
				{Key: "$dateToString", Value: bson.D{
					{Key: "date", Value: "$date"},
					{Key: "format", Value: "%B"},
					{Key: "timezone", Value: "+04:30"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateTrunc_ProjectStage(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("_id"),
			agg.Include("orderDate"),
			agg.Compute(
				"truncatedOrderDate",
				agg.DateTrunc(
					"$orderDate",
					"week",
					agg.WithDateTruncBinSize(2),
					agg.WithDateTruncTimezone("America/Los_Angeles"),
					agg.WithDateTruncStartOfWeek("Monday"),
				),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: int32(1)},
			{Key: "orderDate", Value: int32(1)},
			{Key: "truncatedOrderDate", Value: bson.D{
				{Key: "$dateTrunc", Value: bson.D{
					{Key: "date", Value: "$orderDate"},
					{Key: "unit", Value: "week"},
					{Key: "binSize", Value: 2},
					{Key: "timezone", Value: "America/Los_Angeles"},
					{Key: "startOfWeek", Value: "Monday"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDateTrunc_GroupStage(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage(
			bson.D{{Key: "truncatedOrderDate", Value: agg.DateTrunc("$orderDate", "month",
				agg.WithDateTruncBinSize(6))}},
			agg.Accumulate("sumQuantity", agg.SumAccumulator("$quantity")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "truncatedOrderDate", Value: bson.D{
					{Key: "$dateTrunc", Value: bson.D{
						{Key: "date", Value: "$orderDate"},
						{Key: "unit", Value: "month"},
						{Key: "binSize", Value: 6},
					}},
				}},
			}},
			{Key: "sumQuantity", Value: bson.D{
				{Key: "$sum", Value: "$quantity"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDayOfMonth(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("day", agg.DayOfMonth("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "day", Value: bson.D{
				{Key: "$dayOfMonth", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDayOfWeek(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("dayOfWeek", agg.DayOfWeek("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "dayOfWeek", Value: bson.D{
				{Key: "$dayOfWeek", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDayOfYear(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("dayOfYear", agg.DayOfYear("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "dayOfYear", Value: bson.D{
				{Key: "$dayOfYear", Value: bson.D{
					{Key: "date", Value: "$date"},
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

func TestDeserializeEJSON_DeserializeExtendedJSONDocument(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("Inception"))),
		agg.ProjectStage(
			agg.Compute("original", agg.RootObject()),
			agg.Compute("serialized", agg.SerializeEJSON(agg.RootObject())),
		),
		agg.ProjectStage(
			agg.Compute("title", "$original.title"),
			agg.Compute("deserialized", agg.DeserializeEJSON("$serialized")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "title", Value: bson.D{{Key: "$eq", Value: "Inception"}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "original", Value: "$$ROOT"},
			{Key: "serialized", Value: bson.D{{Key: "$serializeEJSON", Value: bson.D{
				{Key: "input", Value: "$$ROOT"},
			}}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "title", Value: "$original.title"},
			{Key: "deserialized", Value: bson.D{{Key: "$deserializeEJSON", Value: bson.D{
				{Key: "input", Value: "$serialized"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestDeserializeEJSON_ParseJSONStringAndDeserialize when the $documents stage and $convert operator are implemented

func TestDeserializeEJSON_DeserializeSpecificFields(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("Inception"))),
		agg.ProjectStage(
			agg.Include("title"),
			agg.Compute("serializedMetadata", agg.SerializeEJSON(bson.D{
				{Key: "releaseDate", Value: "$released"},
				{Key: "runtime", Value: "$runtime"},
				{Key: "rating", Value: "$imdb.rating"},
			})),
		),
		agg.ProjectStage(
			agg.Include("title"),
			agg.Compute("metadata", agg.DeserializeEJSON("$serializedMetadata")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "title", Value: bson.D{{Key: "$eq", Value: "Inception"}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "serializedMetadata", Value: bson.D{{Key: "$serializeEJSON", Value: bson.D{
				{Key: "input", Value: bson.D{
					{Key: "releaseDate", Value: "$released"},
					{Key: "runtime", Value: "$runtime"},
					{Key: "rating", Value: "$imdb.rating"},
				}},
			}}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "metadata", Value: bson.D{{Key: "$deserializeEJSON", Value: bson.D{
				{Key: "input", Value: "$serializedMetadata"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestDeserializeEJSON_UseOnErrorForErrorHandling(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("result", agg.DeserializeEJSON("$ejsonField",
				agg.WithDeserializeEJSONOnError("Invalid EJSON format"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "result", Value: bson.D{{Key: "$deserializeEJSON", Value: bson.D{
				{Key: "input", Value: "$ejsonField"},
				{Key: "onError", Value: "Invalid EJSON format"},
			}}}},
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

func TestFilter_WithoutFilterOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("items", agg.Filter("$items", agg.Gte("$$this.price", 100))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "items", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$items"},
					{Key: "cond", Value: bson.D{
						{Key: "$gte", Value: bson.A{"$$this.price", 100}},
					}},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFilter_WithFilterLimit(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("items", agg.Filter("$items", agg.Gte("$$this.price", 100), agg.WithFilterLimit(1))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "items", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$items"},
					{Key: "cond", Value: bson.D{
						{Key: "$gte", Value: bson.A{"$$this.price", 100}},
					}},
					{Key: "limit", Value: 1},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFilter_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("items", agg.Filter("$items", agg.Gte("$$item.price", 100), agg.WithFilterAs("item"))),
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

func TestFilter_UsingLimitField(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("items", agg.Filter("$items", agg.Gte("$$item.price", 100), agg.WithFilterAs("item"), agg.WithFilterLimit(1))),
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

func TestFilter_LimitGreaterThanPossibleMatches(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("items", agg.Filter("$items", agg.Gte("$$item.price", 100), agg.WithFilterAs("item"), agg.WithFilterLimit(5))),
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

func TestFunction_UsageExample(t *testing.T) {
	isFoundBody := "function(name) {\n    return hex_md5(name) == \"15b0a220baa16331e8d80e15367677ad\"\n}"
	messageBody := "function(name, scores) {\n    let total = Array.sum(scores);\n    return `Hello ${name}. Your total score is ${total}.`\n}"
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("isFound", agg.Function(isFoundBody, []agg.Expr{"$name"})),
			agg.Assign("message", agg.Function(messageBody, []agg.Expr{"$name", "$scores"})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "isFound", Value: bson.D{{Key: "$function", Value: bson.D{
				{Key: "body", Value: bson.JavaScript(isFoundBody)},
				{Key: "args", Value: bson.A{"$name"}},
				{Key: "lang", Value: "js"},
			}}}},
			{Key: "message", Value: bson.D{{Key: "$function", Value: bson.D{
				{Key: "body", Value: bson.JavaScript(messageBody)},
				{Key: "args", Value: bson.A{"$name", "$scores"}},
				{Key: "lang", Value: "js"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestFunction_NilArgs(t *testing.T) {
	body := "function() { return 42 }"
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("answer", agg.Function(body, nil)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "answer", Value: bson.D{{Key: "$function", Value: bson.D{
				{Key: "body", Value: bson.JavaScript(body)},
				{Key: "args", Value: bson.A{}},
				{Key: "lang", Value: "js"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestFunction_AlternativeToWhere when the $expr query operator is implemented

// TODO: implement $getField tests when $expr is implemented

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

func TestHash_HashAFieldValue(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("filename"),
			agg.Compute("hash", agg.Hash("$filename", "sha256")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "filename", Value: 1},
			{Key: "hash", Value: bson.D{{Key: "$hash", Value: bson.D{
				{Key: "input", Value: "$filename"},
				{Key: "algorithm", Value: "sha256"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestHash_HashALiteralString when the $documents stage is implemented

func TestHash_HashBinData(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("hash", agg.Hash("$data", "sha256")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "hash", Value: bson.D{{Key: "$hash", Value: bson.D{
				{Key: "input", Value: "$data"},
				{Key: "algorithm", Value: "sha256"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestHash_NullOrMissingInput when the $documents stage is implemented

func TestHexHash_HashAFieldValue(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("filename"),
			agg.Compute("hexHash", agg.HexHash("$filename", "sha256")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "filename", Value: 1},
			{Key: "hexHash", Value: bson.D{{Key: "$hexHash", Value: bson.D{
				{Key: "input", Value: "$filename"},
				{Key: "algorithm", Value: "sha256"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestHexHash_NullOrMissingInput when the $documents stage is implemented

func TestHour(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("hour", agg.Hour("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "hour", Value: bson.D{
				{Key: "$hour", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
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

func TestIndexOfArray_WithoutIndexOfOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("index", agg.IndexOfArray("$items", 2)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "index", Value: bson.D{{Key: "$indexOfArray", Value: bson.A{"$items", 2}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfArray_WithIndexOfStart(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("index", agg.IndexOfArray("$items", 2, agg.WithIndexOfStart(1))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "index", Value: bson.D{{Key: "$indexOfArray", Value: bson.A{"$items", 2, 1}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfArray_WithIndexOfEnd(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("index", agg.IndexOfArray("$items", 2, agg.WithIndexOfEnd(4))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "index", Value: bson.D{{Key: "$indexOfArray", Value: bson.A{"$items", 2, 0, 4}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfArray_WithIndexOfOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("index", agg.IndexOfArray("$items", 2, agg.WithIndexOfStart(1), agg.WithIndexOfEnd(4))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "index", Value: bson.D{{Key: "$indexOfArray", Value: bson.A{"$items", 2, 1, 4}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfBytes_WithoutIndexOfOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("byteLocation", agg.IndexOfBytes("$item", "foo")),
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

func TestIndexOfBytes_WithIndexOfStart(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("byteLocation", agg.IndexOfBytes("$item", "foo", agg.WithIndexOfStart(1))),
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

func TestIndexOfBytes_WithIndexOfEnd(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("byteLocation", agg.IndexOfBytes("$item", "foo", agg.WithIndexOfEnd(4))),
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

func TestIndexOfBytes_WithIndexOfOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("byteLocation", agg.IndexOfBytes("$item", "foo", agg.WithIndexOfStart(1), agg.WithIndexOfEnd(4))),
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

func TestIndexOfCP_WithoutIndexOfOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("cpLocation", agg.IndexOfCP("$item", "foo")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "cpLocation", Value: bson.D{{Key: "$indexOfCP", Value: bson.A{"$item", "foo"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfCP_WithIndexOfStart(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("cpLocation", agg.IndexOfCP("$item", "foo", agg.WithIndexOfStart(1))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "cpLocation", Value: bson.D{{Key: "$indexOfCP", Value: bson.A{"$item", "foo", 1}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfCP_WithIndexOfEnd(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("cpLocation", agg.IndexOfCP("$item", "foo", agg.WithIndexOfEnd(4))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "cpLocation", Value: bson.D{{Key: "$indexOfCP", Value: bson.A{"$item", "foo", 0, 4}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIndexOfCP_WithIndexOfOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("cpLocation", agg.IndexOfCP("$item", "foo", agg.WithIndexOfStart(1), agg.WithIndexOfEnd(4))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "cpLocation", Value: bson.D{{Key: "$indexOfCP", Value: bson.A{"$item", "foo", 1, 4}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement $isArray tests when $cond is implemented

func TestIsNumber_CheckIfFieldIsNumeric(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("isNumber", agg.IsNumber("$reading")),
			agg.Assign("hasType", agg.Type("$reading")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "isNumber", Value: bson.D{{Key: "$isNumber", Value: "$reading"}}},
			{Key: "hasType", Value: bson.D{{Key: "$type", Value: "$reading"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestIsNumber_ConditionallyModifyFields when $cond is implemented

func TestIsoDayOfWeek(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("name", "$name"),
			agg.Compute("dayOfWeek", agg.IsoDayOfWeek("$birthday")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "name", Value: "$name"},
			{Key: "dayOfWeek", Value: bson.D{
				{Key: "$isoDayOfWeek", Value: bson.D{
					{Key: "date", Value: "$birthday"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIsoWeek(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("city", "$city"),
			agg.Compute("weekNumber", agg.IsoWeek("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "city", Value: "$city"},
			{Key: "weekNumber", Value: bson.D{
				{Key: "$isoWeek", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestIsoWeekYear(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("yearNumber", agg.IsoWeekYear("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "yearNumber", Value: bson.D{
				{Key: "$isoWeekYear", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
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

func TestLet(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("finalTotal", agg.Let([]agg.SetField{
				agg.Assign("total", agg.Add("$price", "$tax")),
				agg.Assign("discounted", agg.Cond("$applyDiscount", 0.9, 1)),
			},
				agg.Multiply("$$total", "$$discounted"),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "finalTotal", Value: bson.D{{Key: "$let", Value: bson.D{
				{Key: "vars", Value: bson.D{
					{Key: "total", Value: bson.D{{Key: "$add", Value: bson.A{"$price", "$tax"}}}},
					{Key: "discounted", Value: bson.D{{Key: "$cond", Value: bson.D{
						{Key: "if", Value: "$applyDiscount"},
						{Key: "then", Value: 0.9},
						{Key: "else", Value: 1},
					}}}},
				}},
				{Key: "in", Value: bson.D{{Key: "$multiply", Value: bson.A{"$$total", "$$discounted"}}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestLiteral(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("costsOneDollar", agg.Eq("$price", agg.Literal("$1"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "costsOneDollar", Value: bson.D{
				{Key: "$eq", Value: bson.A{
					"$price",
					bson.D{{Key: "$literal", Value: "$1"}},
				}},
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

func TestLtrim_WithoutTrimChars(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("description", agg.Ltrim("$description")),
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

func TestLtrim_WithTrimChars(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("description", agg.Ltrim("$description", agg.WithTrimChars("*"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "description", Value: bson.D{{Key: "$ltrim", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "chars", Value: "*"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMap_AddToEachElementOfAnArray(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("adjustedGrades", agg.Map("$quizzes", agg.Add("$$grade", 2), agg.WithMapAs("grade"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "adjustedGrades", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: "$quizzes"},
				{Key: "as", Value: "grade"},
				{Key: "in", Value: bson.D{{Key: "$add", Value: bson.A{"$$grade", 2}}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMap_TruncateEachArrayElement(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("city", "$city"),
			agg.Compute("integerValues", agg.Map("$distances", agg.Trunc("$$decimalValue"), agg.WithMapAs("decimalValue"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "city", Value: "$city"},
			{Key: "integerValues", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: "$distances"},
				{Key: "as", Value: "decimalValue"},
				{Key: "in", Value: bson.D{{Key: "$trunc", Value: bson.A{"$$decimalValue"}}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMap_ConvertCelsiusTemperaturesToFahrenheit(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("tempsF", agg.Map("$tempsC", agg.Add(agg.Multiply("$$tempInCelsius", 1.8), 32), agg.WithMapAs("tempInCelsius"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "tempsF", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: "$tempsC"},
				{Key: "as", Value: "tempInCelsius"},
				{Key: "in", Value: bson.D{{Key: "$add", Value: bson.A{
					bson.D{{Key: "$multiply", Value: bson.A{"$$tempInCelsius", 1.8}}},
					32,
				}}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMap_UseArrayIndex(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("result", agg.Map("$scores", agg.Add("$$score", "$$idx"), agg.WithMapAs("score"), agg.WithMapArrayIndexAs("idx"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "result", Value: bson.D{{Key: "$map", Value: bson.D{
				{Key: "input", Value: "$scores"},
				{Key: "as", Value: "score"},
				{Key: "arrayIndexAs", Value: "idx"},
				{Key: "in", Value: bson.D{{Key: "$add", Value: bson.A{"$$score", "$$idx"}}}},
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

func TestMedian(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("studentId"),
			agg.Compute("testMedians", agg.Median([]string{"$test01", "$test02", "$test03"})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "studentId", Value: 1},
			{Key: "testMedians", Value: bson.D{{Key: "$median", Value: bson.D{
				{Key: "input", Value: bson.A{"$test01", "$test02", "$test03"}},
				{Key: "method", Value: "approximate"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement $mergeObjects tests when $lookup and $replaceRoot stages are implemented

// TODO: implement TestMeta_TextScore when $text and $search are implemented

func TestMeta_IndexKey(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(
			query.Field("type", query.Eq("apparel")),
		),
		agg.AddFieldsStage(
			agg.Assign("idxKey", agg.Meta("indexKey")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "type", Value: bson.D{{Key: "$eq", Value: "apparel"}}},
		}}},
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "idxKey", Value: bson.D{{Key: "$meta", Value: "indexKey"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestMillisecond(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("millisecond", agg.Millisecond("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "millisecond", Value: bson.D{
				{Key: "$millisecond", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
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

func TestMinute(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("minutes", agg.Minute("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "minutes", Value: bson.D{
				{Key: "$minute", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
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

func TestMonth(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("month", agg.Month("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "month", Value: bson.D{
				{Key: "$month", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
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

func TestObjectToArray_Example(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("dimensions", agg.ObjectToArray("$dimensions")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "dimensions", Value: bson.D{
				{Key: "$objectToArray", Value: "$dimensions"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestObjectToArray_SumNestedFields when $unwind is implemented

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

func TestPercentile(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("studentId"),
			agg.Compute("testPercentiles", agg.Percentile(
				[]string{"$test01", "$test02", "$test03"},
				[]float64{0.5, 0.95},
			),
			),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "studentId", Value: 1},
			{Key: "testPercentiles", Value: bson.D{{Key: "$percentile", Value: bson.D{
				{Key: "input", Value: bson.A{"$test01", "$test02", "$test03"}},
				{Key: "p", Value: bson.A{0.5, 0.95}},
				{Key: "method", Value: "approximate"},
			}}}},
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

// TODO: add MergeStage when $merge is implemented
func TestRand_GenerateRandomDataPoints(t *testing.T) {
	got := agg.Pipeline{
		agg.SetStage(
			agg.Assign("amount", agg.Multiply(agg.Rand(), 100)),
		),
		agg.SetStage(
			agg.Assign("amount", agg.Floor("$amount")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "amount", Value: bson.D{{Key: "$multiply", Value: bson.A{
				bson.D{{Key: "$rand", Value: bson.D{}}},
				100,
			}}}},
		}}},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "amount", Value: bson.D{{Key: "$floor", Value: "$amount"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestRand_SelectRandomItemsFromCollection when $expr is implemented

func TestRange_WithoutRangeOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("city"),
			agg.Compute("Rest stops", agg.Range(0, "$distance")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "city", Value: 1},
			{Key: "Rest stops", Value: bson.D{{Key: "$range", Value: bson.A{0, "$distance"}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRange_WithRangeStep(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("city"),
			agg.Compute("Rest stops", agg.Range(0, "$distance", agg.WithRangeStep(25))),
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

func TestReduce_Multiplication(t *testing.T) {
	got := agg.Pipeline{
		agg.GroupStage(
			"$experimentId",
			agg.Accumulate("probabilityArr", agg.PushAccumulator("$probability")),
		),
		agg.ProjectStage(
			agg.Include("description"),
			agg.Compute("results", agg.Reduce("$probabilityArr", 1, agg.Multiply("$$value", "$$this"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$experimentId"},
			{Key: "probabilityArr", Value: bson.D{{Key: "$push", Value: "$probability"}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "description", Value: 1},
			{Key: "results", Value: bson.D{{Key: "$reduce", Value: bson.D{
				{Key: "input", Value: "$probabilityArr"},
				{Key: "initialValue", Value: 1},
				{Key: "in", Value: bson.D{{Key: "$multiply", Value: bson.A{"$$value", "$$this"}}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReduce_DiscountedMerchandise(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("discountedPrice", agg.Reduce("$discounts", "$price",
				agg.Multiply("$$value", agg.Subtract(1, "$$this")))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "discountedPrice", Value: bson.D{{Key: "$reduce", Value: bson.D{
				{Key: "input", Value: "$discounts"},
				{Key: "initialValue", Value: "$price"},
				{Key: "in", Value: bson.D{{Key: "$multiply", Value: bson.A{
					"$$value",
					bson.D{{Key: "$subtract", Value: bson.A{1, "$$this"}}},
				}}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReduce_StringConcatenation(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("hobbies", query.Gt(bson.A{}))),
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("bio", agg.Reduce("$hobbies", "My hobbies include:",
				agg.Concat(
					"$$value",
					agg.Cond(agg.Eq("$$value", "My hobbies include:"), " ", ", "),
					"$$this",
				))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "hobbies", Value: bson.D{{Key: "$gt", Value: bson.A{}}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "bio", Value: bson.D{{Key: "$reduce", Value: bson.D{
				{Key: "input", Value: "$hobbies"},
				{Key: "initialValue", Value: "My hobbies include:"},
				{Key: "in", Value: bson.D{{Key: "$concat", Value: bson.A{
					"$$value",
					bson.D{{Key: "$cond", Value: bson.D{
						{Key: "if", Value: bson.D{{Key: "$eq", Value: bson.A{"$$value", "My hobbies include:"}}}},
						{Key: "then", Value: " "},
						{Key: "else", Value: ", "},
					}}},
					"$$this",
				}}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReduce_ArrayConcatenation(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("collapsed", agg.Reduce("$arr", bson.A{},
				agg.ConcatArrays("$$value", "$$this"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "collapsed", Value: bson.D{{Key: "$reduce", Value: bson.D{
				{Key: "input", Value: "$arr"},
				{Key: "initialValue", Value: bson.A{}},
				{Key: "in", Value: bson.D{{Key: "$concatArrays", Value: bson.A{"$$value", "$$this"}}}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestReduce_ComputingAMultipleReductions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("results", agg.Reduce("$arr", bson.A{},
				bson.D{
					{Key: "collapsed", Value: agg.ConcatArrays("$$value.collapsed", "$$this")},
					{Key: "firstValues", Value: agg.ConcatArrays("$$value.firstValues", agg.Slice("$$this", 1))},
				})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "results", Value: bson.D{{Key: "$reduce", Value: bson.D{
				{Key: "input", Value: "$arr"},
				{Key: "initialValue", Value: bson.A{}},
				{Key: "in", Value: bson.D{
					{Key: "collapsed", Value: bson.D{{Key: "$concatArrays", Value: bson.A{"$$value.collapsed", "$$this"}}}},
					{Key: "firstValues", Value: bson.D{{Key: "$concatArrays", Value: bson.A{
						"$$value.firstValues",
						bson.D{{Key: "$slice", Value: bson.A{"$$this", 1}}},
					}}}},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestReduce_UseAsValueAsAndArrayIndexAs when $toString is implemented

func TestRegexFind_AndItsOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFind("$description", bson.Regex{Pattern: "line"})),
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
	got := agg.Pipeline{
		// Specify i as part of the Regex type
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFind("$description", bson.Regex{Pattern: "line", Options: "i"})),
		),
		// Specify i in the options field
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFind("$description", "line", agg.WithRegexOptions("i"))),
		),
		// Mix Regex type with options field
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFind("$description", bson.Regex{Pattern: "line"}, agg.WithRegexOptions("i"))),
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
			agg.Assign("returnObject", agg.RegexFindAll("$description", bson.Regex{Pattern: "line"})),
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
	got := agg.Pipeline{
		// Specify i as part of the regex type
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFindAll("$description", bson.Regex{Pattern: "line", Options: "i"})),
		),
		// Specify i in the options field
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFindAll("$description", "line", agg.WithRegexOptions("i"))),
		),
		// Mix Regex type with options field
		agg.AddFieldsStage(
			agg.Assign("returnObject", agg.RegexFindAll("$description", bson.Regex{Pattern: "line"}, agg.WithRegexOptions("i"))),
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
				Pattern: `[a-z0-9_.+-]+@[a-z0-9_.+-]+\.[a-z0-9_.+-]+`, Options: "i"})),
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
			agg.Assign("result", agg.RegexMatch("$description", bson.Regex{Pattern: "line"})),
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
	got := agg.Pipeline{
		// Specify i as part of the Regex type
		agg.AddFieldsStage(
			agg.Assign("result", agg.RegexMatch("$description", bson.Regex{Pattern: "line", Options: "i"})),
		),
		// Specify i in the options field
		agg.AddFieldsStage(
			agg.Assign("result", agg.RegexMatch("$description", "line", agg.WithRegexOptions("i"))),
		),
		// Mix Regex type with options field
		agg.AddFieldsStage(
			agg.Assign("result", agg.RegexMatch("$description", bson.Regex{Pattern: "line"}, agg.WithRegexOptions("i"))),
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

func TestRound_WithoutRoundOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("roundedValue", agg.Round("$value")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "roundedValue", Value: bson.D{
				{Key: "$round", Value: bson.A{"$value"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestRound_WithRoundPlace(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("roundedValue", agg.Round("$value", agg.WithRoundPlace(1))),
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

func TestRtrim_WithoutTrimChars(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("description", agg.Rtrim("$description")),
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

func TestRtrim_WithTrimChars(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("description", agg.Rtrim("$description", agg.WithTrimChars("*"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "description", Value: bson.D{{Key: "$rtrim", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "chars", Value: "*"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSecond(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("seconds", agg.Second("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "seconds", Value: bson.D{
				{Key: "$second", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSerializeEJSON_CanonicalExtendedJSONExample(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("Inception"))),
		agg.ProjectStage(
			agg.Compute("ejson", agg.SerializeEJSON(agg.RootObject())),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "title", Value: bson.D{{Key: "$eq", Value: "Inception"}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "ejson", Value: bson.D{{Key: "$serializeEJSON", Value: bson.D{
				{Key: "input", Value: "$$ROOT"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSerializeEJSON_RelaxedExtendedJSONExample(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("title", query.Eq("Inception"))),
		agg.ProjectStage(
			agg.Compute("ejson", agg.SerializeEJSON(agg.RootObject(), agg.WithSerializeEJSONRelaxed(true))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "title", Value: bson.D{{Key: "$eq", Value: "Inception"}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "ejson", Value: bson.D{{Key: "$serializeEJSON", Value: bson.D{
				{Key: "input", Value: "$$ROOT"},
				{Key: "relaxed", Value: true},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestSerializeEJSON_ConvertToJSONString when the $toString operator is implemented

func TestSerializeEJSON_SerializeSpecificFields(t *testing.T) {
	got := agg.Pipeline{
		agg.MatchStage(query.Field("year", query.Gte(2010))),
		agg.ProjectStage(
			agg.Include("title"),
			agg.Compute("metadataEJSON", agg.SerializeEJSON(bson.D{
				{Key: "releaseDate", Value: "$released"},
				{Key: "runtime", Value: "$runtime"},
				{Key: "imdbRating", Value: "$imdb.rating"},
			})),
		),
	}
	want := bson.A{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "year", Value: bson.D{{Key: "$gte", Value: 2010}}},
		}}},
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "metadataEJSON", Value: bson.D{{Key: "$serializeEJSON", Value: bson.D{
				{Key: "input", Value: bson.D{
					{Key: "releaseDate", Value: "$released"},
					{Key: "runtime", Value: "$runtime"},
					{Key: "imdbRating", Value: "$imdb.rating"},
				}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSerializeEJSON_UseOnErrorForErrorHandling(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("title"),
			agg.Compute("ejson", agg.SerializeEJSON("$customField",
				agg.WithSerializeEJSONOnError(bson.D{{Key: "error", Value: "Serialization failed"}}))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "title", Value: 1},
			{Key: "ejson", Value: bson.D{{Key: "$serializeEJSON", Value: bson.D{
				{Key: "input", Value: "$customField"},
				{Key: "onError", Value: bson.D{{Key: "error", Value: "Serialization failed"}}},
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

// implement $setField tests when $replaceWith is implemented

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

func TestSimilarityCosine(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("raw", agg.SimilarityCosine("$a", "$b")),
			agg.Compute("normalized", agg.SimilarityCosine("$a", "$b", agg.WithSimilarityScore(true))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "raw", Value: bson.D{{Key: "$similarityCosine", Value: bson.D{
				{Key: "vectors", Value: bson.A{"$a", "$b"}},
			}}}},
			{Key: "normalized", Value: bson.D{{Key: "$similarityCosine", Value: bson.D{
				{Key: "vectors", Value: bson.A{"$a", "$b"}},
				{Key: "score", Value: true},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSimilarityDotProduct(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("raw", agg.SimilarityDotProduct("$a", "$b")),
			agg.Compute("normalized", agg.SimilarityDotProduct("$a", "$b", agg.WithSimilarityScore(true))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "raw", Value: bson.D{{Key: "$similarityDotProduct", Value: bson.D{
				{Key: "vectors", Value: bson.A{"$a", "$b"}},
			}}}},
			{Key: "normalized", Value: bson.D{{Key: "$similarityDotProduct", Value: bson.D{
				{Key: "vectors", Value: bson.A{"$a", "$b"}},
				{Key: "score", Value: true},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestSimilarityEuclidean(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("raw", agg.SimilarityEuclidean("$a", "$b")),
			agg.Compute("normalized", agg.SimilarityEuclidean("$a", "$b", agg.WithSimilarityScore(true))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "raw", Value: bson.D{{Key: "$similarityEuclidean", Value: bson.D{
				{Key: "vectors", Value: bson.A{"$a", "$b"}},
			}}}},
			{Key: "normalized", Value: bson.D{{Key: "$similarityEuclidean", Value: bson.D{
				{Key: "vectors", Value: bson.A{"$a", "$b"}},
				{Key: "score", Value: true},
			}}}},
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

func TestSlice_WithoutSliceOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("threeFavorites", agg.Slice("$favorites", 3)),
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

func TestSlice_WithSlicePosition(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("threeFavorites", agg.Slice("$favorites", 3, agg.WithSlicePosition(2))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "threeFavorites", Value: bson.D{{Key: "$slice", Value: bson.A{"$favorites", 2, 3}}}},
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

func TestSubtype(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("result", agg.Subtype("$myBinDataField")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "result", Value: bson.D{
				{Key: "$subtype", Value: "$myBinDataField"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

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

func TestSwitch(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("name"),
			agg.Compute("summary", agg.Switch([]agg.SwitchCase{
				agg.Case(agg.Gte(agg.Avg("$scores"), 90), "Doing great!"),
				agg.Case(agg.And(agg.Gte(agg.Avg("$scores"), 80), agg.Lt(agg.Avg("$scores"), 90)), "Doing pretty well."),
				agg.Case(agg.Lt(agg.Avg("$scores"), 80), "Needs improvement."),
			},
				agg.WithSwitchDefault("No scores found."),
			)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "name", Value: 1},
			{Key: "summary", Value: bson.D{
				{Key: "$switch", Value: bson.D{
					{Key: "branches", Value: bson.A{
						bson.D{
							{Key: "case", Value: bson.D{
								{Key: "$gte", Value: bson.A{
									bson.D{{Key: "$avg", Value: bson.A{"$scores"}}},
									90,
								}},
							}},
							{Key: "then", Value: "Doing great!"},
						},
						bson.D{
							{Key: "case", Value: bson.D{
								{Key: "$and", Value: bson.A{
									bson.D{{Key: "$gte", Value: bson.A{
										bson.D{{Key: "$avg", Value: bson.A{"$scores"}}},
										80,
									}}},
									bson.D{{Key: "$lt", Value: bson.A{
										bson.D{{Key: "$avg", Value: bson.A{"$scores"}}},
										90,
									}}},
								}},
							}},
							{Key: "then", Value: "Doing pretty well."},
						},
						bson.D{
							{Key: "case", Value: bson.D{
								{Key: "$lt", Value: bson.A{
									bson.D{{Key: "$avg", Value: bson.A{"$scores"}}},
									80,
								}},
							}},
							{Key: "then", Value: "Needs improvement."},
						},
					}},
					{Key: "default", Value: "No scores found."},
				}},
			}},
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

func TestToArray_ConvertStringToArray(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("numbers", agg.ToArray("[1, 2, 3]")),
			agg.Compute("documents", agg.ToArray(`[{"a": 1}, {"b": 2}]`)),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "numbers", Value: bson.D{{Key: "$toArray", Value: "[1, 2, 3]"}}},
			{Key: "documents", Value: bson.D{{Key: "$toArray", Value: `[{"a": 1}, {"b": 2}]`}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestToArray_ConvertBinDataToArray(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("original", "$v"),
			agg.Compute("asArray", agg.ToArray("$v")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "original", Value: "$v"},
			{Key: "asArray", Value: bson.D{{Key: "$toArray", Value: "$v"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement $toBool tests when $switch is implemented

func TestToDate(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("convertedDate", agg.ToDate("$order_date")),
		),
		agg.SortStage(
			agg.Sort("convertedDate", agg.Asc),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "convertedDate", Value: bson.D{
				{Key: "$toDate", Value: "$order_date"},
			}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "convertedDate", Value: int32(1)},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestToDecimal(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("convertedPrice", agg.ToDecimal("$price")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "convertedPrice", Value: bson.D{{Key: "$toDecimal", Value: "$price"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestToDouble(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("degrees", agg.ToDouble(agg.SubstrBytes("$temp", 0, 4))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "degrees", Value: bson.D{{Key: "$toDouble", Value: bson.D{
				{Key: "$substrBytes", Value: bson.A{"$temp", 0, 4}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement $toHashedIndexKey tests when $documents is implemented

func TestToInt(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("convertedQty", agg.ToInt("$qty")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "convertedQty", Value: bson.D{{Key: "$toInt", Value: "$qty"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestToLong(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("convertedQty", agg.ToLong("$qty")),
		),
		agg.SortStage(agg.Sort("convertedQty", agg.Desc)),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "convertedQty", Value: bson.D{{Key: "$toLong", Value: "$qty"}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "convertedQty", Value: -1},
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

func TestToObject(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("parsedConfig", agg.ToObject("$config")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "parsedConfig", Value: bson.D{{Key: "$toObject", Value: "$config"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestToObjectId(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("convertedId", agg.ToObjectId("$_id")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "convertedId", Value: bson.D{{Key: "$toObjectId", Value: "$_id"}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestToString(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("convertedZipCode", agg.ToString("$zipcode")),
		),
		agg.SortStage(agg.Sort("convertedZipCode", agg.Asc)),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "convertedZipCode", Value: bson.D{{Key: "$toString", Value: "$zipcode"}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "convertedZipCode", Value: 1},
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

func TestTrim_WithoutTrimChars(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("description", agg.Trim("$description")),
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

func TestTrim_WithTrimChars(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Include("item"),
			agg.Compute("description", agg.Trim("$description", agg.WithTrimChars("*"))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "item", Value: 1},
			{Key: "description", Value: bson.D{{Key: "$trim", Value: bson.D{
				{Key: "input", Value: "$description"},
				{Key: "chars", Value: "*"},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestTrunc_WithoutTruncOptions(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("truncatedValue", agg.Trunc("$value")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "truncatedValue", Value: bson.D{
				{Key: "$trunc", Value: bson.A{"$value"}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestTrunc_WithTruncPlace(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("truncatedValue", agg.Trunc("$value", agg.WithTruncPlace(1))),
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

// TODO: implement $tsIncrement tests when $expr is implemented

func TestTsSecond_ObtainNumberOfSecondsFromTimestampField(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Include("saleTimestamp"),
			agg.Compute("saleSeconds", agg.TsSecond("$saleTimestamp")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: int32(0)},
			{Key: "saleTimestamp", Value: int32(1)},
			{Key: "saleSeconds", Value: bson.D{
				{Key: "$tsSecond", Value: "$saleTimestamp"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestTsSecond_UseInChangeStreamCursorToMonitorCollectionChanges(t *testing.T) {
	got := agg.Pipeline{
		agg.AddFieldsStage(
			agg.Assign("clusterTimeSeconds", agg.TsSecond("$clusterTime")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$addFields", Value: bson.D{
			{Key: "clusterTimeSeconds", Value: bson.D{
				{Key: "$tsSecond", Value: "$clusterTime"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestType(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("a", agg.Type("$a")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "a", Value: bson.D{
				{Key: "$type", Value: "$a"},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement $unsetField tests when $replaceWith is implemented

func TestWeek(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("week", agg.Week("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "week", Value: bson.D{
				{Key: "$week", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
			}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestYear(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Compute("year", agg.Year("$date")),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "year", Value: bson.D{
				{Key: "$year", Value: bson.D{
					{Key: "date", Value: "$date"},
				}},
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
			})),
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

func TestZip_WithZipUseLongestLength(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("transposed", agg.Zip([]any{
				agg.ArrayElemAt("$matrix", 0),
				agg.ArrayElemAt("$matrix", 1),
			}, agg.WithZipUseLongestLength(true))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "transposed", Value: bson.D{{Key: "$zip", Value: bson.D{
				{Key: "inputs", Value: bson.A{
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$matrix", 0}}},
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$matrix", 1}}},
				}},
				{Key: "useLongestLength", Value: true},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

func TestZip_WithZipUseLongestLengthAndZipDefaults(t *testing.T) {
	got := agg.Pipeline{
		agg.ProjectStage(
			agg.Exclude("_id"),
			agg.Compute("transposed", agg.Zip([]any{
				agg.ArrayElemAt("$matrix", 0),
				agg.ArrayElemAt("$matrix", 1),
			}, agg.WithZipUseLongestLength(true), agg.WithZipDefaults(0, 0))),
		),
	}
	want := bson.A{
		bson.D{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "transposed", Value: bson.D{{Key: "$zip", Value: bson.D{
				{Key: "inputs", Value: bson.A{
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$matrix", 0}}},
					bson.D{{Key: "$arrayElemAt", Value: bson.A{"$matrix", 1}}},
				}},
				{Key: "useLongestLength", Value: true},
				{Key: "defaults", Value: bson.A{0, 0}},
			}}}},
		}}},
	}
	assertPipelineEqual(t, got, want)
}

// TODO: implement TestZip_FilteringAndPreservingIndexes when $let is implemented
