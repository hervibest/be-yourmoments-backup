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

func ExtractString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func ExtractINT32(ns sql.NullInt32) int32 {
	if ns.Valid {
		return ns.Int32
	}
	return 0
}

// func ExtractBool(ns sql.NullBool) bool {
// 	if ns.Valid {
// 		return ns.func Benchmark(b *testing.B) {
// 			for i := 0; i < b.N; i++ {

// 			}
// 		}
// 	}
// 	return ""
// }
