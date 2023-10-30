CREATE TABLE `suite_ticket`(
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `ticket` VARCHAR(255) NOT NULL
);

CREATE TABLE `corp_authz`(
    `id` INTEGER PRIMARY KEY AUTOINCREMENT,
    `corp_id`  VARCHAR(255) NOT NULL,
    `corp_name` VARCHAR(255) NOT NULL,
    `perm_code` VARCHAR(512) NOT NULL,
    `auth_info` JSON NOT NULL,
    `perm` JSON NOT NULL
);

CREATE UNIQUE INDEX `idx_corp_authz__corp_id` ON `corp_authz`(`corp_id`);

CREATE UNIQUE INDEX `idx_corp_authz__corp_name` ON `corp_authz`(`corp_name`);