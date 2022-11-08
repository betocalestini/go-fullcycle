package database

import (
	"database/sql"
	"testing"

	"github.com/betocalestini/go-fullcyle/internal/order/entity"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

type OrderRpositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRpositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	_, err = db.Exec("CREATE TABLE orders (id varchar(255), price float NOT NULL, tax FLOAT NOT NULL, final_price FLOAT NOT NULL, PRIMARY KEY (id))")
	suite.NoError(err)
	suite.Db = db

}

func (suite *OrderRpositoryTestSuite) TearDownTest() {
	// _, err := suite.Db.Exec("DELETE FROM orders")
	// suite.NoError(err)
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRpositoryTestSuite))
}

func (suite *OrderRpositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10, 1)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())

	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("SELECT id, price, tax, final_price FROM orders WHERE id = $1", order.ID).Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)

}
