/**
 * Forensic Logger - Captures detailed execution trace for impersonation bug investigation
 * This logger writes all events to a buffer that will be saved to forensic-report.txt
 */

interface ForensicEvent {
	timestamp: number;
	category: 'FRONTEND' | 'BACKEND';
	section: string;
	subsection?: string;
	data: Record<string, any>;
}

class ForensicLogger {
	private events: ForensicEvent[] = [];
	private startTime: number = Date.now();
	private cookieWrites: { file: string; line: number; value: string; timestamp: number }[] = [];
	private cookieRemovals: { file: string; line: number; cookie: string; timestamp: number }[] = [];
	private fetchCalls: { timestamp: number; url: string; method: string; credentials: string; headers: Record<string, string>; cookiesBefore: string; status?: number; headersReceived?: Record<string, string>; bodyReceived?: any }[] = [];
	private navigations: { timestamp: number; type: 'goto' | 'beforeNavigate' | 'afterNavigate'; url?: string; from?: string; to?: string }[] = [];

	private getTimestamp(): number {
		return Date.now() - this.startTime;
	}

	async sendLogToBackend(event: ForensicEvent) {
		try {
			await fetch('http://localhost:8080/api/forensic/log', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(event)
			});
		} catch (error) {
			// Silently fail - don't interfere with the actual flow
		}
	}

	log(category: 'FRONTEND' | 'BACKEND', section: string, subsection: string | undefined, data: Record<string, any>) {
		const event: ForensicEvent = {
			timestamp: this.getTimestamp(),
			category,
			section,
			subsection,
			data
		};
		this.events.push(event);
		// Also send to backend immediately
		this.sendLogToBackend(event);
	}

	recordCookieWrite(file: string, line: number, value: string) {
		this.cookieWrites.push({
			file,
			line,
			value,
			timestamp: this.getTimestamp()
		});
		this.log('FRONTEND', 'COOKIE_WRITE', undefined, { file, line, value });
	}

	recordCookieRemoval(file: string, line: number, cookie: string) {
		this.cookieRemovals.push({
			file,
			line,
			cookie,
			timestamp: this.getTimestamp()
		});
		this.log('FRONTEND', 'COOKIE_REMOVAL', undefined, { file, line, cookie });
	}

	recordFetch(url: string, method: string, credentials: string, headers: Record<string, string>, cookiesBefore: string) {
		this.fetchCalls.push({
			timestamp: this.getTimestamp(),
			url,
			method,
			credentials,
			headers,
			cookiesBefore
		});
		this.log('FRONTEND', 'FETCH', 'BEFORE', { url, method, credentials, headers, cookiesBefore });
	}

	recordFetchResponse(status: number, headersReceived: Record<string, string>, bodyReceived: any) {
		const lastFetch = this.fetchCalls[this.fetchCalls.length - 1];
		if (lastFetch) {
			lastFetch.status = status;
			lastFetch.headersReceived = headersReceived;
			lastFetch.bodyReceived = bodyReceived;
		}
		this.log('FRONTEND', 'FETCH', 'AFTER', { status, headersReceived, bodyReceived });
	}

	recordNavigation(type: 'goto' | 'beforeNavigate' | 'afterNavigate', url?: string, from?: string, to?: string) {
		this.navigations.push({
			timestamp: this.getTimestamp(),
			type,
			url,
			from,
			to
		});
		this.log('FRONTEND', 'NAVIGATION', undefined, { type, url, from, to });
	}

	getCurrentState() {
		return {
			documentCookie: typeof document !== 'undefined' ? document.cookie : 'N/A (SSR)',
			localStorage: typeof localStorage !== 'undefined' ? JSON.stringify(localStorage) : 'N/A (SSR)',
			sessionStorage: typeof sessionStorage !== 'undefined' ? JSON.stringify(sessionStorage) : 'N/A (SSR)'
		};
	}

	generateReport(): string {
		let report = '';
		const separator = '='.repeat(80);

		// FRONTEND SECTION
		report += `${separator}\n`;
		report += `1. FRONTEND\n`;
		report += `${separator}\n\n`;

		// 1.1 Initial State
		report += `1.1 Estado inicial\n`;
		report += '-'.repeat(40) + '\n';
		const initialState = this.events.find(e => e.section === 'INITIAL_STATE');
		if (initialState) {
			report += `Timestamp: ${initialState.timestamp}ms\n`;
			report += `document.cookie: ${initialState.data.documentCookie}\n`;
			report += `localStorage: ${initialState.data.localStorage}\n`;
			report += `sessionStorage: ${initialState.data.sessionStorage}\n`;
		}
		report += '\n';

		// 1.2 Click on "Entrar" button
		report += `1.2 Clique no botão "Entrar"\n`;
		report += '-'.repeat(40) + '\n';
		const clickEvent = this.events.find(e => e.section === 'CLICK_ENTRAR');
		if (clickEvent) {
			report += `Timestamp: ${clickEvent.timestamp}ms\n`;
			report += `Empresa escolhida: ${clickEvent.data.companyName}\n`;
			report += `CompanyID: ${clickEvent.data.companyId}\n`;
		}
		report += '\n';

		// 1.3 requestTenantJWT()
		report += `1.3 requestTenantJWT()\n`;
		report += '-'.repeat(40) + '\n';
		const requestEvents = this.events.filter(e => e.section === 'REQUEST_TENANT_JWT');
		requestEvents.forEach(event => {
			if (event.subsection === 'BEFORE') {
				report += `Timestamp: ${event.timestamp}ms\n`;
				report += `URL: ${event.data.url}\n`;
				report += `Método: ${event.data.method}\n`;
				report += `Headers enviados: ${JSON.stringify(event.data.headers, null, 2)}\n`;
				report += `Authorization enviado: ${event.data.authorization}\n`;
				report += `Credentials utilizado: ${event.data.credentials}\n`;
				report += `Body enviado: ${event.data.body}\n`;
			} else if (event.subsection === 'AFTER') {
				report += `\nTimestamp: ${event.timestamp}ms\n`;
				report += `Status HTTP: ${event.data.status}\n`;
				report += `Headers recebidos: ${JSON.stringify(event.data.headers, null, 2)}\n`;
				report += `Set-Cookie recebido: ${event.data.setCookie || 'N/A'}\n`;
				report += `Body completo recebido: ${JSON.stringify(event.data.body, null, 2)}\n`;
				report += `\nToken JWT recebido: ${event.data.token ? 'SIM' : 'NÃO'}\n`;
				report += `Tamanho: ${event.data.tokenSize || 0} bytes\n`;
				report += `Payload decodificado: ${JSON.stringify(event.data.decodedPayload, null, 2)}\n`;
			}
		});
		report += '\n';

		// 1.4 hydrateContext()
		report += `1.4 hydrateContext()\n`;
		report += '-'.repeat(40) + '\n';
		const hydrateEvents = this.events.filter(e => e.section === 'HYDRATE_CONTEXT');
		hydrateEvents.forEach(event => {
			report += `Timestamp: ${event.timestamp}ms\n`;
			report += `Fase: ${event.subsection}\n`;
			report += `document.cookie: ${event.data.documentCookie}\n`;
			if (event.data.cookiesSeparated) {
				report += `Cookies separados: ${JSON.stringify(event.data.cookiesSeparated, null, 2)}\n`;
			}
			if (event.data.authTokenFound !== undefined) {
				report += `auth_token encontrado?: ${event.data.authTokenFound}\n`;
			}
			if (event.data.platformAuthTokenFound !== undefined) {
				report += `platform_auth_token encontrado?: ${event.data.platformAuthTokenFound}\n`;
			}
			report += '\n';
		});
		report += '\n';

		// 1.5 All document.cookie writes
		report += `1.5 Todas as escritas em document.cookie\n`;
		report += '-'.repeat(40) + '\n';
		if (this.cookieWrites.length === 0) {
			report += 'Nenhuma escrita registrada\n';
		} else {
			this.cookieWrites.forEach(write => {
				report += `Timestamp: ${write.timestamp}ms\n`;
				report += `Arquivo: ${write.file}\n`;
				report += `Linha: ${write.line}\n`;
				report += `Valor escrito: ${write.value}\n\n`;
			});
		}
		report += '\n';

		// 1.6 All cookie removals
		report += `1.6 Todas as remoções de cookie\n`;
		report += '-'.repeat(40) + '\n';
		if (this.cookieRemovals.length === 0) {
			report += 'Nenhuma remoção registrada\n';
		} else {
			this.cookieRemovals.forEach(removal => {
				report += `Timestamp: ${removal.timestamp}ms\n`;
				report += `Arquivo: ${removal.file}\n`;
				report += `Linha: ${removal.line}\n`;
				report += `Cookie removido: ${removal.cookie}\n\n`;
			});
		}
		report += '\n';

		// 1.7 All navigations
		report += `1.7 Todas as navegações\n`;
		report += '-'.repeat(40) + '\n';
		if (this.navigations.length === 0) {
			report += 'Nenhuma navegação registrada\n';
		} else {
			this.navigations.forEach(nav => {
				report += `Timestamp: ${nav.timestamp}ms\n`;
				report += `Tipo: ${nav.type}\n`;
				if (nav.url) report += `URL: ${nav.url}\n`;
				if (nav.from) report += `De: ${nav.from}\n`;
				if (nav.to) report += `Para: ${nav.to}\n`;
				report += '\n';
			});
		}
		report += '\n';

		// 1.8 All fetch requests
		report += `1.8 Todas as requisições fetch()\n`;
		report += '-'.repeat(40) + '\n';
		if (this.fetchCalls.length === 0) {
			report += 'Nenhum fetch registrado\n';
		} else {
			this.fetchCalls.forEach(fetch => {
				report += `Timestamp: ${fetch.timestamp}ms\n`;
				report += `URL: ${fetch.url}\n`;
				report += `Método: ${fetch.method}\n`;
				report += `Credentials: ${fetch.credentials}\n`;
				report += `Headers: ${JSON.stringify(fetch.headers, null, 2)}\n`;
				report += `Cookies antes: ${fetch.cookiesBefore}\n`;
				if (fetch.status !== undefined) {
					report += `Status: ${fetch.status}\n`;
					report += `Headers recebidos: ${JSON.stringify(fetch.headersReceived, null, 2)}\n`;
					report += `Body recebido: ${JSON.stringify(fetch.bodyReceived, null, 2)}\n`;
				}
				report += '\n';
			});
		}
		report += '\n';

		// BACKEND SECTION
		report += `${separator}\n`;
		report += `2. BACKEND\n`;
		report += `${separator}\n\n`;

		// 2.1 impersonation/start
		report += `2.1 impersonation/start\n`;
		report += '-'.repeat(40) + '\n';
		const impersonationEvents = this.events.filter(e => e.category === 'BACKEND' && e.section === 'IMPERSONATION_START');
		impersonationEvents.forEach(event => {
			report += `Timestamp: ${event.timestamp}ms\n`;
			report += `Subseção: ${event.subsection}\n`;
			report += `Dados: ${JSON.stringify(event.data, null, 2)}\n\n`;
		});
		report += '\n';

		// 2.2 Middleware
		report += `2.2 Middleware\n`;
		report += '-'.repeat(40) + '\n';
		const middlewareEvents = this.events.filter(e => e.category === 'BACKEND' && e.section === 'MIDDLEWARE');
		middlewareEvents.forEach(event => {
			report += `Timestamp: ${event.timestamp}ms\n`;
			report += `Subseção: ${event.subsection}\n`;
			report += `Dados: ${JSON.stringify(event.data, null, 2)}\n\n`;
		});
		report += '\n';

		// 2.3 companies/{id}
		report += `2.3 companies/{id}\n`;
		report += '-'.repeat(40) + '\n';
		const companyEvents = this.events.filter(e => e.category === 'BACKEND' && e.section === 'COMPANIES_ID');
		companyEvents.forEach(event => {
			report += `Timestamp: ${event.timestamp}ms\n`;
			report += `Subseção: ${event.subsection}\n`;
			report += `Dados: ${JSON.stringify(event.data, null, 2)}\n\n`;
		});
		report += '\n';

		// SUMMARY SECTION
		report += `${separator}\n`;
		report += `3. RESUMO\n`;
		report += `${separator}\n\n`;

		report += this.generateSummary();

		return report;
	}

	private generateSummary(): string {
		let summary = '';
		let stepNumber = 1;

		const addStep = (description: string) => {
			summary += `PASSO ${stepNumber}\n`;
			summary += `${description}\n\n`;
			stepNumber++;
		};

		// Analyze token lifecycle
		const initialCookies = this.events.find(e => e.section === 'INITIAL_STATE');
		const requestBefore = this.events.find(e => e.section === 'REQUEST_TENANT_JWT' && e.subsection === 'BEFORE');
		const requestAfter = this.events.find(e => e.section === 'REQUEST_TENANT_JWT' && e.subsection === 'AFTER');
		const hydrateEvents = this.events.filter(e => e.section === 'HYDRATE_CONTEXT');

		addStep('Estado inicial capturado');

		if (requestBefore) {
			addStep(`requestTenantJWT() iniciado - Token Platform: ${requestBefore.data.authorization ? 'Presente' : 'Ausente'}`);
		}

		if (requestAfter) {
			if (requestAfter.data.token) {
				addStep(`Tenant JWT recebido do backend - Tamanho: ${requestAfter.data.tokenSize} bytes`);
			} else {
				addStep('Tenant JWT NÃO recebido do backend - FALHA');
			}
		}

		if (hydrateEvents.length > 0) {
			const beforeHydrate = hydrateEvents.find(e => e.subsection === 'ANTES');
			const afterHydrate = hydrateEvents.find(e => e.subsection === 'IMEDIATAMENTE_DEPOIS');
			const after100ms = hydrateEvents.find(e => e.subsection === 'DEPOIS_100MS');
			const after300ms = hydrateEvents.find(e => e.subsection === 'DEPOIS_300MS');
			const after1000ms = hydrateEvents.find(e => e.subsection === 'DEPOIS_1000MS');

			if (beforeHydrate) {
				addStep('hydrateContext() - Antes de definir cookie');
			}

			if (afterHydrate) {
				if (afterHydrate.data.authTokenFound) {
					addStep('hydrateContext() - Cookie auth_token encontrado IMEDIATAMENTE após escrita');
				} else {
					addStep('hydrateContext() - Cookie auth_token NÃO encontrado IMEDIATAMENTE após escrita - POSSÍVEL PROBLEMA');
				}
			}

			if (after100ms) {
				if (after100ms.data.authTokenFound) {
					addStep('hydrateContext() - Cookie auth_token presente após 100ms');
				} else {
					addStep('hydrateContext() - Cookie auth_token AUSENTE após 100ms - COOKIE REMOVIDO');
				}
			}

			if (after300ms) {
				if (after300ms.data.authTokenFound) {
					addStep('hydrateContext() - Cookie auth_token presente após 300ms');
				} else {
					addStep('hydrateContext() - Cookie auth_token AUSENTE após 300ms - COOKIE REMOVIDO');
				}
			}

			if (after1000ms) {
				if (after1000ms.data.authTokenFound) {
					addStep('hydrateContext() - Cookie auth_token presente após 1000ms');
				} else {
					addStep('hydrateContext() - Cookie auth_token AUSENTE após 1000ms - COOKIE REMOVIDO');
				}
			}
		}

		// Analyze cookie writes
		if (this.cookieWrites.length > 0) {
			addStep(`Total de escritas em document.cookie: ${this.cookieWrites.length}`);
			this.cookieWrites.forEach(write => {
				addStep(`Escrita em ${write.file}:${write.line} - ${write.value.substring(0, 50)}...`);
			});
		}

		// Analyze cookie removals
		if (this.cookieRemovals.length > 0) {
			addStep(`Total de remoções de cookie: ${this.cookieRemovals.length}`);
			this.cookieRemovals.forEach(removal => {
				addStep(`Remoção em ${removal.file}:${removal.line} - ${removal.cookie}`);
			});
		}

		// Analyze fetch calls
		if (this.fetchCalls.length > 0) {
			addStep(`Total de requisições fetch: ${this.fetchCalls.length}`);
			this.fetchCalls.forEach(fetch => {
				if (fetch.status && fetch.status >= 400) {
					addStep(`Fetch FALHOU: ${fetch.method} ${fetch.url} - Status: ${fetch.status}`);
				}
			});
		}

		// Conclusion
		summary += 'CONCLUSÃO\n';
		summary += '-'.repeat(40) + '\n';

		const tokenReceived = requestAfter?.data.token;
		const tokenPersistedAfterHydrate = hydrateEvents.find(e => e.subsection === 'IMEDIATAMENTE_DEPOIS')?.data.authTokenFound;
		const tokenPersistedAfter100ms = hydrateEvents.find(e => e.subsection === 'DEPOIS_100MS')?.data.authTokenFound;

		if (!tokenReceived) {
			summary += 'O Tenant JWT NÃO foi recebido do backend.\n';
			summary += 'CAUSA MAIS PROVÁVEL: Backend não gerou ou não enviou o token corretamente.\n';
			summary += 'Confiança: 90%\n';
		} else if (!tokenPersistedAfterHydrate) {
			summary += 'O Tenant JWT foi recebido mas NÃO persistiu no cookie após hydrateContext().\n';
			summary += 'CAUSA MAIS PROVÁVEL: Problema na escrita do cookie ou remoção imediata.\n';
			summary += 'Confiança: 85%\n';
		} else if (!tokenPersistedAfter100ms) {
			summary += 'O Tenant JWT persistiu inicialmente mas foi removido dentro de 100ms.\n';
			summary += 'CAUSA MAIS PROVÁVEL: Outro código removeu o cookie (verificar seções 1.5 e 1.6).\n';
			summary += 'Confiança: 80%\n';
		} else {
			summary += 'O Tenant JWT foi recebido e persistiu corretamente.\n';
			summary += 'O problema pode estar em outra parte do fluxo (autenticação subsequente).\n';
			summary += 'Confiança: 70%\n';
		}

		return summary;
	}

	async saveReport() {
		// Trigger backend to generate the combined report
		try {
			const response = await fetch('http://localhost:8080/api/forensic/generate', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				}
			});
			
			if (response.ok) {
				const result = await response.json();
				console.log(`📋 ${result.message}`);
			} else {
				console.error('Erro ao gerar relatório forense no backend');
			}
		} catch (error) {
			console.error('Erro ao gerar relatório forense:', error);
		}
	}

	async clear() {
		this.events = [];
		this.cookieWrites = [];
		this.cookieRemovals = [];
		this.fetchCalls = [];
		this.navigations = [];
		this.startTime = Date.now();
		
		// Also clear backend logs
		try {
			await fetch('http://localhost:8080/api/forensic/clear', {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				}
			});
		} catch (error) {
			// Silently fail
		}
	}
}

export const forensicLogger = new ForensicLogger();
