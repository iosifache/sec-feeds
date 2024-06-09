# `sec-feeds`

## Description

`sec-feeds` is a repository hosting Atom RSS feeds for different cybersecurity topics in which I'm interested.

All feeds can be accessed at `iosifache.me/feeds/<id>`, where `<id>` is the ID of the feed you want to follow.

## Features

- Hosting of RSS feeds
- RSS feeds generation for the data sources which doesn't have public feeds
- Caching of the content used for the generation of the feeds
- Daily regeneration of the feeds

## Feeds and their data sources

| ID       | Name               | Contained data                                  | Feed location      | Cache location   |
|----------|--------------------|-------------------------------------------------|--------------------|------------------|
| `ubuntu` | Security on Ubuntu | `ubuntu-security` and `ubuntu-meeting` IRC logs | `feeds/ubuntu.xml` | `cache/ubuntu/*` |

