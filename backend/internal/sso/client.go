package sso

import "net/http"

func (a *Auth) Require(handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valid, _, err := a.Validate(r)
		if err != nil || !valid {
			http.Error(w, errorAccessDenied, http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

func (a *Auth) RequireRoles(roles []string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valid, token, err := a.Validate(r)
		if err != nil || !valid {
			http.Error(w, errorAccessDenied, http.StatusUnauthorized)
			return
		}
		if token != nil {
			user, err := a.User(token)
			if err != nil {
				http.Error(w, errorRetrievingClaims, http.StatusInternalServerError)
				return
			}

			for _, role := range roles {
				if user.HasRole(role) {
					handler(w, r)
					return
				}
			}
			http.Error(w, errorMissingRole, http.StatusUnauthorized)

		}
		http.Error(w, errorAccessDenied, http.StatusUnauthorized)
	}
}
