package config

var AllowedOrigins = []string{
	"https://emo-back.onrender.com/",
	"https://emobuddy-495ef.web.app/",
}

var GetAllowedOrigins = func() []string {
	return AllowedOrigins
}