# logparser

Command line log parser.

Allows to extract all the matching lines from a set of files defining filters using regular expressions or logical combinations of simple "contains" expressions.  

It keeps track of the already processed files and lines using an offset file.

## Version

Version 1.0.0

## Usage

```
USAGE: logparser -f=<config_file> -o=<offset_file> -l=<log_file> [-a=true] [-i=<ignore_old_files>]

  -a    Manage all matching files (default: manage only last one)
  -f string
        Configuration file
  -i int
        Ignore files older than N days (default 3)
  -l string
        Log file to parse
  -o string
        Offset file
```

The log file parameter can contain the '*' wildchar to identify multiple files. In this case, the -a parameter allow to either process all the files or only the most up-to-date.

By default, only the last one will be processed.

The offset file will be used to keep track of the already processed files and lines.

## Configuration file

The configuration file uses the YAML syntax. It is required to define a list of filters, each filter has the following properties:

- name: Name of the filter
- filtertype: Type of the filter, allowed values are:
  - REGEXP: Filter based on a single regular expression
  - EXPRESSION: Filter based on a set of "contains" expressions and logical operators to mix them
- pattern: It contains the regular expression (for REGEXP filters) or the "contains" expression (for EXPRESSION filters)

The allowed logical operators for an EXPRESSION filter are:

- && (Logical AND)
- || (Logical OR)
- ! (NOT)

Example:

```
filters:
  - name: Error filter
    filtertype: REGEXP
    pattern: .*stringa*
  - name: Expression Filter
    filtertype: EXPRESSION
    pattern: ("ciccio pluto" || "paperino")
```

## File to parse

The -l option identifies the log file/files to parse.

It is possibile to use standard wildchars (such as *) to define the files to be parsed.

Example:

```
logparser -f="../configs/config.yml" -o="../offset/logparser.off" -l="../data/*.txt" -a=true -i=4
```

## Offset file

The offset file will be used to avoid the re-processing of the same lines within a file. 

It contains one entry for each parsed file using the following syntax:

<file_name> <offset_in_bytes>

The entries within the offset file are automatically removed whenever the file does not exist anymore

## Configuration file

The matching lines wille be printed to stdout enclosed within &lt;matching-lines> and &lt;/matching-lines> tags

Please bear in mind that, if the line matches more than one filter, it will be printed only once
