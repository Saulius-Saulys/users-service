package environment

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEnv(t *testing.T) {
	restore := setEnvVars(map[string]string{
		"ENVIRONMENT": "local",
	})
	defer restore()

	assert.Equal(t, "local", NewEnv().Environment)
}

// setEnvVars sets given environment variables and provides a callback function to restore the variables to their initial values
func setEnvVars(vars map[string]string) (restoreVars func()) {
	initialValues := map[string]string{}
	unsetVars := map[string]bool{}

	for name, value := range vars {
		initialValue, exists := os.LookupEnv(name)
		if exists {
			initialValues[name] = initialValue
		} else {
			unsetVars[name] = true
		}

		err := os.Setenv(name, value)
		if err != nil {
			panic(err)
		}
	}

	return func() {
		for name, value := range initialValues {
			err := os.Setenv(name, value)
			if err != nil {
				panic(err)
			}
		}

		for name := range unsetVars {
			err := os.Unsetenv(name)
			if err != nil {
				panic(err)
			}
		}
	}
}
