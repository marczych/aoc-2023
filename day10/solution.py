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
    NORTH = Coordinate(x=0, y=-1)
    EAST = Coordinate(x=1, y=0)
    SOUTH = Coordinate(x=0, y=1)
    WEST = Coordinate(x=-1, y=0)


@dataclasses.dataclass(frozen=True)
class Node:
    connections: list[Coordinate]

    @property
    def is_starting_location(self) -> bool:
        return len(self.connections) == 4


NODES = {
    "|": Node(connections=[Connection.NORTH.value, Connection.SOUTH.value]),
    "-": Node(connections=[Connection.EAST.value, Connection.WEST.value]),
    "L": Node(connections=[Connection.NORTH.value, Connection.EAST.value]),
    "J": Node(connections=[Connection.NORTH.value, Connection.WEST.value]),
    "7": Node(connections=[Connection.WEST.value, Connection.SOUTH.value]),
    "F": Node(connections=[Connection.EAST.value, Connection.SOUTH.value]),
    "S": Node(
        connections=[
            Connection.NORTH.value,
            Connection.EAST.value,
            Connection.SOUTH.value,
            Connection.WEST.value,
        ]
    ),
}


def main(input_file: str) -> None:
    graph: dict[Coordinate, Node] = {}
    starting_location: Coordinate | None = None

    # Parse the graph.
    with open(input_file, encoding="utf-8") as input:
        y = 0
        for line in (l.rstrip("\n") for l in input):
            for x in range(len(line)):
                if node := NODES.get(line[x]):
                    graph[Coordinate(x=x, y=y)] = node
                    if line[x] == "S":
                        starting_location = Coordinate(x=x, y=y)
            y += 1

    assert starting_location, "Starting location not found!"

    # Remove nodes that don't connect to other nodes.
    for coordinate in list(graph.keys()):

        def is_valid_connection(offset: Coordinate) -> bool:
            other_coordinate = coordinate + offset
            other_node = graph.get(other_coordinate)
            return bool(other_node) and offset.reverse() in other_node.connections

        node = graph[coordinate]
        if node.is_starting_location:
            valid_connections = [
                connection
                for connection in node.connections
                if is_valid_connection(connection)
            ]
            assert len(valid_connections) == 2, "Invalid starting location!"
            graph[coordinate] = Node(connections=valid_connections)
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

    print(len(main_loop) // 2)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
