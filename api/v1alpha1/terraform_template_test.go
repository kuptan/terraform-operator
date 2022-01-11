package v1alpha1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terraform Module", func() {
	expectedFile := `terraform {
	
		required_version = "~> 1.0.2"
	}
	variable "length" {}
	
	## additional-blocks
	
	module "operator" {
		source = "IbraheemAlSaady/test/module"
		version = "0.0.1"
		length = var.length
	}`
	Context("Terraform Template", func() {
		It("should generate the final module", func() {
			run := &Terraform{
				Spec: TerraformSpec{
					TerraformVersion: "1.0.2",
					Module: Module{
						Source:  "IbraheemAlSaady/test/module",
						Version: "0.0.1",
					},
					Variables: []Variable{
						Variable{
							Key:   "length",
							Value: "16",
						},
					},
					Destroy:             false,
					DeleteCompletedJobs: false,
				},
			}

			tpl, err := getTerraformModuleFromTemplate(run)

			tplString := string(tpl)

			Expect(err).ToNot(HaveOccurred())

			Expect(tplString).To(Equal(expectedFile))
		})
	})
})
