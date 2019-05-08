[![CircleCI](https://circleci.com/gh/giantswarm/lighthouse-keeper.svg?style=shield&circle-token=f1d2933a2c4afd35322007b709700216adbd89f1)](https://circleci.com/gh/giantswarm/lighthouse-keeper)

# lighthouse-keeper

Utility to run [lighthouse](https://github.com/GoogleChrome/lighthouse) in a CI environment,
especially against Docker containers, and report back the results into a Pull Request.

## Features

- Create lighthouse audit reports, running lighthouse and Chrome headless in a Docker container
- Print report results in a readable form
- Compare two reports, optionally comment the result into a GitHub pull request

## Installation

```
go get github.com/giantswarm/lighthouse-keeper
go install github.com/giantswarm/lighthouse-keeper
```

## Usage

Check `lighthouse-keeper --help` for details.

### `audit` - Create a lighthouse report

This will create a lighthouse report for `https://example.com/` in `mysite.json`:

```
lighthouse-keeper audit --name mysite --url https://example.com/
```

To audit a site running in a docker container named `container`, use `--docker-link`:

```
lighthouse-keeper audit --url http://container:8000/ --docker-link container:container
```

More flags:

- Use `--form-factor mobile` to emulate a mobile device form factor. Default is `desktop`.
- Use `--ignore-certificate-errors` to check against an HTTPS site using a self-signed or otherwise bad certificate.

Check `lighthouse-keeper audit --help` for details.

### `view` - Pretty-print a report

This will print the complete report results:

```
lighthouse-keeper --input ./report.json
```

To only print the rows with scores below 100, use this:

```
lighthouse-keeper --input ./report.json --omit-done
```

### `compare` - Compare two lighthouse reports

This prints the differences between two lighthouse reports:

```
lighthouse-keeper compare \
  --input before.json --inputlabel before \
  --input after.json --inputlabel after
```

The output looks somewhat like this:

```
+-----------------------------+--------+-------+-------+
|                             | BEFORE | AFTER | DELTA |
+-----------------------------+--------+-------+-------+
| Performance                 |     64 |    66 | +2    |
| - First Contentful Paint    |     36 |    40 | +4    |
| - Speed Index               |     70 |    74 | +4    |
| - Time to Interactive       |     77 |    78 | +1    |
| - Estimated Input Latency   |     98 |   100 | +2    |
| - JavaScript execution time |     93 |    94 | +1    |
+-----------------------------+--------+-------+-------+
```

Here is an example how the result would be commented into a GitHub pull request:

```
lighthouse-keeper compare \
  --input before.json --inputlabel before \
  --input after.json --inputlabel after \
  --github-owner myorg \
  --github-repo myrepo \
  --github-issue 11 \
  --github-token $GITHUB_PERSONAL_ACCESS_TOKEN
  ```

Preview:

---

Comparison of lighthouse reports:

|                             | BEFORE | AFTER |  DELTA  |
|-----------------------------|--------|-------|---------|
| **Performance**             |     64 |    66 | ✅  +2  |
| - First Contentful Paint    |     36 |    40 | ✅   +4 |
| - Speed Index               |     70 |    74 | ✅   +4 |
| - Time to Interactive       |     77 |    78 | ✅   +1 |
| - Estimated Input Latency   |     98 |   100 | ✅   +2 |
| - JavaScript execution time |     93 |    94 | ✅   +1 |

## Misc

`lighthouse-keeper` requires Docker to be installed. It executes the image

    quay.io/giantswarm/lighthouse

to run Chrome and lighthouse independent of the environment. The [image](https://quay.io/repository/giantswarm/lighthouse?tag=latest&tab=tags) is built
automatically based on
[this Dockerfile](https://github.com/giantswarm/lighthouse/blob/master/Dockerfile).

Thanks for the inspiration to:

- https://github.com/andreasonny83/lighthouse-ci
- https://github.com/carlesnunez/lighthouse-gh-reporter
