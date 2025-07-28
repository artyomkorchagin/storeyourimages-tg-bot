-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    tg_id BIGINT NOT NULL UNIQUE,              
    status TEXT                       
);

CREATE TABLE IF NOT EXISTS content_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), 
    tg_id BIGINT NOT NULL,                 
    filepath TEXT,                          
    type TEXT,                              
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),   
    CONSTRAINT fk_tg_id FOREIGN KEY (tg_id) REFERENCES users(tg_id) ON DELETE CASCADE 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS content_data;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
