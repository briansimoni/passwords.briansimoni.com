package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
)


type User struct {
	Id string
	Email string
	Hash string
	Secrets map[string]string
	PageData map[string]string
}

func insertUser(email, hash string) error {

	db, err := sql.Open("mysql", "root:root@/SimoniPassword")
	if err != nil {
		return err
	}

	// insert
	stmt, err := db.Prepare("INSERT users SET email=?,passwordHash=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(email, hash)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(id)

	db.Close()

	return nil
}

func getUser(queryEmail string) User {

	fmt.Println("trying to get the user")
	db, err := sql.Open("mysql", "root:root@/SimoniPassword")
	checkErr(err)

	rows, err := db.Query("SELECT * FROM users WHERE email='" + queryEmail + "'")
	checkErr(err)

	var user User
	for rows.Next() {
		var id string
		var email string
		var passwordHash string
		err = rows.Scan(&id, &email, &passwordHash)
		checkErr(err)
		user = User{id, email, passwordHash, make(map[string]string, 0), make(map[string]string, 0)}

	}

	db.Close()

	user = getSecrets(user)
	return user
}

func insertSecret (application, password, userId string) error {
	db, err := sql.Open("mysql", "root:root@/SimoniPassword")
	if err != nil {
		return err
	}

	encryptedApplication := encrypt(application)
	encryptedPassword := encrypt(password)

	stmt, err := db.Prepare("INSERT secrets SET encrypted_application=?,encrypted_password=?,user_id=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(encryptedApplication, encryptedPassword, userId)
	if err != nil {
		return err
	}

	db.Close()
	return nil
}

func getSecrets(user User) User {
	db, err := sql.Open("mysql", "root:root@/SimoniPassword")
	checkErr(err)

	rows, err := db.Query("SELECT encrypted_application, encrypted_password FROM secrets WHERE user_id = '" + user.Id + "'")
	checkErr(err)

	for rows.Next() {
		var encrypted_application string
		var encrypted_password string
		err = rows.Scan(&encrypted_application, &encrypted_password)
		checkErr(err)
		user.Secrets[encrypted_application] = encrypted_password

	}

	db.Close()
	return user
}

func deleteSecretFromDB(app string, user User) error {
	db, err := sql.Open("mysql", "root:root@/SimoniPassword")
	if err != nil {
		return err
	}

	stmt, err := db.Prepare("DELETE FROM secrets WHERE encrypted_application=? AND user_id=?;")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(app, user.Id)
	if err != nil {
		return err
	}

	db.Close()
	return nil

}

func updateApplicationPassword(password string, app string, user User) error {
	db, err := sql.Open("mysql", "root:root@/SimoniPassword")
	if err != nil {
		return err
	}


	stmt, err := db.Prepare("UPDATE secrets SET encrypted_password=? WHERE encrypted_application=? AND user_id=?;")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(password, app, user.Id)
	if err != nil {
		return err
	}

	db.Close()
	return nil
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}