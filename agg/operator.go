package agg

import "go.mongodb.org/mongo-driver/v2/bson"

// --- arithmetic ---

// Add returns the sum of the given numeric expressions ($add).
// TODO: $add also accepts a date + milliseconds to produce a date;
// that variant is not yet modeled here.
func Add(vals ...NumericExpr) NumericExpr {
	return numericExprVal{v: bson.D{{Key: "$add", Value: exprSlice(vals)}}}
}

// Subtract returns a minus b ($subtract).
// TODO: $subtract also supports date-date → millis and date-millis → date;
// those variants are not yet modeled here.
func Subtract(a, b NumericExpr) NumericExpr {
	return numericExprVal{v: bson.D{{Key: "$subtract", Value: bson.A{a, b}}}}
}

// Multiply returns the product of the given numeric expressions ($multiply).
func Multiply(vals ...NumericExpr) NumericExpr {
	return numericExprVal{v: bson.D{{Key: "$multiply", Value: exprSlice(vals)}}}
}

// Divide returns a divided by b ($divide).
func Divide(a, b NumericExpr) NumericExpr {
	return numericExprVal{v: bson.D{{Key: "$divide", Value: bson.A{a, b}}}}
}

// --- comparison ---
// These accept Expr (not NumericExpr) because MongoDB allows comparing any
// two BSON types in an aggregation context.

// Eq returns true if a equals b ($eq).
func Eq(a, b Expr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$eq", Value: bson.A{a, b}}}}
}

// Ne returns true if a does not equal b ($ne).
func Ne(a, b Expr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$ne", Value: bson.A{a, b}}}}
}

// Gt returns true if a is greater than b ($gt).
func Gt(a, b Expr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$gt", Value: bson.A{a, b}}}}
}

// Gte returns true if a is greater than or equal to b ($gte).
func Gte(a, b Expr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$gte", Value: bson.A{a, b}}}}
}

// Lt returns true if a is less than b ($lt).
func Lt(a, b Expr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$lt", Value: bson.A{a, b}}}}
}

// Lte returns true if a is less than or equal to b ($lte).
func Lte(a, b Expr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$lte", Value: bson.A{a, b}}}}
}

// --- logical ---

// And returns true only when all expressions evaluate to true ($and).
func And(exprs ...BoolExpr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$and", Value: exprSlice(exprs)}}}
}

// Or returns true when any expression evaluates to true ($or).
func Or(exprs ...BoolExpr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$or", Value: exprSlice(exprs)}}}
}

// Not returns the boolean inverse of e ($not).
// $not takes a single-element array in the aggregation expression syntax.
func Not(e BoolExpr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$not", Value: bson.A{e}}}}
}

// --- string ---

// Concat concatenates the given string expressions ($concat).
func Concat(vals ...StringExpr) StringExpr {
	return stringExprVal{v: bson.D{{Key: "$concat", Value: exprSlice(vals)}}}
}

// --- array ---

// In returns true if expr is present in array ($in).
func In(expr Expr, array ArrayExpr) BoolExpr {
	return boolExprVal{v: bson.D{{Key: "$in", Value: bson.A{expr, array}}}}
}

// FilterArray selects elements of input for which cond evaluates to true ($filter).
// as names the variable for each element; pass "" to use the MongoDB default ("this").
// TODO: optional 'limit' parameter from the spec is not yet implemented.
func FilterArray(input ArrayExpr, cond BoolExpr, as string) ArrayExpr {
	args := bson.D{
		{Key: "input", Value: input},
		{Key: "cond", Value: cond},
	}
	if as != "" {
		args = append(args, bson.E{Key: "as", Value: as})
	}
	return arrayExprVal{v: bson.D{{Key: "$filter", Value: args}}}
}

// ArrayToObject converts an array of key-value pairs to a document ($arrayToObject).
func ArrayToObject(array ArrayExpr) ObjectExpr {
	return objectExprVal{v: bson.D{{Key: "$arrayToObject", Value: array}}}
}

// --- conditional / misc ---

// IfNull returns the first non-null expression in vals, or the last expression
// if all preceding ones are null ($ifNull).
func IfNull(vals ...Expr) Expr {
	return genericExprVal{v: bson.D{{Key: "$ifNull", Value: exprSlice(vals)}}}
}

// --- helpers ---

// exprSlice converts a typed slice of Expr sub-interface values to bson.A.
// Used so variadic operators produce a proper BSON array.
func exprSlice[T Expr](vals []T) bson.A {
	a := make(bson.A, len(vals))
	for i, v := range vals {
		a[i] = v
	}
	return a
}
