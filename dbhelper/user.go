package dbhelper

import (
	"awesomeProject/global"
	"awesomeProject/models"
)

func CreateUser(user models.User) error {
	SQL := ` insert into users(name,email,username,password,role_id,is_admin)values($1,$2,$3,$4,$5,$6);`
	_, err := global.DB.Exec(SQL, user.Name, user.Email, user.Username, user.Password, user.RoleID, user.IsAdmin)
	return err
}

func GetUserInfo(username string) (models.User, error) {
	var user models.User
	//todo!!!
	SQL := ` select id,name,username,email,password from users 
                where username=$1
                	and is_deactivated=false
                	and archived_at is null;`
	err := global.DB.Get(&user, SQL, username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
