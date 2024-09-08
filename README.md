# uptime

Simple tool to monitor multiple URLs

## TODO

- [ ] improve create error handling (stop exposing sql query on front-end)
- [ ] normalize URLs
- [ ] improve URL validation
- [ ] add delete action
- [ ] improve styling (use bootstrap)
- [ ] setup basic crawler
- [ ] Home, add redirect when item successfully created
- [ ] add HTMX

## Crawler related stuff

- [ ] Probably add "health-checks" as a separate table.
- [ ] Show "last checked" column for website
- [ ] Detect status, have option to disable crawling for a site
- [ ] Maybe add option to disable crawling from the front-end too
- [ ] Show last crawl time somewhere on the dashboard
- [ ] Detect crawler status, use context to either try to relaunch it or try to reload the entire backend