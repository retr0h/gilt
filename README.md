[![release](https://img.shields.io/github/release/retr0h/gilt.svg?style=for-the-badge)](https://github.com/retr0h/gilt/releases/latest)
[![pypi](https://img.shields.io/pypi/v/python-gilt?style=for-the-badge)](https://pypi.org/project/python-gilt/)
[![codecov](https://img.shields.io/codecov/c/github/retr0h/gilt?token=clAMnFQCEQ&style=for-the-badge)](https://codecov.io/gh/retr0h/gilt)
[![go report card](https://goreportcard.com/badge/github.com/retr0h/gilt?style=for-the-badge)](https://goreportcard.com/report/github.com/retr0h/gilt/v2)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE)
[![build](https://img.shields.io/github/actions/workflow/status/retr0h/gilt/go.yml?style=for-the-badge)](https://github.com/retr0h/gilt/actions/workflows/go.yml)
[![powered by](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
[![conventional commits](https://img.shields.io/badge/Conventional%20Commits-1.0.0-yellow.svg?style=for-the-badge)](https://conventionalcommits.org)
![gitHub commit activity](https://img.shields.io/github/commit-activity/m/retr0h/gilt?style=for-the-badge)

# Gilt

<img src="https://github.com/retr0h/gilt/raw/main/asset/gilt.png" align="left" width="250px" height="250px" />

Gilt is a tool which aims to make repo management, manageable.  Gilt
clones repositories at a particular version, then overlays the repository to
the provided destination.  An alternate approach to "vendoring".

What makes Gilt interesting, is the ability to overlay particular files and/or
directories from the specified repository to given destinations. Originally,
this was quite helpful for those using Ansible, since libraries, plugins, and
playbooks are often shared, but Ansible's [Galaxy][] has no mechanism to handle
this.  Currently, this is proving useful for overlaying [Helm charts][].

<br clear="left"/>

## Documentation

[Installation][] | [Usage][] | [Documentation][]

[Installation]: https://retr0h.github.io/gilt/installation
[Usage]: https://retr0h.github.io/gilt/usage
[Documentation]: https://retr0h.github.io/gilt/

## License

The [MIT][] License.

The logo is licensed under the [Creative Commons NoDerivatives 4.0 License][],
and designed by [@nanotron][].
If you have some other use in mind, contact us.

[Galaxy]: https://docs.ansible.com/ansible/latest/reference_appendices/galaxy.html
[Helm charts]: https://helm.sh/docs/topics/charts/
[MIT]: LICENSE
[Creative Commons NoDerivatives 4.0 License]: https://creativecommons.org/licenses/by-nd/4.0/
[@nanotron]: https://github.com/nanotron
