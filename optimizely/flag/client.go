package flag

type FlagClient interface {
	CreateFlag(feat Flag) (Flag, error)
	CreateRuleset(flag Flag) error
	EnableRuleset(feat Flag) error
	CreateVariation(flag Flag, variation Variation) error
	DeleteFlag(projectId int, flagKey string) error
}
