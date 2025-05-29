package utils

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const csrfCookieName = "csrf_token"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateCSRFToken()(string,error){
	b := make([]byte,32)
	_,err := rand.Read(b)

	if(err!=nil){
		return "",err
	}

	return base64.URLEncoding.EncodeToString(b),nil
}

func SetCSRFCookie(w http.ResponseWriter,token string){
	http.SetCookie(w,&http.Cookie{
		Name:     csrfCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(1 * time.Hour),		
	})
}

func GetCSRFCookie(r * http.Request)(string,error){
	cookie,err := r.Cookie(csrfCookieName)
	if(err!=nil){
		return "",err
	}
	return cookie.Value,nil
}

func VerifyCSRFToken(r *http.Request)bool{
	formToken := r.FormValue("csrf_token")
	cookieToken,err := GetCSRFCookie(r)

	if(err!=nil){
		return false
	}
	return formToken == cookieToken
}