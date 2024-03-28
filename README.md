# nds-converter

## What is it?
A lightweight API-first tool to convert NDS save files (`.sav`) to and from the DeSmuME format (`.dsv`).

## How do I use it?
First, you need a save file, either of the two formats will work.
Then, you can just use `curl`:

```
curl localhost:1111 -F file=@savefile.sav > savefile.dsv
```

The conversion works both ways.
```
curl localhost:1111 -F file=@savefile.dsv > savefile.sav
```
