# pingpong

Executable pingpong reads urls for gifs from a text file, and concurrently converts them to play forwards and then backwards.

The executable that builds from this package accepts three flags:

  - `--urls` path to a local text file containing a newline separated list of urls
  - `--dir` path to a local directory in which the processed gifs should be saved
  - `--trans` optional boolean value that indicates whether transparency should be corrected on reversal, avoiding a ghosting effect on certain gifs

Example:

```
go build
pingpong --urls ./urls.txt --dir ./gifs --trans
```
