create table userinfo(
    id int(11) not null auto_increment,
    user_name varchar(20) not null default ' ' comment "用户账号",
    user_pwd varchar(20) not null default ' ' comment "用户密码 encode by ?",
    email varchar(64) default '' comment "邮箱" ,
    tel varchar(20) default ' ' comment "电话号码",
    email_checked int(1) default 0 comment "邮箱是否验证",
    tel_checked int(1) default 0 comment "电话号码是否验证",
    signup_ts bigint(20) not null comment "注册时间",
    last_active_ts bigint(20) not null comment "上次登录时间",
    profile text comment "属性",
    status int(11) not null default 0 comment "文件当前状态(0可用/1禁用/2标记/3删除)", 
    primary key (id),
    unique key idx_tel (tel),
    key idx_status (status)
)engine=InnoDB auto_increment=5 default charset=utf8mb4;