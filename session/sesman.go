package sesman

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	LoginUuid string
	userName  string
	Email     string
	Password  string
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateSession(db *sql.DB, user user, w http.ResponseWriter, r *http.Request) string {
	sesCookie := CreateSessionCookie()
	http.SetCookie(w, &sesCookie)
	w.WriteHeader(200)
	db.Exec("DELETE FROM Session where auth_uuid='" + user.LoginUuid + "'")
	db.Exec("INSERT INTO Session (uuid,auth_uuid)", sesCookie.Value, user.LoginUuid)
	return sesCookie.Value
}

func Getuser(session_uuid string, db *sql.DB) (user, bool) {
	var err error
	var user user
	if session_uuid != "" {
		err = db.QueryRow("SELECT auth_uuid FROM Session WHERE uuid = ?", session_uuid).Scan(&user.LoginUuid)
		if err == sql.ErrNoRows {
			return user, false
		} else {
			err = db.QueryRow("SELECT * FROM users WHERE uuid = ?", user.LoginUuid).Scan(&user.LoginUuid, &user.userName, &user.Email, &user.Password)
			if err == sql.ErrNoRows {
				return user, false
			} else {
				return user, true
			}
		}
	} else {
		return user, false
	}
}

func CreateSessionCookie() http.Cookie {
	var err error
	u1 := uuid.Must(uuid.NewV4(), err)
	return http.Cookie{
		Name:   "sessionCookie",
		Value:  u1.String(),
		MaxAge: 2 * int(time.Hour),
	}
}

func CheckSession(w http.ResponseWriter, r *http.Request, db *sql.DB) bool {
	// var err error
	cookie, err := r.Cookie("sessionCookie")
	if err != nil {
		return false
	} else {
		// cookie = sessionCookie
		session, _ := db.Query("SELECT * FROM Session WHERE uuid = '" + cookie.Value + "'")
		defer session.Close()
		var id int
		var sessionUuid string
		var authUuid string
		count := 0
		for session.Next() {
			session.Scan(&id, &sessionUuid, &authUuid)
			fmt.Fprintln(w, "session: ", id, sessionUuid, authUuid)
			count++
		}
		if count == 1 {
			return true
		} else {
			return false
		}
	}
}

func DeleteSession(sesid string, db *sql.DB) {
	db.Exec("DELETE FROM Session where uuid='" + sesid + "'")
}
