create table t1 (b char(0));
create table t1 (b char(0) not null);
create table if not exists t1 (b char(0) not null);
create table t2 engine=heap select * from t1;
create table t1 (ordid int(8) not null auto_increment, ord  varchar(50) not null, primary key (ord,ordid)) engine=heap;
create table mysqltest.$test1 (a$1 int, $b int, c$ int);
create table t2 (b int) select a as b, a+1 as b from t1;
create table t1 select if('2002'='2002','Y','N');
create table t1 ( k1 varchar(2), k2 int, primary key(k1,k2));
insert into t1 values ("a", 1), ("b", 2);
create table t2
select
    a,
    ifnull(b,cast(-7                        as signed))   as b,
    ifnull(c,cast(7                         as unsigned)) as c,
    ifnull(d,cast('2000-01-01'              as date))     as d,
    ifnull(e,cast('b'                       as char))     as e,
    ifnull(f,cast('2000-01-01'              as datetime)) as f,
    ifnull(g,cast('5:4:3'                   as time))     as g,
    ifnull(h,cast('yet another binary data' as binary))   as h,
    addtime(cast('1:0:0' as time),cast('1:0:0' as time))  as dd
from t1;
CREATE TABLE t1(id varchar(10) NOT NULL PRIMARY KEY, dsc longtext);
INSERT INTO t1 VALUES ('5000000001', NULL),('5000000003', 'Test'),('5000000004', NULL);
create table t1 (
    a varchar(112) charset utf8 collate utf8_bin not null,
    primary key (a)
) select 'test' as a ;
create table טבלה_של_אריאל
(
    כמות int
);



CREATE TABLE t1(
    c1 INT DEFAULT 12 COMMENT 'column1',
    c2 INT NULL COMMENT 'column2',
    c3 INT NOT NULL COMMENT 'column3',
    c4 VARCHAR(255) CHARACTER SET utf8 NOT NULL DEFAULT 'a',
    c5 VARCHAR(255) COLLATE utf8_unicode_ci NULL DEFAULT 'b',
    c6 VARCHAR(255))
COLLATE latin1_bin;