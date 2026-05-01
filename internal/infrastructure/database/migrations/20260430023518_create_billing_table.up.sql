CREATE TABLE billing_mpesa_transactions (
    id                   UUID          PRIMARY KEY,
    user_id              UUID          NOT NULL,
    phone                VARCHAR(15)   NOT NULL,
    account_name         VARCHAR(255)  NOT NULL,
    checkout_request_id  VARCHAR(100)  NOT NULL UNIQUE,
    merchant_request_id  VARCHAR(100)  NOT NULL,
    status               VARCHAR(20)   NOT NULL DEFAULT 'pending',
    mpesa_receipt_number VARCHAR(50)   NOT NULL DEFAULT '',
    amount               NUMERIC(10,2) NOT NULL,
    created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_mpesa_tx_user_id            ON billing_mpesa_transactions(user_id);
CREATE INDEX idx_mpesa_tx_checkout_req_id    ON billing_mpesa_transactions(checkout_request_id);
CREATE INDEX idx_mpesa_tx_status_created_at  ON billing_mpesa_transactions(status, created_at);



CREATE TABLE billing_cards (
    id               UUID         PRIMARY KEY,
    user_id          UUID         NOT NULL,
    cardholder_name  VARCHAR(255) NOT NULL,
    masked_card      VARCHAR(25)  NOT NULL,
    last_four        CHAR(4)      NOT NULL,
    brand            VARCHAR(20)  NOT NULL,
    expiration_month CHAR(2)      NOT NULL,
    expiration_year  CHAR(4)      NOT NULL,
    billing_address  TEXT         NOT NULL,
    billing_city     VARCHAR(100) NOT NULL,
    billing_state    VARCHAR(100) NOT NULL,
    billing_zip      VARCHAR(20)  NOT NULL,
    verified_at      TIMESTAMPTZ  NOT NULL,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billing_cards_user_id ON billing_cards(user_id);




CREATE TABLE billing_invoices (
    id                   UUID          PRIMARY KEY,
    user_id              UUID          NOT NULL,
    invoice_number       VARCHAR(50)   NOT NULL UNIQUE,
    amount               NUMERIC(10,2) NOT NULL,
    status               VARCHAR(20)   NOT NULL DEFAULT 'pending',
    due_date             TIMESTAMPTZ   NOT NULL,
    created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billing_invoices_user_id ON billing_invoices(user_id);
CREATE INDEX idx_billing_invoices_status ON billing_invoices(status);   



CREATE TABLE billing_invoice_payments (
    id                   UUID          PRIMARY KEY,
    invoice_id           UUID          NOT NULL,
    payment_method       VARCHAR(20)   NOT NULL,
    amount               NUMERIC(10,2) NOT NULL,
    status               VARCHAR(20)   NOT NULL DEFAULT 'pending',
    mpesa_receipt_number VARCHAR(50)   NOT NULL DEFAULT '',
    created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billing_invoice_payments_invoice_id ON billing_invoice_payments(invoice_id);
CREATE INDEX idx_billing_invoice_payments_status ON billing_invoice_payments(status);


CREATE TABLE billing_subscriptions (
    id                   UUID          PRIMARY KEY,
    user_id              UUID          NOT NULL,
    plan_id              UUID          NOT NULL,
    status               VARCHAR(20)   NOT NULL DEFAULT 'active',
    current_period_start TIMESTAMPTZ   NOT NULL,
    current_period_end   TIMESTAMPTZ   NOT NULL,
    created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billing_subscriptions_user_id ON billing_subscriptions(user_id);
CREATE INDEX idx_billing_subscriptions_status ON billing_subscriptions(status);


CREATE TABLE billing_subscription_payments (
    id                   UUID          PRIMARY KEY,
    subscription_id      UUID          NOT NULL,
    payment_method       VARCHAR(20)   NOT NULL,
    amount               NUMERIC(10,2) NOT NULL,
    status               VARCHAR(20)   NOT NULL DEFAULT 'pending',
    mpesa_receipt_number VARCHAR(50)   NOT NULL DEFAULT '',
    created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billing_subscription_payments_subscription_id ON billing_subscription_payments(subscription_id);
CREATE INDEX idx_billing_subscription_payments_status ON billing_subscription_payments(status); 


CREATE TABLE billing_subscription_invoices (
    id                   UUID          PRIMARY KEY,
    subscription_id      UUID          NOT NULL,
    invoice_number       VARCHAR(50)   NOT NULL UNIQUE,
    amount               NUMERIC(10,2) NOT NULL,
    status               VARCHAR(20)   NOT NULL DEFAULT 'pending',
    due_date             TIMESTAMPTZ   NOT NULL,
    created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billing_subscription_invoices_subscription_id ON billing_subscription_invoices(subscription_id);
CREATE INDEX idx_billing_subscription_invoices_status ON billing_subscription_invoices(status); 

CREATE TABLE billings (
    id                   UUID          PRIMARY KEY,
    user_id              UUID          NOT NULL,
    service_name         VARCHAR(255)  NOT NULL,
    amount               NUMERIC(10,2) NOT NULL,
    currency             VARCHAR(3)    NOT NULL,
    billing_date         TIMESTAMPTZ   NOT NULL,
    status               VARCHAR(20)   NOT NULL DEFAULT 'pending',
    created_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_billings_user_id ON billings(user_id);
CREATE INDEX idx_billings_service_name ON billings(service_name);
CREATE INDEX idx_billings_status ON billings(status); 