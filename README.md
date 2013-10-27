fastpack
========

fastpack is intended to be a tool which can pack files into a single archive
where each files is compressed and encrypted independently. The idea is that
single files can be extracted without dealing with the complete archive.

This design also makes it possible to use different compression and encryption
types in each file.

**BIG WARNING NOTE**

This is also partially just an experiement with Go. It is my first attempt
at anything with the language, and code is being written in what inevitably
will become the 'What was I thinking' phase.

## Initial Thoughts
Initially I am using the [snappy] compression algorithm and the encryption
method has not been decided at this time. 

[snappy]: https://code.google.com/p/snappy/ 