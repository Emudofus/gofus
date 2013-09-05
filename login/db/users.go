package db

import (
	"database/sql"
	"fmt"
	"github.com/Blackrush/gofus/shared/db"
	"time"
)

type User struct {
	Id                   uint
	Name                 string
	Password             string
	Nickname             string
	SecretQuestion       string
	SecretAnswer         string
	Rights               db.UserRight
	CommunityId          int
	SubscriptionEnd      time.Time
	CurrentRealmServerId int
}

func (user *User) String() string {
	return fmt.Sprintf("User{ Id: %d, Name: %s, Nickname: %s, Rights: %d, CommunityId: %d, SubscriptionEnd: %s, CurrentRealmServerId: %d }", user.Id, user.Name, user.Nickname, user.Rights, user.CommunityId, user.SubscriptionEnd, user.CurrentRealmServerId)
}

func (user *User) ValidPassword(password string) bool {
	return user.Password == password
}

func (user *User) IsConnected() bool {
	return user.CurrentRealmServerId != -1
}

type Users struct{ *sql.DB }

func scan_user(rows *sql.Rows, user *User) error {
	return rows.Scan(&user.Id, &user.Name, &user.Password, &user.Nickname, &user.SecretQuestion, &user.SecretAnswer, &user.Rights, &user.CommunityId, &user.SubscriptionEnd, &user.CurrentRealmServerId)
}

func (db *Users) find(query string, args ...interface{}) (user *User, err error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	user = &User{}
	if err = scan_user(rows, user); err != nil {
		user = nil
	}
	return
}

func (db *Users) find_many(query string, args ...interface{}) (users []*User, err error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return
	}

	rows, err := stmt.Query(args...)
	if err != nil {
		return
	}

	for rows.Next() {
		user := &User{}
		if err = scan_user(rows, user); err != nil {
			return
		}

		users = append(users, user)
	}

	return
}

func user_values(user *User, with_id, id_last bool) (values []interface{}) {
	if with_id && !id_last {
		values = append(values, user.Id)
	}
	values = append(values, user.Name, user.Password, user.Nickname, user.SecretQuestion, user.SecretAnswer, user.Rights, user.CommunityId, user.SubscriptionEnd, user.CurrentRealmServerId)
	if with_id && id_last {
		values = append(values, user.Id)
	}
	return
}

func (db *Users) FindById(id int) (*User, error) {
	return db.find("select id, name, password, nickname, secret_question, secret_answer, rights, community_id, subscription_end, current_realm_server from users where id=$1", id)
}

func (db *Users) FindByName(name string) (*User, error) {
	return db.find("select id, name, password, nickname, secret_question, secret_answer, rights, community_id, subscription_end, current_realm_server from users where name=$1", name)
}

func (db *Users) FindAll() ([]*User, error) {
	return db.find_many("select id, name, password, nickname, secret_question, secret_answer, rights, community_id, subscription_end, current_realm_server from users")
}

func (db *Users) Insert(user *User) error {
	stmt, err := db.Prepare("insert into users(name, password, nickname, secret_question, secret_answer, rights, community_id, subscription_end, current_realm_server) values($1, $2, $3, $4, $5, $6, $7, $8, $9)")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(user_values(user, false, false)...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.Id = uint(id)
	return nil
}

func (db *Users) Update(user *User) error {
	stmt, err := db.Prepare("update users set name=$1, password=$2, nickname=$3, secret_question=$4, secret_answer=$5, rights=$6, community_id=$7, subscription_end=$8, current_realm_server=$9 where id=$10")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user_values(user, true, true)...)
	if err != nil {
		return err
	}

	return nil
}

func (db *Users) Delete(user *User) error {
	stmt, err := db.Prepare("delete from users where id=$1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (db *Users) ResetCurrentRealm() (err error) {
	_, err = db.Exec("update users set current_realm_server='-1'")
	return
}

func (db *Users) ResetCurrentRealmFor(id int) (err error) {
	_, err = db.Exec("update users set current_realm_server='-1' where current_realm_server=$1", id)
	return
}
