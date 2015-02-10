Based on ``pgtricks`` project [https://github.com/akaihola/pgtricks]

``pg_dump_splitsort.py`` ``mysqldump_splitsort.py`` are a handy scripts for pre-processing PostgreSQL's and MySQL's
``pg_dump`` and ``mysqldump`` output to make it more suitable for diffing and storing in version
control systems.

```
usage: pg_dump_splitsort.py [-h] [-m CHUNK_SIZE_KB] [-d DESTINATION_PATH] [-v]
                            [-c]
                            sql_dump_file

Split database dump file to a chunks.

positional arguments:
  sql_dump_file

optional arguments:
  -h, --help           show this help message and exit
  -m CHUNK_SIZE_KB     Max chunk size of database part, in kb default
                       2014(2Mb)
  -d DESTINATION_PATH  Path, where to store splitted files
  -v                   Verbose dumping output
  -c                   Clean destination
```

```
usage: mysqldump_splitsort.py [-h] [-m CHUNK_SIZE_KB] [-d DESTINATION_PATH]
                              [-v] [-c]
                              sql_dump_file

Split database dump file to a chunks.

positional arguments:
  sql_dump_file

optional arguments:
  -h, --help           show this help message and exit
  -m CHUNK_SIZE_KB     Max chunk size of database part, in kb default
                       2014(2Mb)
  -d DESTINATION_PATH  Path, where to store splitted files
  -v                   Verbose dumping output
  -c                   Clean destination
```

    -v - Verbose output
    -m - Maximum size of dumped data in one chunk [by default - 5MB]
    -c - Clean run. Removes \d\d\d\d_*.sql files from current directory before script run

The script splits the dump into the following files:

    0000_prologue.sql -everything up to the first COPY
    000X_<schema>_<table>_0001.sql
    000X_<schema>_<table>_0002.sql
    :
    :
    00YY_<schemax>_<tabley>_0001.sql - COPY/INSERT data for each table *sorted by the first field, and second fields*, splitted to chunks near specified maximum size.
    zzzz_epilogue.sql - everything after the last COPY
    
For more compact files names numbering are done wit a 36 base numbering sysem 0-9a-z.

For mysql dumps files numbered without schema ``00001_<tabley>_00001.sql``

Mysql dumps have to be prepared with a mysqldump options --skip-opt 

The files for table data are numbered uniquely and table order number is stored persistantly in .pgtricks file.

Backed up files can be used to re-create the database:

    $ cat *.sql | psql <database>


    $ cat *.sql | mysql -d <database>
