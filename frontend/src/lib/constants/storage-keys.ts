/**
 * Storage Keys - HorizonGest
 * 
 * Centralização de todas as chaves de storage utilizadas no sistema.
 * 
 * Regras:
 * - NUNCA usar strings literais para storage keys
 * - SEMPRE importar deste arquivo
 * - Nomes em UPPER_CASE para cookies
 * - Nomes em camelCase para localStorage/sessionStorage
 * 
 * Propósito:
 * - Evitar typos
 * - Facilitar refatoração
 * - Documentar o propósito de cada chave
 * - Centralizar conhecimento de storage
 */

/**
 * Cookies
 * Cookies são utilizados para autenticação e sessões.
 * São enviados automaticamente pelo navegador em requisições HTTP.
 */
export const CookieKeys = {
	/** Platform JWT - Token de autenticação da plataforma */
	PLATFORM_TOKEN: 'platform_auth_token',

	/** Tenant JWT - Token de autenticação do tenant (empresa) */
	TENANT_TOKEN: 'auth_token'
} as const;

/**
 * LocalStorage - Platform
 * Dados específicos da plataforma que persistem entre sessões.
 */
export const PlatformLocalStorageKeys = {
	/** Dados do usuário da plataforma (opcional, se necessário) */
	PLATFORM_USER: 'platform_user',

	/** Preferências da plataforma (tema, idioma, etc.) */
	PLATFORM_PREFERENCES: 'platform_preferences'
} as const;

/**
 * LocalStorage - Tenant
 * Dados específicos do tenant que persistem entre sessões.
 */
export const TenantLocalStorageKeys = {
	/** Dados de impersonação ativa */
	IMPERSONATION: 'impersonation',

	/** Contexto do tenant (empresa atual) */
	TENANT_CONTEXT: 'tenant_context',

	/** Permissões do usuário no tenant */
	TENANT_PERMISSIONS: 'tenant_permissions',

	/** Branding do tenant (logo, cores, etc.) */
	TENANT_BRANDING: 'tenant_branding',

	/** Estado do dashboard do tenant */
	TENANT_DASHBOARD: 'tenant_dashboard',

	/** Cache de dados do tenant (para performance) */
	TENANT_CACHE: 'tenant_cache'
} as const;

/**
 * SessionStorage - Tenant
 * Dados temporários do tenant que são limpos ao fechar o navegador.
 */
export const TenantSessionStorageKeys = {
	/** Estado temporário de navegação do tenant */
	TENANT_NAVIGATION: 'tenant_navigation',

	/** Dados temporários de formulários do tenant */
	TENANT_FORMS: 'tenant_forms',

	/** Estado temporário de filtros do tenant */
	TENANT_FILTERS: 'tenant_filters'
} as const;

/**
 * SessionStorage - Platform
 * Dados temporários da plataforma que são limpos ao fechar o navegador.
 */
export const PlatformSessionStorageKeys = {
	/** Estado temporário de navegação da plataforma */
	PLATFORM_NAVIGATION: 'platform_navigation',

	/** Dados temporários de formulários da plataforma */
	PLATFORM_FORMS: 'platform_forms'
} as const;

/**
 * Tipos
 */
export type CookieKey = typeof CookieKeys[keyof typeof CookieKeys];
export type PlatformLocalStorageKey = typeof PlatformLocalStorageKeys[keyof typeof PlatformLocalStorageKeys];
export type TenantLocalStorageKey = typeof TenantLocalStorageKeys[keyof typeof TenantLocalStorageKeys];
export type TenantSessionStorageKey = typeof TenantSessionStorageKeys[keyof typeof TenantSessionStorageKeys];
export type PlatformSessionStorageKey = typeof PlatformSessionStorageKeys[keyof typeof PlatformSessionStorageKeys];
