# Change Log

## [0.1.0](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/tree/0.1.0) (2018-02-11)
**Implemented enhancements:**

- Avoid unnecessary lmConf updates [\#31](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/issues/31)
- Launch the /usr/sbin/llng-fastcgi-server process [\#28](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/issues/28)
- Add step-by-step documentation [\#27](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/issues/27)
- Include some optional dependencies in the image [\#43](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/43) ([sathieu](https://github.com/sathieu))
- Process [\#39](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/39) ([sathieu](https://github.com/sathieu))
- Renames [\#37](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/37) ([sathieu](https://github.com/sathieu))
- Tests for pkg/controller [\#24](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/24) ([sathieu](https://github.com/sathieu))
- Coveralls integration [\#18](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/18) ([sathieu](https://github.com/sathieu))
- End-to-end tests [\#7](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/7) ([sathieu](https://github.com/sathieu))

**Fixed bugs:**

- lmConf updates should reload LemonLDAP::NG [\#41](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/issues/41)
- Avoid unnecessary lmConf updates [\#31](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/issues/31)
- Fix dataraces [\#29](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/issues/29)
- Launch the /usr/sbin/llng-fastcgi-server process [\#28](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/issues/28)
- Run flaky tests at most 5 times [\#54](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/54) ([sathieu](https://github.com/sathieu))
- Wait LemonLDAP::NG a bit longer [\#51](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/51) ([sathieu](https://github.com/sathieu))
- Avoid unnecessary lmConf updates [\#42](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/42) ([sathieu](https://github.com/sathieu))
- Ensure all map keys are strings [\#25](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/25) ([sathieu](https://github.com/sathieu))
- .travis.yml: Move "docker push" to script [\#19](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/19) ([sathieu](https://github.com/sathieu))
- Fix segfault on Configuration [\#17](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/17) ([sathieu](https://github.com/sathieu))

**Merged pull requests:**

- Ensure lmConf updates reloads LemonLDAP::NG [\#45](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/45) ([sathieu](https://github.com/sathieu))
- Simplify docker-entrypoint.sh [\#40](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/40) ([sathieu](https://github.com/sathieu))
- Workaround for "ERROR: logging before flag.Parse" [\#38](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/38) ([sathieu](https://github.com/sathieu))
- Missins doc update for "Move to pflags [\#36](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/36) ([sathieu](https://github.com/sathieu))
- Move to pflags \(POSIX/GNU-style --flags\) [\#35](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/35) ([sathieu](https://github.com/sathieu))
- Fix FakeFS data race, an gotlint fixes [\#33](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/33) ([sathieu](https://github.com/sathieu))
- Avoid data race with lemonldapng.config.Config [\#26](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/26) ([sathieu](https://github.com/sathieu))
- Fix for vet complaints [\#23](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/23) ([sathieu](https://github.com/sathieu))
- docker tag and push the image without arch suffix [\#22](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/22) ([sathieu](https://github.com/sathieu))
- Invalid ConfigMap should not be fatal [\#21](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/21) ([sathieu](https://github.com/sathieu))
- doc: Command line flags [\#20](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/20) ([sathieu](https://github.com/sathieu))
- Log when config save is failing [\#16](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/16) ([sathieu](https://github.com/sathieu))
- Add ConfigMap support [\#15](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/15) ([sathieu](https://github.com/sathieu))
- Use simpler ListWatch hooking [\#14](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/14) ([sathieu](https://github.com/sathieu))
- Tests [\#13](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/13) ([sathieu](https://github.com/sathieu))
- Handle locationRules [\#12](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/12) ([sathieu](https://github.com/sathieu))
- Fix lmConf-{num}.js Sprintf [\#11](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/11) ([sathieu](https://github.com/sathieu))
- LemonLDAP::NG configuration handling [\#10](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/10) ([sathieu](https://github.com/sathieu))
- Improve controller logging [\#9](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/9) ([sathieu](https://github.com/sathieu))
- Fix syntax error in docker-entrypoint.sh [\#8](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/8) ([sathieu](https://github.com/sathieu))
- Push images to Docker Hub [\#6](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/6) ([sathieu](https://github.com/sathieu))
- Containers [\#5](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/5) ([sathieu](https://github.com/sathieu))
- Add github templates from ingress-nginx [\#4](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/4) ([sathieu](https://github.com/sathieu))
- Static checks [\#3](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/3) ([sathieu](https://github.com/sathieu))
- Move main to cmd/ [\#2](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/2) ([sathieu](https://github.com/sathieu))
- Initial travis configuration [\#1](https://github.com/lemonldap-ng-controller/lemonldap-ng-controller/pull/1) ([sathieu](https://github.com/sathieu))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*