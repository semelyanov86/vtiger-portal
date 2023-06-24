package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/service"
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"net/http"
)

func (h *Handler) initProjectsRoutes(api *gin.RouterGroup) {
	projects := api.Group("/projects")
	{
		projects.GET("/", h.getAllProjects)
		projects.GET("/:id", h.getProject)
		projects.GET("/:id/comments", h.getProjectComments)
		projects.POST("/:id/comments", h.addProjectComment)
		projects.GET("/:id/documents", h.getProjectDocuments)
		projects.POST("/:id/documents", h.uploadProjectDocument)
		projects.GET("/:id/file/:file", h.getProjectFile)
		projects.GET("/:id/tasks", h.getAllProjectTasks)
		projects.PATCH("/:id/tasks/:task", h.updateProjectTask)
		projects.GET("/:id/tasks/:task", h.getProjectTask)
		projects.POST("/:id/tasks", h.createProjectTask)
		projects.GET("/:id/tasks/:task/comments", h.getProjectTaskComments)
		projects.POST("/:id/tasks/:task/comments", h.addProjectTaskComment)
		projects.GET("/:id/tasks/:task/documents", h.getProjectTaskDocuments)
		projects.POST("/:id/tasks/:task/documents", h.uploadProjectTaskDocument)
		projects.GET("/:id/tasks/:task/file/:file", h.getTaskFile)
	}
}

func (h *Handler) getProject(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" {
		return
	}

	project, err := h.services.Projects.GetProjectById(c.Request.Context(), id, true, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	res := AloneDataResponse[domain.Project]{
		Data: project,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getProjectTask(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	taskId := h.getAndValidateId(c, "task")

	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" || taskId == "" {
		newResponse(c, http.StatusForbidden, "wrong user or id")
		return
	}

	_, err := h.services.Projects.GetProjectById(c.Request.Context(), id, true, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	taskModel, err := h.services.ProjectTasks.GetProjectTaskById(c.Request.Context(), taskId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
	}
	res := AloneDataResponse[domain.ProjectTask]{
		Data: taskModel,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getAllProjects(c *gin.Context) {
	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		return
	}
	sortString := c.DefaultQuery("sort", "-project_no")

	if !isSortStringValid(sortString, []string{"project_no", "projectname", "projecttype", "projectstatus", "projectpriority"}) {
		newResponse(c, http.StatusUnprocessableEntity, "sort value "+sortString+" is not allowed")
		return
	}

	projects, count, err := h.services.Projects.GetAll(c.Request.Context(), vtiger.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
		Contact:  userModel.Crmid,
		Sort:     sortString,
		Search:   c.DefaultQuery("search", ""),
	})
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Project]{
		Data:  projects,
		Count: count,
		Page:  page,
		Size:  size,
	})
}

func (h *Handler) getProjectComments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}

	comments, err := h.services.Projects.GetRelatedComments(c.Request.Context(), id, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Comment]{
		Data:  comments,
		Count: len(comments),
		Page:  1,
		Size:  100,
	})
}

func (h *Handler) getProjectDocuments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}

	documents, err := h.services.Projects.GetRelatedDocuments(c.Request.Context(), id, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Document]{
		Data:  documents,
		Count: len(documents),
		Page:  1,
		Size:  100,
	})
}

func (h *Handler) uploadProjectDocument(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		notPermittedResponse(c)
		return
	}
	_, err := h.services.Projects.GetProjectById(c.Request.Context(), id, false, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	document, err := h.services.Documents.AttachFile(c.Request.Context(), file, id, *userModel, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AloneDataResponse[domain.Document]{Data: document})
}

func (h *Handler) uploadProjectTaskDocument(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	taskId := h.getAndValidateId(c, "task")
	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil || taskId == "" {
		notPermittedResponse(c)
		return
	}
	_, err := h.services.Projects.GetProjectById(c.Request.Context(), id, false, userModel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	document, err := h.services.Documents.AttachFile(c.Request.Context(), file, taskId, *userModel, header)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AloneDataResponse[domain.Document]{Data: document})
}

