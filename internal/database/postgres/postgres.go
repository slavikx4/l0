package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/slavikx4/l0/internal/models"
	er "github.com/slavikx4/l0/pkg/error"
	"github.com/slavikx4/l0/pkg/logger"
)

type Postgres struct {
	*pgxpool.Pool
}

func NewPostgres(ctx context.Context, url string) (*Postgres, error) {
	const op = "NewPostgres -> "

	DB, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, &er.Error{Err: err, Code: er.ErrorNoConnect, Message: "не удалось подключиться к DataBase L0", Op: op}
	}
	if err := DB.Ping(context.Background()); err != nil {
		return nil, &er.Error{Err: err, Code: er.ErrorNoPing, Message: "не удалось пингануть к DataBase L0", Op: op}
	}
	logger.Logger.Process.Println("подключён успешно postgres")

	postgres := Postgres{Pool: DB}
	return &postgres, nil
}

func (p *Postgres) SetOrder(ctx context.Context, order *models.Order) error {
	const op = "Postgres.SetOrder -> "

	if err := p.setDelivery(ctx, &order.Delivery); err != nil {
		return er.AddOp(err, op)
	}

	if err := p.setPayment(ctx, &order.Payment); err != nil {
		return er.AddOp(err, op)
	}

	if err := p.setItems(ctx, &order.Items, order.OrderUID); err != nil {
		return er.AddOp(err, op)
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

	if _, err := p.Pool.Exec(ctx, query,
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
		return &er.Error{Err: err, Code: er.ErrorDataBaseLimitation, Message: "не удалось установить значение", Op: op}
	}

	return nil
}

func (p *Postgres) GetOrders(ctx context.Context) (*[]*models.Order, error) {
	const op = "Postgres.GetOrders -> "

	orders := make([]*models.Order, 0)

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

	rows, err := p.Pool.Query(ctx, query)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &er.Error{Err: err, Code: er.ErrorNotFound, Message: "не удалось найти строку", Op: op}

		} else {
			return nil, &er.Error{Err: err, Code: er.ErrorDataBaseIndefinite, Message: "не удалось выполнить запрос", Op: op}
		}
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
			return nil, &er.Error{Err: err, Code: er.ErrorDataBaseIndefinite, Message: "не удалось выполнить запрос", Op: op}
		}
		if err := p.getDelivery(ctx, &order.Delivery); err != nil {
			return nil, er.AddOp(err, op)
		}

		if err := p.getPayment(ctx, &order.Payment); err != nil {
			return nil, er.AddOp(err, op)
		}
		if err := p.getItems(ctx, &order.Items, order.OrderUID); err != nil {
			return nil, er.AddOp(err, op)
		}
		orders = append(orders, &order)
	}

	return &orders, nil
}

func (p *Postgres) setDelivery(ctx context.Context, delivery *models.Delivery) error {
	const op = "Postgres.setDelivery -> "

	query := `INSERT INTO "delivery"(
                       "phone",
                       "name",
                       "zip",
                       "city",
                       "address",
                       "region",
                       "email") VALUES ($1,$2,$3,$4,$5,$6,$7)`

	if _, err := p.Pool.Exec(ctx, query,
		delivery.Phone,
		delivery.Name,
		delivery.Zip,
		delivery.City,
		delivery.Address,
		delivery.Region,
		delivery.Email); err != nil {
		return &er.Error{Err: err, Code: er.ErrorDataBaseLimitation, Message: "не удалось установить значение", Op: op}
	}
	return nil
}

func (p *Postgres) setPayment(ctx context.Context, payment *models.Payment) error {
	const op = "Postgres.setPayment -> "

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

	if _, err := p.Pool.Exec(ctx, query,
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
		return &er.Error{Err: err, Code: er.ErrorDataBaseLimitation, Message: "не удалось установить значение", Op: op}
	}
	return nil
}

func (p *Postgres) setItems(ctx context.Context, items *models.Items, orderUID string) error {
	const op = "Postgres.setItems -> "

	for _, item := range *items {
		if err := p.setItem(ctx, &item); err != nil {
			return er.AddOp(err, op)
		}
		if err := p.setOrderOfItem(ctx, orderUID, item.ChrtID); err != nil {
			return er.AddOp(err, op)
		}
	}
	return nil
}

