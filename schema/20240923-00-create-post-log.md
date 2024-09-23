# スキーマ

```sql
create table nt_post_log (
    id int unsigned not null auto_increment,
    post_id int unsigned not null comment 'nt_post.id',
    url_key varchar(32) not null,
    created_datetime datetime not null,
    updated_datetime datetime not null,
    `text` text not null,
    visibility int not null,

    primary key (`id`),
    key `post_id` (`post_id`)
) comment 'nt_postのログテーブル';
```

# 説明

- nt_post に update を走らせるたびに insert する
- nt_post に insert したときは insert しない
- update, delete はしない