func (h *Handler) getProjectFile(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	fileId := h.getAndValidateId(c, "file")

	if id == "" || userModel == nil || fileId == "" {
		return
	}

	_, err := h.services.Projects.GetProjectById(c.Request.Context(), id, false, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	file, err := h.services.Documents.GetFile(c.Request.Context(), fileId, id)

	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := AloneDataResponse[vtiger.File]{
		Data: file,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) getTaskFile(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	taskId := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	fileId := h.getAndValidateId(c, "file")

	if id == "" || userModel == nil || fileId == "" || taskId == "" {
		return
	}

	_, err := h.services.Projects.GetProjectById(c.Request.Context(), id, false, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	file, err := h.services.Documents.GetFile(c.Request.Context(), fileId, taskId)

	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}

	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := AloneDataResponse[vtiger.File]{
		Data: file,
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) addProjectComment(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)

	if id == "" || userModel == nil {
		return
	}
	var inp createCommentInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}
	comment, err := h.services.Projects.AddComment(c.Request.Context(), inp.Commentcontent, id, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := AloneDataResponse[domain.Comment]{
		Data: comment,
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) getAllProjectTasks(c *gin.Context) {
	id := h.getAndValidateId(c, "id")

	userModel := h.getValidatedUser(c)
	page, size := h.getPageAndSizeParams(c)

	if userModel == nil || page < 0 || size < 0 {
		return
	}

	projectTasks, count, err := h.services.ProjectTasks.GetAllFromProject(c.Request.Context(), vtiger.PaginationQueryFilter{
		Page:     page,
		PageSize: size,
		Client:   userModel.AccountId,
		Contact:  userModel.Crmid,
		Parent:   id,
	}, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.ProjectTask]{
		Data:  projectTasks,
		Count: count,
		Page:  page,
		Size:  size,
	})
}

func (h *Handler) createProjectTask(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	var inp service.ProjectTaskInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}
	userModel := h.getValidatedUser(c)
	if id == "" || userModel == nil {
		newResponse(c, http.StatusBadRequest, "Wrong project or auth user")
		return
	}

	_, err := h.services.Projects.GetProjectById(c.Request.Context(), id, false, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	task, err := h.services.ProjectTasks.CreateProjectTask(c.Request.Context(), inp, id)
	if errors.Is(service.ErrValidation, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "projecttaskname", "message": err.Error()})
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *Handler) updateProjectTask(c *gin.Context) {
	var inp map[string]any
	if err := c.ShouldBindJSON(&inp); err != nil {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "projectstatus", "message": "Incorrect value"})
		return // exit on first error
	}
	id := h.getAndValidateId(c, "id")
	taskId := h.getAndValidateId(c, "task")
	userModel := h.getValidatedUser(c)
	if userModel == nil || id == "" || taskId == "" {
		return
	}

	task, err := h.services.ProjectTasks.Revise(c.Request.Context(), inp, taskId, id, *userModel)
	if errors.Is(service.ErrValidation, err) {
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": "projecttaskstatus", "message": err.Error()})
		return
	}
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusAccepted, task)
}

func (h *Handler) getProjectTaskComments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	taskId := h.getAndValidateId(c, "task")

	userModel := h.getValidatedUser(c)

	if id == "" || taskId == "" || userModel == nil {
		return
	}

	comments, err := h.services.ProjectTasks.GetRelatedComments(c.Request.Context(), taskId, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Comment]{
		Data:  comments,
		Count: len(comments),
		Page:  1,
		Size:  100,
	})
}

func (h *Handler) addProjectTaskComment(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	taskId := h.getAndValidateId(c, "task")

	userModel := h.getValidatedUser(c)

	if id == "" || taskId == "" || userModel == nil {
		return
	}
	var inp createCommentInput
	if err := c.BindJSON(&inp); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": "Validation Error", "field": fieldErr.Field(), "message": fieldErr.Error()})
			return // exit on first error
		}
	}
	comment, err := h.services.ProjectTasks.AddComment(c.Request.Context(), inp.Commentcontent, taskId, *userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	res := AloneDataResponse[domain.Comment]{
		Data: comment,
	}
	c.JSON(http.StatusCreated, res)
}

func (h *Handler) getProjectTaskDocuments(c *gin.Context) {
	id := h.getAndValidateId(c, "id")
	taskId := h.getAndValidateId(c, "task")

	userModel := h.getValidatedUser(c)

	if id == "" || taskId == "" || userModel == nil {
		return
	}

	documents, err := h.services.ProjectTasks.GetRelatedDocuments(c.Request.Context(), taskId, userModel)
	if errors.Is(service.ErrOperationNotPermitted, err) {
		notPermittedResponse(c)
		return
	}
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, DataResponse[domain.Document]{
		Data:  documents,
		Count: len(documents),
		Page:  1,
		Size:  100,
	})
}
