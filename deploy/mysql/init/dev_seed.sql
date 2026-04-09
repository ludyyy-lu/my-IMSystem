-- Development seed data for my-IMSystem.
-- Run this after the services have auto-migrated the schema.

SET NAMES utf8mb4;

INSERT INTO users (
    id, username, password, nickname, avatar, bio, created_at, disabled, gender
) VALUES
    (1, 'alice', '$2a$12$PshG8MpkUd5B1YMWrVRST.s/C5.Cqy7i3HveT3Yqiudr57oxxn0ky', 'Alice', '', 'Hello, I am Alice.', '2026-04-08 08:00:00', 0, 2),
    (2, 'bob',   '$2a$12$OcgrAx0sm.7/DvIAUOYHdeymlFgTMViivDGGVduar/r7EUHBHATWe', 'Bob', '', 'Bob is online for frontend checks.', '2026-04-08 08:01:00', 0, 1),
    (3, 'charlie', '$2a$12$HTcUmAaatk/CVr0.ahd6AuV41xKplnMOBAIXnGR1Egs689H2h0xXa', 'Charlie', '', 'Charlie is available for search and friend request tests.', '2026-04-08 08:02:00', 0, 1)
ON DUPLICATE KEY UPDATE
    username = VALUES(username),
    password = VALUES(password),
    nickname = VALUES(nickname),
    avatar = VALUES(avatar),
    bio = VALUES(bio),
    disabled = VALUES(disabled),
    gender = VALUES(gender);

INSERT INTO friends (
    id, user_id, friend_id, created_at, updated_at
) VALUES
    (1, 1, 2, NOW(), NOW()),
    (2, 2, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE
    user_id = VALUES(user_id),
    friend_id = VALUES(friend_id),
    updated_at = VALUES(updated_at);

INSERT INTO friend_requests (
    id, from_user_id, to_user_id, remark, status, created_at, updated_at
) VALUES
    (1, 3, 1, 'Hi Alice, please add me.', 'pending', NOW(), NOW())
ON DUPLICATE KEY UPDATE
    from_user_id = VALUES(from_user_id),
    to_user_id = VALUES(to_user_id),
    remark = VALUES(remark),
    status = VALUES(status),
    updated_at = VALUES(updated_at);

INSERT INTO message (
    id, from_user_id, to_user_id, content, msg_type, status, created_at
) VALUES
    (1, 2, 1, 'Alice, I sent a test message.', 0, 0, '2026-04-08 08:10:00'),
    (2, 1, 2, 'Received. Frontend and backend are connected.', 0, 1, '2026-04-08 08:11:00'),
    (3, 2, 1, 'This message is unread to test badge count.', 0, 0, '2026-04-08 08:12:00')
ON DUPLICATE KEY UPDATE
    from_user_id = VALUES(from_user_id),
    to_user_id = VALUES(to_user_id),
    content = VALUES(content),
    msg_type = VALUES(msg_type),
    status = VALUES(status);