package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/labstack/echo/v4"
)

var todos []*Todo = []*Todo{
	{
		Name: "Eat",
	},
	{
		Name: "Ski",
	},
}

type Todo struct {
	Name string `form:"name"`
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Home struct {
	Title string
	Name  string
	Todos []*Todo
}

func home(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", &Home{
		Name:  "My Todos",
		Title: "Home",
		Todos: todos,
	})
}

func add(c echo.Context) error {
	t := new(Todo)

	err := c.Bind(t)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	todos = append(todos, t)

	return c.Render(http.StatusOK, "new.html", t)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()
	e.Renderer = t

	e.Static("/public", "public")

	e.GET("/", home)
	e.POST("/add", add)

	e.Logger.Fatal(e.Start(":8080"))
}
