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

// BinarySize returns the size of a string or binary data value's content in bytes ($binarySize).
// The argument must resolve to a string, binary data, or null.
func BinarySize(expr Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$binarySize", Value: expr}}}
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

// BsonSize returns the size in bytes of a document when encoded as BSON ($bsonSize).
// The argument must resolve to an object or null.
func BsonSize[T ObjectResolver](object T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$bsonSize", Value: object}}}
}

// SwitchCase is a single branch of a Switch expression, pairing a boolean
// condition with the value to return when that condition is the first to hold.
type SwitchCase interface{ switchCase() switchCase }

type switchCase struct {
	cond Expr
	then Expr
}

func (c switchCase) switchCase() switchCase { return c }

// Case constructs a single Switch branch that returns then when cond is true ($switch branch).
func Case[T BoolResolver](cond T, then Expr) SwitchCase {
	return switchCase{cond: cond, then: then}
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
// Accepts any expression type (Expr = any) so that string literals, field paths,
// and nested expressions can be mixed in a single call.
func Concat(value Expr, values ...Expr) StringExpr {
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

// Cond evaluates ifExpr and returns then when it is true, otherwise elseExpr ($cond).
func Cond[T BoolResolver](ifExpr T, then Expr, elseExpr Expr) AnyExpr {
	return AnyExpr{expr: bson.D{{Key: "$cond", Value: bson.D{
		{Key: "if", Value: ifExpr},
		{Key: "then", Value: then},
		{Key: "else", Value: elseExpr},
	}}}}
}

type convertOptions struct {
	onError    any
	hasOnError bool
	onNull     any
	hasOnNull  bool
	base       any
}

func WithConvertOnError(onError Expr) Option[convertOptions] {
	return func(o *convertOptions) {
		o.onError = onError
		o.hasOnError = true
	}
}

func WithConvertOnNull(onNull Expr) Option[convertOptions] {
	return func(o *convertOptions) {
		o.onNull = onNull
		o.hasOnNull = true
	}
}

func WithConvertBase(base int32) Option[convertOptions] {
	return func(o *convertOptions) {
		o.base = base
	}
}

// Convert converts the input value to the type named by to ($convert).
// to must resolve to a string type name or numeric type code. Optionally provide a value
// returned on a conversion error via WithConvertOnError, a value returned when the input is
// null or missing via WithConvertOnNull, and the numeric base (2, 8, 10, or 16) used when
// converting between strings and integers via WithConvertBase; base defaults to 10 (MongoDB 8.3+).
func Convert[T StringResolver | NumberResolver](input Expr, to T, opts ...Option[convertOptions]) AnyExpr {
	var o convertOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{
		{Key: "input", Value: input},
		{Key: "to", Value: to},
	}
	if o.hasOnError {
		doc = append(doc, bson.E{Key: "onError", Value: o.onError})
	}
	if o.hasOnNull {
		doc = append(doc, bson.E{Key: "onNull", Value: o.onNull})
	}
	if o.base != nil {
		doc = append(doc, bson.E{Key: "base", Value: o.base})
	}
	return AnyExpr{expr: bson.D{{Key: "$convert", Value: doc}}}
}

// Cos returns the cosine of a value in radians ($cos).
func Cos[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$cos", Value: expr}}}
}

// Cosh returns the hyperbolic cosine of a value in radians ($cosh).
func Cosh[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$cosh", Value: expr}}}
}

// CreateObjectId returns a new randomly generated ObjectId ($createObjectId).
func CreateObjectId() AnyExpr {
	return AnyExpr{expr: bson.D{{Key: "$createObjectId", Value: bson.D{}}}}
}

type dateAddOptions struct {
	timezone any
}

// WithDateAddTimezone sets the timezone for DateAdd. The timezone resolves to an
// Olson timezone identifier or UTC offset string.
func WithDateAddTimezone[T StringResolver](timezone T) Option[dateAddOptions] {
	return func(o *dateAddOptions) { o.timezone = timezone }
}

// DateAdd adds a number of time units to a date ($dateAdd).
// unit is a time unit such as "day" or "hour"; amount resolves to an int or long.
func DateAdd[T DateResolver | TimestampResolver | ObjectIDResolver, A NumberResolver](startDate T, unit string, amount A, opts ...Option[dateAddOptions]) DateExpr {
	var o dateAddOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{
		{Key: "startDate", Value: startDate},
		{Key: "unit", Value: unit},
		{Key: "amount", Value: amount},
	}
	if o.timezone != nil {
		doc = append(doc, bson.E{Key: "timezone", Value: o.timezone})
	}
	return DateExpr{expr: bson.D{{Key: "$dateAdd", Value: doc}}}
}

type dateDiffOptions struct {
	timezone    any
	startOfWeek any
}

// WithDateDiffTimezone sets the timezone for DateDiff.
func WithDateDiffTimezone[T StringResolver](timezone T) Option[dateDiffOptions] {
	return func(o *dateDiffOptions) { o.timezone = timezone }
}

// WithDateDiffStartOfWeek sets the start of the week for DateDiff; used when unit is "week".
func WithDateDiffStartOfWeek[T StringResolver](startOfWeek T) Option[dateDiffOptions] {
	return func(o *dateDiffOptions) { o.startOfWeek = startOfWeek }
}

// DateDiff returns the difference between two dates measured in unit ($dateDiff).
func DateDiff[S DateResolver | TimestampResolver | ObjectIDResolver, E DateResolver | TimestampResolver | ObjectIDResolver](startDate S, endDate E, unit string, opts ...Option[dateDiffOptions]) NumberExpr {
	var o dateDiffOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{
		{Key: "startDate", Value: startDate},
		{Key: "endDate", Value: endDate},
		{Key: "unit", Value: unit},
	}
	if o.timezone != nil {
		doc = append(doc, bson.E{Key: "timezone", Value: o.timezone})
	}
	if o.startOfWeek != nil {
		doc = append(doc, bson.E{Key: "startOfWeek", Value: o.startOfWeek})
	}
	return NumberExpr{expr: bson.D{{Key: "$dateDiff", Value: doc}}}
}

type dateFromPartsOptions struct {
	year         any
	isoWeekYear  any
	month        any
	isoWeek      any
	day          any
	isoDayOfWeek any
	hour         any
	minute       any
	second       any
	millisecond  any
	timezone     any
}

// WithDateFromPartsYear sets the calendar year.
func WithDateFromPartsYear[T NumberResolver](year T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.year = year }
}

// WithDateFromPartsIsoWeekYear sets the ISO week date year.
func WithDateFromPartsIsoWeekYear[T NumberResolver](isoWeekYear T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.isoWeekYear = isoWeekYear }
}

// WithDateFromPartsMonth sets the month.
func WithDateFromPartsMonth[T NumberResolver](month T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.month = month }
}

