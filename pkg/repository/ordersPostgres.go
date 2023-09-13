package repository

import (
	"WBTech_Level0/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

var EmptyDB = errors.New("empty DB")

type OrderPostgres struct {
	db        *sqlx.DB
	stmts     map[string]*sql.Stmt
	validator *validator.Validate
}

func statementError(err error) {
	if err != nil {
		logrus.Fatalf("error occurred while preparing statement: %s", err.Error())
	}
}

// newStatements подготавливает выражения для запросов в БД.
func newStatements(db *sqlx.DB) map[string]*sql.Stmt {
	var err error
	stmts := make(map[string]*sql.Stmt, 10)

	stmts["getAllOrdersUid"], err = db.Prepare(fmt.Sprintf(`SELECT order_uid FROM %s`, OrdersTable))
	statementError(err)

	stmts["getOrderById"], err = db.Prepare(fmt.Sprintf("SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,"+
		" delivery_service, shardkey, sm_id, date_created, oof_shard FROM %s WHERE order_uid=$1", OrdersTable))
	statementError(err)

	stmts["getDeliveryIdFromOrders"], err = db.Prepare(fmt.Sprintf("SELECT delivery FROM %s WHERE order_uid=$1", OrdersTable))
	statementError(err)

	stmts["getDeliveryById"], err = db.Prepare(fmt.Sprintf("SELECT name, phone, zip, city, address, region, email FROM %s "+
		"WHERE id=$1", DeliveryTable))
	statementError(err)

	stmts["getPaymentByOrderId"], err = db.Prepare(fmt.Sprintf("SELECT transaction, request_id, currency, provider, amount, payment_dt,bank, "+
		"delivery_cost, goods_total, custom_fee FROM %s WHERE transaction=$1", TransactionsTable))
	statementError(err)

	stmts["getItemByTrack"], err = db.Prepare(fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, "+
		"status FROM %s WHERE track_number=$1", ItemsTable))
	statementError(err)

	stmts["createDelivery"], err = db.Prepare(fmt.Sprintf("INSERT INTO %s (name, phone, zip, city, address, region, email) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id", DeliveryTable))
	statementError(err)

	stmts["createOrder"], err = db.Prepare(fmt.Sprintf("INSERT INTO %s (order_uid, track_number, entry, delivery, locale, "+
		"internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)", OrdersTable))
	statementError(err)

	stmts["createItem"], err = db.Prepare(fmt.Sprintf("INSERT INTO %s (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, "+
		"brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", ItemsTable))
	statementError(err)

	stmts["createTransaction"], err = db.Prepare(fmt.Sprintf("INSERT INTO %s (transaction, request_id, currency, provider, amount, payment_dt, bank, "+
		"delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", TransactionsTable))
	statementError(err)

	return stmts
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db, stmts: newStatements(db), validator: validator.New(validator.WithRequiredStructEnabled())}
}

// getAllOrders возвращает все заказы из БД (нужен для восстановления кэша после перезапуска приложения)
func (op *OrderPostgres) getAllOrders() (map[string]models.Order, error) {
	orders := make(map[string]models.Order, 100)
	var order models.Order
	rows, err := op.stmts["getAllOrdersUid"].Query()

	if err != nil {
		return nil, err
	}

	var orderId string
	for rows.Next() {
		err = rows.Scan(&orderId)
		if err != nil {
			return nil, err
		}
		order, err = op.GetOrderById(orderId)
		if err != nil {
			return nil, err
		}
		orders[orderId] = order
	}
	if len(orders) == 0 {
		return nil, EmptyDB
	}
	return orders, nil
}

// GetOrderById ищет заказ по его uid в БД.
func (op *OrderPostgres) GetOrderById(uid string) (models.Order, error) {
	// Запрашиваем из БД заказ по его uid
	var o models.Order

	row := op.stmts["getOrderById"].QueryRow(uid)
	if err := row.Scan(&o.Uid, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature, &o.CustomerId,
		&o.DeliveryService, &o.ShardKey, &o.SmId, &o.DateCreated, &o.OofShard); err != nil {
		return o, err
	}

	// Запрашиваем из БД информацию о доставке заказа
	var delivery models.Delivery
	var deliveryId int

	row = op.stmts["getDeliveryIdFromOrders"].QueryRow(uid)
	if err := row.Scan(&deliveryId); err != nil {
		return models.Order{}, err
	}

	row = op.stmts["getDeliveryById"].QueryRow(deliveryId)
	if err := row.Scan(&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
		&delivery.Region, &delivery.Email); err != nil {
		return models.Order{}, err
	}

	o.Delivery = delivery

	// Запрашиваем из БД информацию о оплате заказа
	var payment models.Payment

	row = op.stmts["getPaymentByOrderId"].QueryRow(o.Uid)
	if err := row.Scan(&payment.Transaction, &payment.RequestId, &payment.Currency, &payment.Provider, &payment.Amount, &payment.PaymentDate,
		&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee); err != nil {
		return models.Order{}, err
	}

	o.Payment = payment

	// Запрашиваем из БД информацию о товарах в заказе
	items := make([]models.Item, 0, 10)

	rows, err := op.stmts["getItemByTrack"].Query(o.TrackNumber)
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
	o.Items = items
	return o, nil
}

// CreateOrder создает запись о заказе в БД.
func (op *OrderPostgres) CreateOrder(o models.Order) error {
	if err := op.validator.Struct(o); err != nil {
		return err
	}

	// Создаем транзакцию
	tx, err := op.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Заносим данные в таблицу delivery
	var deliveryId int
	d := o.Delivery

	txStmt := tx.Stmt(op.stmts["createDelivery"])
	raw := txStmt.QueryRow(d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email)
	if err = raw.Scan(&deliveryId); err != nil {
		return fmt.Errorf("error occured inserting data into \"delivery\" table: %w", err)
	}

	// Заносим данные в таблицу orders
	txStmt = tx.Stmt(op.stmts["createOrder"])
	_, err = txStmt.Exec(o.Uid, o.TrackNumber, o.Entry, deliveryId, o.Locale, o.InternalSignature,
		o.CustomerId, o.DeliveryService, o.ShardKey, o.SmId, o.DateCreated, o.OofShard)
	if err != nil {
		return fmt.Errorf("error occured inserting data into \"orders\" table: %w", err)
	}

	// Заносим данные в таблицу items
	for _, im := range o.Items {
		txStmt = tx.Stmt(op.stmts["createItem"])
		_, err = txStmt.Exec(im.ChrtId, im.TrackNumber, im.Price, im.Rid, im.Name, im.Sale, im.Size, im.TotalPrice, im.NmId,
			im.Brand, im.Status)
		if err != nil {
			return fmt.Errorf("error occured inserting data into \"items\" table: %w", err)
		}
	}

	// Заносим данные в таблицу transactions
	p := o.Payment

	txStmt = tx.Stmt(op.stmts["createTransaction"])
	_, err = txStmt.Exec(p.Transaction, p.RequestId, p.Currency, p.Provider, p.Amount, p.PaymentDate.Time, p.Bank, p.DeliveryCost,
		p.GoodsTotal, p.CustomFee)
	if err != nil {
		return fmt.Errorf("error occured inserting data into \"transactions\" table: %w", err)
	}

	return tx.Commit()
}

// CloseStatements закрывает все подготовленные выражения
func (op *OrderPostgres) CloseStatements() error {
	var err error
	for _, stmt := range op.stmts {
		err = stmt.Close()
		if err != nil {
			return err
		}
	}
	return err
}
