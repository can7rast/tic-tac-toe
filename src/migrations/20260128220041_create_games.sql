-- +goose Up
CREATE TABLE IF NOT EXISTS games (
                       id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                       board          JSONB NOT NULL,
                       current_player INTEGER NOT NULL DEFAULT 1,          -- 1 = X (Player), 2 = O
                       player1_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       player2_id     UUID REFERENCES users(id) ON DELETE SET NULL,
                       state          INTEGER NOT NULL DEFAULT 0,          -- 0 = Waiting, 1 = TurnPlayer1, etc.
                       vsai           BOOLEAN NOT NULL DEFAULT FALSE,

    -- Опционально: когда игра создана
                       created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Полезные индексы
CREATE INDEX idx_games_state          ON games(state);
CREATE INDEX idx_games_player1_id     ON games(player1_id);
CREATE INDEX idx_games_vsai_waiting   ON games(vsai, state) WHERE state = 0;


-- +goose Down
DROP TABLE IF EXISTS games;