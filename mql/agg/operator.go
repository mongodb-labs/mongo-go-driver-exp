package agg

import "go.mongodb.org/mongo-driver/v2/bson"

// Abs returns the absolute value of a number ($abs).
func Abs[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$abs", Value: expr}}}
}

// Acos returns the arccosine of a value in radians ($acos).
func Acos[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$acos", Value: expr}}}
}

// Acosh returns the inverse hyperbolic cosine of a value in radians ($acosh).
func Acosh[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$acosh", Value: expr}}}
}

// Add returns the sum of the given numeric expressions ($add).
// TODO: $add also accepts a date + milliseconds to produce a date;
// that variant is not yet modeled here.
func Add[T NumberResolver, U NumberResolver](value T, values ...U) NumberExpr {
	v := make([]any, len(values)+1)
	v[0] = value
	for i := range values {
		v[i+1] = values[i]
	}
	return NumberExpr{
		expr: bson.D{{Key: "$add", Value: v}},
	}
}

// AllElementsTrue returns true if no element of the array evaluates to false ($allElementsTrue).
func AllElementsTrue[T ArrayResolver](expr T) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$allElementsTrue", Value: bson.A{expr}}}}
}

// And returns true only when all expressions evaluate to true ($and).
func And[T BoolResolver](exprs ...T) BoolExpr {
	a := make(bson.A, len(exprs))
	for i, v := range exprs {
		a[i] = v
	}
	return BoolExpr{expr: bson.D{{Key: "$and", Value: a}}}
}

// AnyElementTrue returns true if any element of the array evaluates to true ($anyElementTrue).
func AnyElementTrue[T ArrayResolver](expr T) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$anyElementTrue", Value: bson.A{expr}}}}
}

// ArrayElemAt returns the element at the specified array index ($arrayElemAt).
func ArrayElemAt[T ArrayResolver](array T, idx Expr) AnyExpr {
	return AnyExpr{expr: bson.D{{Key: "$arrayElemAt", Value: bson.A{array, idx}}}}
}

// ArrayToObject converts an array of key-value pairs to a document ($arrayToObject).
func ArrayToObject[T ArrayResolver](array T) ObjectExpr {
	return ObjectExpr{expr: bson.D{{Key: "$arrayToObject", Value: array}}}
}

// Asin returns the arcsine of a value in radians ($asin).
func Asin[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$asin", Value: expr}}}
}

// Asinh returns the inverse hyperbolic sine of a value in radians ($asinh).
func Asinh[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$asinh", Value: expr}}}
}

// Atan returns the arctangent of a value in radians ($atan).
func Atan[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$atan", Value: expr}}}
}

// Atan2 returns the arctangent of y/x in radians ($atan2).
func Atan2[T NumberResolver, U NumberResolver](y T, x U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$atan2", Value: bson.A{y, x}}}}
}

// Atanh returns the inverse hyperbolic tangent of a value in radians ($atanh).
func Atanh[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$atanh", Value: expr}}}
}

// Avg returns the average of numeric expressions ($avg).
func Avg(exprs ...Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$avg", Value: exprs}}}
}

// BitAnd returns the bitwise AND of int or long values ($bitAnd).
// MongoDB requires int or long operands; other numeric types cause a runtime error.
func BitAnd(exprs ...Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$bitAnd", Value: exprs}}}
}

// BitNot returns the bitwise NOT of an int or long value ($bitNot).
// MongoDB requires an int or long operand; other numeric types cause a runtime error.
func BitNot[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$bitNot", Value: expr}}}
}

// BitOr returns the bitwise OR of int or long values ($bitOr).
// MongoDB requires int or long operands; other numeric types cause a runtime error.
func BitOr(exprs ...Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$bitOr", Value: exprs}}}
}

// BitXor returns the bitwise XOR of int or long values ($bitXor).
// MongoDB requires int or long operands; other numeric types cause a runtime error.
func BitXor(exprs ...Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$bitXor", Value: exprs}}}
}

