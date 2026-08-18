CREATE TABLE product (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);

SELECT * FROM product;

INSERT INTO product (name, price) VALUES
('Camiseta', 30.99),
('Calça Jeans', 89.99),
('Bermuda', 25.49);
