package helpers

import "github.com/jackc/pgtype"

func ConvertInt4ArrayToInt32(elems pgtype.Int4Array) []int32 {
	if len(elems.Elements) == 0 || elems.Elements[0].Status == pgtype.Null {
		return []int32{}
	}
	res := make([]int32, 0, len(elems.Elements))
	for i := range elems.Elements {
		res = append(res, elems.Elements[i].Int)
	}
	return res
}
