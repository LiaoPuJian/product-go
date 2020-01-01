package repositories

import (
	"database/sql"
	"fmt"
	"product-go/common"
	"product-go/models"
	"strconv"
)

//第一步，先开发对应的接口
//第二部，实现定义的接口

type IOrder interface {
	//连接数据库
	Conn() error
	Insert(*models.Order) (int64, error)
	Delete(int64) bool
	Update(*models.Order) error
	SelectByKey(int64) (*models.Order, error)
	SelectAll() ([]*models.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderManager(table string, db *sql.DB) IOrder {
	return &OrderManager{table: table, mysqlConn: db}
}

//连接mysql
func (o *OrderManager) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderManager) Insert(order *models.Order) (id int64, err error) {
	//判断mysql连接是否正常
	if err = o.Conn(); err != nil {
		return
	}
	sqlStr := fmt.Sprintf("INSERT INTO `%s` (user_id, product_id, order_status) VALUES (?, ?, ?)", o.table)
	stmt, err := o.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return
	}
	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (o *OrderManager) Delete(id int64) bool {
	//判断连接是否正常
	if err := o.Conn(); err != nil {
		return false
	}
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE id = ?", o.table)
	stmt, err := o.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}
	return true
}

func (o *OrderManager) Update(order *models.Order) error {
	//判断连接是否正常
	if err := o.Conn(); err != nil {
		return err
	}
	sqlStr := fmt.Sprintf("UPDATE %s SET user_id =?, product_id =?, order_status =? WHERE id = ?", o.table)
	stmt, err := o.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.UserId, order.ProductId, order.OrderStatus, order.ID)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderManager) SelectByKey(id int64) (*models.Order, error) {
	//判断连接是否正常
	if err := o.Conn(); err != nil {
		return &models.Order{}, err
	}
	sqlStr := fmt.Sprintf("SELECT * FROM %s WHERE id = %s", o.table, strconv.FormatInt(id, 10))
	row, err := o.mysqlConn.Query(sqlStr)
	defer row.Close()

	if err != nil {
		return &models.Order{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &models.Order{}, nil
	}
	orderResult := &models.Order{}
	common.DataToStruct(result, orderResult)
	return orderResult, nil
}

func (o *OrderManager) SelectAll() (result []*models.Order, err error) {
	if err = o.Conn(); err != nil {
		return
	}
	sqlStr := fmt.Sprintf("SELECT * FROM %s", o.table)
	rows, err := o.mysqlConn.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		return
	}
	res := common.GetResultRows(rows)
	if len(res) == 0 {
		return
	}
	for _, v := range res {
		product := &models.Order{}
		common.DataToStruct(v, product)
		result = append(result, product)
	}
	return
}

//将对应产品的信息也取出
func (o *OrderManager) SelectAllWithInfo() (result map[int]map[string]string, err error) {
	if err = o.Conn(); err != nil {
		return nil, err
	}
	sqlStr := fmt.Sprintf("SELECT a.id,a.order_status,b.product_name FROM `%s` as a LEFT JOIN `%s` as b ON a.product_id = b.id", o.table, "product")
	fmt.Println(sqlStr)
	rows, err := o.mysqlConn.Query(sqlStr)
	//defer rows.Close()
	if err != nil {
		return nil, err
	}
	return common.GetResultRows(rows), nil
}
