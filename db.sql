CREATE TABLE "order" (
	"id" serial NOT NULL,
	"order_uid" character varying NOT NULL UNIQUE,
	"track_number" character varying NOT NULL,
	"entry" character varying NOT NULL,
	"delivery" character varying NOT NULL,
	"payment" character varying NOT NULL,
	"locale" character varying NOT NULL,
	"internal_signature" character varying,
	"customer_id" character varying NOT NULL,
	"delivery_service" character varying NOT NULL,
	"shardkey" character varying NOT NULL,
	"sm_id" int NOT NULL,
	"date_created" character varying NOT NULL,
	"oof_shard" character varying NOT NULL,
	CONSTRAINT "order_pk" PRIMARY KEY ("id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "delivery" (
	"phone" character varying NOT NULL UNIQUE,
	"name" character varying NOT NULL,
	"zip" character varying NOT NULL,
	"city" character varying NOT NULL,
	"address" character varying NOT NULL,
	"region" character varying NOT NULL,
	"email" character varying NOT NULL,
	CONSTRAINT "delivery_pk" PRIMARY KEY ("phone")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "payment" (
	"transaction" character varying NOT NULL,
	"request_id" character varying,
	"currency" character varying NOT NULL,
	"provider" character varying NOT NULL,
	"amount" int NOT NULL,
	"payment_dt" int NOT NULL,
	"bank" character varying NOT NULL,
	"delivery_cost" int NOT NULL,
	"goods_total" int NOT NULL,
	"custom_fee" int NOT NULL,
	CONSTRAINT "payment_pk" PRIMARY KEY ("transaction")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "items" (
	"chrt_id" int NOT NULL,
	"track_number" character varying NOT NULL,
	"price" int NOT NULL,
	"rid" character varying NOT NULL,
	"name" character varying NOT NULL,
	"sale" int NOT NULL,
	"size" character varying NOT NULL,
	"total_price" int NOT NULL,
	"nm_id" int NOT NULL,
	"brand" character varying NOT NULL,
	"status" int NOT NULL,
	CONSTRAINT "items_pk" PRIMARY KEY ("chrt_id")
) WITH (
  OIDS=FALSE
);



CREATE TABLE "order_item" (
	"order_uid" character varying NOT NULL,
	"chrt_id" int NOT NULL,
	CONSTRAINT "order_item_pk" PRIMARY KEY ("order_uid")
) WITH (
  OIDS=FALSE
);



ALTER TABLE "order" ADD CONSTRAINT "order_fk0" FOREIGN KEY ("order_uid") REFERENCES "order_item"("order_uid");
ALTER TABLE "order" ADD CONSTRAINT "order_fk1" FOREIGN KEY ("delivery") REFERENCES "delivery"("phone");
ALTER TABLE "order" ADD CONSTRAINT "order_fk2" FOREIGN KEY ("payment") REFERENCES "payment"("transaction");




ALTER TABLE "order_item" ADD CONSTRAINT "order_item_fk0" FOREIGN KEY ("chrt_id") REFERENCES "items"("chrt_id");






