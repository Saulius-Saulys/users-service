package environment

import (
	"os"
	"reflect"
)

const (
	varTag = "var"
)

type Env struct {
	Environment              string `var:"ENVIRONMENT"`
	PostgreSQLDBHostname     string `var:"POSTGRESQL_DB_HOSTNAME"`
	PostgreSQLDBUser         string `var:"POSTGRESQL_DB_USER"`
	PostgreSQLDBPassword     string `var:"POSTGRESQL_DB_PASSWORD"`
	RabbitMQHostname         string `var:"RABBITMQ_HOSTNAME"`
	RabbitMQUser             string `var:"RABBITMQ_USER"`
	RabbitMQPassword         string `var:"RABBITMQ_PASSWORD"`
	OpenAIAzureAPIKey        string `var:"OPENAI_AZURE_API_KEY"`
	PromptBuilderAPIKey      string `var:"PROMPT_BUILDER_API_KEY"`
	HashedCompletionsAPIKeys string `var:"COMPLETIONS_API_KEYS"`
	ContentSafetyAPIKey      string `var:"CONTENT_SAFETY_API_KEY"`
}

func NewEnv() Env {
	var env Env

	v := reflect.ValueOf(&env).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envVar := os.Getenv(field.Tag.Get(varTag))
		v.FieldByName(field.Name).SetString(envVar)
	}

	return env
}
