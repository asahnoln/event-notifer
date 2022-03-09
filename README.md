![example workflow](https://github.com/asahnoln/event-notifier/actions/workflows/go.yml/badge.svg)
[![codecov](https://codecov.io/gh/asahnoln/event-notifier/branch/master/graph/badge.svg?token=KHGEDM6X34)](https://codecov.io/gh/asahnoln/event-notifier)

# Calendar Event Notifier

## What is it?

Calendar Event Notifier collects events from google calendar and sends them somewhere else in a text format.

## Sends events where?

Currently Telegram Bot and Discord Webhook senders are implemented.

## What text format?

Just usual text like: `"Musical Rehearsal" is going to happen in "Public Theatre" at 2022.03.22 19:00. Arthur, Alina and Danila are needed!`

For the moment, the text is hard-coded.

## Why?

Even when we use a calendar service such as Google Calendar, people might miss notifications sent to their email or phone. A simple reminder in a group chat saves the day. And, of course, it can be automated. A bot reminder.
