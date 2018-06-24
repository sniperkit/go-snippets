
CSV2SQL, An util to convert csv data to SQL scripts and load to MySQL.

## How it works

Pipeline:

[Read the xlsx file] -> [Convert to SQL scripts] -> [Execute SQL scripts]

## How to use

1. Prepare the xlsx file with your data

Assume the file is dropped into `test/import.xlsx`, you can download it from [here](https://github.com/medcl/csv2sql/blob/master/test/import.xlsx)

<img width="800"  src="https://github.com/medcl/csv2sql/raw/master/doc/assets/img/Snip20180505_8.png">

2. Config the `csv2sql.yml`

```
modules:
- name: pipeline
  enabled: true
  runners:
    - name: process_excel
      enabled: true
      input_queue: primary
      max_go_routine: 1
      threshold_in_ms: 0
      timeout_in_ms: 5000
      default_config:
        start:
          joint: read_csv
          enabled: true
          parameters:
            file_name: "test/import.xlsx"
        process:
        - joint: convert_sql
          enabled: true
          parameters:
            sheet_name: fish_information
            data_start_from_index: 3
            column_name:
            - id
            - outer_code
            - common_name
            - scientific_name
            - english_name
            - chinese_name
            - region_name
            - aquatic_category_id
            - category_name
            - is_homemade
            - aquatic_region_id
            - inner_code
            - produce_pattern
            - feed_pattern
            - catch_pattern
            row_format:
            - 'INSERT INTO `aquatic_base_info` (`id`, `outer_code`, `common_name`, `scientific_name`, `english_name`, `chinese_name`, `region_name`, `aquatic_category_id`)'
            - 'VALUES (<{id: }>, <{outer_code: }>, <{common_name: }>, <{scientific_name: }>, <{english_name: }>, <{chinese_name: }>, <{region_name: }>, <{aquatic_category_id: }>);'
            - 'INSERT INTO `aquatic_source` (`inner_code`, `aquatic_base_info_id`, `is_homemade`, `aquatic_region_id`, `produce_pattern`, `feed_pattern`, `catch_pattern`) '
            - 'VALUES (<{inner_code: }>, <{id: }>, <{is_homemade: }>, <{aquatic_region_id: }>, <{produce_pattern: }>, <{feed_pattern: }>, <{catch_pattern: }>);'

        - joint: convert_sql
          enabled: false
          parameters:
            sheet_name: aquatic_region
            data_start_from_index: 1
            column_name:
            - id
            - code
            - name
            - is_homemade
            row_format:
            - 'INSERT INTO `aquatic_region` (`id`, `code`, `name`, `is_homemade`) VALUES (<{id: }>, <{code: }>, <{name: }>, <{is_homemade: }>);'

        - joint: import_sql
          enabled: true
          parameters:
            mysql_conn: root:password@tcp(localhost:3306)/ifish?charset=utf8
            rollback_enabled: true

        - joint: logging
          enabled: true

        error:
          joint: on_error
          enabled: true
```

Note, there are more than one joint config with name: `convert_sql`, which means you can import multi data sheet with different config,
in this example the second `convert_sql` joint has been disabled, you can enable it if you wish, and also you can add more `convert_sql` joints.

And also as you can see, in the first `convert_sql` joint, with the data sheet `fish_information`, the config `row_format` is a array,
and have more than one SQL template, separated with `;` , which means you can generate multi SQL from one single data sheet,
map one data row to multi mysql data records, and then we can inset the data into different mysql tables.

The config `row_format` is a SQL template, and the config `column_name` is how your data sheet will be used in your template,
like this template variable `<{is_homemade: }>`, it will looking for the data from column `is_homemade` which we have already configured in section `column_name`,
you can get the SQL template by using MySQLWorkBench quickly([select db]->[select table]->[Copy to Clipboard]->[Insert Statement])
<img width="800"  src="https://github.com/medcl/csv2sql/raw/master/doc/assets/img/Snip20180505_9.png">

Note, once you change the column in data sheet, you must keep the `column_name` updated.

Note, please also change MySQL connection in joint `mysql_conn` config, ie: `your_mysql_user:you_password@tcp(your_mysql_host:3306)/your_mysql_db?charset=utf8`

And, the example is jus a example, you can use it to map any data sheet to any mysql table, update the config and import the data.

3. Start to import the data
```
➜  csv2sql git:(master) ✗ ./bin/csv2sql
  _____  _______      _____   _____  ____  _
 / ____|/ ____\ \    / /__ \ / ____|/ __ \| |
| |    | (___  \ \  / /   ) | (___ | |  | | |
| |     \___ \  \ \/ /   / / \___ \| |  | | |
| |____ ____) |  \  /   / /_ ____) | |__| | |____
 \_____|_____/    \/   |____|_____/ \___\_\______|
[CSV2SQL] An util to convert csv data to SQL scripts.
0.1.0_SNAPSHOT,  67026ee, Sat May 5 13:23:42 2018 +0800, medcl, support import multi datasheet

[05-05 15:29:38] [INF] [instance.go:23] workspace: data/APP/nodes/0
[05-05 15:29:38] [INF] [pipeline.go:67] pipeline: process_excel started with 1 instances
[05-05 15:29:38] [INF] [api.go:147] api server listen at: http://0.0.0.0:2900
[05-05 15:29:38] [INF] [ui.go:149] http server listen at: http://127.0.0.1:9001
[05-05 15:29:38] [INF] [import_sql.go:87] sql execute success, 1 rows affected, lastInsertID: 1
[05-05 15:29:38] [INF] [import_sql.go:87] sql execute success, 1 rows affected, lastInsertID: 1
^C
[CSV2SQL] got signal:interrupt, start shutting down
                         |    |
   _` |   _ \   _ \   _` |     _ \  |  |   -_)
 \__, | \___/ \___/ \__,_|   _.__/ \_, | \___|
 ____/                             ___/
[CSV2SQL] 0.1.0_SNAPSHOT, uptime:5.773621s

```

4. Now, check out the database, you will see 2 new records in two different table

<img width="800"  src="https://github.com/medcl/csv2sql/raw/master/doc/assets/img/Snip20180505_5.png">
<img width="800"  src="https://github.com/medcl/csv2sql/raw/master/doc/assets/img/Snip20180505_6.png">


License
=======
Released under the [Apache License, Version 2.0](https://github.com/medcl/csv2sql/blob/master/LICENSE) .

Powered by [https://github.com/infinitbyte/framework](https://github.com/infinitbyte/framework)


