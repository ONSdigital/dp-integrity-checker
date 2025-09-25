# dp-integrity-checker

Periodic nomad job for checking integrity of zebedee workspace

## Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Tools

To run some of our tests you will need additional tooling:

#### Audit

We use `dis-vulncheck` to do auditing, which you will [need to install](https://github.com/ONSdigital/dis-vulncheck).

### Configuration

| Environment variable          | Default             | Description                                                |
|-------------------------------|---------------------|------------------------------------------------------------|
| ZEBEDEE_ROOT                  | "content"           | Root of the zebedee-content                                |
| CHECK_PUBLISHED_PREVIOUS_DAYS | 1                   | Number of previous days to check published collections for |
| SLACK_ENABLED                 | false               | Whether to send a slack message on failed checks           |
| SLACK_API_TOKEN               | ""                  | A valid slack api token (suppressed from logs)             |
| SLACK_USER_NAME               | "Integrity Checker" | User name to be used for slack messages                    |
| SLACK_ALARM_CHANNEL           | "#sandbox-alarm"    | Slack channel to send alarm messages to                    |
| SLACK_ALARM_EMOJI             | ":rotating_light:"  | Emoji to use for alarm messages                            |

A valid Slack token with `chat:write` and `chat:write.customize` permissions is required if Slack notification is to be
enabled.

NB. For developers the zebedee root is usually specified in the lowercase `zebedee_root` env so this service aliases
this in the `make debug` target to make local development more straightforward.

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

Copyright © 2023, Office for National Statistics <https://www.ons.gov.uk>

Released under MIT license, see [LICENSE](LICENSE.md) for details.
