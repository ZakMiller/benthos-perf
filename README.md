# A simple tool to investigate how benthos works

- We create a new message every second and produce it to kafka/redpanda
- Benthos takes those messages and sends them to our api.
- We return a 500 for some of the messages.
- We're interested in the throughput we get.

Watch stdout. There are three columns:
`<timestamp> | <status code>          <guid>`


 You'll notice that when a 500 happens everything stops for 5 seconds. This is because we're only allowing one message through at a time (max_in_flight=1), and have retry_period set to 5s.

This is intended behavior from the perspective of benthos, but _will_ cause head of line blocking if you have some messages that can't be processed.

If, for example, you have 3 retries, a retry_period of 1 second, and 5 messages in a row that can't be processed, then that means that no messages will get through for 15 seconds.