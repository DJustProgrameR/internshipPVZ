CREATE TABLE users (
                       id UUID PRIMARY KEY,
                       email TEXT UNIQUE NOT NULL,
                       password TEXT NOT NULL,
                       role TEXT NOT NULL CHECK (role IN ('employee', 'moderator'))
);

CREATE TABLE pvz (
                     id UUID PRIMARY KEY,
                     registration_date TIMESTAMPTZ NOT NULL,
                     city TEXT NOT NULL CHECK (city IN ('Москва', 'Санкт-Петербург', 'Казань'))
);

CREATE TABLE receptions (
                            id UUID PRIMARY KEY,
                            pvz_id UUID NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
                            date_time TIMESTAMPTZ NOT NULL,
                            status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
);

CREATE TABLE products (
                          id UUID PRIMARY KEY,
                          reception_id UUID NOT NULL REFERENCES receptions(id) ON DELETE CASCADE,
                          date_time TIMESTAMPTZ NOT NULL,
                          type TEXT NOT NULL CHECK (type IN ('электроника', 'одежда', 'обувь'))
);

CREATE TABLE logs (
                          id UUID PRIMARY KEY,
                          user_id UUID REFERENCES users(id),
                          reception_id UUID REFERENCES receptions(id),
                          product_id UUID REFERENCES products(id),
                          date_time TIMESTAMPTZ NOT NULL,
                          eventType TEXT NOT NULL CHECK (type IN ('pvzCreate',
                                                                  'receptionOpened',
                                                                  'receptionClosed',
                                                                  'productAdded',
                                                                  'productDeleted',
                                                                  'userCreated',
                                                                  'moderatorCreated'))
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_receptions_pvz_id ON receptions(pvz_id);
CREATE INDEX idx_products_reception_id ON products(reception_id);