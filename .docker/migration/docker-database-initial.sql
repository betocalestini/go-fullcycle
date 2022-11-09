CREATE TABLE IF NOT EXISTS "orders" (
  "id" varchar(255), 
  "price" FLOAT NOT NULL, 
  "tax" FLOAT NOT NULL, 
  "final_price" FLOAT NOT NULL, 
  PRIMARY KEY (id));