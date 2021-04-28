package gatherers

type FactCategory string

const(
	Application FactCategory = "Application"
	ApplicationTechnology FactCategory = "Application technology"
	Docker FactCategory = "Docker configuration"
	Drupal FactCategory = "Drupal configuration"
	EnvVar FactCategory = "Environment variable"
	Lagoon FactCategory = "Lagoon configuration"
	Laravel FactCategory = "Laravel configuration"
	Network FactCategory = "Network configuration"
	Platform FactCategory = "Platform"
	ProgrammingLanguage FactCategory = "Programming language"
)
