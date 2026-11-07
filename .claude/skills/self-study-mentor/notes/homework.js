/* Homework sets for Athena Study Notes.
 *
 * Every interview question in progress-log.md's Interview Question Bank appears in exactly one set's
 * "interview" part, verbatim, with refs to the chapter sections that answer it. Objective parts (match,
 * blank, mc) carry their keys so they can self-check; open parts (short, interview) carry no answers on
 * purpose, only where to look. Refs use the router's "chNN-sN" form.
 */
window.ATHENA_HW = [

// =================================================================================================
{
  id: "hw01", chapter: "ch01", title: "Bounded Queues and Worker Pools",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch01-s2", pairs: [
      ["Producer / consumer queue", "Handlers put readings in; workers take them out."],
      ["Bounded queue", "Has a fixed capacity, so overload shows up as refusals instead of memory growth."],
      ["Ring buffer", "A fixed array whose read and write positions wrap back to slot 0."],
      ["Backpressure", "Slowing the producer down when the consumer can't keep up."],
      ["Load shedding", "Refusing work you can't handle, quickly and on purpose."],
      ["Drain", "On shutdown, finishing what's already queued before exiting."],
      ["At-most-once", "Each accepted reading is written once or not at all."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "Inside the Go runtime, a channel is a struct called ___.", answers: ["hchan"], ref: "ch01-s3" },
      { text: "Little's Law, L = &lambda;W: at 100,000 readings/s, each waiting 2 ms, about ___ readings are queued on average.", answers: ["200"], ref: "ch01-s7" },
      { text: "Sending on a closed channel causes a ___.", answers: ["panic"], ref: "ch01-s6" },
      { text: "By default, Kubernetes waits ___ seconds after SIGTERM before it kills the pod.", answers: ["30"], ref: "ch01-s6" },
      { text: "A full queue answers with HTTP status ___ and a <code>Retry-After</code> header.", answers: ["503"], ref: "ch01-s4" },
      { text: "With one queue per worker, a very busy device fills its own queue while others sit empty: a hot ___.", answers: ["partition"], ref: "ch01-s5" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "With N workers reading an unbuffered channel, the queue's effective capacity is:", options: ["N", "Zero, always", "However many workers are idle right now", "Unbounded"], answer: 2, ref: "ch01-s3" },
      { q: "Why does validation run before the offer to the queue?", options: ["Validation is slow, so it should go first", "So a malformed reading never takes a slot a good one could use", "The queue can't hold invalid structs", "So that 400 responses are faster than 503 responses"], answer: 1, ref: "ch01-s4" },
      { q: "During a sustained 20k/s overload, doubling the queue's capacity:", options: ["Prevents shedding", "Delays shedding but doesn't prevent it", "Halves request latency", "Doubles throughput"], answer: 1, ref: "ch01-s7" },
      { q: "Who closes the queue's channel?", options: ["Any handler that sees an error", "Each worker when it goes idle", "Whoever started the queue, once every sender is gone", "Nobody; the garbage collector does"], answer: 2, ref: "ch01-s6" },
      { q: "When a Go scheduler's local run queue (a 256-slot ring) fills, it:", options: ["Blocks the goroutine that's spawning", "Spills half of it to a global queue", "Drops the oldest goroutines", "Grows without bound"], answer: 1, ref: "ch01-s3" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "Sketch Fig 1.5's shutdown order from memory. What goes wrong if the queue is closed before the HTTP server stops accepting requests?", refs: ["ch01-s6"] },
      { q: "Using Fig 1.6's numbers, explain why a buffer buys time and not throughput.", refs: ["ch01-s7"] },
      { q: "Why is a shared queue a better load balancer than per-worker queues, and why doesn't routing to the least-loaded worker fix per-worker queues?", refs: ["ch01-s5"] },
    ]},
    { type: "interview", title: "Interview questions", items: [
      { q: "What happens when your consumer is slower than your producer?", refs: ["ch01-s4", "ch01-s7", "ch04-s5"] },
      { q: "How did you size that buffer? <em>(Little's Law: L = &lambda;W)</em>", refs: ["ch01-s7"] },
      { q: "How many workers would you need once the store is SQS? Why is that a question about your target throughput and per-call latency rather than about SQS's capacity? What does batching do to the number?", refs: ["ch01-s7"] },
      { q: "SRE ch. 22 says keep queues &le; 50% of the thread pool. Why doesn't that rule transfer to a fire-and-forget queue?", refs: ["ch01-s7"] },
      { q: "Shared queue vs. a queue per worker: which one balances load better, and why does routing to the least-loaded worker give up the only thing per-worker queues were for?", refs: ["ch01-s5"] },
      { q: "Who closes the channel, and why can't it be the senders?", refs: ["ch01-s6"] },
      { q: "What's the production failure mode of an unbounded queue?", refs: ["ch01-s7"] },
      { q: "When your queue is full you return <code>503</code>, not <code>429</code>. Why? When would <code>429</code> be the right answer?", refs: ["ch01-s4", "ch10-s4"] },
      { q: "How do you test that the service sheds load? How did you make the queue full on purpose?", refs: ["ch01-s9", "ch02-s3"] },
    ]},
  ],
},

