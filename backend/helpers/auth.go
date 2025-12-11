package helpers

func GetAuthenticatedUser(c *gin.Context) (*users.User, error) {
	val, exists := c.Get(UserKey)
	if !exists {
		return nil, fmt.Errorf("user key missing")
	}

	user, ok := val.(*users.User)
	if !ok || user == nil {
		return nil, fmt.Errorf("invalid user type")
	}

	return user, nil
}
