import argparse
import dataclasses


@dataclasses.dataclass(frozen=True)
class Node:
    left: str
    right: str


def main(input_file: str) -> None:
    instructions: list[str] = []
    map: dict[str, Node] = {}

    with open(input_file, encoding="utf-8") as input:
        for line in (l.rstrip("\n") for l in input):
            if not instructions:
                instructions = [*line]
            elif line:
                (node, _, left, right) = line.split(" ")
                # Remove extraneous parens and commas.
                map[node] = Node(left=left[1:-1], right=right[:-1])

    instruction_count = 0
    current_node = "AAA"
    while current_node != "ZZZ":
        if instructions[instruction_count % len(instructions)] == "L":
            current_node = map[current_node].left
        else:
            current_node = map[current_node].right
        instruction_count += 1

    print(instruction_count)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