func (p *Postgres) setItem(ctx context.Context, item *models.Item) error {
	const op = "Postgres.setItem -> "

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

	if _, err := p.Pool.Exec(ctx, query,
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
		return &er.Error{Err: err, Code: er.ErrorDataBaseLimitation, Message: "не удалось установить значение", Op: op}
	}

	return nil
}

func (p *Postgres) setOrderOfItem(ctx context.Context, orderUID string, chrtID int) error {
	const op = "Postgres.setOrderOfItem -> "

	query := `INSERT INTO "order_item"(
                         "order_uid",
                         "chrt_id") VALUES ($1, $2)`
	if _, err := p.Pool.Exec(ctx, query, orderUID, chrtID); err != nil {
		return &er.Error{Err: err, Code: er.ErrorDataBaseLimitation, Message: "не удалось установить значение", Op: op}
	}
	return nil
}

func (p *Postgres) getDelivery(ctx context.Context, delivery *models.Delivery) error {
	const op = "Postgres.getDelivery -> "

	query := `SELECT 
				"phone",
				"name",
				"zip",
				"city",
				"address",
				"region",
				"email"
			FROM "delivery"`

	if err := p.Pool.QueryRow(ctx, query).Scan(
		&delivery.Phone,
		&delivery.Name,
		&delivery.Zip,
		&delivery.City,
		&delivery.Address,
		&delivery.Region,
		&delivery.Email); err != nil {
		if err == pgx.ErrNoRows {
			return &er.Error{Err: err, Code: er.ErrorNotFound, Message: "не удалось найти строку", Op: op}

		} else {
			return &er.Error{Err: err, Code: er.ErrorDataBaseIndefinite, Message: "не удалось выполнить запрос ", Op: op}
		}
	}

	return nil
}

func (p *Postgres) getPayment(ctx context.Context, payment *models.Payment) error {
	const op = "Postgres.getPayment -> "

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

	if err := p.Pool.QueryRow(ctx, query).Scan(
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
		if err == pgx.ErrNoRows {
			return &er.Error{Err: err, Code: er.ErrorNotFound, Message: "не удалось найти строку", Op: op}

		} else {
			return &er.Error{Err: err, Code: er.ErrorDataBaseIndefinite, Message: "не удалось выполнить запрос ", Op: op}
		}
	}

	return nil
}

func (p *Postgres) getItems(ctx context.Context, items *models.Items, orderUID string) error {
	const op = "Postgres.getItems -> "

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

	chrtIDs, err := p.getItemChrtIDsWithOrderUID(ctx, orderUID)
	if err != nil {
		return er.AddOp(err, op)
	}

	for _, chrtID := range chrtIDs {
		item := models.Item{}

		if err := p.Pool.QueryRow(ctx, query, chrtID).Scan(
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
			if err == pgx.ErrNoRows {
				return &er.Error{Err: err, Code: er.ErrorNotFound, Message: "не удалось найти строку", Op: op}

			} else {
				return &er.Error{Err: err, Code: er.ErrorDataBaseIndefinite, Message: "не удалось выполнить запрос", Op: op}
			}
		}

		*items = append(*items, item)
	}

	return nil
}

func (p *Postgres) getItemChrtIDsWithOrderUID(ctx context.Context, orderUID string) ([]string, error) {
	const op = "Postgres.getItemChrtIDsWitchOrderUID -> "

	query := `SELECT
					"chrt_id"
				FROM "order_item" WHERE ("order_uid"=$1)`

	chrtIDs := make([]string, 0, 1)

	rows, err := p.Pool.Query(ctx, query, orderUID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, &er.Error{Err: err, Code: er.ErrorNotFound, Message: "не удалось найти строку", Op: op}

		} else {
			return nil, &er.Error{Err: err, Code: er.ErrorDataBaseIndefinite, Message: "не удалось получить список chrtID", Op: op}
		}
	}

	var chrtID string
	for rows.Next() {
		if err := rows.Scan(&chrtID); err != nil {
			return nil, &er.Error{Err: err, Code: er.ErrorDataBaseIndefinite, Message: "ошибка при сканировании", Op: op}
		}
		chrtIDs = append(chrtIDs, chrtID)
	}

	return chrtIDs, nil
}
