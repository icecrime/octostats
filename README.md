Octostats
=========

GitHub repository statistics to Graphite

## Usage

    $> docker build -t icecrime/octostats .
    $> docker run --rm -t -v `pwd -P`/.gittoken:/.gittoken -v `pwd -P`/octostats.json:/octostats.json icecrime/octostats --config /octostats.json
