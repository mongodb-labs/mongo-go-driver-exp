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

// And returns true only when all expressions evaluate to true ($and).
func And[T BoolResolver](exprs ...T) BoolExpr {
	a := make(bson.A, len(exprs))
	for i, v := range exprs {
		a[i] = v
	}
	return BoolExpr{expr: bson.D{{Key: "$and", Value: a}}}
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
func FilterArray[T ArrayResolver, U BoolResolver](input T, cond U, as string, limit ...NumberExpr) ArrayExpr {
	args := bson.D{
		{Key: "input", Value: input},
		{Key: "cond", Value: cond},
	}
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

// Round rounds a number to a whole integer or to a specified decimal place ($round).
// When place is omitted the array form is still used: [$number] (equivalent to place 0).
func Round[T NumberResolver](number T, place ...int) NumberExpr {
	if len(place) == 0 {
		return NumberExpr{expr: bson.D{{Key: "$round", Value: bson.A{number}}}}
	}
	return NumberExpr{expr: bson.D{{Key: "$round", Value: bson.A{number, place[0]}}}}
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

// Trunc truncates a number to a whole integer or to a specified decimal place ($trunc).
// When place is omitted the array form is still used: [$number] (equivalent to place 0).
func Trunc[T NumberResolver](number T, place ...int) NumberExpr {
	if len(place) == 0 {
		return NumberExpr{expr: bson.D{{Key: "$trunc", Value: bson.A{number}}}}
	}
	return NumberExpr{expr: bson.D{{Key: "$trunc", Value: bson.A{number, place[0]}}}}
}