// WithDateFromPartsIsoWeek sets the ISO week of the year.
func WithDateFromPartsIsoWeek[T NumberResolver](isoWeek T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.isoWeek = isoWeek }
}

// WithDateFromPartsDay sets the day of the month.
func WithDateFromPartsDay[T NumberResolver](day T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.day = day }
}

// WithDateFromPartsIsoDayOfWeek sets the ISO day of the week (Monday 1 - Sunday 7).
func WithDateFromPartsIsoDayOfWeek[T NumberResolver](isoDayOfWeek T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.isoDayOfWeek = isoDayOfWeek }
}

// WithDateFromPartsHour sets the hour.
func WithDateFromPartsHour[T NumberResolver](hour T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.hour = hour }
}

// WithDateFromPartsMinute sets the minute.
func WithDateFromPartsMinute[T NumberResolver](minute T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.minute = minute }
}

// WithDateFromPartsSecond sets the second.
func WithDateFromPartsSecond[T NumberResolver](second T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.second = second }
}

// WithDateFromPartsMillisecond sets the millisecond.
func WithDateFromPartsMillisecond[T NumberResolver](millisecond T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.millisecond = millisecond }
}

// WithDateFromPartsTimezone sets the timezone.
func WithDateFromPartsTimezone[T StringResolver](timezone T) Option[dateFromPartsOptions] {
	return func(o *dateFromPartsOptions) { o.timezone = timezone }
}

// DateFromParts constructs a Date from its constituent parts ($dateFromParts).
// Provide either the calendar parts (year, month, day) or the ISO week date
// parts (isoWeekYear, isoWeek, isoDayOfWeek), not both.
func DateFromParts(opts ...Option[dateFromPartsOptions]) DateExpr {
	var o dateFromPartsOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{}
	if o.year != nil {
		doc = append(doc, bson.E{Key: "year", Value: o.year})
	}
	if o.isoWeekYear != nil {
		doc = append(doc, bson.E{Key: "isoWeekYear", Value: o.isoWeekYear})
	}
	if o.month != nil {
		doc = append(doc, bson.E{Key: "month", Value: o.month})
	}
	if o.isoWeek != nil {
		doc = append(doc, bson.E{Key: "isoWeek", Value: o.isoWeek})
	}
	if o.day != nil {
		doc = append(doc, bson.E{Key: "day", Value: o.day})
	}
	if o.isoDayOfWeek != nil {
		doc = append(doc, bson.E{Key: "isoDayOfWeek", Value: o.isoDayOfWeek})
	}
	if o.hour != nil {
		doc = append(doc, bson.E{Key: "hour", Value: o.hour})
	}
	if o.minute != nil {
		doc = append(doc, bson.E{Key: "minute", Value: o.minute})
	}
	if o.second != nil {
		doc = append(doc, bson.E{Key: "second", Value: o.second})
	}
	if o.millisecond != nil {
		doc = append(doc, bson.E{Key: "millisecond", Value: o.millisecond})
	}
	if o.timezone != nil {
		doc = append(doc, bson.E{Key: "timezone", Value: o.timezone})
	}
	return DateExpr{expr: bson.D{{Key: "$dateFromParts", Value: doc}}}
}

type dateFromStringOptions struct {
	format     any
	timezone   any
	onError    any
	hasOnError bool
	onNull     any
	hasOnNull  bool
}

// WithDateFromStringFormat sets the date format of the input string.
func WithDateFromStringFormat[T StringResolver](format T) Option[dateFromStringOptions] {
	return func(o *dateFromStringOptions) { o.format = format }
}

// WithDateFromStringTimezone sets the timezone used to parse the string.
func WithDateFromStringTimezone[T StringResolver](timezone T) Option[dateFromStringOptions] {
	return func(o *dateFromStringOptions) { o.timezone = timezone }
}

// WithDateFromStringOnError sets the value returned if the string cannot be parsed.
func WithDateFromStringOnError(onError Expr) Option[dateFromStringOptions] {
	return func(o *dateFromStringOptions) { o.onError = onError; o.hasOnError = true }
}

// WithDateFromStringOnNull sets the value returned if the string is null or missing.
func WithDateFromStringOnNull(onNull Expr) Option[dateFromStringOptions] {
	return func(o *dateFromStringOptions) { o.onNull = onNull; o.hasOnNull = true }
}

// DateFromString converts a date/time string to a Date ($dateFromString).
func DateFromString[S StringResolver](dateString S, opts ...Option[dateFromStringOptions]) DateExpr {
	var o dateFromStringOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "dateString", Value: dateString}}
	if o.format != nil {
		doc = append(doc, bson.E{Key: "format", Value: o.format})
	}
	if o.timezone != nil {
		doc = append(doc, bson.E{Key: "timezone", Value: o.timezone})
	}
	if o.hasOnError {
		doc = append(doc, bson.E{Key: "onError", Value: o.onError})
	}
	if o.hasOnNull {
		doc = append(doc, bson.E{Key: "onNull", Value: o.onNull})
	}
	return DateExpr{expr: bson.D{{Key: "$dateFromString", Value: doc}}}
}

type dateSubtractOptions struct {
	timezone any
}

// WithDateSubtractTimezone sets the timezone for DateSubtract.
func WithDateSubtractTimezone[T StringResolver](timezone T) Option[dateSubtractOptions] {
	return func(o *dateSubtractOptions) { o.timezone = timezone }
}

// DateSubtract subtracts a number of time units from a date ($dateSubtract).
// unit is a time unit such as "day" or "hour"; amount resolves to an int or long.
func DateSubtract[T DateResolver | TimestampResolver | ObjectIDResolver, A NumberResolver](startDate T, unit string, amount A, opts ...Option[dateSubtractOptions]) DateExpr {
	var o dateSubtractOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{
		{Key: "startDate", Value: startDate},
		{Key: "unit", Value: unit},
		{Key: "amount", Value: amount},
	}
	if o.timezone != nil {
		doc = append(doc, bson.E{Key: "timezone", Value: o.timezone})
	}
	return DateExpr{expr: bson.D{{Key: "$dateSubtract", Value: doc}}}
}

type dateToPartsOptions struct {
	timezone any
	iso8601  *bool
}

// WithDateToPartsTimezone sets the timezone for DateToParts.
func WithDateToPartsTimezone[T StringResolver](timezone T) Option[dateToPartsOptions] {
	return func(o *dateToPartsOptions) { o.timezone = timezone }
}

// WithDateToPartsIso8601 selects ISO week date fields in the output document.
func WithDateToPartsIso8601(iso8601 bool) Option[dateToPartsOptions] {
	return func(o *dateToPartsOptions) { o.iso8601 = &iso8601 }
}

