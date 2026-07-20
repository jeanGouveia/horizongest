package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jeanGouveia/pratoOnline/backend/internal/service"
)

type PlanHandler struct {
	planService *service.PlanService
}

func NewPlanHandler(planService *service.PlanService) *PlanHandler {
	return &PlanHandler{planService: planService}
}

// POST /api/platform/plans
func (h *PlanHandler) CreatePlan(w http.ResponseWriter, r *http.Request) {
	var input service.CreatePlanInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}
	if err := validatePlatform.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	plan, err := h.planService.CreatePlan(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrPlanAlreadyExists) {
			jsonError(w, "plano com este slug já existe", http.StatusConflict)
			return
		}
		jsonError(w, "não foi possível criar plano", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, plan)
}

// GET /api/platform/plans
func (h *PlanHandler) ListPlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.planService.ListPlans(r.Context())
	if err != nil {
		jsonError(w, "não foi possível listar planos", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"plans": plans,
	})
}

// GET /api/platform/plans/active
func (h *PlanHandler) ListActivePlans(w http.ResponseWriter, r *http.Request) {
	plans, err := h.planService.ListActivePlans(r.Context())
	if err != nil {
		jsonError(w, "não foi possível listar planos ativos", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"plans": plans,
	})
}

// GET /api/platform/plans/:id
func (h *PlanHandler) GetPlan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		jsonError(w, "ID do plano inválido", http.StatusBadRequest)
		return
	}

	plan, err := h.planService.GetPlan(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrPlanNotFound) {
			jsonError(w, "plano não encontrado", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível obter plano", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, plan)
}

// PUT /api/platform/plans/:id
func (h *PlanHandler) UpdatePlan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		jsonError(w, "ID do plano inválido", http.StatusBadRequest)
		return
	}

	var input service.UpdatePlanInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		jsonError(w, "formato dos dados inválido", http.StatusBadRequest)
		return
	}
	if err := validatePlatform.Struct(input); err != nil {
		jsonValidationError(w, err)
		return
	}

	plan, err := h.planService.UpdatePlan(r.Context(), uint(id), input)
	if err != nil {
		if errors.Is(err, service.ErrPlanNotFound) {
			jsonError(w, "plano não encontrado", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível atualizar plano", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, plan)
}

// DELETE /api/platform/plans/:id
func (h *PlanHandler) DeletePlan(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		jsonError(w, "ID do plano inválido", http.StatusBadRequest)
		return
	}

	if err := h.planService.DeletePlan(r.Context(), uint(id)); err != nil {
		if errors.Is(err, service.ErrPlanNotFound) {
			jsonError(w, "plano não encontrado", http.StatusNotFound)
			return
		}
		jsonError(w, "não foi possível deletar plano", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{
		"message": "plano deletado com sucesso",
	})
}
