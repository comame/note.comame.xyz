# 目的

visibility カラムを permission へ移行する

# クエリ

## alter

```sql
ALTER TABLE `nt_post`
ADD COLUMN `permission` VARCHAR(16) AFTER `visibility`;

ALTER TABLE `nt_post_log`
ADD COLUMN `permission` VARCHAR(16) AFTER `visibility`;
```

## 既存の値を変換

```sql
UPDATE `nt_post`
SET `permission` = CASE
    WHEN `visibility` = 0 THEN 'private'
    WHEN `visibility` = 1 THEN 'url'
    WHEN `visibility` = 2 THEN 'public'
    ELSE 'private'
END;

UPDATE `nt_post_log`
SET `permission` = CASE
    WHEN `visibility` = 0 THEN 'private'
    WHEN `visibility` = 1 THEN 'url'
    WHEN `visibility` = 2 THEN 'public'
    ELSE 'private'
END;
```

## permission カラムへの NOT NULL 制約付与と visibility カラムの削除

```sql
ALTER TABLE `nt_post`
MODIFY COLUMN `permission` VARCHAR(16) NOT NULL,
DROP COLUMN `visibility`;

ALTER TABLE `nt_post_log`
MODIFY COLUMN `permission` VARCHAR(16) NOT NULL,
DROP COLUMN `visibility`;
```
