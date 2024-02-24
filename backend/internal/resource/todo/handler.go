package todo

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	common "github.com/charukak/todo-app-htmx/common/pkg"
	"github.com/gin-gonic/gin"
)

type ITodoService interface {
	GetAll() func(c *gin.Context)
	GetByID() func(c *gin.Context)
	Create() func(c *gin.Context)
	Update() func(c *gin.Context)
	Delete() func(c *gin.Context)
}

type TodoService struct {
	db *sql.DB
}

func (s *TodoService) GetAll() func(c *gin.Context) {
	return func(c *gin.Context) {

		rows, err := s.db.Query("SELECT * FROM todos")

		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		defer rows.Close()

		var todos []common.Todo = []common.Todo{}

		for rows.Next() {
			var todo common.Todo

			err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status)

			if err != nil {
				c.JSON(500, gin.H{"message": err.Error()})
				return
			}

			todos = append(todos, todo)
		}

		c.JSON(200, todos)
	}
}

func (s *TodoService) GetByID() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		var todo common.Todo

		row := s.db.QueryRow("SELECT * FROM todos WHERE id = ?", id)

		err = row.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status)

		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		c.JSON(200, todo)
	}
}

func (s *TodoService) Create() func(c *gin.Context) {
	return func(c *gin.Context) {
		var todo common.Todo

		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		result, err := s.db.Exec(
			"INSERT INTO todos (title, description, status) VALUES (?, ?, ?)",
			todo.Title,
			todo.Description,
			todo.Status)

		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		lastID, err := result.LastInsertId()

		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		todo.ID = int(lastID)

		c.Header("Location", fmt.Sprintf("/todos/%d", todo.ID))
		c.Status(201)
	}
}

func (s *TodoService) Update() func(c *gin.Context) {
	return func(c *gin.Context) {
		var todo common.Todo

		if err := c.ShouldBindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		todo.ID = id

		cols := []string{}

		if todo.Title != "" {
			cols = append(cols, "title = '"+todo.Title+"'")
		}

		if todo.Description != "" {
			cols = append(cols, "description = '"+todo.Description+"'")
		}

		cols = append(cols, "status = "+strconv.FormatBool(todo.Status))

		_, err = s.db.Exec(
			fmt.Sprintf("UPDATE todos SET %s WHERE id = ?", strings.Join(cols, ",")),
			todo.ID)

		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		c.Header("Location", fmt.Sprintf("/todos/%d", todo.ID))
		c.Status(204)
	}
}

func (s *TodoService) Delete() func(c *gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		_, err = s.db.Exec("DELETE FROM todos WHERE id = ?", id)

		if err != nil {
			c.JSON(500, gin.H{"message": err.Error()})
			return
		}

		c.Status(204)
	}
}

func NewTodoService(db *sql.DB) ITodoService {
	return &TodoService{db}
}

func RegisterTodoHandler(r *gin.Engine, db *sql.DB) {
	service := NewTodoService(db)

	todoRoutes := r.Group("/todos")
	{
		todoRoutes.GET("", service.GetAll())
		todoRoutes.GET(":id", service.GetByID())
		todoRoutes.POST("", service.Create())
		todoRoutes.PUT(":id", service.Update())
		todoRoutes.DELETE(":id", service.Delete())
	}

}
