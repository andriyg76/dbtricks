Based on ``pgtricks`` project [https://github.com/akaihola/pgtricks]

``pg_dump_splitsort.py`` ``mysqldump_splitsort.py`` are a handy scripts for pre-processing PostgreSQL's and MySQL's
``pg_dump`` and ``mysqldump`` output to make it more suitable for diffing and storing in version
control.

Usage::

    python (pg_dump|mysdlump)_splitsort.py -v -m=max_datachunk_in_bytes -c <filename>.sql
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
    9999_epilogue.sql - everything after the last COPY

for mysql files will be numbered without schema ``00001_<tabley>_00001.sql``

The files for table data are numbered uniquely, and first, order number assign is stored in .pgtricks file.
files can be used to re-create the database::

    $ cat *.sql | psql <database>


    $ cat *.sql | mysql -d <database>
