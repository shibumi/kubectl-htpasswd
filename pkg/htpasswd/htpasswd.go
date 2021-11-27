package htpasswd

const (
	// FieldSeparator specifies the separator between user and hash in the htpasswd file
	FieldSeparator = ":"
)

// BuildEntry creates a htpasswd compatible entry with a user, password, algorithm
// and the cost for the algorithm. If the algorithm do not support cost, cost will
// be ignored.
func BuildEntry(user, password, algorithm string, cost int) (string, error) {
	switch algorithm {
	case "bcrypt":
		fallthrough
	default:
		hashedPw, err := Bcrypt([]byte(password), cost)
		if err != nil {
			return "", err
		}
		return user + FieldSeparator + string(hashedPw), nil
	}
}