// Bottom returns the bottom element within an array according to the specified sort order ($bottom).
// This is the expression operator (MongoDB 7.0+) that takes an input array.
// See BottomAccumulator for the $group/$setWindowFields accumulator form.
func Bottom[T ArrayResolver](input T, output Expr, sortBy ...SortField) AnyExpr {
	sortDoc := make(bson.D, len(sortBy))
	for i, f := range sortBy {
		sf := f.sortField()
		sortDoc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return AnyExpr{expr: bson.D{{Key: "$bottom", Value: bson.D{
		{Key: "sortBy", Value: sortDoc},
		{Key: "output", Value: output},
		{Key: "input", Value: input},
	}}}}
}

// BottomN returns the bottom n elements within an array according to the specified sort order ($bottomN).
// This is the expression operator (MongoDB 7.0+) that takes an input array.
// See BottomNAccumulator for the $group/$setWindowFields accumulator form.
func BottomN[T ArrayResolver](n Expr, input T, output Expr, sortBy ...SortField) ArrayExpr {
	sortDoc := make(bson.D, len(sortBy))
	for i, f := range sortBy {
		sf := f.sortField()
		sortDoc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return ArrayExpr{expr: bson.D{{Key: "$bottomN", Value: bson.D{
		{Key: "n", Value: n},
		{Key: "sortBy", Value: sortDoc},
		{Key: "output", Value: output},
		{Key: "input", Value: input},
	}}}}
}

// Ceil returns the smallest integer greater than or equal to the number ($ceil).
func Ceil[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$ceil", Value: expr}}}
}

// Cmp returns 0 if a == b, 1 if a > b, and -1 if a < b ($cmp).
func Cmp(a Expr, b Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$cmp", Value: bson.A{a, b}}}}
}

// Concat concatenates the given string expressions ($concat).
func Concat[T StringResolver, U StringResolver](value T, values ...U) StringExpr {
	v := make([]any, len(values)+1)
	v[0] = value
	for i := range values {
		v[i+1] = values[i]
	}
	return StringExpr{expr: bson.D{{Key: "$concat", Value: v}}}
}

// ConcatArrays concatenates arrays to return the concatenated array ($concatArrays).
func ConcatArrays[T ArrayResolver, U ArrayResolver](array T, arrays ...U) ArrayExpr {
	v := make([]any, len(arrays)+1)
	v[0] = array
	for i := range arrays {
		v[i+1] = arrays[i]
	}
	return ArrayExpr{expr: bson.D{{Key: "$concatArrays", Value: v}}}
}

// Cos returns the cosine of a value in radians ($cos).
func Cos[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$cos", Value: expr}}}
}

// Cosh returns the hyperbolic cosine of a value in radians ($cosh).
func Cosh[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$cosh", Value: expr}}}
}

// DegreesToRadians converts a value from degrees to radians ($degreesToRadians).
func DegreesToRadians[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$degreesToRadians", Value: expr}}}
}

// Divide returns a divided by b ($divide).
func Divide[T NumberResolver, U NumberResolver](a T, b U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$divide", Value: bson.A{a, b}}}}
}

// Eq returns true if a equals b ($eq).
func Eq(a Expr, b Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$eq", Value: bson.A{a, b}}}}
}

// Exp raises Euler's number e to the specified exponent ($exp).
func Exp[T NumberResolver](exponent T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$exp", Value: exponent}}}
}

// FilterArray selects elements of input for which cond evaluates to true ($filter).
// as names the variable for each element; pass "" to use the MongoDB default ("this").
func FilterArray[T ArrayResolver, U BoolResolver](input T, as string, cond U, limit ...Expr) ArrayExpr {
	args := bson.D{bson.E{Key: "input", Value: input}}
	if as != "" {
		args = append(args, bson.E{Key: "as", Value: as})
	}
	args = append(args,
		bson.E{Key: "cond", Value: cond},
	)
	if len(limit) > 1 {
		panic("FilterArray: at most one limit expression may be provided")
	}
	if len(limit) == 1 {
		args = append(args, bson.E{Key: "limit", Value: limit[0]})
	}
	return ArrayExpr{expr: bson.D{{Key: "$filter", Value: args}}}
}

// First returns the first element of the array expression ($first).
// This is the array expression operator (MongoDB 4.4+).
// See FirstAccumulator for the $group/$setWindowFields accumulator form.
func First[T ArrayResolver](expr T) AnyExpr {
	return AnyExpr{expr: bson.D{{Key: "$first", Value: expr}}}
}

// FirstN returns the first n elements from an array ($firstN).
// This is the array expression operator (MongoDB 5.1+).
// See FirstNAccumulator for the $group/$setWindowFields accumulator form.
func FirstN[T ArrayResolver](n Expr, input T) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$firstN", Value: bson.D{
		{Key: "n", Value: n},
		{Key: "input", Value: input},
	}}}}
}

// Floor returns the largest integer less than or equal to the number ($floor).
func Floor[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$floor", Value: expr}}}
}

// Gt returns true if a is greater than b ($gt).
func Gt(a Expr, b Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$gt", Value: bson.A{a, b}}}}
}

// Gte returns true if a is greater than or equal to b ($gte).
func Gte(a Expr, b Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$gte", Value: bson.A{a, b}}}}
}

