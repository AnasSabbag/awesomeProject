package middlewares

import (
	"awesomeProject/global"
	"awesomeProject/models"
	"awesomeProject/util"
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	userContext       = "user_context"
	permissionContext = "permission_context"
)

func EnableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)

	})
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Token is empty")
			return
		}
		user := models.User{}
		err := global.DB.Get(&user, `
				select u.id,name,email,username,role_id,is_admin
				from users u
				join 
				sessions s
				on s.user_id = u.id
				where s.id = $1
				and u.is_deactivated = false
			`, token)

		if err != nil {
			log.Println("auth-handler user not found : ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		var expiryTime time.Time
		err = global.DB.Get(&expiryTime, `
				select expiry_time
				from sessions
				where id =$1
				and user_id = $2
				`, token, user.ID)
		if err != nil {
			log.Println("auth-handler: expiry time not getting:  ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if expiryTime.Before(time.Now()) {

			log.Println("auth-handler: session expired ", expiryTime, " current time ", time.Now())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		timeOut := os.Getenv("TIME_OUT")
		sessionTimeOut, err := strconv.Atoi(timeOut)
		if err != nil {
			log.Println("auth-handler: session timeout not getting: ", err)
			expiryTime = time.Now().Add(30 * time.Minute)
		} else {
			expiryTime = time.Now().Add(time.Duration(sessionTimeOut) * time.Minute)
		}
		_, err = global.DB.Exec(`update sessions
				set 
				    expiry_time = $1
				    where 
				        user_id =$2
					and id =$3 
				`, expiryTime, user.ID, token)
		if err != nil {
			log.Println("auth-handler: session expiry time update error: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userContext, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) *models.User {
	user, ok := ctx.Value(userContext).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func PermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r.Context())
		permission, err := util.GetPermissionObject(user.RoleID, user.ID)
		if err != nil {
			log.Println("permission-middleware: get permission object error: ", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), permissionContext, &permission)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func GetPermission(ctx context.Context) *[]models.Permission {
	permission, ok := ctx.Value(permissionContext).(*[]models.Permission)
	if !ok {
		return nil
	}
	return permission
}
