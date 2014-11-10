Octostats
=========

GitHub repository statistics to Graphite

## Usage

    $> docker build -t icecrime/octostats .
    $> docker run --rm -t -v ~/.gittoken:/.gittoken icecrime/octostats --repository=icecrime/octostats --token-file=/.gittoken
