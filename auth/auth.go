package auth

const (
	OAuth = iota
	NoAuth
)

var AuthMethod int

func SetAuthMethod(method int) {
	AuthMethod = method
}

