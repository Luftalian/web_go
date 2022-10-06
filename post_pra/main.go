package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type jsonData struct {
	Number int    `json:"number,omitempty"`
	String string `json:"string,omitempty"`
	Bool   bool   `json:"bool,omitempty"`
}

// type studentData struct {
// 	Student_number int    `json:"student_number,omitempty"`
// 	Name           string `json:"name,omitempty"`
// }

// type classData struct {
// 	Class_number int           `json:"class_number,omitempty"`
// 	Students     []studentData `json:"students,omitempty"`
// }

// type schoolData struct {
// 	Class1 []classData `json:"class1,omitempty"`
// 	Class2 []classData `json:"class2,omitempty"`
// 	Class3 []classData `json:"class3,omitempty"`
// 	Class4 []classData `json:"class4,omitempty"`
// }

// type returnData struct {
// 	Class_number int    `json:"class_number,omit"`
// 	Name         string `json:"name,omitempty"`
// 	Error        string `json:"error,omitempty"`
// }

func main() {
	e := echo.New()

	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World.\n")
	})

	e.GET("/Luftalian", func(c echo.Context) error {
		return c.String(http.StatusOK, "Luftalianです。\nインサイダー取引がしたいです。")
	})

	e.GET("/json", jsonHandler)

	e.POST("/post", postHandler)

	e.GET("/hello/:username/:ID", helloHandler)

	e.GET("/ping", pingHandler)

	e.GET("/fizzbuzz", fizzbuzzHandler)

	// e.GET("/students/:class/:studentNumber", studentsHandler)

	e.Logger.Fatal(e.Start(":8080"))
	// ここを前述の通り自分のポートにすること(例: e.Start(":10100"))
}

func jsonHandler(c echo.Context) error {
	res := jsonData{
		Number: 10,
		String: "hoge",
		Bool:   false,
	}

	return c.JSON(http.StatusOK, &res)
}

func postHandler(c echo.Context) error {
	data := &jsonData{}
	err := c.Bind(data)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("%+v", data))
	}
	return c.JSON(http.StatusOK, data)
}

func helloHandler(c echo.Context) error {
	userID := c.Param("username")
	ID := c.Param("ID")
	return c.String(http.StatusOK, "Hello, "+userID+","+ID+".\n")
}

func pingHandler(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func fizzbuzzHandler(c echo.Context) error {
	count := c.QueryParam("count")
	intCount, err := strconv.Atoi(count)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	var returnText string
	for i := 1; i <= intCount; i++ {
		if i%3 == 0 {
			if i%5 == 0 {
				returnText += "FizzBuzz"
			} else {
				returnText += "Fizz"
			}
		} else {
			if i%5 == 0 {
				returnText += "Buzz"
			} else {
				returnText += strconv.Itoa(i)
			}
		}
		returnText += "\n"
	}
	return c.String(http.StatusOK, returnText)
}

// func studentsHandler(c echo.Context) error {
// 	class := c.Param("class")
// 	studentNumber := c.Param("studentNumber")
// 	res := []byte(`
//         [
//             {"class_number": 1, "students": [
//                 {"student_number": 1, "name": "hijiki51"},
//                 {"student_number": 2, "name": "logica"},
//                 {"student_number": 3, "name": "Ras"}
//             ]},
//             {"class_number": 2, "students": [
//                 {"student_number": 1, "name": "asari"},
//                 {"student_number": 2, "name": "irori"},
//                 {"student_number": 3, "name": "itt"},
//                 {"student_number": 4, "name": "mehm8128"}
//             ]},
//             {"class_number": 3, "students": [
//                 {"student_number": 1, "name": "reyu"},
//                 {"student_number": 2, "name": "yukikurage"},
//                 {"student_number": 3, "name": "anko"}
//             ]},
//             {"class_number": 4, "students": [
//                 {"student_number": 1, "name": "Uzaki"},
//                 {"student_number": 2, "name": "yashu"}
//             ]}
//         ]
//     `)
// 	var returnText schoolData
// 	if err := json.Unmarshal(res, &returnText); err != nil {
// 		panic(err)
// 	}

// 	classInfo := returnText.Class1
// 	for i := 0; i < classInfo.length; i++ {
// 	}
// 	name := classInfo
// 	answer := &returnData
// 	if name == nil || answer == nil {
// 		answer["Error"] = "Student Not Found"
// 		return c.JSON(http.StatusNotFound, answer)
// 	}
// 	answer["Class_number"] = class
// 	answer["Student_number"] = name
// 	return c.JSON(http.StatusOK, answer)
// }
