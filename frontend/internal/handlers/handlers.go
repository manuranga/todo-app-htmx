package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	common "github.com/charukak/todo-app-htmx/common/pkg"
	"github.com/charukak/todo-app-htmx/frontend/internal/app"
	"github.com/charukak/todo-app-htmx/frontend/pkg/log"
	"github.com/charukak/todo-app-htmx/frontend/web/templates"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	todoAppClient *app.TodoAppClient
}

func (h *Handler) Hello(w http.ResponseWriter, r *http.Request) {
	component := templates.Hello("yay")
	component.Render(context.Background(), w)
}

func (h *Handler) TodoList(w http.ResponseWriter, r *http.Request) {
	log.Info("GET TodoList")
	todos, err := h.todoAppClient.GetTodos()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	component := templates.TodoList(todos)
	component.Render(context.Background(), w)
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	log.Info("POST CreateTodo")

	title := r.FormValue("title")
	description := r.FormValue("description")

	err := h.todoAppClient.CreateTodo(title, description)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "todos-changed")

	// component := templates.TodoItem(*todo)

}

func (h *Handler) UpdateTodoStatus(w http.ResponseWriter, r *http.Request) {
	log.Info("POST UpdateTodoStatus")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status := false
	if r.FormValue("status") == "on" {
		status = true
	}

	fmt.Println(status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.todoAppClient.UpdateTodo(&common.Todo{
		ID:     id,
		Status: status,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "todos-changed")
}

func (h *Handler) DeleteTodoById(w http.ResponseWriter, r *http.Request) {
	log.Info("DELETE DeleteTodoById")

	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.todoAppClient.DeleteTodoById(strconv.Itoa(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", "todos-changed")
}

func NewHandler() *Handler {

	return &Handler{
		todoAppClient: app.NewTodoAppClient(os.Getenv("TODO_APP_URL")),
	}
}
