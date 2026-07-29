package util

import (
	"strings"
)

const (
	maskLength = 8
	maskChar   = "*"
)

// MaskToken mascara tokens sensíveis (JWT, API keys, etc.)
// Mantém os primeiros e últimos caracteres visíveis para debugging
func MaskToken(token string) string {
	if token == "" {
		return ""
	}
	
	length := len(token)
	if length <= maskLength*2 {
		return strings.Repeat(maskChar, length)
	}
	
	return token[:maskLength] + strings.Repeat(maskChar, length-maskLength*2) + token[length-maskLength:]
}

// MaskEmail mascara emails para logs
// Exemplo: jo***@example.com
func MaskEmail(email string) string {
	if email == "" {
		return ""
	}
	
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return MaskToken(email)
	}
	
	username := parts[0]
	domain := parts[1]
	
	if len(username) <= 3 {
		username = strings.Repeat(maskChar, len(username))
	} else {
		username = username[:3] + strings.Repeat(maskChar, len(username)-3)
	}
	
	return username + "@" + domain
}

// MaskSecret mascara secrets sensíveis
// Retorna apenas asteriscos
func MaskSecret(secret string) string {
	if secret == "" {
		return ""
	}
	return strings.Repeat(maskChar, len(secret))
}

// MaskAuthorizationHeader mascara headers de autorização
// Remove o valor do Bearer token
func MaskAuthorizationHeader(header string) string {
	if header == "" {
		return ""
	}
	
	if strings.HasPrefix(header, "Bearer ") {
		return "Bearer " + MaskToken(header[7:])
	}
	
	return MaskToken(header)
}

// MaskCookieValue mascara valores de cookies
// Mantém o nome do cookie mas mascara o valor
func MaskCookieValue(cookieName, cookieValue string) string {
	if cookieValue == "" {
		return ""
	}
	return cookieName + "=" + MaskToken(cookieValue)
}
