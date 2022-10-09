package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/srinathgs/mysqlstore"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type City struct {
	ID          int            `json:"id,omitempty"  db:"ID"`
	Name        sql.NullString `json:"name,omitempty"  db:"Name"`
	CountryCode sql.NullString `json:"countryCode,omitempty"  db:"CountryCode"`
	District    sql.NullString `json:"district,omitempty"  db:"District"`
	Population  sql.NullInt64  `json:"population,omitempty"  db:"Population"`
}

type Country struct {
	Code           sql.NullString `json:"code,omitempty" db:"Code"`
	Name           sql.NullString `json:"name,omitempty" db:"Name"`
	Continent      sql.NullString `json:"continent,omitempty" db:"Continent"`
	Region         sql.NullString `json:"region,omitempty" db:"Region"`
	SurfaceArea    sql.NullString `json:"surfaceArea,omitempty" db:"SurfaceArea"`
	IndepYear      int            `json:"indepYear,omitempty" db:"IndepYear"`
	Population     float32        `json:"population,omitempty" db:"Population"`
	LifeExpectancy float32        `json:"lifeExpectancy,omitempty" db:"LifeExpectancy"`
	GNP            float32        `json:"GNP,omitempty" db:"GNP"`
	GNPOld         float32        `json:"GNPOld,omitempty" db:"GNPOld"`
	LocalName      sql.NullString `json:"localName,omitempty" db:"LocalName"`
	GovernmentForm sql.NullString `json:"governmentForm,omitempty" db:"GovernmentForm"`
	HeadOfState    sql.NullString `json:"headOfState,omitempty" db:"HeadOfState"`
	Capital        int            `json:"capital,omitempty" db:"Capital"`
	Code2          sql.NullString `json:"code2,omitempty"  db:"Code2"`
}

var (
	db *sqlx.DB
)

func main() {
	_db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}
	db = _db

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	e.POST("/login", postLoginHandler)
	// e.GET("/logout", logoutHandler)
	e.POST("/signup", postSignUpHandler)

	withLogin := e.Group("")
	withLogin.Use(checkLogin)
	withLogin.GET("/cities/:cityName", getCityInfoHandler)
	withLogin.GET("/countries", getCountryInfoHandler)
	withLogin.GET("/countriesss/:countryCode", getCountryCitiesInfoHandler)
	withLogin.GET("/whoami", getWhoAmIHandler)

	e.GET("/check", checkHandler)

	e.Start(":8088")
}

type LoginRequestBody struct {
	Username string `json:"username,omitempty"  form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type User struct {
	Username   string `json:"username,omitempty" db:"Username"`
	HashedPass string `json:"-" db:"HashedPass"`
}

type Me struct {
	Username string `json:"username,omitempty" db:"username"`
}

type Name_Code struct {
	Name string `json:"name,omitempty" db:"country"`
	Code string `json:"country,omitempty" db:"code"`
}

func postSignUpHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	// もう少し真面目に
	if req.Password == "" || req.Username == "" {
		// エラーは真面目に
		return c.String(http.StatusBadRequest, "項目が空です")
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("bcrypt generate error: %v", err))
	}

	var count int

	err = db.Get(&count, "SELECT COUNT(*) FROM users WHERE Username=?", req.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	if count > 0 {
		return c.String(http.StatusConflict, "ユーザーがすでに存在しています")
	}

	_, err = db.Exec("INSERT INTO users (Username, HashedPass) VALUES (?, ?)", req.Username, hashedPass)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}
	return c.NoContent(http.StatusCreated)
}

func postLoginHandler(c echo.Context) error {
	req := LoginRequestBody{}
	c.Bind(&req)

	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE username=?", req.Username)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return c.NoContent(http.StatusForbidden)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["userName"] = req.Username
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

// func logoutHandler(c echo.Context) error {
// 	req := LoginRequestBody{}
// 	c.Bind(&req)

// 	myUsername := getWhoAmIHandler.Me

// 	user := User{}
// 	err := db.Get(&user, "SELECT * FROM users WHERE username=?", myUsername)
// 	if err != nil {
// 		return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
// 	}

// 	// err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password))
// 	// if err != nil {
// 	// 	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
// 	// 		return c.NoContent(http.StatusForbidden)
// 	// 	} else {
// 	// 		return c.NoContent(http.StatusInternalServerError)
// 	// 	}
// 	// }

// 	sess, err := session.Get("sessions", c)
// 	if err != nil {
// 		fmt.Println(err)
// 		return c.String(http.StatusInternalServerError, "something wrong in getting session")
// 	}
// 	// sess.Values["userName"] = req.Username
// 	sess.Values = remove(sess.Values, myUsername)
// 	sess.Save(c.Request(), c.Response())

// 	return c.NoContent(http.StatusOK)
// }

func checkLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}

		if sess.Values["userName"] == nil {
			return c.String(http.StatusForbidden, "please login")
		}
		c.Set("userName", sess.Values["userName"].(string))

		return next(c)
	}
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")

	city := City{}
	db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName)
	if !city.Name.Valid {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, city)
}

func getCountryCitiesInfoHandler(c echo.Context) error {
	countryCode := c.Param("countryCode")

	city := City{}
	var cityNames []Name_Code
	rows, _ := db.Queryx("SELECT * FROM city WHERE CountryCode=?", countryCode)
	for rows.Next() {
		err := rows.StructScan(&city)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}
		name_code := Name_Code{city.Name.String, "city.ID.String"}
		cityNames = append(cityNames, name_code)
	}

	return c.JSON(http.StatusOK, cityNames)
}

func getCountryInfoHandler(c echo.Context) error {
	country := Country{}
	var countryNames []Name_Code
	rows, _ := db.Queryx("SELECT Name, Code FROM country")
	for rows.Next() {
		err := rows.StructScan(&country)
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("db error: %v", err))
		}
		name_code := Name_Code{country.Name.String, country.Code.String}
		countryNames = append(countryNames, name_code)
	}

	return c.JSON(http.StatusOK, countryNames)
}

func checkHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	log.Println(sess)
	return c.JSON(http.StatusOK, sess.Values["userName"])
}

func getWhoAmIHandler(c echo.Context) error {
	sess, _ := session.Get("sessions", c)

	return c.JSON(http.StatusOK, Me{
		Username: sess.Values["userName"].(string),
	})
}
