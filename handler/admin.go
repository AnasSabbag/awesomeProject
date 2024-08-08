package handler

import (
	"awesomeProject/dbhelper"
	"awesomeProject/global"
	"awesomeProject/middlewares"
	"awesomeProject/models"
	"awesomeProject/util"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
)

func AdminRoutes(router *mux.Router) *mux.Router {

	router.Handle("/create-role", http.HandlerFunc(CreateRole)).Methods("POST")
	router.Handle("/update-role/{roleId}", http.HandlerFunc(UpdateRole)).Methods("POST")
	router.Handle("/create-permission", http.HandlerFunc(CreatePermission)).Methods("POST")
	//get-all role
	//update-role
	//create-permission
	//update-permission
	//get-all-permission
	//update-user-password
	//activate-user
	//deactivate-user

	return router
}

func CreateRole(w http.ResponseWriter, r *http.Request) {
	//permissions := middlewares.GetPermission(r.Context())
	//permissionRole := 25
	//if !util.HasPermission(*permissions, permissionRole) {
	//	util.RespondError(w, http.StatusForbidden, nil, "permission forbidden")
	//	return
	//}
	newRole := models.NewRolePayload{}
	if err := util.ParseBody(r.Body, &newRole); err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "body parse error")
		return
	}
	v := validator.New()
	err := v.Struct(newRole)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "body validate error")
		return
	}

	err = dbhelper.AddNewRole(newRole)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err, "add new role error")
		return
	}
	util.ResponseJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "role created successfully",
	})

}
func CreatePermission(w http.ResponseWriter, r *http.Request) {
	permissions := middlewares.GetPermission(r.Context())

	if !util.HasPermission(*permissions, util.CanAddPermission) {
		util.RespondError(w, http.StatusForbidden, errors.New("user not have permission create permission"), "permission create permission")
		return
	}
	payloadPermission := models.PermissionPayload{}
	if err := util.ParseBody(r.Body, &payloadPermission); err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "body parse error")
		return
	}
	v := validator.New()
	err := v.Struct(payloadPermission)
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "body validate error")
		return
	}
	err = dbhelper.CreatePermission(payloadPermission)
	if err != nil {
		util.RespondError(w, http.StatusInternalServerError, err, "add new permission error")
		return
	}
}
func UpdateRole(w http.ResponseWriter, r *http.Request) {
	permissions := middlewares.GetPermission(r.Context())
	if !util.HasPermission(*permissions, util.CanUpdateRole) {
		util.RespondError(w, http.StatusForbidden, errors.New("user not have permission update permission"), "permission update permission")
		return
	}
	roleId, err := strconv.Atoi(r.URL.Query().Get("roleId"))
	if err != nil {
		util.RespondError(w, http.StatusBadRequest, err, "body parse error")
		return
	}

	payloadPermission := models.NewRolePayload{}
	txErr := global.Tx(func(tx *sqlx.Tx) error {
		err = dbhelper.ArchivedOldPermissionsOfRole(tx, roleId)
		if err != nil {
			return err
		}
		err = dbhelper.AddNewRoleWithPermission(tx, roleId, payloadPermission.PermissionID)
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		util.RespondError(w, http.StatusInternalServerError, txErr, "tx error")
		return
	}
	//archived old permission and add new permission relation

}