// IfNull returns the first non-null expression among val and more, or fallback
// if all preceding ones are null ($ifNull).
func IfNull(val Expr, fallback Expr, more ...Expr) AnyExpr {
	v := make([]any, len(more)+2)
	v[0] = val
	v[1] = fallback
	for i := range more {
		v[i+2] = more[i]
	}
	return AnyExpr{expr: bson.D{{Key: "$ifNull", Value: v}}}
}

// In returns true if expr is present in array ($in).
func In[U ArrayResolver](expr Expr, array U) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$in", Value: bson.A{expr, array}}}}
}

// IndexOfArray searches an array for a value and returns the index of the first occurrence ($indexOfArray).
// Optionally provide start and end index bounds via IndexOfOptions.
func IndexOfArray[T ArrayResolver](array T, search Expr, start Expr, end Expr) NumberExpr {
	args := bson.A{array, search}
	if start != nil || end != nil {
		if start == nil {
			start = 0
		}
		args = append(args, start)
		if end != nil {
			args = append(args, end)
		}
	}
	return NumberExpr{expr: bson.D{{Key: "$indexOfArray", Value: args}}}
}

// Last returns the last element of the array expression ($last).
// This is the array expression operator (MongoDB 4.4+).
// See LastAccumulator for the $group/$setWindowFields accumulator form.
func Last[T ArrayResolver](expr T) AnyExpr {
	return AnyExpr{expr: bson.D{{Key: "$last", Value: expr}}}
}

// LastN returns the last n elements from an array ($lastN).
// This is the array expression operator (MongoDB 5.1+).
// See LastNAccumulator for the $group/$setWindowFields accumulator form.
func LastN[T ArrayResolver](n Expr, input T) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$lastN", Value: bson.D{
		{Key: "n", Value: n},
		{Key: "input", Value: input},
	}}}}
}

// Ln calculates the natural logarithm of a number ($ln).
func Ln[T NumberResolver](number T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$ln", Value: number}}}
}

// Log calculates the log of number in the specified base ($log).
func Log[T NumberResolver, U NumberResolver](number T, base U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$log", Value: bson.A{number, base}}}}
}

// Log10 calculates the log base 10 of a number ($log10).
func Log10[T NumberResolver](number T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$log10", Value: number}}}
}

// Lt returns true if a is less than b ($lt).
func Lt(a Expr, b Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$lt", Value: bson.A{a, b}}}}
}

// Lte returns true if a is less than or equal to b ($lte).
func Lte(a Expr, b Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$lte", Value: bson.A{a, b}}}}
}

// Max returns the maximum value among the given expressions ($max).
// Accepts any expression type (Expr = any).
func Max(value Expr, values ...Expr) AnyExpr {
	v := make([]any, len(values)+1)
	v[0] = value
	for i := range values {
		v[i+1] = values[i]
	}
	return AnyExpr{expr: bson.D{{Key: "$max", Value: v}}}
}

// MaxN returns the n largest values in an array ($maxN).
// This is the array expression operator (MongoDB 5.1+).
// See MaxNAccumulator for the $group/$setWindowFields accumulator form.
func MaxN[T ArrayResolver](n Expr, input T) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$maxN", Value: bson.D{
		{Key: "n", Value: n},
		{Key: "input", Value: input},
	}}}}
}

// Min returns the minimum value among the given expressions ($min).
// Accepts any expression type (Expr = any).
func Min(value Expr, values ...Expr) AnyExpr {
	v := make([]any, len(values)+1)
	v[0] = value
	for i := range values {
		v[i+1] = values[i]
	}
	return AnyExpr{expr: bson.D{{Key: "$min", Value: v}}}
}

// MinN returns the n smallest values in an array ($minN).
// This is the array expression operator (MongoDB 5.1+).
// See MinNAccumulator for the $group/$setWindowFields accumulator form.
func MinN[T ArrayResolver](n Expr, input T) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$minN", Value: bson.D{
		{Key: "n", Value: n},
		{Key: "input", Value: input},
	}}}}
}

// Mod returns the remainder of dividing dividend by divisor ($mod).
func Mod[T NumberResolver, U NumberResolver](dividend T, divisor U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$mod", Value: bson.A{dividend, divisor}}}}
}

// Multiply returns the product of the given numeric expressions ($multiply).
func Multiply[T NumberResolver, U NumberResolver](value T, values ...U) NumberExpr {
	v := make([]any, len(values)+1)
	v[0] = value
	for i := range values {
		v[i+1] = values[i]
	}
	return NumberExpr{expr: bson.D{{Key: "$multiply", Value: v}}}
}

