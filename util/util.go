package util

import (
	"awesomeProject/global"
	"awesomeProject/models"
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/teris-io/shortid"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"math/big"
	"net/http"
	"time"
)

const (
	CanAddPermission     = 1
	CanAddRole           = 2
	CanCreateUser        = 3
	CanDeactivateUser    = 4
	CanUpdateRole        = 5
	CanAddUpdateUserRole = 6
	ListOfPermission     = 7
	ListOfRole           = 8
	EditProfile          = 9
)

type clientError struct {
	ID            string `json:"id"`
	Message       string `json:"message"`
	Err           string `json:"err"`
	StatusCode    int    `json:"statusCode"`
	IsClientError bool   `json:"isClientError"`
}

var generator *shortid.Shortid

const generatorSeed = 1000

func init() {
	n, err := rand.Int(rand.Reader, big.NewInt(generatorSeed))
	if err != nil {
		log.Panicf("failed to initialize utilitites with random seed. %v", err)
		return
	}
	g, err := shortid.New(1, shortid.DefaultABC, n.Uint64())
	if err != nil {
		log.Panicf("failed to initialize utilitites package  with random seed. %v", err)
	}
	generator = g

}

func GetPermissionObject(roleID, userID int) ([]models.Permission, error) {
	SQL := `
			select up.id,name,description,is_deleted
			from user_permission up 
			join user_role_permission_relation urpr
			on
			up.id = urpr.permission_id
			where role_id =$1 
			and 
			    is_deleted =false
				order by up.id
		`
	permissions := make([]models.Permission, 0)
	err := global.DB.Select(&permissions, SQL, roleID)
	return permissions, err
}

func RespondError(w http.ResponseWriter, statusCode int, err error, message string) {
	//log.Error(message)
	clientErr := newClientError(err, statusCode, message)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientErr); err != nil {
		log.Println("Respond Error:failed to send message to caller with error ", err)
	}
}

func newClientError(err error, statusCode int, message string) *clientError {
	errorId, _ := generator.Generate()
	var errString string
	if err != nil {
		errString = err.Error()
	}

	return &clientError{
		ID:            errorId,
		Message:       message,
		StatusCode:    statusCode,
		IsClientError: true,
		Err:           errString,
	}

}

func EmailAvailable(email string) (bool, error) {
	var checkEmail string

	SQL := `select email from users 
             where email = trim(lower($1))
             and archived_at is null and is_deactivated = false
				limit 1	;`
	err := global.DB.Get(&checkEmail, SQL, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return false, err
	}
	if checkEmail == email {
		return true, nil
	}
	return false, nil
}

func UsernameAvailable(username string) (bool, error) {
	var checkUsername string
	SQL := `select username from users 
             where username = trim(lower($1))
             and archived_at is null and is_deactivated = false;`
	err := global.DB.Get(&checkUsername, SQL, username)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if checkUsername == username {
		return true, nil
	}
	return false, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func ResponseJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Println("Response JSON  Error:failed to send message to caller with error ", err)
		}
	}
}

func GenerateSession() string {
	token := HashString(time.Now().String())
	return token
}

func HashString(s string) string {
	sha := sha512.New()
	sha.Write([]byte(s))
	return hex.EncodeToString(sha.Sum(nil))
}

func CreateSession(token string, userId int) error {
	SQL := `insert into sessions(id,user_id)values($1,$2);`
	_, err := global.DB.Exec(SQL, token, userId)
	return err
}

func UpdateUserLastLogin(userId int) error {
	SQL := `update users set last_login = now() where id = $1;`
	_, err := global.DB.Exec(SQL, userId)
	return err
}

func HasPermission(permissions []models.Permission, permission int) bool {
	for _, per := range permissions {
		if per.ID == permission {
			return true
		}
	}
	return false
}
func ParseBody(body io.Reader, out interface{}) error {
	return json.NewDecoder(body).Decode(out)
}
