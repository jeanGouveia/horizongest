package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type ForensicHandler struct {
	mu     sync.Mutex
	events []ForensicEvent
}

type ForensicEvent struct {
	Timestamp  int                    `json:"timestamp"`
	Category   string                 `json:"category"`
	Section    string                 `json:"section"`
	Subsection string                 `json:"subsection,omitempty"`
	Data       map[string]interface{} `json:"data"`
}

type ForensicReport struct {
	Events []ForensicEvent `json:"events"`
}

func NewForensicHandler() *ForensicHandler {
	return &ForensicHandler{
		events: make([]ForensicEvent, 0),
	}
}

func (h *ForensicHandler) RegisterRoutes(r chi.Router) {
	r.Post("/log", h.ReceiveLog)
	r.Post("/generate", h.GenerateReport)
	r.Post("/clear", h.ClearLogs)
}

// ReceiveLog receives forensic events from the frontend
func (h *ForensicHandler) ReceiveLog(w http.ResponseWriter, r *http.Request) {
	var event ForensicEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.events = append(h.events, event)
	h.mu.Unlock()

	w.WriteHeader(http.StatusOK)
	jsonResponse(w, http.StatusOK, map[string]string{"status": "logged"})
}

// GenerateReport generates the forensic report and saves it to forensic-report.txt
func (h *ForensicHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	reportPath := "forensic-report.txt"
	report := h.generateReportContent()

	if err := os.WriteFile(reportPath, []byte(report), 0644); err != nil {
		log.Printf("Error writing forensic report: %v", err)
		http.Error(w, "failed to write report", http.StatusInternalServerError)
		return
	}

	log.Printf("📋 Relatório forense salvo: %s", reportPath)

	jsonResponse(w, http.StatusOK, map[string]string{
		"status":  "success",
		"path":    reportPath,
		"message": fmt.Sprintf("Relatório forense salvo em: %s", reportPath),
	})
}

// ClearLogs clears all forensic events
func (h *ForensicHandler) ClearLogs(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	h.events = make([]ForensicEvent, 0)
	h.mu.Unlock()

	jsonResponse(w, http.StatusOK, map[string]string{"status": "cleared"})
}