// DateToParts returns a document containing the constituent parts of a date ($dateToParts).
func DateToParts[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[dateToPartsOptions]) ObjectExpr {
	var o dateToPartsOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "date", Value: date}}
	if o.timezone != nil {
		doc = append(doc, bson.E{Key: "timezone", Value: o.timezone})
	}
	if o.iso8601 != nil {
		doc = append(doc, bson.E{Key: "iso8601", Value: *o.iso8601})
	}
	return ObjectExpr{expr: bson.D{{Key: "$dateToParts", Value: doc}}}
}

type dateToStringOptions struct {
	format    any
	timezone  any
	onNull    any
	hasOnNull bool
}

// WithDateToStringFormat sets the output date format.
func WithDateToStringFormat[T StringResolver](format T) Option[dateToStringOptions] {
	return func(o *dateToStringOptions) { o.format = format }
}

// WithDateToStringTimezone sets the timezone used to format the date.
func WithDateToStringTimezone[T StringResolver](timezone T) Option[dateToStringOptions] {
	return func(o *dateToStringOptions) { o.timezone = timezone }
}

// WithDateToStringOnNull sets the value returned if the date is null or missing.
func WithDateToStringOnNull(onNull Expr) Option[dateToStringOptions] {
	return func(o *dateToStringOptions) { o.onNull = onNull; o.hasOnNull = true }
}

// DateToString returns the date as a formatted string ($dateToString).
func DateToString[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[dateToStringOptions]) StringExpr {
	var o dateToStringOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "date", Value: date}}
	if o.format != nil {
		doc = append(doc, bson.E{Key: "format", Value: o.format})
	}
	if o.timezone != nil {
		doc = append(doc, bson.E{Key: "timezone", Value: o.timezone})
	}
	if o.hasOnNull {
		doc = append(doc, bson.E{Key: "onNull", Value: o.onNull})
	}
	return StringExpr{expr: bson.D{{Key: "$dateToString", Value: doc}}}
}

type dateTruncOptions struct {
	binSize     any
	timezone    any
	startOfWeek any
}

// WithDateTruncBinSize sets the number of units per bin; defaults to 1.
func WithDateTruncBinSize[T NumberResolver](binSize T) Option[dateTruncOptions] {
	return func(o *dateTruncOptions) { o.binSize = binSize }
}

// WithDateTruncTimezone sets the timezone for DateTrunc.
func WithDateTruncTimezone[T StringResolver](timezone T) Option[dateTruncOptions] {
	return func(o *dateTruncOptions) { o.timezone = timezone }
}

// WithDateTruncStartOfWeek sets the start of the week; used when unit is "week".
func WithDateTruncStartOfWeek[T StringResolver](startOfWeek T) Option[dateTruncOptions] {
	return func(o *dateTruncOptions) { o.startOfWeek = startOfWeek }
}

// DateTrunc truncates a date to the given unit ($dateTrunc).
func DateTrunc[T DateResolver | TimestampResolver | ObjectIDResolver](date T, unit string, opts ...Option[dateTruncOptions]) DateExpr {
	var o dateTruncOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{
		{Key: "date", Value: date},
		{Key: "unit", Value: unit},
	}
	if o.binSize != nil {
		doc = append(doc, bson.E{Key: "binSize", Value: o.binSize})
	}
	if o.timezone != nil {
		doc = append(doc, bson.E{Key: "timezone", Value: o.timezone})
	}
	if o.startOfWeek != nil {
		doc = append(doc, bson.E{Key: "startOfWeek", Value: o.startOfWeek})
	}
	return DateExpr{expr: bson.D{{Key: "$dateTrunc", Value: doc}}}
}

type datePartOptions struct {
	timezone any
}

// WithDatePartTimezone sets the timezone for the single-date extraction operators
// (DayOfMonth, DayOfWeek, DayOfYear, Hour, IsoDayOfWeek, IsoWeek, IsoWeekYear,
// Millisecond, Minute, Month, Second, Week, Year). It resolves to an Olson
// timezone identifier or UTC offset string.
func WithDatePartTimezone[T StringResolver](timezone T) Option[datePartOptions] {
	return func(o *datePartOptions) { o.timezone = timezone }
}

// datePart builds the verbose {op: {date, timezone?}} form shared by the
// single-date extraction operators, all of which return an int.
func datePart(op string, date Expr, opts ...Option[datePartOptions]) NumberExpr {
	var o datePartOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "date", Value: date}}
	if o.timezone != nil {
		doc = append(doc, bson.E{Key: "timezone", Value: o.timezone})
	}
	return NumberExpr{expr: bson.D{{Key: op, Value: doc}}}
}

// DayOfMonth returns the day of the month for a date, from 1 to 31 ($dayOfMonth).
func DayOfMonth[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$dayOfMonth", date, opts...)
}

// DayOfWeek returns the day of the week for a date, from 1 (Sunday) to 7 (Saturday) ($dayOfWeek).
func DayOfWeek[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$dayOfWeek", date, opts...)
}

// DayOfYear returns the day of the year for a date, from 1 to 366 ($dayOfYear).
func DayOfYear[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$dayOfYear", date, opts...)
}

// DegreesToRadians converts a value from degrees to radians ($degreesToRadians).
func DegreesToRadians[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$degreesToRadians", Value: expr}}}
}

type deserializeEJSONOptions struct {
	onError    any
	hasOnError bool
}

// WithDeserializeEJSONOnError sets the value returned if $deserializeEJSON encounters a conversion error.
func WithDeserializeEJSONOnError(onError Expr) Option[deserializeEJSONOptions] {
	return func(o *deserializeEJSONOptions) {
		o.onError = onError
		o.hasOnError = true
	}
}

// DeserializeEJSON converts an Extended JSON document into native BSON values ($deserializeEJSON).
// Optionally provide a fallback value via WithDeserializeEJSONOnError.
func DeserializeEJSON(input Expr, opts ...Option[deserializeEJSONOptions]) ObjectExpr {
	var o deserializeEJSONOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "input", Value: input}}
	if o.hasOnError {
		doc = append(doc, bson.E{Key: "onError", Value: o.onError})
	}
	return ObjectExpr{expr: bson.D{{Key: "$deserializeEJSON", Value: doc}}}
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

type filterOptions struct {
	as           any
	arrayIndexAs any
	limit        any
}

func WithFilterAs(as string) Option[filterOptions] {
	return func(o *filterOptions) {
		o.as = as
	}
}

// WithFilterArrayIndexAs names the variable that represents the index of the current element (MongoDB 8.3+).
func WithFilterArrayIndexAs(arrayIndexAs string) Option[filterOptions] {
	return func(o *filterOptions) {
		o.arrayIndexAs = arrayIndexAs
	}
}

