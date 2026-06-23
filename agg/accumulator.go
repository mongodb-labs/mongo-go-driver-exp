package agg

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Accumulator struct {
	doc bson.D
}

func (a Accumulator) MarshalBSON() ([]byte, error) {
	return bson.Marshal(a.doc)
}

// --- $accumulator ---

// CustomAccumulator defines a custom accumulator function using JavaScript ($accumulator).
// Pass nil for initArgs or finalize to omit those optional fields.
func CustomAccumulator[A ArrayTypes](init, accumulate string, accumulateArgs A, merge, lang string, initArgs []any, finalize *string) Accumulator {
	doc := bson.D{{Key: "init", Value: init}}
	if initArgs != nil {
		doc = append(doc, bson.E{Key: "initArgs", Value: initArgs})
	}
	doc = append(doc,
		bson.E{Key: "accumulate", Value: accumulate},
		bson.E{Key: "accumulateArgs", Value: accumulateArgs},
		bson.E{Key: "merge", Value: merge},
	)
	if finalize != nil {
		doc = append(doc, bson.E{Key: "finalize", Value: *finalize})
	}
	doc = append(doc, bson.E{Key: "lang", Value: lang})
	return Accumulator{doc: bson.D{{Key: "$accumulator", Value: doc}}}
}

// --- $addToSet ---

// AddToSetAccumulator returns an array of unique expression values for each group ($addToSet).
func AddToSetAccumulator(expr Expr) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$addToSet", Value: expr}}}
}

// --- $avg ---

// AvgAccumulator returns the average of the given numeric expression across documents in the group ($avg).
func AvgAccumulator[T NumberTypes](expr T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$avg", Value: expr}}}
}

// --- $bottom ---

// BottomAccumulator returns the bottom element within a group according to the specified sort order ($bottom).
func BottomAccumulator(output Expr, sortBy ...SortField) Accumulator {
	sortDoc := make(bson.D, len(sortBy))
	for i, f := range sortBy {
		sf := f.sortField()
		sortDoc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return Accumulator{doc: bson.D{{Key: "$bottom", Value: bson.D{
		{Key: "sortBy", Value: sortDoc},
		{Key: "output", Value: output},
	}}}}
}

// --- $bottomN ---

// BottomNAccumulator returns the bottom n elements within a group according to the specified sort order ($bottomN).
func BottomNAccumulator[T NumberTypes](n T, output Expr, sortBy ...SortField) Accumulator {
	sortDoc := make(bson.D, len(sortBy))
	for i, f := range sortBy {
		sf := f.sortField()
		sortDoc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return Accumulator{doc: bson.D{{Key: "$bottomN", Value: bson.D{
		{Key: "n", Value: n},
		{Key: "sortBy", Value: sortDoc},
		{Key: "output", Value: output},
	}}}}
}

// --- $concatArrays ---

// ConcatArraysAccumulator concatenates the arrays from each document in the group ($concatArrays).
func ConcatArraysAccumulator[T ArrayTypes](expr T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$concatArrays", Value: expr}}}
}

// --- $count ---

// CountAccumulator returns the number of documents in the group ($count).
func CountAccumulator() Accumulator {
	return Accumulator{doc: bson.D{{Key: "$count", Value: bson.D{}}}}
}

// --- $covariancePop ---

// CovariancePopAccumulator returns the population covariance of two numeric expressions ($covariancePop).
func CovariancePopAccumulator[T, U NumberTypes](expr1 T, expr2 U) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$covariancePop", Value: bson.A{expr1, expr2}}}}
}

// --- $covarianceSamp ---

// CovarianceSampAccumulator returns the sample covariance of two numeric expressions ($covarianceSamp).
func CovarianceSampAccumulator[T, U NumberTypes](expr1 T, expr2 U) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$covarianceSamp", Value: bson.A{expr1, expr2}}}}
}

// --- $denseRank ---

// DenseRankAccumulator returns the document rank with no gaps in the $setWindowFields stage partition ($denseRank).
func DenseRankAccumulator() Accumulator {
	return Accumulator{doc: bson.D{{Key: "$denseRank", Value: bson.D{}}}}
}

// --- $derivative ---

