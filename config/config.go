package config

func Get(key string) interface{} {
	switch key {
	case "jwtsecret":
		return []byte("secret")
	default:
		return nil
	}

}