func WithFilterLimit(limit Expr) Option[filterOptions] {
	return func(o *filterOptions) {
		o.limit = limit
	}
}

// Filter selects elements of input for which cond evaluates to true ($filter).
// Optionally name the per-element variable via WithFilterAs (defaults to "this") and the element
// index variable via WithFilterArrayIndexAs, and cap the number of matching elements via WithFilterLimit.
func Filter[T ArrayResolver, U BoolResolver](input T, cond U, opts ...Option[filterOptions]) ArrayExpr {
	var o filterOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.D{bson.E{Key: "input", Value: input}}
	if o.as != nil {
		args = append(args, bson.E{Key: "as", Value: o.as})
	}
	if o.arrayIndexAs != nil {
		args = append(args, bson.E{Key: "arrayIndexAs", Value: o.arrayIndexAs})
	}
	args = append(args, bson.E{Key: "cond", Value: cond})
	if o.limit != nil {
		args = append(args, bson.E{Key: "limit", Value: o.limit})
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

type functionOptions struct {
	lang    string
	hasLang bool
}

// WithFunctionLang sets the language of the $function body; defaults to "js".
func WithFunctionLang(lang string) Option[functionOptions] {
	return func(o *functionOptions) {
		o.lang = lang
		o.hasLang = true
	}
}

// Function defines and invokes a custom function ($function).
// args are passed to the function body in order; the language defaults to "js"
// unless overridden via WithFunctionLang.
func Function(body string, args []Expr, opts ...Option[functionOptions]) AnyExpr {
	var o functionOptions
	for _, opt := range opts {
		opt(&o)
	}
	lang := "js"
	if o.hasLang {
		lang = o.lang
	}
	if args == nil {
		args = []Expr{}
	}
	return AnyExpr{expr: bson.D{{Key: "$function", Value: bson.D{
		{Key: "body", Value: bson.JavaScript(body)},
		{Key: "args", Value: args},
		{Key: "lang", Value: lang},
	}}}}
}

type getFieldOptions struct {
	input any
}

func WithGetFieldInput(input Expr) Option[getFieldOptions] {
	return func(o *getFieldOptions) {
		o.input = input
	}
}

// GetField returns the value of the specified field from a document ($getField).
// field must resolve to a string constant. Optionally provide the document to read from via
// WithGetFieldInput; it defaults to the document currently being processed ($$CURRENT).
func GetField[T StringResolver](field T, opts ...Option[getFieldOptions]) AnyExpr {
	var o getFieldOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "field", Value: field}}
	if o.input != nil {
		doc = append(doc, bson.E{Key: "input", Value: o.input})
	}
	return AnyExpr{expr: bson.D{{Key: "$getField", Value: doc}}}
}

// Gt returns true if a is greater than b ($gt).
func Gt(a Expr, b Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$gt", Value: bson.A{a, b}}}}
}

// Gte returns true if a is greater than or equal to b ($gte).
func Gte(a Expr, b Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$gte", Value: bson.A{a, b}}}}
}

// Hash generates a binary hash (BinData) of a string or binary input using the named algorithm ($hash).
func Hash(input Expr, algorithm string) AnyExpr {
	return AnyExpr{expr: bson.D{{Key: "$hash", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "algorithm", Value: algorithm},
	}}}}
}

// HexHash generates an uppercase hexadecimal string hash of a string or binary input using the named algorithm ($hexHash).
func HexHash(input Expr, algorithm string) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$hexHash", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "algorithm", Value: algorithm},
	}}}}
}

// Hour returns the hour for a date, from 0 to 23 ($hour).
func Hour[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$hour", date, opts...)
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

type indexOfOptions struct {
	start any
	end   any
}

func WithIndexOfStart[T NumberResolver](start T) Option[indexOfOptions] {
	return func(o *indexOfOptions) {
		o.start = start
	}
}

func WithIndexOfEnd[T NumberResolver](end T) Option[indexOfOptions] {
	return func(o *indexOfOptions) {
		o.end = end
	}
}

// IndexOfArray searches an array for a value and returns the index of the first occurrence ($indexOfArray).
// Optionally bound the search via WithIndexOfStart and WithIndexOfEnd; if only the
// end is given, the start defaults to 0.
func IndexOfArray[T ArrayResolver](array T, search Expr, opts ...Option[indexOfOptions]) NumberExpr {
	var o indexOfOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.A{array, search}
	if o.start != nil || o.end != nil {
		if o.start != nil {
			args = append(args, o.start)
		} else {
			args = append(args, 0) // end requires start; default to 0
		}
	}
	if o.end != nil {
		args = append(args, o.end)
	}
	return NumberExpr{expr: bson.D{{Key: "$indexOfArray", Value: args}}}
}

// IndexOfBytes searches a string for a substring and returns the UTF-8 byte index of the first occurrence, or -1 if not found ($indexOfBytes).
// Optionally bound the search via WithIndexOfStart and WithIndexOfEnd; if only the
// end is given, the start defaults to 0.
func IndexOfBytes[T StringResolver, U StringResolver](str T, substring U, opts ...Option[indexOfOptions]) NumberExpr {
	var o indexOfOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.A{str, substring}
	if o.start != nil || o.end != nil {
		if o.start != nil {
			args = append(args, o.start)
		} else {
			args = append(args, 0)
		}
	}
	if o.end != nil {
		args = append(args, o.end)
	}
	return NumberExpr{expr: bson.D{{Key: "$indexOfBytes", Value: args}}}
}

// IndexOfCP searches a string for a substring and returns the UTF-8 code point index of the first occurrence, or -1 if not found ($indexOfCP).
// Optionally bound the search via WithIndexOfStart and WithIndexOfEnd; if only the
// end is given, the start defaults to 0.
func IndexOfCP[T StringResolver, U StringResolver](str T, substring U, opts ...Option[indexOfOptions]) NumberExpr {
	var o indexOfOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.A{str, substring}
	if o.start != nil || o.end != nil {
		if o.start != nil {
			args = append(args, o.start)
		} else {
			args = append(args, 0)
		}
	}
	if o.end != nil {
		args = append(args, o.end)
	}
	return NumberExpr{expr: bson.D{{Key: "$indexOfCP", Value: args}}}
}

// IsArray returns true if the operand resolves to an array ($isArray).
func IsArray(expr Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$isArray", Value: bson.A{expr}}}}
}

// IsNumber returns true if the expression resolves to an integer, decimal, double, or long ($isNumber).
func IsNumber(expr Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$isNumber", Value: expr}}}
}

// IsoDayOfWeek returns the ISO 8601 weekday number, from 1 (Monday) to 7 (Sunday) ($isoDayOfWeek).
func IsoDayOfWeek[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$isoDayOfWeek", date, opts...)
}

// IsoWeek returns the ISO 8601 week number, from 1 to 53 ($isoWeek).
func IsoWeek[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$isoWeek", date, opts...)
}

// IsoWeekYear returns the year number in ISO 8601 format ($isoWeekYear).
func IsoWeekYear[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$isoWeekYear", date, opts...)
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

// Let binds variables for use within the scope of the in subexpression and returns its result ($let).
// Build the variable bindings with Assign.
func Let(vars []SetField, in Expr) AnyExpr {
	varsDoc := make(bson.D, len(vars))
	for i, v := range vars {
		sf := v.setField()
		varsDoc[i] = bson.E{Key: sf.name, Value: sf.expr}
	}
	return AnyExpr{expr: bson.D{{Key: "$let", Value: bson.D{
		{Key: "vars", Value: varsDoc},
		{Key: "in", Value: in},
	}}}}
}

// Literal returns a value without parsing it as an expression ($literal).
func Literal(value Expr) AnyExpr {
	return AnyExpr{expr: bson.D{{Key: "$literal", Value: value}}}
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

// Ltrim removes whitespace or the specified characters from the beginning of a string ($ltrim).
// Optionally specify the characters to remove via WithTrimChars; when omitted, whitespace is removed.
func Ltrim[T StringResolver](input T, opts ...Option[trimOptions]) StringExpr {
	var o trimOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.D{{Key: "input", Value: input}}
	if o.chars != nil {
		args = append(args, bson.E{Key: "chars", Value: o.chars})
	}
	return StringExpr{expr: bson.D{{Key: "$ltrim", Value: args}}}
}

type mapOptions struct {
	as           any
	arrayIndexAs any
}

// WithMapAs names the variable that represents each element of the input array; defaults to "this".
func WithMapAs(as string) Option[mapOptions] {
	return func(o *mapOptions) {
		o.as = as
	}
}

// WithMapArrayIndexAs names the variable that represents the index of the current element (MongoDB 8.3+).
func WithMapArrayIndexAs(arrayIndexAs string) Option[mapOptions] {
	return func(o *mapOptions) {
		o.arrayIndexAs = arrayIndexAs
	}
}

// Map applies the in subexpression to each element of the input array and returns the resulting array ($map).
// Optionally name the per-element variable via WithMapAs (defaults to "this") and the element
// index variable via WithMapArrayIndexAs.
func Map[T ArrayResolver](input T, in Expr, opts ...Option[mapOptions]) ArrayExpr {
	var o mapOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.D{{Key: "input", Value: input}}
	if o.as != nil {
		args = append(args, bson.E{Key: "as", Value: o.as})
	}
	if o.arrayIndexAs != nil {
		args = append(args, bson.E{Key: "arrayIndexAs", Value: o.arrayIndexAs})
	}
	args = append(args, bson.E{Key: "in", Value: in})
	return ArrayExpr{expr: bson.D{{Key: "$map", Value: args}}}
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

// Median returns an approximation of the median (50th percentile) as a scalar ($median).
func Median(input Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$median", Value: bson.D{
		{Key: "input", Value: input},
		// Currently the method must always be "approximate". If this changes, we need to
		// add an argument for this.
		{Key: "method", Value: "approximate"},
	}}}}
}

