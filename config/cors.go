package config

var AllowedOrigins = []string{
	"http://localhost:1506",
	"http://localhost:5173",
}

var GetAllowedOrigins = func() []string {
	return AllowedOrigins
}