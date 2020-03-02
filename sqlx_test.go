package datatypes

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var sqlxconn *sqlx.DB

func TestInitSQLx(t *testing.T) {
	fmt.Println("\n\nTestInitSQLx ***")
	var err error
	sqlxconn, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		t.Errorf("Connection failed. %v", err)
		t.FailNow()
	}

	//Drop previous table, error expected, ignore that
	sql := "drop table datatestX;"
	sqlxconn.Exec(sql)

	//Create table daltest
	sql = "Create table datatestX( " +
		"ID int not null," +
		"name character varying(50)," +
		"tstamp time without time zone NOT NULL," +
		"Ip inet NOT NULL," +
		"dob timestamp without time zone," +
		"cint2 smallint NULL," +
		"cint4 int NULL," +
		"cint8 bigint NULL," +
		"cfloat4 real NULL," +
		"cfloat8 numeric(10, 2) NULL," +
		"verified boolean NULL" +
		");"
	_, err = sqlxconn.Exec(sql)
	if err != nil {
		t.Errorf("NonQuery 'Create table' failed. %v", err)
		t.FailNow()
	}

	// insert
	sql = "insert into datatestX(id, name, tstamp, Ip, dob, cint2, cint4, cint8, cfloat4, cfloat8, verified) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);"
	_, err = sqlxconn.Exec(sql, 1, "row-1", time.Now(), "192.168.1.20", time.Now().AddDate(-18, -2, -11), 10, 45856, 9415423285, 21.6745, 94154232.850091, true)
	if err != nil {
		t.Errorf("Insert 'row-1' failed. %v", err)
		t.FailNow()
	}

	sql = "insert into datatestX(id, name, tstamp, ip) values($1, $2, $3, $4);"
	_, err = sqlxconn.Exec(sql, 2, "row-2", time.Now(), "0.0.0.0")
	if err != nil {
		t.Errorf("Insert 'row-2' failed. %v", err)
		t.FailNow()
	}
}

func TestWithValuesSQLx(t *testing.T) {
	fmt.Println("\n\nTestWithValuesSQLx ***")

	// fetch from db
	sql := "select id, name, tstamp, Ip, dob, cint2, cint4, cint8, cfloat4, cfloat8, verified from datatestX where ID=$1"
	rows, err := sqlxconn.Query(sql, 1)
	if err != nil {
		t.Errorf("Query failed. %v", err)
		t.FailNow()
	}
	defer rows.Close()

	// scan to struct
	row := DataTest{}
	if rows.Next() {
		err = rows.Scan(&row.ID, &row.Name, &row.Tstamp, &row.IP, &row.Dob, &row.I2, &row.I4, &row.I8, &row.F4, &row.F8, &row.Verified)
		if err != nil {
			t.Errorf("Scan failed. %v", err)
			t.FailNow()
		}
	}
	if row.ID != 1 {
		t.Errorf("\n\tExpected: %d\tReceived: %d\n", 1, row.ID)
	}
	fmt.Printf("struct: %#v\n", row)
	fmt.Printf("JSON: %#v\n", toJSON(row).String())
}

func TestWithNullValuesSQLx(t *testing.T) {
	fmt.Println("\n\nTestWithNullValuesSQLx ***")

	// fetch from db
	sql := "select id, name, tstamp, Ip, dob, cint2, cint4, cint8, cfloat4, cfloat8, verified from datatestX where ID=$1"
	rows, err := sqlxconn.Query(sql, 2)
	if err != nil {
		t.Errorf("Query failed. %v", err)
		t.FailNow()
	}
	defer rows.Close()

	// scan to struct
	row := DataTest{}
	if rows.Next() {
		err = rows.Scan(&row.ID, &row.Name, &row.Tstamp, &row.IP, &row.Dob, &row.I2, &row.I4, &row.I8, &row.F4, &row.F8, &row.Verified)
		if err != nil {
			t.Errorf("Scan failed. %v", err)
			t.FailNow()
		}
	}
	if row.ID != 2 {
		t.Errorf("\n\tExpected: %d\tReceived: %d\n", 2, row.ID)
	}
	fmt.Printf("struct: %#v\n", row)
	fmt.Printf("JSON: %#v\n", toJSON(row).String())
}

func TestCloseSQLx(t *testing.T) {
	sqlxconn.Close()
}