// MergeObjects combines multiple documents into a single document ($mergeObjects).
// Each argument must resolve to a document.
func MergeObjects(documents ...Expr) ObjectExpr {
	return ObjectExpr{expr: bson.D{{Key: "$mergeObjects", Value: documents}}}
}

// Meta accesses per-document metadata (such as "textScore" or "indexKey") for the aggregation ($meta).
func Meta(keyword string) AnyExpr {
	return AnyExpr{expr: bson.D{{Key: "$meta", Value: keyword}}}
}

// Millisecond returns the millisecond portion of a date, from 0 to 999 ($millisecond).
func Millisecond[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$millisecond", date, opts...)
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

// Minute returns the minute for a date, from 0 to 59 ($minute).
func Minute[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$minute", date, opts...)
}

// Mod returns the remainder of dividing dividend by divisor ($mod).
func Mod[T NumberResolver, U NumberResolver](dividend T, divisor U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$mod", Value: bson.A{dividend, divisor}}}}
}

// Month returns the month for a date, from 1 (January) to 12 (December) ($month).
func Month[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$month", date, opts...)
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

// ObjectToArray converts a document to an array of key-value pair documents ($objectToArray).
// The argument must resolve to a document; only top-level fields are converted.
func ObjectToArray[T ObjectResolver](object T) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$objectToArray", Value: object}}}
}

// Or returns true when any expression evaluates to true ($or).
func Or[T BoolResolver](exprs ...T) BoolExpr {
	a := make(bson.A, len(exprs))
	for i, v := range exprs {
		a[i] = v
	}
	return BoolExpr{expr: bson.D{{Key: "$or", Value: a}}}
}

// Percentile returns an array of scalar values corresponding to the requested percentiles ($percentile).
// p lists the percentiles to compute (each in the range 0.0 to 1.0).
func Percentile(input Expr, p Expr) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$percentile", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "p", Value: p},
		// Currently the method must always be "approximate". If this changes, we need to
		// add an argument for this.
		{Key: "method", Value: "approximate"},
	}}}}
}

// Pow raises number to the specified exponent ($pow).
func Pow[T NumberResolver, U NumberResolver](number T, exponent U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$pow", Value: bson.A{number, exponent}}}}
}

// RadiansToDegrees converts a value from radians to degrees ($radiansToDegrees).
func RadiansToDegrees[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$radiansToDegrees", Value: expr}}}
}

// Rand returns a random float between 0 and 1 ($rand).
func Rand() NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$rand", Value: bson.D{}}}}
}

type rangeOptions struct {
	step any
}

func WithRangeStep[T NumberResolver](step T) Option[rangeOptions] {
	return func(o *rangeOptions) {
		o.step = step
	}
}

// Range outputs an array of integers from start (inclusive) to end (exclusive) ($range).
// Optionally provide a step to control the increment via WithRangeStep; defaults to 1.
func Range[T NumberResolver, U NumberResolver](start T, end U, opts ...Option[rangeOptions]) ArrayExpr {
	var o rangeOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.A{start, end}
	if o.step != nil {
		args = append(args, o.step)
	}
	return ArrayExpr{expr: bson.D{{Key: "$range", Value: args}}}
}

type reduceOptions struct {
	as           any
	valueAs      any
	arrayIndexAs any
}

// WithReduceAs names the variable that represents each element of the input array; defaults to "this" (MongoDB 8.3+).
func WithReduceAs(as string) Option[reduceOptions] {
	return func(o *reduceOptions) {
		o.as = as
	}
}

