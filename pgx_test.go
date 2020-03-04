package datatypes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
)

const connStr = "postgres://testuser:testuser@192.168.1.20:5432/testdb"

var pgxconn *pgx.Conn

type DataTest struct {
	ID       int
	Name     NullString
	Tstamp   OnlyTime
	IP       IPAddr
	Dob      UnixTime
	I2       NullInt16
	I4       NullInt32
	I8       NullInt64
	F4       NullFloat32
	F8       NullFloat64
	Verified NullBool
}

func TestInitPgx(t *testing.T) {
	fmt.Println("\n\nTestInitPgx ***")
	var err error
	pgxconn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		t.Errorf("Connection failed. %v", err)
		t.FailNow()
	}
	err = pgxconn.Ping(context.Background())
	if err != nil {
		t.Errorf("Ping failed. %v", err)
		t.FailNow()
	}

	//Drop previous table, error expected, ignore that
	sql := "drop table datatest;"
	pgxconn.Exec(context.Background(), sql)

	//Create table daltest
	sql = "Create table datatest( " +
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
	_, err = pgxconn.Exec(context.Background(), sql)
	if err != nil {
		t.Errorf("NonQuery 'Create table' failed. %v", err)
		t.FailNow()
	}

	// insert
	sql = "insert into datatest(id, name, tstamp, Ip, dob, cint2, cint4, cint8, cfloat4, cfloat8, verified) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);"
	_, err = pgxconn.Exec(context.Background(), sql, 1, "row-1", time.Now(), "192.168.1.20", time.Now().AddDate(-18, -2, -11), 10, 45856, 9415423285, 21.6745, 94154232.850091, true)
	if err != nil {
		t.Errorf("Insert 'row-1' failed. %v", err)
		t.FailNow()
	}

	sql = "insert into datatest(id, name, tstamp, ip) values($1, $2, $3, $4);"
	_, err = pgxconn.Exec(context.Background(), sql, 2, "row-2", time.Now(), "0.0.0.0")
	if err != nil {
		t.Errorf("Insert 'row-2' failed. %v", err)
		t.FailNow()
	}
}

func TestWithValuesPgx(t *testing.T) {
	fmt.Println("\n\nTestWithValuesPgx ***")

	// fetch from db
	sql := "select id, name, tstamp, Ip, dob, cint2, cint4, cint8, cfloat4, cfloat8, verified from datatest where ID=$1"
	rows, err := pgxconn.Query(context.Background(), sql, 1)
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

func TestWithNullValuesPgx(t *testing.T) {
	fmt.Println("\n\nTestWithNullValuesPgx ***")

	// fetch from db
	sql := "select id, name, tstamp, Ip, dob, cint2, cint4, cint8, cfloat4, cfloat8, verified from datatest where ID=$1"
	rows, err := pgxconn.Query(context.Background(), sql, 2)
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

func TestJsonBind(t *testing.T) {
	fmt.Println("\n\nTestJsonBind ***")

	jsonstr := `{"ID":1,"Name":"sam","Dob":null}`
	var dtest DataTest
	if err := json.Unmarshal([]byte(jsonstr), &dtest); err != nil {
		t.Error(err)
	}
	fmt.Printf("struct: %#v\n", dtest)
}

func TestClosePgx(t *testing.T) {
	pgxconn.Close(context.Background())
}

func toJSON(d interface{}) *bytes.Buffer {
	data, err := json.Marshal(d)
	if err != nil {
		fmt.Println("marshal error: ", err)
		return nil
	}
	return bytes.NewBuffer(data)
}
