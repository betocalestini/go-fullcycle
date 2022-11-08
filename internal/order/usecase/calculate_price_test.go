package usecase

import (
	"database/sql"
	"testing"

	"github.com/betocalestini/go-fullcyle/internal/order/entity"
	"github.com/betocalestini/go-fullcyle/internal/order/infra/database"
	"github.com/stretchr/testify/suite"

	_ "github.com/mattn/go-sqlite3"
)

type CalculatePriceUseCaseTestSuite struct {
	suite.Suite
	OrderRepository entity.OrderRepositoryInterface
	Db              *sql.DB
}

func (suite *CalculatePriceUseCaseTestSuite) SetupTest() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)

	//create table orders
	_, err = db.Exec("CREATE TABLE orders (id varchar(255), price float NOT NULL, tax FLOAT NOT NULL, final_price FLOAT NOT NULL, PRIMARY KEY (id))")
	suite.NoError(err)
	suite.Db = db
	suite.OrderRepository = database.NewOrderRepository(db)

}

func (suite *CalculatePriceUseCaseTestSuite) TearDownTest() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(CalculatePriceUseCaseTestSuite))
}

func (suite *CalculatePriceUseCaseTestSuite) TestCalculatePrice() {
	order, err := entity.NewOrder("123", 10, 1)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())

	calculateFinalPriceInput := OrderInpuDTO{
		ID:    order.ID,
		Price: order.Price,
		Tax:   order.Tax,
	}

	calculateFinalPriceUseCase := NewCalculateFinalPriceUseCase(suite.OrderRepository)
	output, err := calculateFinalPriceUseCase.Execute(calculateFinalPriceInput)
	suite.NoError(err)
	suite.Equal(order.FinalPrice, output.FinalPrice)
	suite.Equal(order.ID, output.ID)
	suite.Equal(order.Price, output.Price)
	suite.Equal(order.Tax, output.Tax)

}