// WithReduceValueAs names the variable that represents the cumulative value; defaults to "value" (MongoDB 8.3+).
func WithReduceValueAs(valueAs string) Option[reduceOptions] {
	return func(o *reduceOptions) {
		o.valueAs = valueAs
	}
}

// WithReduceArrayIndexAs names the variable that represents the index of the current element (MongoDB 8.3+).
func WithReduceArrayIndexAs(arrayIndexAs string) Option[reduceOptions] {
	return func(o *reduceOptions) {
		o.arrayIndexAs = arrayIndexAs
	}
}

// Reduce applies the in subexpression to each element of the input array, accumulating a single
// result starting from initialValue ($reduce). The in expression may reference $$value (the
// accumulated value) and $$this (the current element). Rename those variables and expose the
// element index via WithReduceValueAs, WithReduceAs, and WithReduceArrayIndexAs.
func Reduce[T ArrayResolver](input T, initialValue Expr, in Expr, opts ...Option[reduceOptions]) AnyExpr {
	var o reduceOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{
		{Key: "input", Value: input},
		{Key: "initialValue", Value: initialValue},
		{Key: "in", Value: in},
	}
	if o.as != nil {
		doc = append(doc, bson.E{Key: "as", Value: o.as})
	}
	if o.valueAs != nil {
		doc = append(doc, bson.E{Key: "valueAs", Value: o.valueAs})
	}
	if o.arrayIndexAs != nil {
		doc = append(doc, bson.E{Key: "arrayIndexAs", Value: o.arrayIndexAs})
	}
	return AnyExpr{expr: bson.D{{Key: "$reduce", Value: doc}}}
}

type regexOptions struct {
	options any
}

func WithRegexOptions(options string) Option[regexOptions] {
	return func(o *regexOptions) {
		o.options = options
	}
}

// RegexFind applies a regular expression to a string and returns information on the first matched substring ($regexFind).
// Optionally set the regex option flags via WithRegexOptions.
func RegexFind[T StringResolver](input T, regex Expr, opts ...Option[regexOptions]) ObjectExpr {
	var o regexOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.D{{Key: "input", Value: input}, {Key: "regex", Value: regex}}
	if o.options != nil {
		args = append(args, bson.E{Key: "options", Value: o.options})
	}
	return ObjectExpr{expr: bson.D{{Key: "$regexFind", Value: args}}}
}

// RegexFindAll applies a regular expression to a string and returns information on all matched substrings ($regexFindAll).
// Optionally set the regex option flags via WithRegexOptions.
func RegexFindAll[T StringResolver](input T, regex Expr, opts ...Option[regexOptions]) ArrayExpr {
	var o regexOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.D{{Key: "input", Value: input}, {Key: "regex", Value: regex}}
	if o.options != nil {
		args = append(args, bson.E{Key: "options", Value: o.options})
	}
	return ArrayExpr{expr: bson.D{{Key: "$regexFindAll", Value: args}}}
}

// RegexMatch applies a regular expression to a string and returns true if a match is found ($regexMatch).
// Optionally set the regex option flags via WithRegexOptions.
func RegexMatch[T StringResolver](input T, regex Expr, opts ...Option[regexOptions]) BoolExpr {
	var o regexOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.D{{Key: "input", Value: input}, {Key: "regex", Value: regex}}
	if o.options != nil {
		args = append(args, bson.E{Key: "options", Value: o.options})
	}
	return BoolExpr{expr: bson.D{{Key: "$regexMatch", Value: args}}}
}

// ReplaceAll replaces all instances of a search string in an input string with a replacement string ($replaceAll).
func ReplaceAll[T StringResolver, R StringResolver](input T, find Expr, replacement R) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$replaceAll", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "find", Value: find},
		{Key: "replacement", Value: replacement},
	}}}}
}

// ReplaceOne replaces the first instance of a matched string in a given input ($replaceOne).
func ReplaceOne[T StringResolver, R StringResolver](input T, find Expr, replacement R) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$replaceOne", Value: bson.D{
		{Key: "input", Value: input},
		{Key: "find", Value: find},
		{Key: "replacement", Value: replacement},
	}}}}
}

// ReverseArray returns an array with the elements in reverse order ($reverseArray).
func ReverseArray[T ArrayResolver](expr T) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$reverseArray", Value: expr}}}
}

type roundOptions struct {
	place any
}

func WithRoundPlace[T NumberResolver](place T) Option[roundOptions] {
	return func(o *roundOptions) {
		o.place = place
	}
}

// Round rounds a number to a whole integer or to a specified decimal place ($round).
// Optionally provide a decimal place via WithRoundPlace; when omitted the array form
// is still used: [$number] (equivalent to place 0).
func Round[T NumberResolver](number T, opts ...Option[roundOptions]) NumberExpr {
	var o roundOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.A{number}
	if o.place != nil {
		args = append(args, o.place)
	}
	return NumberExpr{expr: bson.D{{Key: "$round", Value: args}}}
}

// Rtrim removes whitespace characters or the specified characters from the end of a string ($rtrim).
// Optionally specify the characters to remove via WithTrimChars; when omitted, whitespace is removed.
func Rtrim[T StringResolver](input T, opts ...Option[trimOptions]) StringExpr {
	var o trimOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.D{{Key: "input", Value: input}}
	if o.chars != nil {
		args = append(args, bson.E{Key: "chars", Value: o.chars})
	}
	return StringExpr{expr: bson.D{{Key: "$rtrim", Value: args}}}
}

// Second returns the seconds for a date, from 0 to 60 (leap seconds) ($second).
func Second[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$second", date, opts...)
}

type serializeEJSONOptions struct {
	relaxed    any
	onError    any
	hasOnError bool
}

// WithSerializeEJSONRelaxed selects Relaxed Extended JSON output when true; defaults to Canonical.
func WithSerializeEJSONRelaxed[T BoolResolver](relaxed T) Option[serializeEJSONOptions] {
	return func(o *serializeEJSONOptions) {
		o.relaxed = relaxed
	}
}

// WithSerializeEJSONOnError sets the value returned if $serializeEJSON encounters a conversion error.
func WithSerializeEJSONOnError(onError Expr) Option[serializeEJSONOptions] {
	return func(o *serializeEJSONOptions) {
		o.onError = onError
		o.hasOnError = true
	}
}

// SerializeEJSON converts a BSON value into an Extended JSON document ($serializeEJSON).
// Optionally select Relaxed output via WithSerializeEJSONRelaxed and a fallback via WithSerializeEJSONOnError.
func SerializeEJSON(input Expr, opts ...Option[serializeEJSONOptions]) ObjectExpr {
	var o serializeEJSONOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "input", Value: input}}
	if o.relaxed != nil {
		doc = append(doc, bson.E{Key: "relaxed", Value: o.relaxed})
	}
	if o.hasOnError {
		doc = append(doc, bson.E{Key: "onError", Value: o.onError})
	}
	return ObjectExpr{expr: bson.D{{Key: "$serializeEJSON", Value: doc}}}
}

