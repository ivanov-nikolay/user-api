package storage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/ivanov-nikolay/user-api/internal/dto"
	"github.com/ivanov-nikolay/user-api/internal/entity"
	"github.com/ivanov-nikolay/user-api/internal/filters"
)

type Storage interface {
	CreateUserStorage(user entity.User) (int64, error)
	DeleteUserStorage(ID int64) (bool, error)
	UpdateUserStorage(user entity.User) (bool, error)
	GetUserByIDStorage(ID int64) (*entity.User, error)
	SearchUsersStorage(filters filters.Filter) ([]entity.User, error)
}

type DBStorage struct {
	db         *sql.DB
	redisConn  redis.Conn
	expireTime int
}

func New(db *sql.DB, redisConn redis.Conn) *DBStorage {
	return &DBStorage{
		db:         db,
		redisConn:  redisConn,
		expireTime: 24 * 60 * 2,
	}
}

func (ps *DBStorage) CreateUserStorage(user entity.User) (int64, error) {
	var lastInsertId int64
	query := "INSERT INTO users (surname, name, gender, status, join_date"
	values := []interface{}{user.Surname, user.Name, user.Gender, user.Status, user.JoinDate}

	if user.Patronymic != "" {
		query += ", patronymic"
		values = append(values, user.Patronymic)
	}

	if !user.Birthday.IsZero() {
		query += ", birthday"
		values = append(values, user.Birthday)
	}

	query += ") VALUES ("
	for i := range values {
		query += fmt.Sprintf("$%d", i+1)
		if i < len(values)-1 {
			query += ", "
		}
	}
	query += ") RETURNING id"

	err := ps.db.QueryRow(query, values...).Scan(&lastInsertId)
	if err != nil {
		return 0, err
	}
	user.ID = lastInsertId
	go ps.saveUserToRedis(user)
	return lastInsertId, nil

}

func (ps *DBStorage) DeleteUserStorage(ID int64) (bool, error) {
	result, err := ps.db.Exec(
		"DELETE FROM users WHERE id = $1",
		ID,
	)
	if err != nil {
		return false, err
	}
	num, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if num > 0 {
		go ps.deleteUserFromRedis(ID)
		return true, nil
	}
	return false, nil
}

func (ps *DBStorage) UpdateUserStorage(user entity.User) (bool, error) {
	res, err := ps.db.Exec(
		`UPDATE users SET 
		"surname" = $1,
		"name" = $2,
		"patronymic" = $3,
		"gender" = $4,
		"status" = $5,
		"birthday" = $6
		WHERE id = $7`,
		user.Surname,
		user.Name,
		getNullOrStr(user.Patronymic),
		user.Gender,
		user.Status,
		getNullOrTime(user.Birthday),
		user.ID,
	)
	if err != nil {
		return false, err
	}
	num, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if num > 0 {
		go ps.saveUserToRedis(user)
		return true, nil
	}
	return false, nil
}

func getNullOrTime(tm time.Time) interface{} {
	if tm.IsZero() {
		return nil
	}
	return tm
}

func getNullOrStr(str string) interface{} {
	if str == "" {
		return nil
	}
	return str
}

func (ps *DBStorage) GetUserByIDStorage(ID int64) (*entity.User, error) {
	userRD, err := ps.getUserFromRedis(ID)
	if err == nil || userRD != nil {
		fmt.Println("User got from redis")
		return userRD, nil
	}
	user := &dto.UserDB{}
	err = ps.db.
		QueryRow(`SELECT id, name, surname, patronymic, gender, status, birthday, join_date FROM users WHERE id = $1`, ID).
		Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Gender, &user.Status, &user.Birthday, &user.JoinDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	convertedUser := user.ConvertToUser()
	return &convertedUser, nil
}

func (ps *DBStorage) SearchUsersStorage(filter filters.Filter) ([]entity.User, error) {
	query := "SELECT id, name, surname, patronymic, gender, status, birthday, join_date FROM users WHERE 1=1"
	var values []interface{}

	if filter.Gender != "" {
		query += " AND gender = $" + strconv.Itoa(len(values)+1)
		values = append(values, filter.Gender)
	}

	if filter.Status != "" {
		query += " AND status = $" + strconv.Itoa(len(values)+1)
		values = append(values, filter.Status)
	}

	if filter.FullName != "" {
		query += " AND (name || ' ' || surname || COALESCE(' ' || patronymic, '')) ILIKE $" + strconv.Itoa(len(values)+1)
		values = append(values, "%"+filter.FullName+"%")
	}

	if filter.AttributesToSort != "" {
		query += " ORDER BY " + filter.AttributesToSort
		if filter.SortDesc {
			query += " DESC"
		} else if filter.SortAsk {
			query += " ASC"
		}
	}

	if filter.Limit != 0 {
		query += " LIMIT $" + strconv.Itoa(len(values)+1)
		values = append(values, filter.Limit)
	}
	if filter.Offset != 0 {
		query += " OFFSET $" + strconv.Itoa(len(values)+1)
		values = append(values, filter.Offset)
	}

	rows, err := ps.db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Printf("error in closing db rows: %s", err)
		}
	}()

	var users []entity.User
	for rows.Next() {
		var user dto.UserDB
		err = rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Gender, &user.Status, &user.Birthday, &user.JoinDate)
		if err != nil {
			return nil, err
		}
		userConverted := user.ConvertToUser()
		users = append(users, userConverted)
	}

	return users, nil
}

func (ps *DBStorage) saveUserToRedis(user entity.User) {
	userJSON, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error marshalling user: %s\n", err)
		return
	}

	_, err = ps.redisConn.Do("HSET", "users", user.ID, userJSON)
	if err != nil {
		fmt.Printf("Error saving user to Redis: %s\n", err)
		return
	}

	_, err = ps.redisConn.Do("EXPIRE", "users", ps.expireTime)
	if err != nil {
		fmt.Printf("Error setting expire time for users hashset: %s\n", err)
	}

	fmt.Println("User added to Redis hashset")
}

func (ps *DBStorage) getUserFromRedis(userID int64) (*entity.User, error) {
	userJSON, err := redis.Bytes(ps.redisConn.Do("HGET", "users", userID))
	if err != nil {
		return nil, fmt.Errorf("error getting user from Redis: %s", err)
	}

	var user entity.User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling user JSON: %s", err)
	}

	return &user, nil
}

func (ps *DBStorage) deleteUserFromRedis(userID int64) error {
	_, err := ps.redisConn.Do("HDEL", "users", userID)
	if err != nil {
		return fmt.Errorf("error deleting user from Redis: %s", err)
	}

	fmt.Println("User deleted from Redis hashset")
	return nil
}
