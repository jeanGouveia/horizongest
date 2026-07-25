// Svelte 5 Runes: $state reativo e global via module singleton
// Importe { companyStore } em qualquer componente para acessar/mutar a empresa atual

interface CompanySettings {
	name: string;
	// Adicione outros campos conforme necessário
}

function createCompanyStore() {
  let company = $state<CompanySettings | null>(null);
  let loading = $state(true);

  return {
    get company() { return company; },
    get loading() { return loading; },

    setCompany(c: CompanySettings | null) { company = c; },
    setLoading(v: boolean) { loading = v; },

    clear() { company = null; }
  };
}

export const companyStore = createCompanyStore();
