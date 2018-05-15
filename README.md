# butter

Useful command line tools for interacting with BTrDB. We pronounce BTrDB as 
*"Better DB"*, so the name "butter" is a play on that.

## Installation

Run this to install butter:

```
git clone git@github.com:PingThingsIO/butter.git
go get
go install
```

## Usage

Butter currently has 4 sub commands: `cp`, `ls`, `rm`, and `tail`. You can read
about them in this README, or run `butter -h` for a help page:
```
$ butter -h

Usage: butter COMMAND [arg...]

Useful BTrDB CLI tools for development

Commands:
  ls           List collections for a BTrDB endpoint. If only one collection is returne
d, its streams will be listed.
  rm           Remove a stream from BTrDB
  cp           Copy a collection from one BTrDB server to another
  tail         Prints the latest values inserted into BTrDB

Run 'butter COMMAND --help' for more information on a command.

```


### Copy

The `cp` subcommand will copy streams from one BTrDB server to another. It is
essentially a fork of (btrdbcp)[https://github.com/BTrDB/smartgridstore/tree/master/tools/btrdbcp].

Here's the usage text:
```
$ butter cp -h

Usage: butter cp FROMSERVER TOSERVER [-sea] STREAMCONFIG...

Copy a collection from one BTrDB server to another

Arguments:
  FROMSERVER     BTrDB endpoint to copy from
  TOSERVER       BTrDB endpoint to copy to
  STREAMCONFIG   Config for the streams to copy (follows the format src_collection,dest
_collection,tagname=tagvalue) (default [])

Options:
  -s, --start    Start time of the range to copy (in format 2006-01-02T15:04:05+07:00)
  -e, --end      End time of the range to copy (in format 2006-01-02T15:04:05+07:00)
  -a, --abort    Abort the copy if the collection already exists

```

Here's an example of how to use it:
```
butter cp example.com:4410 localhost:4410 \
    --start 2017-12-10T15:04:05+07:00 \
    --end 2017-12-11T15:04:05+07:00 \
    pingthingsio/90807,pingthingsio/90807,name=L1MAG
```

The output should look something like this:

```
pingthingsio/90807/"name"="L1MAG", has 10368000 points
pingthingsio/90807/"name"="L1MAG", 855565 / 10368000 [=>---------------]   8.25% 01m26s
```

Explanation of the example:
 * The first two positional arguments specify the source and destination servers.
 * The `-s` and `-e` flags specify the range from start and end to copy. 
 * The last arguments are the stream config, which is a comma separated list
   of the source collection, dest collection, and tag key value pairs separated
   by equal signs. Multiple stream configs can be specified, which must be separated
   by a space.


### List

The `ls` subcommand will list collections of a BTrDB server.

Here's it's usage:
```
Usage: butter ls [ENDPOINT] [PREFIX]

List collections for a BTrDB endpoint. If only one collection is returned, its streams
will be listed.

Arguments:
  ENDPOINT     The BTrDB endpoint to list (default "localhost:4410")
  PREFIX       A prefix to filter collections.
```

It's output looks like this:

```
Collection name                      Stream count
pingthingsio/sensor1                 19
pingthingsio/sensor2                 16
pingthingsio/sensor3                 19
```

You can use the prefix argument to narrow down result. If only one
collection is returned from the query, stream result will be displayed:

```
$ butter ls localhost:4410 pingthingsio/sensor1

Collection: pingthingsio/sensor1:
Streams:
 * UUID: b66fa23a-4abd-53b7-99f9-651f5f2fa3b1
 * Tags:
     - name: R1HNG
 * Annontations:
     - None

 * UUID: bf29bd31-0d8f-56c6-af3c-33251adb1009
 * Tags:
     - name: F2MAG
 * Annontations:
     - None
```


### Remove

The `rm` subcommand can be used to delete a stream. This should really
only be used for development, please only delete data in production
databases through the proper channels (e.g. the smartgridstore admin interface)

Here's the usage for `rm`:
```
$ butter rm -h

Usage: butter rm [ENDPOINT] UUID [-y]

Remove a stream from BTrDB

Arguments:
  ENDPOINT     The BTrDB endpoint to list (default "localhost:4410")
  UUID         UUID of the stream to delete

Options:
  -y, --yes    Skip confirmation prompt
```

Here's an example:

```
butter rm b66fa23a-4abd-53b7-99f9-651f5f2fa3b1
```

This will prompt for confirmation:

```
Are you sure you want to delete stream b66fa23a-4abd-53b7-99f9-651f5f2fa3b1 from collection pingthingsio/90807? [y/n]
```
