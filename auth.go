package main

func BearerAuth(accessToken string) *Auth {
	return &Auth{
		BearerToken: accessToken,
	}
}
