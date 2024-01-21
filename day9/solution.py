import argparse


def get_next_value(values: list[int]) -> int:
    total = 0

    for end_position in range(len(values) - 1, 0, -1):
        # Don't touch the last value.
        for position in range(0, end_position):
            values[position] = values[position + 1] - values[position]

        total += values[end_position]

    return total


def main(input_file: str) -> None:
    total = 0

    with open(input_file, encoding="utf-8") as input:
        for line in (l.rstrip("\n") for l in input):
            values = [int(x) for x in line.split(" ")]
            total += get_next_value(values)

    print(total)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
