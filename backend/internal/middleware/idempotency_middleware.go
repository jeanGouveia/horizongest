package middleware

import (
	"bytes"
	"errors"
	"log"
	"net/http"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

const (
	IdempotencyKeyHeader = "Idempotency-Key"
)

// IdempotencyMiddleware gerencia idempotência para operações mutáveis
// Deve ser aplicado apenas em endpoints POST, PUT, PATCH, DELETE
type IdempotencyMiddleware struct {
	idempotencyService *service.IdempotencyService
}

func NewIdempotencyMiddleware(idempotencyService *service.IdempotencyService) *IdempotencyMiddleware {
	return &IdempotencyMiddleware{
		idempotencyService: idempotencyService,
	}
}

// Handler wraps http.Handler com lógica de idempotência
func (m *IdempotencyMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extrair idempotency key do header
		idempotencyKey := r.Header.Get(IdempotencyKeyHeader)

		// Se não tem chave, prosseguir sem idempotência (ou retornar erro)
		// Aqui optamos por prosseguir para manter backward compatibility
		if idempotencyKey == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Obter companyID do contexto (tenant)
		tenantCtxValue := r.Context().Value("tenant")
		if tenantCtxValue == nil {
			jsonError(w, "contexto de tenant não encontrado", http.StatusInternalServerError)
			return
		}
		tenantCtx, ok := tenantCtxValue.(*domain.TenantContext)
		if !ok {
			jsonError(w, "contexto de tenant inválido", http.StatusInternalServerError)
			return
		}
		companyID := tenantCtx.CompanyID

		// Ler body da requisição
		body, err := service.ReadBody(r)
		if err != nil {
			jsonError(w, "erro ao ler body da requisição", http.StatusBadRequest)
			return
		}

		// Extrair parâmetros da requisição
		params := service.ExtractRequestParams(r)

		// Verificar idempotência
		result, recordID, err := m.idempotencyService.CheckAndCreate(
			r.Context(),
			companyID,
			idempotencyKey,
			body,
			params,
		)

		if err != nil {
			// Tratar erros específicos de idempotência
			if errors.Is(err, service.ErrIdempotencyKeyReuse) {
				jsonError(w, "idempotency_key já foi utilizada com payload diferente. Use uma nova chave.", http.StatusBadRequest)
				return
			}
			if errors.Is(err, service.ErrIdempotencyProcessing) {
				jsonError(w, "operação com esta chave ainda está em processamento. Aguarde e tente novamente.", http.StatusConflict)
				return
			}
			// Outros erros
			log.Printf("[IdempotencyMiddleware] Erro ao verificar idempotência: %v", err)
			jsonError(w, "erro ao verificar idempotência", http.StatusInternalServerError)
			return
		}

		// Se é replayable, retornar resposta armazenada
		if result.Replayable && result.Response != nil {
			log.Printf("[IdempotencyMiddleware] Replay response: key=%s", idempotencyKey)
			// Restaurar headers
			for k, v := range result.Response.Headers {
				w.Header().Set(k, v)
			}
			w.WriteHeader(result.Response.StatusCode)
			w.Write([]byte(result.Response.Body))
			return
		}

		// Criar response wrapper para capturar a resposta
		rw := &idempotencyResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Executar handler
		next.ServeHTTP(rw, r)

		// Se foi criado um novo registro, registrar resultado
		if recordID > 0 {
			// Capturar body da resposta
			responseBody := rw.body.String()

			// Se status code for 2xx, registrar sucesso
			if rw.statusCode >= 200 && rw.statusCode < 300 {
				if err := m.idempotencyService.RecordSuccess(
					r.Context(),
					recordID,
					rw.statusCode,
					w.Header(),
					[]byte(responseBody),
				); err != nil {
					log.Printf("[IdempotencyMiddleware] Erro ao registrar sucesso: %v", err)
				}
			} else {
				// Registrar falha
				if err := m.idempotencyService.RecordFailure(
					r.Context(),
					recordID,
					responseBody,
				); err != nil {
					log.Printf("[IdempotencyMiddleware] Erro ao registrar falha: %v", err)
				}
			}
		}
	})
}

// idempotencyResponseWriter captura status code e body da resposta
type idempotencyResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rw *idempotencyResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *idempotencyResponseWriter) Write(b []byte) (int, error) {
	if rw.body == nil {
		rw.body = &bytes.Buffer{}
	}
	return rw.body.Write(b)
}
