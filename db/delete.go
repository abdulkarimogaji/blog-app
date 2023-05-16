package db

import (
	"context"
	"fmt"
	"log"
)

func (d *DBStruct) DeleteRow(ctx context.Context, tableName string, id int) (int, error) {
	stmt, err := d.DB.PrepareContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = ?;", tableName))
	if err != nil {
		log.Println("here? preparing")
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	log.Println("here? executing")
	if err != nil {
		return 0, err
	}

	return id, nil
}