// =================================================================================================
{
  id: "hw02", chapter: "ch02", title: "Testing Outside-In",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch02-s3", pairs: [
      ["Dummy", "Passed to fill a parameter and never used."],
      ["Fake", "A working implementation with a shortcut that makes it unfit for production."],
      ["Stub", "Controls an indirect input: canned answers, or when a call returns."],
      ["Spy", "Records what was called, for assertions afterward."],
      ["Mock", "Pre-programmed with expectations; verifies behaviour itself."],
      ["Contract test", "One test body run against many implementations of an interface."],
      ["Sociable test", "The unit under test uses its real collaborators."],
      ["Walking skeleton", "The thinnest real slice, built, deployed, and tested end to end before any feature."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "<code>TelemetrySpec</code> is a contract on a ___ port; <code>StorerContract</code> is a contract on a driven port.", answers: ["driving"], ref: "ch02-s2" },
      { text: "Acceptance cases need no reset because each one makes its own device with a fresh ___ as its ID.", answers: ["uuid", "random id", "random uuid"], ref: "ch02-s4" },
      { text: "The test pyramid inverted, with most coverage in slow end-to-end tests, is called the ___.", answers: ["ice-cream cone", "ice cream cone"], ref: "ch02-s5" },
      { text: "The acceptance module points at the local service with a ___ directive in <code>go.mod</code>, not a <code>go.work</code> file.", answers: ["replace"], ref: "ch02-s9" },
      { text: "The standard library's <code>testing/fstest.___</code> checks any <code>fs.FS</code> implementation.", answers: ["TestFS"], ref: "ch02-s4" },
      { text: "GOOS ch. 8: only mock types that you ___.", answers: ["own"], ref: "ch02-s3" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "In this repo, what determines a test's level?", options: ["The test body", "The entry point and what it binds the body to", "The file name", "Whether it uses testify"], answer: 1, ref: "ch02-s1" },
      { q: "The gated <code>Storer</code> whose <code>Insert</code> blocks until the test releases it is a:", options: ["Fake", "Stub", "Spy", "Mock"], answer: 1, ref: "ch02-s3" },
      { q: "Why isn't <code>Validate</code> a method on <code>models.Telemetry</code>?", options: ["Methods on exported structs are slower", "A driver could validate client-side and pass the Spec without the server doing anything", "JSON can't marshal types with methods", "A lint rule forbids it"], answer: 1, ref: "ch02-s9" },
      { q: "Which test proves the Dockerfile and env-var config actually work?", options: ["The domain-level Spec run", "The httptest handler tests", "The container acceptance test", "StorerContract"], answer: 2, ref: "ch02-s5" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "Give one example from this repo of a test failing for the wrong reason, and one of a test passing for the wrong reason.", refs: ["ch02-s7"] },
      { q: "Why does a fake need a contract test when a mock doesn't?", refs: ["ch02-s3"] },
      { q: "List three costs of putting <code>specifications</code> in an exported, non-test package.", refs: ["ch02-s4"] },
    ]},
    { type: "interview", title: "Interview questions: cross-cutting", items: [
      { q: "Walk me through the kinds of test in this repo. How many distinct kinds are there, and what does each buy that the others can't?", refs: ["ch02-s2", "ch02-s1"] },
      { q: "<code>TelemetrySpec</code> runs at three levels. What does each level catch that the one beneath it can't? If the suite got too slow, which would you delete first, and what risk would you be accepting?", refs: ["ch02-s5"] },
      { q: "<code>specifications.TelemetrySpec</code> and <code>data.DataStoreSpec</code> are both called &ldquo;spec&rdquo; and are not the same pattern. What's the difference? <em>(The storage body is now <code>StorerContract</code>.)</em>", refs: ["ch02-s2"] },
      { q: "Why does <code>specifications</code> live in an exported, non-test package instead of a <code>_test.go</code> file? What does that cost you, and who else does this?", refs: ["ch02-s4"] },
      { q: "Your validation table has seven cases; the Spec has one invalid-input case. Why don't all seven run at the acceptance level?", refs: ["ch02-s5"] },
      { q: "When <code>telemetry.Query</code> is tested against a real in-memory store, is that a mock, a stub, or a fake? Why does the distinction matter here?", refs: ["ch02-s3"] },
      { q: "Why outside-in rather than inside-out? What goes wrong if you build the store first?", refs: ["ch02-s6"] },
      { q: "What does &ldquo;red for the right reason&rdquo; mean? Give a case where a test of yours failed for the <em>wrong</em> reason and how you noticed.", refs: ["ch02-s7"] },
      { q: "Your acceptance cases need no reset between runs. How did you get that, and why is it better than a reset hook on the interface?", refs: ["ch02-s4"] },
      { q: "A test was awkward to write and you changed the design instead of the test. Walk me through it.", refs: ["ch02-s8", "ch05-s2"] },
      { q: "What's an ice-cream cone, and what in this repo is most likely to produce one?", refs: ["ch02-s5"] },
    ]},
    { type: "interview", title: "Interview questions: M1, walking skeleton", items: [
      { q: "Why build a walking skeleton before any feature?", refs: ["ch02-s6"] },
      { q: "What does your container acceptance test prove that the <code>httptest</code> one doesn't? What does it cost per run?", refs: ["ch02-s9", "ch02-s5"] },
      { q: "Why does <code>httpserver.Driver</code> hardcode <code>/v1/telemetry</code> instead of importing <code>api.TelemetryPath</code>?", refs: ["ch02-s4"] },
      { q: "Why is <code>Validate</code> not a method on <code>models.Telemetry</code>?", refs: ["ch02-s9"] },
      { q: "There's no <code>go.work</code> despite the project notes saying &ldquo;Go workspaces.&rdquo; What does the <code>replace</code> directive do instead, and why does the dependency only flow one way?", refs: ["ch02-s9"] },
    ]},
  ],
},

