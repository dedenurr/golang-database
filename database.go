package golang_database

import (
	"database/sql"
	"time"
)

func GetConnection() *sql.DB {
	db, err := sql.Open("mysql","root:@tcp(localhost:3306)/belajar_golang_database?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)// jumlah maksimal koneksi yang dibuat 10
	db.SetMaxOpenConns(100)// maksimal open connetion 100
	db.SetConnMaxIdleTime(5 * time.Minute)// jika dalam 5 menit tidak ada aktifitas, koneksi akan di close dalam 5 menit
	db.SetConnMaxLifetime(60 * time.Minute)// koneksi baru akan dibuat setelah 60 menit

	return db

}