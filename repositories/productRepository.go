package repositories

import (
	"database/sql"
	"product-go/common"
	"product-go/models"
	"strconv"
)

//第一步，先开发对应的接口
//第二部，实现定义的接口

type IProduct interface {
	//连接数据库
	Conn() error
	Insert(*models.Product) (int64, error)
	Delete(int64) bool
	Update(*models.Product) error
	SelectByKey(int64) (*models.Product, error)
	SelectAll() ([]*models.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

//连接mysql
func (p *ProductManager) Conn() error {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return nil
}

func (p *ProductManager) Insert(product *models.Product) (id int64, err error) {
	//判断mysql连接是否正常
	if err = p.Conn(); err == nil {
		return
	}
	sqlStr := "INSERT INTO product (product_name, product_num, product_image, product_url) VALUES(?, ?, ?, ?)"
	stmt, err := p.mysqlConn.Prepare(sqlStr)
	if err != nil {
		return
	}
	result, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return
	}
	return result.LastInsertId()
}

func (p *ProductManager) Delete(id int64) bool {
	//判断连接是否正常
	if err := p.Conn(); err != nil {
		return false
	}
	sqlStr := "DELETE FROM product WHERE id = ?"
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

func (p *ProductManager) Update(product *models.Product) error {
	//判断连接是否正常
	if err := p.Conn(); err != nil {
		return err
	}
	sqlStr := "UPDATE product SET product_name =?, product_num =?, product_image =?, product_url =? WHERE id = ?"
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

func (p *ProductManager) SelectByKey(id int64) (*models.Product, error) {
	//判断连接是否正常
	if err := p.Conn(); err != nil {
		return &models.Product{}, err
	}
	sqlStr := "SELECT * FROM " + p.table + "WHERE ID = " + strconv.FormatInt(id, 10)
	row, err := p.mysqlConn.Query(sqlStr)
	if err != nil {
		return &models.Product{}, err
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &models.Product{}, nil
	}
	productResult := &models.Product{}
	common.DataToStructByTagSql(result, productResult)
	return productResult, nil
}

func (p *ProductManager) SelectAll() ([]*models.Product, error) {
	panic("implement me")
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}