// =================================================================================================
{
  id: "hw03", chapter: "ch03", title: "Shared Memory, Races, and Ownership",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch03-s1", pairs: [
      ["Data race", "Unordered concurrent access to the same memory where at least one access writes."],
      ["Atomicity violation", "Each step is locked, but the multi-step operation as a whole isn't."],
      ["Aliasing", "A reference to guarded data leaves the lock and is used without it."],
      ["Happens-before", "The ordering, created by locks and channels, that the race detector tracks."],
      ["Single-writer principle", "Only one goroutine ever touches the state, so no lock is needed."],
      ["Slice header", "A pointer to a backing array, a length, and a capacity."],
      ["Deep copy", "A new outer slice plus a fresh copy of each record's nested Metrics slice."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "Go's race detector is built on ___.", answers: ["ThreadSanitizer", "tsan"], ref: "ch03-s2" },
      { text: "Typical race-detector overhead: 5&ndash;10&times; memory and 2&ndash;___&times; execution time.", answers: ["20"], ref: "ch03-s2" },
      { text: "<code>slices.___</code> makes only a shallow copy, so each record's Metrics array is still shared.", answers: ["Clone"], ref: "ch03-s3" },
      { text: "To let several readers hold the lock at once, use <code>sync.___</code>.", answers: ["RWMutex"], ref: "ch03-s5" },
      { text: "A buffer moves waiting from request latency to ___: the time a reading spends in the queue.", answers: ["staleness"], ref: "ch03-s6" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "A <code>-race</code> run reports nothing. You can conclude:", options: ["The code has no data races", "No data race occurred in the paths this run executed", "There are no deadlocks", "The code is thread-safe"], answer: 1, ref: "ch03-s2" },
      { q: "&ldquo;Check whether the device exists&rdquo; and &ldquo;insert&rdquo; as two separately locked calls is caught by:", options: ["The race detector", "The compiler", "Design review, or a test that forces the interleaving", "<code>go vet</code>"], answer: 2, ref: "ch03-s1" },
      { q: "Unbuffered channel, non-blocking send, and a burst arrives while all N workers are busy:", options: ["Requests block until a worker frees up", "Every request in that moment is shed", "N requests are absorbed, the rest shed", "The send panics"], answer: 1, ref: "ch03-s6" },
      { q: "For a map with insert and lookup, the Go wiki's guidance points to:", options: ["A mutex", "An owning goroutine", "Atomics on each entry", "A buffered channel"], answer: 0, ref: "ch03-s5" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "Sketch Fig 3.1, then explain in two sentences why the mutex doesn't protect a returned slice.", refs: ["ch03-s3"] },
      { q: "Name three kinds of evidence that the store is safe, and which failure each one covers.", refs: ["ch03-s4"] },
      { q: "Name two situations where you'd switch from a mutex to an owning goroutine.", refs: ["ch03-s5"] },
    ]},
    { type: "interview", title: "Interview questions", items: [
      { q: "Buffered or unbuffered channel, and why? What does each do to tail latency under a burst?", refs: ["ch03-s6", "ch01-s3"] },
      { q: "How do you know your store is thread-safe? Name two different failures and which tool catches each.", refs: ["ch03-s4", "ch03-s1"] },
      { q: "What does <code>-race</code> actually prove? What does it miss?", refs: ["ch03-s2"] },
      { q: "<code>GetByDeviceID</code> hands back a slice &mdash; what's the risk, and why is a mutex not enough?", refs: ["ch03-s3"] },
      { q: "Why a mutex here rather than a single owning goroutine? When would you switch?", refs: ["ch03-s5"] },
    ]},
  ],
},

