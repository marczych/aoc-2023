import argparse
from collections import defaultdict


def _puzzle_value_helper(puzzle: list[str]) -> int | None:
    pivot_matches: dict[int, int] = defaultdict(int)

    for line in puzzle:
        for pivot in range(1, len(line)):
            size = min(pivot, len(line) - pivot)
            left = line[pivot - size : pivot]
            right = line[pivot + size - 1 : pivot - 1 : -1]

            if left == right:
                pivot_matches[pivot] += 1

    for pivot, count in pivot_matches.items():
        if count == len(puzzle) - 1:
            return pivot

    return None


def rotate_puzzle(puzzle: list[str]) -> list[str]:
    return ["".join(line[i] for line in puzzle) for i in range(len(puzzle[0]))]


def get_puzzle_value(puzzle: list[str]) -> int:
    if value := _puzzle_value_helper(puzzle):
        return value

    if value := _puzzle_value_helper(rotate_puzzle(puzzle)):
        return value * 100

    assert False, "Couldn't find puzzle value!"


def main(input_file: str) -> None:
    with open(input_file, encoding="utf-8") as input:
        total = 0
        puzzle = []
        for line in (l.rstrip("\n") for l in input):
            if line:
                puzzle.append(line)
            else:
                total += get_puzzle_value(puzzle)
                puzzle = []

    total += get_puzzle_value(puzzle)
    print(total)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
