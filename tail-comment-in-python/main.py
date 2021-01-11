import argparse
import queue


def tail(filename, n):
    q = queue.Queue()
    size = 0

    with open(filename) as fh:
        for line in fh:
            q.put(line.strip())
            if size >= n:
                q.get()
            else:
                size += 1      

    for i in range(size):
        print(q.get())

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Print last n lines from a file.')
    parser.add_argument('file', type=str, help='File to read from')
    parser.add_argument('-n', type=int, default=10, help='The last n lines to be printed')
    args = parser.parse_args()

    tail(args.file, args.n)