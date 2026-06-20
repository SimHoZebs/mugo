# Database Transaction Evolution: From TxDatabase to Context-Aware Repositories

This document tracks the architectural evolution of the LazyFood/Mugo database layer, exploring the trade-offs between strict type-safety and developer velocity.

---

### Q1: What was the "TxDatabase" and why did we have it?
**Initial State:** We used two distinct types: `Database` (bound to a connection pool) and `TxDatabase` (bound to a specific `pgx.Tx`). Both implemented the `DB` interface.

**The Reasoning:** This provided **Strict Type Safety**. If a function required a transaction, you had to pass it a `*TxDatabase`. The compiler prevented you from accidentally running a non-transactional query inside a transactional block.

---

### Q2: What was the "Maintenance Burden" of this approach?
**The Problem:** Because `sqlc` generates a `Queries` struct that must be explicitly "re-bound" to a transaction, adding a new repository (e.g., `UserRepository`) required updating four separate locations:
1.  **The `DB` Interface:** Add `Users() UserRepository`.
2.  **The `Database` Struct:** Add the field and the getter.
3.  **The `TxDatabase` Struct:** Add the field and the getter (The "Mirroring" problem).
4.  **The `WithTx` Factory:** Update the logic to initialize the new repository with the transaction.

This "Double Maintenance" created friction for every new feature, as the specialized types were essentially identical in capability but different in internal state.

---

### Q3: Why is the `q(ctx)` helper necessary in the new model?
**The Shift:** We moved the transaction state from the *Struct* into the *Context*.

**The Fundamental Logic:**
*   **`sqlc` Requirement:** `Queries` objects are immutable regarding their database connection. To switch from a pool to a transaction, you *must* call `queries.WithTx(tx)`.
*   **The Helper:** `q(ctx)` acts as a "Traffic Controller." It checks the `context.Context` for a transaction "receipt" (placed there by `WithTx`).
    *   **Found?** It returns a temporary `Queries` object bound to that transaction.
    *   **Not Found?** It returns the default `Queries` object bound to the pool.

This allows the **Business Logic** to be "Transaction Blind"—it calls the same method on the same object, and the repository handles the routing internally.

---

### Q4: If the Context is the source of truth, why does `WithTx` still pass a `db` argument?
**The Redundancy:** Technically, since `q(ctx)` looks at the context, you could use the original global `db` object inside a `WithTx` block and it would still be transactional.

**The Trade-off:**
1.  **Interface Consistency:** It keeps the `DB` interface standardized.
2.  **Explicit Intent:** Passing `txDB` into the callback signals to the developer that "this block is transactional."
3.  **Future-Proofing:** It allows us to wrap the database in a logging or restricted interface later without changing every call site.

---

### Q5: Does specialization truly avoid overlap? 
**The Intuition:** It seems like having two types (`Database` and `TxDatabase`) would naturally separate concerns.

**The Reality in Go + sqlc:** In our case, the "Specialization" wasn't changing the *capabilities* of the database, only its *internal state* (Pool vs. Transaction). Because we wanted our business logic (routes/services) to be "Transaction Blind," we forced both types to implement the exact same interface.

**The Overlap:** This created "Fake Specialization." We had two different types that:
1.  Had the same methods.
2.  Provided access to the same repositories.
3.  Returned the same models.

The only difference was one field inside the struct. This is where the overlap occurred: we were maintaining two identical "maps" of the database structure in code, which led to the "Mirroring" problem.

---

### Q6: What are the true costs of specialization vs. unification?
**The Cost of Specialization (Strict Type Safety):**
*   **Maintenance:** High "friction." Every new repository requires 4+ file changes.
*   **Initialization:** Every `WithTx` call must re-allocate and re-initialize *every* repository in the system (e.g., `NewUserRepository(tx)`, `NewMealRepository(tx)`), even if the transaction only uses one of them. While memory allocation is cheap, this is wasted work.
*   **Cognitive Load:** Developers must constantly decide which object to pass around.

**The Cost of Unification (Context-Aware):**
*   **Type Safety:** The compiler no longer "forces" a transaction. If a developer passes `context.Background()` instead of the context from `WithTx`, the query silently runs outside the transaction.
*   **Implicit Behavior:** The "magic" happens inside `q(ctx)`, which can be less obvious to a new developer than a dedicated `TxDatabase` type.

---

### Summary: The Unified Model vs. Specialized Model

| Feature | Specialized (`TxDatabase`) | Unified (`q(ctx)`) |
| :--- | :--- | :--- |
| **Safety** | Compiler-enforced transactions. | Trust-based (relies on `ctx` propagation). |
| **Maintenance** | High (Update 4+ files per repo). | Low (Update 1 file per repo). |
| **Cleanliness** | Verbose, duplicate code. | Minimal, idiomatic Go. |
| **Performance** | Constant re-initialization of all repos. | Lazy re-binding via `WithTx(tx)` only when needed. |

**Conclusion:** We chose the **Unified Model** because the "specialization" we had was not providing unique functionality—it was only providing a safety check at the cost of significant structural duplication. The unified model is more "Lazy" (only re-binding `sqlc.Queries` when a method is actually called) and significantly faster to extend.
