import argparse
import dataclasses

Position = tuple[int, int]


@dataclasses.dataclass(frozen=True)
class Instruction:
    direction: str
    length: int


DIRECTIONS: dict[str, Position] = {
    "3": (0, -1),
    "0": (1, 0),
    "1": (0, 1),
    "2": (-1, 0),
}


def main(input_file: str) -> None:
    instructions: list[Instruction] = []
    with open(input_file, encoding="utf-8") as input:
        for line in (l.rstrip("\n") for l in input):
            (_, _, raw_color) = line.split(" ")
            color = raw_color[2:-1]
            instructions.append(
                Instruction(direction=color[-1], length=int(color[0:-1], 16))
            )

    xpos = 0
    perimeter = 0
    area = 0
    for instruction in instructions:
        direction = DIRECTIONS[instruction.direction]

        dx = direction[0] * instruction.length
        dy = direction[1] * instruction.length

        # No need to keep track of y-position.
        xpos = xpos + dx

        perimeter += instruction.length
        area += xpos * dy

    print(area + perimeter // 2 + 1)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
