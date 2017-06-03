package lib

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-ozzo/ozzo-routing"
	r "gopkg.in/gorethink/gorethink.v3"
)

func GetJWTclaims(c *routing.Context) error {
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
	sorter := []Sorter{}
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
	sorter := []Sorter{}
	json.Unmarshal([]byte(sortstr), &sorter)

	q = q.Slice(start, limit*page)

	c.Set("q", q)
	return nil
}
func AddFilter(c *routing.Context) error {
	q := c.Get("q").(r.Term)
	filterstr := c.Request.FormValue("filter")
	filter := []Filter{}
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
func ParseBasicAuth(auth string) (username, password string) {
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

type JsonReturnArray struct {
	Data    []interface{} `json:"data"`
	Success bool          `json:"success"`
	Length  int           `json:"total"`
}

type JsonReturn struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}
type JsonReturnError struct {
	Errors  interface{} `json:"errors"`
	Success bool        `json:"success"`
}
type Filter struct {
	Property string `json:"property"`
	Value    string `json:"value"`
	Operator string `json:"operator"`
}
type Sorter struct {
	Property  string `json:"property"`
	Direction string `json:"direction"`
}
