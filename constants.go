package pg

type EnvironmentType string

const (
	// Production environment type for go live to real customer
	Production EnvironmentType = "production"

	// SandBox is staging (for development)
	SandBox EnvironmentType = "sandBox"
)
