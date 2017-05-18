package db

// User represents a row in the users table.
type User struct {
	Email    string
	Password []byte
}

const usersSchema = `
	email VARCHAR(255) NOT NULL,
	password BINARY(60) NOT NULL,
	PRIMARY KEY (email)
`

// GetUser gets the user with the given email.
func GetUser(email string) *User {
	row := db.QueryRow(`SELECT * FROM users WHERE email=?`, email)
	if row != nil {
		var email string
		var password []byte
		row.Scan(&email, &password)
		return &User{Email: email, Password: password}
	}
	return nil
}

// GetAuthTokens gets the auth tokens of the user.
func (user *User) GetAuthTokens() []AuthToken {
	results, err := db.Query(`SELECT * FROM authtokens WHERE user=? AND is_recovery=0`, user.Email)
	if err != nil {
		return []AuthToken{}
	}
	return scanAuthTokens(results)
}

// GetRecoveryTokens gets the password recovery tokens of the user.
func (user *User) GetRecoveryTokens() []AuthToken {
	results, err := db.Query(`SELECT * FROM authtokens WHERE user=? AND is_recovery=1`, user.Email)
	if err != nil {
		return []AuthToken{}
	}
	return scanAuthTokens(results)
}

// CheckAuthToken checks if the given authentication token is valid for this user.
func (user *User) CheckAuthToken(token string) bool {
	for _, at := range user.GetAuthTokens() {
		if at.Token == token {
			return true
		}
	}
	return false
}

// CheckRecoveryToken checks if the given recovery token is valid for this user.
func (user *User) CheckRecoveryToken(token string) bool {
	for _, at := range user.GetRecoveryTokens() {
		if at.Token == token {
			return true
		}
	}
	return false
}

// GetPermissions returns the permissions this user has.
func (user *User) GetPermissions() []*Permission {
	results, err := db.Query(`SELECT * FROM permissions WHERE owner=?`, user.Email)
	if err != nil {
		return []*Permission{}
	}
	return scanPermissions(results)
}
