package port

type UserInfo struct {
	ID            string
	Email         string
	VerifiedEmail bool
	Name          string
	GivenName     string
	FamilyName    string
	Picture       string
}

type OAuthService interface {
	GetAccessToken(code string) (string, error)
	GetUserInfo(accessToken string) (*UserInfo, error)
	RevokeAccessToken(accessToken string) error
}
