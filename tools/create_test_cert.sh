#!/bin/bash

OUTDIR=keys/
if [ ! -d "$OUTDIR" ]; then
    mkdir -p ${OUTDIR}
fi
openssl req -x509 -newkey rsa:2048 -keyout ${OUTDIR}key.pem -out ${OUTDIR}cert.pem -days 365 -nodes
