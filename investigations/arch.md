### High level architecture of client - Short term

Within the scope of PAP

##### Basic requirements:

1. Should be able to handle PAP flow.
2. Dynamic handling of Access-Accept/Access-Reject.
3. Packet validation feature for 1,2,3.
4. Metrics to be supported:
    Latency
    TPS
    Total Requests
    ErrorRate
    (others)
5. Timeouts and Retries to be supported
6. Metric export (JSON,CSV) or to prometheus/graphana. Maybe even a metrics endpoint.
7. Maybe support for input files that can save AVPs and connection information etc; like JMETER.
8. Rate limiting

##### Routine architecture:

Need more exploration



