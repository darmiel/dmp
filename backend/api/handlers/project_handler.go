package handlers

import (
	"errors"
	"github.com/darmiel/dmp/api/presenter"
	"github.com/darmiel/dmp/api/services"
	"github.com/darmiel/dmp/pkg/model"
	"github.com/darmiel/dmp/pkg/util"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	gofiberfirebaseauth "github.com/ralf-life/gofiber-firebaseauth"
	"go.uber.org/zap"
)

var ErrNoAccess = errors.New("no access")
var ErrNotFound = errors.New("not found")

type ProjectHandler struct {
	srv       services.ProjectService
	logger    *zap.SugaredLogger
	validator *validator.Validate
}

func NewProjectHandler(
	srv services.ProjectService,
	logger *zap.SugaredLogger,
	validator *validator.Validate,
) *ProjectHandler {
	return &ProjectHandler{srv, logger, validator}
}

type projectDto struct {
	Name        string `validate:"required,proj-extended,min=3,max=36,startsnotwith= ,endsnotwith= " json:"name,omitempty"`
	Description string `validate:"proj-extended,max=256" json:"description,omitempty"`
}

// ProjectAccessMiddleware checks if the requesting user has access to the project
// and preloads the project data to the "project" local
func (h *ProjectHandler) ProjectAccessMiddleware(ctx *fiber.Ctx) error {
	u := ctx.Locals("user").(gofiberfirebaseauth.User)
	projectID, err := ctx.ParamsInt("project_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse(err))
	}

	// find project to check if the requester is the owner
	project, err := h.srv.FindProject(uint(projectID), "Users")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse(err))
	}

	// owner always has access
	if project.OwnerID == u.UserID || util.HasAccess(project, u.UserID) {
		ctx.Locals("project", *project)
		return ctx.Next()
	}

	return ctx.Status(fiber.StatusUnauthorized).JSON(presenter.ErrorResponse(ErrNoAccess))
}

// AddProject creates a new project
func (h *ProjectHandler) AddProject(ctx *fiber.Ctx) error {
	u := ctx.Locals("user").(gofiberfirebaseauth.User)
	h.logger.Infof("user %s is requesting to create a new project", u.UserID)

	var payload projectDto
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse(err))
	}
	if err := h.validator.Struct(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse(err))
	}

	project, err := h.srv.CreateProject(payload.Name, payload.Description, u.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse(err))
	}

	h.logger.Infof("created project %d for user %s (%s)", project.ID, u.UserID, project.Name)
	return ctx.Status(fiber.StatusCreated).JSON(presenter.SuccessResponse("project created", project))
}

// GetProjects returns a list of all projects the user has access to
func (h *ProjectHandler) GetProjects(ctx *fiber.Ctx) error {
	u := ctx.Locals("user").(gofiberfirebaseauth.User)
	h.logger.Infof("user %s is requesting a list of projects", u.UserID)

	projects := make(map[uint]model.Project)

	ownerProjects, err := h.srv.FindProjectsByOwner(u.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse(err))
	}
	for _, p := range ownerProjects {
		projects[p.ID] = p
	}

	userProjects, err := h.srv.FindProjectsByUserAccess(u.UserID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse(err))
	}
	for _, p := range userProjects {
		projects[p.ID] = p
	}

	result := make([]model.Project, 0, len(projects))
	for _, v := range projects {
		result = append(result, v)
	}
	return ctx.Status(fiber.StatusOK).JSON(presenter.SuccessResponse("", result))
}

// DeleteProject deletes a project
func (h *ProjectHandler) DeleteProject(ctx *fiber.Ctx) error {
	p := ctx.Locals("project").(model.Project)
	if err := h.srv.DeleteProject(p.ID); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse(err))
	}
	return ctx.Status(fiber.StatusOK).JSON(presenter.SuccessResponse("project deleted", nil))
}

// EditProject edits the name and description of a project
func (h *ProjectHandler) EditProject(ctx *fiber.Ctx) error {
	var payload projectDto
	if err := ctx.BodyParser(&payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse(err))
	}
	h.logger.Infof("Validating %+v :: %v", payload, h.validator.Struct(payload))
	if err := h.validator.Struct(payload); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(presenter.ErrorResponse(err))
	}
	p := ctx.Locals("project").(model.Project)
	if err := h.srv.EditProject(p.ID, payload.Name, payload.Description); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(presenter.ErrorResponse(err))
	}
	return ctx.Status(fiber.StatusOK).JSON(presenter.SuccessResponse("project updated", nil))
}
