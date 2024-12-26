package middleware

import (
	"errors"
	"fmt"
	"gomailer/utils"
	"net/http"
)

func AuthMiddlewareUser(next func(w http.ResponseWriter, r *http.Request))func(w http.ResponseWriter, r *http.Request){

	return func(w http.ResponseWriter, r *http.Request){
		cookie, err := r.Cookie("auth_jwt")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				fmt.Println("No cookie found")
				http.Redirect(w,r,"/login",http.StatusSeeOther)
			default:
				fmt.Println(err)
				http.Redirect(w,r,"/login",http.StatusSeeOther)
			}
			return
		}
		
		user_email,jwtErr := utils.VerifyJWT(cookie.Value)


		if jwtErr != nil{
			fmt.Println(jwtErr)
			http.Redirect(w,r,"/login",http.StatusSeeOther)
			return
		}

		r.Header.Add("user_id",user_email)
		next(w,r)
	}
}