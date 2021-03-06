dte
===

Django template executor

See [The Django template language | Django documentation | Django]( https://docs.djangoproject.com/en/dev/topics/templates/ ) for template syntax.

## Download binary

OS X

```
curl -O https://hnakamur.github.io/dte/download/darwin_amd64/dte
```

Linux 64bit

```
curl -O https://hnakamur.github.io/dte/download/linux_amd64/dte
```

Linux 32bit

```
curl -O https://hnakamur.github.io/dte/download/linux_386/dte
```

Windows 64bit

```
curl -O https://hnakamur.github.io/dte/download/windows_amd64/dte
```

Windows 32bit

```
curl -O https://hnakamur.github.io/dte/download/windows_386/dte
```

## Usage

```
% ./dte -h
Usage of ./dte:
  -j="-": json filename (default stdin)
  -o="-": output filename (default stdout)
  -v=false: show version and exit
```

## Example

example/data.json

```
{
  "persons": [
    {"name": "Alice"},
    {"name": "Bob"},
    {"name": "Charlie"}
  ]
}
```

example/hello.tpl

```
{% for person in persons %}{% if not forloop.First %}
{% endif %}Hello, {{ person.name }}!{% endfor %}
```

Example session

```
$ ./dte -j example/data.json example/hello.tpl
Hello, Alice!
Hello, Bob!
Hello, Charlie!
```

