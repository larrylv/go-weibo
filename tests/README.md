go-weibo tests
===============

This directory contains additional test suites beyond the unit tests already in
[../weibo](../weibo).  Whereas the unit tests run very quickly (since they
don't make any network calls) and are run by Travis on every commit, the tests
in this directory are only run manually.

The test packages are:

fields
------

This will identify the fields being returned by the live weibo API that are
not currently being mapped into the relevant Go data type.  Sometimes fields
are deliberately not mapped, so the results of this tool should just be taken
as a hint.

While the tests will try to be well-behaved in terms of what data they modify,
it is **strongly** recommended that these tests only be run using a dedicated
test account.

Run the fields tool using:

    WEIBO_AUTH_TOKEN=XXX go run ./fields/fields.go
