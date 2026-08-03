package handler

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jeanGouveia/horizongest/backend/internal/domain"
	"github.com/jeanGouveia/horizongest/backend/internal/service"
)

func TestMediaHandler_UploadMedia_RequiresAuth(t *testing.T) {
	svc := service.NewMediaService(nil)
	handler := NewMediaHandler(svc)

	// Create a request without tenant context
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("entity_type", "product")
	writer.WriteField("entity_id", "1")
	writer.Close()

	req := httptest.NewRequest("POST", "/api/media/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler.UploadMedia(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rr.Code)
	}
}

func TestMediaHandler_UploadMedia_WithTenantContext(t *testing.T) {
	// This test verifies that tenant context validation passes
	// We'll test the validation logic directly without calling the service

	tenantCtx := &domain.TenantContext{
		UserID:    1,
		CompanyID: 123,
	}

	// Create a mock context with tenant context
	ctx := context.WithValue(context.Background(), domain.ContextKeyTenant, tenantCtx)

	// Verify tenant context can be retrieved
	retrievedCtx, ok := domain.GetTenantContextFromContext(ctx)
	if !ok {
		t.Error("expected to retrieve tenant context from context")
	}

	if retrievedCtx.CompanyID != 123 {
		t.Errorf("expected CompanyID 123, got %d", retrievedCtx.CompanyID)
	}
}

func TestMediaHandler_ServeFile_DirectoryTraversal(t *testing.T) {
	svc := service.NewMediaService(nil)
	handler := NewMediaHandler(svc)

	testCases := []struct {
		name     string
		path     string
		expected int
	}{
		{"double dot", "../etc/passwd", http.StatusBadRequest},
		{"double dot encoded", "%2e%2e%2fetc%2fpasswd", http.StatusBadRequest},
		{"backslash", "..\\windows\\system32", http.StatusBadRequest},
		{"double dot in middle", "products/../etc/passwd", http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/uploads/"+tc.path, nil)
			rr := httptest.NewRecorder()

			handler.ServeFile(rr, req)

			if rr.Code != tc.expected {
				t.Errorf("expected status %d, got %d", tc.expected, rr.Code)
			}
		})
	}
}

func TestMediaHandler_ServeFile_PathValidation(t *testing.T) {
	svc := service.NewMediaService(nil)
	handler := NewMediaHandler(svc)

	// Test that a valid path doesn't trigger path validation errors
	req := httptest.NewRequest("GET", "/uploads/products/test.jpg", nil)
	rr := httptest.NewRecorder()

	handler.ServeFile(rr, req)

	// Should not return 400 (path validation passed)
	// It will return 404 because file doesn't exist, but that's OK
	if rr.Code == http.StatusBadRequest && strings.Contains(rr.Body.String(), "caminho inválido") {
		t.Error("path validation should accept valid paths")
	}
}
