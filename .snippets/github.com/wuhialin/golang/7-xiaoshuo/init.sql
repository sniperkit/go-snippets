DROP TABLE source;
CREATE TABLE "source" (
  "id" serial4,
  "name" varchar(50) NOT NULL DEFAULT '',
  "domain" varchar(30) NOT NULL DEFAULT '',
  "url" varchar(50) NOT NULL DEFAULT '',
  "main_action" varchar(50) NOT NULL DEFAULT '',
  "state" int2 NOT NULL DEFAULT 1,
  "create_at" int4 NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
COMMENT ON COLUMN "source"."id" IS '主键ID';
COMMENT ON COLUMN "source"."name" IS '名称';
COMMENT ON COLUMN "source"."domain" IS '域名|来源';
COMMENT ON COLUMN "source"."url" IS '地址';
COMMENT ON COLUMN "source"."main_action" IS '入口地址';
COMMENT ON COLUMN "source"."state" IS '状态';
COMMENT ON COLUMN "source"."create_at" IS '添加时间';
CREATE UNIQUE INDEX source_uniq_key ON source(
  domain
);

DROP TABLE types;
CREATE TABLE "types" (
  "id" serial4,
  "name" varchar(255) NOT NULL DEFAULT '',
  "code" varchar(50) NOT NULL DEFAULT '',
  "table_name" varchar(50) NOT NULL DEFAULT '',
  "create_at" int4 NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
COMMENT ON COLUMN "types"."id" IS '主键ID';
COMMENT ON COLUMN "types"."name" IS '名称';
COMMENT ON COLUMN "types"."code" IS '代号';
COMMENT ON COLUMN "types"."table_name" IS '表名，记录入库指定名称的表';
COMMENT ON COLUMN "types"."create_at" IS '添加时间';
CREATE UNIQUE INDEX types_uniq_key ON types (
  code
);

DROP TABLE selector;
CREATE TABLE "selector" (
  "id" serial4,
  "source_id" int4 NOT NULL DEFAULT 0,
  "type_id" int4 NOT NULL DEFAULT 0,
  "selector" text NOT NULL DEFAULT '',
  "eq" int2 NOT NULL DEFAULT -1,
  "sorting" int2 NOT NULL DEFAULT 0,
  "pid" int4 NOT NULL DEFAULT 0,
  "state" int2 NOT NULL DEFAULT 1,
  "create_at" int4 NOT NULL DEFAULT 0,
  "update_at" int4 NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
COMMENT ON COLUMN "selector"."id" IS '主键ID';
COMMENT ON COLUMN "selector"."source_id" IS '来源';
COMMENT ON COLUMN "selector"."type_id" IS 'type id';
COMMENT ON COLUMN "selector"."selector" IS '选择器';
COMMENT ON COLUMN "selector"."eq" IS '定位';
COMMENT ON COLUMN "selector"."sorting" IS '排序 小到大';
COMMENT ON COLUMN "selector"."pid" IS '父ID';
COMMENT ON COLUMN "selector"."state" IS '状态';
COMMENT ON COLUMN "selector"."create_at" IS '添加时间';
COMMENT ON COLUMN "selector"."update_at" IS '修改时间';


DROP TABLE url;
CREATE TABLE "url" (
    "id" serial4,
    "url" varchar(1024) NOT NULL DEFAULT '',
    "source_id" int4 NOT NULL DEFAULT 0,
    "state" int2 NOT NULL DEFAULT 0,
    "create_at" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY("id")
);
COMMENT ON COLUMN "url"."id" IS '主键ID';
COMMENT ON COLUMN "url"."url" IS '地址';
COMMENT ON COLUMN "url"."source_id" IS '来源';
COMMENT ON COLUMN "url"."state" IS '状态';
COMMENT ON COLUMN "url"."create_at" IS '添加时间';
CREATE UNIQUE INDEX "url_uniq_url" ON "url"(
  "url"
);
CREATE INDEX "url_idx_source" ON "url" (
  "source_id",
  "create_at" DESC
);

DROP TABLE html;
CREATE TABLE "html" (
    "id" serial4,
    "data" text NOT NULL DEFAULT '',
    "url_id" int4 NOT NULL DEFAULT 0,
    "create_at" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY("id")
);
COMMENT ON COLUMN "html"."id" IS '主键ID';
COMMENT ON COLUMN "html"."data" IS '页面内容';
COMMENT ON COLUMN "html"."url_id" IS '地址ID';
COMMENT ON COLUMN "html"."create_at" IS '添加时间';
CREATE UNIQUE INDEX "html_idx_url" ON "html" (
    "url_id"
);

DROP TABLE parsed;
CREATE TABLE "parsed" (
    "id" serial4,
    "html_id" int4 NOT NULL DEFAULT 0,
    "selector_id" int4 NOT NULL DEFAULT 0,
    "create_at" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY("id")
);
COMMENT ON COLUMN "parsed"."id" IS '主键ID';
COMMENT ON COLUMN "parsed"."html_id" IS '页面内容';
COMMENT ON COLUMN "parsed"."selector_id" IS '选择器';
COMMENT ON COLUMN "parsed"."create_at" IS '添加时间';
CREATE INDEX "parsed_idx_key" ON "parsed" (
    "html_id",
    "selector_id"
);

DROP TABLE category;
CREATE TABLE "category" (
	"id" serial4,
	"name" varchar(20) NOT NULL DEFAULT '',
	"state" int2 NOT NULL DEFAULT 1,
	"create_at" int4 NOT NULL DEFAULT 0,
	PRIMARY KEY("id")
);
COMMENT ON COLUMN "category"."id" IS '主键ID';
COMMENT ON COLUMN "category"."name" IS '名称';
COMMENT ON COLUMN "category"."state" IS '状态';
COMMENT ON COLUMN "category"."create_at" IS '添加时间';
CREATE UNIQUE INDEX "category_uniq_key" ON "category"(
    name
);

DROP TABLE tag;
CREATE TABLE "tag" (
	"id" serial4,
	"name" varchar(50) NOT NULL DEFAULT '',
	"state" int2 NOT NULL DEFAULT 1,
	"create_at" int4 NOT NULL DEFAULT 0,
	PRIMARY KEY("id")
);
COMMENT ON COLUMN "tag"."id" IS '主键ID';
COMMENT ON COLUMN "tag"."name" IS '名称';
COMMENT ON COLUMN "tag"."state" IS '状态';
COMMENT ON COLUMN "tag"."create_at" IS '添加时间';

DROP TABLE author;
CREATE TABLE "author" (
	"id" serial4,
	"name" varchar(50) NOT NULL DEFAULT '',
	"state" int2 NOT NULL DEFAULT 1,
	"create_at" int4 NOT NULL DEFAULT 0,
	PRIMARY KEY("id")
);
COMMENT ON COLUMN "author"."id" IS '主键ID';
COMMENT ON COLUMN "author"."name" IS '名称';
COMMENT ON COLUMN "author"."state" IS '状态';
COMMENT ON COLUMN "author"."create_at" IS '添加时间';

DROP TABLE book;
CREATE TABLE "book" (
	  "id" serial4,
	  "name" varchar(50) NOT NULL DEFAULT '',
    "author_id" int4 NOT NULL DEFAULT 0,
    "desc" text,
	  "state" int2 NOT NULL DEFAULT 1,
	  "create_at" int4 NOT NULL DEFAULT 0,
	  PRIMARY KEY("id")
);
COMMENT ON COLUMN "book"."id" IS '主键ID';
COMMENT ON COLUMN "book"."name" IS '名称';
COMMENT ON COLUMN "book"."author_id" IS '作者';
COMMENT ON COLUMN "book"."desc" IS '描述';
COMMENT ON COLUMN "book"."state" IS '状态';
COMMENT ON COLUMN "book"."create_at" IS '添加时间';

DROP TABLE book_tag;
CREATE TABLE "book_tag"(
    "id" serial4,
    "book_id" int4 NOT NULL DEFAULT 0,
    "tag_id" int4 NOT NULL DEFAULT 0,
    "state" int2 not NULL DEFAULT 1,
    "create_at" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY("id")
);
COMMENT ON COLUMN "book_tag"."id" IS '主键ID';
COMMENT ON COLUMN "book_tag"."book_id" IS '书名';
COMMENT ON COLUMN "book_tag"."tag_id" IS '标签';
COMMENT ON COLUMN "book_tag"."state" IS '状态';
COMMENT ON COLUMN "book_tag"."create_at" IS '添加时间';

DROP TABLE book_category;
CREATE TABLE "book_category"(
    "id" serial4,
    "book_id" int4 NOT NULL DEFAULT 0,
    "category_id" int4 NOT NULL DEFAULT 0,
    "state" int2 NOT NULL DEFAULT 1,
    "create_at" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY("id")
);
COMMENT ON COLUMN "book_category"."id" IS '主键ID';
COMMENT ON COLUMN "book_category"."book_id" IS '书名';
COMMENT ON COLUMN "book_category"."category_id" IS '分类';
COMMENT ON COLUMN "book_category"."state" IS '状态';
COMMENT ON COLUMN "book_category"."create_at" IS '添加时间';

DROP TABLE book_pages;
CREATE TABLE "book_pages"(
    "id" serial4,
    "book_id" int4 NOT NULL DEFAULT 0,
    "title" varchar(50) NOT NULL DEFAULT '',
    "sorting" int2 NOT NULL DEFAULT 0,
    "state" int2 NOT NULL DEFAULT 1,
    "detail" text NOT NULL,
    "create_at" int4 NOT NULL DEFAULT 0,
    PRIMARY KEY("id")
);
COMMENT ON COLUMN "book_pages"."id" IS '主键ID';
COMMENT ON COLUMN "book_pages"."book_id" IS '书名';
COMMENT ON COLUMN "book_pages"."title" IS '标题';
COMMENT ON COLUMN "book_pages"."sorting" IS '排序';
COMMENT ON COLUMN "book_pages"."state" IS '状态';
COMMENT ON COLUMN "book_pages"."detail" IS '详情';
COMMENT ON COLUMN "book_pages"."create_at" IS '添加时间';

DROP TABLE task;
CREATE TABLE "task" (
    id serial4,
    type_id int4 NOT NULL DEFAULT 0,
    title varchar(255) NOT NULL DEFAULT '',
    state int2 NOT NULL DEFAULT 1,
    create_at int4 NOT NULL DEFAULT 0,
    start_at int4 NOT NULL DEFAULT 0,
    end_at int4 NOT NULL DEFAULT 0,
    PRIMARY KEY(id)
);
COMMENT ON COLUMN "task"."id" IS '主键ID';
COMMENT ON COLUMN "task"."type_id" IS '类型';
COMMENT ON COLUMN "task"."title" IS '标题';
COMMENT ON COLUMN "task"."state" IS '状态';
COMMENT ON COLUMN "task"."create_at" IS '添加时间';
COMMENT ON COLUMN "task"."end_at" IS '结束时间';