// SetDifference returns elements in the first set but not the second ($setDifference).
func SetDifference[T ArrayResolver, U ArrayResolver](expr1 T, expr2 U) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$setDifference", Value: bson.A{expr1, expr2}}}}
}

// SetDocField adds, updates, or removes the specified field in a document ($setField).
// field must resolve to a string constant and input must resolve to a document.
// Set value to the "$$REMOVE" system variable to remove field from the document.
func SetDocField[T StringResolver, U ObjectResolver](field T, input U, value Expr) ObjectExpr {
	return ObjectExpr{expr: bson.D{{Key: "$setField", Value: bson.D{
		{Key: "field", Value: field},
		{Key: "input", Value: input},
		{Key: "value", Value: value},
	}}}}
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

type similarityOptions struct {
	score any
}

// WithSimilarityScore normalizes the similarity result to a value between 0 and 1 when true.
func WithSimilarityScore(score bool) Option[similarityOptions] {
	return func(o *similarityOptions) {
		o.score = score
	}
}

// SimilarityCosine returns the cosine similarity between two equal-length vectors ($similarityCosine).
// Optionally normalize the result for use as a vector search score via WithSimilarityScore.
func SimilarityCosine[T ArrayResolver, U ArrayResolver](a T, b U, opts ...Option[similarityOptions]) NumberExpr {
	return similarity("$similarityCosine", a, b, opts...)
}

// SimilarityDotProduct returns the dot product similarity between two equal-length vectors ($similarityDotProduct).
// Optionally normalize the result for use as a vector search score via WithSimilarityScore.
func SimilarityDotProduct[T ArrayResolver, U ArrayResolver](a T, b U, opts ...Option[similarityOptions]) NumberExpr {
	return similarity("$similarityDotProduct", a, b, opts...)
}

// SimilarityEuclidean returns the Euclidean similarity between two equal-length vectors ($similarityEuclidean).
// Optionally normalize the result for use as a vector search score via WithSimilarityScore.
func SimilarityEuclidean[T ArrayResolver, U ArrayResolver](a T, b U, opts ...Option[similarityOptions]) NumberExpr {
	return similarity("$similarityEuclidean", a, b, opts...)
}

func similarity[T ArrayResolver, U ArrayResolver](op string, a T, b U, opts ...Option[similarityOptions]) NumberExpr {
	var o similarityOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "vectors", Value: bson.A{a, b}}}
	if o.score != nil {
		doc = append(doc, bson.E{Key: "score", Value: o.score})
	}
	return NumberExpr{expr: bson.D{{Key: op, Value: doc}}}
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

type sliceOptions struct {
	position any
}

func WithSlicePosition[T NumberResolver](position T) Option[sliceOptions] {
	return func(o *sliceOptions) {
		o.position = position
	}
}

// Slice returns n elements of an array ($slice).
// Optionally provide a starting index via WithSlicePosition; otherwise elements are
// taken from the beginning.
func Slice[T ArrayResolver](expression T, n Expr, opts ...Option[sliceOptions]) ArrayExpr {
	var o sliceOptions
	for _, opt := range opts {
		opt(&o)
	}
	if o.position != nil {
		return ArrayExpr{expr: bson.D{{Key: "$slice", Value: bson.A{expression, o.position, n}}}}
	}
	return ArrayExpr{expr: bson.D{{Key: "$slice", Value: bson.A{expression, n}}}}
}

// Split splits a string into substrings based on a delimiter and returns an array of substrings ($split).
func Split[T StringResolver](str T, delimiter Expr) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$split", Value: bson.A{str, delimiter}}}}
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

// Strcasecmp performs case-insensitive string comparison and returns 0 if equivalent, 1 if the first is greater, and -1 if less ($strcasecmp).
func Strcasecmp[T StringResolver, U StringResolver](expr1 T, expr2 U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$strcasecmp", Value: bson.A{expr1, expr2}}}}
}

// StrLenBytes returns the number of UTF-8 encoded bytes in a string ($strLenBytes).
func StrLenBytes[T StringResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$strLenBytes", Value: expr}}}
}

// StrLenCP returns the number of UTF-8 code points in a string ($strLenCP).
func StrLenCP[T StringResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$strLenCP", Value: expr}}}
}

// Substr returns a substring of a string. Deprecated; use SubstrBytes or SubstrCP ($substr).
func Substr[T StringResolver, U NumberResolver, V NumberResolver](str T, start U, length V) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$substr", Value: bson.A{str, start, length}}}}
}

// SubstrBytes returns the substring of a string starting at the specified UTF-8 byte index for the specified number of bytes ($substrBytes).
func SubstrBytes[T StringResolver, U NumberResolver, V NumberResolver](str T, start U, length V) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$substrBytes", Value: bson.A{str, start, length}}}}
}

// SubstrCP returns the substring of a string starting at the specified UTF-8 code point index for the specified number of code points ($substrCP).
func SubstrCP[T StringResolver, U NumberResolver, V NumberResolver](str T, start U, length V) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$substrCP", Value: bson.A{str, start, length}}}}
}

// Subtract returns a minus b ($subtract).
// TODO: $subtract also supports date-date → millis and date-millis → date;
// those variants are not yet modeled here.
func Subtract[T NumberResolver, U NumberResolver](a T, b U) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$subtract", Value: bson.A{a, b}}}}
}

// Subtype returns the subtype of a value as an integer ($subtype).
// The argument must resolve to a BinData value.
func Subtype(expr Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$subtype", Value: expr}}}
}

// Sum returns the sum of numeric expressions ($sum).
func Sum(exprs ...Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$sum", Value: exprs}}}
}

type switchOptions struct {
	defaultVal any
	hasDefault bool
}

// WithSwitchDefault sets the value returned when no branch case evaluates to true.
func WithSwitchDefault(defaultVal Expr) Option[switchOptions] {
	return func(o *switchOptions) {
		o.defaultVal = defaultVal
		o.hasDefault = true
	}
}

// Switch evaluates each branch case in order and returns the then value of the first case that is
// true ($switch). Build branches with Case. Optionally provide a fallback via WithSwitchDefault;
// without one, an unmatched input produces a runtime error.
func Switch(branches []SwitchCase, opts ...Option[switchOptions]) AnyExpr {
	var o switchOptions
	for _, opt := range opts {
		opt(&o)
	}
	b := make(bson.A, len(branches))
	for i, br := range branches {
		c := br.switchCase()
		b[i] = bson.D{{Key: "case", Value: c.cond}, {Key: "then", Value: c.then}}
	}
	doc := bson.D{{Key: "branches", Value: b}}
	if o.hasDefault {
		doc = append(doc, bson.E{Key: "default", Value: o.defaultVal})
	}
	return AnyExpr{expr: bson.D{{Key: "$switch", Value: doc}}}
}

