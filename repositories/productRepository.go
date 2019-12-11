package repositories

import (
	"database/sql"
	"product-go/models"
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

func (p *ProductManager) Conn() error {

	panic("implement me")
}

func (p *ProductManager) Insert(*models.Product) (int64, error) {
	panic("implement me")
}

func (p *ProductManager) Delete(int64) bool {
	panic("implement me")
}

func (p *ProductManager) Update(*models.Product) error {
	panic("implement me")
}

func (p *ProductManager) SelectByKey(int64) (*models.Product, error) {
	panic("implement me")
}

func (p *ProductManager) SelectAll() ([]*models.Product, error) {
	panic("implement me")
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}
