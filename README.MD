# SB-Counter 色逼数

**alpha版, 谨慎部署**<sub>_~~作者都部署了, 还怕什么?~~_</sub>

好耶, 是色逼数!
<sub>_~~你心里, 没有逼数!~~_</sub>

## demo
| 数据 |
|:----:|
|<img src="https://sbc-io.xyz:81/_sbcounter/test/rw/card" width="100%">|

## 用途

记录访问量

## 怎么用

0. `clone`本项目, 或下载压缩包并解压

1. 创建数据库, 目前只支持`postgresql`:
```sql
DROP TABLE public.counter;

CREATE TABLE IF NOT EXISTS public.counter
(
    id text COLLATE pg_catalog."default" NOT NULL,
    db text COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone,
    ip text COLLATE pg_catalog."default" NOT NULL
)

TABLESPACE [表空间, 默认为 pg_default];

ALTER TABLE public.counter
    OWNER to "[数据表owner]";
```

2. 在config目录下, 复制`config.example.toml`到`config.toml`, 按需编辑`config.toml`:
+ 更改数据库ip, 端口号, 名称, 用户名和密码
+ 更改监听端口
+ 更改`sb-counter`设置

3. 编译
+ 准备好`golang`和`gocli`环境
+ `gf build main.go --pack public,template`

4. 部署
+ 运行编译好的`main`程序

5. 使用
+ 写入模式, 每次访问都会计数: `http://127.0.0.1:端口号/根目录/用户名/rw/card`
+ 只读模式, 只能查看, 不会记录: `http://127.0.0.1:端口号/根目录/用户名/ro/card`

## 框架

[GoFrame](https://goframe.org)

## 许可证

Apache 2.0
