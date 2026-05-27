-- Fix super admin password hash (run after initial schema)
UPDATE `users`
SET `password_hash` = '$2b$12$aDcf2hFro7SZ1XRH.s.l4eH.7pD87fyvHJFMVBN/uGwiytW4vbgVS',
    `role` = 9,
    `status` = 1,
    `email_verified` = 1,
    `updated_at` = NOW()
WHERE `email` = 'admin@trbbtw.com';

-- If not exists, insert
INSERT INTO `users`
  (`uuid`,`username`,`email`,`phone`,`password_hash`,`display_name`,`role`,`status`,`email_verified`,`created_at`,`updated_at`)
SELECT UUID(),'superadmin','admin@trbbtw.com','0000000000',
   '$2b$12$aDcf2hFro7SZ1XRH.s.l4eH.7pD87fyvHJFMVBN/uGwiytW4vbgVS','TRBB 超級管理員',9,1,1,NOW(),NOW()
WHERE NOT EXISTS (SELECT 1 FROM `users` WHERE email='admin@trbbtw.com');