func (h *ForensicHandler) generateReportContent() string {
	separator := "================================================================================"

	report := separator + "\n"
	report += "RELATÓRIO FORENSE - INVESTIGAÇÃO DO BUG DE IMPERSONAÇÃO\n"
	report += fmt.Sprintf("Gerado em: %s\n", time.Now().Format("2006-01-02 15:04:05"))
	report += separator + "\n\n"

	// Group events by category
	frontendEvents := make([]ForensicEvent, 0)
	backendEvents := make([]ForensicEvent, 0)

	for _, event := range h.events {
		if event.Category == "FRONTEND" {
			frontendEvents = append(frontendEvents, event)
		} else if event.Category == "BACKEND" {
			backendEvents = append(backendEvents, event)
		}
	}

	// FRONTEND SECTION
	report += separator + "\n"
	report += "1. FRONTEND\n"
	report += separator + "\n\n"

	// 1.1 Initial State
	report += "1.1 Estado inicial\n"
	report += "----------------------------------------\n"
	for _, event := range frontendEvents {
		if event.Section == "INITIAL_STATE" {
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("document.cookie: %v\n", event.Data["documentCookie"])
			report += fmt.Sprintf("localStorage: %v\n", event.Data["localStorage"])
			report += fmt.Sprintf("sessionStorage: %v\n", event.Data["sessionStorage"])
		}
	}
	report += "\n"

	// 1.2 Click on "Entrar" button
	report += "1.2 Clique no botão \"Entrar\"\n"
	report += "----------------------------------------\n"
	for _, event := range frontendEvents {
		if event.Section == "CLICK_ENTRAR" {
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("Empresa escolhida: %v\n", event.Data["companyName"])
			report += fmt.Sprintf("CompanyID: %v\n", event.Data["companyId"])
		}
	}
	report += "\n"

	// 1.3 requestTenantJWT()
	report += "1.3 requestTenantJWT()\n"
	report += "----------------------------------------\n"
	for _, event := range frontendEvents {
		if event.Section == "REQUEST_TENANT_JWT" {
			if event.Subsection == "BEFORE" {
				report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
				report += fmt.Sprintf("URL: %v\n", event.Data["url"])
				report += fmt.Sprintf("Método: %v\n", event.Data["method"])
				report += fmt.Sprintf("Headers enviados: %v\n", event.Data["headers"])
				report += fmt.Sprintf("Authorization enviado: %v\n", event.Data["authorization"])
				report += fmt.Sprintf("Credentials utilizado: %v\n", event.Data["credentials"])
				report += fmt.Sprintf("Body enviado: %v\n", event.Data["body"])
			} else if event.Subsection == "AFTER" {
				report += fmt.Sprintf("\nTimestamp: %dms\n", event.Timestamp)
				report += fmt.Sprintf("Status HTTP: %v\n", event.Data["status"])
				report += fmt.Sprintf("Headers recebidos: %v\n", event.Data["headers"])
				report += fmt.Sprintf("Set-Cookie recebido: %v\n", event.Data["setCookie"])
				report += fmt.Sprintf("Body completo recebido: %v\n", event.Data["body"])
				report += fmt.Sprintf("\nToken JWT recebido: %v\n", event.Data["token"])
				report += fmt.Sprintf("Tamanho: %v bytes\n", event.Data["tokenSize"])
				report += fmt.Sprintf("Payload decodificado: %v\n", event.Data["decodedPayload"])
			}
		}
	}
	report += "\n"

	// 1.4 hydrateContext()
	report += "1.4 hydrateContext()\n"
	report += "----------------------------------------\n"
	for _, event := range frontendEvents {
		if event.Section == "HYDRATE_CONTEXT" {
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("Fase: %s\n", event.Subsection)
			report += fmt.Sprintf("document.cookie: %v\n", event.Data["documentCookie"])
			if event.Data["cookiesSeparated"] != nil {
				report += fmt.Sprintf("Cookies separados: %v\n", event.Data["cookiesSeparated"])
			}
			if event.Data["authTokenFound"] != nil {
				report += fmt.Sprintf("auth_token encontrado?: %v\n", event.Data["authTokenFound"])
			}
			if event.Data["platformAuthTokenFound"] != nil {
				report += fmt.Sprintf("platform_auth_token encontrado?: %v\n", event.Data["platformAuthTokenFound"])
			}
			report += "\n"
		}
	}
	report += "\n"

	// 1.5 Cookie writes
	report += "1.5 Todas as escritas Em document.cookie\n"
	report += "----------------------------------------\n"
	cookieWrites := 0
	for _, event := range frontendEvents {
		if event.Section == "COOKIE_WRITE" {
			cookieWrites++
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("Arquivo: %v\n", event.Data["file"])
			report += fmt.Sprintf("Linha: %v\n", event.Data["line"])
			report += fmt.Sprintf("Valor escrito: %v\n\n", event.Data["value"])
		}
	}
	if cookieWrites == 0 {
		report += "Nenhuma escrita registrada\n"
	}
	report += "\n"

	// 1.6 Cookie removals
	report += "1.6 Todas as remoções de cookie\n"
	report += "----------------------------------------\n"
	cookieRemovals := 0
	for _, event := range frontendEvents {
		if event.Section == "COOKIE_REMOVAL" {
			cookieRemovals++
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("Arquivo: %v\n", event.Data["file"])
			report += fmt.Sprintf("Linha: %v\n", event.Data["line"])
			report += fmt.Sprintf("Cookie removido: %v\n\n", event.Data["cookie"])
		}
	}
	if cookieRemovals == 0 {
		report += "Nenhuma remoção registrada\n"
	}
	report += "\n"

	// 1.7 Navigations
	report += "1.7 Todas as navegações\n"
	report += "----------------------------------------\n"
	navigations := 0
	for _, event := range frontendEvents {
		if event.Section == "NAVIGATION" {
			navigations++
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("Tipo: %v\n", event.Data["type"])
			if event.Data["url"] != nil {
				report += fmt.Sprintf("URL: %v\n", event.Data["url"])
			}
			if event.Data["from"] != nil {
				report += fmt.Sprintf("De: %v\n", event.Data["from"])
			}
			if event.Data["to"] != nil {
				report += fmt.Sprintf("Para: %v\n", event.Data["to"])
			}
			report += "\n"
		}
	}
	if navigations == 0 {
		report += "Nenhuma navegação registrada\n"
	}
	report += "\n"

	// 1.8 Fetch calls
	report += "1.8 Todas as requisições fetch()\n"
	report += "----------------------------------------\n"
	fetchCalls := 0
	for _, event := range frontendEvents {
		if event.Section == "FETCH" {
			fetchCalls++
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("URL: %v\n", event.Data["url"])
			report += fmt.Sprintf("Método: %v\n", event.Data["method"])
			report += fmt.Sprintf("Credentials: %v\n", event.Data["credentials"])
			report += fmt.Sprintf("Headers: %v\n", event.Data["headers"])
			report += fmt.Sprintf("Cookies antes: %v\n", event.Data["cookiesBefore"])
			if event.Data["status"] != nil {
				report += fmt.Sprintf("Status: %v\n", event.Data["status"])
				report += fmt.Sprintf("Headers recebidos: %v\n", event.Data["headersReceived"])
				report += fmt.Sprintf("Body recebido: %v\n", event.Data["bodyReceived"])
			}
			report += "\n"
		}
	}
	if fetchCalls == 0 {
		report += "Nenhum fetch registrado\n"
	}
	report += "\n"

	// BACKEND SECTION
	report += separator + "\n"
	report += "2. BACKEND\n"
	report += separator + "\n\n"

	// 2.1 impersonation/start
	report += "2.1 impersonation/start\n"
	report += "----------------------------------------\n"
	for _, event := range backendEvents {
		if event.Section == "IMPERSONATION_START" {
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("Subseção: %s\n", event.Subsection)
			report += fmt.Sprintf("Dados: %v\n\n", event.Data)
		}
	}
	report += "\n"

	// 2.2 Middleware
	report += "2.2 Middleware\n"
	report += "----------------------------------------\n"
	for _, event := range backendEvents {
		if event.Section == "MIDDLEWARE" {
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("Subseção: %s\n", event.Subsection)
			report += fmt.Sprintf("Dados: %v\n\n", event.Data)
		}
	}
	report += "\n"

	// 2.3 companies/{id}
	report += "2.3 companies/{id}\n"
	report += "----------------------------------------\n"
	for _, event := range backendEvents {
		if event.Section == "COMPANIES_ID" {
			report += fmt.Sprintf("Timestamp: %dms\n", event.Timestamp)
			report += fmt.Sprintf("Subseção: %s\n", event.Subsection)
			report += fmt.Sprintf("Dados: %v\n\n", event.Data)
		}
	}
	report += "\n"

	// SUMMARY SECTION
	report += separator + "\n"
	report += "3. RESUMO\n"
	report += separator + "\n\n"

	report += h.generateSummary(frontendEvents)

	return report
}

