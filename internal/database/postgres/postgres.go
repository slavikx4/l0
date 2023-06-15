package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/slavikx4/l0/internal/models"
	//"github.com/slavikx4/l0/pkg/error"
	"github.com/slavikx4/l0/pkg/logger"
)

type Postgres struct {
	*pgxpool.Pool
}

func NewPostgres(url string) (*Postgres, error) {
	DB, err := pgxpool.New(context.Background(), url)
	if err != nil {
		logger.Logger.Error.Println("не удалось подключиться к DataBase L0: ", err)
		return nil, err
	}
	if err := DB.Ping(context.Background()); err != nil {
		logger.Logger.Error.Println("не удалось пингануть к DataBase L0: ", err)
		return nil, err
	}
	logger.Logger.Process.Println("подключён успешно postgres")

	postgres := Postgres{DB}
	return &postgres, nil
}

func (p *Postgres) SetOrder(order *models.Order) error {
	const op = "Postgres.SetOrder"

	if err := p.setDelivery(&order.Delivery); err != nil {
		return err
	}

	if err := p.setPayment(&order.Payment); err != nil {
		return err
	}

	if err := p.setItems(&order.Items, order.OrderUID); err != nil {
		return err
	}

	query := `INSERT INTO "order" (
                     "order_uid",
                     "track_number",
                     "entry",
                     "delivery", 
                     "payment",
                     "locale",
                     "internal_signature",
                     "customer_id",
                     "delivery_service",
                     "shardkey",
                     "sm_id",
                     "date_created",
                     "oof_shard") VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`

	if _, err := p.Pool.Exec(context.TODO(), query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Delivery.Phone,
		order.Payment.Transaction,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) GetOrders() (*[]models.Order, error) {
	orders := make([]models.Order, 0)

	query := `SELECT 
    				"order_uid",
    				"track_number",
    				"entry",
    				"delivery",
    				"payment",
    				"locale",
    				"internal_signature",
    				"customer_id",
    				"delivery_service",
    				"shardkey",
    				"sm_id",
    				"date_created",
    				"oof_shard"
    		  FROM "order"`

	rows, err := p.Pool.Query(context.TODO(), query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		order := models.Order{
			Delivery: models.Delivery{},
			Payment:  models.Payment{},
			Items:    models.Items{},
		}
		if err := rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Delivery.Phone,
			&order.Payment.Transaction,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard); err != nil {
			return nil, err
		}
		if err := p.getDelivery(&order.Delivery); err != nil {
			return nil, err
		}

		if err := p.getPayment(&order.Payment); err != nil {
			return nil, err
		}
		if err := p.getItems(&order.Items, order.OrderUID); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return &orders, nil
}

func (p *Postgres) setDelivery(delivery *models.Delivery) error {
	query := `INSERT INTO "delivery"(
                       "phone",
                       "name",
                       "zip",
                       "city",
                       "address",
                       "region",
                       "email") VALUES ($1,$2,$3,$4,$5,$6,$7)`

	if _, err := p.Pool.Exec(context.TODO(), query,
		delivery.Phone,
		delivery.Name,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) setPayment(payment *models.Payment) error {
	query := `INSERT INTO "payment"(
                      "transaction",
                      "request_id",
                      "currency",
                      "provider",
                      "amount",
                      "payment_dt",
                      "bank",
                      "delivery_cost",
                      "goods_total",
                      "custom_fee") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	if _, err := p.Pool.Exec(context.TODO(), query,
		payment.Transaction,
		payment.RequestID,
		payment.Currency,
		payment.Provider,
		payment.Amount,
		payment.PaymentDt,
		payment.Bank,
		payment.DeliveryCost,
		payment.GoodsTotal,
		payment.CustomFee); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) setItems(items *models.Items, orderUID string) error {
	for _, item := range *items {
		if err := p.setItem(&item); err != nil {
			return err
		}
		if err := p.setOrderOfItem(orderUID, item.ChrtID); err != nil {
			return err
		}
	}
	return nil
}

func (p *Postgres) setItem(item *models.Item) error {
	query := `INSERT INTO "items"(
                    "chrt_id",
                    "track_number",
                    "price",
                    "rid",
                    "name",
                    "sale",
                    "size",
                    "total_price",
                    "nm_id",
                    "brand",
                    "status") VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9,$10,$11)`

	if _, err := p.Pool.Exec(context.TODO(), query,
		item.ChrtID,
		item.TrackNumber,
		item.Price,
		item.Rid,
		item.Name,
		item.Sale,
		item.Size,
		item.TotalPrice,
		item.NmID,
		item.Brand,
		item.Status); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) setOrderOfItem(orderUID string, chrtID int) error {
	query := `INSERT INTO "order_item"(
                         "order_uid",
                         "chrt_id") VALUES ($1, $2)`
	if _, err := p.Pool.Exec(context.TODO(), query, orderUID, chrtID); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) getDelivery(delivery *models.Delivery) error {
	query := `SELECT 
				"phone",
				"name",
				"zip",
				"city",
				"address",
				"region",
				"email"
			FROM "delivery"`

	if err := p.Pool.QueryRow(context.TODO(), query).Scan(
		&delivery.Phone,
		&delivery.Name,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) getPayment(payment *models.Payment) error {
	query := `SELECT 
					"transaction",
					"request_id",
					"currency",
					"provider",
					"amount",
					"payment_dt",
					"bank",
					"delivery_cost",
					"goods_total",
					"custom_fee"
				FROM "payment"`

	if err := p.Pool.QueryRow(context.TODO(), query).Scan(
		&payment.Transaction,
		&payment.RequestID,
		&payment.Currency,
		&payment.Provider,
		&payment.Amount,
		&payment.PaymentDt,
		&payment.Bank,
		&payment.DeliveryCost,
		&payment.GoodsTotal,
		&payment.CustomFee); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) getItems(items *models.Items, orderUID string) error {
	query := `SELECT
					"track_number",
					"price",
					"rid",
					"name",
					"sale",
					"size",
					"total_price",
					"nm_id",
					"brand",
					"status"
			FROM "items" WHERE ("chrt_id"=$1)`

	chrtIDs, err := p.getItemChrtIDsWithOrderUID(orderUID)
	if err != nil {
		return err
	}

	for _, chrtID := range chrtIDs {
		item := models.Item{}

		if err := p.Pool.QueryRow(context.TODO(), query, chrtID).Scan(
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status); err != nil {
			return err
		}

		*items = append(*items, item)
	}

	return nil
}

func (p *Postgres) getItemChrtIDsWithOrderUID(orderUID string) ([]string, error) {
	query := `SELECT
					"chrt_id"
				FROM "order_item" WHERE ("order_uid"=$1)`

	chrtIDs := make([]string, 0, 1)

	rows, err := p.Pool.Query(context.TODO(), query, orderUID)
	if err != nil {
		return nil, err
	}

	var chrtID string
	for rows.Next() {
		if err := rows.Scan(&chrtID); err != nil {
			return nil, err
		}
		chrtIDs = append(chrtIDs, chrtID)
	}

	return chrtIDs, nil
}
