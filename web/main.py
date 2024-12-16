import argparse
import logging

from src.show import show


def main():
    parser = argparse.ArgumentParser(description='SZTUEA')
    parser.add_argument('--debug', action='store_true')

    args = parser.parse_args()
    debug = args.debug
    logging.info(f"Debug mode is {'on' if debug else 'off'}")

    if debug:
        logging.basicConfig(level=logging.DEBUG)
    show(debug)


if __name__ == "__main__":
    main()
