package golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	
	script := "INSERT INTO customer(id,name) VALUES('2','budi'),('3','boy')"
	_, err := db.ExecContext(ctx,script) // ExecContext disarankan hanya untuk perintah SQL yang tidak membutuhkan hasil/result seperti INSERT,UPDATE,DELETE

	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")

}

func TestQuerySql(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	
	script := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx,script) // QueryContext disarankan hanya untuk perintah SQL yang membutuhkan hasil/result seperti SELECT

	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}

		fmt.Println("Id: ",id)
		fmt.Println("Name: ",name)
	}
	defer rows.Close()

}

func TestQuerySqlComplex(t *testing.T)  {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	
	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx,script) // QueryContext disarankan hanya untuk perintah SQL yang membutuhkan hasil/result seperti SELECT

	if err != nil {
		panic(err)
	}

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
	defer rows.Close()
}