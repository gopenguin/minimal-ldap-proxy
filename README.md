# minimal-ldap-proxy
Proxy ldap authentication requests to a database backend

## Getting started

For local testing you can use a `sqlite` database. You can create with the following commands:

```sh
sqlite3 test.db '
CREATE TABLE "users"
(
  id INTEGER not null primary key autoincrement,
  name TEXT not null unique,
  password TEXT not null,
  gname TEXT not null,
  sname TEXT not null,
  email TEXT NULL
);
CREATE TABLE groups
(
    id integer PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL
);
CREATE TABLE user_groups
(
    user_id integer NOT NULL,
    group_id integer NOT NULL,
    CONSTRAINT user_groups_user_id_group_id_pk PRIMARY KEY (user_id, group_id)
);
sqlite3 test.db ''
'
```



```yaml
driver: sqlite3
conn: "./test.db"
authQuery: "select password from users where name = ?"
searchQuery: "select u.name  as cn, u.gname as gn, u.sname as sn, u.email as mail, g.name  as memberOf from users as u join user_groups as ug on (u.id = ug.user_id) join groups as g on (g.id = ug.group_id) where u.name = ?"
attributes:
  - cn
  - gn
  - sn
  - mail
  - memberOf
baseDn: "ou=People,dc=example,dc=com"
rdn: "cn"
```

