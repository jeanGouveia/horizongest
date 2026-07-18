export interface ErrorResponse {
	code: string;
	message: string;
	details?: any;
	timestamp: string;
	requestId: string;
}

export class ErrorService {
	static formatError(error: any): string {
		if (typeof error === 'string') return error;
		
		// Se já é um erro formatado do backend
		if (error?.code && error?.message) {
			return error.message;
		}
		
		// Se é erro da API
		if (error?.error) {
			return error.error;
		}
		
		// Se é erro nativo
		if (error?.message) {
			return error.message;
		}
		
		return 'Erro desconhecido. Tente novamente.';
	}

	static getErrorCode(error: any): string {
		if (error?.code) return error.code;
		if (error?.error) return 'API_ERROR';
		return 'UNKNOWN_ERROR';
	}

	static getRequestId(error: any): string | null {
		return error?.requestId || null;
	}

	static getDetails(error: any): any {
		return error?.details || null;
	}

	static isValidationError(error: any): boolean {
		const code = this.getErrorCode(error);
		return code === 'VALIDATION_ERROR' || code === 'INVALID_INPUT';
	}

	static isAuthError(error: any): boolean {
		const code = this.getErrorCode(error);
		return code === 'AUTH_ERROR' || code === 'UNAUTHORIZED' || code === 'FORBIDDEN';
	}

	static isNotFoundError(error: any): boolean {
		const code = this.getErrorCode(error);
		return code === 'NOT_FOUND';
	}

	static isServerError(error: any): boolean {
		const code = this.getErrorCode(error);
		return code === 'INTERNAL_ERROR' || code === 'SERVER_ERROR';
	}

	static getErrorVariant(error: any): 'error' | 'warning' | 'info' {
		if (this.isValidationError(error)) return 'warning';
		if (this.isAuthError(error)) return 'error';
		if (this.isNotFoundError(error)) return 'warning';
		if (this.isServerError(error)) return 'error';
		return 'error';
	}
}

export const errorService = ErrorService;
