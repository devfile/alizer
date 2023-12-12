# Contributing

Thank you for your interest in contributing to Alizer! We welcome your additions to this project.

## Code of Conduct

Before contributing to this repository for the first time, please review our project's [Code of Conduct](https://github.com/devfile/api/blob/main/CODE_OF_CONDUCT.md).

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](DCO) file for details.

## How to contribute:

### Issues

If you spot a problem with devfile alizer, [search if an issue already exists](https://github.com/devfile/api/issues?q=is%3Aissue+is%3Aopen+label%3Aarea%2Falizer). If a related issue doesn't exist, you can open a new issue using a relevant [issue form](https://github.com/devfile/api/issues/new/choose).

You can tag Alizer related issues with the `/area alizer` text in your issue.

### Development

#### Repository Format

The `alizer` repository includes different components:

- [CLI](./README.md#cli)
- [Alizer Library](./README.md#library-package)

As a result, `alizer` can be used both as a cli tool and imported as a package inside other projects. More information for the repository can be found [here](./docs/public/alizer-spec.md).

#### Building locally

More information for building & running locally the project can be found [here](./README.md#usage).

#### Testing

Apart from testing your changes locally with an updated Alizer CLI, someone can test their changes by running `make test`. This will test the updates against all existing test cases.

### Submitting Pull Request

**Note:** All commits must be signed off with the footer:

```
Signed-off-by: First Lastname <email@email.com>
```

You can easily add this footer to your commits by adding `-s` when running `git commit`. When you think the code is ready for review, create a pull request and link the issue associated with it.

Owners of the repository will watch out for and review new PRs.

By default for each change in the PR, GitHub Actions and OpenShift CI will run checks against your changes (linting, unit testing, and integration tests).

If comments have been given in a review, they have to be addressed before merging.

After addressing review comments, don't forget to add a comment in the PR with the reviewer mentioned afterward, so they get notified by Github to provide a re-review.

# Contact us

If you have any questions, please visit us the `#devfile` channel under the [Kubernetes Slack](https://slack.k8s.io) workspace.
