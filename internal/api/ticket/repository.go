package ticket

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Selahattinn/ticketAllocating-Purchasing/internal/api/ticket/model"
	"github.com/Selahattinn/ticketAllocating-Purchasing/pkg/mysql"
)

type ITicketRepository interface {
	Create(ctx context.Context, ticket model.CreateTicketRequest) (int64, error)
	Get(ctx context.Context, ticket model.GetTicketRequest) (*model.GetTicketResource, error)
	Purchase(ctx context.Context, ticket model.PurchaseTicketRequest) error
}

type ticketRepository struct {
	mysqlDB mysql.IMysqlInstance
}

func NewTicketRepository(pg mysql.IMysqlInstance) ITicketRepository {
	return &ticketRepository{
		mysqlDB: pg,
	}
}

func (t *ticketRepository) Create(ctx context.Context, ticket model.CreateTicketRequest) (int64, error) {
	stmt, err := t.mysqlDB.Database().Prepare("insert into ticket (name,description,quantity) values (?,?,?)")
	if err != nil {
		return -1, err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, ticket.Name, ticket.Description, ticket.Allocation)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (t *ticketRepository) Get(ctx context.Context, ticket model.GetTicketRequest) (*model.GetTicketResource, error) {
	q := "SELECT id, name ,description ,quantity FROM ticket where id=?"

	res, err := t.mysqlDB.Database().QueryContext(ctx, q, ticket.ID)
	if err != nil {
		return nil, err
	}

	var ticketRes model.GetTicketResource
	for res.Next() {
		if err := res.Scan(&ticketRes.ID, &ticketRes.Name, &ticketRes.Description, &ticketRes.Allocation); err != nil {
			return nil, err
		}
	}

	if ticketRes.IsEmpty() {
		return nil, fmt.Errorf("ticket_not_found")
	}

	return &ticketRes, nil
}

func (t *ticketRepository) Purchase(ctx context.Context, ticket model.PurchaseTicketRequest) error {
	tx, err := t.mysqlDB.Database().BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ticketRes := t.getTicketBeforePurchase(tx, ticket.ID)

	if ticketRes.IsEmpty() {
		return fmt.Errorf("ticket_not_found")
	}

	if ticketRes.Allocation < ticket.Quantity {
		err = tx.Rollback()
		if err != nil {
			return err
		}

		return fmt.Errorf("not_enough_tickets")
	}

	stmt, err := tx.Prepare("UPDATE ticket SET quantity=quantity-? WHERE id=?")
	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(ticket.Quantity, ticket.ID)
	if err != nil {
		return err
	}

	_, err = t.storePurchase(tx, ticket)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (t *ticketRepository) getTicketBeforePurchase(tx *sql.Tx, id int64) *model.GetTicketResource {
	q := "SELECT id,quantity FROM ticket where id=?"

	var ticket model.GetTicketResource
	res, err := tx.Query(q, id)
	if err != nil {
		return nil
	}

	for res.Next() {
		if err := res.Scan(&ticket.ID, &ticket.Allocation); err != nil {
			return nil
		}
	}

	return &ticket
}

func (t *ticketRepository) storePurchase(tx *sql.Tx, ticket model.PurchaseTicketRequest) (string, error) {
	stmt, err := tx.Prepare("insert into purchase (user_id,quantity) values (?,?)")
	if err != nil {
		return "", err
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(context.Background(), ticket.UserID, ticket.Quantity)
	if err != nil {
		return "", err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", id), nil
}
