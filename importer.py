import argparse
from string import join
import sys

__author__ = 'andriy'

__verbose = False


def verbose(*objs):
    if __verbose:
        sys.stderr.write(" ".join([str(i) for i in objs]) + '\n')
        sys.stderr.flush()


def create_argsparser():
    parser = argparse.ArgumentParser(description='Split database dump file to a chunks.')
    parser.add_argument('-m', dest='chunk_size_kb',
                        default=2 * 1024,
                        type=int,
                        help='Max chunk size of database part, in kb default 2014(2Mb)')
    parser.add_argument('-d', dest='destination_path', help='Path, where to store splitted files')
    parser.add_argument('-v', action="store_true",
                        help='Verbose dumping output')
    parser.add_argument('-c', action="store_true", help='Clean destination')
    parser.add_argument('sql_dump_file')
    args = parser.parse_args()
    args.chunk_size = args.chunk_size_kb * 1024

    global __verbose
    __verbose = args.v

    return args
