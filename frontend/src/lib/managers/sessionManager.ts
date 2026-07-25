import { goto } from '$app/navigation';
import { userStore } from '$lib/stores/userStore.svelte';
import { companyStore } from '$lib/stores/companyStore.svelte';
import { rbacStore } from '$lib/stores/rbacStore.svelte';
import { themeStore } from '$lib/stores/themeStore.svelte';
import { brandStore } from '$lib/stores/brandStore';
import { toast } from '$lib/stores/toast';
import { tenantSessionManager } from './tenantSessionManager';
import { CookieKeys, TenantLocalStorageKeys } from '$lib/constants/storage-keys';

interface SessionValidationResult {
	valid: boolean;
	sessionType: 'platform' | 'tenant' | 'none';
	error?: string;
}

interface LogoutResult {
	success: boolean;
	error?: string;
}

/**
 * SessionManager
 * 
 * Gerenciador centralizado de todas as sessões do HorizonGest.
 * 
 * Responsabilidades:
 * - Validar sessão na inicialização
 * - Gerenciar Platform Session
 * - Gerenciar Tenant Session (via TenantSessionManager)
 * - Executar logout completo
 * - Destruir todas as sessões
 * 
 * Política Oficial do HorizonGest:
 * 
 * Platform Session:
 * - Exige login
 * - Possui expiração
 * - Nunca pode ser restaurada sem validação do backend
 * - Perde validade após reinício do backend
 * - Perde validade após logout
 * - Perde validade após expiração do JWT
 * 
 * Tenant Session:
 * - Só existe enquanto existir uma Platform Session válida
 * - Nunca pode sobreviver sozinha
 * - Nunca pode sobreviver ao logout
 * - Nunca pode sobreviver ao reinício do backend
 * - Nunca pode sobreviver à troca de empresa
 */
class SessionManager {
	private isValidating = false;

	/**
	 * Valida a sessão na inicialização da aplicação
	 * 
	 * Fluxo:
	 * 1. Verifica se existe token
	 * 2. Valida sessão no backend
	 * 3. Se válida, hidrata stores
	 * 4. Se inválida, destrói sessão e redireciona para login
	 */
	async validateSession(): Promise<SessionValidationResult> {
		if (this.isValidating) {
			return { valid: false, sessionType: 'none', error: 'Já está validando sessão' };
		}

		this.isValidating = true;

		try {
			// Verificar Platform Session
			const platformToken = this.getPlatformToken();
			const tenantToken = this.getTenantToken();

			if (!platformToken && !tenantToken) {
				// Nenhuma sessão
				return { valid: false, sessionType: 'none' };
			}

			if (tenantToken && !platformToken) {
				// Tenant Session sem Platform Session - inválido
				await this.destroyAllSessions();
				return { valid: false, sessionType: 'none', error: 'Tenant Session sem Platform Session' };
			}

			// Validar Platform Session
			if (platformToken) {
				const platformValid = await this.validatePlatformSession();
				if (!platformValid) {
					await this.destroyAllSessions();
					return { valid: false, sessionType: 'none', error: 'Platform Session inválida' };
				}
			}

			// Validar Tenant Session
			if (tenantToken) {
				const tenantValid = await this.validateTenantSession();
				if (!tenantValid) {
					// Tenant Session inválida, mas Platform Session pode estar válida
					// Destruir apenas Tenant Session
					await this.destroyTenantSession();
					return { valid: true, sessionType: 'platform' };
				}

				return { valid: true, sessionType: 'tenant' };
			}

			return { valid: true, sessionType: 'platform' };
		} catch (error) {
			console.error('Erro ao validar sessão:', error);
			await this.destroyAllSessions();
			return { 
				valid: false, 
				sessionType: 'none', 
				error: error instanceof Error ? error.message : 'Erro ao validar sessão' 
			};
		} finally {
			this.isValidating = false;
		}
	}

	/**
	 * Executa logout completo
	 * 
	 * Fluxo:
	 * 1. Encerrar impersonation
	 * 2. Encerrar Platform Session
	 * 3. Limpar cookies
	 * 4. Limpar LocalStorage
	 * 5. Limpar SessionStorage
	 * 6. Limpar todas as Stores
	 * 7. Limpar todos os caches
	 * 8. Redirecionar para Login
	 */
	async logout(): Promise<LogoutResult> {
		try {
			// 1. Encerrar impersonation (se existir)
			if (this.getTenantToken()) {
				await this.endImpersonation();
			}

			// 2. Destruir todas as sessões
			await this.destroyAllSessions();

			// 3. Redirecionar para login
			goto('/login');

			return { success: true };
		} catch (error) {
			console.error('Erro ao fazer logout:', error);
			return { 
				success: false, 
				error: error instanceof Error ? error.message : 'Erro ao fazer logout' 
			};
		}
	}

	/**
	 * Destroi todas as sessões
	 * 
	 * Responsabilidades:
	 * - Destruir Tenant Session (cookie auth_token, localStorage impersonation)
	 * - Destruir Platform Session (cookie platform_auth_token)
	 * - Limpar todas as stores Svelte
	 * - Limpar todos os caches (localStorage, sessionStorage)
	 * 
	 * Quando é chamado:
	 * - Logout explícito do usuário
	 * - Validação de sessão falha (token inválido)
	 * - Erro crítico de autenticação
	 * 
	 * Efeito:
	 * - Usuário é redirecionado para /login
	 * - Todas as sessões são completamente limpas
	 */
	async destroyAllSessions(): Promise<void> {
		// Destruir Tenant Session
		await this.destroyTenantSession();

		// Destruir Platform Session
		this.destroyPlatformSession();

		// Limpar todas as stores
		this.clearAllStores();

		// Limpar todos os caches
		this.clearAllCaches();
	}

