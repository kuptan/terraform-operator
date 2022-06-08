package v1alpha1

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// returns a bool on whether a string is available in a given array of string
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// removes a string from a given array of string
func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// generates a random alphanumeric based on the length provided
func random(n int64) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz123456790")

	b := make([]rune, n)

	for i := range b {
		generated, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))

		b[i] = letters[generated.Int64()]
	}
	return string(b)
}

// returns common labels to be attached to children resources
func getCommonLabels(name string, runID string) map[string]string {
	return map[string]string{
		"terraformRunName": name,
		"terraformRunId":   runID,
		"component":        "Terraform-run",
		"owner":            "run.terraform-operator.io",
	}
}

func truncateResourceName(s string, i int) string {
	name := s
	if len(s) > i {
		name = s[0:i]
		// End in alphanum, Assume only "-" and "." can be in name
		name = strings.TrimRight(name, "-")
		name = strings.TrimRight(name, ".")
	}
	return name
}

// creates a name for the terraform Run job
func getUniqueResourceName(name string, runID string) string {
	// return fmt.Sprintf("tf-apply-%s-%s", name, runId)

	return fmt.Sprintf("%s-%s", truncateResourceName(name, 220), runID)
}
