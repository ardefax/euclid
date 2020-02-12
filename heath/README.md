# Heath

Contains the XML text for Heath's version of Euclid's Elements and tool(s) to
transform into a JSON representation. Said JSON is a more useful intermediate format
for generating the site via Hugo.

There have been some hand modifications to the XML files downloaded from the Perseus
Project where errors or significant inconsistencies exist. However, attempts have
been made to limit these by making the transforms programmatic.

## Development

Build and run the go binary

```sh
go build
./heath -d ../data/heath vol?.xml
```
