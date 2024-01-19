import math
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
    periods: list[int] = []
    current_nodes = [node for node in map if node.endswith("A")]
    while current_nodes:
        instruction = instructions[instruction_count % len(instructions)]
        instruction_count += 1

        i = 0
        while i < len(current_nodes):
            if instruction == "L":
                current_nodes[i] = map[current_nodes[i]].left
            else:
                current_nodes[i] = map[current_nodes[i]].right

            if current_nodes[i].endswith("Z"):
                periods.append(instruction_count)
                del current_nodes[i]
            else:
                i += 1

    print(math.lcm(*periods))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