// Ne returns true if a does not equal b ($ne).
func Ne(a Expr, b Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$ne", Value: bson.A{a, b}}}}
}

// Not returns the boolean inverse of e ($not).
// $not takes a single-element array in the aggregation expression syntax.
func Not[T BoolResolver](e T) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$not", Value: bson.A{e}}}}
}

// Or returns true when any expression evaluates to true ($or).
func Or[T BoolResolver](exprs ...T) BoolExpr {
	a := make(bson.A, len(exprs))
	for i, v := range exprs {
		a[i] = v
	}
	return BoolExpr{expr: bson.D{{Key: "$or", Value: a}}}
}

// Pow raises number to the specified exponent ($pow).
func Pow[T NumberResolver, U NumberResolver](number T, exponent U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$pow", Value: bson.A{number, exponent}}}}
}

// RadiansToDegrees converts a value from radians to degrees ($radiansToDegrees).
func RadiansToDegrees[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$radiansToDegrees", Value: expr}}}
}

// Range outputs an array of integers from start (inclusive) to end (exclusive) ($range).
// An optional step controls the increment; defaults to 1.
func Range[T NumberResolver, U NumberResolver, S NumberResolver](start T, end U, step *S) ArrayExpr {
	args := bson.A{start, end}
	if step != nil {
		args = append(args, *step)
	}
	return ArrayExpr{expr: bson.D{{Key: "$range", Value: args}}}
}

// ReverseArray returns an array with the elements in reverse order ($reverseArray).
func ReverseArray[T ArrayResolver](expr T) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$reverseArray", Value: expr}}}
}

// Round rounds a number to a whole integer or to a specified decimal place ($round).
// When place is omitted the array form is still used: [$number] (equivalent to place 0).
func Round[T NumberResolver](number T, place ...int) NumberExpr {
	if len(place) == 0 {
		return NumberExpr{expr: bson.D{{Key: "$round", Value: bson.A{number}}}}
	}
	return NumberExpr{expr: bson.D{{Key: "$round", Value: bson.A{number, place[0]}}}}
}

// SetDifference returns elements in the first set but not the second ($setDifference).
func SetDifference[T ArrayResolver, U ArrayResolver](expr1 T, expr2 U) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$setDifference", Value: bson.A{expr1, expr2}}}}
}

// SetEquals returns true if all input sets have the same distinct elements ($setEquals).
func SetEquals[T ArrayResolver](exprs ...T) BoolExpr {
	a := make(bson.A, len(exprs))
	for i, v := range exprs {
		a[i] = v
	}
	return BoolExpr{expr: bson.D{{Key: "$setEquals", Value: a}}}
}

// SetIntersection returns elements that appear in all of the input sets ($setIntersection).
func SetIntersection[T ArrayResolver](exprs ...T) ArrayExpr {
	a := make(bson.A, len(exprs))
	for i, v := range exprs {
		a[i] = v
	}
	return ArrayExpr{expr: bson.D{{Key: "$setIntersection", Value: a}}}
}

// SetIsSubset returns true if all elements of expr1 appear in expr2 ($setIsSubset).
func SetIsSubset[T ArrayResolver, U ArrayResolver](expr1 T, expr2 U) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$setIsSubset", Value: bson.A{expr1, expr2}}}}
}

// SetUnion returns elements that appear in any of the input sets ($setUnion).
func SetUnion[T ArrayResolver](exprs ...T) ArrayExpr {
	a := make(bson.A, len(exprs))
	for i, v := range exprs {
		a[i] = v
	}
	return ArrayExpr{expr: bson.D{{Key: "$setUnion", Value: a}}}
}

// Sigmoid returns 1 / (1 + e^(-x)) ($sigmoid).
func Sigmoid[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$sigmoid", Value: expr}}}
}

// Sin returns the sine of a value in radians ($sin).
func Sin[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$sin", Value: expr}}}
}

// Sinh returns the hyperbolic sine of a value in radians ($sinh).
func Sinh[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$sinh", Value: expr}}}
}

// Size returns the number of elements in the array ($size).
func Size[T ArrayResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$size", Value: expr}}}
}

// Slice returns n elements of an array ($slice).
// Pass a non-nil start to specify a starting index; otherwise elements are taken from the beginning.
func Slice[T ArrayResolver](expression T, n Expr, start Expr) ArrayExpr {
	if start != nil {
		return ArrayExpr{expr: bson.D{{Key: "$slice", Value: bson.A{expression, start, n}}}}
	}
	return ArrayExpr{expr: bson.D{{Key: "$slice", Value: bson.A{expression, n}}}}
}