// DerivativeAccumulator returns the average rate of change within the specified window ($derivative).
// Pass a non-nil unit to specify a time unit (e.g. "week").
func DerivativeAccumulator(input Expr, unit *string) Accumulator {
	doc := bson.D{{Key: "input", Value: input}}
	if unit != nil {
		doc = append(doc, bson.E{Key: "unit", Value: *unit})
	}
	return Accumulator{doc: bson.D{{Key: "$derivative", Value: doc}}}
}

// --- $documentNumber ---

// DocumentNumberAccumulator returns the position of a document in the $setWindowFields stage partition ($documentNumber).
func DocumentNumberAccumulator() Accumulator {
	return Accumulator{doc: bson.D{{Key: "$documentNumber", Value: bson.D{}}}}
}

// --- $expMovingAvg ---

// ExpMovingAvgAccumulator returns the exponential moving average of the input expression ($expMovingAvg).
// Pass a non-nil n to use a historical document count, or a non-nil alpha to use an exponential decay value.
// Exactly one of n or alpha must be non-nil.
func ExpMovingAvgAccumulator[T NumberTypes](input T, n *int32, alpha *float64) Accumulator {
	doc := bson.D{{Key: "input", Value: input}}
	if n != nil {
		doc = append(doc, bson.E{Key: "N", Value: *n})
	}
	if alpha != nil {
		doc = append(doc, bson.E{Key: "alpha", Value: *alpha})
	}
	return Accumulator{doc: bson.D{{Key: "$expMovingAvg", Value: doc}}}
}

// --- $first ---

// FirstAccumulator returns the result of an expression for the first document in a group or window ($first).
func FirstAccumulator(expr Expr) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$first", Value: expr}}}
}

// --- $firstN ---

// FirstNAccumulator returns the first n elements within a group ($firstN).
func FirstNAccumulator[T NumberTypes](input Expr, n T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$firstN", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "n", Value: n},
	}}}}
}

// --- $integral ---

// IntegralAccumulator returns the approximation of the area under a curve ($integral).
// Pass a non-nil unit to specify a time unit (e.g. "week").
func IntegralAccumulator(input Expr, unit *string) Accumulator {
	doc := bson.D{{Key: "input", Value: input}}
	if unit != nil {
		doc = append(doc, bson.E{Key: "unit", Value: *unit})
	}
	return Accumulator{doc: bson.D{{Key: "$integral", Value: doc}}}
}

// --- $last ---

// LastAccumulator returns the result of an expression for the last document in a group or window ($last).
func LastAccumulator(expr Expr) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$last", Value: expr}}}
}

// --- $lastN ---

// LastNAccumulator returns the last n elements of input across documents in the group ($lastN).
func LastNAccumulator[T NumberTypes](input Expr, n T) Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$lastN", Value: bson.D{
			{Key: "input", Value: input},
			{Key: "n", Value: n},
		}}},
	}
}

// --- $linearFill ---

// LinearFillAccumulator fills null and missing fields using linear interpolation based on surrounding values ($linearFill).
func LinearFillAccumulator[T NumberTypes](expr T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$linearFill", Value: expr}}}
}

// --- $locf ---

// LocfAccumulator carries forward the last non-null value for null and missing fields ($locf).
func LocfAccumulator(expr Expr) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$locf", Value: expr}}}
}

// --- $max ---

// MaxAccumulator returns the maximum value of expr across documents in the group ($max).
func MaxAccumulator(expr Expr) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$max", Value: expr}}}
}

// --- $maxN ---

// MaxNAccumulator returns the n largest values in an array ($maxN).
func MaxNAccumulator[T NumberTypes](input Expr, n T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$maxN", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "n", Value: n},
	}}}}
}

// --- $median ---

// MedianAccumulator returns an approximation of the median (50th percentile) value ($median).
func MedianAccumulator[T NumberTypes](input T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$median", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "method", Value: "approximate"},
	}}}}
}

// --- $mergeObjects ---

// MergeObjectsAccumulator combines multiple documents into a single document ($mergeObjects).
func MergeObjectsAccumulator(expr Expr) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$mergeObjects", Value: expr}}}
}

// --- $min ---

// MinAccumulator returns the minimum value of expr across documents in the group ($min).
func MinAccumulator(expr Expr) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$min", Value: expr}}}
}

// --- $minMaxScaler ---

// MinMaxScalerAccumulator normalizes a numeric expression within a window of values to the range [0, 1] ($minMaxScaler).
func MinMaxScalerAccumulator[T NumberTypes](input T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$minMaxScaler", Value: bson.D{
		{Key: "input", Value: input},
	}}}}
}

