package domain

// Module represents a registered module in the platform
type Module struct {
	Key          string   // Unique module key (e.g., "finance", "inventory")
	Name         string   // Display name (e.g., "Financeiro")
	Description  string   // Module description
	Route        string   // Base route (e.g., "/finance")
	Icon         string   // Icon name (e.g., "dollar-sign")
	FeatureFlag  string   // Feature flag key in GlobalConfig (e.g., "enable_finance")
	Version      string   // Module version (e.g., "1.0.0")
	Dependencies []string // List of module keys this module depends on
	Status       string   // Status: "active", "deprecated", "experimental"
}

// ModuleRegistry is the central registry for all platform modules
// Every new module must be registered here
var ModuleRegistry = map[string]Module{
	"finance": {
		Key:          "finance",
		Name:         "Financeiro",
		Description:  "Gestão financeira, contas a pagar/receber, fluxo de caixa",
		Route:        "/finance",
		Icon:         "dollar-sign",
		FeatureFlag:  "enable_finance",
		Version:      "1.0.0",
		Dependencies: []string{},
		Status:       "active",
	},
	"purchasing": {
		Key:          "purchasing",
		Name:         "Compras",
		Description:  "Gestão de compras, fornecedores, pedidos de compra",
		Route:        "/purchasing",
		Icon:         "shopping-cart",
		FeatureFlag:  "enable_purchasing",
		Version:      "1.0.0",
		Dependencies: []string{},
		Status:       "active",
	},
	"inventory": {
		Key:          "inventory",
		Name:         "Estoque",
		Description:  "Gestão de estoque, produtos, categorias, movimentações",
		Route:        "/inventory",
		Icon:         "package",
		FeatureFlag:  "enable_inventory",
		Version:      "1.0.0",
		Dependencies: []string{},
		Status:       "active",
	},
	"crm": {
		Key:          "crm",
		Name:         "CRM",
		Description:  "Gestão de relacionamento com clientes",
		Route:        "/crm",
		Icon:         "users",
		FeatureFlag:  "enable_crm",
		Version:      "1.0.0",
		Dependencies: []string{},
		Status:       "active",
	},
	"calendar": {
		Key:          "calendar",
		Name:         "Agenda",
		Description:  "Gestão de agenda, compromissos, eventos",
		Route:        "/calendar",
		Icon:         "calendar",
		FeatureFlag:  "enable_calendar",
		Version:      "1.0.0",
		Dependencies: []string{},
		Status:       "active",
	},
	"pos": {
		Key:          "pos",
		Name:         "PDV",
		Description:  "Ponto de venda para atendimento ao cliente",
		Route:        "/pos",
		Icon:         "monitor",
		FeatureFlag:  "enable_pos",
		Version:      "1.0.0",
		Dependencies: []string{"inventory"},
		Status:       "active",
	},
	"ai": {
		Key:          "ai",
		Name:         "IA",
		Description:  "Inteligência artificial para previsões e recomendações",
		Route:        "/ai",
		Icon:         "brain",
		FeatureFlag:  "enable_ai",
		Version:      "1.0.0",
		Dependencies: []string{"inventory", "finance"},
		Status:       "experimental",
	},
	"delivery": {
		Key:          "delivery",
		Name:         "Delivery",
		Description:  "Gestão de entregas, rastreamento, logística",
		Route:        "/delivery",
		Icon:         "truck",
		FeatureFlag:  "enable_delivery",
		Version:      "1.0.0",
		Dependencies: []string{"inventory"},
		Status:       "active",
	},
	"marketplace": {
		Key:          "marketplace",
		Name:         "Marketplace",
		Description:  "Marketplace para integração com fornecedores",
		Route:        "/marketplace",
		Icon:         "store",
		FeatureFlag:  "enable_marketplace",
		Version:      "1.0.0",
		Dependencies: []string{"purchasing", "inventory"},
		Status:       "experimental",
	},
}

// GetModule retrieves a module by key
func GetModule(key string) (Module, bool) {
	module, exists := ModuleRegistry[key]
	return module, exists
}

// GetAllModules returns all registered modules
func GetAllModules() []Module {
	modules := make([]Module, 0, len(ModuleRegistry))
	for _, module := range ModuleRegistry {
		modules = append(modules, module)
	}
	return modules
}

// GetActiveModules returns all active modules
func GetActiveModules() []Module {
	modules := make([]Module, 0)
	for _, module := range ModuleRegistry {
		if module.Status == "active" {
			modules = append(modules, module)
		}
	}
	return modules
}

// GetModuleDependencies returns all dependencies for a module (recursive)
func GetModuleDependencies(key string) []string {
	dependencies := make([]string, 0)
	visited := make(map[string]bool)
	
	var visit func(string)
	visit = func(k string) {
		if visited[k] {
			return
		}
		visited[k] = true
		
		module, exists := GetModule(k)
		if !exists {
			return
		}
		
		for _, dep := range module.Dependencies {
			if !visited[dep] {
				dependencies = append(dependencies, dep)
				visit(dep)
			}
		}
	}
	
	visit(key)
	return dependencies
}
