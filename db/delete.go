package db

import (
	"fmt"
	"log"
)

func (d *DBStruct) DeleteRow(tableName string, id int) (int, error) {
	stmt, err := d.DB.Prepare(fmt.Sprintf("DELETE FROM %s WHERE id = ?;", tableName))
	if err != nil {
		log.Println("here? preparing")
		return 0, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	log.Println("here? executing")
	if err != nil {
		return 0, err
	}

	return id, nil
}
