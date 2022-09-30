package base

var JwtConfig = jwtConfig{
	Days:      24,
	SecretKey: "default key",
	Issuer:    "default issuer",
}

type jwtConfig struct {
	SecretKey string
	Issuer    string
	Days      int
}
