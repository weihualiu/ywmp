-- dbanme ywmp
-- 建立表空间
drop tablespace ywmp_ts;
create tablespace ywmp_ts location '/home/liuweihua/workspace/ywmpdb';

/*
 * 日志类型表
 */
drop table lg_type;
create table lg_type (
	tid integer primary key,
	tname varchar(100)
) tablespace ywmp_ts;

/*
 * 项目接口信息表 program table
 */
drop table pg_inf_info;
create table pg_inf_info (
	infId varchar(11) primary key,
	infName varchar(100)
) tablespace ywmp_ts;
/*
 * 接口响应时间
 */
drop table lg_inf_usetime;
create table lg_inf_usetime (
	id bigserial primary key,
	infId varchar(11),
	-- 毫秒数
	userTime integer,
	nodeId integer
) tablespace ywmp_ts;

/*
 * 客户端崩溃日志信息表
 */
drop table lg_cli_crash;
create table lg_cli_crash (
	id serial primary key,
	deviceType varchar(20),
	deviceId varchar(100),
	osVer varchar(20),
	appVer varchar(9),
	-- IP
	cliAddr cidr,
	loginName varchar(50),
	-- 崩溃数据产生时间
	createDate timestamp,
	bundleId varchar(50),
	-- 崩溃摘要
	summary varchar(500),
	-- 崩溃全文
	content text,
	-- 节点编号
	nodeId integer,
	-- 入库时间
	stamp timestamp
) tablespace ywmp_ts;

/*
 * 服务器节点配置信息
 */
drop table node_conf;
create table node_conf (
	id integer primary key,
	name varchar(100),
	addr cidr,
	extand varchar(100)
) tablespace ywmp_ts;

/*
 * 功能系统信息表
 */
drop table program;
create table program (
	id varchar(20) primary key,
	name varchar(100),
	description text
) tablespace ywmp_ts;

/*
 * 项目与接口关系表
 */
drop table pg_inf_rel;
create table pg_inf_rel (
	pgid varchar(20),
	infId varchar(20)
) tablespace ywmp_ts;

