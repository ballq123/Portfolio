# Distributed Bitcoin Miner – Fault-Tolerant Client–Server System

**Timeline:** Oct 2024  
**Focus:** Distributed systems, concurrency, fault tolerance  
**Tech Stack:** Go (goroutines, channels), UDP networking, custom reliability protocol

A distributed system written in Go that parallelizes compute-intensive jobs across multiple worker nodes with **dynamic load balancing, concurrency, and automated failure recovery**.  
Implements a lightweight reliability protocol over UDP, ensuring correct execution even under node churn and network failures.  
Demonstrates skill in **distributed systems, concurrency, networking, and protocol design**.

---

## Files

- [`server_excerpt.go`](code_snippets/server_excerpt.go) → Server excerpt  
  - Assigns job ranges to workers using a **round-robin scheduler with chunking**.  
  - Detects dropped workers and **reassigns incomplete work** to idle nodes.  
  - Aggregates results and returns the final minimum-hash solution to the client.  
  - Demonstrates **fault tolerance and stateful scheduling**.  

- [`miner_excerpt.go`](code_snippets/miner_excerpt.go) → Worker excerpt  
  - Joins the system, receives job assignments, executes work in a nonce range, and returns results.  
  - Demonstrates **concurrent worker logic, result reporting, and resiliency under reassignment**.

- [`client_excerpt.go`](code_snippets/client_excerpt.go) → Request client excerpt  
  - Submits a compute job and waits for the aggregated final result from the server.  
  - Demonstrates **request orchestration, JSON message encoding, and reliable communication with a distributed backend**.  

*(Note: Course-specific networking libraries have been redacted and replaced with minimal interfaces/types. Excerpts show original scheduling, concurrency, and fault-tolerance logic.)*

---

## Skills Demonstrated

- **Fault Tolerance & Reliability**  
  - Detects failed or disconnected workers and reassigns tasks without manual intervention.  
  - Preserves partial progress by requeuing interrupted work.  
  - Implements a lightweight reliability protocol on top of UDP with retransmissions and acknowledgments.  

- **Concurrency & Systems Programming**  
  - Used Go **goroutines + channels** for scalable, lock-free concurrency.  
  - Managed job scheduling, worker churn, and result aggregation without bottlenecks.  
  - Scaled to multiple concurrent workers while maintaining correctness.  

- **Load Balancing & Scheduling**  
  - Implemented **round-robin job distribution with chunking**, ensuring fair allocation across workers.  
  - Balanced workload dynamically as workers joined or left the system.  
  - Aggregated partial results into a consistent client output. 

---

## Impact

- Built a **fault-tolerant distributed job execution framework** that continued operating correctly under worker crashes and network dropouts.  
- Demonstrated **scalable throughput**: distributing work to multiple miners reduced completion time significantly compared to sequential execution.  
- Showcased robust **distributed systems engineering**: concurrency, custom networking, fault recovery, and dynamic load balancing.   