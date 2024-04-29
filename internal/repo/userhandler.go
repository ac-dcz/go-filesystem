package repo

import (
	"fmt"
	"go-fs/internal/repo/user"
	"log"
)

func InsertUser(user *user.UserInfo) error {
	sql := fmt.Sprintf(
		"insert into %s (user_name,user_pwd,email,tel,signup_ts,last_active_ts,status) values (%s)",
		userinfoSession.TableName(),
		bindVars(7),
	)
	vars := []any{user.Name, user.PassWD, user.Email, user.Tel, user.SignUpTs, user.LastActiveTs, user.Status}

	_, err := userinfoSession.Raw(sql, vars...).Exex()
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateUser(user *user.UserInfo) error {
	sql := fmt.Sprintf(
		"update %s set user_pwd = ?,email = ?,tel = ?,last_active_ts =?,email_checked = ?, tel_checked = ?,profile = ? ,status = ? where user_name = ?",
		userinfoSession.TableName(),
	)
	vars := []any{user.PassWD, user.Email, user.Tel, user.LastActiveTs, user.EmailChecked, user.TelChecked, user.Profile, user.Status, user.Name}

	_, err := userinfoSession.Raw(sql, vars...).Exex()
	if err != nil {
		log.Println(err)
	}
	return err
}

func SelectUser(name string) (*user.UserInfo, error) {
	sql := fmt.Sprintf(
		"select user_name,user_pwd,email,tel,email_checked,tel_checked,signup_ts,last_active_ts,profile,status from %s where user_name = ?",
		userinfoSession.TableName(),
	)

	row := userinfoSession.Raw(sql, name).QueryRow()
	if err := row.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	info := &user.UserInfo{}
	err := row.Scan(
		&info.Name, &info.PassWD, &info.Email, &info.Tel, &info.EmailChecked,
		&info.TelChecked, &info.SignUpTs, &info.LastActiveTs, &info.Profile, &info.Status,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return info, nil
}
