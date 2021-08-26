package db

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	sq "github.com/Masterminds/squirrel"
)

func Remove(c *gin.Context) error {
	db, err := sql.Open("sqlite3", "../proxy/ent2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	selectBuilder := sq.Delete(c.Param("typeName")).Where(sq.Eq{"id": c.Param("id")})
	_, err = selectBuilder.RunWith(db).Exec()
	if err != nil {
		return err
	}

	c.JSON(200, gin.H{
		"message": "The item successfully deleted!",
	})

	return nil
}

func doFetch(colLen int, rows *sql.Rows) string {
	result := ""
	vals := make([]interface{}, colLen)
	for rows.Next() {
		for i := 0; i < colLen; i++ {
			vals[i] = new(string)
		}
		err := rows.Scan(vals...)
		if err != nil {
			log.Fatal(err.Error()) // if wrong type
		}
		for _, value := range vals {
			result += *(value.(*string)) + "|"
		}
		result += "\n"
	}

	return result
}

func FetchAll(c *gin.Context) error {
	db, err := sql.Open("sqlite3", "../proxy/ent2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	selectBuilder := sq.Select("*").From(c.Param("typeName"))
	rows, err := selectBuilder.RunWith(db).Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	result := doFetch(len(columns), rows)

	c.JSON(200, gin.H{
		"data": result,
	})

	return nil
}

func Fetch(c *gin.Context) error {
	db, err := sql.Open("sqlite3", "../proxy/ent2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	selectBuilder := sq.Select("*").From(c.Param("typeName")).Where(sq.Eq{"id": c.Param("id")})
	rows, err := selectBuilder.RunWith(db).Query()
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	colLen := len(columns)

	result := doFetch(colLen, rows)

	c.JSON(200, gin.H{
		"data": strings.TrimRight(result, "|"),
	})

	return nil
}

func Create(c *gin.Context) error {
	db, err := sql.Open("sqlite3", "../proxy/ent2.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bodyData, _ := ioutil.ReadAll(c.Request.Body)
	var data map[string]string

	if err := json.Unmarshal(bodyData, &data); err != nil {
		panic(err)
	}

	log.Println("data is ", data)
	keys := make([]string, 0, len(data))
	values := make([]string, 0, len(data))
	for k, v := range data {
		keys = append(keys, k)
		values = append(values, v)
	}
	log.Println("keys are ", strings.Join(keys, ","))
	log.Println("values are ", strings.Join(values, ","))

	for i, s := range values {
		values[i] = "'" + s + "'"
	}
	sqlStr := "INSERT INTO " + c.Param("typeName") + "(" + strings.Join(keys, ",") + ") VALUES " + "(" + strings.Join(values, ",") + ")"
	log.Println("sqlStr is ", sqlStr)
	_, err = db.Exec(sqlStr)
	if err != nil {
		return err
	}

	c.JSON(200, gin.H{
		"message": "New reading item successfully created!",
		"data":    sqlStr,
	})

	return nil
}
