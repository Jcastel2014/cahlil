package main

import (
	"net/url"

	"sysadmin.com/final/pkg/models"
)

type TemplateData struct {
	User            []*models.User
	ErrorsFromForm  map[string]string
	Flash           string
	FormData        url.Values
	IsAuthenticated bool
}
