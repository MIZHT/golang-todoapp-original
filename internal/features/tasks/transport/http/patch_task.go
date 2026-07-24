package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/MIZHT/golang-todoapp/internal/core/domain"
	core_logger "github.com/MIZHT/golang-todoapp/internal/core/logger"
	core_http_request "github.com/MIZHT/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/MIZHT/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/MIZHT/golang-todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"купить елку"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"сделать до рождества"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean" example:"false"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must b e between 1 and 100 sumbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 100 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be NULL")
		}
	}
	return nil
}

type PatchUserResponse TaskDTOResponse

// PatchTask 	godoc
// @Summary		Изменение задачи
// @Description Изменение информации об уже существующем задачи
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**: `description` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**:  `"description": "Вечером пойти за кладом"` - устанавливает новое описание задачи
// @Description 3. **Передан null**: `"description": null` - очищает поле в БД(set to NULL)
// @Description ОГРАНИЧЕНИЯ: `title` и `completed` не может быть выставлен как null
// @Tags 		tasks
// @Produce 	json
// @Accept		json
// @Param 		id path int true "ID изменяемой задачи"
// @Param 		request body PatchTaskRequest true "PatchTask тело запроса"
// @Success 	200		{object} PatchUserResponse "Успешно изменена задача"
// @Failure 	400		{object} core_http_response.ErrorResponse "Bad request"
// @Failure 	404		{object} core_http_response.ErrorResponse "User not found"
// @Failure 	500		{object} core_http_response.ErrorResponse "Internal server error"
// @Router 		/tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)
		return
	}
	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	taskPatch := taskPatcFromhRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)

		return
	}
	response := PatchUserResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatcFromhRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomin(),
		request.Description.ToDomin(),
		request.Completed.ToDomin(),
	)
}
