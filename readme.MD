## ![](RackMultipart20210320-4-iaa6th_html_5142fe91682defd7.png)

# **Given a stream of ride measurements calculate design a high performance script to calculate ride fares**


# Different solutions of the problem:

- The naive approach will be loading the whole file into memory and performing all the operations on it and then writing it to csv. This is a synchronous approach and it does not scale. Memory and speed both will be a bottleneck here.Also I/O is a bottleneck as we will be pushing large csv&#39;s.
- The second solution is instead of reading the whole file we can buffer the file and read the data line by line. Instead of doing all the filtering and price estimation synchronously we can do it through an asynchronous pipeline made of goroutines and channels.

-
# Approach to solve the problem:

We take the asynchronous approach to solve the problem. We will stream the data into our system and the data will flow into different stages where the transformations and estimations will take place.At last the fares for each ride gets consumed from the channel and written into the output csv file.All the stages are defined in **pipelines.go** file. Let&#39;s talk about the stages now:

1. **emitStructuredRecords:** This stage reads the data using a buffered reader , converts them into its corresponding Point struct and emits them into a channel.
2. **groupUniqueRides:** Consumes Points and groups unique rides into unique channels and emits them into a different channel.
3. **filterInvalidPoints:** Consumes pointes in a ride and filters out according to the filtering logic given in assignment.
4. **estimateFare:** Consumes points in a ride and estimates the fare according to the rules given. Then emits the fare into a different channel.

-
# Scope of improvement:
- Deployment script is not written as of now for the whole system.
- Currently this script takes ~7 seconds to process a 470mb csv data file in my local machine. I believe this can still be improved by playing around with the channel buffer size and stages.
- As of now the input file path and output file path is placed in constants. But Ideally it should be configurable(Ex. Command line args , Environment variable etc.)
- We are ignoring invalid entries right now. But we can dump those invalid data into someplace for better analytics and insights.
- Current test coverage is 88.0% as of now. We can increase the test coverage by writing more test cases.

-
# Scalability:
- We can scale the solution using vertical scaling. Here we need to increase the machine power(RAM,CPU etc). But it&#39;s not recommended as it has a upper limit and it is very expensive.
- Another approach is scaling it horizontally. In this approach we will use multiple machines. We can add or remove machines depending on the latency requirements. In this approach there isn&#39;t any upper limit(Except the cost). The monolith csv will be broken into smaller pieces in this approach and they will be assigned to specific machines which will take care of processing it. In this approach we also need an aggregator which will aggregate all the results if required.
-
