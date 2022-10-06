package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type City struct {
	ID          int     `json:"id,omitempty"  db:"ID"`
	Name        string  `json:"name,omitempty"  db:"Name"`
	CountryCode string  `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string  `json:"district,omitempty"  db:"District"`
	Population  float32 `json:"population,omitempty"  db:"Population"`
}

var slice []int

// type Country struct {
// 	Code           string  `json:"code,omitempty" db:"Code"`
// 	Name           string  `json:"name,omitempty" db:"Name"`
// 	Continent      string  `json:"continent,omitempty" db:"Continent"`
// 	Region         string  `json:"region,omitempty" db:"Region"`
// 	SurfaceArea    string  `json:"surfaceArea,omitempty" db:"SurfaceArea"`
// 	IndepYear      int     `json:"indepYear,omitempty" db:"IndepYear"`
// 	Population     float32 `json:"population,omitempty" db:"Population"`
// 	LifeExpectancy float32 `json:"lifeExpectancy,omitempty" db:"LifeExpectancy"`
// 	GNP            float32 `json:"GNP,omitempty" db:"GNP"`
// 	GNPOld         float32 `json:"GNPOld,omitempty" db:"GNPOld"`
// 	LocalName      string  `json:"localName,omitempty" db:"LocalName"`
// 	GovernmentForm string  `json:"governmentForm,omitempty" db:"GovernmentForm"`
// 	HeadOfState    string  `json:"headOfState,omitempty" db:"HeadOfState"`
// 	Capital        int     `json:"capital,omitempty" db:"Capital"`
// 	Code2          string  `json:"code2,omitempty"  db:"Code2"`
// }

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	fmt.Println("Connected!")

	db = _db

	e := echo.New()
	e.GET("/cities/:cityName", getCityInfoHandler)

	e.POST("/post", postHandler)
	e.GET("/delete", deleteHandler)

	e.Start(":8080")
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")
	fmt.Println(cityName)

	var city City
	if err := db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName); errors.Is(err, sql.ErrNoRows) {
		log.Printf("No Such City Name=%s", cityName)
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}
	return c.JSON(http.StatusOK, city)
}

func postHandler(c echo.Context) error {
	data := &City{}
	err := c.Bind(data)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%+v", data))
	}

	var city City
	sql_1 := fmt.Sprintf("SELECT * FROM city WHERE ID ='%s'", data.ID)
	if err := db.Get(&city, sql_1); errors.Is(err, sql.ErrNoRows) {
		add := fmt.Sprintf("INSERT INTO city (Name, CountryCode, District, Population, ID) VALUES ('%s', '%s', '%s', %v, '%v');", data.Name, data.CountryCode, data.District, data.Population, data.ID)
		if _, err := db.Exec(add); err != nil {
			log.Fatalf("Exec error: %s", err)
		}
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	slice = append(slice, data.ID)

	return c.JSON(http.StatusOK, data.Name)
}

func deleteHandler(c echo.Context) error {
	for _, cityID := range slice {
		var city City
		sql_1 := fmt.Sprintf("SELECT * FROM city WHERE ID='%d'", cityID)
		if err := db.Get(&city, sql_1); errors.Is(err, sql.ErrNoRows) {
			log.Printf("Already no such city ID = %d", cityID)
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		} else {
			del := fmt.Sprintf("DELETE FROM city WHERE ID = %d;", cityID)
			if _, err := db.Exec(del); err != nil {
				log.Fatalf("Exec error: %s", err)
			}
			log.Printf("Deleted: %d", cityID)
		}
	}
	slice = slice[0:0]
	log.Printf("All Added Data Deleted")
	return c.JSON(http.StatusOK, "Deleted Add Data")
}
