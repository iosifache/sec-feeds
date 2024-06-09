import sys
import typing
from datetime import date, timedelta

from loguru import logger


def generate_entries(since_date: typing.Optional[str]) -> None:
    if since_date:
        day, month, year = since_date.split(".")
        start_date = date(year, month, day)
    else:
        start_date = date.now()

    end_date = date.now()

    delta = timedelta(days=1)
    while start_date <= end_date:
        print(start_date.strftime("%Y-%m-%d"))
        start_date += delta


def main() -> None:
    since_date = None
    if len(sys.argv) == 2:
        since_date = sys.argv[1]
    elif len(sys.argv) > 2:
        exit(1)

    logger.debug(f"Set date is: {since_date}")

    entries = generate_entries(since_date)


if __name__ == "__main__":
    main()
