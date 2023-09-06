package repository

import (
	"WBTech_Level0/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (op *OrderPostgres) GetOrderById(uid string) (models.Order, error) {
	return models.Order{}, nil
}

func (op *OrderPostgres) CreateOrder(o models.Order) error { //нужно ли возвращать id?
	// Создаем транзакцию
	tx, err := op.db.Begin()
	if err != nil {
		return err
	}

	// Заносим данные в таблицу delivery
	var deliveryId int
	d := o.Delivery
	query := fmt.Sprintf("INSERT INTO %s (name, phone, zip, city, address, region, email) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", DeliveryTable)
	raw := tx.QueryRow(query, d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email)
	if err = raw.Scan(&deliveryId); err != nil {
		return fmt.Errorf("error occured inserting data into \"delivery\" table: %w", err)
	}

	// Заносим данные в таблицу payment
	var paymentId int
	p := o.Payment
	query = fmt.Sprintf("INSERT INTO %s (transaction, request_id, currency, provider, amount, payment_dt, bank, "+
		"delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id", TransactionsTable)
	raw = tx.QueryRow(query, p.Transaction, p.RequestId, p.Currency, p.Provider, p.Amount, p.PaymentDate.Time, p.Bank, p.DeliveryCost,
		p.GoodsTotal, p.CustomFee)
	if err = raw.Scan(&paymentId); err != nil {
		return fmt.Errorf("error occured inserting data into \"payment\" table: %w", err)
	}

	// Заносим данные в таблицу orders
	query = fmt.Sprintf("INSERT INTO %s (order_uid, track_number, entry, delivery, payment, locale, "+
		"internal_signature, shardkey, sm_id, date_created, oof_shard) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", OrdersTable) //исправить sharedkey

	_, err = tx.Exec(query, o.Uid, o.TrackNumber, o.Entry, deliveryId, paymentId, o.Locale, o.InternalSignature,
		o.ShardKey, o.SmId, o.DateCreated, o.OofShard)
	if err != nil {
		return fmt.Errorf("error occured inserting data into \"orders\" table: %w", err)
	}

	// Заносим данные в таблицу items
	for _, im := range o.Items {
		query = fmt.Sprintf("INSERT INTO %s (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, "+
			"brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", ItemsTable)
		_, err = tx.Exec(query, im.ChrtId, im.TrackNumber, im.Price, im.Rid, im.Name, im.Sale, im.Size, im.TotalPrice, im.NmId,
			im.Brand, im.Status)
		if err != nil {
			return fmt.Errorf("error occured inserting data into \"items\" table: %w", err)
		}
	}

	return tx.Commit()
}
