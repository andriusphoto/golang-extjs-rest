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
		UseDB,
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
		q := c.Get("q").(r.Term)

		res, err := q.Table("users").Filter(map[string]interface{}{
			"username": username,
			"password": password,
		}).Run(c.Get("session").(r.QueryExecutor))
		if err != nil {
			fmt.Print(err)

			return nil
		}
		var users []map[string]interface{}
		err = res.All(&users)
		if len(users) > 0 {
			token, err := auth.NewJWT(jwt.MapClaims{
				"username": users[0]["username"],
			}, signingKey)
			if err != nil {
				fmt.Print(err)

				return nil
			}
			ret = jsonReturn{token, true}
		} else {
			ret = jsonReturn{nil, false}
		}

		json, _ := json.Marshal(ret)

		return c.Write(string(json))
	})

	// api.Options("/<table>", conect, useTable, Get)

	api.Get("/restricted", auth.JWT(signingKey), getJWTclaims, restricted)
	api.Get("/<table>", UseTable, AddFilter, Total, AddSorter, AddPagination, Get)
	api.Get("/<table>/<id>", auth.JWT(signingKey), UseTable, GetOne)
	api.Delete("/<table>/<id>", auth.JWT(signingKey), UseTable, Delete)
	api.Post("/<table>", auth.JWT(signingKey), UseTable, Create)
	api.Put("/<table>/<id>", auth.JWT(signingKey), UseTable, Update)
	print("Server started... port:7776")
	http.Handle("/", router)
	http.ListenAndServe(":7776", nil)

}
func getJWTclaims(c *routing.Context) error {
	claims := c.Get("JWT").(*jwt.Token).Claims.(jwt.MapClaims)

	c.Set("claims", claims)
	return nil

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
func UseDB(c *routing.Context) error {

	c.Set("q", r.DB("test"))

	return nil
}
func UseTable(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	c.Set("q", q.Table(c.Param("table")))

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

		return nil
	}
	start, err := strconv.Atoi(c.Request.FormValue("start"))
	if err != nil {

		return nil
	}
	limit, err := strconv.Atoi(c.Request.FormValue("limit"))
	if err != nil {

		return nil
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
