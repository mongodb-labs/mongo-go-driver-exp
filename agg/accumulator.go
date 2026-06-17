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

// LastNAccumulator returns the last n elements of input across documents in the group ($lastN).
func LastNAccumulator[T ArrayTypes, U NumberTypes](input T, n U) Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$lastN", Value: bson.D{
			{Key: "input", Value: input},
			{Key: "n", Value: n},
		}}},
	}
}

// AvgAccumulator returns the average of the given numeric expression across documents in the group ($avg).
func AvgAccumulator[T NumberTypes](expr T) Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$avg", Value: expr}},
	}
}

// SumAccumulator returns the sum of the given numeric expression across documents in the group ($sum).
func SumAccumulator[T NumberTypes](expr T) Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$sum", Value: expr}},
	}
}

// MinAccumulator returns the minimum value of expr across documents in the group ($min).
func MinAccumulator[T Expr](expr T) Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$min", Value: expr}},
	}
}

// MaxAccumulator returns the maximum value of expr across documents in the group ($max).
func MaxAccumulator[T Expr](expr T) Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$max", Value: expr}},
	}
}

// PercentileAccumulator returns the percentile values of input at the given probabilities p
// across documents in the group ($percentile).
// p must be an array of numeric values between 0 and 1.
func PercentileAccumulator[T NumberTypes, U ArrayTypes | []float32 | []float64](input T, p U) Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$percentile", Value: bson.D{
			{Key: "input", Value: input},
			{Key: "p", Value: p},
			// TODO: Currently the method must always be "approximate". Do we need an argument for that?
			{Key: "method", Value: "approximate"},
		}}},
	}
}

// CountAccumulator returns the number of documents in the group ($count).
func CountAccumulator() Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$count", Value: bson.D{}}},
	}
}

// PushAccumulator appends expr to an array of values accumulated across documents in the group ($push).
func PushAccumulator[T Expr](expr T) Accumulator {
	return Accumulator{
		doc: bson.D{{Key: "$push", Value: expr}},
	}
}