func (h *ForensicHandler) generateSummary(events []ForensicEvent) string {
	summary := ""
	stepNumber := 1

	addStep := func(description string) {
		summary += fmt.Sprintf("PASSO %d\n", stepNumber)
		summary += fmt.Sprintf("%s\n\n", description)
		stepNumber++
	}

	// Analyze token lifecycle
	requestBefore := findEvent(events, "REQUEST_TENANT_JWT", "BEFORE")
	requestAfter := findEvent(events, "REQUEST_TENANT_JWT", "AFTER")
	hydrateEvents := filterEvents(events, "HYDRATE_CONTEXT", "")

	addStep("Estado inicial capturado")

	if requestBefore != nil {
		authPresent := requestBefore.Data["authorization"] != nil && requestBefore.Data["authorization"] != ""
		addStep(fmt.Sprintf("requestTenantJWT() iniciado - Token Platform: %v", authPresent))
	}

	if requestAfter != nil {
		tokenReceived := requestAfter.Data["token"] != nil && requestAfter.Data["token"] == true
		if tokenReceived {
			addStep(fmt.Sprintf("Tenant JWT recebido do backend - Tamanho: %v bytes", requestAfter.Data["tokenSize"]))
		} else {
			addStep("Tenant JWT NÃO recebido do backend - FALHA")
		}
	}

	if len(hydrateEvents) > 0 {
		beforeHydrate := findEvent(hydrateEvents, "HYDRATE_CONTEXT", "ANTES")
		afterHydrate := findEvent(hydrateEvents, "HYDRATE_CONTEXT", "IMEDIATAMENTE_DEPOIS")
		after100ms := findEvent(hydrateEvents, "HYDRATE_CONTEXT", "DEPOIS_100MS")
		after300ms := findEvent(hydrateEvents, "HYDRATE_CONTEXT", "DEPOIS_300MS")
		after1000ms := findEvent(hydrateEvents, "HYDRATE_CONTEXT", "DEPOIS_1000MS")

		if beforeHydrate != nil {
			addStep("hydrateContext() - Antes de definir cookie")
		}

		if afterHydrate != nil {
			authFound := afterHydrate.Data["authTokenFound"] != nil && afterHydrate.Data["authTokenFound"] == true
			if authFound {
				addStep("hydrateContext() - Cookie auth_token encontrado IMEDIATAMENTE após escrita")
			} else {
				addStep("hydrateContext() - Cookie auth_token NÃO encontrado IMEDIATAMENTE após escrita - POSSÍVEL PROBLEMA")
			}
		}

		if after100ms != nil {
			authFound := after100ms.Data["authTokenFound"] != nil && after100ms.Data["authTokenFound"] == true
			if authFound {
				addStep("hydrateContext() - Cookie auth_token presente após 100ms")
			} else {
				addStep("hydrateContext() - Cookie auth_token AUSENTE após 100ms - COOKIE REMOVIDO")
			}
		}

		if after300ms != nil {
			authFound := after300ms.Data["authTokenFound"] != nil && after300ms.Data["authTokenFound"] == true
			if authFound {
				addStep("hydrateContext() - Cookie auth_token presente após 300ms")
			} else {
				addStep("hydrateContext() - Cookie auth_token AUSENTE após 300ms - COOKIE REMOVIDO")
			}
		}

		if after1000ms != nil {
			authFound := after1000ms.Data["authTokenFound"] != nil && after1000ms.Data["authTokenFound"] == true
			if authFound {
				addStep("hydrateContext() - Cookie auth_token presente após 1000ms")
			} else {
				addStep("hydrateContext() - Cookie auth_token AUSENTE após 1000ms - COOKIE REMOVIDO")
			}
		}
	}

	// Analyze cookie writes
	cookieWrites := filterEvents(events, "COOKIE_WRITE", "")
	if len(cookieWrites) > 0 {
		addStep(fmt.Sprintf("Total de escritas em document.cookie: %d", len(cookieWrites)))
		for _, write := range cookieWrites {
			value := fmt.Sprintf("%v", write.Data["value"])
			if len(value) > 50 {
				value = value[:50] + "..."
			}
			addStep(fmt.Sprintf("Escrita em %v:%v - %s", write.Data["file"], write.Data["line"], value))
		}
	}

	// Analyze cookie removals
	cookieRemovals := filterEvents(events, "COOKIE_REMOVAL", "")
	if len(cookieRemovals) > 0 {
		addStep(fmt.Sprintf("Total de remoções de cookie: %d", len(cookieRemovals)))
		for _, removal := range cookieRemovals {
			addStep(fmt.Sprintf("Remoção em %v:%v - %v", removal.Data["file"], removal.Data["line"], removal.Data["cookie"]))
		}
	}

	// Analyze fetch calls
	fetchCalls := filterEvents(events, "FETCH", "")
	if len(fetchCalls) > 0 {
		addStep(fmt.Sprintf("Total de requisições fetch: %d", len(fetchCalls)))
		for _, fetch := range fetchCalls {
			if fetch.Data["status"] != nil {
				status := fmt.Sprintf("%v", fetch.Data["status"])
				if status[0] == '4' || status[0] == '5' {
					addStep(fmt.Sprintf("Fetch FALHOU: %v %v - Status: %v", fetch.Data["method"], fetch.Data["url"], fetch.Data["status"]))
				}
			}
		}
	}

	// Conclusion
	summary += "CONCLUSÃO\n"
	summary += "----------------------------------------\n"

	tokenReceived := false
	if requestAfter != nil && requestAfter.Data["token"] != nil && requestAfter.Data["token"] == true {
		tokenReceived = true
	}

	tokenPersistedAfterHydrate := false
	if afterHydrate := findEvent(hydrateEvents, "HYDRATE_CONTEXT", "IMEDIATAMENTE_DEPOIS"); afterHydrate != nil {
		if afterHydrate.Data["authTokenFound"] != nil && afterHydrate.Data["authTokenFound"] == true {
			tokenPersistedAfterHydrate = true
		}
	}

	tokenPersistedAfter100ms := false
	if after100ms := findEvent(hydrateEvents, "HYDRATE_CONTEXT", "DEPOIS_100MS"); after100ms != nil {
		if after100ms.Data["authTokenFound"] != nil && after100ms.Data["authTokenFound"] == true {
			tokenPersistedAfter100ms = true
		}
	}

	if !tokenReceived {
		summary += "O Tenant JWT NÃO foi recebido do backend.\n"
		summary += "CAUSA MAIS PROVÁVEL: Backend não gerou ou não enviou o token corretamente.\n"
		summary += "Confiança: 90%\n"
	} else if !tokenPersistedAfterHydrate {
		summary += "O Tenant JWT foi recebido mas NÃO persistiu no cookie após hydrateContext().\n"
		summary += "CAUSA MAIS PROVÁVEL: Problema na escrita do cookie ou remoção imediata.\n"
		summary += "Confiança: 85%\n"
	} else if !tokenPersistedAfter100ms {
		summary += "O Tenant JWT persistiu inicialmente mas foi removido dentro de 100ms.\n"
		summary += "CAUSA MAIS PROVÁVEL: Outro código removeu o cookie (verificar seções 1.5 e 1.6).\n"
		summary += "Confiança: 80%\n"
	} else {
		summary += "O Tenant JWT foi recebido e persistiu corretamente.\n"
		summary += "O problema pode estar em outra parte do fluxo (autenticação subsequente).\n"
		summary += "Confiança: 70%\n"
	}

	return summary
}

func findEvent(events []ForensicEvent, section, subsection string) *ForensicEvent {
	for _, event := range events {
		if event.Section == section && (subsection == "" || event.Subsection == subsection) {
			return &event
		}
	}
	return nil
}

func filterEvents(events []ForensicEvent, section, subsection string) []ForensicEvent {
	result := make([]ForensicEvent, 0)
	for _, event := range events {
		if event.Section == section && (subsection == "" || event.Subsection == subsection) {
			result = append(result, event)
		}
	}
	return result
}
