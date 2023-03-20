package db

import "fmt"

func getIntClause(fieldName string, val *int) string {
	if val == nil {
		return "1"
	}
	return fmt.Sprintf("%s = %v", fieldName, *val)
}

func getLikeClause(fieldName string, val *string) string {
	if val == nil {
		return "1"
	}
	return fmt.Sprintf("%s LIKE '%%%s%%'", fieldName, *val)
}

func getTimeClause(fieldName, sign string, val *string) string {
	if val == nil {
		return "1"
	}
	return fmt.Sprintf("%s %s '%v'", fieldName, sign, *val)
}

func getPaginationStr(params PaginationParams) string {
	offset := params.Limit * params.Page
	if params.Limit == 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %v OFFSET %v", params.Limit, offset)
}
