package app

import "github.com/golang-sychan/allinonerest/cmd/app"

const (
	commandDesc = `The application is a WEB server. In addition to provide services for front-end
           to get , update , create , delete user info `
)

func NewApp(name string) *app.App {
	a := app.NewApp(name, "allinonerest", app.WithDescription(commandDesc))
	return a
}
