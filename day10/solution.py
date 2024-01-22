from enum import Enum
import argparse
import dataclasses


@dataclasses.dataclass(frozen=True)
class Coordinate:
    x: int
    y: int

    def __add__(self, other: "Coordinate") -> "Coordinate":
        return Coordinate(x=self.x + other.x, y=self.y + other.y)

    def __lt__(self, other: "Coordinate") -> bool:
        if self.x < other.x:
            return True

        if self.x == other.x:
            return self.y < other.y

        return False

    def reverse(self) -> "Coordinate":
        return Coordinate(-self.x, -self.y)


class Connection(Enum):
    WEST = Coordinate(x=-1, y=0)
    NORTH = Coordinate(x=0, y=-1)
    SOUTH = Coordinate(x=0, y=1)
    EAST = Coordinate(x=1, y=0)


@dataclasses.dataclass(frozen=True)
class Node:
    connections: list[Coordinate]
    is_closed_corner: bool = False

    def __post_init__(self) -> None:
        assert self.connections == sorted(self.connections), "Connections not sorted!"

    @property
    def is_starting_location(self) -> bool:
        return len(self.connections) == 4


NODES = {
    "|": Node(connections=[Connection.NORTH.value, Connection.SOUTH.value]),
    "-": Node(connections=[Connection.WEST.value, Connection.EAST.value]),
    "L": Node(
        connections=[Connection.NORTH.value, Connection.EAST.value],
        is_closed_corner=True,
    ),
    "J": Node(connections=[Connection.WEST.value, Connection.NORTH.value]),
    "7": Node(
        connections=[Connection.WEST.value, Connection.SOUTH.value],
        is_closed_corner=True,
    ),
    "F": Node(connections=[Connection.SOUTH.value, Connection.EAST.value]),
    "S": Node(
        connections=[
            Connection.WEST.value,
            Connection.NORTH.value,
            Connection.SOUTH.value,
            Connection.EAST.value,
        ]
    ),
}


def main(input_file: str) -> None:
    graph: dict[Coordinate, Node] = {}
    starting_location: Coordinate | None = None
    size_x = 0
    size_y = 0

    # Parse the graph.
    with open(input_file, encoding="utf-8") as input:
        for line in (l.rstrip("\n") for l in input):
            size_x = len(line)
            for x in range(len(line)):
                if node := NODES.get(line[x]):
                    graph[Coordinate(x=x, y=size_y)] = node
                    if line[x] == "S":
                        starting_location = Coordinate(x=x, y=size_y)
            size_y += 1

    assert starting_location, "Starting location not found!"

    # Remove nodes that don't connect to other nodes.
    for coordinate in list(graph.keys()):

        def is_valid_connection(offset: Coordinate) -> bool:
            other_coordinate = coordinate + offset
            other_node = graph.get(other_coordinate)
            return bool(other_node) and offset.reverse() in other_node.connections

        node = graph[coordinate]
        if node.is_starting_location:
            valid_connections = sorted(
                [
                    connection
                    for connection in node.connections
                    if is_valid_connection(connection)
                ]
            )
            assert len(valid_connections) == 2, "Invalid starting location!"
            graph[coordinate] = Node(
                connections=valid_connections,
                is_closed_corner=any(
                    node.is_closed_corner and valid_connections == node.connections
                    for node in NODES.values()
                ),
            )
        elif not all(is_valid_connection(offset) for offset in node.connections):
            del graph[coordinate]

    # Continue removing nodes without valid neighbors until we don't remove any
    # more.  We don't need to check if the neighbors also connect with us
    # because that was checked up above.
    removed_node_count = 1
    while removed_node_count != 0:
        removed_node_count = 0
        for coordinate in list(graph.keys()):
            node = graph[coordinate]
            if not (
                coordinate == starting_location
                or all((coordinate + offset) in graph for offset in node.connections)
            ):
                del graph[coordinate]
                removed_node_count += 1

    current_node = starting_location
    main_loop: dict[Coordinate, Node] = {}
    while current_node != starting_location or not main_loop:
        for connection in graph[current_node].connections:
            next_coordinate = current_node + connection
            if next_coordinate not in main_loop and (
                next_coordinate != starting_location or len(main_loop) > 1
            ):
                current_node = next_coordinate
                main_loop[current_node] = graph[current_node]
                break
        else:
            assert "Failed to find next node!"

    # Use a raytracing method to count the number of intersections with the
    # loop to determine how many tiles are inside the loop. Because we are
    # raytracing "down and to the right", ignore "L" and "7" as they are
    # "closed corners".
    enclosed_tile_count = 0
    for x in range(size_x):
        for y in range(size_y):
            if Coordinate(x=x, y=y) in main_loop:
                continue
            intersection_count = 0
            offset = 0

            while x + offset < size_x and y + offset < size_y:
                offset += 1
                node = main_loop.get(Coordinate(x=x + offset, y=y + offset))
                # Ignore closed corners because the rays go in and out.
                if node and not node.is_closed_corner:
                    intersection_count += 1
            if intersection_count % 2 == 1:
                enclosed_tile_count += 1

    print(enclosed_tile_count)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
