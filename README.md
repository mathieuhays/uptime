# uptime

Simple tool to monitor multiple URLs


## TODO

- [ ] Website View: add option to toggle cache bypass setting?
- [ ] Website View: add option to disable crawling
- [ ] Unit test crawler
- [ ] Record website outages in different table
- [ ] Record a marker on deployment
- [ ] Website: Multi-add through comma separated list
- [ ] Website: Import from CSV
- [ ] Website: add ability to change the name


## Ideas

- [ ] Setup websockets for instant updates
- [ ] Have an image attached to each website (automated). either screenshot, social media meta image or favicon. screenshot will probably require nodejs.
- [ ] Website view: maybe infer page colour theme from target website?
- [ ] website setting: response time budget, then color code health check historic to reflect it


## Development

### Tools

- Air: `go install github.com/air-verse/air@latest`
- GoSec: `go install github.com/securego/gosec/v2/cmd/gosec@latest`
- StaticCheck: `go install honnef.co/go/tools/cmd/staticcheck@latest`

### Running locally

I use the `air` command to run the app locally. It'll make the HTTP server accessible on http://localhost:8081 and handle auto-reloading.

The air command should generate a tmp folder. Running the app should create an uptime.db file in there and run all the migrations.