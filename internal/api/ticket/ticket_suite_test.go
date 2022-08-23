package ticket_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/constants"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/mysql"
)

var (
	mysqlInstance mysql.IMysqlInstance
)

func TestRent(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ticket Suite")
}

var _ = BeforeSuite(func() {
	var err error

	mysqlInstance, err = mysql.InitMysql(mysql.Config{
		URL: constants.MysqlTestURL,
	})

	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	defer func() {
		_ = mysqlInstance.Database().Close()
	}()
})
