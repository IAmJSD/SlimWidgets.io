# [SlimWidgets.io](https://slimwidgets.io)

A fast, slim, safe, open-source and customisable widget for Discord. Here's why it is all those things:

- **Fast** - We use Go for all the backend which gets turned into Assembly. This means that each component in this application is extremely optimised. Additionally, we call NATS which is a very fast event handler to distribute events.
- **Slim** - This widget takes up less space than the Discord widget does in its default configuration.
- **Safe** - The Discord widget makes the invite code easily accessable. This means that any scrapers that visit the widget can easily get a invite to your guild which is a common method of data mining. SafeWidgets.io allows you to turn off invites on the widget, and if you have them on, users have to go through a CAPTCHA to get a one time invite, preventing scraping.
- **Open-source** - Every part of this service is open source in this repository!
- **Customisable** - Want invites on? You can turn them on. Want a description? You can set a description. The choice is yours.

## Setup
All of our Kubernetes files are provided. Note that for this service to work, you will need both RethinkDB and NATS installed on your cluster. We suggest using Helm to do this.

You will also need to make a RethinkDB database called `slimwidgets`. In this database, you can simply make a table called `guilds`. Setup complete!
