package config

var AllowedOrigins = []string{
	"https://emo-back.onrender.com",
	"https://emobuddy-495ef.web.app",
	"http://localhost:1506",
	"http://localhost:5173",
}

var GetAllowedOrigins = func() []string {
	return AllowedOrigins
}