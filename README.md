# Firelink

Shorten your link using Firebase Dynamic Link service.

## Installation

Use `go get`.

```bash
go get github.com/nguyenvanduocit/firelink
```

## Usage

```python
firelink --link=very_long_long_link 
```
### Arguments

| Name          | Short name | Default              | Description                     |
|---------------|------------|----------------------|---------------------------------|
| --config      | -c         | $HOME/.firelink.yaml | Link to your configuration file |
| --key         | -k         |                      | Your Firebase web API key       |
| --prefix      | -p         |                      | Your Dynamic Link prefix        |
| --link        | -l         |                      | Very long long link             |
| --title       | -t         |                      | Custom social title             |
| --description | -d         |                      | Custom social description       |
| --imageLink   | -i         |                      | Custom social image link        |

If `--link` is not provided, the program will try to get from clipboard.

### Configuration

Put all your configurations in `$HOME/.firelink.yaml`:

```yaml
webApiKey: RazaSydAizT1EUAcLH2R9BZ_Ah_AMeUtfbdssZq
domainUriPrefix: https://xn--gh-fja.vn
```

if `domainUriPrefix` is unicode, please use it in `idna` format.

## Contributingz
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](https://choosealicense.com/licenses/mit/)