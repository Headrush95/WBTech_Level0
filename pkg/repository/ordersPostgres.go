package repository

import (
	"WBTech_Level0/models"
	"fmt"
	"github.com/jmoiron/sqlx"
)

// TODO реализовать подготовленные выражения для ускорения роботы БД

type OrderPostgres struct {
	db *sqlx.DB
	//stmts map[string]*sqlx.Stmt
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {

	return &OrderPostgres{db: db}
}

// GetOrderById ищет заказ по его uid в БД и если находит, то возвращает структуру models.Order, если нет - ошибку.
func (op *OrderPostgres) GetOrderById(uid string) (models.Order, error) {
	// Запрашиваем из БД заказ по его uid
	var order models.Order
	query := fmt.Sprintf("SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,"+
		" delivery_service, shardkey, sm_id, date_created, oof_shard FROM %s WHERE order_uid=$1", OrdersTable)
	if err := op.db.Get(&order, query, uid); err != nil {
		return order, err
	}

	// Запрашиваем из БД информацию о доставке заказа
	var delivery models.Delivery
	var deliveryId int
	query = fmt.Sprintf("SELECT delivery FROM %s WHERE order_uid=$1", OrdersTable)
	row := op.db.QueryRow(query, uid)
	if err := row.Scan(&deliveryId); err != nil {
		return models.Order{}, err
	}
	query = fmt.Sprintf("SELECT name, phone, zip, city, address, region, email FROM %s "+
		"WHERE id=$1", DeliveryTable)
	if err := op.db.Get(&delivery, query, deliveryId); err != nil {
		return models.Order{}, nil
	}
	order.Delivery = delivery

	// Запрашиваем из БД информацию о оплате заказа
	var payment models.Payment
	query = fmt.Sprintf("SELECT transaction, request_id, currency, provider, amount, payment_dt,bank, "+
		"delivery_cost, goods_total, custom_fee FROM %s WHERE transaction=$1", TransactionsTable)
	if err := op.db.Get(&payment, query, order.Uid); err != nil {
		return models.Order{}, err
	}

	order.Payment = payment

	// Запрашиваем из БД информацию о товарах в заказе
	items := make([]models.Item, 0, 10)
	query = fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, "+
		"status FROM %s WHERE track_number=$1", ItemsTable)
	rows, err := op.db.Query(query, order.TrackNumber)
	if err != nil {
		return models.Order{}, err
	}
	defer rows.Close()
	var i models.Item
	for rows.Next() {
		err = rows.Scan(&i.ChrtId, &i.TrackNumber, &i.Price, &i.Rid, &i.Name, &i.Sale, &i.Size,
			&i.TotalPrice, &i.NmId, &i.Brand, &i.Status)
		if err != nil {
			return models.Order{}, err
		}
		items = append(items, i)
	}

	if rows.Err() != nil {
		return models.Order{}, rows.Err()
	}
	order.Items = items
	return order, nil
}

// CreateOrder создает запись о заказе в БД.
func (op *OrderPostgres) CreateOrder(o models.Order) error {
	// Создаем транзакцию
	tx, err := op.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// TODO заменить явные запросы на подготовленные statements
	// Заносим данные в таблицу delivery
	var deliveryId int
	d := o.Delivery
	query := fmt.Sprintf("INSERT INTO %s (name, phone, zip, city, address, region, email) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", DeliveryTable)
	raw := tx.QueryRow(query, d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email)
	if err = raw.Scan(&deliveryId); err != nil {
		return fmt.Errorf("error occured inserting data into \"delivery\" table: %w", err)
	}

	// Заносим данные в таблицу orders
	query = fmt.Sprintf("INSERT INTO %s (order_uid, track_number, entry, delivery, locale, "+
		"internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", OrdersTable)

	_, err = tx.Exec(query, o.Uid, o.TrackNumber, o.Entry, deliveryId, o.Locale, o.InternalSignature,
		o.CustomerId, o.DeliveryService, o.ShardKey, o.SmId, o.DateCreated, o.OofShard)
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

	// Заносим данные в таблицу transactions
	p := o.Payment
	query = fmt.Sprintf("INSERT INTO %s (transaction, request_id, currency, provider, amount, payment_dt, bank, "+
		"delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", TransactionsTable)
	_, err = tx.Exec(query, p.Transaction, p.RequestId, p.Currency, p.Provider, p.Amount, p.PaymentDate.Time, p.Bank, p.DeliveryCost,
		p.GoodsTotal, p.CustomFee)
	if err != nil {
		return fmt.Errorf("error occured inserting data into \"transactions\" table: %w", err)
	}

	return tx.Commit()
}
