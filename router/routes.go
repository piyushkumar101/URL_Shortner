package router

import (
	"net/http"
	"urlshortner/constant"
	"urlshortner/controller"
)

var urlShortner = Routes{
	Route{"Url Shortner Service", http.MethodPost, constant.UrlShortnerPath, controller.ShortUrl},
	Route{"Redirect to url", http.MethodGet, constant.RedirectUrlPath, controller.RedirectURL},
}
