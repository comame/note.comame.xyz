create table nt_post(
    id INT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
    url_key VARCHAR(32) NOT NULL,
    created_datetime DATETIME NOT NULL,
    updated_datetime DATETIME NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,

    unique url_key (url_key)
)
;

create table nt_post_visibility (
    post_id INT UNSIGNED NOT NULL,
    visibility INT NOT NULL COMMENT '0=非公開, 1=限定公開, 2=全体公開',

    KEY post_id (post_id)
)
;
