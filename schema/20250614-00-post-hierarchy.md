## スキーマ

```sql
ALTER TABLE `nt_post`
  ADD COLUMN `parent` int unsigned DEFAULT 0 COMMENT '0の場合、最上位とする' AFTER `id`,
  ADD COLUMN `permission_inherited` boolean DEFAULT TRUE AFTER `permission`;

ALTER TABLE `nt_post_log`
  ADD COLUMN `parent` int unsigned DEFAULT 0 COMMENT '0の場合、最上位とする' AFTER `id`,
  ADD COLUMN `permission_inherited` boolean DEFAULT TRUE AFTER `permission`
  ;

ALTER TABLE `nt_post`
  ALTER COLUMN `permission_inherited` DROP DEFAULT,
  ALTER COLUMN `parent` DROP DEFAULT;

ALTER TABLE `nt_post_log`
  ALTER COLUMN `permission_inherited` DROP DEFAULT,
  ALTER COLUMN `parent` DROP DEFAULT;
```

## 意図

- parent カラムは、親記事の ID を入れる
- permission_inherited は、親記事の権限を引き継ぐ場合に TRUE が入る。
    - これがFALSEになっているときだけ、自身の PERMISSION を参照する。
    - これが TRUE になっているときは、permission_inherited=false が設定されているかあるいは最上位記事に到達するまで、再帰的に探索する。
      - クエリの回数は多くなるが、ネスト数はせいぜい ~10 だろうから、いまのところはいったん気にしないことにする
