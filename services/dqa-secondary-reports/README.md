# Service for DQA secondary report information

Service that exposes a endpoint for fetching PEDSnet DQA secondary report information for different data model versions, e.g.

```
/pedsnet/2.0.0
/pedsnet/1.0.0
```

## Installation

To install the prerequisites, run
```
pip3 install -r requirements.txt 
```

## Run

To get usage information, run:
```
python3 main.py -h
```

To run the program with default connection settings, run
```
python3 main.py --token=abc123
```
where abc123 is your GitHub authorization token. 
As an alternative to supplying the token through the --token option, 
`GITHUB_AUTH_TOKEN` environment variable can be set
