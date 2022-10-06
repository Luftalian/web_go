package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type City struct {
	ID          int     `json:"id,omitempty"  db:"ID"`
	Name        string  `json:"name,omitempty"  db:"Name"`
	CountryCode string  `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string  `json:"district,omitempty"  db:"District"`
	Population  float32 `json:"population,omitempty"  db:"Population"`
}

type Country struct {
	Code           string  `json:"code,omitempty" db:"Code"`
	Name           string  `json:"name,omitempty" db:"Name"`
	Continent      string  `json:"continent,omitempty" db:"Continent"`
	Region         string  `json:"region,omitempty" db:"Region"`
	SurfaceArea    string  `json:"surfaceArea,omitempty" db:"SurfaceArea"`
	IndepYear      int     `json:"indepYear,omitempty" db:"IndepYear"`
	Population     float32 `json:"population,omitempty" db:"Population"`
	LifeExpectancy float32 `json:"lifeExpectancy,omitempty" db:"LifeExpectancy"`
	GNP            float32 `json:"GNP,omitempty" db:"GNP"`
	GNPOld         float32 `json:"GNPOld,omitempty" db:"GNPOld"`
	LocalName      string  `json:"localName,omitempty" db:"LocalName"`
	GovernmentForm string  `json:"governmentForm,omitempty" db:"GovernmentForm"`
	HeadOfState    string  `json:"headOfState,omitempty" db:"HeadOfState"`
	Capital        int     `json:"capital,omitempty" db:"Capital"`
	Code2          string  `json:"code2,omitempty"  db:"Code2"`
}

func main() {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	fmt.Println("Connected!")

	var city City
	var country Country
	var cityName string
	if len(os.Args) == 1 {
		log.Fatalf("Please write a city name")
	}
	cityName = os.Args[1]

	if err := db.Get(&city, "SELECT * FROM city WHERE Name='"+cityName+"'"); errors.Is(err, sql.ErrNoRows) {
		log.Printf("no such city Name = %s", cityName)
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	if err := db.Get(&country, "SELECT * FROM country WHERE CODE='"+city.CountryCode+"'"); errors.Is(err, sql.ErrNoRows) {
		log.Printf("no such country Code = %s", city.CountryCode)
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	cities := []City{}
	db.Select(&cities, "SELECT * FROM city WHERE CountryCode='JPN'")
	var city2 City
	if err := db.Get(&city2, "SELECT * FROM city WHERE ID='4080'"); errors.Is(err, sql.ErrNoRows) {
		add := fmt.Sprintf("INSERT INTO city (Name, CountryCode, District, Population, ID) VALUES ('oookayama', 'JPN', 'Tokyo', 5000, '4080');")
		if _, err := db.Exec(add); err != nil {
			log.Fatalf("Exec error: %s", err)
		}
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	var populationRation float32
	populationRation = city.Population * 100 / country.Population

	fmt.Printf("%sの人口は%d人です\n%sには%sの人口のうち%v%%の人が住んでいます。\n", cityName, int(city.Population), cityName, country.Name, populationRation)

	// fmt.Println("日本の都市一覧")
	// for _, city := range cities {
	// 	fmt.Printf("都市名：%s, 人口：%d人\n", city.Name, int(city.Population))
	// }

	if len(os.Args) > 2 && os.Args[2] == "delete" {
		if err := db.Get(&city2, "SELECT * FROM city WHERE ID='4080'"); errors.Is(err, sql.ErrNoRows) {
			log.Printf("Already no such city ID = %s", "4080")
		} else if err != nil {
			log.Fatalf("DB Error: %s", err)
		} else {
			del := fmt.Sprintf("DELETE FROM city WHERE ID = 4080;")
			if _, err := db.Exec(del); err != nil {
				log.Fatalf("Exec error: %s", err)
			}
		}
	}

	if err := db.Get(&city2, "SELECT * FROM city WHERE ID='4080'"); errors.Is(err, sql.ErrNoRows) {
		log.Printf("no such ID = %s", "4080")
		city2 = City{}
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}
	fmt.Println(city2)
}
