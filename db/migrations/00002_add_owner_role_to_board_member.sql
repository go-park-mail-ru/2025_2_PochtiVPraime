-- +goose Up
-- +goose StatementBegin
ALTER TABLE board_member 
DROP CONSTRAINT board_member_member_role_check;

ALTER TABLE board_member 
ADD CONSTRAINT board_member_member_role_check 
CHECK (member_role IN ('owner', 'admin', 'member', 'observer'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