// SortArray sorts the elements of an array by the specified document fields ($sortArray).
// Use SortArrayByValue for arrays of scalars.
func SortArray[T ArrayResolver](input T, sortBy ...SortField) ArrayExpr {
	sortDoc := make(bson.D, len(sortBy))
	for i, f := range sortBy {
		sf := f.sortField()
		sortDoc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return ArrayExpr{expr: bson.D{{Key: "$sortArray", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "sortBy", Value: sortDoc},
	}}}}
}

// SortArrayByValue sorts a scalar-element array in the given direction ($sortArray).
// Use SortArray for arrays of documents.
func SortArrayByValue[T ArrayResolver](input T, order SortOrder) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$sortArray", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "sortBy", Value: order.bsonValue()},
	}}}}
}

// Sqrt calculates the square root of a number ($sqrt).
func Sqrt[T NumberResolver](number T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$sqrt", Value: number}}}
}

// StdDevPop calculates the population standard deviation of numeric expressions ($stdDevPop).
func StdDevPop(exprs ...Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$stdDevPop", Value: exprs}}}
}

// StdDevSamp calculates the sample standard deviation of numeric expressions ($stdDevSamp).
func StdDevSamp(exprs ...Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$stdDevSamp", Value: exprs}}}
}

// Subtract returns a minus b ($subtract).
// TODO: $subtract also supports date-date → millis and date-millis → date;
// those variants are not yet modeled here.
func Subtract[T NumberResolver, U NumberResolver](a T, b U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$subtract", Value: bson.A{a, b}}}}
}

// Sum returns the sum of numeric expressions ($sum).
func Sum(exprs ...Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$sum", Value: exprs}}}
}

// Tan returns the tangent of a value in radians ($tan).
func Tan[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$tan", Value: expr}}}
}

// Tanh returns the hyperbolic tangent of a value in radians ($tanh).
func Tanh[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$tanh", Value: expr}}}
}

// Top returns the top element within an array according to the specified sort order ($top).
// This is the expression operator (MongoDB 7.0+) that takes an input array.
// See TopAccumulator for the $group/$setWindowFields accumulator form.
func Top[T ArrayResolver](input T, output Expr, sortBy ...SortField) AnyExpr {
	sortDoc := make(bson.D, len(sortBy))
	for i, f := range sortBy {
		sf := f.sortField()
		sortDoc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return AnyExpr{expr: bson.D{{Key: "$top", Value: bson.D{
		{Key: "sortBy", Value: sortDoc},
		{Key: "output", Value: output},
		{Key: "input", Value: input},
	}}}}
}

// TopN returns the top n elements within an array according to the specified sort order ($topN).
// This is the expression operator (MongoDB 7.0+) that takes an input array.
// See TopNAccumulator for the $group/$setWindowFields accumulator form.
func TopN[T ArrayResolver](n Expr, input T, output Expr, sortBy ...SortField) ArrayExpr {
	sortDoc := make(bson.D, len(sortBy))
	for i, f := range sortBy {
		sf := f.sortField()
		sortDoc[i] = bson.E{Key: sf.name, Value: sf.order.bsonValue()}
	}
	return ArrayExpr{expr: bson.D{{Key: "$topN", Value: bson.D{
		{Key: "n", Value: n},
		{Key: "sortBy", Value: sortDoc},
		{Key: "output", Value: output},
		{Key: "input", Value: input},
	}}}}
}

// Trunc truncates a number to a whole integer or to a specified decimal place ($trunc).
// When place is omitted the array form is still used: [$number] (equivalent to place 0).
func Trunc[T NumberResolver](number T, place ...int) NumberExpr {
	if len(place) == 0 {
		return NumberExpr{expr: bson.D{{Key: "$trunc", Value: bson.A{number}}}}
	}
	return NumberExpr{expr: bson.D{{Key: "$trunc", Value: bson.A{number, place[0]}}}}
}

// Zip merges arrays together into an array of arrays ($zip).
// When useLongestLength is true, the output length is determined by the longest input array;
// pass default values for shorter arrays via defaults. When useLongestLength is false (default),
// the output length is the shortest input array and defaults must be empty.
func Zip(inputs Expr, useLongestLength bool, defaults ...Expr) ArrayExpr {
	doc := bson.D{{Key: "inputs", Value: inputs}}
	if useLongestLength {
		doc = append(doc, bson.E{Key: "useLongestLength", Value: true})
	}
	if len(defaults) > 0 {
		doc = append(doc, bson.E{Key: "defaults", Value: defaults})
	}
	return ArrayExpr{expr: bson.D{{Key: "$zip", Value: doc}}}
}