// =================================================================================================
{
  id: "hw04", chapter: "ch04", title: "Overload, Retries, and What 202 Promises",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch04-s6", pairs: [
      ["Group commit", "Acknowledge only after the batch containing this request is written."],
      ["Open system", "Arrivals don't depend on how fast you respond."],
      ["Closed system", "A fixed set of clients, each waiting for a response before sending again."],
      ["Thundering herd", "Clients retrying in lockstep, knocking the server over in waves."],
      ["Jitter", "Randomizing each backoff wait so retries spread out."],
      ["Retry budget", "A cap on retries, per request or as a fraction of normal traffic."],
      ["End-to-end argument", "Only the endpoints can fully implement a guarantee; lower layers only improve the odds."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "RFC 9110: after a <code>202</code>, the request may or may not eventually be ___ upon.", answers: ["acted"], ref: "ch04-s1" },
      { text: "If three layers each try up to four times, one request can become ___ requests at the bottom.", answers: ["64"], ref: "ch04-s7" },
      { text: "Retrying after a timeout turns at-most-once into at-___-once.", answers: ["least"], ref: "ch04-s7" },
      { text: "<code>SendMessageBatch</code> takes up to ___ messages per call.", answers: ["10", "ten"], ref: "ch04-s8" },
      { text: "Per-device ___ numbers let the destination detect every lost reading, including ones a crash took.", answers: ["sequence"], ref: "ch04-s3" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "Under sustained overload, which policy hides loss from the client?", options: ["Block", "Shed", "Drop", "None of them"], answer: 2, ref: "ch04-s5" },
      { q: "Why does blocking fail for a fleet of thermostats?", options: ["Go can't block an HTTP handler", "It's an open system: arrivals don't slow down when you do", "Thermostats can't handle slow responses", "Blocking uses more CPU than shedding"], answer: 1, ref: "ch04-s6" },
      { q: "Which response should a device <em>not</em> retry?", options: ["<code>503</code>", "<code>429</code>", "<code>400</code>", "A timeout"], answer: 2, ref: "ch04-s7" },
      { q: "Brooker's &ldquo;full jitter&rdquo; picks each wait:", options: ["Exactly equal to the backoff", "Uniformly between zero and the backoff", "The backoff plus one second", "Uniformly between 0 and 10 s, whatever the backoff"], answer: 1, ref: "ch04-s7" },
      { q: "Where does Athenas get <em>real</em> backpressure at M5?", options: ["<code>503</code> with <code>Retry-After</code>", "Consumers pulling, and pausing partitions when behind", "The SQS visibility timeout", "Redis rate limits"], answer: 1, ref: "ch04-s6" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "List three ways to measure how much telemetry you lose, cheapest first, and what each can't see.", refs: ["ch04-s3"] },
      { q: "Summarize Saltzer, Reed and Clark's careful-file-transfer example and map it onto your <code>202</code>.", refs: ["ch04-s4"] },
      { q: "Name three production clients that buffer in-process, and the trade each makes.", refs: ["ch04-s8"] },
    ]},
    { type: "interview", title: "Interview questions", items: [
      { q: "Your endpoint returns <code>202 Accepted</code> &mdash; what exactly have you promised the caller? What can still be lost?", refs: ["ch04-s1", "ch01-s6"] },
      { q: "You chose fire-and-forget over group commit. What did you trade away, and what would make you switch? How would you know how much telemetry you're actually losing?", refs: ["ch04-s2", "ch04-s3"] },
      { q: "Which errors should a telemetry sender retry, and which should it drop? Why does exponential backoff need jitter?", refs: ["ch04-s7"] },
      { q: "What is the end-to-end argument, and how does it apply to your <code>202</code>?", refs: ["ch04-s4"] },
      { q: "Block, shed, or drop: what does each do to latency and memory under sustained overload?", refs: ["ch04-s5"] },
      { q: "Is a <code>503</code> with <code>Retry-After</code> backpressure or load shedding? Why does it matter that your devices are an <em>open</em> system (they keep sending whether or not you answer)? Where in this project do you get real backpressure, and why there?", refs: ["ch04-s6"] },
      { q: "Why buffer in memory at all when you could write straight to SQS? What does each give up: latency, cost per request, durability of the <code>202</code>? Which production clients buffer in-process anyway?", refs: ["ch04-s8"] },
    ]},
  ],
},

// =================================================================================================
{
  id: "hw05", chapter: "ch05", title: "Configuration and the Composition Root",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch05-s4", pairs: [
      ["Composition root", "The one place the application's object graph is assembled: <code>main</code>."],
      ["12-factor III", "Store config in the environment."],
      ["12-factor V", "Keep build, release, and run as separate stages."],
      ["Functional options", "Variadic option functions that modify defaults in a constructor."],
      ["Zero value", "What a Go field holds when the JSON leaves it out."],
      ["Fail fast", "Invalid config stops the process at startup, before it serves anything."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "Precedence, lowest to highest: defaults, then the <code>.env</code> file, then ___.", answers: ["environment variables", "env vars", "env", "environment", "env variables"], ref: "ch05-s1" },
      { text: "A missing <code>.env</code> file is a ___, not a failure.", answers: ["warning"], ref: "ch05-s1" },
      { text: "POSIX <code>sh</code> already gives the variable ___ a meaning, so the app uses <code>APP_ENV</code>.", answers: ["ENV"], ref: "ch05-s4" },
      { text: "<code>PORT</code> must be between 1024 and ___.", answers: ["65535"], ref: "ch05-s1" },
      { text: "<code>omitempty</code> only affects ___, not decoding.", answers: ["encoding", "marshaling", "marshalling"], ref: "ch05-s7" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "Why does a config test that sets <code>APP_ENV=dev</code> prove little?", options: ["<code>dev</code> is invalid", "<code>dev</code> is the default, so the test passes even if the variable is never read", "Tests can't set env vars", "viper ignores lowercase values"], answer: 1, ref: "ch05-s2" },
      { q: "gin's release mode changes:", options: ["Request throughput significantly", "gin's own debug output, such as the route listing and warnings", "Middleware order", "The JSON encoder"], answer: 1, ref: "ch05-s6" },
      { q: "Where does &ldquo;prod means release mode&rdquo; live?", options: ["<code>internal/api</code>", "<code>internal/config</code>", "<code>main</code>", "The Dockerfile"], answer: 2, ref: "ch05-s5" },
      { q: "Rejecting <code>Timestamp</code> 0 costs you:", options: ["The ability to store negative values", "A device that genuinely reports epoch 0 is rejected", "Any timestamp before 2000", "Nothing at all"], answer: 1, ref: "ch05-s7" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "Describe the config-seam bug in three sentences: what the seam was, what it hid, and what changed.", refs: ["ch05-s2"] },
      { q: "Give two ways to make accidentally logging the config struct safe.", refs: ["ch05-s3"] },
      { q: "Calling <code>gin.SetMode</code> inside <code>NewServer</code>: what does that hide, and where should the cost be written down?", refs: ["ch05-s6"] },
    ]},
    { type: "interview", title: "Interview questions", items: [
      { q: "Where does your service get its configuration, and in what precedence order? What happens if the config file is missing?", refs: ["ch05-s1"] },
      { q: "Your first config test passed while production ignored every env var. Why, and what did you change so the test couldn't miss it again?", refs: ["ch05-s2"] },
      { q: "What's the risk of logging the parsed config at startup?", refs: ["ch05-s3"] },
      { q: "Why is the environment baked into your image as <code>prod</code>, and what does 12-factor say about that? Why is it <code>APP_ENV</code> and not <code>ENV</code>?", refs: ["ch05-s4"] },
      { q: "Where does &ldquo;prod means release mode&rdquo; live, and why there rather than in the HTTP layer?", refs: ["ch05-s5"] },
      { q: "A timestamp of <code>0</code> is now rejected. Why is Go's zero value a validation problem at a JSON boundary? What other options were there (pointer fields, <code>omitempty</code>, a presence check), and what do you lose by treating zero as &ldquo;missing&rdquo;?", refs: ["ch05-s7"] },
      { q: "Your library only offers a global setting (gin's mode). Where did you put the call that changes it, and what did that choice hide or expose? What does gin's release mode actually change?", refs: ["ch05-s6"] },
    ]},
  ],
},

