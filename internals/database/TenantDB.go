package database

import "errors"

func ResolveTenant(tenant string) (string, string, error) {

	if tenant == "" {
		return "", "", errors.New("tenant code is missing")
	}

	switch tenant {

	case "staff", "test":
		return tenant, "postgres", nil

	case "db007", "db008":
		return tenant, "mysql", nil

	default:
		return "", "", errors.New("invalid tenant code")
	}
}
