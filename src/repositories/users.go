package repositories

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

type users struct {
	db *sql.DB
}

func NewRepositoryUsers(db *sql.DB) *users {
	return &users{db}
}

func (repository users) Create(user models.User) (id uint64, err error) {
	statement, err := repository.db.Prepare("INSERT INTO users (name, nick, email, password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return
	}

	defer statement.Close()

	result, err := statement.Exec(user.Name, user.Nick, user.Email, user.Password)
	if err != nil {
		return
	}

	LastId, err := result.LastInsertId()
	if err != nil {
		return
	}

	id = uint64(LastId)

	return
}

func (repository users) List(nameOrNick string) (users []models.User, err error) {
	nameOrNick = fmt.Sprintf("%%%s%%", nameOrNick) //%nameOrNick%

	rows, err := repository.db.Query("SELECT id, name, nick, email, created_at FROM users WHERE name LIKE ? OR nick LIKE ?", nameOrNick, nameOrNick)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err = rows.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt); err != nil {
			return
		}
		users = append(users, user)
	}
	return
}

func (repository users) FindUser(ID uint64) (user models.User, err error) {
	rows, err := repository.db.Query("SELECT id, name, nick, email, created_at FROM users WHERE id = ?", ID)
	if err != nil {
		return
	}

	defer rows.Close()
	if rows.Next() {
		if err = rows.Scan(&user.ID, &user.Name, &user.Nick, &user.Email, &user.CreatedAt); err != nil {
			return
		}
	}

	return
}

func (repository users) DeleteUser(ID uint64) (err error) {
	statement, err := repository.db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return
	}

	defer statement.Close()

	_, err = statement.Exec(ID)
	if err != nil {
		return
	}
	return
}

func (repository users) UpdateUser(ID uint64, user models.User) (err error) {
	statement, err := repository.db.Prepare("UPDATE users SET name = ?, nick = ?, email = ? WHERE id = ?")
	if err != nil {
		return
	}

	defer statement.Close()

	_, err = statement.Exec(user.Name, user.Nick, user.Email, ID)
	if err != nil {
		return
	}

	return
}

func (repository users) FindUserByEmail(email string) (user models.User, err error) {
	row, err := repository.db.Query("SELECT id, password FROM users WHERE email = ?", email)
	if err != nil {
		return
	}

	defer row.Close()

	if row.Next() {
		if err = row.Scan(&user.ID, &user.Password); err != nil {
			return
		}
	}

	return
}

func (repository users) Flollow(userID uint64, followerID uint64) (err error) {
	statement, err := repository.db.Prepare(
		"insert ignore into followers (user_id, follower_id) values (?, ?)",
	)
	if err != nil {
		return
	}

	defer statement.Close()

	if _, err = statement.Exec(userID, followerID); err != nil {
		return
	}

	return
}

func (repository users) Unflollow(userID uint64, followerID uint64) (err error) {
	statement, err := repository.db.Prepare(
		"delete from followers where user_id = ? and follower_id = ?",
	)
	if err != nil {
		return
	}

	defer statement.Close()

	if _, err = statement.Exec(userID, followerID); err != nil {
		return
	}

	return
}

func (repository users) ListFollowers(userID uint64) (followers []models.User, err error) {
	rows, err := repository.db.Query(`SELECT
										u.id,
										u.name,
										u.nick,
										u.email,
										u.created_at
									FROM
										users u
									join followers f on
										(u.id = f.follower_id)
									WHERE
										f.user_id = ?`, userID)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var follower models.User
		if err = rows.Scan(&follower.ID, &follower.Name, &follower.Nick, &follower.Email, &follower.CreatedAt); err != nil {
			return
		}
		followers = append(followers, follower)
	}
	return
}