// MinMaxScalerRangeAccumulator normalizes a numeric expression within a window of values to a custom [min, max] range ($minMaxScaler).
func MinMaxScalerRangeAccumulator[T NumberTypes, U Number](input T, min, max U) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$minMaxScaler", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "min", Value: min},
		{Key: "max", Value: max},
	}}}}
}

// --- $minN ---

// MinNAccumulator returns the n smallest values in an array ($minN).
func MinNAccumulator[T NumberTypes](input Expr, n T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$minN", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "n", Value: n},
	}}}}
}

// --- $percentile ---

// PercentileAccumulator returns the percentile values of input at the given probabilities p
// across documents in the group ($percentile).
// p must be an array of numeric values between 0 and 1.
func PercentileAccumulator[T NumberTypes, U ArrayTypes | []float32 | []float64](input T, p U) Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$percentile", Value: bson.D{
			{Key: "input", Value: input},
			{Key: "p", Value: p},
			// Currently the method must always be "approximate". If this changes, we need to
			// add an argument for this.
			{Key: "method", Value: "approximate"},
		}}},
	}
}

// --- $push ---

// PushAccumulator appends expr to an array of values accumulated across documents in the group ($push).
func PushAccumulator(expr Expr) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$push", Value: expr}}}
}

// --- $rank ---

// RankAccumulator returns the document rank relative to other documents in the $setWindowFields stage partition ($rank).
func RankAccumulator() Accumulator {
	return Accumulator{doc: bson.D{{Key: "$rank", Value: bson.D{}}}}
}

// --- $setUnion ---

// SetUnionAccumulator returns an array of unique values from the given array field across all documents in the group ($setUnion).
func SetUnionAccumulator[T ArrayTypes](expr T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$setUnion", Value: expr}}}
}

// --- $shift ---

// ShiftAccumulator returns the value from a document in a specified position relative to the current document ($shift).
// Pass a non-nil defaultExpr to specify a value for out-of-bounds positions; it must evaluate to a constant.
func ShiftAccumulator(output Expr, by int32, defaultExpr any) Accumulator {
	doc := bson.D{
		{Key: "output", Value: output},
		{Key: "by", Value: by},
	}
	if defaultExpr != nil {
		doc = append(doc, bson.E{Key: "default", Value: defaultExpr})
	}
	return Accumulator{doc: bson.D{{Key: "$shift", Value: doc}}}
}

// --- $stdDevPop ---

// StdDevPopAccumulator calculates the population standard deviation of the input values ($stdDevPop).
func StdDevPopAccumulator[T NumberTypes](expr T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$stdDevPop", Value: expr}}}
}

// --- $stdDevSamp ---

// StdDevSampAccumulator calculates the sample standard deviation of the input values ($stdDevSamp).
func StdDevSampAccumulator[T NumberTypes](expr T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$stdDevSamp", Value: expr}}}
}

// --- $sum ---

// SumAccumulator returns the sum of the given numeric expression across documents in the group ($sum).
func SumAccumulator[T NumberTypes](expr T) Accumulator {
	return Accumulator{doc: bson.D{{Key: "$sum", Value: expr}}}
}

// --- $top ---

// TopAccumulator returns the top element within a group according to the specified sort order ($top).
func TopAccumulator(output Expr, sortBy ...SortField) Accumulator {
	sortDoc := make(bson.D, len(sortBy))
	for i, f := range sortBy {
		sf := f.sortField()
		sortDoc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return Accumulator{doc: bson.D{{Key: "$top", Value: bson.D{
		{Key: "sortBy", Value: sortDoc},
		{Key: "output", Value: output},
	}}}}
}

// --- $topN ---

// TopNAccumulator returns the top n elements within a group according to the specified sort order ($topN).
func TopNAccumulator[T NumberTypes](n T, output Expr, sortBy ...SortField) Accumulator {
	sortDoc := make(bson.D, len(sortBy))
	for i, f := range sortBy {
		sf := f.sortField()
		sortDoc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return Accumulator{doc: bson.D{{Key: "$topN", Value: bson.D{
		{Key: "n", Value: n},
		{Key: "sortBy", Value: sortDoc},
		{Key: "output", Value: output},
	}}}}
}
