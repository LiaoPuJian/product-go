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
func (p *OrderManager) Conn() error {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "order"
	}
	return nil
}

func (p *OrderManager) Insert(order *models.Order) (id int64, err error) {
	//判断mysql连接是否正常
	if err = p.Conn(); err != nil {
		return
	}
	sqlStr := fmt.Sprintf("INSERT INTO %s (user_id, product_id, order_status) VALUES (?, ?, ?)", p.table)
	stmt, err := p.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return
	}
	result, err := stmt.Exec(order.UserId, order.ProductId, order.OrderStatus)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (p *OrderManager) Delete(id int64) bool {
	//判断连接是否正常
	if err := p.Conn(); err != nil {
		return false
	}
	sqlStr := fmt.Sprintf("DELETE FROM %s WHERE id = ?", p.table)
	stmt, err := p.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return false
	}
	return true
}

func (p *OrderManager) Update(order *models.Order) error {
	//判断连接是否正常
	if err := p.Conn(); err != nil {
		return err
	}
	sqlStr := fmt.Sprintf("UPDATE %s SET product_name =?, product_num =?, product_image =?, product_url =? WHERE id = ?", p.table)
	stmt, err := p.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (p *OrderManager) SelectByKey(id int64) (*models.Order, error) {
	//判断连接是否正常
	if err := p.Conn(); err != nil {
		return &models.Product{}, err
	}
	sqlStr := fmt.Sprintf("SELECT * FROM %s WHERE ID = %s", p.table, strconv.FormatInt(id, 10))

	row, err := p.mysqlConn.Query(sqlStr)
	defer row.Close()

	if err != nil {
		return &models.Product{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &models.Product{}, nil
	}
	productResult := &models.Product{}
	common.DataToStruct(result, productResult)
	return productResult, nil
}

func (p *OrderManager) SelectAll() (result []*models.Order, err error) {
	if err = p.Conn(); err != nil {
		return
	}
	sqlStr := fmt.Sprintf("SELECT * FROM %s", p.table)
	rows, err := p.mysqlConn.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		return
	}
	res := common.GetResultRows(rows)
	if len(res) == 0 {
		return
	}
	for _, v := range res {
		product := &models.Product{}
		common.DataToStruct(v, product)
		result = append(result, product)
	}
	return
}

func (p *OrderManager) SelectAllWithInfo() (map[int]map[string]string, error) {

}
