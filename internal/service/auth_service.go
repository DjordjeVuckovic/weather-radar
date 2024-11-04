package service

type AuthCredentials struct {
	Username string
	Password string
}

type AuthService struct {
	basicAuthAdminCreds AuthCredentials
}

func NewAuthService(adminCreds AuthCredentials) *AuthService {
	return &AuthService{
		basicAuthAdminCreds: adminCreds,
	}
}

func (a *AuthService) ValidateBasicAuth(creds AuthCredentials) bool {
	return creds.Username == a.basicAuthAdminCreds.Username && creds.Password == a.basicAuthAdminCreds.Password
}
