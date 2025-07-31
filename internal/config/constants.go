package config

import "os"

var (
	PROTOCOL_NAME = os.Getenv("PROTOCOL_NAME")
	SITE_FOR_TEST = os.Getenv("SITE_FOR_TEST")
)

const (
	FACEBOOK_URL                   = "https://www.facebook.com/"
	FACEBOOK_LOGIN_SELECTOR        = "#email"
	FACEBOOK_PASSWORD_SELECTOR     = "#pass"
	FACEBOOK_LOGIN_BUTTON_SELECTOR = "button[name='login']"
)
