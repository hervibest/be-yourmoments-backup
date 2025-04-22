package nullable

import (
	"database/sql"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func GRPCtoSQLDouble(v *wrapperspb.DoubleValue) sql.NullFloat64 {
	if v == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: v.Value, Valid: true}
}

func GRPCtoSQLString(v *wrapperspb.StringValue) sql.NullString {
	if v == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: v.Value, Valid: true}
}
