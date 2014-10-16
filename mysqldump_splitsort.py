# !/usr/bin/env python2.7
import heapq
import json
import os
from pprint import pprint

import re
import tempfile
import importer


def _try_as_float(s):
    if not s or s[0] not in '0123456789.-':
        # optimization
        return s
    try:
        return float(s)
    except ValueError:
        return s


def _clean_start_spaces(_line):
    while _line and (_line[0] == ' ' or _line[0] == '\t'):
        _line = _line[1:]
    return _line


def _get_first_value(line):
    line = _clean_start_spaces(line)

    if not line:
        return None, None
    if line[0] == '\'' or line[0] == '\"':

        last = 1
        while True:
            pos = line.find(line[0], last)
            if pos < 0:
                raise ValueError('Can\'t split')
            if pos == 1:
                head, separator, tail = line[2:].rpartition(',')
                return '', _clean_start_spaces(tail)
            elif line[pos - 1] == '\\':
                last = pos + 1
                continue
            elif pos + 1 < len(line) and line[pos + 1] == line[0]:
                last = pos + 2
                continue
            else:
                value = line[1:pos - 1]
                head, separator, tail = line[pos + 1:].rpartition(',')
                return value, _clean_start_spaces(tail)
        # finding closing quote
    else:
        head, separator, tail = line.rpartition(',')
        if not head and not separator:
            return tail, None
        else:
            return head, _clean_start_spaces(tail)
# def _get_first_value


def _lines_compare(l1, l2):
    p1, tail1 = _get_first_value(l1)
    p2, tail2 = _get_first_value(l2)
    result = cmp(_try_as_float(p1), _try_as_float(p2))
    if not result and tail1 and tail2:
        return _lines_compare(tail1, tail2)
    return result


def key(line):
    key, tail = _get_first_value(line)
    if tail:
        return _try_as_float(key), _try_as_float(_get_first_value(tail))
    else:
        return _try_as_float(key), None


class _Dumper:
    def __init__(self):
        self._output = None
        self._buf = []

    def flush(self):
        if not self._output:
            raise ValueError("Output is no opened during flush")
        self._output.writelines(self._buf)
        self._output.flush()
        self._buf = []

    def new_output(self, path):
        if self._output:
            self.flush()
            self._output.close()
        self._output = file(path, 'w')

    def pop_last_lines(self, lines=1):
        accum = []
        if lines > 0:
            while lines > 0 and len(self._buf) > 0:
                accum.append(self._buf.pop())
                lines -= 1
            accum.reverse()
        return accum

    def add_lines(self, lines):
        self._buf.extend(lines)

    def append(self, line):
        self._buf.append(line)

    def __enter__(self):
        return self

    def __exit__(self, exc_type, exc_val, exc_tb):
        if self._output:
            self.flush()
            self._output.close()
            self._output = None


# class _Dumper


class _DataHandler:
    def __init__(self, max_chunk_size, start_line, table_name, counter):
        self._max_chunk_size = max_chunk_size

        self._start_line = start_line
        self._table_name = table_name
        self._counter = counter

        self._buf = []
        self._buf_size = 0
        self._chunks = []
        self._flushed = False

    def add_line(self, line):
        if self._flushed:
            raise ValueError("Can't add table data, table interaction already flushed")
        self._buf.append(line)
        self._buf_size += len(line)
        if self._buf_size > self._max_chunk_size:
            importer.verbose("Splitting %s temporary data, %d-th part, %d lines, size %d" %
                             (self._table_name, len(self._chunks) + 1, len(self._buf), self._buf_size))

            self._buf.sort(cmp=_lines_compare)
            chunk = tempfile.TemporaryFile()
            chunk.writelines(self._buf)
            chunk.seek(0)

            self._chunks.append(chunk)
            self._buf_size = 0
            self._buf = []

    # add_line

    def flush_data(self, dumper):
        if self._flushed:
            raise ValueError("Can't add table data, table interaction already flushed")
        importer.verbose("Storing %s data, %d lines in memory, size %d" %
                         (self._table_name, len(self._buf), self._buf_size))
        self._buf.sort(cmp=_lines_compare)

        if len(self._chunks):

            # multiple chunks
            sequence = 1
            dumper.append(self._start_line)
            output_size = 0
            sorted_memory = [(key(line), line) for line in self._buf]
            sorted_chunks = [[(key(line), line) for line in chunk]
                             for chunk in self._chunks]
            for _key, _line in heapq.merge(sorted_memory, *sorted_chunks):
                dumper.append(_line + '\n')
                output_size += len(_line)
                if output_size > self._max_chunk_size:
                    dumper.flush()

                    dumper.new_output('{counter:04}_{table_name}_{sequence:05}.sql'.format(
                        counter=self._counter, table_name=self._table_name, sequence=sequence))
                    output_size = 0
                    sequence += 1

                    dumper.append(self._start_line)
                    # for _key _ine in sorted data
        else:
            dumper.append(self._start_line)
            for _line in self._buf:
                dumper.append(_line + '\n')
            dumper.flush()

        dumper.flush()

        for chunk in self._chunks:
            chunk.close()
        self._chunks = []
        self._buf_size = 0
        self._buf = []

        self._flushed = True
        # flush_data

# class _DataHandler


TABLE_STRUCTURE_RE = re.compile(r'^-- Table structure for table `(?P<table>.*?)`')
INSERT_INTO_RE = re.compile(r'^INSERT INTO .* VALUES$')
VALUES_RE = re.compile(r'^\((?P<data>.*?)\)[;,]$')


def __do_split(args, sql_dump_file, order):
    with _Dumper() as dumper:
        counter = 0
        previous_table = None
        dumper.new_output('0000_prologue.sql')
        table_name = None
        epilogue = False
        data_handler = None

        for line in sql_dump_file:
            if epilogue:
                dumper.append(line)
            if TABLE_STRUCTURE_RE.match(line):
                table_name = TABLE_STRUCTURE_RE.match(line).groupdict()['table']
                counter = importer.get_order_number(order, table_name, previous_table)

                backup = dumper.pop_last_lines(2)
                dumper.flush()
                dumper.new_output('{counter:04}_{table_name}.sql'.format(
                    counter=counter, table_name=table_name))
                dumper.add_lines(backup)
                dumper.append(line)
            elif INSERT_INTO_RE.match(line):
                if not data_handler:
                    data_handler = _DataHandler(args.chunk_size, line, table_name, counter)
            elif VALUES_RE.match(line):
                data_handler.add_line(VALUES_RE.match(line).groupdict()['data'])
            elif data_handler and line == '\n':
                pass
            elif data_handler:
                data_handler.flush_data(dumper)
                data_handler = None
            else:
                dumper.append(line)
        #foreach line

        if data_handler:
            data_handler.flush_data(dumper)
            data_handler = None

        dumper.flush()


# def __do_split


def split_sql_file(args):
    importer.verbose("Opening sql dump file:", args.sql_dump_file)
    with file(args.sql_dump_file) as sql_dump_file:
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
            with open('.order', 'r') as infile:
                order = json.load(infile)
                importer.verbose("Read .order file", order)
        else:
            importer.verbose("Can't find .order file, starting from scratch")
            order = {}
            # if .order exist

        __do_split(args, sql_dump_file, order)

        importer.verbose("Writing updated .order file")
        with open('.order', 'w') as outfile:
            json.dump(order, fp=outfile, sort_keys=True, indent=4, separators=(',', ': '))
            # with sql_dump_file

# def split_sql_file

if __name__ == '__main__':
    split_sql_file(importer.create_argsparser())
