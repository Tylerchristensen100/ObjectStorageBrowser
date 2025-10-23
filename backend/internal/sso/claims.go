package sso

import "strings"

type Claims map[string]interface{}

func (c *Claims) HasRole(role string) bool {
	for _, r := range c.roles() {
		if strings.EqualFold(r, role) {
			return true
		}
	}
	return false
}

func (c *Claims) roles() []string {
	var roles []string
	claims := *c
	r := claims["roles"].([]string)

	roles = append(roles, r...)

	return roles
}
