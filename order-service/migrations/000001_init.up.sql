CREATE TYPE status AS ENUM ('awaitingPayment', 'inProcess', 'completed');

CREATE TABLE "order" (
    id UUID DEFAULT gen_random_uuid(),
    status status NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    user_id UUID NOT NULL
);

CREATE TABLE wallet (
    id UUID DEFAULT gen_random_uuid(),
    balance NUMERIC NOT NULL,
    user_id UUID NOT NULL
);

CREATE TABLE order_item (
    order_id UUID NOT NULL,
    FOREIGN KEY (order_id) REFERENCES "order"(id),
    cover_id UUID NOT NULL
)