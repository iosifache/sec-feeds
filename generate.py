import sys
import typing
from datetime import datetime, date, timedelta
from dataclasses import dataclass
import requests
import os

from loguru import logger
from feedgen.feed import FeedGenerator


SECURITY_CHANNELS = ["ubuntu-security", "ubuntu-meeting"]
IRC_CHANNEL_URL_FORMAT = "https://irclogs.ubuntu.com/{}/{}/{}/%23{}.txt"
CACHE_FOLDER = "cache/ubuntu"
ATOM_FEED_PATH = "feeds/ubuntu.xml"


@dataclass
class UbuSecRSSEntry:
    name: str
    date: date
    content: str


def generate_entries(
    since_date: typing.Optional[str],
) -> typing.Generator[UbuSecRSSEntry, None, None]:
    yesterday = (datetime.now() - timedelta(1)).date()

    if since_date:
        day, month, year = since_date.split(".")
        start_date = date(int(year), int(month), int(day))
    else:
        start_date = yesterday

    end_date = yesterday

    logger.debug(f"Generating entries since {start_date}")

    delta = timedelta(days=1)
    while start_date <= end_date:
        logger.debug(f"Generating an entry for {start_date.strftime('%Y-%m-%d')}")

        yield generate_entry_for_date(start_date)

        start_date += delta


def generate_entry_for_date(entry_date: date) -> typing.Optional[UbuSecRSSEntry]:
    entry_content = ""
    channels = []
    for channel in SECURITY_CHANNELS:
        header = generate_channel_title(channel)
        content = get_and_cache_channel_content_for_date(channel, entry_date)

        if len(content) > 0:
            entry_content += header + "\n\n" + content + "\n"
            channels.append(channel)

    if len(channels) > 0:
        title = generate_rss_entry_title(entry_date, channels)

        return UbuSecRSSEntry(title, entry_date, entry_content)

    return None


def generate_channel_title(channel_name: str) -> str:
    return "# " + channel_name.upper()


def get_and_cache_channel_content_for_date(
    channel_name: str, desired_date: date
) -> str:
    day = str(desired_date.day).rjust(2, "0")
    month = str(desired_date.month).rjust(2, "0")
    year = str(desired_date.year).rjust(2, "0")

    cache_file = f"{CACHE_FOLDER}/{year}-{month}-{day}-{channel_name}.txt"
    if os.path.isfile(cache_file):
        logger.debug(f"Using cache file: {cache_file}")

        content = open(cache_file, "r").read()
    else:
        url = IRC_CHANNEL_URL_FORMAT.format(year, month, day, channel_name)
        logger.debug(f"Getting content from the URL: {url}")

        content = requests.get(url).text

        open(cache_file, "w").write(content)

    logger.debug(f"Got content for date {desired_date} with length of {len(content)}")

    return content


def generate_rss_entry_title(entry_date: date, channels: typing.List[str]) -> str:
    return entry_date.strftime("%d.%m.%Y") + ": " + ", ".join(channels)


def generate_and_save_rss_feed(
    since_date: typing.Optional[str],
) -> None:
    feed = FeedGenerator()
    feed.id("ubuntu-security-feeds")
    feed.title("Ubuntu Security Feeds")
    feed.logo("https://assets.ubuntu.com/v1/a7e3c509-Canonical%20Ubuntu.svg")
    feed.language("en")

    logger.debug("Adding entries in an RSS feed")
    for entry in generate_entries(since_date):
        if not entry:
            continue

        rss_entry = feed.add_entry()
        rss_entry.id(entry.name)
        rss_entry.title(entry.name)
        rss_entry.description(entry.content)

        logger.debug(f'Added entry with the title: "{entry.name}"')

    feed.atom_file(ATOM_FEED_PATH)


def main() -> None:
    since_date = None
    if len(sys.argv) == 2:
        since_date = sys.argv[1]
    elif len(sys.argv) > 2:
        exit(1)

    logger.debug(f"Set date is: {since_date}")

    generate_and_save_rss_feed(since_date)


if __name__ == "__main__":
    main()
