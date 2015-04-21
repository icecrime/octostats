Octostats [![Build Status](https://travis-ci.org/icecrime/octostats.svg)](https://travis-ci.org/icecrime/octostats)
=========

GitHub repository statistics to InfuxDB.

## Usage

    $> docker build -t icecrime/octostats .
    $> docker run --rm -t -v `pwd -P`/.gittoken:/.gittoken -v `pwd -P`/octostats.json:/octostats.json icecrime/octostats --config /octostats.json