// Tan returns the tangent of a value in radians ($tan).
func Tan[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$tan", Value: expr}}}
}

// Tanh returns the hyperbolic tangent of a value in radians ($tanh).
func Tanh[T NumberResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$tanh", Value: expr}}}
}

// ToArray converts a value to an array ($toArray).
// Returns null if the value is null or missing; errors if it cannot be converted.
func ToArray(expr Expr) ArrayExpr {
	return ArrayExpr{expr: bson.D{{Key: "$toArray", Value: expr}}}
}

// ToBool converts a value to a boolean ($toBool).
func ToBool(expr Expr) BoolExpr {
	return BoolExpr{expr: bson.D{{Key: "$toBool", Value: expr}}}
}

// ToDate converts a value to a Date ($toDate).
func ToDate(expr Expr) DateExpr {
	return DateExpr{expr: bson.D{{Key: "$toDate", Value: expr}}}
}

// ToDecimal converts a value to a Decimal128 ($toDecimal).
func ToDecimal(expr Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$toDecimal", Value: expr}}}
}

// ToDouble converts a value to a double ($toDouble).
func ToDouble(expr Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$toDouble", Value: expr}}}
}

// ToHashedIndexKey computes the hash of the input using MongoDB's hashed-index hash function ($toHashedIndexKey).
func ToHashedIndexKey(value Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$toHashedIndexKey", Value: value}}}
}

// ToInt converts a value to an integer ($toInt).
func ToInt(expr Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$toInt", Value: expr}}}
}

// ToLong converts a value to a long ($toLong).
func ToLong(expr Expr) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$toLong", Value: expr}}}
}

// ToLower converts a string to lowercase ($toLower).
func ToLower[T StringResolver](expr T) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$toLower", Value: expr}}}
}

// ToObject converts a string to an object ($toObject).
// Returns null if the value is null or missing; errors if it cannot be converted.
func ToObject(expr Expr) ObjectExpr {
	return ObjectExpr{expr: bson.D{{Key: "$toObject", Value: expr}}}
}

// ToObjectId converts a value to an ObjectId ($toObjectId).
func ToObjectId(expr Expr) AnyExpr {
	return AnyExpr{expr: bson.D{{Key: "$toObjectId", Value: expr}}}
}

// ToString converts a value to a string ($toString).
func ToString(expr Expr) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$toString", Value: expr}}}
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

// ToUpper converts a string to uppercase ($toUpper).
func ToUpper[T StringResolver](expr T) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$toUpper", Value: expr}}}
}

type trimOptions struct {
	chars any
}

func WithTrimChars[T StringResolver](chars T) Option[trimOptions] {
	return func(o *trimOptions) {
		o.chars = chars
	}
}

// Trim removes whitespace or the specified characters from the beginning and end of a string ($trim).
// Optionally specify the characters to remove via WithTrimChars; when omitted, whitespace is removed.
func Trim[T StringResolver](input T, opts ...Option[trimOptions]) StringExpr {
	var o trimOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.D{{Key: "input", Value: input}}
	if o.chars != nil {
		args = append(args, bson.E{Key: "chars", Value: o.chars})
	}
	return StringExpr{expr: bson.D{{Key: "$trim", Value: args}}}
}

type truncOptions struct {
	place any
}

func WithTruncPlace[T NumberResolver](place T) Option[truncOptions] {
	return func(o *truncOptions) {
		o.place = place
	}
}

// Trunc truncates a number to a whole integer or to a specified decimal place ($trunc).
// Optionally provide a decimal place via WithTruncPlace; when omitted the array form
// is still used: [$number] (equivalent to place 0).
func Trunc[T NumberResolver](number T, opts ...Option[truncOptions]) NumberExpr {
	var o truncOptions
	for _, opt := range opts {
		opt(&o)
	}
	args := bson.A{number}
	if o.place != nil {
		args = append(args, o.place)
	}
	return NumberExpr{expr: bson.D{{Key: "$trunc", Value: args}}}
}

// Type returns the BSON data type of the expression as a string ($type).
func Type(expr Expr) StringExpr {
	return StringExpr{expr: bson.D{{Key: "$type", Value: expr}}}
}

// UnsetDocField removes the specified field from a document ($unsetField).
// field must resolve to a string constant and input must resolve to a document.
func UnsetDocField[T StringResolver, U ObjectResolver](field T, input U) ObjectExpr {
	return ObjectExpr{expr: bson.D{{Key: "$unsetField", Value: bson.D{
		{Key: "field", Value: field},
		{Key: "input", Value: input},
	}}}}
}

type zipOptions struct {
	useLongestLength any
	defaults         any
}

func WithZipUseLongestLength(useLongestLength bool) Option[zipOptions] {
	return func(o *zipOptions) {
		o.useLongestLength = useLongestLength
	}
}

func WithZipDefaults(defaults ...Expr) Option[zipOptions] {
	return func(o *zipOptions) {
		o.defaults = defaults
	}
}

// TsIncrement returns the incrementing ordinal from a timestamp as a long ($tsIncrement).
func TsIncrement[T TimestampResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$tsIncrement", Value: expr}}}
}

// TsSecond returns the seconds from a timestamp as a long ($tsSecond).
func TsSecond[T TimestampResolver](expr T) NumberExpr {
	return NumberExpr{expr: bson.D{{Key: "$tsSecond", Value: expr}}}
}

// Week returns the week number for a date, from 0 to 53 ($week).
func Week[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$week", date, opts...)
}

// Year returns the year for a date ($year).
func Year[T DateResolver | TimestampResolver | ObjectIDResolver](date T, opts ...Option[datePartOptions]) NumberExpr {
	return datePart("$year", date, opts...)
}

// Zip merges arrays together into an array of arrays ($zip).
// When WithZipUseLongestLength(true) is provided, the output length is determined by the
// longest input array; pass default values for shorter arrays via WithZipDefaults. Otherwise
// the output length is the shortest input array and defaults must be empty.
func Zip(inputs Expr, opts ...Option[zipOptions]) ArrayExpr {
	var o zipOptions
	for _, opt := range opts {
		opt(&o)
	}
	doc := bson.D{{Key: "inputs", Value: inputs}}
	if o.useLongestLength != nil {
		doc = append(doc, bson.E{Key: "useLongestLength", Value: o.useLongestLength})
	}
	if o.defaults != nil {
		doc = append(doc, bson.E{Key: "defaults", Value: o.defaults})
	}
	return ArrayExpr{expr: bson.D{{Key: "$zip", Value: doc}}}
}