	/**
	 * Destroi apenas a Tenant Session
	 */
	async destroyTenantSession(): Promise<void> {
		// Limpar cookie tenant
		this.clearTenantCookie();

		// Limpar localStorage tenant
		this.clearTenantLocalStorage();

		// Limpar stores tenant
		this.clearTenantStores();

		// Reset TenantSessionManager
		(tenantSessionManager as any).currentCompanyId = null;
	}

	/**
	 * Destroi a Platform Session
	 */
	destroyPlatformSession(): void {
		// Limpar cookie platform
		this.clearPlatformCookie();

		// Limpar localStorage platform
		this.clearPlatformLocalStorage();
	}

	/**
	 * Valida a Platform Session no backend
	 */
	private async validatePlatformSession(): Promise<boolean> {
		try {
			const token = this.getPlatformToken();
			if (!token) return false;

			const response = await fetch('http://localhost:8080/api/platform/me', {
				headers: {
					'Authorization': `Bearer ${token}`
				}
			});

			return response.ok;
		} catch (error) {
			console.error('Erro ao validar Platform Session:', error);
			return false;
		}
	}

	/**
	 * Valida a Tenant Session no backend
	 */
	private async validateTenantSession(): Promise<boolean> {
		try {
			const token = this.getTenantToken();
			if (!token) return false;

			const response = await fetch('http://localhost:8080/api/me', {
				headers: {
					'Cookie': `auth_token=${token}`
				}
			});

			return response.ok;
		} catch (error) {
			console.error('Erro ao validar Tenant Session:', error);
			return false;
		}
	}

	/**
	 * Encerra a impersonation no backend
	 */
	private async endImpersonation(): Promise<void> {
		try {
			const platformToken = this.getPlatformToken();
			if (!platformToken) return;

			const response = await fetch('http://localhost:8080/api/platform/impersonation/end', {
				method: 'POST',
				headers: {
					'Authorization': `Bearer ${platformToken}`,
					'Content-Type': 'application/json'
				}
			});

			if (!response.ok) {
				console.error('Erro ao encerrar impersonation:', await response.text());
			}
		} catch (error) {
			console.error('Erro ao encerrar impersonation:', error);
		}
	}

	/**
	 * Limpa todas as stores
	 */
	private clearAllStores(): void {
		userStore.logout();
		companyStore.clear();
		rbacStore.reset();
		themeStore.clear();
		brandStore.clear();
		toast.clear();
	}

	/**
	 * Limpa apenas as stores tenant
	 */
	private clearTenantStores(): void {
		userStore.logout();
		companyStore.clear();
		rbacStore.reset();
		themeStore.clear();
		// brandStore e toast são compartilhados, não limpar
	}

	/**
	 * Limpa todos os caches de forma granular
	 * 
	 * Responsabilidades:
	 * - Limpar apenas dados que pertencem ao contexto específico
	 * - Nunca apagar dados globais do navegador
	 * - Nunca apagar caches internos do SvelteKit
	 * 
	 * Caches limpos:
	 * - localStorage: impersonation (tenant)
	 * - sessionStorage: todos (temporal)
	 */
	private clearAllCaches(): void {
		this.clearTenantLocalStorage();
		this.clearPlatformLocalStorage();
		sessionStorage.clear();
	}

	/**
	 * Obtém o Platform Token
	 */
	private getPlatformToken(): string | null {
		return document.cookie.split('; ').find(row => row.startsWith(`${CookieKeys.PLATFORM_TOKEN}=`))?.split('=')[1] || null;
	}

	/**
	 * Obtém o Tenant Token
	 */
	private getTenantToken(): string | null {
		return document.cookie.split('; ').find(row => row.startsWith(`${CookieKeys.TENANT_TOKEN}=`))?.split('=')[1] || null;
	}

	/**
	 * Limpa o cookie Platform
	 */
	private clearPlatformCookie(): void {
		document.cookie = `${CookieKeys.PLATFORM_TOKEN}=; path=/; max-age=0; SameSite=Lax`;
	}

	/**
	 * Limpa o cookie Tenant
	 */
	private clearTenantCookie(): void {
		document.cookie = `${CookieKeys.TENANT_TOKEN}=; path=/; max-age=0; SameSite=Lax`;
	}

	/**
	 * Limpa o localStorage Platform
	 */
	private clearPlatformLocalStorage(): void {
		// Adicionar chaves platform-specific se necessário
	}

	/**
	 * Limpa o localStorage Tenant
	 */
	private clearTenantLocalStorage(): void {
		localStorage.removeItem(TenantLocalStorageKeys.IMPERSONATION);
	}

	/**
	 * Verifica se existe uma Platform Session válida
	 */
	hasPlatformSession(): boolean {
		return this.getPlatformToken() !== null;
	}

	/**
	 * Verifica se existe uma Tenant Session válida
	 */
	hasTenantSession(): boolean {
		return this.getTenantToken() !== null;
	}

	/**
	 * Verifica se está em uma sessão ativa
	 */
	hasActiveSession(): boolean {
		return this.hasPlatformSession() || this.hasTenantSession();
	}
}

// Export singleton instance
export const sessionManager = new SessionManager();
