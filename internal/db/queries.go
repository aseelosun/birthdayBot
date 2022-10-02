package db

import (
	entity "birthdayBot/internal/entities"
	"database/sql"
)

func InsertIntoBirthdayList(db *sql.DB, name string, birthdate string, status int) (int, error) {
	var id int
	sqlStmt := `INSERT INTO birthdayList (bName, birthDate,bStatus) VALUES ($1, $2, $3) returning b_id;`
	err := db.QueryRow(sqlStmt, name, birthdate, status).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil
}

func GetUserStatus(db *sql.DB, id int) (int, error) {
	var (
		status int
	)
	sqlStatement := "SELECT bStatus FROM birthdayList WHERE b_id = $1;"
	err := db.QueryRow(sqlStatement, id).Scan(&status)
	if err != nil {
		return status, err
	}
	return status, nil
}

func UpdateUserStatus(db *sql.DB, id int, status int) error {
	sqlStatement := `UPDATE birthdayList SET bstatus = $2 WHERE b_id = $1;`
	_, err := db.Exec(sqlStatement, id, status)
	if err != nil {
		return err
	}
	return nil
}

func UpdateName(db *sql.DB, id int, name string) error {
	sqlStatement := `UPDATE birthdayList SET bName= $2 WHERE b_id = $1;`
	_, err := db.Exec(sqlStatement, id, name)
	if err != nil {
		return err
	}
	return nil
}
func UpdateBirthDate(db *sql.DB, id int, bithDate string) error {
	sqlStatement := `UPDATE birthdayList SET birthDate= $2 WHERE b_id = $1;`
	_, err := db.Exec(sqlStatement, id, bithDate)
	if err != nil {
		return err
	}
	return nil
}

func GetBirthdayList(db *sql.DB) ([]entity.MyBirthdayList, error) {
	var (
		obj          entity.MyBirthdayList
		bName        string
		birthdayDate string
		objArray     []entity.MyBirthdayList
	)
	rows, err := db.Query("select bName, birthDate from birthdayList where bStatus = 2 and isDeleted = false;")

	if err != nil {
		if err != sql.ErrNoRows {
			return objArray, nil
		}
		return objArray, err
	}

	for rows.Next() {
		err := rows.Scan(&bName, &birthdayDate)
		if err != nil {
			if err != sql.ErrNoRows {
				return objArray, nil
			}
			return objArray, err
		}
		obj.BirthDate = birthdayDate[:len(birthdayDate)-10]
		obj.Name = bName

		objArray = append(objArray, obj)

	}
	return objArray, nil

}
func DeleteFromBirthdayList(db *sql.DB, name string) (bool, error) {
	sqlStatement := `UPDATE birthdayList SET isDeleted = true WHERE bName = $1;`
	_, err := db.Exec(sqlStatement, name)
	if err != nil {
		return false, err

	}
	return true, nil
}
