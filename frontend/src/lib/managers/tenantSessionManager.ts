import { goto } from '$app/navigation';
import { userStore } from '$lib/stores/userStore.svelte';
import { companyStore } from '$lib/stores/companyStore.svelte';
import { rbacStore } from '$lib/stores/rbacStore.svelte';
import { themeStore } from '$lib/stores/themeStore.svelte';
import { brandStore } from '$lib/stores/brandStore';
import { toast } from '$lib/stores/toast';
import { showSuccess, showError } from '$lib/stores/toast';
import { CookieKeys, TenantLocalStorageKeys, TenantSessionStorageKeys } from '$lib/constants/storage-keys';
import { forensicLogger } from '$lib/forensic/forensic-logger';

/**
 * Error types for better error handling
 */
class SessionError extends Error {
	constructor(message: string, public type: 'session' | 'infrastructure' | 'backend' | 'ui') {
		super(message);
		this.name = 'SessionError';
	}
}

class InfrastructureError extends SessionError {
	constructor(message: string) {
		super(message, 'infrastructure');
		this.name = 'InfrastructureError';
	}
}

class SessionValidationError extends SessionError {
	constructor(message: string) {
		super(message, 'session');
		this.name = 'SessionValidationError';
	}
}

class BackendError extends SessionError {
	constructor(message: string, public status?: number) {
		super(message, 'backend');
		this.name = 'BackendError';
	}
}

class UIError extends SessionError {
	constructor(message: string) {
		super(message, 'ui');
		this.name = 'UIError';
	}
}

interface Company {
	ID: number;
	Name: string;
	Slug: string;
}

interface EnterCompanyResult {
	success: boolean;
	error?: string;
}

interface LeaveCompanyResult {
	success: boolean;
	error?: string;
}

interface SwitchCompanyResult {
	success: boolean;
	error?: string;
}

/**
 * TenantSessionManager
 * 
 * Gerenciador centralizado do ciclo de vida da sessão Tenant.
 * 
 * Responsabilidades:
 * - Gerenciar entrada em empresa (enterCompany)
 * - Gerenciar saída de empresa (leaveCompany)
 * - Gerenciar troca de empresa (switchCompany)
 * - Destruir contexto completo (destroy)
 * 
 * Este é o ÚNICO lugar onde deve ocorrer troca de empresa.
 * Nenhum outro componente deve trocar empresa diretamente.
 */
class TenantSessionManager {
	private currentCompanyId: number | null = null;
	private isEntering = false;
	private isLeaving = false;

	/**
	 * Entra em uma empresa
	 * 
	 * Fluxo:
	 * 1. Valida empresa
	 * 2. Valida se já está entrando em uma empresa (previne race conditions)
	 * 3. Valida se já está na empresa selecionada
	 * 4. Encerra impersonation anterior (caso exista)
	 * 5. Destrói contexto anterior
	 * 6. Solicita Tenant JWT
	 * 7. Limpar caches
	 * 8. Hidrata contexto
	 * 9. Carrega branding
	 * 10. Carrega permissões
	 * 11. Carrega empresa
	 * 12. Navegar para Dashboard
	 */
	async enterCompany(companyId: number): Promise<EnterCompanyResult> {

		// Validações de estado
		if (this.isEntering) {
			throw new SessionValidationError('Já está entrando em uma empresa');
		}

		if (this.currentCompanyId === companyId) {
			throw new SessionValidationError('Já está nesta empresa');
		}

		this.isEntering = true;

		// FORENSIC: Log initial state
		forensicLogger.log('FRONTEND', 'INITIAL_STATE', undefined, forensicLogger.getCurrentState());

		try {
			// 1. Destruir contexto anterior
			await this.destroy();

			// 2. Solicitar Tenant JWT
			// O backend encerra automaticamente qualquer impersonation ativa antes de criar uma nova
			// e define o cookie auth_token via Set-Cookie
			const success = await this.requestTenantJWT(companyId);
			if (!success) {
				throw new BackendError('Falha ao iniciar impersonation', 401);
			}


			// 3. Limpar caches
			this.clearCaches();

			// 4. Hidratar contexto
			await this.hydrateContext(companyId);


			// 5. Carregar branding
			await this.loadBranding();

			// 6. Carregar permissões
			await this.loadPermissions();

			// 7. Carregar empresa
			await this.loadCompany();


			// 8. Atualizar estado
			this.currentCompanyId = companyId;

			// 9. Navegar para Dashboard
			forensicLogger.recordNavigation('goto', '/dashboard');
			await goto('/dashboard');


			return { success: true };
		} catch (error) {
			// Tratamento diferenciado por tipo de erro
			if (error instanceof SessionError) {
				switch (error.type) {
					case 'infrastructure':
						console.error('[Infrastructure Error]', error.message);
						return { success: false, error: `Erro de conexão: ${error.message}. Verifique sua internet.` };
					case 'session':
						console.error('[Session Validation Error]', error.message);
						return { success: false, error: error.message };
					case 'backend':
						console.error('[Backend Error]', error.message);
						if (error instanceof BackendError && error.status === 401) {
							return { success: false, error: 'Sessão expirada. Faça login novamente.' };
						}
						return { success: false, error: `Erro do servidor: ${error.message}` };
					case 'ui':
						console.error('[UI Error]', error.message);
						return { success: false, error: error.message };
				}
			}

			// Tratamento de erros genéricos
			if (error instanceof TypeError && error.message.includes('fetch')) {
				console.error('[Network Error]', error.message);
				return { success: false, error: 'Erro de conexão. Verifique se o backend está rodando em http://localhost:8080' };
			}

			console.error('[Unknown Error]', error);
			return { 
				success: false, 
				error: error instanceof Error ? error.message : 'Erro desconhecido ao entrar na empresa' 
			};
		} finally {
			this.isEntering = false;
		}
	}

