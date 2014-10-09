#!/usr/bin/env python2.7
import heapq
import json
import os
from pprint import pprint

import re
import tempfile
import importer


COPY_RE = re.compile(r'COPY .*? \(.*?\) FROM stdin;\n$')


def try_float(s):
    if not s or s[0] not in '0123456789.-':
        # optimization
        return s
    try:
        return float(s)
    except ValueError:
        return s


def lines_compare(l1, l2):
    p1 = l1.split('\t', 1)
    p2 = l2.split('\t', 1)
    result = cmp(try_float(p1[0]), try_float(p2[0]))
    if not result and len(p1) > 1 and len(p2) > 1:
        return lines_compare(p1[1], p2[1])
    return result


DATA_COMMENT_RE = r'-- Data for Name: (?P<table>.*?); Type: TABLE DATA; Schema: (?P<schema>.*?);'
END_COPY_LINE = '\\.\n'

def key(line):
    keys = line.split('\t')
    if len(keys) >= 2:
        return try_float(keys[0]), try_float(keys[1])
    else:
        return try_float(keys[0]), None


def split_sql_file(args):
    sql_file_path = args.sql_dump_file

    if args.destination_path:
        importer.verbose("Changing dir to:", args.destination_path)
        os.chdir(args.destination_path)
    # if -d path

    if args.c:
        importer.verbose("Removing previous sql chunks...")
        sql_file_part_re = re.compile("^\d+_.*[sS][qQ][lL]$")
        for f in os.listdir("."):
            if sql_file_part_re.match(f):
                importer.verbose("Removing file", f)
                os.remove(f)
                # for f
    # if -c

    if os.path.exists(".order"):
        order = json.loads(file(".order").read())
        importer.verbose("Read .order file", order)
    else:
        importer.verbose("Can't find .order file, starting from scratch")
        order = {}
    # if .order exist

    output = None
    buf = []

    def flush():
        output.writelines(buf)
        buf[:] = []

    def new_output(path):
        if output:
            output.close()
        return file(path, 'w')

    copy_lines = None
    copy_line = None
    copy_size = 0
    counter = 0
    output = new_output('0000_prologue.sql')

    for line in file(sql_file_path):
        if copy_lines is None:
            if line in ('\n', '--\n'):
                buf.append(line)
            elif line.startswith('SET search_path = '):
                flush()
                buf.append(line)
            else:
                matcher = re.match(DATA_COMMENT_RE, line)
                if matcher:
                    counter += 1
                    schema = matcher.groupdict()['schema']
                    table = matcher.groupdict()['table']
                    output = new_output(
                        '{counter:04}_{schema}.{table}.sql'.format(
                            counter=counter,
                            schema=schema,
                            table=table))
                    buf.append(line)
                elif COPY_RE.match(line):
                    copy_line = line
                    copy_lines = []
                    chunks = []
                    copy_size = 0
                elif 1 <= counter < 9999:
                    counter = 9999
                    output = new_output('%04d_epilogue.sql' % counter)
                    buf.append(line)
                flush()
        else:  # if copy_lines
            if line == END_COPY_LINE:
                copy_lines.sort(cmp=lines_compare)
                if not len(chunks):
                    importer.verbose("Storing %d lines directly to %s.%s table output" %
                                     (len(copy_lines), schema, table))
                    buf.append(copy_line)
                    buf.extend(copy_lines)
                    buf.append(END_COPY_LINE)
                    flush()
                else:
                    if len(copy_lines):
                        importer.verbose("Storing %d lines to last, %d-th chunk" %
                                         (len(copy_lines), len(chunks) + 1))
                        temp_file = tempfile.TemporaryFile("r+w")
                        temp_file.writelines(copy_lines)
                        temp_file.seek(0)
                        chunks.append(temp_file)

                    importer.verbose("Merging %d chunks of table %s.%s" %
                                     (len(chunks), schema, table))

                    buf.append(copy_line)

                    output_size = 0
                    sequence = 0
                    for _key, _line in heapq.merge(*[
                        [(key(line), line) for line in chunk]
                        for chunk in chunks]):
                        buf.append(_line)
                        output_size += len(_line)
                        if output_size > args.chunk_size:
                            buf.append(END_COPY_LINE)
                            flush()
                            output = new_output(
                                '{counter:04}_{schema}.{table}_{sequence:04}.sql'.format(
                                    counter=counter,
                                    schema=schema,
                                    table=table,
                                    sequence=sequence))

                            sequence += 1
                            output_size = 0
                            buf.append(copy_line)
                            # if chunk should be flushed
                    # for key, line in all_lines
                    buf.append(END_COPY_LINE)
                    flush()

                # if not chunks

                copy_lines = None
                copy_size = 0
                copy_line = None
                chunks = []
            else:
                copy_lines.append(line)
                copy_size += len(line)
                if copy_size > args.chunk_size:
                    importer.verbose("Storing %d lines to %d-th chunk" %
                                     (len(copy_lines), len(chunks) + 1))
                    temp_file = tempfile.TemporaryFile("r+w")
                    copy_lines.sort(cmp=lines_compare)
                    temp_file.writelines(copy_lines)
                    temp_file.seek(0)
                    chunks.append(temp_file)
                    copy_size = 0
                    copy_lines = []
            # if line = END_COPY_LINE (end table data)
        # copy_lines is none
    # for every line

# split_sql_file()

if __name__ == '__main__':
    split_sql_file(importer.create_argsparser())
