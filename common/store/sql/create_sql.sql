/* create database fileDB; */

create table filemeta(
    id int(11) not null auto_increment,
    file_sha1 char(40) not null default ' ' comment "文件哈希",
    file_name varchar(200) not null default ' ' comment "文件名",
    file_size bigint(20) default 0 comment "文件大小" ,
    file_local_path varchar(1024) not null default ' ' comment "文件在本地存储的路径",
    upload_timestamp bigint(20) not null comment "文件上传时间",
    last_modify_timestamp bigint(20) not null comment "文件最近修改时间",
    status int(11) not null default 0 comment "文件当前状态(0可用/1禁用/2删除)", 
    ext1 int(11) default 0 comment "备用字段1",
    ext2 int(11) default 0 comment "备用字段2",
    primary key (id),
    unique key idx_file_hash (file_sha1),
    key idx_status (status)
)engine=InnoDB default charset=utf8;