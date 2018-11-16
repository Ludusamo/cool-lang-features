# Cool Language Features
## Brendan Horng

Website for Parallel Processing Course creating a website.

## Getting Started

### Running

You can run the code inside `cool-lang-features/`:

```
go run frontend/frontend.go [--listen {port} --backend {backendAddress}]
go run backend/backend.go [--listen {port}]
```

### Testing

Currently the routes are as follows:
- `/`: Serves the listing of all the features. From here you can navigate to delete, modify, or add new features
- `/add`: Form to add a new feature
- `/modify/{id}`: Form to modify a feature identified by id
- `/api/feature`: REST API defined here for GET and POST commands to retrieve all features or to add a new one
- `/api/feature/{id}`: REST API defined her for GET, PATCH, and DELETE to retrieve, modify, or delete specific features based on id

#### Vegeta Testing

Inside my `vegeta` folder, I created a Python script, `attack_gen.py` to generate test traffic for me. It is not particularly robust, and it is tailored to go along with my script `attack.sh`.
Running `attack.sh` will compile the frontend and the backend and run both in the background. It will then use the Python script to generate an attack and the perform that attack with Vegeta and prints out a report.
The background processes are killed and the temporary data folder is cleaned up after the attack.

## Other Notes

### Current Status

The front end code is a little more unclean than I would like it to be, but the overall functionality is there. The web server, written in Go, is capable of doing add, delete, read, and modify on "features".

#### Update for Part 2

The above is what I wrote for part 1 and still relevant so I left it as is. The frontend can be spun up separately from the backend (see **Running**).

I don't like how I ended up structuring the code base. I will probably refactor it in the near future. Particularly what I don't like is because both my frontend and my backend use **main** as their package, I can't do a `go build`.

To properly do this, I would probably have to separate out everything even further from the current folder hierarchy that I have. For simplicity, I left it as is for now.

I also didn't stay on top of my test cases as much as I would have liked, so I have omitted them this time.

#### Update for Part 3

Important Design Decisions:

1. I provided safe, concurrent, performant access to my data store by locking the data store during operations.
In order to make it more performant, I locked the data store in "buckets".
This way, certain chunks of the data store could be locked at a given time while others could still be used.
I achieved this by having a set of locks (specifically 10), and the entire data store was partitioned into these 10 buckets based on `id` number.
I also used `RWMutex`s to achieve simultaneous reads to improve performance.
There are a couple of downsides to my approach that I would have liked to change given more time.
    * The number of buckets is fixed. If I had more time I would have this scale based on how many elements are in the data store.
    * For adding a new element or reading the entire set, you still need to get all of the locks.
I considered locking the entire data store. This would be much simpler to implement, but would also increase lock contention.
I also considered streamlining the commands through a set of goroutines that could process these in a more intelligent way (to be determined). This seemed like it would have taken more time to get right, so I decided not to go down this route.
The performance metric that I found to be important was the average latency from the Vegeta report. I wanted the average to be low. I constructed attacks that were more heavy on the reading of data rather than update, write, or delete because that would be a more realistic use case of this application.
The average latency prior to my optimazations were ~288 milliseconds.
The average latency after my optimatizations were ~368 microseconds.
2. I implemented the failure detector using a heartbeat approach.
I had the frontend send a message asking the backend to subscribe to a heartbeat. After this point, the backend would send the frontend a heartbeat every 10 seconds.
These heartbeats were sequenced. If the front end does not receive an incremented heartbeat after 35 seconds, it presumes the backend is dead. It then attempts to reconnect.
I considered ping-ack as well, but I figured that heartbeat would be better for longer running processes. I thought ping-ack would be more network traffic than necessary.
A ping-ack would have to be done on an interval just like the heartbeat, but at no point do we want to stop receiving confirmation. So instead of sending a request one way and receiving one back, I decided it was better just to subscribe once and continually receive confirmation.

### Additional Resources

- http://elm-lang.org/
    - I used Elm to create the front end of this application. I am still pretty new to it so the code code be improved.
    - You should not have to do anything with it, I built the one html file that the Go server will be hosting.
- https://www.python.org/
    - I used Python3 to help me generate traffic for my Vegeta attacks
- https://www.gnu.org/software/bash/
    - I used bash scripting to setup and teardown for an attack with Vegeta

### Additional Thoughts

It took longer than I thought it would. I enjoy these kinds of projects though, and I thoroughly enjoyed it because I was using two languages I don't typically get to.

#### Update for Part 2

No additional insights.

#### Update for Part 3

I really like this assignment the further we get into it. It is fun to have a challenging architecture to model, and I like having a project that could actually be useful after I am done writing it.
It has the potential to not be throwaway work, and I find the assignment very practical.
