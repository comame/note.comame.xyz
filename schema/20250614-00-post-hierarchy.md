## スキーマ

```sql
ALTER TABLE `nt_post`
  ADD COLUMN `parent` int unsigned COMMENT 'nullの場合、最上位とする' AFTER `id`,
  ADD COLUMN `permission_inherited` boolean DEFAULT TRUE AFTER `permission`;

ALTER TABLE `nt_post_log`
  ADD COLUMN `parent` int unsigned COMMENT 'nullの場合、最上位とする' AFTER `id`,
  ADD COLUMN `permission_inherited` boolean DEFAULT TRUE AFTER `permission`
  ;

ALTER TABLE `nt_post`
  ALTER COLUMN `permission_inherited` DROP DEFAULT;

ALTER TABLE `nt_post_log`
  ALTER COLUMN `permission_inherited` DROP DEFAULT;
```
