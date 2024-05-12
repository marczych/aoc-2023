import argparse
import collections


def main(input_file: str) -> None:
    graph: dict[str, set[str]] = collections.defaultdict(set)

    with open(input_file, encoding="utf-8") as input:
        for line in (l.rstrip("\n") for l in input):
            split_line = line.split(": ")
            initial_node = split_line[0]
            connected_nodes = split_line[1].split(" ")
            for node in connected_nodes:
                graph[initial_node].add(node)
                graph[node].add(initial_node)

    new_set: set[str] = set(graph)

    def get_connections_to_old_graph(node: str) -> int:
        return len(graph[node] - new_set)

    def done() -> bool:
        connections_to_old_graph = 0
        for node in new_set:
            connections_to_old_graph += get_connections_to_old_graph(node)
            if connections_to_old_graph > 3:
                break

        return connections_to_old_graph == 3

    # It's possible that an invalid node is selected first which makes this
    # algorithm fail because it moves nodes from across the minimum cut into
    # the new set.
    while not done():
        new_set.remove(max(new_set, key=get_connections_to_old_graph))

    print(len(new_set) * (len(graph) - len(new_set)))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
