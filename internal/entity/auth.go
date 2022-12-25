package entity

type Credentials struct {
	Logins map[string]string
}

func (c *Credentials) CheckCredentials(login, pwd string) bool {
	if realPwd, ok := c.Logins[login]; ok {
		if pwd == realPwd {
			return true
		}
	}
	return false
}
