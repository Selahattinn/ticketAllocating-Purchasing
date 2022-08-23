package model_test

import (
	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TicketModel", func() {
	Describe("IsEmpty", func() {
		It("should return true", func() {
			ticket := model.GetTicketResource{
				ID: int64(0),
			}

			Expect(ticket.IsEmpty()).To(BeTrue())
		})
		It("should return false", func() {
			ticket := model.GetTicketResource{
				ID: int64(1),
			}

			Expect(ticket.IsEmpty()).To(BeFalse())
		})
	})
})
