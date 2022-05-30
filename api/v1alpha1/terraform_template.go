package v1alpha1

import (
	"bytes"
	"text/template"
)

// getTerraformModuleFromTemplate generates the Terraform module template
func getTerraformModuleFromTemplate(run *Terraform) ([]byte, error) {
	tfTemplate, err := template.New("main.tf").Parse(`terraform {
		{{- if .Spec.Backend }}
		{{.Spec.Backend}}
		{{- end}}
	
		required_version = "~> {{.Spec.TerraformVersion}}"
	}

	{{- if .Spec.ProvidersConfig }}
	{{.Spec.ProvidersConfig}}
	{{- end}}
	
	{{- range .Spec.Variables}}
	{{- if not .EnvironmentVariable }}
	variable "{{.Key}}" {}
	{{- end}}
	{{- end}}
	
	## additional-blocks
	
	module "operator" {
		source = "{{.Spec.Module.Source}}"
		
		{{- if .Spec.Module.Version }}
		version = "{{.Spec.Module.Version}}"
		{{- end}}
	
		{{- range .Spec.Variables}}
		{{- if not .EnvironmentVariable }}
		{{.Key}} = var.{{.Key}}
		{{- end}}
		{{- end}}
	}
	
	{{- range .Spec.Outputs}}
	output "{{.Key}}" {
		value = module.operator.{{.ModuleOutputName}}
	}
	{{- end}}`)

	if err != nil {
		return nil, err
	}
	var tpl bytes.Buffer

	if err := tfTemplate.Execute(&tpl, run); err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}
