# Functional vs OOP: Runner Package Evolution

Explored whether a purely functional approach (stateless functions, no objects) could replace the OOP interface+struct pattern for the ADK runner wrapper. Started with OOP, went functional, came back to a struct — and learned where the line is.

---

### Q1: Why did the runner originally have an interface + struct?

The `AgentRunner` interface bundled two related operations (`Run`, `CreateSession`) that share state (`*runner.Runner`, `session.Service`, `appName`). Routes depended on the interface, mocks implemented it for testing. Standard OOP strategy pattern.

---

### Q2: What happened when we removed the interface and struct?

Replaced with stateless functions: `NewRunner`, `Run(ctx, *runner.Runner, ...)`, `CreateSession(ctx, session.Service, appName, ...)`. State passed as arguments. No objects, no methods.

**What worked:** The runner package became pure. Easy to read, easy to reason about. Each function did one thing.

**What broke:** The route signature. To call `runner.Run` and `runner.CreateSession` in a handler, the route needed `*runner.Runner` + `session.Service` + `appName` as three separate parameters. These three values are related — they're all part of "how to run this agent" — but they traveled separately. Every route that used a runner had the same ugly three-param signature.

---

### Q3: Why not use closures to bind the state?

A closure capturing `*runner.Runner` + `session.Service` + `appName` and returning a function would clean up the route signature. But a closure is an object with no structural enforcement. At 2-3 related operations, the lack of structure hurts: no compile-time check that all operations are present, no named method set, adding an operation means adding a parameter at every call site.

---

### Q4: Why not pass `*runner.Runner` directly and let the handler call `runner.Run`?

`*runner.Runner` is a concrete Google type with unexported fields. It absorbs `appName` at construction but doesn't expose it. `CreateSession` needs `appName` to call the ADK session service. So `appName` must be passed separately — the runner has it but won't share it. This is a structural constraint of the dependency, not a design choice.

---

### Q5: What's the actual tradeoff line?

| Operations sharing state | Better approach |
|---|---|
| 1 | Stateless function — no bundling needed |
| 2-3 | Struct with methods — bundling earns its keep |
| 5+ (like DB repositories) | Interface with named method set — structure is essential |

The runner is at 2-3 operations. The struct is the right call. The DB layer is at 5+ — the interface is essential there.

---

### Q6: Is trying functional approach wasted effort if we end up back at OOP?

No. The exercise taught *why* the struct exists, not just *that* it exists. Before, the interface+struct was there because that's how it was vibecoded. Now it's there because we felt the friction of not having it — ugly route signatures, redundant parameter passing, no compile-time enforcement. That understanding transfers: when designing new packages, the question isn't "OOP or FP?" but "how many operations share state, and does bundling pay off?"

---

### Q7: What did we land on?

A struct `AgentRunner` with `Runner`, `SessionService`, and `AppName` fields, plus `Run` and `CreateSession` methods that delegate to stateless functions (`runner.Run`, `runner.CreateSession`). The stateless functions still exist for testing and direct use. The struct is a convenience wrapper that bundles related state — not an OOP hierarchy, just a bag of values with methods that forward to pure functions.

```
runner.NewAgentRunner("meal_orchestrator", agent, ss) → *AgentRunner
runner.Run(ctx, *runner.Runner, ...)                   → *RunResult  (stateless)
runner.CreateSession(ctx, ss, appName, ...)            → error        (stateless)

agentRunner.Run(ctx, uid, sid, text)                   → delegates to runner.Run
agentRunner.CreateSession(ctx, uid, sid)               → delegates to runner.CreateSession
```

The key difference from the original: no interface. The struct is concrete, routes take `*runner.AgentRunner` directly. The interface was only needed for mocking, and `CreateMeal` (the extracted service function) is now tested without mocks — it takes `*runner.RunResult` (a plain struct) as input.
