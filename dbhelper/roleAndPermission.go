package dbhelper

import (
	"awesomeProject/global"
	"awesomeProject/models"
	"github.com/jmoiron/sqlx"
)

func AddNewRole(newRole models.NewRolePayload) error {
	args := []interface{}{newRole.Name,
		newRole.Description}

	SQL := `insert into user_role(name,description) values($1,$2)`
	_, err := global.DB.DB.Exec(SQL, args...)
	return err
}

func CreatePermission(permission models.PermissionPayload) error {
	SQL := `insert into user_permission(name,description) values($1,$2);`
	_, err := global.DB.DB.Exec(SQL, permission.Name, permission.Description)
	return err
}
func ArchivedOldPermissionsOfRole(db sqlx.Ext, roleId int) error {
	SQL := `update user_role_permission set archived=true;`
	if _, err := db.Exec(SQL); err != nil {
		return err
	}
	return nil
}

func AddNewRoleWithPermission(tx *sqlx.Tx, roleId int, permissionId []int) error {

	if len(permissionId) == 0 {
		return nil
	}
	size := 65535 / len(permissionId)

	//go get github.com/thoas/go-funk
	chunkData := funk.Chunk(permissionId, size).([][]int)
	var err error

	for i := range chunkData {
		SQL := `insert into user_role_permission_relation(role_id,permission_id)values %s`
		SQL = global.SetupBindVars(SQL, "(?,?)", len(chunkData[i]))
		values := make([]interface{}, 0)

		for j := range chunkData[i] {
			values = append(values, roleId, chunkData[i][j])
		}
		_, err = tx.Exec(SQL, values...)
		if err != nil {
			return err
		}
	}

	return nil
}