	/**
	 * Sai da empresa atual
	 * 
 * Fluxo:
	 * 1. Destrói contexto Tenant
	 * 2. Limpa stores Tenant
	 * 3. Limpa caches Tenant
	 * 4. Mantém Platform Session
	 * 5. Retorna para Plataforma
	 */
	async leaveCompany(): Promise<LeaveCompanyResult> {
		if (this.isLeaving) {
			return { success: false, error: 'Já está saindo da empresa' };
		}

		this.isLeaving = true;

		try {
			// 1. Destruir contexto Tenant
			await this.destroy();

			// 2. Limpar caches
			this.clearCaches();

			// 3. Atualizar estado
			this.currentCompanyId = null;

			// 4. Navegar para Plataforma
			forensicLogger.recordNavigation('goto', '/platform/admin');
			await goto('/platform/admin');

			return { success: true };
		} catch (error) {
			console.error('Erro ao sair da empresa:', error);
			return { 
				success: false, 
				error: error instanceof Error ? error.message : 'Erro desconhecido ao sair da empresa' 
			};
		} finally {
			this.isLeaving = false;
		}
	}

	/**
	 * Troca de empresa
	 * 
 * Executa leaveCompany() seguido de enterCompany()
	 */
	async switchCompany(companyId: number): Promise<SwitchCompanyResult> {
		try {
			// Sair da empresa atual
			const leaveResult = await this.leaveCompany();
			if (!leaveResult.success) {
				return { success: false, error: leaveResult.error };
			}

			// Entrar na nova empresa
			const enterResult = await this.enterCompany(companyId);
			if (!enterResult.success) {
				return { success: false, error: enterResult.error };
			}

			return { success: true };
		} catch (error) {
			console.error('Erro ao trocar de empresa:', error);
			return { 
				success: false, 
				error: error instanceof Error ? error.message : 'Erro desconhecido ao trocar de empresa' 
			};
		}
	}

	/**
	 * Destroi completamente o contexto do tenant
	 * 
	 * Responsabilidades:
	 * - Limpar todos os stores Svelte (user, company, rbac, theme, brand, toast)
	 * - Limpar cookies de autenticação tenant (auth_token)
	 * - Limpar localStorage tenant (impersonation)
	 * 
	 * NÃO limpa:
	 * - Platform session (platform_auth_token)
	 * - Dados globais do navegador
	 * - Caches internos do SvelteKit
	 * 
	 * Quando é chamado:
	 * - Ao entrar em uma nova empresa (antes de carregar novo contexto)
	 * - Ao sair de uma empresa
	 * - Ao trocar de empresa
	 */
	async destroy(): Promise<void> {
		console.log("⚠️ destroy() CHAMADO - Destruindo contexto tenant");
		console.log("Stack trace:", new Error().stack);

		// Limpar stores
		this.clearStores();

		// Limpar cookies tenant
		this.clearTenantCookies();

		// Limpar localStorage tenant
		this.clearTenantLocalStorage();
	}

