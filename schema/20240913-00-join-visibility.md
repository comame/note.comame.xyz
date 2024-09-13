```sql
alter table nt_post
add column visibility int not null default 0 comment '0=非公開, 1=限定公開, 2=全体公開'
;

alter table nt_post
add index visibility (visibility)
;

update nt_post, nt_post_visibility
set nt_post.visibility = nt_post_visibility.visibility
where nt_post.id = nt_post_visibility.post_id
;

select
    self.id,
    self.visibility as 'dst',
    v.visibility as 'src'
from nt_post as self
inner join nt_post_visibility as v
on v.post_id = self.id
;

drop table nt_post_visibility;

alter table nt_post alter column visibility drop default;
```

# やったこと

- JOIN で事故りそうなので、 `nt_post_visibility` を `nt_post` にくっつける
