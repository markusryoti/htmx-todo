package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var todos []*Todo

type Todo struct {
	Id   string `param:"id"`
	Name string `form:"name"`
	Done bool
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
		Name:  "A killer htmx todo app",
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

	t.Id = uuid.New().String()

	todos = append(todos, t)

	return c.Render(http.StatusOK, "edited-todo.html", t)
}

func toggle(c echo.Context) error {
	t := new(Todo)

	err := c.Bind(t)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	for _, todo := range todos {
		if todo.Id == t.Id {
			todo.Done = !todo.Done
			t = todo
		}
	}

	return c.Render(http.StatusOK, "edited-todo.html", t)
}

func remove(c echo.Context) error {
	t := new(Todo)

	err := c.Bind(t)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	updatedTodos := make([]*Todo, 0)

	for _, todo := range todos {
		if todo.Id != t.Id {
			updatedTodos = append(updatedTodos, todo)
		}
	}

	todos = updatedTodos

	return nil
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()
	e.Renderer = t

	e.Static("/public", "public")

	e.GET("/", home)
	e.POST("/todos/add", add)
	e.PATCH("/todos/:id", toggle)
	e.DELETE("/todos/:id", remove)

	e.Logger.Fatal(e.Start(":8080"))
}