	/**
	 * Solicita Tenant JWT ao backend
	 * 
	 * O backend gera o JWT e define o cookie auth_token via Set-Cookie
	 * Este método apenas verifica se a operação foi bem-sucedida
	 */
	private async requestTenantJWT(companyId: number): Promise<boolean> {

		try {
			const platformToken = document.cookie
				.split('; ')
				.find(row => row.startsWith(`${CookieKeys.PLATFORM_TOKEN}=`))
				?.split('=')[1];

			const url = 'http://localhost:8080/api/platform/impersonation/start';
			const method = 'POST';
			const credentials = 'include';
			const headers = {
				'Content-Type': 'application/json',
				'Authorization': `Bearer ${platformToken}`
			};
			const body = JSON.stringify({ companyId });

			// FORENSIC: Log before request
			forensicLogger.log('FRONTEND', 'REQUEST_TENANT_JWT', 'BEFORE', {
				url,
				method,
				headers,
				authorization: headers.Authorization,
				credentials,
				body
			});

			forensicLogger.recordFetch(url, method, credentials, headers, document.cookie);

			if (!platformToken) {
				throw new Error('Token da plataforma não encontrado');
			}

			const response = await fetch(url, {
				method,
				credentials,
				headers,
				body
			});

			// FORENSIC: Log response
			const responseHeaders: Record<string, string> = {};
			response.headers.forEach((value, key) => {
				responseHeaders[key] = value;
			});
			const setCookieHeader = response.headers.get('Set-Cookie');

			forensicLogger.recordFetchResponse(response.status, responseHeaders, null);

			if (!response.ok) {
				const errorData = await response.json();
				throw new Error(errorData.error || 'Erro ao iniciar impersonation');
			}

			const data = await response.json();
			forensicLogger.recordFetchResponse(response.status, responseHeaders, data);

			// FORENSIC: Log after request
			forensicLogger.log('FRONTEND', 'REQUEST_TENANT_JWT', 'AFTER', {
				status: response.status,
				headers: responseHeaders,
				setCookie: setCookieHeader,
				body: data,
				success: data.success === true
			});

			return data.success === true;
		} catch (error) {
			return false;
		}
	}

	/**
	 * Limpa todas as stores
	 */
	private clearStores(): void {
		userStore.logout();
		companyStore.clear();
		rbacStore.reset();
		themeStore.clear();
		brandStore.clear();
		toast.clear();
	}

	/**
	 * Limpa todos os caches do tenant de forma granular
	 * 
	 * Responsabilidades:
	 * - Limpar apenas dados que pertencem ao contexto Tenant
	 * - Nunca apagar dados da Platform
	 * - Nunca apagar informações globais
	 * - Nunca apagar caches do navegador
	 * 
	 * Caches limpos:
	 * - localStorage: impersonation, tenant_context, tenant_permissions, tenant_branding, tenant_dashboard, tenant_cache
	 * - sessionStorage: tenant_navigation, tenant_forms, tenant_filters
	 * 
	 * Caches NÃO limpos:
	 * - Dados da Platform (platform_user, platform_preferences)
	 * - Dados globais do navegador
	 * - Caches internos do SvelteKit
	 */
	private clearCaches(): void {
		// Limpar localStorage tenant - dados de impersonação e contexto
		this.clearTenantLocalStorage();

		// Limpar sessionStorage tenant - dados temporários específicos
		this.clearTenantSessionStorage();
	}

	/**
	 * Limpa localStorage do tenant de forma granular
	 * 
	 * Remove apenas chaves que pertencem ao contexto Tenant:
	 * - impersonation: dados de impersonação ativa
	 * - tenant_context: contexto do tenant
	 * - tenant_permissions: permissões do usuário
	 * - tenant_branding: branding do tenant
	 * - tenant_dashboard: estado do dashboard
	 * - tenant_cache: cache de dados
	 */
	private clearTenantLocalStorage(): void {
		localStorage.removeItem(TenantLocalStorageKeys.IMPERSONATION);
		localStorage.removeItem(TenantLocalStorageKeys.TENANT_CONTEXT);
		localStorage.removeItem(TenantLocalStorageKeys.TENANT_PERMISSIONS);
		localStorage.removeItem(TenantLocalStorageKeys.TENANT_BRANDING);
		localStorage.removeItem(TenantLocalStorageKeys.TENANT_DASHBOARD);
		localStorage.removeItem(TenantLocalStorageKeys.TENANT_CACHE);
	}

