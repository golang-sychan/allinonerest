package userservice

import (
	restful "github.com/emicklei/go-restful/v3"
	"net/http"
)

//User domain object
type User struct {
	Id   string
	Name string
}

func FindUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("user-id")
	usr := &User{
		Id:   id,
		Name: "Sychan Joey",
	}
	response.WriteEntity(usr)
}

func UpdateUser(request restful.Request, response restful.Response) {
	usr := new(User)
	err := request.ReadEntity(&usr)
	if err != nil {
		usr.Name += ".sychan"
		response.WriteEntity(usr)
		return
	}
	response.WriteError(http.StatusInternalServerError, err)
}