// =================================================================================================
{
  id: "hw06", chapter: "ch06", title: "Metrics, Logs, and Cardinality",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch06-s1", pairs: [
      ["Series", "One metric name plus one exact set of label values."],
      ["Cardinality", "The number of series."],
      ["Stream (Loki)", "All log lines that share one unique label set."],
      ["RED", "Rate, errors, duration: for request-driven services."],
      ["USE", "Utilization, saturation, errors: for resources with capacity."],
      ["Histogram", "Counts per latency bucket, which can be summed across pods."],
      ["Summary", "Quantiles computed inside each process; can't be aggregated."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "Label HTTP metrics by <code>c.___()</code>, the route template, never the raw URL.", answers: ["FullPath"], ref: "ch06-s2" },
      { text: "Compute a quantile from histogram buckets at query time with ___.", answers: ["histogram_quantile"], ref: "ch06-s5" },
      { text: "A common convention for an application's metrics port is ___.", answers: ["2112"], ref: "ch06-s7" },
      { text: "Avoid port ___ for <code>/metrics</code>: it's the Prometheus server's own default.", answers: ["9090"], ref: "ch06-s7" },
      { text: "For the queue, the saturation signal is queue ___ relative to capacity.", answers: ["depth"], ref: "ch06-s4" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "3 routes &times; 6 status codes &times; 4 pods is how many series?", options: ["13", "72", "24", "4"], answer: 1, ref: "ch06-s1" },
      { q: "Pod A: 1,000 requests, p99 10 ms. Pod B: 10 requests, p99 2 s. Averaging the p99s gives about 1 s. What's wrong?", options: ["Nothing; it's a fair summary", "Percentiles aren't linear, so the average describes nobody's experience", "You should take the median instead", "p99 needs at least 100 requests"], answer: 1, ref: "ch06-s5" },
      { q: "In a Loki query, the label index is used to:", options: ["Search the text of every line", "Pick streams, whose chunks are then scanned", "Compute quantiles", "Deduplicate lines"], answer: 1, ref: "ch06-s3" },
      { q: "When should shedding page someone?", options: ["On any shed request", "When it's sustained and burning error budget in both a long and a short window", "Whenever queue depth is above zero", "Never; shedding is by design"], answer: 1, ref: "ch06-s6" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "If <code>device_id</code> can't be a label, name three places it can live.", refs: ["ch06-s2"] },
      { q: "Map USE's three signals onto the ring, one sentence each.", refs: ["ch06-s4"] },
      { q: "Compare Loki with a full-text index on ingest cost, query speed, and how much label choice matters.", refs: ["ch06-s3"] },
    ]},
    { type: "interview", title: "Interview questions", items: [
      { q: "Why must <code>device_id</code> never be a Prometheus label? What exactly grows, and what does it cost? Where <em>can</em> it live instead?", refs: ["ch06-s2"] },
      { q: "Loki indexes labels but not log content. What does that buy compared with a full-text index, and how does it change which labels you choose?", refs: ["ch06-s3"] },
      { q: "RED vs. USE: which applies to your HTTP handlers, and which to your queue? What's the saturation signal for the ring?", refs: ["ch06-s4"] },
      { q: "Histogram or summary for request latency? Why can't you average p99s across pods?", refs: ["ch06-s5"] },
      { q: "Your queue sheds with <code>503</code>. Which metric and alert tells you that's happening, and at what rate should someone be woken up?", refs: ["ch06-s6"] },
      { q: "Why serve <code>/metrics</code> on a separate port?", refs: ["ch06-s7"] },
    ]},
  ],
},

