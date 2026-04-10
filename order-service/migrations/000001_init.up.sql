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

CREATE TABLE orders_history (
    user_id UUID NOT NULL,
    order_id UUID NOT NULL
);

CREATE TABLE cart (
    id UUID DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL
);

CREATE TABLE carts_orders (
    cart_id UUID NOT NULL,
    cover_id UUID NOT NULL
);

