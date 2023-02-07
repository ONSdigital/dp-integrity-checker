# dp-integrity-checker

Periodic nomad job for checking integrity of zebedee workspace

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable | Default   | Description                 |
|----------------------|-----------|-----------------------------|
| ZEBEDEE_ROOT         | "content" | Root of the zebedee-content |

NB. For developers the zebedee root is usually specified in the lowercase `zebedee_root` env so this service aliases
this in the `make debug` target to make local development more straightforward.

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2023, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

