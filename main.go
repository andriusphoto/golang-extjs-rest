package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/auth"
	"github.com/go-ozzo/ozzo-routing/cors"
	"github.com/go-ozzo/ozzo-routing/fault"
	"github.com/go-ozzo/ozzo-routing/slash"
	r "gopkg.in/gorethink/gorethink.v3"
)

type jsonReturnArray struct {
	Data    []interface{} `json:"data"`
	Success bool          `json:"success"`
	Length  int           `json:"total"`
}
type jsonReturn struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

type filter struct {
	Property string `json:"property"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}
type sorter struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
}

func main() {

	router := routing.New()
	signingKey := "secre-key"
	router.Use(
		slash.Remover(http.StatusMovedPermanently),
		fault.Recovery(log.Printf),
		cors.Handler(cors.Options{
			AllowOrigins: "*",
			AllowHeaders: "*",
			AllowMethods: "*",
		}),
		Connect,
		UseTable,
		// func(c *routing.Context) error {
		// 	// id, err := authenticate(c)
		// 	// if err != nil {
		// 	// 	return err
		// 	// }
		// 	username, password := parseBasicAuth(c.Request.Header.Get("Authorization"))
		// 	if username == "demo" && password == "foo" {
		// 		// auth.JWT(signingKey)
		// 	}
		// 	token, err := auth.NewJWT(jwt.MapClaims{
		// 		"id": "10000",
		// 	}, signingKey)
		// 	if err != nil {

		// 		return err

		// 	}

		// 	return c.Write(token)

		// },
	)
	api := router.Group("/api")
	api.Post("/login", func(c *routing.Context) error {
		// id, err := authenticate(c)
		// if err != nil {
		// 	return err
		// }
		var ret interface{}
		username, password := parseBasicAuth(c.Request.Header.Get("Authorization"))
		if username == "demo" && password == "foo" {
			token, err := auth.NewJWT(jwt.MapClaims{
				"id": "10000",
			}, signingKey)
			if err != nil {
				return err
			}
			ret = jsonReturn{token, true}
		} else {
			ret = jsonReturn{nil, false}
		}

		json, _ := json.Marshal(ret)

		return c.Write(string(json))
	})

	// api.Options("/<table>", conect, useTable, Get)

	router.Get("/restricted", func(c *routing.Context) error {
		claims := c.Get("JWT").(*jwt.Token).Claims.(jwt.MapClaims)
		return c.Write(fmt.Sprint("Welcome, %v!", claims["id"]))
	})
	api.Get("/<table>", AddFilter, Total, AddSorter, AddPagination, Get)
	api.Get("/<table>/<id>", GetOne)
	api.Delete("/<table>/<id>", Delete)
	api.Post("/<table>", Create)
	api.Put("/<table>/<id>", Update)
	print("Server started... port:7776")
	http.Handle("/", router)
	http.ListenAndServe(":7776", nil)

}
func Connect(c *routing.Context) error {

	session, err := r.Connect(r.ConnectOpts{
		Address: "localhost",
	})
	c.Set("session", session)
	if err != nil {

		log.Fatalln(err)
	}
	return nil
}
func UseTable(c *routing.Context) error {
	c.Set("q", r.DB("test").Table(c.Param("table")))

	return nil
}
func Total(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	var data []interface{}
	res, err := q.Run(c.Get("session").(r.QueryExecutor))
	err = res.All(&data)
	if err != nil {
		log.Fatalln(err)
	}
	c.Set("total", len(data))
	return nil
}
func AddSorter(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	sortstr := c.Request.FormValue("sort")
	sorter := []sorter{}
	json.Unmarshal([]byte(sortstr), &sorter)

	for _, item := range sorter {
		if item.Direction == "ASC" {
			q = q.OrderBy(item.Property)

		}
	}
	c.Set("q", q)
	return nil
}
func AddPagination(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	page, err := strconv.Atoi(c.Request.FormValue("page"))
	if err != nil {

		log.Fatalln(err)
	}
	start, err := strconv.Atoi(c.Request.FormValue("start"))
	if err != nil {

		log.Fatalln(err)
	}
	limit, err := strconv.Atoi(c.Request.FormValue("limit"))
	if err != nil {

		log.Fatalln(err)
	}
	fmt.Println(limit)
	sortstr := c.Request.FormValue("sort")
	sorter := []sorter{}
	json.Unmarshal([]byte(sortstr), &sorter)

	q = q.Slice(start, limit*page)

	c.Set("q", q)
	return nil
}
func AddFilter(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	filterstr := c.Request.FormValue("filter")
	filter := []filter{}
	json.Unmarshal([]byte(filterstr), &filter)
	for _, item := range filter {
		if item.Operator == "like" {
			q = q.Filter(func(el r.Term) r.Term {
				return el.Field(item.Property).Match(item.Value)
			})
		}

	}
	c.Set("q", q)
	return nil
}
func parseBasicAuth(auth string) (username, password string) {
	if strings.HasPrefix(auth, "Basic ") {
		if bytes, err := base64.StdEncoding.DecodeString(auth[6:]); err == nil {
			str := string(bytes)
			if i := strings.IndexByte(str, ':'); i >= 0 {
				return str[:i], str[i+1:]
			}
		}
	}
	return
}
