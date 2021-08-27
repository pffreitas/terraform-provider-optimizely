package audience

type AudienceClient interface {
	CreateAudience(aud Audience) (Audience, error)
	GetAudience(audId string) (Audience, error)
	ArchiveAudience(audId string) (Audience, error)
	UpdateAudience(aud Audience) (Audience, error)
}
