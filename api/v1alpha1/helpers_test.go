package v1alpha1

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Helpers", func() {
	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	Context("Helpers functions", func() {
		It("should return true if string is available in an array of strings", func() {
			arr := []string{"abc", "def"}

			Expect(containsString(arr, "abc")).To(BeTrue())
			Expect(containsString(arr, "123")).To(BeFalse())
		})

		It("should remove a string from an array", func() {
			arr := []string{"abc", "def"}

			arr = removeString(arr, "abc")

			Expect(arr).To(HaveLen(1))
			Expect(arr[0]).To(Equal("def"))
		})

		It("should generate a random string with a specified length", func() {
			str := random(5)

			Expect(str).To(HaveLen(5))
		})

		It("should return common labels", func() {
			labels := getCommonLabels("foo", "1234")

			Expect(labels["terraformRunName"]).ToNot(BeEmpty())
			Expect(labels["terraformRunId"]).ToNot(BeEmpty())
			Expect(labels["component"]).ToNot(BeEmpty())
			Expect(labels["owner"]).ToNot(BeEmpty())
		})

		It("should return the unique resource name", func() {
			name := getUniqueResourceName("foo", "1234")
			Expect(name).To(Equal("foo-1234"))
		})

		It("should return the the name of the secret output", func() {
			name := getOutputSecretname("foo")
			Expect(name).To(Equal("foo-outputs"))
		})
	})
})
