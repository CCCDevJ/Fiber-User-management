package model

import "fmt"

type User struct {
	ID           uint64 `json:"id"`
	Name         string `json:"name,omitempty"`
	IsActive     uint64 `json:"is_active"`
	Role         uint64 `json:"role"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	IsDelete     uint64 `json:"is_delete"`
	ApprovedById uint64 `json:"approved_by_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func GetAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT `id`, `name`, `is_active`, `role`, `email`,`approved_by_id`, `created_at`, `updated_at` FROM `tbl_user` WHERE `is_delete` = 0 ORDER BY `id` DESC;")
	if err != nil {
		return []User{}, err
	}
	// We close the resource
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.IsActive, &user.Role, &user.Email, &user.ApprovedById, &user.CreatedAt, &user.UpdatedAt)

		users = append(users, user)
	}

	return users, nil
}

func GetAllUsersForAdmin(ID string) ([]User, error) {
	rows, err := db.Query("SELECT `id`, `name`, `is_active`, `role`, `email`,`approved_by_id`, `created_at`, `updated_at`, `is_delete` FROM `tbl_user` WHERE `id` != " + ID + " ORDER BY `id` DESC;")
	if err != nil {
		return []User{}, err
	}
	// We close the resource
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		user := User{}
		rows.Scan(&user.ID, &user.Name, &user.IsActive, &user.Role, &user.Email, &user.ApprovedById, &user.CreatedAt, &user.UpdatedAt, &user.IsDelete)

		users = append(users, user)
	}

	return users, nil
}

func GetUserById(ID uint64) (User, error) {
	stmt, err := db.Prepare("SELECT `name`, `is_active`, `role`, `email`, `is_delete`, `approved_by_id`, `created_at`, `updated_at` FROM `tbl_user` WHERE `id` = ?")
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	var user User
	user.ID = ID
	err = stmt.QueryRow(ID).Scan(
		&user.Name,
		&user.IsActive,
		&user.Role,
		&user.Email,
		&user.IsDelete,
		&user.ApprovedById,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetUserByEmail(email string) (User, error) {
	stmt, err := db.Prepare("SELECT `id`, `name`, `is_active`, `role`, `password`, `is_delete`, `approved_by_id`, `created_at`, `updated_at` FROM `tbl_user` WHERE `email` = ?")
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	defer stmt.Close()

	var user User
	user.Email = email
	err = stmt.QueryRow(email).Scan(
		&user.ID,
		&user.Name,
		&user.IsActive,
		&user.Role,
		&user.Password,
		&user.IsDelete,
		&user.ApprovedById,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	return user, nil
}

func CreateUser(user User) (bool, error) {
	stmt, err := db.Prepare("INSERT INTO `tbl_user`(`name`, `email`, `password`) VALUES (?, ?, ?)")
	if err != nil {
		return false, err
	}

	defer stmt.Close()

	_, err = stmt.Exec(user.Name, user.Email, user.Password)

	// rowsAffec, _ := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return true, nil
}

func UpdateUserByID(user User) (bool, error) {
	res, err := db.Exec("UPDATE `tbl_user` SET `name`=?,`is_active`=?,`role`=?,`password`=?,`is_delete`=?,`approved_by_id`=? WHERE `id` = ?",
		user.Name,
		user.IsActive,
		user.Role,
		user.Password,
		user.IsDelete,
		user.ApprovedById,
		user.ID,
	)
	if err != nil {
		return false, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}

	return true, err
}

func UpdateUserByEmail(user User) (bool, error) {
	res, err := db.Exec("UPDATE `tbl_user` SET `name`=?,`is_active`=?,`role`=?,`password`=?,`is_delete`=?,`approved_by_id`=? WHERE `email` = ?",
		user.Name,
		user.IsActive,
		user.Role,
		user.Password,
		user.IsDelete,
		user.ApprovedById,
		user.Email,
	)
	if err != nil {
		return false, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}

	return true, err
}

func DeleteUserByID(ID uint64) (bool, error) {
	res, err := db.Exec("DELETE FROM `tbl_user` WHERE `id` = ?", ID)
	if err != nil {
		return false, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}

	return true, nil
}

func DeleteUserByEmail(Email string) (bool, error) {
	res, err := db.Exec("DELETE FROM `tbl_user` WHERE `email` = ?", Email)
	if err != nil {
		return false, err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return false, err
	}

	return true, nil
}
