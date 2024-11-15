package postgresql

const (
	InsertUserStmt = `
		INSERT INTO my_user(username, password, email, first_name, last_name) 
		VALUES(@username, @password, @email, @firstName, @lastName)
		RETURNING id;
	`
	GetUserByName = `
		SELECT id FROM my_user WHERE username=@username;
	`
	GetUserCredStmt = `
		SELECT id, username, password, role FROM my_user WHERE username=@username;
	`
	UpdateUserRoleStmt = `
		UPDATE my_user SET role=@role WHERE username=@username;
	`
)