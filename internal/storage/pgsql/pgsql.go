package pgsql

import (
	"database/sql"
	"log"
	"time"

	"wbTechL0/models"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) SaveOrder(order models.Order) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.Exec(`
	INSERT INTO orders (
		order_id, order_uid, track_number, entry, delivery_name, delivery_phone, delivery_zip,
		delivery_city, delivery_address, delivery_region, delivery_email, payment_transaction,
		payment_request_id, payment_currency, payment_provider, payment_amount, payment_dt,
		payment_bank, payment_delivery_cost, payment_goods_total, payment_custom_fee, locale,
		internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES (
		DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
		$21, $22, $23, $24, $25, $26, $27, $28
	) RETURNING order_id`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region,
		order.Delivery.Email, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, time.Unix(order.Payment.PaymentDT, 0),
		order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
		order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey,
		order.SMID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		log.Println("ошибка в orders, такой order уже существует")
		tx.Rollback()
		return err
	}

	var orderID int
	err = tx.QueryRow("SELECT LASTVAL()").Scan(&orderID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
			INSERT INTO order_items (
				order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
			)`,
			orderID, item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name, item.Sale,
			item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			log.Println("ошибка в items")
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *Database) GetOrder(orderID int) (*models.Order, error) {
	row := d.db.QueryRow(`
		SELECT *
		FROM orders
		WHERE order_id = $1`,
		orderID,
	)

	order := &models.Order{}
	delivery := &models.Delivery{}
	payment := &models.Payment{}

	err := row.Scan(
		&order.OrderID, &order.OrderUID, &order.TrackNumber, &order.Entry,
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
		&delivery.Region, &delivery.Email, &payment.Transaction, &payment.RequestID, &payment.Currency,
		&payment.Provider, &payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost,
		&payment.GoodsTotal, &payment.CustomFee, &order.Locale, &order.InternalSignature,
		&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SMID, &order.DateCreated,
		&order.OofShard,
	)
	if err != nil {
		return nil, err
	}

	order.Delivery = *delivery
	order.Payment = *payment

	rows, err := d.db.Query(`
		SELECT *
		FROM order_items
		WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.OrderItem{}
		err := rows.Scan(
			&item.ItemID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, *item)
	}

	return order, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) MakeMigrations() {
	driver, err := postgres.WithInstance(d.db, &postgres.Config{})
	if err != nil {
		log.Fatal("Error creating database driver instance:", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("Error creating migration instance:", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migrations:", err)
	}

	log.Println("Migrations applied successfully!")
}

func (d *Database) GetAllOrders() ([]models.Order, error) {
	rows, err := d.db.Query(`
		SELECT * FROM orders`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		order := models.Order{}
		delivery := models.Delivery{}
		payment := models.Payment{}
		var paymentDTUnix time.Time
		err := rows.Scan(
			&order.OrderID, &order.OrderUID, &order.TrackNumber, &order.Entry,
			&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
			&delivery.Region, &delivery.Email, &payment.Transaction, &payment.RequestID, &payment.Currency,
			&payment.Provider, &payment.Amount, &paymentDTUnix, &payment.Bank, &payment.DeliveryCost,
			&payment.GoodsTotal, &payment.CustomFee, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.ShardKey, &order.SMID, &order.DateCreated,
			&order.OofShard,
		)
		if err != nil {
			log.Println("ошибка в orders")
			return nil, err
		}
		payment.PaymentDT = paymentDTUnix.Unix()
		order.Delivery = delivery
		order.Payment = payment

		order.Items, err = d.getOrderItems(order.OrderID)
		if err != nil {
			log.Println("ошибка в order_items")
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (d *Database) getOrderItems(orderID int) ([]models.OrderItem, error) {
	rows, err := d.db.Query(`
		SELECT * FROM order_items
		WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem

	for rows.Next() {
		item := models.OrderItem{}
		err := rows.Scan(
			&item.ItemID, &item.OrderID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.RID,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}
