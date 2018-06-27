
# Network Service Scanner

A tool to help Network Users test access to services.

To test access to network services.  Two design goals are: to be able to build a single standalone binary and to have concurrent tests running at the same time.

The tests are defined in csv files and are in the tests directory.  The first line of the file is an Area name.  Then each line is a test

```
Area
Name,Hostname,IP,ports
```


