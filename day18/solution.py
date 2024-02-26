import argparse

Position = tuple[int, int]

DIRECTIONS: dict[str, Position] = dict(
    U=(0, -1),
    R=(1, 0),
    D=(0, 1),
    L=(-1, 0),
)


def grid_debug_print(grid: set[Position], start: int, end: int) -> None:
    for y in range(start, end + 1):
        for x in range(start, end + 1):
            print("#" if (x, y) in grid else ".", end="")
        print()


def main(input_file: str) -> None:
    grid: set[Position] = {(0, 0)}
    xvel = 0
    yvel = 0

    xpos = 0
    ypos = 0

    xmin = 0
    xmax = 0
    ymin = 0
    ymax = 0

    with open(input_file, encoding="utf-8") as input:
        for line in (l.rstrip("\n") for l in input):
            (direction, length_string, _) = line.split(" ")
            length = int(length_string)
            xvel, yvel = DIRECTIONS[direction]
            for _ in range(1, length + 1):
                xpos += xvel
                ypos += yvel
                xmin = min(xmin, xpos)
                xmax = max(xmax, xpos)
                ymin = min(ymin, ypos)
                ymax = max(ymax, ypos)

                grid.add((xpos, ypos))

    xpos = min(xmin, ymin)
    ypos = xpos
    while True:
        if (xpos, ypos) in grid:
            # We found an edge.

            """
            ##
             #
            """
            if (xpos - 1, ypos) in grid and (xpos, ypos + 1) in grid:
                assert False, "Can't handle corner type 1!"
            """
            #
            ##
            """
            if (xpos + 1, ypos) in grid and (xpos, ypos - 1) in grid:
                assert False, "Can't handle corner type 2!"

            xpos += 1
            ypos += 1
            break

        xpos += 1
        ypos += 1

        if xpos > xmax or ypos > ymax:
            assert False, "Couldn't find inside to start filling!"

    # In fill.
    positions_to_fill: set[Position] = {(xpos, ypos)}
    while positions_to_fill:
        xpos, ypos = positions_to_fill.pop()
        for xoffset in range(-1, 2):
            for yoffset in range(-1, 2):
                position = (xpos + xoffset, ypos + yoffset)
                if position not in grid:
                    grid.add(position)
                    positions_to_fill.add(position)

    # grid_debug_print(grid, min(xmin, ymin), max(xmax, ymax))
    print(len(grid))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
