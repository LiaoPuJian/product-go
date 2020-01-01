package repositories

import (
	"database/sql"
	"fmt"
	"product-go/common"
	"product-go/models"
	"strconv"

	"github.com/pkg/errors"
)

type UserRepository interface {
	Conn() error
	Select(userName string) (user *models.User, err error)
	Insert(user *models.User) (userId int64, err error)
	SelectByID(id int64) (user *models.User, err error)
}

type UserManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewUserRepository(table string, db *sql.DB) UserRepository {
	return &UserManagerRepository{
		table:     table,
		mysqlConn: db,
	}
}

func (u *UserManagerRepository) Conn() error {
	if u.mysqlConn == nil {
		newConn, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		u.mysqlConn = newConn
	}
	if u.table == "" {
		u.table = "user"
	}
	return nil
}

func (u *UserManagerRepository) Select(userName string) (user *models.User, err error) {
	if userName == "" {
		return user, errors.New("查询条件不能为空")
	}
	if err = u.Conn(); err != nil {
		return user, err
	}
	sqlStr := fmt.Sprintf("SELECT * FROM `%s` WHERE user_name = '%s'", u.table, userName)
	rows, err := u.mysqlConn.Query(sqlStr)
	//defer rows.Close()
	if err != nil {
		return user, err
	}

	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return user, errors.New("用户不存在!")
	}
	user = &models.User{}

	common.DataToStruct(result, user)
	return user, nil
}

func (u *UserManagerRepository) Insert(user *models.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return userId, err
	}

	sqlStr := fmt.Sprintf("INSERT INTO `%s` (nick_name, user_name, password) VALUES (?, ?, ?)", u.table)
	stmt, err := u.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return userId, err
	}
	//这里的password需经过加密
	result, err := stmt.Exec(user.NickName, user.UserName, user.Password)
	if err != nil {
		return userId, err
	}
	return result.LastInsertId()
}

func (u *UserManagerRepository) SelectByID(id int64) (user *models.User, err error) {
	if err = u.Conn(); err != nil {
		return user, err
	}
	sqlStr := fmt.Sprintf("SELECT * FROM `%s` WHERE id = %s", u.table, strconv.Itoa(int(id)))
	rows, err := u.mysqlConn.Query(sqlStr)
	//defer rows.Close()
	if err != nil {
		return user, err
	}
	result := common.GetResultRow(rows)
	if len(result) == 0 {
		return user, errors.New("用户不存在")
	}
	common.DataToStruct(result, user)
	return user, nil
}
