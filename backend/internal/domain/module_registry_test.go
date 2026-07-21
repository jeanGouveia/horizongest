package domain

import (
	"testing"
)

func TestGetModule(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    bool
	}{
		{
			name: "existing module",
			key:  "finance",
			want: true,
		},
		{
			name: "non-existing module",
			key:  "nonexistent",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, exists := GetModule(tt.key)
			if exists != tt.want {
				t.Errorf("GetModule() exists = %v, want %v", exists, tt.want)
			}
		})
	}
}

func TestGetAllModules(t *testing.T) {
	modules := GetAllModules()
	if len(modules) == 0 {
		t.Error("GetAllModules() returned empty list")
	}

	// Verify all expected modules are present
	expectedModules := []string{
		"finance", "purchasing", "inventory", "crm", "calendar",
		"pos", "ai", "delivery", "marketplace",
	}

	moduleMap := make(map[string]bool)
	for _, module := range modules {
		moduleMap[module.Key] = true
	}

	for _, expected := range expectedModules {
		if !moduleMap[expected] {
			t.Errorf("GetAllModules() missing expected module: %s", expected)
		}
	}
}

func TestGetActiveModules(t *testing.T) {
	modules := GetActiveModules()
	if len(modules) == 0 {
		t.Error("GetActiveModules() returned empty list")
	}

	// Verify no experimental modules are included
	for _, module := range modules {
		if module.Status == "experimental" {
			t.Errorf("GetActiveModules() included experimental module: %s", module.Key)
		}
	}
}

func TestGetModuleDependencies(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected []string
	}{
		{
			name:     "module with no dependencies",
			key:      "finance",
			expected: []string{},
		},
		{
			name:     "module with direct dependencies",
			key:      "pos",
			expected: []string{"inventory"},
		},
		{
			name:     "module with transitive dependencies",
			key:      "ai",
			expected: []string{"inventory", "finance"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deps := GetModuleDependencies(tt.key)
			
			// Check if all expected dependencies are present
			expectedMap := make(map[string]bool)
			for _, exp := range tt.expected {
				expectedMap[exp] = true
			}
			
			for _, dep := range deps {
				if !expectedMap[dep] {
					t.Errorf("GetModuleDependencies() returned unexpected dependency: %s", dep)
				}
			}
			
			// Check if all expected dependencies are returned
			for _, exp := range tt.expected {
				found := false
				for _, dep := range deps {
					if dep == exp {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetModuleDependencies() missing expected dependency: %s", exp)
				}
			}
		})
	}
}
