SACAPI
======

This is a API middleware for a smart alarm clock.

It connects calendar APIs (google calendar so far)
to google maps and provides calculated data via
a REST json API.

Since its purpose is to offer a REST api to the
smart alarm clock, which hard and software you don't
have, the code is probably not that useful for most people.

Still it is licensed under Apache 2.0 license.

It uses Go with (among others) the very interesting
packages:

- github.com/ant0ine/go-json-rest

and

- golang.org/x/net/context

so if you're interested, have a look.

To set the service up just run

    make

in this direcrtory.

This also creates a snakeoil cert/key  pair, so you're ready to play around with it.
