import argparse


def main(input_file: str) -> None:
    with open(input_file, encoding="utf-8") as input:
        for line in (l.rstrip("\n") for l in input):
            print(line)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--file", required=True)
    args = parser.parse_args()
    main(args.file)
