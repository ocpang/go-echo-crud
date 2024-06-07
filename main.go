package main

import (
	"fmt"
	_ "log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	sqlx "github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

type Employee struct {
	Id      int    `json:"id" db:"id"`
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Address string `json:"address" db:"address"`
}

type Response struct {
	Message string
	Status  bool
}

func main() {
	db, err := sqlx.Connect("mysql", "root@tcp(127.0.0.1:3306)/go-book")

	// if err = db.Ping(); err != nil {
	// 	fmt.Println(err)
	// }

	if err != nil {
		panic(err)
	}

	respon := Response{
		Message: "Success executing the query",
		Status:  true,
	}

	responError := Response{
		Message: "Error executing the query",
		Status:  false,
	}

	e := echo.New()

	e.Use(middleware.CORS())

	//get all users
	e.GET("/users", func(c echo.Context) error {
		rows, _ := db.Queryx("select * from users")

		var users []Employee

		for rows.Next() {
			place := Employee{}
			rows.StructScan(&place)
			users = append(users, place)
		}

		return c.JSON(http.StatusOK, users)
	})

	//get spesific user
	e.GET("/users/:id", func(c echo.Context) error {

		id, _ := strconv.Atoi(c.Param("id"))

		user := Employee{}
		sql := fmt.Sprintf(`SELECT * FROM users WHERE id = %d`, id)
		err = db.Get(&user, sql)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, user)

	})

	//add user
	e.POST("/users", func(c echo.Context) error {
		reqBody := Employee{}
		c.Bind(&reqBody)

		_, err = db.NamedExec("insert into users(name, phone, address) values (:name, :phone, :address)", reqBody)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(http.StatusOK, respon)
	})

	//update user
	e.PUT("/users/update/:id", func(c echo.Context) error {
		reqBody := Employee{}

		// c.Bind(&reqBody)
		if err := c.Bind(&reqBody); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		id, _ := strconv.Atoi(c.Param("id")) // Konversi id dari string ke int
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
		}
		reqBody.Id = id
		_, errQuery := db.NamedExec("update users SET name= :name, phone= :phone, address= :address WHERE id= :id", reqBody)
		if errQuery != nil {
			return c.JSON(http.StatusInternalServerError, responError)
		}

		return c.JSON(http.StatusOK, respon)
	})

	//delete user
	e.DELETE("/users/delete/:id", func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))

		sql := fmt.Sprintf(`DELETE FROM users WHERE id = %d`, id)
		_, err = db.Exec(sql)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, respon)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
