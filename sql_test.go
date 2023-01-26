package golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()// connect ke db
	defer db.Close()// close db jika sudah tidak digunakan, menggunakan defer agar kode dibawahnya tetap dijalankan terlebih dahulu

	ctx := context.Background()// proses cancellation dan passing value
	
	script := "INSERT INTO customer(id,name) VALUES('2','budi'),('3','boy')"
	_, err := db.ExecContext(ctx,script) // ExecContext disarankan hanya untuk perintah SQL yang tidak membutuhkan hasil/result seperti INSERT,UPDATE,DELETE

	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")

}

func TestQuerySql(t *testing.T)  {
	db := GetConnection()// connect ke db
	defer db.Close()// close db jika sudah tidak digunakan, menggunakan defer agar kode dibawahnya tetap dijalankan terlebih dahulu

	ctx := context.Background()// proses cancellation dan passing value
		
	script := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx,script) // QueryContext disarankan hanya untuk perintah SQL yang membutuhkan hasil/result seperti SELECT

	if err != nil {
		panic(err)
	}

	defer rows.Close()
	
	// proses menampilkan data
	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("Id: ",id)
		fmt.Println("Name: ",name)
	}
	

}

func TestQuerySqlComplex(t *testing.T)  {
	db := GetConnection()// connect ke db
	defer db.Close()// close db jika sudah tidak digunakan, menggunakan defer agar kode dibawahnya tetap dijalankan terlebih dahulu

	ctx := context.Background()// proses cancellation dan passing value

	
	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx,script) // QueryContext disarankan hanya untuk perintah SQL yang membutuhkan hasil/result seperti SELECT

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	// proses menampilkan data
	for rows.Next() {
		var id, name string
		var email sql.NullString //untuk tipe data jika isian data kosong
		var balance int32
		var rating float64
		var birthDate sql.NullTime //untuk tipe data jika isian data kosong
		var createdAt time.Time 
		var married bool
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("==================")
		fmt.Println("Id: ",id)
		fmt.Println("Name: ",name)
		if email.Valid{ // validitas jika email kosong datanya
			fmt.Println("Email: ",email.String)
		}	
		fmt.Println("Balance: ",balance)
		fmt.Println("Rating: ",rating)
		if birthDate.Valid{  // validitas jika birth date kosong datanya
			fmt.Println("Birth Date: ",birthDate)
		}
		fmt.Println("Married: ",married)
		fmt.Println("Created At: ",createdAt)
		
	}
	
}

func TestSqlInjecton(t *testing.T)  {
	db := GetConnection()// connect ke db
	defer db.Close()// close db jika sudah tidak digunakan, menggunakan defer agar kode dibawahnya tetap dijalankan terlebih dahulu

	ctx := context.Background()// proses cancellation dan passing value
	
	username := "admin'; #"
	password := "salah"

	// ini tidak aman karena bisa kena hack
	script := "SELECT username FROM user WHERE username = '"+ username +"' AND password = '"+ password +"' LIMIT 1"

	fmt.Println(script)

	rows, err := db.QueryContext(ctx,script) // QueryContext disarankan hanya untuk perintah SQL yang membutuhkan hasil/result seperti SELECT

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	// Cek apakah username ada?
	if rows.Next(){
		var username string
		
		err := rows.Scan(&username)

		if err != nil{
			panic(err)
		}

		fmt.Println("Sukses Login",username)
	}else{
		fmt.Println("Gagal Login")
	}
	
}

func TestSqlInjectonSafe(t *testing.T)  {
	db := GetConnection()// connect ke db
	defer db.Close()// close db jika sudah tidak digunakan, menggunakan defer agar kode dibawahnya tetap dijalankan terlebih dahulu

	ctx := context.Background()// proses cancellation dan passing value
	
	username := "admin'; #"
	password := "salah"

	// ini aman karena menggunakan sql paramater
	script := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"

	fmt.Println(script)

	rows, err := db.QueryContext(ctx, script, username, password) // QueryContext disarankan hanya untuk perintah SQL yang membutuhkan hasil/result seperti SELECT

	if err != nil {
		panic(err)
	}

	defer rows.Close()

	// Cek apakah username ada?
	if rows.Next(){
		var username string
		
		err := rows.Scan(&username)

		if err != nil{
			panic(err)
		}

		fmt.Println("Sukses Login",username)
	}else{
		fmt.Println("Gagal Login")
	}
	
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()// connect ke db
	defer db.Close()// close db jika sudah tidak digunakan, menggunakan defer agar kode dibawahnya tetap dijalankan terlebih dahulu

	ctx := context.Background()// proses cancellation dan passing value
	
	username := "dede'; DROP TABLE user; #"
	password := "dede"

	// aman dari sql injection karena menggunakan parameter ? ?
	script := "INSERT INTO user(username,password) VALUES(?,?)"
	_, err := db.ExecContext(ctx, script, username, password) // ExecContext disarankan hanya untuk perintah SQL yang tidak membutuhkan hasil/result seperti INSERT,UPDATE,DELETE

	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new user")

}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()// connect ke db
	defer db.Close()// close db jika sudah tidak digunakan, menggunakan defer agar kode dibawahnya tetap dijalankan terlebih dahulu

	ctx := context.Background()// proses cancellation dan passing value
	
	email := "dede@gmail.com"
	comment := "dede ganteng"

	// aman dari sql injection karena menggunakan parameter ? ?
	script := "INSERT INTO comments(email,comment) VALUES(?,?)"
	result, err := db.ExecContext(ctx, script, email, comment) // ExecContext disarankan hanya untuk perintah SQL yang tidak membutuhkan hasil/result seperti INSERT,UPDATE,DELETE

	if err != nil {
		panic(err)
	}
	// menambahkan id auto dengan LastInsertId() 
	insertId, err := result.LastInsertId()

	if err != nil {
		panic(err)
	}


	fmt.Println("Success insert new comment with ID", insertId)

}

// prepare statement secara otomatis akan mengenali koneksi database yang digunakan, Sehingga ketika kita mengeksekusi Prepare Statement berkali-kali, maka akan menggunakan koneksi yang sama dan lebih efisien karena pembuatan prepare statement nya hanya sekali diawal saja

func TestPrepareStatement(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "INSERT INTO comments(email,comment) VALUES(?,?)"
	statment, err := db.PrepareContext(ctx,script)

	if err != nil {
		panic(err)
	}

	defer statment.Close()


	
	// memasukan data berulng kedalam database dengan preparestatement
	for i := 0; i < 10; i++ {
		email := "zahra" + strconv.Itoa(i) + "@gmail.com"
		comment := "Komentar ke " + strconv.Itoa(1)

		result, err := statment.ExecContext(ctx,email,comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("Comment Id:", id )

	}
}

func TestTransaction(t *testing.T)  {
	db := GetConnection()
	defer db.Close()
	
	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	} 

	script := "INSERT INTO comments(email,comment) VALUES(?,?)"

	// do transaction
	for i := 0; i < 10; i++ {
		email := "zahra" + strconv.Itoa(i) + "@gmail.com"
		comment := "Komentar ke " + strconv.Itoa(1)

		result, err := tx.ExecContext(ctx,script,email,comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("Comment Id:", id )

	}

	err = tx.Rollback()
	
	if err != nil {
		panic(err)
	} 

}