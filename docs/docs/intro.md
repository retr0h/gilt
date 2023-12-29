---
slug: /
sidebar_position: 1
title: Home
---

# Gilt

<img src="img/gilt.png" align="left" width="250px" height="250px" />

Gilt is a tool which aims to make repo management, manageable.  Gilt
clones repositories at a particular version, then overlays the repository to
the provided destination.  An alternate approach to "vendoring".

What makes Gilt interesting, is the ability to overlay particular files and/or
directories from the specified repository to given destinations. Originally,
this was quite helpful for those using Ansible, since libraries, plugins, and
playbooks are often shared, but Ansible's [Galaxy][] has no mechanism to handle
this.  Currently, this is proving useful for overlaying [Helm charts].

<br clear="left"/>

## Alternatives

* [Repo][]
* [Git submodules][]
* [Git subtree][]
* [Gilt][]

## History

This project is a golang port of [Gilt][], and aims to correct poor decisions
made in the python version, primarially around config syntax, portability,
and reproducibility.

[Galaxy]: https://docs.ansible.com/ansible/latest/reference_appendices/galaxy.html
[Helm charts]: https://helm.sh/docs/topics/charts/
[Repo]: https://gerrit.googlesource.com/git-repo/+/refs/heads/master/README.md
[Git submodules]: https://git-scm.com/book/en/v2/Git-Tools-Submodules
[Git subtree]: https://github.com/git/git/blob/master/contrib/subtree/git-subtree.txt
[Gilt]: http://gilt.readthedocs.io/en/latest/
