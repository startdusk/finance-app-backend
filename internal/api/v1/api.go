package v1

import (
	"net/http"

	"github.com/startdusk/finance-app-backend/internal/api/auth"
)

type API struct {
	Path            string
	Method          string
	Func            http.HandlerFunc
	permissionTypes []auth.PermissionType
}

func NewAPI(method string, path string, handlerFunc http.HandlerFunc, permissionTypes ...auth.PermissionType) API {
	return API{
		Path:            path,
		Method:          method,
		Func:            handlerFunc,
		permissionTypes: permissionTypes,
	}
}
