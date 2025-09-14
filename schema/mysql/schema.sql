CREATE TABLE products
(
    id                 INT PRIMARY KEY,
    name               VARCHAR(265)                          NOT NULL,
    description        TEXT,
    price              DECIMAL(10, 2)                        NOT NULL,
    stock              INT                                   NOT NULL
);