// =================================================================================================
{
  id: "hw07", chapter: "ch07", title: "Video I: Codecs, GOPs, and Containers",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch07-s2", pairs: [
      ["I-frame", "Decodable on its own; the largest frame type."],
      ["P-frame", "Predicted from earlier frames."],
      ["B-frame", "Predicted from earlier and later frames, which adds reordering delay."],
      ["GOP", "The run of frames from one keyframe to the next."],
      ["IDR", "An H.264 keyframe that also tells the decoder to discard earlier references."],
      ["Quantization", "The lossy step, where bitrate is traded for quality."],
      ["Transmux", "Rewrap compressed frames in another container without re-encoding."],
      ["Transcode", "Decode to raw frames and encode again."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "Raw 1080p at 30 fps with 4:2:0 chroma is about ___ Mbit/s.", answers: ["750", "746"], ref: "ch07-s1" },
      { text: "A 2-second GOP at 30 fps is ___ frames.", answers: ["60"], ref: "ch07-s2" },
      { text: "MPEG-TS packets are ___ bytes long.", answers: ["188"], ref: "ch07-s4" },
      { text: "Kinesis Video Streams stores media as ___ fragments.", answers: ["MKV", "Matroska"], ref: "ch07-s4" },
      { text: "The ffmpeg flag that transmuxes instead of transcoding is <code>-c ___</code>.", answers: ["copy"], ref: "ch07-s5" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "Why is plain MP4 risky as a recording format?", options: ["It can't hold H.264", "Its index is often written at the end, so a crash leaves the file unplayable", "It can't be streamed over HTTP", "It only supports one audio track"], answer: 1, ref: "ch07-s4" },
      { q: "Real-time WebRTC encoders avoid B-frames because:", options: ["They're larger than P-frames", "They need a future frame first, which adds latency", "Browsers can't decode them", "They break SRTP"], answer: 1, ref: "ch07-s2" },
      { q: "For a battery camera, the codec decision usually starts with:", options: ["Licensing cost", "What the camera's chip can encode in hardware", "Whichever codec is newest", "What the CDN prefers"], answer: 1, ref: "ch07-s6" },
      { q: "Which makes a new viewer wait longer to see a picture?", options: ["A short GOP", "A long GOP", "Using fMP4 instead of MPEG-TS", "Using Opus for audio"], answer: 1, ref: "ch07-s3" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "Name three tactics that let a fleet avoid transcoding on the server.", refs: ["ch07-s5"] },
      { q: "Why must HLS segment boundaries line up with keyframes, and how does that clash with WebRTC's keyframe habits?", refs: ["ch07-s3"] },
      { q: "How could switching the fleet to HEVC turn a bandwidth saving into a compute cost?", refs: ["ch07-s6"] },
    ]},
    { type: "interview", title: "Interview questions", items: [
      { q: "What is a GOP, and how does keyframe interval affect join time, seek, bandwidth, and recovery after loss?", refs: ["ch07-s2", "ch07-s3"] },
      { q: "Codec vs. container: H.264 in MPEG-TS, fMP4, MKV, IVF. Which container does each system in the stack use, and why does that matter for recording?", refs: ["ch07-s4"] },
      { q: "Transcode vs. transmux: what does each cost per stream, and how would you design a fleet so the server rarely transcodes?", refs: ["ch07-s5"] },
      { q: "H.264 vs. VP8/VP9 vs. AV1 vs. HEVC for a battery-powered camera: hardware encode, licensing, browser support.", refs: ["ch07-s6"] },
    ]},
  ],
},

// =================================================================================================
{
  id: "hw08", chapter: "ch08", title: "Video II: WebRTC and Real-Time Transport",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch08-s1", pairs: [
      ["STUN", "Tells a client what its public address looks like from outside."],
      ["TURN", "A server that relays every media packet between peers."],
      ["ICE", "Gathers candidate addresses and tests pairs to find one that works."],
      ["SDP", "Describes codecs, ICE credentials, and the DTLS certificate fingerprint."],
      ["DTLS-SRTP", "Keys from the DTLS handshake encrypt the RTP media."],
      ["SSRC", "A random identifier for one RTP stream."],
      ["Jitter buffer", "Holds packets briefly to reorder and smooth them before decoding."],
      ["SFU", "Forwards packets to each viewer without decoding them."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "WHIP, a standard for WebRTC signaling over HTTP, is RFC ___.", answers: ["9725"], ref: "ch08-s1" },
      { text: "In TCP, ___ blocking means one lost packet delays everything behind it.", answers: ["head-of-line", "head of line", "hol"], ref: "ch08-s3" },
      { text: "A receiver that can't decode asks for a keyframe with a PLI or an ___.", answers: ["FIR"], ref: "ch08-s5" },
      { text: "On a constrained uplink, <code>outbound-rtp.qualityLimitationReason</code> becomes ___.", answers: ["bandwidth"], ref: "ch08-s6" },
      { text: "GCC's delay-based estimate relies on transport-wide congestion control feedback, abbreviated ___.", answers: ["TWCC"], ref: "ch08-s6" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "Which part of session setup does the WebRTC standard leave to you?", options: ["The DTLS handshake", "ICE connectivity checks", "The signaling transport", "SRTP encryption"], answer: 2, ref: "ch08-s1" },
      { q: "Video packets are arriving, <code>framesDecoded</code> is flat, and <code>pliCount</code> is rising. Most likely:", options: ["The network is down", "The decoder is waiting for a keyframe", "Audio is muted", "TURN is in use"], answer: 1, ref: "ch08-s7" },
      { q: "FEC is the wrong tool when:", options: ["Loss is light and random", "The link is congested", "Round-trip time is long", "The stream is audio"], answer: 1, ref: "ch08-s5" },
      { q: "Two family members watch a camera over P2P. What changes?", options: ["Nothing", "The camera's uplink doubles", "Latency halves", "The SFU's load doubles"], answer: 1, ref: "ch08-s8" },
      { q: "STUN alone fails with:", options: ["A typical home NAT", "Symmetric NAT", "IPv6", "Wi-Fi"], answer: 1, ref: "ch08-s2" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "List the stages of glass-to-glass latency in order, with a rough figure for each.", refs: ["ch08-s9"] },
      { q: "Explain the escalation from NACK to FEC to PLI: when each applies.", refs: ["ch08-s5"] },
      { q: "Explain how GCC is backpressure in Chapter 4's sense, and which queue it's protecting.", refs: ["ch08-s6", "ch04-s6"] },
    ]},
    { type: "interview", title: "Interview questions", items: [
      { q: "Where does glass-to-glass latency come from? Name each stage, from capture and encode through the network and jitter buffer to decode and render, and which one you'd attack first.", refs: ["ch08-s9"] },
      { q: "UDP vs. TCP for real-time media: why does WebRTC choose loss over delay? When does it fall back to TCP or TURN-over-TLS, and what happens to quality?", refs: ["ch08-s3", "ch08-s2"] },
      { q: "A viewer's video freezes but audio keeps playing. Walk through the likely causes (lost keyframe, PLI/FIR, decoder stall) and what stats would tell them apart.", refs: ["ch08-s7"] },
      { q: "Walk through a WebRTC session setup: signaling, SDP offer/answer, ICE candidates, STUN, TURN, DTLS, SRTP. Which parts does the standard leave to you?", refs: ["ch08-s1"] },
      { q: "What fraction of sessions need TURN, and what does TURN cost you (bandwidth, latency, servers)? How would you find out?", refs: ["ch08-s2"] },
      { q: "What do NACK, FEC, and PLI each do, and when is each the wrong tool?", refs: ["ch08-s5"] },
      { q: "How does WebRTC congestion control react to a bad uplink from a home camera? What would you see in <code>getStats()</code>?", refs: ["ch08-s6"] },
      { q: "P2P vs. SFU vs. MCU: when does a home-camera product want each? What changes when two family members watch at once?", refs: ["ch08-s8"] },
    ]},
  ],
},

