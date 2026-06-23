package agg

import "go.mongodb.org/mongo-driver/v2/bson"

// --- arithmetic ---

// Add returns the sum of the given numeric expressions ($add).
// TODO: $add also accepts a date + milliseconds to produce a date;
// that variant is not yet modeled here.
func Add[T NumberTypes, U NumberTypes](value T, values ...U) NumberExpr {
	v := make([]any, len(values)+1)
	v[0] = value
	for i := range values {
		v[i+1] = values[i]
	}
	return NumberExpr{
		expr: bson.D{{Key: "$add", Value: v}},
	}
}

// Subtract returns a minus b ($subtract).
// TODO: $subtract also supports date-date → millis and date-millis → date;
// those variants are not yet modeled here.
func Subtract[T NumberTypes, U NumberTypes](a T, b U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$subtract", Value: bson.A{a, b}}}}
}

// Multiply returns the product of the given numeric expressions ($multiply).
func Multiply[T NumberTypes, U NumberTypes](value T, values ...U) NumberExpr {
	v := make([]any, len(values)+1)
	v[0] = value
	for i := range values {
		v[i+1] = values[i]
	}
	return NumberExpr{expr: bson.D{{Key: "$multiply", Value: v}}}
}

// Divide returns a divided by b ($divide).
func Divide[T NumberTypes, U NumberTypes](a T, b U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$divide", Value: bson.A{a, b}}}}
}

// --- comparison ---
// These accept AnyExpr because MongoDB allows comparing any two BSON types
// in an aggregation context.

// Eq returns true if a equals b ($eq).
func Eq(a AnyExpr, b AnyExpr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$eq", Value: bson.A{a, b}}}}
}

// Ne returns true if a does not equal b ($ne).
func Ne(a AnyExpr, b AnyExpr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$ne", Value: bson.A{a, b}}}}
}

// Gt returns true if a is greater than b ($gt).
func Gt(a AnyExpr, b AnyExpr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$gt", Value: bson.A{a, b}}}}
}

// Gte returns true if a is greater than or equal to b ($gte).
func Gte(a AnyExpr, b AnyExpr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$gte", Value: bson.A{a, b}}}}
}

// Lt returns true if a is less than b ($lt).
func Lt(a AnyExpr, b AnyExpr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$lt", Value: bson.A{a, b}}}}
}

// Lte returns true if a is less than or equal to b ($lte).
func Lte(a AnyExpr, b AnyExpr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$lte", Value: bson.A{a, b}}}}
}

// --- logical ---

// And returns true only when all expressions evaluate to true ($and).
func And[T BoolTypes](exprs ...T) BoolExpr {
	a := make(bson.A, len(exprs))
	for i, v := range exprs {
		a[i] = v
	}
	return BoolExpr{expr: bson.D{{Key: "$and", Value: a}}}
}

// Or returns true when any expression evaluates to true ($or).
func Or[T BoolTypes](exprs ...T) BoolExpr {
	a := make(bson.A, len(exprs))
	for i, v := range exprs {
		a[i] = v
	}
	return BoolExpr{expr: bson.D{{Key: "$or", Value: a}}}
}

// Not returns the boolean inverse of e ($not).
// $not takes a single-element array in the aggregation expression syntax.
func Not[T BoolTypes](e T) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$not", Value: bson.A{e}}}}
}

// --- string ---

// Concat concatenates the given string expressions ($concat).
func Concat[T StringTypes, U StringTypes](value T, values ...U) StringExpr {
	v := make([]any, len(values)+1)
	v[0] = value
	for i := range values {
		v[i+1] = values[i]
	}
	return StringExpr{expr: bson.D{{Key: "$concat", Value: v}}}
}

// --- array ---

// In returns true if expr is present in array ($in).
func In[U ArrayTypes](expr AnyExpr, array U) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$in", Value: bson.A{expr, array}}}}
}

// FilterArray selects elements of input for which cond evaluates to true ($filter).
// as names the variable for each element; pass "" to use the MongoDB default ("this").
func FilterArray[T ArrayTypes, U BoolTypes](input T, cond U, as string, limit ...NumberExpr) ArrayExpr {
	args := bson.D{
		{Key: "input", Value: input},
		{Key: "cond", Value: cond},
	}
	if as != "" {
		args = append(args, bson.E{Key: "as", Value: as})
	}
	if len(limit) > 1 {
		panic("FilterArray: at most one limit expression may be provided")
	}
	if len(limit) == 1 {
		args = append(args, bson.E{Key: "limit", Value: limit[0]})
	}
	return ArrayExpr{expr: bson.D{{Key: "$filter", Value: args}}}
}

// ArrayToObject converts an array of key-value pairs to a document ($arrayToObject).
func ArrayToObject[T ArrayTypes](array T) ObjectExpr {
	return ObjectExpr{expr: bson.D{{Key: "$arrayToObject", Value: array}}}
}

// --- conditional / misc ---

// IfNull returns the first non-null expression among val and more, or fallback
// if all preceding ones are null ($ifNull).
func IfNull(val AnyExpr, fallback AnyExpr, more ...AnyExpr) AnyExpr {
	v := make([]any, len(more)+2)
	v[0] = val
	v[1] = fallback
	for i := range more {
		v[i+2] = more[i]
	}
	return AnyExpr{expr: bson.D{{Key: "$ifNull", Value: v}}}
}

