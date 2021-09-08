package flag

type FlagClient interface {
	CreateFlag(flag Flag) (Flag, error)
	GetFlag(projectId int, flagKey string) (Flag, error)
	DeleteFlag(projectId int, flagKey string) error

	CreateRuleset(flag Flag) error
	GetRuleset(flag Flag) (map[string]FeatureEnvironment, error)
	EnableRuleset(flag Flag) error
	DisableRuleset(flag Flag) error

	CreateVariation(flag Flag, variation Variation) error
	GetVariation(projectId int, flagKey string) ([]Variation, error)
}