	/**
	 * Limpa sessionStorage do tenant de forma granular
	 * 
	 * Remove apenas chaves que pertencem ao contexto Tenant:
	 * - tenant_navigation: estado de navegação
	 * - tenant_forms: dados de formulários
	 * - tenant_filters: filtros aplicados
	 * 
	 * NÃO limpa dados da Platform ou globais.
	 */
	private clearTenantSessionStorage(): void {
		sessionStorage.removeItem(TenantSessionStorageKeys.TENANT_NAVIGATION);
		sessionStorage.removeItem(TenantSessionStorageKeys.TENANT_FORMS);
		sessionStorage.removeItem(TenantSessionStorageKeys.TENANT_FILTERS);
	}

	/**
	 * Limpa cookies tenant
	 * 
	 * NOTA: O backend deve limpar os cookies via Set-Cookie com expiração no passado
	 * Este método é mantido apenas para limpeza local em caso de necessidade
	 */
	private clearTenantCookies(): void {
		console.log("⚠️ clearTenantCookies CHAMADO - REMOVENDO auth_token");
		// Limpar auth_token (usado pelo backend para tenant)
		const authTokenValue = `auth_token=; path=/; max-age=0; SameSite=Lax`;
		document.cookie = authTokenValue;
		forensicLogger.recordCookieRemoval('tenantSessionManager.ts', 411, 'auth_token');
	}

	/**
	 * Hidrata o contexto do tenant
	 * 
	 * Responsabilidades:
	 * - Atualizar stores
	 * - Carregar branding
	 * - Carregar permissões
	 * - Carregar empresa
	 * 
	 * NÃO grava cookies - o backend já definiu o auth_token via Set-Cookie
	 */
	private async hydrateContext(companyId: number): Promise<void> {
		// Cookie já foi definido pelo backend via Set-Cookie
		// Apenas preparamos o contexto frontend
	}

	/**
	 * Carrega branding
	 */
	private async loadBranding(): Promise<void> {
		try {
			await brandStore.load();
		} catch (error) {
			console.error('Erro ao carregar branding:', error);
		}
	}

	/**
	 * Carrega permissões
	 */
	private async loadPermissions(): Promise<void> {
		try {
			await rbacStore.load();
		} catch (error) {
			console.error('Erro ao carregar permissões:', error);
		}
	}

	/**
	 * Carrega empresa
	 * 
	 * Usa o endpoint seguro GET /api/me/company que obtém a empresa
	 * exclusivamente do contexto tenant (JWT), sem receber CompanyID por URL
	 */
	private async loadCompany(): Promise<void> {
		try {
			const url = 'http://localhost:8080/api/me/company';
			const options = {
				credentials: 'include' as RequestCredentials,
				headers: {
					'Content-Type': 'application/json'
				}
			};

			const response = await fetch(url, options);

			if (!response.ok) {
				throw new Error('Erro ao carregar empresa');
			}

			const company: Company = await response.json();

			companyStore.setCompany({ name: company.Name });

			// Atualizar localStorage impersonation (metadados, não JWT)
			localStorage.setItem('impersonation', JSON.stringify({
				isImpersonating: true,
				companyName: company.Name,
				companyId: company.ID
			}));

		} catch (error) {
			throw error;
		}
	}

	/**
	 * Encerra a impersonation anterior no backend
	 */
	private async endPreviousImpersonation(): Promise<void> {
		try {
			const platformToken = document.cookie
				.split('; ')
				.find(row => row.startsWith(`${CookieKeys.PLATFORM_TOKEN}=`))
				?.split('=')[1];

			if (!platformToken) return;

			const response = await fetch('http://localhost:8080/api/platform/impersonation/end', {
				method: 'POST',
				credentials: 'include',
				headers: {
					'Authorization': `Bearer ${platformToken}`,
					'Content-Type': 'application/json'
				}
			});


			if (!response.ok) {
				console.error('Erro ao encerrar impersonation anterior:', await response.text());
			}
		} catch (error) {
			console.error('Erro ao encerrar impersonation anterior:', error);
		}
	}

	/**
	 * Retorna a empresa atual
	 */
	getCurrentCompanyId(): number | null {
		return this.currentCompanyId;
	}

	/**
	 * Verifica se está em uma empresa
	 */
	isInCompany(): boolean {
		return this.currentCompanyId !== null;
	}
}

// Export singleton instance
export const tenantSessionManager = new TenantSessionManager();
