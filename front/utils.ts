export function sleep(n: number, signal?: AbortSignal) {
  using stack = new DisposableStack();
  const p = Promise.withResolvers<void>()

  stack.adopt(setTimeout(() => {
    p.resolve();
  }, n), timeout => clearTimeout(timeout));

  if(signal) {
    const handler = () => {
      p.reject(new Error("Aborted"));
    }
    signal.addEventListener("abort", handler);
    stack.defer(() => {
      signal.removeEventListener("abort", handler);
    })
  }

  return p.promise
}

export function spawn(cb: (stack: DisposableStack) => PromiseLike<void>, handleError?: (e: any) => void) {
  const stack = new DisposableStack();
  try {
    const res = cb(stack);
    if(res instanceof Promise) {
      res.catch(e => {
        stack.dispose();
        if (handleError) {
          handleError(e);
        } else {
          throw e;
        }
      })
    }
  } catch (e) {
    stack.dispose();
    if (handleError) {
      handleError(e);
    } else {
      throw e;
    }
  }

  return stack;
}