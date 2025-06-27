package integration

import "context"

func ClearUsers() {
	ctx := context.Background()
	_, err := dbConfig.ExecContext(ctx, "DELETE FROM users")
	if err != nil {
		logs.Log("Error clearing users: " + err.Error())
	} else {
		logs.Log("Users cleared successfully")
	}
}

func ClearUserProfiles() {
	ctx := context.Background()
	_, err := dbConfig.ExecContext(ctx, "DELETE FROM user_profiles")
	if err != nil {
		logs.Log("Error clearing user profiles: " + err.Error())
	} else {
		logs.Log("User profiles cleared successfully")
	}
}

func ClearEmailVerifications() {
	ctx := context.Background()
	_, err := dbConfig.ExecContext(ctx, "DELETE FROM email_verifications")
	if err != nil {
		logs.Log("Error clearing user profiles: " + err.Error())
	} else {
		logs.Log("User profiles cleared successfully")
	}
}