// =================================================================================================
{
  id: "hw09", chapter: "ch09", title: "Video III: Delivery, Media Servers, and Scale",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch09-s1", pairs: [
      ["Segment", "A few seconds of video, published as one file."],
      ["Playlist", "The <code>.m3u8</code> file listing available segments."],
      ["Partial segment", "An LL-HLS part, a few hundred milliseconds long."],
      ["Blocking playlist reload", "The server holds the request until the next part exists."],
      ["Media server", "Any protocol in, any protocol out, mostly by transmuxing."],
      ["Pre-roll", "Seconds of video kept on the camera so a clip includes the moment before the event."],
      ["Cumulative counter", "A value that only grows, so a lost sample costs only resolution."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "In KVS, <code>DataRetentionInHours</code> = ___ means no storage, pass-through only.", answers: ["0", "zero"], ref: "ch09-s5" },
      { text: "HLS players typically buffer about ___ segments before starting playback.", answers: ["3", "three"], ref: "ch09-s1" },
      { text: "Low-Latency HLS typically gets latency down to roughly 2&ndash;___ seconds.", answers: ["5", "five"], ref: "ch09-s2" },
      { text: "Prometheus's ___() function turns a counter into a per-second rate and handles resets.", answers: ["rate"], ref: "ch09-s6" },
      { text: "With &sect;8's assumptions, 24/7 recording stores ___ times as much as event-triggered recording.", answers: ["72"], ref: "ch09-s8" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "Why can a CDN cache HLS?", options: ["HLS is encrypted per viewer", "Every viewer requests the same files over plain HTTP", "HLS uses UDP", "Segments are under 1 MB"], answer: 1, ref: "ch09-s1" },
      { q: "A camera's upload is slower than its encoder. What should give first?", options: ["Keyframes", "Audio", "The encoder's bitrate", "The whole stream"], answer: 2, ref: "ch09-s7" },
      { q: "When you ship cumulative counters, a lost sample costs:", options: ["That interval's data, permanently", "Only resolution", "The whole session's data", "Nothing at all, ever"], answer: 1, ref: "ch09-s6" },
      { q: "Which of these does KVS <em>not</em> do for you?", options: ["Durable, time-indexed storage", "HLS playback", "Know that a doorbell event happened", "STUN/TURN for KVS WebRTC"], answer: 2, ref: "ch09-s5" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "Redo &sect;8's math with your own assumptions: 15 events a day of 20 s each, at 1.5 Mbit/s. Storage per camera per day, fleet per day, and 30-day total for 1M cameras.", refs: ["ch09-s8"] },
      { q: "What does the user see when the app switches from live to recorded, and what control would you add?", refs: ["ch09-s3"] },
      { q: "Explain how the camera's pre-roll buffer is the same structure as Chapter 1's ring.", refs: ["ch09-s5", "ch01-s3"] },
    ]},
    { type: "interview", title: "Interview questions", items: [
      { q: "WebRTC vs. HLS for live view: what latency does each give, why, and what does each cost to serve to 10,000 viewers? Why can a CDN cache one and not the other?", refs: ["ch09-s1"] },
      { q: "Live vs. recorded: why might the same camera use WebRTC to viewers and HLS (from KVS) for playback? Where does the switch happen in the UI, and what does the user see at the boundary?", refs: ["ch09-s3"] },
      { q: "What is LL-HLS, and what does it trade to get from ~10 s to ~2&ndash;3 s? Why isn't it WebRTC?", refs: ["ch09-s2"] },
      { q: "What does a media server like MistServer actually do? Why put one between a camera and viewers instead of going P2P?", refs: ["ch09-s4", "ch08-s8"] },
      { q: "What does Kinesis Video Streams handle for you (ingest, durable time-indexed storage, retention, HLS/DASH/clip playback, WebRTC signaling + TURN), and what's left for you to build?", refs: ["ch09-s5"] },
      { q: "How would you pick a KVS retention period? What does retention cost, and how does event-triggered recording change the bill?", refs: ["ch09-s5", "ch09-s8"] },
      { q: "Build vs. buy: Mist (self-hosted) for live and KVS (managed) for recording. What does each choice trade, and what would make you consolidate?", refs: ["ch09-s5"] },
      { q: "Which stream-quality metrics would you put on a dashboard, and which would wake someone up? Why does <code>device_id</code> (or SSRC) as a label break Prometheus?", refs: ["ch09-s6", "ch06-s2"] },
      { q: "<code>getStats()</code> counters are cumulative. Why ship counters rather than rates, and what does that do to your tolerance for losing telemetry?", refs: ["ch09-s6"] },
      { q: "A camera's upload is slower than its encoder. Where's the queue, what's the backpressure signal, and what should drop first? <em>(the M2 question, wearing a video costume)</em>", refs: ["ch09-s7"] },
      { q: "Design live view and event recording for 1M home cameras: capacity math for ingest bandwidth, storage per day, and concurrent viewers.", refs: ["ch09-s8"] },
    ]},
  ],
},

