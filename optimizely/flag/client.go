package flag

type FlagClient interface {
	CreateFlag(flag Flag) (Flag, error)
	DeleteFlag(projectId int, flagKey string) error

	CreateRuleset(flag Flag) error
	EnableRuleset(flag Flag) error
	DisableRuleset(flag Flag) error

	CreateVariation(flag Flag, variation Variation) error
}
