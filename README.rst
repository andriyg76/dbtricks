``pg_dump_splitsort.py`` is a handy script for pre-processing PostgreSQL's
``pg_dump`` output to make it more suitable for diffing and storing in version
control.

Usage::

    python (pg_dump|mysdlump)_splitsort.py -v --max=max_datachunk_in_bytes <filename>.sql

The script splits the dump into the following files:

Before script run all \d\d\d\d.*.sql files will be dropped from 

| ``0000_prologue.sql``:
    everything up to the first COPY
| ``000X_<schema>_<table>_00001.sql``
| ``000X_<schema>_<table>_00002.sql``
| :
| :
| ``00YY_<schemax>_<tabley>_00001.sql``
    COPY data for each table *sorted by the first field, and second field*
| ``9999_epilogue.sql``:
    everything after the last COPY

The files for table data are numbered uniquely, and order number is stored in .pgtricks file.
files can be used to re-create the database::

    $ cat *.sql | (psql <database>|mysql -d <database>)

I've found that a good way to take a quick look at differences between dumps is
to use the `meld` tool on the whole directory::

    $ meld old-dump/ new-dump/

Storing the dump in version control also gives a decent view on the
differences. Here's how to configure git to use color in diffs::

    # ~/.gitconfig
    [color]
            diff = true
    [color "diff"]
            frag = white blue bold
            meta = white green bold
            commit = white red bold

**Note:** If you have created/dropped/renamed tables, remember to delete all
`.sql` files before post-processing the new dump.
