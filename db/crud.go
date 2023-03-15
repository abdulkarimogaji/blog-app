package db

type DB_GetAll_Response struct {
	Data []any
}

func (db *DBService) GetAll() (DB_GetAll_Response, error) {
	return DB_GetAll_Response{Data: []any{}}, nil
}