// =================================================================================================
{
  id: "hw10", chapter: "ch10", title: "The Road Ahead: Delivery to DORA",
  parts: [
    { type: "match", title: "Match each term to its meaning", ref: "ch10-s1", pairs: [
      ["At-least-once", "Retry until acknowledged; duplicates are possible, loss isn't."],
      ["Effectively once", "At-least-once delivery plus idempotent processing."],
      ["Visibility timeout", "How long a received SQS message stays hidden from other consumers."],
      ["Dead-letter queue", "Where a message goes after too many failed receives."],
      ["Two-phase commit", "A coordinator asks every participant to prepare, then to commit."],
      ["Consumer lag", "The log's latest offset minus the group's committed offset."],
      ["Linearizable", "Each operation appears to take effect at one instant, consistent with real time."],
    ]},
    { type: "blank", title: "Fill in the blank", items: [
      { text: "SQS FIFO queues deduplicate within a ___-minute window.", answers: ["5", "five"], ref: "ch10-s1" },
      { text: "A consumer extends its hold on a message with <code>___</code>.", answers: ["ChangeMessageVisibility"], ref: "ch10-s1" },
      { text: "A consumer group's parallelism is capped at the number of ___.", answers: ["partitions"], ref: "ch10-s3" },
      { text: "The two DORA keys that need incident data are change failure rate and time to ___ service.", answers: ["restore"], ref: "ch10-s5" },
      { text: "Moving a message to a dead-letter queue happens after <code>___</code> receives.", answers: ["maxReceiveCount"], ref: "ch10-s1" },
    ]},
    { type: "mc", title: "Multiple choice", items: [
      { q: "Why not use two-phase commit across Postgres and SQS?", options: ["It's too fast to observe", "SQS doesn't take part as an XA resource, and 2PC blocks if the coordinator fails after prepare", "Postgres doesn't support transactions", "It requires Kafka"], answer: 1, ref: "ch10-s2" },
      { q: "The outbox relay publishes events:", options: ["Exactly once", "At least once", "At most once", "In a single distributed transaction"], answer: 1, ref: "ch10-s2" },
      { q: "The caveat on calling a Redis counter linearizable:", options: ["Redis runs commands on many threads", "Replication is asynchronous, so a failover can lose recent increments", "INCR isn't atomic", "Lua scripts run concurrently"], answer: 1, ref: "ch10-s4" },
      { q: "With eager rebalancing:", options: ["Only the partitions that move are paused", "Every consumer gives up all its partitions while the group reorganizes", "Nothing pauses", "Offsets are reset to zero"], answer: 1, ref: "ch10-s3" },
    ]},
    { type: "short", title: "Short answer", items: [
      { q: "Pick a definition of &ldquo;successful deploy&rdquo; and defend it in three sentences.", refs: ["ch10-s5"] },
      { q: "For each of M3&ndash;M6, name what plays the role of Chapter 1's queue.", refs: ["ch10-s1", "ch10-s2", "ch10-s3", "ch10-s4", "ch01-s8"] },
      { q: "Why does your sort-on-read decision mean you don't need SQS FIFO queues?", refs: ["ch10-s1", "ch01-s5"] },
    ]},
    { type: "interview", title: "Interview questions (seeded per milestone)", items: [
      { q: "<strong>M3:</strong> At-least-once vs. exactly-once: what's the difference, and which one can you actually build?", refs: ["ch10-s1"] },
      { q: "<strong>M3:</strong> Why must consumers be idempotent?", refs: ["ch10-s1"] },
      { q: "<strong>M3:</strong> What does a visibility timeout buy, and what does it cost?", refs: ["ch10-s1"] },
      { q: "<strong>M4:</strong> Why not two-phase commit?", refs: ["ch10-s2"] },
      { q: "<strong>M4:</strong> What does the outbox cost you in exchange for atomicity?", refs: ["ch10-s2"] },
      { q: "<strong>M5:</strong> Partition ordering vs. throughput: what's the trade?", refs: ["ch10-s3"] },
      { q: "<strong>M5:</strong> What does consumer lag measure?", refs: ["ch10-s3"] },
      { q: "<strong>M5:</strong> Why is rebalancing the hard part?", refs: ["ch10-s3"] },
      { q: "<strong>M6:</strong> What does linearizable mean for rate-limit state?", refs: ["ch10-s4"] },
      { q: "<strong>M6:</strong> Why is coordination expensive?", refs: ["ch10-s4"] },
      { q: "<strong>M7:</strong> What are the four keys?", refs: ["ch10-s5"] },
      { q: "<strong>M7:</strong> What counts as a successful deploy?", refs: ["ch10-s5"] },
      { q: "<strong>M7:</strong> Which two of the four keys need incident data?", refs: ["ch10-s5"] },
    ]},
  ],
},

];
